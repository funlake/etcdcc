package services

import (
	"etcdcc/apiserver/pkg/dao"
	"etcdcc/apiserver/pkg/dto"
)

type Config struct{}

func (c Config) List(start int, limit int, q []string) (interface{}, int64) {
	configDao := &dao.CenterConfig{}
	configDao.SetPageParams(start, limit)
	configDao.SetSearchCdt(q)
	configDao.SetSearchMap(dto.CONFIG_SEARCH)
	return configDao.List()
}
func (c Config) Create(cdto *dto.ConfigAddDto) (int64, error) {
	configDao := &dao.CenterConfig{
		Key:  cdto.Key,
		Val:  cdto.Val,
		Env:  cdto.Env,
		Mod:  cdto.Mod,
		Type: cdto.Type,
	}
	return configDao.Create()
}
func (c Config) Update(cdto *dto.ConfigEditDto) error {
	configDao := &dao.CenterConfig{
		Id: cdto.Id,
	}
	err := configDao.Find()
	if err != nil {
		return err
	}
	configDao.Key = cdto.Key
	configDao.Val = cdto.Val
	configDao.Env = cdto.Env
	configDao.Mod = cdto.Mod
	configDao.Type = cdto.Type
	//configDao := &dao.CenterConfig{
	//	Id: cdto.Id,
	//	Key: cdto.Key,
	//	Val: cdto.Val,
	//	Env: cdto.Env,
	//	Mod: cdto.Mod,
	//}
	return configDao.Update()
}
func (c Config) Delete(cdto *dto.ConfigDelDto) error {
	configDao := &dao.CenterConfig{
		Id: cdto.Id,
	}
	err := configDao.Find()
	if err != nil {
		return err
	}
	return configDao.Delete()
}
