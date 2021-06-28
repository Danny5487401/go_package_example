package rmq

//配置文件结构
type tCfg struct {
	Connects  []Connect  `json:Connects`
	Channels  []Channel  `json:Channels`
	Exchanges []Exchange `json:Exchanges`
	Queue     []Queue    `json:Queue`
	Pusher    []Pusher   `json:Pusher`
	Poper     []Poper    `json:Poper`
}

//连接结构
type Connect struct {
	Name string `json:Name`
	Addr string `json:Addr`
}

//信道结构
type Channel struct {
	Name     string `json:Name`
	Connect  string `json:Connect`
	QosCount int    `json:QosCount`
	QosSize  int    `json:QosSize`
}

//交换机绑定结构
type EBind struct {
	Destination string `json:Destination`
	Key         string `json:Key`
	NoWait      bool   `json:NoWait`
}

//交换机结构
type Exchange struct {
	Name        string                 `json:Name`
	Channel     string                 `json:Channel`
	Type        string                 `json:Type`
	Durable     bool                   `json:Durable`
	AutoDeleted bool                   `json:AutoDeleted`
	Internal    bool                   `json:Internal`
	NoWait      bool                   `json:NoWait`
	Bind        []EBind                `json:Bind`
	Args        map[string]interface{} `json:Args`
}

//队列绑定结构
type QBind struct {
	ExchangeName string `json:ExchangeName`
	Key          string `json:Key`
	NoWait       bool   `json:NoWait`
}

//队列结构
type Queue struct {
	Name       string                 `json:Name`
	Channel    string                 `json:Channel`
	Durable    bool                   `json:Durable`
	AutoDelete bool                   `json:AutoDelete`
	Exclusive  bool                   `json:Exclusive`
	NoWait     bool                   `json:NoWait`
	Bind       []QBind                `json:Bind`
	Args       map[string]interface{} `json:Args`
}

//发送者配置
type Pusher struct {
	Name         string `json:Name`
	Channel      string `json:Channel`
	Exchange     string `json:Exchange`
	Key          string `json:Key`
	Mandtory     bool   `json:Mandtory`
	Immediate    bool   `json:Immediate`
	ContentType  string `json:ContentType`
	DeliveryMode uint8  `json:DeliveryMode`
}

//接收者配置
type Poper struct {
	Name      string `json:Name`
	QName     string `json:QName`
	Channel   string `json:Channel`
	Consumer  string `json:Consumer`
	AutoACK   bool   `json:AutoACK`
	Exclusive bool   `json:Exclusive`
	NoLocal   bool   `json:NoLocal`
	NoWait    bool   `json:NoWait`
}
