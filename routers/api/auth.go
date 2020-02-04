package api

import (
	"gin-blog/models"
	"gin-blog/pkg/exception"
	"gin-blog/pkg/mylog"
	"gin-blog/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	data := make(map[string]interface{})
	code := exception.INVALID_PARAMS
	if ok {
		isExist := models.CheckAuth(username, password)
		if isExist {
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = exception.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token

				code = exception.SUCCESS
			}

		} else {
			code = exception.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			mylog.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  exception.GetMsg(code),
		"data": data,
	})
}
