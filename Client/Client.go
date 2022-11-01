/*
 * @Author: git config SnowFish && git config 3200401354@qq.com
 * @Date: 2022-10-28 11:24:43
 * @LastEditors: git config SnowFish && git 3200401354@qq.com
 * @LastEditTime: 2022-11-01 17:10:57
 * @FilePath: \IM_V3\Client\Client.go
 * @Description:
 *
 * Copyright (c) 2022 by snow-fish 3200401354@qq.com, All Rights Reserved.
 */
package main

import (
	f "demo/theme/fonts"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

/*变量声明*/
var (
	client *Client

	serverIp   string
	serverPort int

	msginput chan string
	// orderinput chan string

	Exit      [1]string
	Selecting [1]bool

	Tmp []string //消息列表转化

	functionlb = widget.NewLabel("功能选择")                              //功能标签
	selections = widget.NewSelect([]string{"公聊", "私聊", "修改用户名"}, nil) //功能选择标签
	//清空命令行
	OrderclearButton = widget.NewButton("清空命令行", func() {
		textOutOrder.Text = ""
		textOutOrder.Refresh()
	})
	textOutOrder    = widget.NewEntry()       //命令输出行
	OrderinputEntry = widget.NewEntry()       //命令输入行
	onlinelb        = widget.NewLabel("当前在线") //在线状态栏
	onlineList      = widget.NewList(         //在线列表
		func() int {
			return len(Tmp)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("nil")
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(Tmp[id])
		})
	clearButton = widget.NewButton("清空消息列表", func() { //清空消息列表
		msgoutEntry.Text = ""
		msgoutEntry.Refresh()
	})
	//结束当前服务按钮
	cut = widget.NewButton("结束当前服务", func() {
		Exit[0] = "$cut$"
		msginput <- Exit[0]
	})
	statuslb      = widget.NewLabel("连接服务器失败") //连接状态栏
	msginputEntry = widget.NewEntry()          //输入框
	msgoutEntry   = widget.NewEntry()          //聊天框

)

//客户端结构体
type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

/*客户端↓*/

//初始化服务端地址
func init() {
	// flag.StringVar(&serverIp, "ip", "101.43.19.232", "设置服务器IP地址(默认为101.43.19.232)")
	// flag.IntVar(&serverPort, "port", 82, "设置服务器端口(默认为82)")
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认为127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认为8888)")
}

//初始化客户端
func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	conn, error := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if error != nil {
		fmt.Println("net dial error...")
		return nil
	}
	client.conn = conn
	return client
}

/*客户端参数生成函数群↓*/
//修改用户名
func (c *Client) UpdateName() {

	textOutOrder.Text += "请输入修改的用户名:\n"
	textOutOrder.MultiLine = true
	textOutOrder.Refresh()

	c.Name = <-msginput

	if c.Name == "$cut$" {
		textOutOrder.Text += "已结束修改用户名进程\n"
		textOutOrder.MultiLine = true
		textOutOrder.Refresh()

		Selecting[0] = false

		return
	}

	if c.Name != "$cut$" {
		sendMsg := "rename|" + c.Name + "\n"
		_, error := c.conn.Write([]byte(sendMsg))
		if error != nil {
			textOutOrder.Text += "连接服务器失败...\n"
			textOutOrder.MultiLine = true
			textOutOrder.Refresh()
			return
		}

		textOutOrder.Text += "修改完成\n"
		textOutOrder.MultiLine = true
		textOutOrder.Refresh()
	}

}

//公聊
func (c *Client) PublicChat() {
	var chatMsg string

	textOutOrder.Text += "请在右侧聊天框输入聊天内容\n"
	textOutOrder.MultiLine = true
	textOutOrder.Refresh()

	chatMsg = <-msginput

	if chatMsg == "$cut$" {
		textOutOrder.Text += "已结束公聊进程\n"
		textOutOrder.MultiLine = true
		textOutOrder.Refresh()

		Selecting[0] = false

		return
	}

	for chatMsg != "$cut$" {
		if len(chatMsg) != 0 {
			msg := "all|" + chatMsg + "\n\n"
			_, error := c.conn.Write([]byte(msg))
			if error != nil {
				textOutOrder.Text += "连接服务器失败...\n"
				textOutOrder.MultiLine = true
				textOutOrder.Refresh()
				break
			}
		}

		chatMsg = ""

		textOutOrder.Text += "请在右侧聊天框输入聊天内容\n"
		textOutOrder.MultiLine = true
		textOutOrder.Refresh()

		chatMsg = <-msginput

		if chatMsg == "$cut$" {
			textOutOrder.Text += "已结束公聊进程\n"
			textOutOrder.MultiLine = true
			textOutOrder.Refresh()

			Selecting[0] = false

			break
		}

	}

}

