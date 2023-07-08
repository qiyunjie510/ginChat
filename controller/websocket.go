package controller

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Chat(c *gin.Context) {
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("socket connected!")
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("socket close!")
		}
	}(ws)
	var wg sync.WaitGroup
	wg.Add(2)
	// 接收客户端消息
	go func(ws *websocket.Conn, c *gin.Context) {
		for {
			_, p, err := ws.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
			tm := time.Now().Format("2006-01-02 15:04:05")
			m := fmt.Sprintf("[ws][%s]:%s", tm, string(p))
			err = utils.Publish(c, utils.PublishKey, m)
			if err != nil {
				fmt.Println(err)
				break
			}
		}
		wg.Done()
	}(ws, c)
	// 向客户端发消息
	go func(ws *websocket.Conn, c *gin.Context) {
		for {
			msg, err := utils.Subscribe(c, utils.PublishKey)
			if err != nil {
				fmt.Println(err)
				break
			}
			tm := time.Now().Format("2006-01-02 15:04:05")
			m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
			err = ws.WriteMessage(1, []byte(m))
			if err != nil {
				fmt.Println(err)
				break
			}
		}
		wg.Done()
	}(ws, c)
	wg.Wait()
}

func ChatV2(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	users := models.SearchFriend(c.PostForm("userId"))
	utils.RespOKList(c.Writer, users, len(users))
}