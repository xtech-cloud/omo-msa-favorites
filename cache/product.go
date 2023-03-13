package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

type ProductInfo struct {
	Type   uint8
	Status uint8
	BaseInfo
	Key      string
	Entry    string
	Menus    string
	Templet  string
	Catalogs string
	Remark   string
	Revises  []string
	Effects  []*proxy.ProductEffect
}

func (mine *cacheContext) CreateProduct(name, key, entry, menus, templet, remark, operator string, tp uint8, revises []string, effects []*proxy.ProductEffect) (*ProductInfo, error) {
	db := new(nosql.Product)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetWordsNextID()
	db.CreatedTime = time.Now()
	db.Name = name
	db.Key = key
	db.Entry = entry
	db.Menus = menus
	db.Type = uint8(tp)
	db.Creator = operator
	db.Operator = operator
	db.Revises = revises
	db.Templet = templet
	db.Effects = effects
	db.Status = 1
	db.Catalogs = ""
	db.Remark = remark

	err := nosql.CreateProduct(db)
	if err == nil {
		info := new(ProductInfo)
		info.initInfo(db)
		return info, err
	}
	return nil, err
}

func (mine *cacheContext) GetProduct(uid string) (*ProductInfo, error) {
	db, err := nosql.GetProduct(uid)
	if err != nil {
		return nil, err
	}
	info := new(ProductInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetAllProducts() []*ProductInfo {
	dbs, err := nosql.GetProducts()
	if err != nil {
		return nil
	}
	arr := make([]*ProductInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(ProductInfo)
		info.initInfo(db)
		arr = append(arr, info)
	}
	return arr
}

func (mine *cacheContext) RemoveProduct(uid, operator string) error {
	return nosql.RemoveProduct(uid, operator)
}

func (mine *ProductInfo) initInfo(db *nosql.Product) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.Name = db.Name
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Type = db.Type
	mine.Status = db.Status
	mine.Menus = db.Menus
	mine.Key = db.Key
	mine.Entry = db.Entry
	mine.Remark = db.Remark
	mine.Catalogs = db.Catalogs
	mine.Templet = db.Templet
	mine.Revises = db.Revises
	mine.Effects = db.Effects
	if mine.Effects == nil {
		mine.Effects = make([]*proxy.ProductEffect, 0, 5)
	}
	if mine.Revises == nil {
		mine.Revises = make([]string, 0, 5)
	}
}

func (mine *ProductInfo) UpdateCatalogs(catalogs, operator string) error {
	err := nosql.UpdateProductCatalog(mine.UID, catalogs, operator)
	if err == nil {
		mine.Catalogs = catalogs
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ProductInfo) UpdateMenus(msg, operator string) error {
	err := nosql.UpdateProductMenus(mine.UID, msg, operator)
	if err == nil {
		mine.Menus = msg
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ProductInfo) UpdateRevises(operator string, list []string) error {
	err := nosql.UpdateProductRevises(mine.UID, operator, list)
	if err == nil {
		mine.Revises = list
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ProductInfo) UpdateTemplet(msg, operator string) error {
	err := nosql.UpdateProductTemplet(mine.UID, msg, operator)
	if err == nil {
		mine.Templet = msg
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ProductInfo) UpdateEntry(msg, operator string) error {
	err := nosql.UpdateProductEntry(mine.UID, msg, operator)
	if err == nil {
		mine.Entry = msg
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ProductInfo) UpdateBase(name, key, remark, operator string) error {
	err := nosql.UpdateProductBase(mine.UID, operator, name, key, remark)
	if err == nil {
		mine.Name = name
		mine.Key = key
		mine.Remark = remark
		mine.UpdateTime = time.Now()
	}
	return err
}
