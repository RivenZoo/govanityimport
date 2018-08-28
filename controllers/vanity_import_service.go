package controllers

import (
	"context"
	"govanityimport/proto/apidef"
	"govanityimport/errorcode"
	"govanityimport/headers"
	"govanityimport/model"
	"govanityimport/zaplog"
	"strings"
	"github.com/go-redis/redis"
)

const (
	UpdateModuleAction = "set"
	DeleteModuleAction = "del"
)

var (
	instance *controller
)

const (
	_                          = iota
	queryImportMetaInfoMethod  = iota
	updateModuleMetaInfoMethod
)

type controller struct {
}

func GetController() *controller {
	return instance
}

func init() {
	instance = &controller{}
}

func (c *controller) SetErrorResponse(err error, resp interface{}, method int) {
	e := errorcode.OK
	switch err.(type) {
	case errorcode.InnerError:
		e = err.(errorcode.InnerError).BaseError
	case errorcode.BaseError:
		e = err.(errorcode.BaseError)
	default:
		e = errorcode.ErrServerError
	}
	switch method {
	case queryImportMetaInfoMethod:
		ret := resp.(*apidef.ImportMetaInfoResp)
		ret.Ret = int32(e.Ret)
		ret.Msg = e.Msg
	case updateModuleMetaInfoMethod:
		ret := resp.(*apidef.UpdateModuleMetaInfoResp)
		ret.Ret = int32(e.Ret)
		ret.Msg = e.Msg
	}
}

func (c *controller) queryModuleMetaInfo(importPath string) (*apidef.ModuleMetaInfo, error) {
	m := model.GetImportMetaModel()
	metaInfo, err := m.GetModuleMetaInfo(importPath)
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
		// assume it is sub dir
		root, err := m.GetSubDirRoot(importPath)
		if err == nil {
			metaInfo, err = m.GetModuleMetaInfo(root)
		}
		if err != nil {
			return nil, err
		}
	}
	return metaInfo, nil
}

func (c *controller) QueryImportMetaInfo(ctx context.Context, req *apidef.ImportMetaInfoReq) (*apidef.ImportMetaInfoResp, error) {
	resp := &apidef.ImportMetaInfoResp{}
	requestID := ctx.Value(headers.HeaderRequestID).(string)
	method := ctx.Value(headers.ContextKeyRPCMethod)

	log := zaplog.GetSugarLogger()

	importPath := strings.TrimRight(req.ImportPath, "/")

	metaInfo, err := c.queryModuleMetaInfo(importPath)
	if err != nil {
		log.Errorw("get module meta info fail",
			"request_id", requestID,
			"method", method,
			"module", importPath,
			"error", err)
		c.SetErrorResponse(errorcode.ErrNoModuleInfo, resp, queryImportMetaInfoMethod)
		return resp, nil
	}

	c.SetErrorResponse(errorcode.OK, resp, queryImportMetaInfoMethod)
	resp.TraceId = requestID
	resp.MetaInfo = metaInfo
	return resp, nil
}

func (c *controller) UpdateModuleMetaInfo(ctx context.Context, req *apidef.UpdateModuleMetaInfoReq) (*apidef.UpdateModuleMetaInfoResp, error) {
	resp := &apidef.UpdateModuleMetaInfoResp{}

	switch req.Action {
	case UpdateModuleAction:
		c.setModuleMetaInfo(ctx, req, resp)
	case DeleteModuleAction:
		c.deleteModuleMetaInfo(ctx, req, resp)
	default:
		c.SetErrorResponse(errorcode.ErrBadRequest, resp, updateModuleMetaInfoMethod)
	}
	return resp, nil
}

func (c *controller) deleteModuleMetaInfo(ctx context.Context, req *apidef.UpdateModuleMetaInfoReq, resp *apidef.UpdateModuleMetaInfoResp) {
	requestID := ctx.Value(headers.HeaderRequestID).(string)
	method := ctx.Value(headers.ContextKeyRPCMethod)

	m := model.GetImportMetaModel()
	log := zaplog.GetSugarLogger()

	importPath := strings.TrimRight(req.MetaInfo.ImportInfo.ModuleImportPath, "/")
	err := m.DeleteModuleMetaInfo(importPath)
	if err != nil {
		log.Errorw("delete module fail",
			"request_id", requestID,
			"method", method,
			"module", importPath,
			"error", err)
		c.SetErrorResponse(errorcode.ErrNoModuleInfo, resp, updateModuleMetaInfoMethod)
		return
	}
	c.SetErrorResponse(errorcode.OK, resp, updateModuleMetaInfoMethod)
}

func (c *controller) assignModuleMetaInfo(src, dst *apidef.ModuleMetaInfo) {
	if src == nil || dst == nil {
		return
	}
	if len(src.SubImportDirs) != 0 {
		dst.SubImportDirs = src.SubImportDirs
	}
	if src.ImportInfo != nil {
		if dst.ImportInfo == nil {
			dst.ImportInfo = src.ImportInfo
		} else {
			assignNonEmptyByFieldName(src.ImportInfo, dst.ImportInfo,
				"ModuleImportPath", "RepoUrl", "Vcs")
		}
	}
	if src.SourceInfo != nil {
		if dst.SourceInfo == nil {
			dst.SourceInfo = src.SourceInfo
		} else {
			assignNonEmptyByFieldName(src.SourceInfo, dst.SourceInfo,
				"ModuleImportPath", "DirPattern", "DocHost", "FilePattern", "HomeUrl")
		}
	}
}

func (c *controller) setModuleMetaInfo(ctx context.Context, req *apidef.UpdateModuleMetaInfoReq, resp *apidef.UpdateModuleMetaInfoResp) {
	requestID := ctx.Value(headers.HeaderRequestID).(string)
	method := ctx.Value(headers.ContextKeyRPCMethod)
	resp.TraceId = requestID

	m := model.GetImportMetaModel()
	log := zaplog.GetSugarLogger()

	importPath := strings.TrimRight(req.MetaInfo.ImportInfo.ModuleImportPath, "/")
	metaInfo, err := c.queryModuleMetaInfo(importPath)
	if err != nil {
		if err != redis.Nil {
			log.Errorw("get module meta info fail",
				"request_id", requestID,
				"method", method,
				"module", importPath,
				"error", err)
			c.SetErrorResponse(errorcode.ErrServerError, resp, updateModuleMetaInfoMethod)
			return
		}
		// first set meta info
		err = m.SetModuleMetaInfo(importPath, req.MetaInfo)
		if err != nil {
			log.Errorw("set module meta info fail",
				"request_id", requestID,
				"method", method,
				"module", importPath,
				"error", err)
			c.SetErrorResponse(errorcode.ErrServerError, resp, updateModuleMetaInfoMethod)
			return
		}
		c.SetErrorResponse(errorcode.OK, resp, updateModuleMetaInfoMethod)
		return
	}
	c.assignModuleMetaInfo(req.MetaInfo, metaInfo)
	err = m.SetModuleMetaInfo(importPath, metaInfo)
	if err != nil {
		log.Errorw("set module meta info fail",
			"request_id", requestID,
			"method", method,
			"module", importPath,
			"error", err)
		c.SetErrorResponse(errorcode.ErrServerError, resp, updateModuleMetaInfoMethod)
		return
	}
	c.SetErrorResponse(errorcode.OK, resp, updateModuleMetaInfoMethod)
}
