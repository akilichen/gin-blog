package app

import (
	"gin-blog/pkg/mylog"
	"github.com/astaxie/beego/validation"
)

func MakeErrors(errors []*validation.Error) {
	for _, err := range errors {
		mylog.Info(err.Key, err.Message)
	}
	return
}
