package common

import (
	"errors"
	"github.com/BabySid/gobase"
	"github.com/BabySid/gorpc"
	"github.com/BabySid/gorpc/api"
	"github.com/BabySid/gorpc/codec"
	"github.com/BabySid/proto/sodor"
	"reflect"
	"sync"
)

type sodorClient struct {
	fatCtrl FatCtrl
	handle  api.Client

	mutex sync.RWMutex

	methods map[SodorMethod]parameterType
}

type parameterType struct {
	method   string
	reqType  reflect.Type
	respType reflect.Type
}

var (
	errNullHandle = errors.New("null handle")
)

var (
	sodorHandle     *sodorClient
	sodorHandleOnce sync.Once
)

func GetSodorClient() *sodorClient {
	sodorHandleOnce.Do(func() {
		sodorHandle = &sodorClient{}
		sodorHandle.registerMethod()
	})
	return sodorHandle
}

type SodorMethod int

const (
	AddThomas SodorMethod = iota
	CreateAlertGroup
	CreateAlertPluginInstance
	CreateJob
	DeleteJob
	DeleteAlertGroup
	DeleteAlertPluginInstance
	DropThomas
	ListAlertGroup
	ListAlertPluginInstances
	ListJobs
	ListThomas
	RunJob
	SelectJob
	SelectJobInstances
	SelectTaskInstances
	ShowAlertPluginInstanceHistories
	ShowThomas
	UpdateAlertGroup
	UpdateAlertPluginInstance
	UpdateJob
	UpdateThomas
	maxMethod
)

