package model

import (
	"github.com/go-redis/redis"
	"govanityimport/proto/apidef"
	"govanityimport/config"
	"govanityimport/errorcode"
	"github.com/golang/protobuf/jsonpb"
	"fmt"
	"bytes"
)

var (
	metaInfoModel    *ImportMetaInfoModel
	defaultMarshaler *jsonpb.Marshaler
)

type ImportMetaInfoModel struct {
	cli *redis.Client
}

func subDirKey(subDir string) string {
	return fmt.Sprintf("subdir:%s", subDir)
}

func importPathKey(importPath string) string {
	return fmt.Sprintf("module:%s", importPath)
}

func (m *ImportMetaInfoModel) SetSubDirsToRootModule(subDirs []string, root string) error {
	if len(subDirs) <= 0 {
		return errorcode.ErrInnerWrongParam
	}
	kv := make([]interface{}, 0, len(subDirs)*2)
	for i, _ := range subDirs {
		kv = append(kv, subDirKey(subDirs[i]), root)
	}
	res := m.cli.MSet(kv...)
	return res.Err()
}

func (m *ImportMetaInfoModel) DeleteSubDirs(subDirs []string) error {
	if len(subDirs) <= 0 {
		return errorcode.ErrInnerWrongParam
	}
	key := make([]string, 0, len(subDirs))
	for i, _ := range subDirs {
		key = append(key, subDirKey(subDirs[i]))
	}
	res := m.cli.Del(key...)
	return res.Err()
}

func (m *ImportMetaInfoModel) GetSubDirRoot(subDir string) (string, error) {
	if len(subDir) == 0 {
		return "", errorcode.ErrInnerWrongParam
	}
	res := m.cli.Get(subDirKey(subDir))
	return res.Result()
}

func (m *ImportMetaInfoModel) SetModuleMetaInfo(importPath string, info *apidef.ModuleMetaInfo) error {
	if info == nil || info.ImportInfo == nil || importPath != info.ImportInfo.ModuleImportPath {
		return errorcode.ErrInnerWrongParam
	}
	key := importPathKey(importPath)
	buf := bytes.NewBuffer(make([]byte, 0, 64))
	err := defaultMarshaler.Marshal(buf, info)
	if err != nil {
		return err
	}
	res := m.cli.Set(key, (interface{})(buf.Bytes()), 0)
	return res.Err()
}

func (m *ImportMetaInfoModel) DeleteModuleMetaInfo(importPath string) error {
	if len(importPath) == 0 {
		return errorcode.ErrInnerWrongParam
	}
	key := importPathKey(importPath)
	return m.cli.Del(key).Err()
}

func (m *ImportMetaInfoModel) GetModuleMetaInfo(importPath string) (*apidef.ModuleMetaInfo, error) {
	if len(importPath) == 0 {
		return nil, errorcode.ErrInnerWrongParam
	}
	key := importPathKey(importPath)
	res, err := m.cli.Get(key).Bytes()
	if err != nil {
		return nil, err
	}

	info := &apidef.ModuleMetaInfo{}
	err = jsonpb.Unmarshal(bytes.NewReader(res), info)
	return info, err
}

func GetImportMetaModel() *ImportMetaInfoModel {
	return metaInfoModel
}

func InitModel() error {
	cfg := config.GetConfig()
	cli := redis.NewClient(&redis.Options{
		Addr:       cfg.MetaInfoRedis.Addr,
		DB:         cfg.MetaInfoRedis.DB,
		Password:   cfg.MetaInfoRedis.Password,
		MaxRetries: 3,
	})
	metaInfoModel = &ImportMetaInfoModel{
		cli: cli,
	}
	defaultMarshaler = &jsonpb.Marshaler{
		EmitDefaults: true,
		OrigName:     true,
	}
	return nil
}

func Close() error {
	if metaInfoModel != nil {
		metaInfoModel.cli.Close()
	}
	return nil
}
