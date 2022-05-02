# 粘包解包方式

### fix length
固定长度，服务器和客户端约定一个固定长度，发送和接收时都使用固定长度的缓存发送和接收。
优点：实现简单，适合底层协议如MTU
缺点：效率低，当发送的数据长度小于约定长度时，需要填充无效字符，浪费传输资源，性能低

### delimiter based
固定分隔符，服务器和客户端约定一个分隔符，比如;;等，使用分隔符检测一个数据包的结束。
优点：实现简单
缺点：需要避免分隔符和数据内容重复导致的误切割

### length field based frame decoder
定义自己的数据包头和包体，在包头上标注包体的长度，接收方通过按相同规则解析包头和包体。
缺点：实现较复杂
优点：灵活，包头上除了定义包体长度还能够定义一些其他的字段，比如序列号等，支持并发传输等高级功能

# 实现一个从 socket connection 中解码出 goim 协议的解码器