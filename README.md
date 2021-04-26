# go-syntax-ComSys
Golang语法学习，实现即时通信系统的基本功能  
### v0.1 实现基础server（监听，连接功能）
1. 端口8888
2. 实现创建监听
### v0.2 用户上线功能实现（用户上线，向全体发送广播）
1. User结构体记录用户的Name Addr Channel Conn
2. 每个User有一个Go程，通过Channel收到信息后，通过Conn发送出去
3. Server实现一个listenMessage监听广播用Message的Channel
4. 每个用户上线即调用一次Handel，用于广播上线消息
### v0.3 用户消息广播
1. Handle中加入循环读取Conn发来的信息，读取到则调用Broadcast广播
### v0.4 用户业务封装
1. 优化程序逻辑，将Handle移除，对应逻辑放置在User内，同时User需要保存Server指针用于操作在线图
2. 当前逻辑为，Server保存在线图，提供广播方法，User提供上下线时的操作，以及使用两个Go程分别进行信息收发
