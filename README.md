
<h1 align="center">pi-pi-ne 👋</h1>

<p align="center">
  <img src="https://img.shields.io/badge/linux-epoll-yellowgreen" />
  <img src="https://img.shields.io/badge/linux-golang-red" />
  <img src="https://img.shields.io/badge/linux-pi__pi__net-orange" />
  <img src="https://img.shields.io/badge/linux-pi--pi--miao-brightgreen" />
</p>

 🎉
```azure
pi-pi-net 是一个在linux环境下封装epoll的网络库,可以基于此库非常方便的实现Reactor网络模型,或者web，rpc，websocket等网络框架的基础框架
```


## ✨ [详细文档点这里](https://pkg.go.dev/github.com/pi-pi-miao/pi_pi_net)
## ✨ [Develop detailed documentation](https://pkg.go.dev/github.com/pi-pi-miao/pi_pi_net)


## 服务端使用

### 👊 服务端使用有两种方式,方式一
```go
首先创建context :  ctx := pi_pi_net.NewContext()
然后直接读取  :  go readBlock(ctx)
最后运行即可:   ctx.Run("tcp","127.0.0.1:10000")


完整代码如下

func readBlock(ctx *pi_pi_net.Context){
    i := 0
    for messageServer := range ctx.ReadBlock() {
        fmt.Println("[read data]",messageServer.ReadString())
        messageServer.WriteString(fmt.Sprintf("hello client%v",i))
        i ++
        messageServer.Close()
    }
}

func main(){
    ctx := pi_pi_net.NewContext()
    fmt.Println("run")
    go readBlock(ctx)
    if err := ctx.Run("tcp","127.0.0.1:10000");err != nil {
        fmt.Println("run err",err)
        return
    }
}
```

### 👊 服务端使用有两种方式,方式二

```go

1 首先监听地址  ctx,err := pi_pi_net.NewContext().Listen("tcp","127.0.0.1:10000")
2 异步读取  :  go read(ctx)
3 获取连接然后运行: ctx.Accept()

func main(){
	ctx,err := pi_pi_net.NewContext().Listen("tcp","127.0.0.1:10000")
	if err != nil {
		fmt.Println("[listen err]",err)
		return
	}
	go read(ctx)
	fmt.Println("run listen")
	if err := ctx.Accept();err != nil {
		fmt.Println("[accept err]",err)
	}
}

func read(ctx *pi_pi_net.Context){
	obj := ctx.ReadServerConnect()
	fmt.Println("[read data]",obj.ReadString())
	if n,err := obj.WriteString("hello client");err != nil {
		fmt.Println("write err",err)
	}else {
		fmt.Println("write n is ",n)
	}
	obj.Close()
}
```

## 👊 客户端使用

```go
1 客户端连接服务端 	cli,err := Dail("tcp","127.0.0.1:10000")
2 读取数据 data,err := cli.ReadString()
3 写入数据  cli.WriteString("hello world")


	cli,err := Dail("tcp","127.0.0.1:10000")
	if err != nil {
		t.Error("dail err",err)
		return
	}
	_,err = cli.WriteString("hello world")
	if err != nil {
		t.Error("cli writeString err",err)
		return
	}
	//t.Log("write len is ",n)
	data,err := cli.ReadString()
	if err != nil {
		t.Error("read err",err)
		return
	}
	t.Log("data is",data)
```

## 👊 客户端可以api

```go
Read()   :读取byte数组数据,阻塞
ReadString()  : 读取字符串,阻塞
Write(data []byte)  :写入数据
WriteString(data string) :写入数据
Close()   :关闭连接

```

## 🔨 服务端api

```go
ReadServerConnect()  :单次读取
ReadBlock()          :阻塞读取,返回chan
ReadByte()           :按照byte数组读取
ReadString()         :按照字符串读取
Write(data []byte)    :按照byte数组写入
WriteString(data string)  : 按照字符串写入
WriteByteAndClose(data []byte)  : 写入byte数组之后关闭连接
WriteStringAndClose(data string) :写入字符串后关闭连接
Close()  :关闭连接
```




## 🔋 PR
> 我当然希望你能够有更好的想法或贡献。

🚥 我可能需要更多的灵感和建议、如果您有更好的想法，请联系我 🙌🙌


















