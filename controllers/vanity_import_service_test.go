package controllers

import (
	"context"
	"github.com/RivenZoo/govanityimport/proto/apidef"
	"reflect"
	"testing"
)

func TestGetController(t *testing.T) {
	tests := []struct {
		name string
		want *controller
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetController(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_controller_SetErrorResponse(t *testing.T) {
	type args struct {
		err    error
		resp   interface{}
		method int
	}
	tests := []struct {
		name string
		c    *controller
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.setErrorResponse(tt.args.err, tt.args.resp, tt.args.method)
		})
	}
}

func Test_controller_QueryImportMetaInfo(t *testing.T) {
	type args struct {
		ctx context.Context
		req *apidef.ImportMetaInfoReq
	}
	tests := []struct {
		name    string
		c       *controller
		args    args
		want    *apidef.ImportMetaInfoResp
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.QueryImportMetaInfo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("controller.QueryImportMetaInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("controller.QueryImportMetaInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_controller_UpdateModuleMetaInfo(t *testing.T) {
	type args struct {
		ctx context.Context
		req *apidef.UpdateModuleMetaInfoReq
	}
	tests := []struct {
		name    string
		c       *controller
		args    args
		want    *apidef.UpdateModuleMetaInfoResp
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.UpdateModuleMetaInfo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("controller.UpdateModuleMetaInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("controller.UpdateModuleMetaInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_init(t *testing.T) {
	tests := []struct {
		name string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_controller_queryModuleMetaInfo(t *testing.T) {
	type args struct {
		importPath string
	}
	tests := []struct {
		name    string
		c       *controller
		args    args
		want    *apidef.ModuleMetaInfo
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &controller{}
			got, err := c.queryModuleMetaInfo(tt.args.importPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("controller.queryModuleMetaInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("controller.queryModuleMetaInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_controller_deleteModuleMetaInfo(t *testing.T) {
	type args struct {
		ctx  context.Context
		req  *apidef.UpdateModuleMetaInfoReq
		resp *apidef.UpdateModuleMetaInfoResp
	}
	tests := []struct {
		name string
		c    *controller
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &controller{}
			c.deleteModuleMetaInfo(tt.args.ctx, tt.args.req, tt.args.resp)
		})
	}
}

func Test_controller_assignModuleMetaInfo(t *testing.T) {
	type args struct {
		src *apidef.ModuleMetaInfo
		dst *apidef.ModuleMetaInfo
	}
	tests := []struct {
		name string
		c    *controller
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &controller{}
			c.assignModuleMetaInfo(tt.args.src, tt.args.dst)
		})
	}
}

func Test_controller_setModuleMetaInfo(t *testing.T) {
	type args struct {
		ctx  context.Context
		req  *apidef.UpdateModuleMetaInfoReq
		resp *apidef.UpdateModuleMetaInfoResp
	}
	tests := []struct {
		name string
		c    *controller
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &controller{}
			c.setModuleMetaInfo(tt.args.ctx, tt.args.req, tt.args.resp)
		})
	}
}
