package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	FromId   int64  // 发送者
	TargetId int64  // 接收者
	Type     int    // 消息类型 ： 群聊、私聊、广播
	Media    int    // 消息类型 文字 图片 音频
	Content  string // 消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int // 其他数据统计
}

type MessageV2 struct {
	DstId   int64  `json:"dstid"`
	Cmd     int    `json:"cmd"`
	UserId  int64  `json:"userid"`
	Media   int    `json:"media"`
	Content string `json:"content"`
}

//
func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	// 获取参数，校验合法性
	query := request.URL.Query()
	id := query.Get("id")
	userId, _ := strconv.ParseInt(id, 10, 64)
	// token := query.Get("token")
	// targetId := query.Get("target_id")
	// Content := query.Get("Content")
	// msgType := query.Get("type")
	isValid := true // checkToken()
	conn, err := (&websocket.Upgrader{
		// token校验
		CheckOrigin: func(r *http.Request) bool {
			return isValid
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 获取连接
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	// 用户关系判断
	// userid和node绑定，加锁
	rwLocker.Lock()
	fmt.Printf("id=%d用户已经连接..\n", userId)
	clientMap[userId] = node
	rwLocker.Unlock()
	// 发送逻辑
	go sendProc(node)
	// 接受逻辑
	go receiveProc(node)

	sendMsg(userId, []byte("欢迎来到聊天室"))

}

// 发送消息给客户端
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Printf("从客户端接受到消息")
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 从客户端接受消息
func receiveProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws] :", string(data))
	}
}

var udpSendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpSendChan <- data
}

func init() {
	go udpSendProc()
	go udpReceiveProc()
}

// 从广播中拿到数据，然后发送数据到udp
func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpSendChan:
			fmt.Printf("udp 发送消息%s", string(data))
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据接收协程
func udpReceiveProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	for {
		var buf [512]byte
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("udp 接到消息%s\n", string(buf[0:n]))
		dispatch(buf[0:n])

	}
}

// 后端调度逻辑
func dispatch(data []byte) {
	msg := MessageV2{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("msg的目标是：%d", msg.DstId)
	switch msg.Cmd {
	case 1: // 私信
		sendMsg(msg.DstId, data)
	case 11: // 群发
		err := sendGroupMsg(data)
		if err != nil {
			fmt.Println("发送群消息失败")
		}
		// case 3: // 广播
		// 	sendAllMsg()
		// default
	}
}

func sendMsg(userId int64, data []byte) {
	rwLocker.Lock()
	node, ok := clientMap[userId]
	fmt.Printf("send to %d\n", userId)
	rwLocker.Unlock()
	if ok {
		fmt.Println("sendMsg done")
		node.DataQueue <- data
	}
}

func sendGroupMsg(data []byte) error {
	// 通过群，查到群里面的人，依次给每一个人发送消息
	message := MessageV2{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		return errors.New("解析消息失败")
	}
	groupId := message.DstId
	contactList, err := getGroupNumberList(uint(groupId))
	if err != nil {
		return errors.New("加载群成员列表失败")
	}
	for _, contact := range contactList {
		sendMsg(int64(contact.OwnerId), data)
	}
	return nil
}
