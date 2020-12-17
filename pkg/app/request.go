package app

import (
	"filfox_data/pkg/logging"
	"github.com/astaxie/beego/validation"
)

// MarkErrors logs error logs
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}

	return
}
