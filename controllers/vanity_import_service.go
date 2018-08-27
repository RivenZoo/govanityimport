package controllers

import (
	"context"
	"govanityimport/proto/apidef"
)

var (
	instance *controller
)

type controller struct {
}

func GetController() *controller {
	return instance
}

func init() {
	instance = &controller{}
}

func (c *controller) QueryImportMetaInfo(context.Context, *apidef.ImportMetaInfoReq) (*apidef.ImportMetaInfoResp, error) {
	return nil, nil
}

func (c *controller) UpdateModuleMetaInfo(context.Context, *apidef.UpdateModuleMetaInfoReq) (*apidef.UpdateModuleMetaInfoResp, error) {
	return nil, nil
}
