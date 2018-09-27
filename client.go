package vishnu

type Client struct {
	bridge string
	port   int
}

func NewClient(bridge string, port int) *Client {
	return &Client{
		bridge: bridge,
		port:   port,
	}
}

func (c *Client) Run() {
	// 1. 连接到bridge，记录本地连接端口
	// 2. 向bridge发起连接server请求
	// 3. listen本地连接端口，同时连接bridge返回的server地址
	// 4. 连接成功后，listen port提供服务
}
