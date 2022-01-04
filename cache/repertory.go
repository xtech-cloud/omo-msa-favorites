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