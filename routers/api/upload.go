package api

import (
	"gin-blog/pkg/exception"
	"gin-blog/pkg/mylog"
	"gin-blog/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadImage(c *gin.Context) {
	code := exception.SUCCESS
	data := make(map[string]interface{})

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		mylog.Warning(err)
		code = exception.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"data": data,
			"msg":  exception.GetMsg(code),
		})
	}

	if image != nil {
		imageName := upload.GetImageName(image.Filename)
		imageSavePath := upload.GetImagePath()
		imageFullPath := upload.GetImageFullPath()

		src := imageFullPath + imageName

		if !upload.CheckImageExt(image.Filename) || !upload.CheckImageSize(file) {
			code = exception.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(imageFullPath)
			if err != nil {
				mylog.Warning(err)
				code = exception.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				mylog.Warning(err)
				code = exception.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = imageSavePath + imageName
			}
		}
	} else {
		code = exception.INVALID_PARAMS
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  exception.GetMsg(code),
		"data": data,
	})
}