//私聊
func (c *Client) PrivateChat() {
	var remoteUser string
	var chatMsg string

	textOutOrder.Text += "请在右侧聊天框输入私聊对象\n"
	textOutOrder.MultiLine = true
	textOutOrder.Refresh()

	remoteUser = <-msginput

	if remoteUser == "$cut$" {
		textOutOrder.Text += "已结束私聊进程\n"
		textOutOrder.MultiLine = true
		textOutOrder.Refresh()
		Selecting[0] = false
		return
	}

	//循环标签
OuterLoop:
	for remoteUser != "exit" && chatMsg != "$cut$" {

		textOutOrder.Text += "请在右侧聊天框输入聊天内容,输入exit返回上一级\n"
		textOutOrder.MultiLine = true
		textOutOrder.Refresh()

		chatMsg = <-msginput

		if chatMsg == "$cut$" {
			textOutOrder.Text += "已结束私聊进程\n"
			textOutOrder.MultiLine = true
			textOutOrder.Refresh()
			Selecting[0] = false
			break OuterLoop
		}

		for chatMsg != "exit" && chatMsg != "$cut$" {

			if len(chatMsg) != 0 {
				msg := "to|" + remoteUser + "|" + chatMsg + "\n\n"
				_, error := c.conn.Write([]byte(msg))
				if error != nil {
					textOutOrder.Text += "连接服务器失败...\n"
					textOutOrder.MultiLine = true
					textOutOrder.Refresh()
					break
				}
			}

			chatMsg = ""
			textOutOrder.Text += "请在右侧聊天框输入聊天内容,输入exit返回上一级\n"
			textOutOrder.MultiLine = true
			textOutOrder.Refresh()

			chatMsg = <-msginput

			if chatMsg == "$cut$" {
				textOutOrder.Text += "已结束私聊进程\n"
				textOutOrder.MultiLine = true
				textOutOrder.Refresh()
				Selecting[0] = false
				break OuterLoop
			}
		}

		remoteUser = ""
		textOutOrder.Text += "请在右侧聊天框输入私聊对象\n"
		textOutOrder.MultiLine = true
		textOutOrder.Refresh()

		remoteUser = <-msginput

		if chatMsg == "$cut$" {
			textOutOrder.Text += "已结束私聊进程\n"
			textOutOrder.MultiLine = true
			textOutOrder.Refresh()
			Selecting[0] = false
			break OuterLoop
		}
	}

}

//查询在线用户
func (c *Client) SelectUsers() {

	sendMsg := "onlist|/$/\n"
	_, error := c.conn.Write([]byte(sendMsg))
	if error != nil {
		textOutOrder.Text += "连接服务器失败...\n"
		textOutOrder.MultiLine = true
		textOutOrder.Refresh()
		return
	}
	time.Sleep(100 * time.Millisecond)

}

/*客户端参数生成函数群↑*/

//处理server返回的消息
func (c *Client) DealResponse() {
	buf := make([]byte, 4096)

	n, error := c.conn.Read(buf)
	if n > 0 {
		msg := string(buf[:n])

		if len(msg) > 8 && msg[:7] == "%BC|/$/" { //广播
			tmsg := strings.Split(msg, "/$/")[1]
			msgoutEntry.Text += tmsg
			msgoutEntry.Refresh()
		} else if len(msg) > 8 && msg[:7] == "%SC|/$/" { //私聊
			tmsg := strings.Split(msg, "/$/")[1]
			msgoutEntry.Text += tmsg
			msgoutEntry.Refresh()
		} else if len(msg) > 8 && msg[:7] == "%OL|/$/" { //在线列表
			onlineSlice := strings.Split(msg, "%OL|/$/")
			Tmp = onlineSlice[1:]
		} else { //其他返回处理
			textOutOrder.Text += msg
			textOutOrder.Refresh()
		}
	}

	if error != nil && error != io.EOF {
		textOutOrder.Text += "读取服务器返回数据出错\n"
		textOutOrder.Refresh()
	}

}

//服务选择
func Selections() {
	selections.OnChanged = func(s string) {
		switch s {
		case "公聊":
			Selecting[0] = true
			textOutOrder.Text = "已进入公聊模式"
			textOutOrder.Refresh()
			go client.PublicChat()
		case "私聊":
			Selecting[0] = true
			textOutOrder.Text = "已进入私聊模式"
			textOutOrder.Refresh()
			go client.PrivateChat()
		case "修改用户名":
			Selecting[0] = true
			textOutOrder.Text = "现在可以修改用户名"
			textOutOrder.Refresh()
			go client.UpdateName()
		}
	}

}

