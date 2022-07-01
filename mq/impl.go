package mq

type Handler func(data interface{})

// AlarmHandler 创建一个报警处理函数
type AlarmHandler func() error

type Body struct {
	Raw   interface{} `json:"raw"`
	Retry int         `json:"retry"`
}

// MQueue 声明队列实现接口
type MQueue interface {
	// Register 注册一个主题
	Register(topics ...string) error
	// Publish 发布消息
	Publish(topic string, data ...interface{}) error
	// Subscribe 订阅消息, 并且指定工作协程数量
	Subscribe(topic string, handler Handler) error
	//Length 获取消息长度
	Length(topic string) (int64, error)
	//后台服务
	daemonServer() error
}
