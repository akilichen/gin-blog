package upload

import (
	"fmt"
	"gin-blog/pkg/file"
	"gin-blog/pkg/mylog"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// 图片在图片目录的地址
func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefix + "/" + GetImagePath() + name
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMd5(fileName)

	return fileName + ext
}

func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

//图片在整个项目中的地址
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToLower(ext) == allowExt {
			return true
		}
	}
	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		mylog.Warning(err)
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.GetWd error :%v\n", err)
	}
	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.isNotExistMkDir error %v\n", err)
	}

	permission := file.CheckPermission(src)
	if permission == true {
		return fmt.Errorf("file.CheckPermission permission denied src :%v\n", src)
	}
	return nil
}
