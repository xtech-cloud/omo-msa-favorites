package cache

import (
	"omo.msa.favorite/proxy/nosql"
)

type FavoriteInfo struct {
	BaseInfo
	Type uint8
	Owner string
	Cover string
	Remark string
	entities []string
}

func GetFavorite(uid string) (*OwnerInfo,*FavoriteInfo) {
	for i := 0;i < len(cacheCtx.boxes);i += 1{
		info := cacheCtx.boxes[i].GetFavorite(uid)
		if info != nil {
			return cacheCtx.boxes[i], info
		}
	}
	db,err := nosql.GetFavorite(uid)
	if err == nil{
		info:= new(FavoriteInfo)
		info.initInfo(db)
		owner := GetOwner(info.Owner)
		if owner != nil {
			owner.favorites = append(owner.favorites, info)
		}
		return owner,info
	}
	return nil,nil
}

func (mine *FavoriteInfo)initInfo(db *nosql.Favorite)  {
	mine.UID = db.UID.Hex()
	mine.Remark = db.Remark
	mine.ID = db.ID
	mine.Name = db.Name
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Cover = db.Cover
	mine.Type = db.Type
	mine.Owner = db.Owner
	mine.entities = db.Entities
}

func (mine *FavoriteInfo)GetEntities() []string {
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

func (mine *FavoriteInfo) UpdateEntities(operator string, list []string) error {
	err := nosql.UpdateFavoriteEntity(mine.UID, operator, list)
	if err == nil {
		mine.entities = list
	}
	return err
}

func (mine *FavoriteInfo)HadEntity(uid string) bool {
	for _, item := range mine.entities {
		if item == uid {
			return true
		}
	}
	return false
}

func (mine *FavoriteInfo)AppendEntity(uid string) error {
	if mine.HadEntity(uid) {
		return nil
	}
	er := nosql.AppendFavoriteEntity(mine.UID, uid)
	if er == nil {
		mine.entities = append(mine.entities, uid)
	}
	return er
}

func (mine *FavoriteInfo)SubtractEntity(uid string) error {
	if !mine.HadEntity(uid) {
		return nil
	}
	er := nosql.SubtractFavoriteEntity(mine.UID, uid)
	if er == nil {
		for i := 0;i < len(mine.entities);i += 1 {
			if mine.entities[i] == uid {
				mine.entities = append(mine.entities[:i], mine.entities[i+1:]...)
				break
			}
		}
	}
	return er
}