//初始化客户端
func initWindow() {
	flag.Parse()
	client = NewClient(serverIp, serverPort)
	if client == nil {
		statuslb.Text = "连接服务器失败"
		statuslb.Refresh()
		OrderinputEntry.Disable()
		msginputEntry.Disable()
		OrderinputEntry.Refresh()
		msginputEntry.Refresh()
		return
	} else {
		statuslb.Text = "连接服务器成功"
		statuslb.Refresh()
	}

	//初始化channel
	msginput = make(chan string)
	// orderinput = make(chan string)

	//发送请求在线用户
	go SenduserGet()

	//处理返回信息
	go DealwithMsg()

	// //防止channel阻塞
	go ChanJudge()

	//启动功能选择
	Selections()

}

//在线用户
func SenduserGet() {

	for {
		// //发送列表获取指令   这个位置发送指令就卡住了公聊私聊
		client.SelectUsers()
		//更新在线列表
		updateOnlineList()
	}

}

//处理返回信息
func DealwithMsg() {
	for {
		//处理返回值
		client.DealResponse()
	}
}

//检测channel阻塞
func ChanJudge() {
	for {
		if !Selecting[0] {
			DeleteStr := <-msginput
			if DeleteStr != "" {
				textOutOrder.Text += "channel阻塞,已清空阻塞\n"
				textOutOrder.Refresh()
				fmt.Println(DeleteStr)
			}

		}
	}
}

/*客户端↑*/

/*GUI↓*/

//用户list更新
func updateOnlineList() {

	onlineList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
		item.(*widget.Label).SetText(Tmp[id])
	}
	onlineList.Refresh()

}

//指令行输入(后续更新)
func InputSubmitf(s string) {
	if s != "" {

		// orderinput <- s

		//回显(可删除)
		// textOutOrder.Text += "[Order]" + s + "\n"
		// textOutOrder.MultiLine = true //多行输入
		// textOutOrder.Refresh()

		//清空指令输入框
		OrderinputEntry.Text = ""
		OrderinputEntry.Refresh()

	}
}

func InputSubmits(s string) {
	if s != "" {

		msginput <- s

		//回显
		// msgoutEntry.Text += "[me]" + s + "\n"
		// msgoutEntry.MultiLine = true //多行输入
		// msgoutEntry.Refresh()

		//清空消息输入框
		msginputEntry.Text = ""
		msginputEntry.Refresh()

	}
}

func initGUI() {
	//初始化窗口
	a := app.New()
	a.Settings().SetTheme(&f.MyTheme{})

	w := a.NewWindow("Client")
	w.SetFixedSize(true)

	w.Resize(fyne.NewSize(800, 600))

	//元素初始化
	selections.PlaceHolder = "选择功能"

	textOutOrder.Wrapping = fyne.TextWrapBreak //修饰新行带有滚动条
	textOutOrder.MultiLine = true
	textOutOrder.SetPlaceHolder("命令输出行") //背景字
	// textOutOrder.Disable()               //禁用输入

	OrderinputEntry.SetPlaceHolder("指令输入行") //背景字
	OrderinputEntry.OnSubmitted = InputSubmitf

	msgoutEntry.Wrapping = fyne.TextWrapBreak //修饰新行带有滚动条
	msgoutEntry.MultiLine = true
	msgoutEntry.SetPlaceHolder("消息列表") //背景字
	// msgoutEntry.Disable()              //禁用输入

	msginputEntry.SetPlaceHolder("消息输入,回车发送") //背景字
	msginputEntry.OnSubmitted = InputSubmits

	//布局
	part1 := container.NewVBox(functionlb, widget.NewSeparator(), selections, cut)
	part2 := container.NewBorder(OrderclearButton, OrderinputEntry, nil, nil, textOutOrder)
	containlp := container.New(layout.NewGridLayout(2), part1, part2)

	leafTopBox := container.NewVBox(clearButton)
	leftBox := container.NewBorder(onlinelb, leafTopBox, nil, nil, onlineList)
	containText := container.NewBorder(statuslb, msginputEntry, nil, nil, msgoutEntry)
	containrg := container.New(layout.NewGridLayout(2), leftBox, containText)

	containMain := container.NewGridWithColumns(2, containlp, containrg)

	w.SetContent(containMain)

	//初始化连接服务器
	initWindow()

	w.ShowAndRun()
}

/*GUI↑*/

//入口函数
//所有 goroutine 会在 main() 函数结束时一同结束
func main() {
	initGUI()
}
