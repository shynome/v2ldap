package v2ray

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"v2ray.com/core/app/proxyman/command"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

func (v2 V2ray) getClient() (grpcClient command.HandlerServiceClient, err error) {
	if v2.grpcClient != nil {
		return v2.grpcClient, nil
	}
	cc, err := grpc.Dial(v2.GrpcAddr, grpc.WithInsecure())
	if err != nil {
		return
	}
	v2.grpcClient = command.NewHandlerServiceClient(cc)
	grpcClient = v2.grpcClient
	return
}

func (v2 V2ray) addUser(user User) (err error) {
	hsClient, err := v2.getClient()
	if err != nil {
		return
	}
	v2User := &protocol.User{
		Email: user.Email,
		Account: serial.ToTypedMessage(&vmess.Account{
			Id:      user.UUID,
			AlterId: 64,
		}),
	}
	_, err = hsClient.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag:       v2.Tag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{User: v2User}),
	})
	if err != nil {
		return
	}
	return
}

var mux = sync.Mutex{}

func (v2 V2ray) loopUsers(fn func(user User) error) func(users []User) []error {
	return func(users []User) (errs []error) {
		mux.Lock()
		defer mux.Unlock()
		wait := make(chan int)
		errs = []error{}
		var finshedTaskCount = 0
		var allTaskCount = len(users)
		for _, user := range users {
			wrapper := func() {
				err := fn(user)
				finshedTaskCount++
				if err != nil {
					errs = append(errs, err)
				}
				if allTaskCount == finshedTaskCount {
					wait <- 1
				}
			}
			go wrapper()
		}
		<-wait
		return
	}
}

// AddUsers v2ray
func (v2 V2ray) AddUsers(users []User) (errs []error) {
	return v2.loopUsers(v2.addUser)(users)
}

func (v2 V2ray) removeUsers(user User) (err error) {
	hsClient, err := v2.getClient()
	if err != nil {
		return
	}
	_, err = hsClient.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag:       v2.Tag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{Email: user.Email}),
	})
	if err != nil {
		return
	}
	return
}

// RemoveUsers v2ray
func (v2 V2ray) RemoveUsers(users []User) (errs []error) {
	return v2.loopUsers(v2.removeUsers)(users)
}
