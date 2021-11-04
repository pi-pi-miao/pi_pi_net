package main

import (
	"fmt"
	"pi_pi_net"
)

func main(){
	//runExample()
	listenExample()
}

func listenExample(){
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

func readBlock(ctx *pi_pi_net.Context){
	i := 0
	for messageServer := range ctx.ReadBlock() {
		fmt.Println("[read data]",messageServer.ReadString())
		messageServer.WriteString(fmt.Sprintf("hello client%v",i))
		i ++
		messageServer.Close()
	}
}

func runExample(){
	ctx := pi_pi_net.NewContext()
	fmt.Println("run")
	go readBlock(ctx)
	if err := ctx.Run("tcp","127.0.0.1:10000");err != nil {
		fmt.Println("run err",err)
		return
	}
}

