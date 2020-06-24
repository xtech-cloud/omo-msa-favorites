package cache

import (
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
)

type FavoriteInfo struct {
	BaseInfo
	Cover string
	Remark string
	entities []proxy.EntityInfo
}

func (mine *FavoriteInfo)initInfo(db *nosql.Favorite)  {
	mine.UID = db.UID.Hex()
	mine.Remark = db.Remark
	mine.ID = db.ID
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Cover = db.Cover
	mine.entities = db.Entities
}

func (mine *FavoriteInfo)GetEntities() []proxy.EntityInfo {
	return mine.entities
}

func (mine *FavoriteInfo)UpdateBase(name, remark string) error {
	err := nosql.UpdateFavoriteBase(mine.UID, name, remark)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
	}
	return err
}

func (mine *FavoriteInfo)UpdateCover(cover string) error {
	err := nosql.UpdateFavoriteCover(mine.UID, cover)
	if err == nil {
		mine.Cover = cover
	}
	return err
}

func (mine *FavoriteInfo)AppendEntity(uid, name string) error {
	err := nosql.AppendFavoriteEntity(mine.UID, proxy.EntityInfo{UID: uid, Name: name})
	if err == nil {
		mine.entities = append(mine.entities, proxy.EntityInfo{UID: uid, Name: name})
	}
	return err
}

func (mine *FavoriteInfo)SubtractEntity(uid string) error {
	err := nosql.SubtractFavoriteEntity(mine.UID, uid)
	if err == nil {
		for i := 0;i < len(mine.entities);i += 1 {
			if mine.entities[i].UID == uid {
				mine.entities = append(mine.entities[:i], mine.entities[i+1:]...)
				break
			}
		}
	}
	return err
}

