# GoProxy
轻量的socks5代理服务器,代码简单易读，就像sock原始协议一样



# 安装
```go
go get github.com/Xumany/GoProxy
```

# 案例
```go
s := socks5.NewServer(socks5.ServerCfg{
	    Port: 1080,
	    LogLevel:"error",
        FilePath:"./GoProxy.log"
})
err := s.Run()
if err != nil {
	log.Fatalln(err)
}
```
# 鸣谢
<https://github.com/wzshiming/socks5>