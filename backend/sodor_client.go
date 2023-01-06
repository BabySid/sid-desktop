package backend

import (
	"sid-desktop/common"
	"sync"
)

type sodorClient struct {
	fatCtrl common.FatCtrl
}

var (
	sodorHandle *sodorClient
	sodorOnce   sync.Once
)

func GetSodorClient() *sodorClient {
	sodorOnce.Do(func() {
		sodorHandle = &sodorClient{}
	})
	return sodorHandle
}

func (c *sodorClient) SetFatCtrlAddr(addr common.FatCtrl) {
	c.fatCtrl = addr
}

func (c *sodorClient) GetFatCrl() common.FatCtrl {
	return c.fatCtrl
}
