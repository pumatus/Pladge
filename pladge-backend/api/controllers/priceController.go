package controllers

import (
	"net/http"
	"pledge-backend/api/models/ws"
	"pledge-backend/log"
	"pledge-backend/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type PriceController struct {
}

// NewPrice 处理 WebSocket 连接请求，实现价格实时推送
func (c *PriceController) NewPrice(ctx *gin.Context) {
	// 1. 异常恢复：捕获当前协程可能出现的 Panic，防止因单个连接崩溃导致整个进程退出
	defer func() {
		recoverRes := recover()
		if recoverRes != nil {
			log.Logger.Sugar().Error("new price recover ", recoverRes)
		}
	}()
	// 2. 配置 Upgrader（升级器）：将 HTTP 协议提升为 WebSocket 协议
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:   1024,            // 读取缓冲区大小
		WriteBufferSize:  1024,            // 写入缓冲区大小
		HandshakeTimeout: 5 * time.Second, // 握手超时时间
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许跨域（在生产环境建议限制特定域名）
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Logger.Sugar().Error("websocket request err:", err)
		return
	}

	// 3. 生成唯一连接 ID (randomId)：
	// 优先使用用户的远程 IP（点号替换为下划线）拼上随机字符串，方便在后台管理连接
	randomId := ""
	// remoteIP, ok := ctx.RemoteIP()
	// if ok {
	// 	randomId = strings.Replace(remoteIP.String(), ".", "_", -1) + "_" + utils.GetRandomString(23)
	// } else {
	// 	randomId = utils.GetRandomString(32) // 如果获取不到 IP，则生成 32 位随机字符串
	// }
	remoteIP := ctx.RemoteIP() // 直接返回 string
	if remoteIP != "" {
		// 因为已经是 string 类型，不需要再调用 .String()
		randomId = strings.Replace(remoteIP, ".", "_", -1) + "_" + utils.GetRandomString(23)
	} else {
		randomId = utils.GetRandomString(32)
	}

	// 4. 初始化 WebSocket 服务对象
	server := &ws.Server{
		Id:       randomId,               // 连接唯一标识
		Socket:   conn,                   // 物理连接实例
		Send:     make(chan []byte, 800), // 消息待发缓冲区（通道长度 800）
		LastTime: time.Now().Unix(),      // 记录最后一次活跃时间（用于心跳检测）
	}

	// 5. 启动异步协程处理读写
	// 这个方法通常会在后台循环监听消息并从 Send 通道取数据下发
	go server.ReadAndWrite()
}