func (c *sodorClient) registerMethod() {
	ns := "rpc."
	c.methods = make(map[SodorMethod]parameterType)

	c.methods[AddThomas] = parameterType{
		method:   ns + "AddThomas",
		reqType:  reflect.TypeOf(&sodor.ThomasInfo{}),
		respType: reflect.TypeOf(&sodor.ThomasReply{}),
	}
	c.methods[CreateAlertGroup] = parameterType{
		method:   ns + "CreateAlertGroup",
		reqType:  reflect.TypeOf(&sodor.AlertGroup{}),
		respType: reflect.TypeOf(&sodor.AlertGroupReply{}),
	}
	c.methods[CreateAlertPluginInstance] = parameterType{
		method:   ns + "CreateAlertPluginInstance",
		reqType:  reflect.TypeOf(&sodor.AlertPluginInstance{}),
		respType: reflect.TypeOf(&sodor.AlertPluginReply{}),
	}
	c.methods[CreateJob] = parameterType{
		method:   ns + "CreateJob",
		reqType:  reflect.TypeOf(&sodor.Job{}),
		respType: reflect.TypeOf(&sodor.JobReply{}),
	}
	c.methods[DeleteJob] = parameterType{
		method:   ns + "DeleteJob",
		reqType:  reflect.TypeOf(&sodor.Job{}),
		respType: reflect.TypeOf(&sodor.JobReply{}),
	}
	c.methods[DeleteAlertGroup] = parameterType{
		method:   ns + "DeleteAlertGroup",
		reqType:  reflect.TypeOf(&sodor.AlertGroup{}),
		respType: reflect.TypeOf(&sodor.AlertGroupReply{}),
	}
	c.methods[DeleteAlertPluginInstance] = parameterType{
		method:   ns + "DeleteAlertPluginInstance",
		reqType:  reflect.TypeOf(&sodor.AlertPluginInstance{}),
		respType: reflect.TypeOf(&sodor.AlertPluginReply{}),
	}
	c.methods[DropThomas] = parameterType{
		method:   ns + "DropThomas",
		reqType:  reflect.TypeOf(&sodor.ThomasInfo{}),
		respType: reflect.TypeOf(&sodor.ThomasReply{}),
	}
	c.methods[ListAlertGroup] = parameterType{
		method:   ns + "ListAlertGroup",
		reqType:  reflect.TypeOf(nil),
		respType: reflect.TypeOf(&sodor.AlertGroups{}),
	}
	c.methods[ListAlertPluginInstances] = parameterType{
		method:   ns + "ListAlertPluginInstances",
		reqType:  reflect.TypeOf(nil),
		respType: reflect.TypeOf(&sodor.AlertPluginInstances{}),
	}
	c.methods[ListJobs] = parameterType{
		method:   ns + "ListJobs",
		reqType:  reflect.TypeOf(nil),
		respType: reflect.TypeOf(&sodor.Jobs{}),
	}
	c.methods[ListThomas] = parameterType{
		method:   ns + "ListThomas",
		reqType:  reflect.TypeOf(nil),
		respType: reflect.TypeOf(&sodor.ThomasInfos{}),
	}
	c.methods[RunJob] = parameterType{
		method:   ns + "RunJob",
		reqType:  reflect.TypeOf(&sodor.Job{}),
		respType: reflect.TypeOf(&sodor.JobReply{}),
	}
	c.methods[SelectJob] = parameterType{
		method:   ns + "SelectJob",
		reqType:  reflect.TypeOf(&sodor.Job{}),
		respType: reflect.TypeOf(&sodor.Job{}),
	}
	c.methods[SelectJobInstances] = parameterType{
		method:   ns + "SelectJobInstances",
		reqType:  reflect.TypeOf(&sodor.JobInstance{}),
		respType: reflect.TypeOf(&sodor.JobInstances{}),
	}
	c.methods[SelectTaskInstances] = parameterType{
		method:   ns + "SelectTaskInstances",
		reqType:  reflect.TypeOf(&sodor.TaskInstance{}),
		respType: reflect.TypeOf(&sodor.TaskInstances{}),
	}
	c.methods[ShowAlertPluginInstanceHistories] = parameterType{
		method:   ns + "ShowAlertPluginInstanceHistories",
		reqType:  reflect.TypeOf(&sodor.AlertPluginInstanceHistory{}),
		respType: reflect.TypeOf(&sodor.AlertPluginInstanceHistories{}),
	}
	c.methods[ShowThomas] = parameterType{
		method:   ns + "ShowThomas",
		reqType:  reflect.TypeOf(&sodor.ThomasInfo{}),
		respType: reflect.TypeOf(&sodor.ThomasInstance{}),
	}
	c.methods[UpdateAlertGroup] = parameterType{
		method:   ns + "UpdateAlertGroup",
		reqType:  reflect.TypeOf(&sodor.AlertGroup{}),
		respType: reflect.TypeOf(&sodor.AlertGroupReply{}),
	}
	c.methods[UpdateAlertPluginInstance] = parameterType{
		method:   ns + "UpdateAlertPluginInstance",
		reqType:  reflect.TypeOf(&sodor.AlertPluginInstance{}),
		respType: reflect.TypeOf(&sodor.AlertPluginReply{}),
	}
	c.methods[UpdateJob] = parameterType{
		method:   ns + "UpdateJob",
		reqType:  reflect.TypeOf(&sodor.Job{}),
		respType: reflect.TypeOf(&sodor.JobReply{}),
	}
	c.methods[UpdateThomas] = parameterType{
		method:   ns + "UpdateThomas",
		reqType:  reflect.TypeOf(&sodor.ThomasInfo{}),
		respType: reflect.TypeOf(&sodor.ThomasReply{}),
	}

	gobase.True(len(c.methods) == int(maxMethod))
}

func (c *sodorClient) SetFatCtrlAddr(addr FatCtrl) error {
	handle, err := gorpc.Dial(addr.Addr, api.ClientOption{
		Codec: codec.ProtobufCodec,
	})
	if err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.handle = handle
	c.fatCtrl = addr

	return nil
}

func (c *sodorClient) GetFatCrl() FatCtrl {
	return c.fatCtrl
}

func (c *sodorClient) getHandle() api.Client {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.handle
}

func (c *sodorClient) Call(method SodorMethod, request interface{}, resp interface{}) error {
	params, ok := c.methods[method]
	gobase.True(ok)
	gobase.True(params.reqType == reflect.TypeOf(request))
	gobase.True(params.respType == reflect.TypeOf(resp))

	handle := c.getHandle()
	if handle == nil {
		return errNullHandle
	}
	return handle.CallJsonRpc(resp, params.method, request)
}
