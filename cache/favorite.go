package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

type FavoriteInfo struct {
	BaseInfo
	Status  uint8
	Type    uint8  //
	Owner   string //该展览所属用户等
	Cover   string
	Remark  string
	Meta    string //
	Tags    []string
	Keys    []string
}

func (mine *cacheContext) CreateFavorite(info *FavoriteInfo) error {
	db := new(nosql.Favorite)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetFavoriteNextID()
	db.CreatedTime = time.Now()
	db.Cover = info.Cover
	db.Name = info.Name
	db.Remark = info.Remark
	db.Owner = info.Owner
	db.Type = info.Type
	db.Status = info.Status
	db.Creator = info.Creator
	db.Operator = info.Operator
	db.Tags = info.Tags
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	db.Keys = info.Keys
	if db.Keys == nil {
		db.Keys = make([]string, 0, 1)
	}

	err := nosql.CreateFavorite(db)
	if err == nil {
		info.UID = db.UID.Hex()
		info.CreateTime = db.CreatedTime
		info.ID = db.ID
		info.UpdateTime = db.UpdatedTime
	}
	return err
}

func (mine *cacheContext) HadFavoriteByName(owner, name string, tp uint8) bool {
	fav, err := nosql.GetFavoriteByName(owner, name, tp)
	if err != nil {
		return false
	}
	if fav != nil {
		return true
	} else {
		return true
	}
}

func (mine *cacheContext) RemoveFavorite(uid, operator string) error {
	err := nosql.RemoveFavorite(uid, operator)
	return err
}

func (mine *cacheContext) GetFavorite(uid string) *FavoriteInfo {
	db, err := nosql.GetFavorite(uid)
	if err == nil {
		info := new(FavoriteInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetFavoritesByOwner(uid string) []*FavoriteInfo {
	array, err := nosql.GetFavoritesByOwner(uid)
	if err == nil {
		list := make([]*FavoriteInfo, 0, len(array))
		for _, item := range array {
			info := new(FavoriteInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetFavoritesByStatus(uid string, st uint8) []*FavoriteInfo {
	array, err := nosql.GetFavoritesByStatus(uid, st)
	if err == nil {
		list := make([]*FavoriteInfo, 0, len(array))
		for _, item := range array {
			info := new(FavoriteInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetFavoritesByType(owner string, kind uint8) []*FavoriteInfo {
	var array []*nosql.Favorite
	var err error
	if kind == 1 {
		array, err = nosql.GetFavoritesByType(kind)
	} else {
		array, err = nosql.GetFavoritesByOwnerTP(owner, kind)
	}
	if err == nil {
		list := make([]*FavoriteInfo, 0, len(array))
		for _, item := range array {
			info := new(FavoriteInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetFavoritesByList(array []string) []*FavoriteInfo {
	if array == nil || len(array) < 1 {
		return make([]*FavoriteInfo, 0, 1)
	}
	list := make([]*FavoriteInfo, 0, 1)
	for _, s := range array {
		db, _ := nosql.GetFavorite(s)
		if db != nil {
			info := new(FavoriteInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}

	return list
}

func (mine *FavoriteInfo) initInfo(db *nosql.Favorite) {
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
	mine.Tags = db.Tags
	mine.Keys = db.Keys
	mine.Status = db.Status
	if mine.Keys == nil {
		mine.Keys = make([]string, 0, 1)
	}
}

func (mine *FavoriteInfo) GetKeys() []string {
	return mine.Keys
}

func (mine *FavoriteInfo) UpdateBase(name, remark, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
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

func (mine *FavoriteInfo) UpdateMeta(operator, meta string) error {
	return nil
}

func (mine *FavoriteInfo) UpdateTags(operator string, tags []string) error {
	if tags == nil {
		return errors.New("the list of tags is nil")
	}
	err := nosql.UpdateFavoriteTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo) UpdateCover(cover, operator string) error {
	err := nosql.UpdateFavoriteCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo) UpdateStatus(st uint8, operator string) error {
	err := nosql.UpdateFavoriteState(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
		if st == FavStatusPublish {
			_ = cacheCtx.updateRecord(mine.Owner, ObserveFav, 1)
		}
	}
	return err
}

func (mine *FavoriteInfo) UpdateEntities(operator string, list []string) error {
	var err error
	if list == nil || len(list) < 1{
		err = nosql.UpdateFavoriteKeys(mine.UID, operator, make([]string, 0, 1))
		if err == nil {
			mine.Keys = make([]string, 0, 1)
			mine.Operator = operator
		}
	}else{
		err = nosql.UpdateFavoriteKeys(mine.UID, operator, list)
		if err == nil {
			mine.Keys = list
			mine.Operator = operator
		}
	}
	return err
}

func (mine *FavoriteInfo) HadKey(uid string) bool {
	for _, item := range mine.Keys {
		if item == uid {
			return true
		}
	}
	return false
}

func (mine *FavoriteInfo) AppendKey(uid string) error {
	if mine.HadKey(uid) {
		return nil
	}
	er := nosql.AppendFavoriteKey(mine.UID, uid)
	if er == nil {
		mine.Keys = append(mine.Keys, uid)
	}
	return er
}

func (mine *FavoriteInfo) SubtractKey(uid string) error {
	if !mine.HadKey(uid) {
		return nil
	}
	er := nosql.SubtractFavoriteKey(mine.UID, uid)
	if er == nil {
		for i := 0; i < len(mine.Keys); i += 1 {
			if mine.Keys[i] == uid {
				if i == len(mine.Keys)-1 {
					mine.Keys = append(mine.Keys[:i])
				} else {
					mine.Keys = append(mine.Keys[:i], mine.Keys[i+1:]...)
				}
				break
			}
		}
	}
	return er
}
