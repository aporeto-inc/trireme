package rpcWrapper

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"time"

	"github.com/aporeto-inc/trireme/cache"
)

//RPCHdl exported
type RPCHdl struct {
	Client  *rpc.Client
	Channel string
}

type RPCWrapper struct {
	rpcClientMap *cache.Cache
}

func NewRPCWrapper() *RPCWrapper {
	rpcWrapper := new(RPCWrapper)
	rpcWrapper.rpcClientMap = cache.NewCache(nil)
	return rpcWrapper
}

const (
	maxRetries = 100
)

//NewRPCClient exported
//Will worry about locking later ... there is a small case where two callers
//call NewRPCClient from a different thread
func (r *RPCWrapper) NewRPCClient(contextID string, channel string) error {
	//establish new connection to context/container
	RegisterTypes()
	numRetries := 0
	client, err := rpc.DialHTTP("unix", channel)

	for err != nil {
		time.Sleep(5 * time.Millisecond)
		err = nil
		numRetries = numRetries + 1
		if numRetries < maxRetries {
			client, err = rpc.DialHTTP("unix", channel)
		} else {
			return err
		}
	}
	return r.rpcClientMap.Add(contextID, &RPCHdl{Client: client, Channel: channel})

}

//GetRPCClient exported
func (r *RPCWrapper) GetRPCClient(contextID string) (*RPCHdl, error) {
	val, err := r.rpcClientMap.Get(contextID)
	if err == nil {
		return val.(*RPCHdl), err
	}
	return nil, err
}

func getsharedKey() []byte {
	var sharedKey = []byte("sharedsecret")
	return sharedKey
}

//RemoteCall exported
func (r *RPCWrapper) RemoteCall(contextID string, methodName string, req *Request, resp *Response) error {
	var rpcBuf bytes.Buffer
	binary.Write(&rpcBuf, binary.BigEndian, req.Payload)
	digest := hmac.New(sha256.New, getsharedKey())
	digest.Write(rpcBuf.Bytes())
	req.HashAuth = digest.Sum(nil)
	rpcClient, err := r.GetRPCClient(contextID)
	if err != nil {
		fmt.Println("Cant find rpc handle")
		return err
	}
	return rpcClient.Client.Call(methodName, req, resp)

}

//CheckValidity exported
func (r *RPCWrapper) CheckValidity(req *Request) bool {
	var rpcBuf bytes.Buffer
	binary.Write(&rpcBuf, binary.BigEndian, req.Payload)
	digest := hmac.New(sha256.New, getsharedKey())
	digest.Write(rpcBuf.Bytes())
	return hmac.Equal(req.HashAuth, digest.Sum(nil))
}
func NewRPCServer() RPCServer {
	return new(RPCWrapper)
}

//StartServer exported
func (r *RPCWrapper) StartServer(protocol string, path string, handler interface{}) error {
	RegisterTypes()
	rpc.Register(handler)
	rpc.HandleHTTP()
	if len(path) == 0 {
		panic("Sock param not passed in environment")
	}
	listen, err := net.Listen(protocol, path)

	if err != nil {
		return err
	}
	go http.Serve(listen, nil)
	defer func() {
		listen.Close()
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		os.Remove(path)
	}
	return nil
}

func (r *RPCWrapper) DestroyRPCClient(contextID string) {
	rpcHdl, _ := r.rpcClientMap.Get(contextID)
	rpcHdl.(*RPCHdl).Client.Close()
	os.Remove(rpcHdl.(*RPCHdl).Channel)
}

func (r *RPCWrapper) ProcessMessage(req *Request) bool {
	return r.CheckValidity(req)
}

//RegisterTypes exported
func RegisterTypes() {
	gob.Register(InitRequestPayload{})
	gob.Register(InitResponsePayload{})
	gob.Register(InitSupervisorPayload{})

	gob.Register(EnforcePayload{})
	gob.Register(UnEnforcePayload{})

	gob.Register(SuperviseRequestPayload{})
	gob.Register(UnSupervisePayload{})
}
