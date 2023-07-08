package controller

import (
	"fmt"
	"ginchat/common"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// List
// @Tags 用户模块
// @Summary 用户列表
// @Success 200 {string} json {"code","message"}
// @Router /api/user/list [get]
func List(c *gin.Context) {
	data := models.GetUserList()
	utils.Success(c, "获取列表成功", data)
}

// DeleteUser
// @Tags 用户模块
// @Summary 删除用户
// @param id query string false "id"
// @Success 200 {string} json {"code","message"}
// @Router /api/user/delete [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	utils.Success(c, "删除用户成功", nil)
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json {"code","message"}
// @Router /api/user/update [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		utils.Error(c, common.ERROR, "修改参数不符合规则！", nil)
		return
	}
	models.UpdateUser(user)
	utils.Success(c, "修改用户成功", user)
}

// Login
// @Tags 登录模块
// @Summary 登录
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @Success 200 {string} json {"code","message"}
// @Router /api/user/login [post]
func Login(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")

	user := models.FindUserByName(name)
	if user.Name == "" {
		utils.Error(c, common.ERROR, "用户不存在", nil)
		return
	}

	ok := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !ok {
		utils.Error(c, common.ERROR, "密码错误", nil)
		return
	}
	pwd := utils.MakePassword(password, user.Salt)
	user = models.FindUserByNameAndPwd(name, pwd)
	utils.Success(c, "登录成功", user)
}

// Register
// @Tags 登录模块
// @Summary 用户注册
// @param name formData string false "用户名"
// @param phone formData string false "手机号"
// @param password formData string false "密码"
// @param real_password formData string false "确认密码"
// @Success 200 {string} json {"code","message"}
// @Router /api/user/register [post]
func Register(c *gin.Context) {
	user := models.UserBasic{}

	name := c.PostForm("name")
	passward := c.PostForm("password")
	phone := c.PostForm("phone")
	realPassward := c.PostForm("real_password")
	if name == "" || passward == "" || realPassward == "" || phone == "" {
		utils.Error(c, common.ERROR, "用户名、密码、手机号不能为空", nil)
		return
	}

	if passward != realPassward {
		utils.Error(c, common.ERROR, "两次密码不一致！", nil)
		return
	}

	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(name)
	if data.Name != "" {
		utils.Error(c, common.ERROR, "用户名已经注册", nil)
		return
	}
	user.PassWord = utils.MakePassword(passward, salt)
	user.Salt = salt
	user.Phone = phone
	user.Name = name
	models.CreateUser(user)
	utils.Success(c, "新增用户成功", user)
}
