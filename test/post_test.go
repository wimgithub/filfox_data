package test

import (
	"encoding/json"
	model "filfox_data/models"
	"filfox_data/pkg/http_util"
	"filfox_data/pkg/util"
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"testing"
)

func TestGet(t *testing.T) {
	url := "https://filfox.info/api/v1/address/f3u5xnumgzr2h4ysnejnrket7boj3457vyh22s4wjnfhukefzgw5n6zi3kp5slufat3dpvag3eifcklb5vx2iq/transfers?pageSize=100&page=0"
	byte, err := http_util.Get(url)
	if err != nil {
		panic(err)
	}

	var data *model.Resp
	err = json.Unmarshal(byte, &data)
	if err != nil {
		panic(err)
	}
	for _, v := range data.Transfers {
		fmt.Println(v)
	}

}

func Test(t *testing.T) {
	var a string
	a = "-182957707651276349"
	parseInt, err := strconv.ParseInt(a, 0, 64)
	if err != nil {
		t.Error(err)
	}
	fil := util.ToFil(decimal.NewFromInt(parseInt))
	fmt.Println(parseInt)
	fmt.Println(fil.String())
}

func TestRange(t *testing.T)  {
	for i := 10; i > 0; i-- {
		fmt.Println(i-1)
	}
}
