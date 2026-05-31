package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Hub WebSocket 中心，管理iTick连接和前端客户端
type Hub struct {
	baseURL  string
	token    string
	itickWS  *websocket.Conn
	itickMu  sync.Mutex

	clients   map[*Client]bool
	clientsMu sync.RWMutex

	writeMu sync.Mutex

	currentCodes []string
	currentRegion string
	codesMu      sync.RWMutex
}

// Client 前端WebSocket客户端
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

// ITickQuote iTick WebSocket推送的quote数据
type ITickQuote struct {
	Code   int            `json:"code"`
	Data   ITickQuoteData `json:"data"`
}

// ITickQuoteData quote数据字段
type ITickQuoteData struct {
	Symbol       string  `json:"s"`    // 股票代码
	Region       string  `json:"r"`    // 市场区域
	LastPrice    float64 `json:"ld"`   // 最新价
	Open         float64 `json:"o"`    // 开盘价
	High         float64 `json:"h"`    // 最高价
	Low          float64 `json:"l"`    // 最低价
	PrevClose    float64 `json:"p"`    // 前收盘价
	Volume       float64 `json:"v"`    // 成交量
	Turnover     float64 `json:"tu"`   // 成交额
	Change       float64 `json:"ch"`   // 涨跌额
	ChangePercent float64 `json:"chp"` // 涨跌幅百分比
	Timestamp    int64   `json:"t"`    // 时间戳(ms)
	Type         string  `json:"type"` // 数据类型: quote
}

// WSMessage 前端WebSocket消息
type WSMessage struct {
	Action  string   `json:"action"`  // subscribe / unsubscribe
	Region  string   `json:"region"`  // 市场区域
	Codes   []string `json:"codes"`   // 股票代码列表
}

// WSStockData 推送给前端的股票数据
type WSStockData struct {
	Type          string  `json:"type"`           // quote
	Code          string  `json:"code"`           // 股票代码
	Region        string  `json:"region"`         // 市场区域
	Price         float64 `json:"price"`          // 最新价
	Change        float64 `json:"change"`         // 涨跌额
	ChangePercent float64 `json:"change_percent"` // 涨跌幅
	Volume        float64 `json:"volume"`         // 成交量
	Open          float64 `json:"open"`           // 开盘价
	High          float64 `json:"high"`           // 最高价
	Low           float64 `json:"low"`            // 最低价
	PrevClose     float64 `json:"prev_close"`     // 前收盘价
}

// NewHub 创建WebSocket Hub
func NewHub(baseURL, token string) *Hub {
	return &Hub{
		baseURL: baseURL,
		token:   token,
		clients: make(map[*Client]bool),
	}
}

// Run 启动Hub，连接iTick WebSocket
func (h *Hub) Run() error {
	wsBaseURL := strings.Replace(h.baseURL, "https://", "wss://", 1)
	wsBaseURL = strings.Replace(wsBaseURL, "http://", "ws://", 1)

	header := http.Header{}
	header.Set("token", h.token)

	conn, _, err := websocket.DefaultDialer.Dial(wsBaseURL+"/stock", header)
	if err != nil {
		return err
	}
	h.itickWS = conn
	log.Printf("✅ 已连接iTick WebSocket: %s/stock", wsBaseURL)

	// 启动读取iTick数据的协程（包含认证消息）
	go h.readITick()

	return nil
}

// readITick 持续读取iTick推送的数据
// 第一个消息是认证响应，后续是实时数据
func (h *Hub) readITick() {
	defer func() {
		log.Printf("⚠️ iTick WebSocket 读取协程退出")
	}()

	// 1. 读取认证响应
	_, reader, err := h.itickWS.NextReader()
	if err != nil {
		log.Printf("❌ 读取iTick认证消息失败: %v", err)
		return
	}
	var authMsg map[string]interface{}
	if err := json.NewDecoder(reader).Decode(&authMsg); err != nil {
		log.Printf("❌ 解析iTick认证消息失败: %v", err)
		return
	}
	log.Printf("📡 iTick 认证响应: %v", authMsg)

	// 2. 循环读取实时数据
	for {
		_, reader, err := h.itickWS.NextReader()
		if err != nil {
			log.Printf("❌ 读取iTick数据失败: %v", err)
			return
		}

		// 先用原始map读取，以便记录所有消息
		var rawMsg map[string]interface{}
		if err := json.NewDecoder(reader).Decode(&rawMsg); err != nil {
			log.Printf("⚠️ 解析iTick数据失败: %v", err)
			continue
		}

		// 检查是否是订阅确认
		if resAc, ok := rawMsg["resAc"]; ok {
			log.Printf("📡 iTick 响应: resAc=%v, msg=%v, code=%v", resAc, rawMsg["msg"], rawMsg["code"])
			continue
		}

		// 解析为quote数据
		var quote ITickQuote
		rawBytes, _ := json.Marshal(rawMsg)
		if err := json.Unmarshal(rawBytes, &quote); err != nil {
			continue
		}

		if quote.Code != 1 || quote.Data.Type != "quote" {
			log.Printf("📡 iTick 其他消息: %s", string(rawBytes))
			continue
		}

		stockData := WSStockData{
			Type:          "quote",
			Code:          quote.Data.Symbol,
			Region:        quote.Data.Region,
			Price:         quote.Data.LastPrice,
			Change:        quote.Data.Change,
			ChangePercent: quote.Data.ChangePercent,
			Volume:        quote.Data.Volume,
			Open:          quote.Data.Open,
			High:          quote.Data.High,
			Low:           quote.Data.Low,
			PrevClose:     quote.Data.PrevClose,
		}

		data, _ := json.Marshal(stockData)
		h.broadcast(data)
	}
}

// SubscribeITick 向iTick订阅股票
func (h *Hub) SubscribeITick(region string, codes []string) error {
	h.itickMu.Lock()
	defer h.itickMu.Unlock()

	if h.itickWS == nil {
		return nil
	}

	// 构造订阅参数: AAPL$US,NVDA$US,...
	params := make([]string, len(codes))
	for i, code := range codes {
		params[i] = code + "$" + region
	}

	subMsg := map[string]interface{}{
		"ac":     "subscribe",
		"params": strings.Join(params, ","),
		"types":  "quote",
	}

	log.Printf("📡 向iTick订阅 %d 只股票: region=%s", len(codes), region)
	return h.itickWS.WriteJSON(subMsg)
}

// broadcast 广播数据给所有前端客户端
func (h *Hub) broadcast(data []byte) {
	h.clientsMu.RLock()
	defer h.clientsMu.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- data:
		default:
			// 客户端缓冲区满，跳过
		}
	}
}

// Register 注册前端客户端
func (h *Hub) Register(client *Client) {
	h.clientsMu.Lock()
	h.clients[client] = true
	h.clientsMu.Unlock()
	log.Printf("🔌 新客户端连接，当前连接数: %d", len(h.clients))
}

// Unregister 注销前端客户端
func (h *Hub) Unregister(client *Client) {
	h.clientsMu.Lock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.Send)
	}
	h.clientsMu.Unlock()
	log.Printf("🔌 客户端断开，当前连接数: %d", len(h.clients))
}

// ReadPump 读取前端客户端消息（订阅请求）
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		if msg.Action == "subscribe" {
			log.Printf("📡 前端请求订阅: region=%s, codes=%d只", msg.Region, len(msg.Codes))
			c.Hub.SubscribeITick(msg.Region, msg.Codes)
		}
	}
}

// WritePump 推送数据给前端客户端
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for data := range c.Send {
		c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			break
		}
	}
}
