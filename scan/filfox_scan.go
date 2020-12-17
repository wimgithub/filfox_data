package scan

import (
	"encoding/json"
	model "filfox_data/models"
	"filfox_data/models/mysql"
	"filfox_data/pkg/gredis"
	"filfox_data/pkg/http_util"
	"filfox_data/pkg/logging"
	"filfox_data/pkg/util"
	"fmt"
	"github.com/panjf2000/ants"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

type FilFoxScan struct {
	url      string
	page     int
	total    int64
	ch       chan []*model.FilFoxResponse
	errPages *model.PageData
}

func NewFilFoxScan() *FilFoxScan {
	return &FilFoxScan{
		url: "https://filfox.info/api/v1/address/f3u5xnumgzr2h4ysnejnrket7boj3457vyh22s4wjnfhukefzgw5n6zi3kp5slufat3dpvag3eifcklb5vx2iq/transfers?pageSize=100&page=",
	}
}

func (f *FilFoxScan) Start() {
	go f.DataHandle()
	go f.runSnapshots()
	//f.AllDataHandle()
}

func (f *FilFoxScan) AllDataHandle() {
	//rand.Seed(time.Now().UnixNano())
	pool, _ := ants.NewPool(5)
	defer pool.Release()

	data, err := f.GetData(0)
	if err != nil {
		return
	}
	f.ch = make(chan []*model.FilFoxResponse, data.TotalCount)
	if data.TotalCount > f.total {
		count := data.TotalCount - f.total
		totalPage := (count / 100) + 1
		logging.Info("总量：", count, " 总页数: ", totalPage)
		for i := totalPage; i >= 0; i-- {
			_ = pool.Submit(f.submit(i))
		}
	}
	for {
		select {
		case data := <-f.ch:
			f.ResponseHandler(data)
			if len(f.ch) == 0 {
				fmt.Println("Channel OK! ")
				return
			}
		}
	}
}

func (f *FilFoxScan) submit(i int64) func() {
	return func() {
		fmt.Println("协程 ：", i, " 开始处理")
		//time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		data, err := f.GetData(i)
		if err != nil || data.Transfers == nil {
			logging.Error("数据获取失败: ", i)
			return
		}
		f.ch <- data.Transfers
		fmt.Println("协程 ：", i, " 处理完成")
	}
}

func (f *FilFoxScan) DataHandle() {
	data, err := f.GetData(0)
	if err != nil {
		return
	}
	if data.TotalCount > f.total {
		//var allData []*model.FilFoxResponse
		count := data.TotalCount - f.total
		totalPage := (count / 100) + 1
		for i := totalPage; i >= 0; i-- {
			getData, err := f.GetData(i)
			if err != nil || getData.Transfers == nil {
				logging.Error("第 ", i, " 页获取失败!")
				f.errPages.Page = append(f.errPages.Page, i)
				continue
			}
			f.ResponseHandler(getData.Transfers)
			//allData = append(allData, getData.Transfers...)
			fmt.Println("第: ", i, " 页数据: ", len(getData.Transfers))
			time.Sleep(1 * time.Second)
		}
		// 数据处理入库
		//f.ResponseHandler(allData[int64(len(allData))-count:])
		// 更新 f.total
	}
}

func (f *FilFoxScan) ResponseHandler(data []*model.FilFoxResponse) {
	var d []*model.Data
	for _, v := range data {
		parseInt, _ := strconv.ParseInt(v.Value, 0, 64)
		d = append(d, &model.Data{
			Time:    v.Timestamp,
			FilFrom: v.From,
			Height:  v.Height,
			Message: v.Message,
			FilTo:   v.To,
			Type:    f.typeHandler(v.Type),
			Value:   util.ToFil(decimal.NewFromInt(parseInt)).String(),
		})
	}
	err := mysql.SharedStore().AddFilData(d)
	if err != nil {
		logging.Error("数据入库失败!")
		return
	}
}

func (f *FilFoxScan) GetData(page int64) (data *model.Resp, err error) {
	url := f.url + fmt.Sprint(page)
	fmt.Println("url: ", url)
	byte, err := http_util.Get(url)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(byte, &data)
	return
}

func (f *FilFoxScan) typeHandler(t string) string {
	switch t {
	case "send":
		return "发出"
	case "burn-fee":
		return "销毁手续费"
	case "miner-fee":
		return "矿工手续费"
	default:
		return t
	}
}

func (f *FilFoxScan) runSnapshots() {
	for {
		select {
		case _, ok := <-gredis.RedisSnapshot:
			if !ok {
				marshal, _ := json.Marshal(f.errPages)
				logging.Info("redis json: ", string(marshal))
				err := gredis.SharedSnapshotStore().SetJson(gredis.ErrPages, marshal)
				if err != nil {
					logging.Error("err pages 备份失败: ", err)
				}
				return
			}
		}
	}
}
