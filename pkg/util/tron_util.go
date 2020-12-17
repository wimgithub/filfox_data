package util

import (
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"golang.org/x/crypto/sha3"
	"math/big"
	"strings"
)

func GetMethodID(str string) []byte {
	transferFnSignature := []byte(str)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	return hash.Sum(nil)[:4]
}

// 处理合约获取余额参数
func ProcessBalanceOfParameter(addr string, methodID []byte,isSend bool,amount int64) (data []byte) {
	add, _ := common.DecodeCheck(addr)
	paddedAddress := common.LeftPadBytes(add[1:], 32)
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	if isSend {
		amountBig := new(big.Int).SetInt64(amount)
		paddedAmount := common.LeftPadBytes(amountBig.Bytes(), 32)
		data = append(data, paddedAmount...)
	}
	return
}

// 处理合约转账参数
func ProcessTransferParameter(methodID []byte, amount int64) (data []byte) {
	amountBig := new(big.Int).SetInt64(amount)
	paddedAmount := common.LeftPadBytes(amountBig.Bytes(), 32)
	data = append(data, methodID...)
	data = append(data, paddedAmount...)
	return
}

// 地址处理
func AddressDealWith(address string) string {
	if strings.Contains(address, "0x41") {
		return address
	}
	return strings.Replace(address, "0x", "0x41", -1)
}
