# YABL api
## CLI
yabl-client为命令行程序，CLI交互使用[cobra](https://github.com/spf13/cobra)库进行配置和生成。  
yabl-client主要有两种用法
1. yabl version  
展示yabl服务器与解释器的版本号、运行OS信息、go工具链版本和编译时间  
1. yabl [-a/--address] [connect address] [-p/--port] [connect port]  
连接yabl服务器，连接服务器的IP地址和端口可由参数指定。之后将由命令行与服务器进行交互。  

还有completion、help等命令，--help参数等，为cobra库默认提供的功能，此处省略说明。
## ws接口
为连接yabl服务器与解释器，yabl-client使用websocket协议作为接口，连接api路径的设为/ws。  
yabl服务器与解释器和客户端通过websocket交换websocket标准TextMessage，内容为json格式的文本，是一个结构的序列化的json。该结构如下：  
```go
type MsgStruct struct {
	Timestamp int64  `json:"timestamp"`
	Content   string `json:"content"`
}
```
序列化后的文本，例如：
```json
{
    "timestamp": 1639889248,
    "content": "你好"
}
```
Timestamp为unix时间戳，记录发送方发送时间。
Content为具体文本内容。