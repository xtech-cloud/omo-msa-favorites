package cache

import (
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
)

type FavoriteInfo struct {
	BaseInfo
	Type uint8
	Owner string
	Cover string
	Remark string
	entities []proxy.EntityInfo
}

func GetFavorite(uid string) *FavoriteInfo {
	for i := 0;i < len(cacheCtx.boxes);i += 1{
		info := cacheCtx.boxes[i].GetFavorite(uid)
		if info != nil {
			return info
		}
	}
	return nil
}

func (mine *FavoriteInfo)initInfo(db *nosql.Favorite)  {
	mine.UID = db.UID.Hex()
	mine.Remark = db.Remark
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Cover = db.Cover
	mine.Type = db.Type
	mine.Owner = db.Owner
	mine.entities = db.Entities
}

func (mine *FavoriteInfo)GetEntities() []proxy.EntityInfo {
	return mine.entities
}

func (mine *FavoriteInfo)UpdateBase(name, remark,operator string) error {
	if len(name) <1 {
		name = mine.Name
	}
	if len(remark) <1 {
		remark = mine.Remark
	}
	err := nosql.UpdateFavoriteBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo)UpdateCover(cover, operator string) error {
	err := nosql.UpdateFavoriteCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo) UpdateEntities(operator string, list []proxy.EntityInfo) error {
	err := nosql.UpdateFavoriteEntity(mine.UID, operator, list)
	if err == nil {
		mine.entities = list
	}
	return err
}

