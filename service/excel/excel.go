package excel

import (
	model "filfox_data/models"
	"filfox_data/models/mysql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"sort"
	"time"
)

func (a Datas) Len() int           { return len(a) }
func (a Datas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Datas) Less(i, j int) bool { return a[i].Time < a[j].Time }

type Datas []*model.Data

var Rows = []string{"时间", "消息ID", "发送方", "接收方", "净收入", "类型"}

func GetExcel(c *gin.Context, begin, end, height int64, msg, to, t string) (string, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return "", err
	}
	first := sheet.AddRow()
	first.SetHeightCM(1) //设置每行的高度
	for _, r := range Rows {
		cell := first.AddCell()
		cell.Value = r
	}

	data, _ := mysql.SharedStore().GetFilFoxData(begin, end, height, msg, to, t)
	sort.Sort(Datas(data))
	for _, v := range data {
		row := sheet.AddRow()
		row.SetHeightCM(1)                                              //设置每行的高度
		cell := row.AddCell()                                           //cell 1
		cell.Value = time.Unix(v.Time, 0).Format("2006-01-02 15:04:05") //时间
		cell = row.AddCell()                                            //cell 2
		cell.Value = v.Message                                          //消息ID
		cell = row.AddCell()                                            //cell 3
		cell.Value = v.FilFrom                                          //发送方
		cell = row.AddCell()                                            //cell 4
		cell.Value = v.FilTo                                            //接收方
		cell = row.AddCell()                                            //cell 5
		cell.Value = v.Value                                            // 净收入
		cell = row.AddCell()                                            //cell 6
		cell.Value = v.Type                                             // 类型
	}
	name := fmt.Sprint("./excel/", time.Now(), ".xlsx")
	err = file.Save(name)
	if err != nil {
		return "", err
	}
	return name, nil
}
