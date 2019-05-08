package services

import (
	"etcdcc/apiserver/pkg/dto"
	"etcdcc/apiserver/pkg/models"
)

type Config struct {}
func (c Config) List(start int, limit int, q []string)(interface{},int64){
	configModel := &models.CenterConfig{}
	configModel.SetPageParams(start,limit)
	configModel.SetSearchCdt(q)
	configModel.SetSearchMap(dto.CONFIG_SEARCH)
	return configModel.List()
}
func (c Config) Create(cdto *dto.ConfigAddDto)(int64,error){
	configModel := &models.CenterConfig{
		Key: cdto.Key,
		Val: cdto.Val,
		Env: cdto.Env,
		Mod: cdto.Mod,
		Type: cdto.Type,
	}
	return configModel.Create()
}
func (c Config) Update(cdto *dto.ConfigEditDto) error {
	configModel := &models.CenterConfig{
		Id: cdto.Id,
	}
	err := configModel.Find()
	if err != nil{
		return err
	}
	configModel.Key = cdto.Key
	configModel.Val = cdto.Val
	configModel.Env = cdto.Env
	configModel.Mod = cdto.Mod
	configModel.Type = cdto.Type
	//configModel := &models.CenterConfig{
	//	Id: cdto.Id,
	//	Key: cdto.Key,
	//	Val: cdto.Val,
	//	Env: cdto.Env,
	//	Mod: cdto.Mod,
	//}
	return configModel.Update()
}
func (c Config) Delete(cdto *dto.ConfigDelDto) error {
	configModel := &models.CenterConfig{
		Id: cdto.Id,
	}
	err := configModel.Find()
	if err != nil{
		return err
	}
	return configModel.Delete()
}

