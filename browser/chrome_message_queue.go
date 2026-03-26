package browser

// 定义全局消息通道
type mess struct {
	ID      int
	Content string
}

var messageQueue = make(chan mess, 100) // 缓冲队列

var ConnTabDone = make(chan struct{})

// GetNextMsgID 获取自增的消息ID（线程安全）
func GetNextMsgID() int {
	mu.Lock()
	defer mu.Unlock()
	id := chromeInstance.NextID
	chromeInstance.NextID++
	return id
}
