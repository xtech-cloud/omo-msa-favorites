package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

const (
	OwnerTypePerson = 1
	OwnerTypeUnit = 0
)

type RepertoryInfo struct {
	UID string
}

func (mine *cacheContext)GetRepertory(owner string) (*RepertoryInfo,error) {
	db,err := nosql.GetRepertoryByOwner(owner)
	if err != nil {
		return nil, err
	}
	info := new(RepertoryInfo)
	info.initInfo(db)
	return info,nil
}

func (mine *cacheContext)CreateFavorite(info *FavoriteInfo, person bool) error {
	table := getFavoriteTable(person)
	db := new(nosql.Favorite)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetFavoriteNextID(table)
	db.CreatedTime = time.Now()
	db.Cover = info.Cover
	db.Name = info.Name
	db.Remark = info.Remark
	db.Owner = info.Owner
	db.Type = info.Type
	db.Origin = info.Origin
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

	err := nosql.CreateFavorite(table, db)
	if err == nil {
		info.UID = db.UID.Hex()
		info.CreateTime = db.CreatedTime
		info.ID = db.ID
		info.UpdateTime = db.UpdatedTime
	}
	return err
}

func (mine *cacheContext)HadFavoriteByName(owner, name string, tp uint8, person bool) bool {
	table := getFavoriteTable(person)
	fav, err := nosql.GetFavoriteByName(table, owner, name, tp)
	if err != nil {
		return false
	}
	if fav != nil {
		return true
	}else{
		return true
	}
}

func (mine *cacheContext)RemoveFavorite(uid, operator string, person bool) error {
	table := getFavoriteTable(person)
	err := nosql.RemoveFavorite(table, uid, operator)
	return err
}

func (mine *RepertoryInfo) initInfo(db *nosql.Repertory) {

}

func (mine *RepertoryInfo)createRepertory(uid string) (*nosql.Repertory,error) {
	db := new(nosql.Repertory)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRepertoryNextID()
	db.CreatedTime = time.Now()
	db.Owner = uid
	err := nosql.CreateRepertory(db)
	if err != nil {
		return nil, err
	}else{
		return db, nil
	}
}

func (mine *RepertoryInfo)HadBag(uid string) bool {

	return false
}

func (mine *RepertoryInfo) AppendAsset(uid string) error {
	if len(uid) < 1 {
		return errors.New("the asset uid is empty")
	}
	if mine.HadBag(uid) {
		return nil
	}
	er := nosql.AppendRepertoryBag(mine.UID, uid)
	return er
}

func (mine *RepertoryInfo) SubtractAsset(uid string) error {
	if !mine.HadBag(uid) {
		return nil
	}
	er := nosql.SubtractRepertoryBag(mine.UID, uid)
	return er
}

func (mine *RepertoryInfo)UpdateBags(list []string, operator string) error {
	er := nosql.UpdateRepertoryBags(mine.UID, operator, list)
	return er
}