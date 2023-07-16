package service

import (
	"fmt"
	"ginchat/utils"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Upload(c *gin.Context) {
	request := c.Request
	writer := c.Writer
	file, m, err := request.FormFile("file")
	if err != nil {
		utils.RespFail(writer, err.Error())
	}
	suffix := ".png"
	oFileName := m.Filename
	temp := strings.Split(oFileName, ".")
	if len(temp) > 1 {
		suffix = "." + temp[len(temp)-1]
	}

	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFile, er := os.Create("./asset/upload/" + fileName)
	if er != nil {
		utils.RespFail(writer, er.Error())
	}
	_, err = io.Copy(dstFile, file)
	if err != nil {
		utils.RespFail(writer, err.Error())
	}
	url := "/asset/upload/" + fileName
	utils.RespOK(writer, url, "发送图片成功")
}
