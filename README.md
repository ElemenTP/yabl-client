# YABL-Client
A CLI yabl client implementation in go & A yabl client in browser.
## Usage example
Using yabl as the server.  
Script:
```yaml
#!/usr/bin/yabl
address: 127.0.0.1
port: 8080
func main:
  - hello = "你好，"
  - hello = invoke joinfunc hello
  - postmsg hello
  - loop
  - answer = getmsg
  - flag = answer equal "测试"
  - if flag
  - break
  - fi
  - postmsg "试试告诉我测试"
  - pool
  - postmsg "测试结束，再见"

func joinfunc hello:
  - temp = hello join "世界"
  - return temp
```
Result:  
CLI client  
```
Connecting to server ws://127.0.0.1:8080/ws
2021-12-19T11:42:49+08:00 你好，世界
你好
2021-12-19T11:42:51+08:00 试试告诉我测试
测试？
2021-12-19T11:42:53+08:00 试试告诉我测试
测试
2021-12-19T11:42:55+08:00 测试结束，再见
Disconnected from server ws://127.0.0.1:8080/ws
```
web client  
```
下午5:36:08 你好，世界
下午5:36:13 你好
下午5:36:13 试试告诉我测试
下午5:36:16 测试？
下午5:36:16 试试告诉我测试
下午5:36:18 测试
下午5:36:18 测试结束，再见
Connection closed.
```
## Thanks
CLI client  
[spf13/cobra](https://github.com/spf13/cobra)  
[gorilla/websocket](https://github.com/gorilla/websocket)  
web client  
Copy paste code from [gorilla/websocket/examples/chat/home.html](https://github.com/gorilla/websocket/blob/master/examples/chat/home.html), just change some code in message sending and receiving to adapt yabl server message type.
