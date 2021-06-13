package app

import (
	"github.com/astaxie/beego/validation"

	"github.com/jamsa/gin-k8s/pkg/logging"
)

// MarkErrors logs error logs
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}

	return
}


// 分页查询参数
type Page struct{
	PageNum int
	PageSize int
}
