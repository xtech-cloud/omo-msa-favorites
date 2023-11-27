package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

func (mine *cacheContext) createRecord(owner, begin string, tp uint8, count uint32) error {
	db := new(nosql.Record)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRecordNextID()
	db.CreatedTime = time.Now()
	db.Owner = owner
	db.Creator = ""
	db.Begin = begin
	db.Type = tp
	db.Count = count

	return nosql.CreateRecord(db)
}

func (mine *cacheContext) getRecord(owner string, tp uint8) *nosql.Record {
	db, _ := nosql.GetRecordsByType(owner, tp)
	return db
}

func (mine *cacheContext) GetRecordCount(owner string, tp uint8) uint32 {
	db, _ := nosql.GetRecordsByType(owner, tp)
	if db == nil {
		return 0
	}
	return db.Count
}

func (mine *cacheContext) updateRecord(owner string, tp uint8, offset uint32) error {
	db, _ := nosql.GetRecordsByType(owner, tp)
	if db != nil {
		return nosql.UpdateRecordCount(db.UID.Hex(), db.Count+offset)
	} else {
		return mine.createRecord(owner, "", tp, offset)
	}
}

func (mine *cacheContext) GetClickCount(owner, begin string) uint32 {
	db, _ := nosql.GetRecordsByDate(owner, begin, RecodeClick)
	if db == nil {
		return 0
	}
	return db.Count
}

func (mine *cacheContext) UpdateClickCount(owner, begin string, offset uint32) error {
	db, _ := nosql.GetRecordsByDate(owner, begin, RecodeClick)
	if db != nil {
		return nosql.UpdateRecordCount(db.UID.Hex(), db.Count+offset)
	} else {
		return mine.createRecord(owner, begin, RecodeClick, offset)
	}
}
