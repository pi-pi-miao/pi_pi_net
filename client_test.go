package pi_pi_net

import (
	"testing"
)

func TestClient_Dail(t *testing.T) {
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
}
