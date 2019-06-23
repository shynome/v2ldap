package v2ray

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	UUID "github.com/google/uuid"
	"github.com/parnurzeal/gorequest"
	"google.golang.org/grpc"
	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/command"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/socks"
)

var apiPort uint32 = 3001
var socksPort = 3002
var checkServerPort = 3003
var checkServerAddr = fmt.Sprintf("127.0.0.1:%v", checkServerPort)
var checkServerResp = UUID.New().String()

func initCheckServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(checkServerResp))
	})
	server := &http.Server{
		Addr:    checkServerAddr,
		Handler: mux,
	}
	server.ListenAndServe()
}

var config *core.Config

func startV2ray() (cmd *exec.Cmd, err error) {
	config = getV2rayConfig(apiPort)
	cmd = exec.Command("v2ray", "-config=stdin:", "-format=pb")
	var pbconfig []byte
	if pbconfig, err = proto.Marshal(config); err != nil {
		return
	}
	cmd.Stdin = bytes.NewBuffer(pbconfig)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		return
	}
	<-time.After(time.Second * 1) // 等待 1s v2ray 启动
	return
}

func addSocksProxy() (body string, err error) {
	addr := fmt.Sprintf("127.0.0.1:%v", apiPort)
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return
	}
	hsClient := command.NewHandlerServiceClient(cc)
	req := &command.AddInboundRequest{
		Inbound: &core.InboundHandlerConfig{
			Tag: "check",
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortRange: net.SinglePortRange(net.Port(socksPort)),
				Listen:    net.NewIPOrDomain(net.LocalHostIP),
			}),
			ProxySettings: serial.ToTypedMessage(&socks.ServerConfig{AuthType: socks.AuthType_NO_AUTH}),
		},
	}
	_, err = hsClient.AddInbound(context.Background(), req)
	if err != nil {
		return
	}
	return
}

func Test_V2rayConfig(t *testing.T) {
	blackoutListBak := blackoutList
	blackoutList = []string{}
	var err error
	var cmd *exec.Cmd

	go initCheckServer()

	check := func(checkServerResp string) {

		if cmd, err = startV2ray(); err != nil {
			t.Error(err)
			return
		}
		clean := func() {
			cmd.Process.Kill()
		}
		defer clean()

		var body string
		if body, err = addSocksProxy(); err != nil {
			t.Error(err)
			return
		}

		proxy := fmt.Sprintf("socks5://127.0.0.1:%v", socksPort)
		request := gorequest.New().Proxy(proxy)
		_, body, errs := request.Get(fmt.Sprintf("http://%v", checkServerAddr)).End()
		if errs != nil {
			return
		}

		if body != checkServerResp {
			t.Errorf("expect: %v , but get %v", checkServerResp, body)
			return
		}

	}

	check(checkServerResp)

	<-time.After(time.Second * 2)

	blackoutList = blackoutListBak
	check("")

}
