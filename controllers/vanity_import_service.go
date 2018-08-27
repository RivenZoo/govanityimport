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

func (c *controller) QueryImportMetaInfo(ctx context.Context, req *apidef.ImportMetaInfoReq) (*apidef.ImportMetaInfoResp, error) {
	return nil, nil
}

func (c *controller) UpdateModuleMetaInfo(ctx context.Context, req *apidef.UpdateModuleMetaInfoReq) (*apidef.UpdateModuleMetaInfoResp, error) {
	return nil, nil
}
