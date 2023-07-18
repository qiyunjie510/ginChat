package controller

import (
	"fmt"
	"ginchat/models"
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 页面
// @Summary 首页目前是登录
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "view/chat/head.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "index")
	if err != nil {
		fmt.Println(err)
	}
}

// GetRegister
// @Tags 页面
// @Summary 注册页面
// @Success 200 {string} welcome
// @Router /register [get]
func GetRegister(c *gin.Context) {
	ind, err := template.ParseFiles("view/user/register.html")
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	err = ind.Execute(c.Writer, "register")
	if err != nil {
		fmt.Println(err)
	}
}

func GetCommunity(c *gin.Context) {
	ind, err := template.ParseFiles("view/chat/createcom.html", "view/chat/head.html")
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	err = ind.Execute(c.Writer, "createcom")
	if err != nil {
		fmt.Println(err)
	}
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles(
		"view/chat/index.html",
		"view/chat/head.html",
		"view/chat/tabmenu.html",
		"view/chat/concat.html",
		"view/chat/group.html",
		"view/chat/profile.html",
		"view/chat/main.html",
		"view/chat/foot.html",
		"view/chat/createcom.html")
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")

	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	fmt.Println(user)
	err = ind.Execute(c.Writer, user)
	if err != nil {
		fmt.Println(err)
	}
}
