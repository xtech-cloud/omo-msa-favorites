package cache

import (
	"github.com/micro/go-micro/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

const (
	WordsTypeTemplate WordsType = 0 //
	WordsTypeBless    WordsType = 1 //
	WordsTypeImage    WordsType = 2
	WordsTypeOther    WordsType = 3
)

type WordsType uint8

type WordsInfo struct {
	BaseInfo
	Type   WordsType
	Owner  string //
	Words  string
	Target string //
	Device string
	Weight int32
	Quote  string
	Assets []string
}

func (mine *cacheContext) CreateWords(words, owner, target, sn, operator, quote string, assets []string, tp WordsType) (*WordsInfo, error) {
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
	db.Assets = assets
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

func (mine *cacheContext) GetWordsByToday(owner, user, device string) *WordsInfo {
	dbs, err := nosql.GetWordsByCreator(owner, user, device, 1)
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

func (mine *cacheContext) GetWordsByDate(owner, date string, before bool) []*WordsInfo {
	now, er := time.Parse("2006-01-02", date)
	if er != nil {
		logger.Error("GetWordsByDate ...." + er.Error())
		return nil
	}
	var array []*nosql.Words
	var err error
	if before {
		array, err = nosql.GetWordsBeforeDate(owner, now)
	} else {
		array, err = nosql.GetWordsAfterDate(owner, now)
	}

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

func (mine *cacheContext) GetWordsByBetweenDate(owner, from, to string, tp WordsType) []*WordsInfo {
	first, er := time.Parse("2006-01-02", from)
	if er != nil {
		logger.Error("GetWordsByDate .... from error that " + er.Error())
		return nil
	}
	second, er := time.Parse("2006-01-02", to)
	if er != nil {
		logger.Error("GetWordsByDate ....to error that " + er.Error())
		return nil
	}
	array, err := nosql.GetWordsBetweenDate(owner, first, second, uint8(tp))
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

func (mine *cacheContext) GetWordsCountByDevice(device string) (uint32, error) {
	count, err := nosql.GetWordsCountByDevice(device)
	return uint32(count), err
}

func (mine *cacheContext) GetWordsCountByScene(scene string) (uint32, error) {
	count, err := nosql.GetWordsCountByScene(scene)
	return uint32(count), err
}

func (mine *cacheContext) GetWordsCountByToday(device string) (uint32, error) {
	dbs, err := nosql.GetWordsByDevice(device)
	if err != nil {
		return 0, err
	}
	var count uint32 = 0
	now := time.Now()
	for _, db := range dbs {
		if db.CreatedTime.Format("2006-01-02") == now.Format("2006-01-02") {
			count = count + 1
		}
	}
	return count, nil
}

func (mine *cacheContext) GetWordsCountBetween(owner, from, to string, tp WordsType) (uint32, error) {
	dbs := mine.GetWordsByBetweenDate(owner, from, to, tp)
	return uint32(len(dbs)), nil
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
	mine.Quote = db.Quote
	mine.Device = db.Device
	mine.Target = db.Target
	mine.Assets = db.Assets
	mine.Weight = db.Weight
	mine.Type = WordsType(db.Type)
}

func (mine *WordsInfo) UpdateWeight(weight int32, operator string) error {
	err := nosql.UpdateWordsState(mine.UID, operator, weight)
	if err == nil {
		mine.Weight = weight
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *WordsInfo) UpdateBase(words, operator string) error {
	err := nosql.UpdateWordsBase(mine.UID, words, operator)
	if err == nil {
		mine.Words = words
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *WordsInfo) UpdateAssets(assets []string, operator string) error {
	err := nosql.UpdateWordsAssets(mine.UID, operator, assets)
	if err == nil {
		mine.Assets = assets
		mine.UpdateTime = time.Now()
	}
	return err
}
