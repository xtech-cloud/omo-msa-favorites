package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

const (
	WordsTypeBless  WordsType = 0 //
	WordsTypePerson WordsType = 1 //
	WordsTypeImage  WordsType = 2
	WordsTypeOther  WordsType = 3
)

type WordsType uint8

type WordsInfo struct {
	BaseInfo
	Type   WordsType
	Owner  string //
	Words  string
	Target string //
	Asset  string
	Device string
	Weight int32
	Quote  string
}

func (mine *cacheContext) CreateWords(words, owner, target, sn, operator,quote,asset string, tp WordsType) (*WordsInfo, error) {
	db := new(nosql.Words)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetWordsNextID()
	db.CreatedTime = time.Now()
	db.Name = ""
	db.Words = words
	db.Owner = owner
	db.Type = uint8(tp)
	db.Creator = operator
	db.Operator = operator
	db.Target = target
	db.Quote = quote
	db.Asset = asset
	db.Device = sn
	db.Weight = 0
	err := nosql.CreateWords(db)
	if err == nil {
		info := new(WordsInfo)
		info.initInfo(db)
		return info, err
	}
	return nil, err
}

func (mine *cacheContext) RemoveWords(uid, operator string) error {
	err := nosql.RemoveWords(uid, operator)
	return err
}

func (mine *cacheContext) GetWords(uid string) *WordsInfo {
	db, err := nosql.GetWords(uid)
	if err == nil {
		info := new(WordsInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetWordsByToday(owner, user, target string) *WordsInfo {
	dbs, err := nosql.GetWordsByCreator(owner, user, target, 1)
	if err == nil {
		now := time.Now()
		for _, db := range dbs {
			if db.CreatedTime.Year() == now.Year() && db.CreatedTime.Month() == now.Month() && db.CreatedTime.Day() == now.Day() {
				info := new(WordsInfo)
				info.initInfo(db)
				return info
			}
		}
	}
	return nil
}

func (mine *cacheContext) GetWordsByOwnerTP(owner string, tp WordsType) []*WordsInfo {
	array, err := nosql.GetWordsByOwnerType(owner, uint8(tp))
	if err == nil {
		list := make([]*WordsInfo, 0, len(array))
		for _, item := range array {
			info := new(WordsInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetWordsByOwner(uid string) []*WordsInfo {
	array, err := nosql.GetWordsByOwner(uid)
	if err == nil {
		list := make([]*WordsInfo, 0, len(array))
		for _, item := range array {
			info := new(WordsInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetWordsByQuote(owner, quote string) []*WordsInfo {
	array, err := nosql.GetWordsByQuote(owner, quote)
	if err == nil {
		list := make([]*WordsInfo, 0, len(array))
		for _, item := range array {
			info := new(WordsInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetWordsByTarget(uid string) []*WordsInfo {
	array, err := nosql.GetWordsByTarget(uid)
	if err == nil {
		list := make([]*WordsInfo, 0, len(array))
		for _, item := range array {
			info := new(WordsInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetWordsByUser(uid string) []*WordsInfo {
	array, err := nosql.GetWordsByUser(uid)
	if err == nil {
		list := make([]*WordsInfo, 0, len(array))
		for _, item := range array {
			info := new(WordsInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *WordsInfo) initInfo(db *nosql.Words) {
	mine.UID = db.UID.Hex()
	mine.Words = db.Words
	mine.ID = db.ID
	mine.Name = db.Name
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Owner = db.Owner
	mine.Device = db.Device
	mine.Target = db.Target
	mine.Asset = db.Asset
	mine.Weight = db.Weight
	mine.Type = WordsType(db.Type)
}

func (mine *WordsInfo) UpdateWeight(weight int32, operator string) error {
	err := nosql.UpdateWordsState(mine.UID, operator, weight)
	if err == nil {
		mine.Weight = weight
	}
	return err
}
