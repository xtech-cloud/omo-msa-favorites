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
)

type WordsType uint8

type WordsInfo struct {
	BaseInfo
	Type   WordsType
	Owner  string //该展览表所属用户等
	Words  string
	Target string //
	Asset  string
	Weight int32
	Quote  string
}

func (mine *cacheContext) CreateWords(words, owner, target, operator,quote,asset string, tp WordsType) (*WordsInfo, error) {
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
	mine.Target = db.Target
	mine.Asset = db.Asset
	mine.Weight = db.Weight
	mine.Type = WordsType(db.Type)
}
