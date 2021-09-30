package cache

import (
	"omo.msa.favorite/proxy/nosql"
)

type FavoriteInfo struct {
	BaseInfo
	Type   uint8
	Owner  string
	Cover  string
	Remark string
	Origin string
	Meta string
	table string
	Tags   []string
	Keys   []string
}

func getFavoriteTable(person bool) string {
	if person {
		return nosql.TableFavorite+ "_person"
	}else{
		return nosql.TableFavorite+"_scene"
	}
}

func (mine *cacheContext)GetFavorite(uid string, person bool) *FavoriteInfo {
	table := getFavoriteTable(person)
	db,err := nosql.GetFavorite(table, uid)
	if err == nil{
		info:= new(FavoriteInfo)
		info.initInfo(db, table)
		return info
	}
	return nil
}

func (mine *cacheContext)GetFavoriteByOrigin(user, uid string, person bool) *FavoriteInfo {
	table := getFavoriteTable(person)
	db,err := nosql.GetFavoriteByOrigin(table, user, uid)
	if err == nil{
		info:= new(FavoriteInfo)
		info.initInfo(db, table)
		return info
	}
	return nil
}

func (mine *cacheContext)GetFavoritesByOwner(uid string, person bool) []*FavoriteInfo {
	table := getFavoriteTable(person)
	array,err := nosql.GetFavoritesByOwner(table, uid)
	if err == nil{
		list := make([]*FavoriteInfo, 0, len(array))
		for _, item := range array {
			info := new(FavoriteInfo)
			info.initInfo(item, table)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext)GetFavoritesByType(owner string, kind uint8, person bool) []*FavoriteInfo {
	table := getFavoriteTable(person)
	var array []*nosql.Favorite
	var err error
	if kind == 1 {
		array,err = nosql.GetFavoritesByType(table, kind)
	}else{
		array,err = nosql.GetFavoritesByOwnerTP(table, owner, kind)
	}
	if err == nil{
		list := make([]*FavoriteInfo, 0, len(array))
		for _, item := range array {
			info := new(FavoriteInfo)
			info.initInfo(item, table)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext)GetFavoritesByList(person bool, array []string) []*FavoriteInfo {
	if array == nil || len(array) < 1 {
		return make([]*FavoriteInfo, 0, 1)
	}
	table := getFavoriteTable(person)
	list := make([]*FavoriteInfo, 0, 1)
	for _, s := range array {
		db,_ := nosql.GetFavorite(table, s)
		if db != nil {
			info := new(FavoriteInfo)
			info.initInfo(db, table)
			list = append(list, info)
		}
	}

	return list
}

func (mine *FavoriteInfo)initInfo(db *nosql.Favorite, table string)  {
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
	mine.Origin = db.Origin
	mine.Tags = db.Tags
	mine.Keys = db.Keys
	mine.table = table
}

func (mine *FavoriteInfo) GetKeys() []string {
	return mine.Keys
}

func (mine *FavoriteInfo)UpdateBase(name, remark,operator string) error {
	if len(name) <1 {
		name = mine.Name
	}
	if len(remark) <1 {
		remark = mine.Remark
	}
	err := nosql.UpdateFavoriteBase(mine.table, mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo)UpdateTags(operator string, tags []string) error {
	err := nosql.UpdateFavoriteTags(mine.table, mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo)UpdateCover(cover, operator string) error {
	err := nosql.UpdateFavoriteCover(mine.table, mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo) UpdateEntities(operator string, list []string) error {
	err := nosql.UpdateFavoriteEntity(mine.table, mine.UID, operator, list)
	if err == nil {
		mine.Keys = list
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo)HadEntity(uid string) bool {
	for _, item := range mine.Keys {
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
	er := nosql.AppendFavoriteEntity(mine.table, mine.UID, uid)
	if er == nil {
		mine.Keys = append(mine.Keys, uid)
	}
	return er
}

func (mine *FavoriteInfo)SubtractEntity(uid string) error {
	if !mine.HadEntity(uid) {
		return nil
	}
	er := nosql.SubtractFavoriteEntity(mine.table, mine.UID, uid)
	if er == nil {
		for i := 0;i < len(mine.Keys);i += 1 {
			if mine.Keys[i] == uid {
				mine.Keys = append(mine.Keys[:i], mine.Keys[i+1:]...)
				break
			}
		}
	}
	return er
}

