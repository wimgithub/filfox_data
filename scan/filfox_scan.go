package scan

import (
	"encoding/json"
	model "filfox_data/models"
	"filfox_data/models/mysql"
	"filfox_data/pkg/http_util"
	"filfox_data/pkg/logging"
	"filfox_data/pkg/util"
	"fmt"
	"github.com/shopspring/decimal"
	"sort"
	"strconv"
	"time"
)

func (a Datas) Len() int           { return len(a) }
func (a Datas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Datas) Less(i, j int) bool { return a[i].Timestamp < a[j].Timestamp }

type Datas []*model.FilFoxResponse

type FilFoxScan struct {
	url      string
	total    int64
	ch       chan []*model.FilFoxResponse
	errPages *model.PageData
}

func NewFilFoxScan() *FilFoxScan {
	count, _ := mysql.SharedStore().GetFilFoxCount()
	return &FilFoxScan{
		url:   "https://filfox.info/api/v1/address/f3u5xnumgzr2h4ysnejnrket7boj3457vyh22s4wjnfhukefzgw5n6zi3kp5slufat3dpvag3eifcklb5vx2iq/transfers?pageSize=100&page=",
		total: count,
	}
}

func (f *FilFoxScan) Start() {
	go f.DataHandle()
}

func (f *FilFoxScan) AllDataHandle() {
	data, err := f.GetData(0)
	if err != nil {
		return
	}
	if data.TotalCount > f.total {
		count := data.TotalCount - f.total
		totalPage := (count / 100) + 1
		logging.Info("当前数据总量：", data.TotalCount, " 已加载总量：", f.total, " 还需加载总量：", count, " 需加载总页数: ", totalPage)
		for i := totalPage; i >= 0; i-- {
			getData, err := f.GetData(i)
			if err != nil || getData.Transfers == nil {
				logging.Error("第 ", i, " 页获取失败!")
				continue
			}
			sort.Sort(Datas(getData.Transfers))
			f.ResponseHandler(getData.Transfers)
			fmt.Println("第: ", i, " 页数据: ", len(getData.Transfers))
		}
		f.total += count
	}
}

func (f *FilFoxScan) DataHandle() {
	for {
		data, err := f.GetData(0)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		if data.TotalCount > f.total {
			var allData []*model.FilFoxResponse
			count := data.TotalCount - f.total
			totalPage := (count / 100) + 1
			logging.Info("当前数据总量：", data.TotalCount, " 已加载总量：", f.total, " 还需加载总量：", count, " 需加载总页数: ", totalPage)
			for i := totalPage; i >= 0; i-- {
				getData, err := f.GetData(i)
				if err != nil || getData.Transfers == nil {
					logging.Error("第 ", i, " 页获取失败!")
					f.errPages.Page = append(f.errPages.Page, i)
					continue
				}
				allData = append(allData, getData.Transfers...)
				fmt.Println("第: ", i, " 页数据: ", len(getData.Transfers))
				time.Sleep(1 * time.Second)
			}
			sort.Sort(Datas(allData))
			// 数据处理入库
			f.ResponseHandler(allData[int64(len(allData))-count:])
			// 更新 f.total
			f.total += count
		}
		time.Sleep(3 * time.Second)
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
		fmt.Println("time: ", v.Timestamp)
		if len(d) >= 300 {
			err := mysql.SharedStore().AddFilData(d)
			if err != nil {
				logging.Error("mysql insert error!", err)
			}
			fmt.Println(">300-入库： ", len(d))
			d = nil
		}
	}
	err := mysql.SharedStore().AddFilData(d)
	if err != nil {
		logging.Error("mysql insert error!", err)
	}
	fmt.Println("<300入库： ", len(d))
}

func (f *FilFoxScan) GetData(page int64) (data *model.Resp, err error) {
	url := f.url + fmt.Sprint(page)
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
