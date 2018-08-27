package model

import (
	"gopkg.in/redis.v6"
	"govanityimport/proto/apidef"
)

type ImportMetaInfoModel struct {
	cli *redis.Client
}

func (m *ImportMetaInfoModel) SetSubDirsToRootModule(subDirs []string, root string) error {

}

func (m *ImportMetaInfoModel) DeleteSubDirs(subDirs []string) error {

}

func (m *ImportMetaInfoModel) GetSubDirRoot(subDir string) (string, error) {

}

func (m *ImportMetaInfoModel) SetModuleMetaInfo(info *apidef.ModuleMetaInfo) error {

}

func (m *ImportMetaInfoModel) DeleteModuleMetaInfo(importPath string) error {

}

func (m *ImportMetaInfoModel) GetModuleMetaInfo(importPath string) (*apidef.ModuleMetaInfo, error) {

}

func GetImportMetaModel() *ImportMetaInfoModel {

}

func InitModel() error {

}

func Close() error {

}
