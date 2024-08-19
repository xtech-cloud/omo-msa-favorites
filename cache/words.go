package cache

import (
	"encoding/json"
	"github.com/micro/go-micro/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

const (
	WordsTypeTemplate WordsType = 0 //
	WordsTypeBless    WordsType = 1 //祝福
	WordsTypeImage    WordsType = 2
	WordsTypeOther    WordsType = 3 //建议
	WordsTypeComment  WordsType = 4 //评论
	WordsTypeMessage  WordsType = 5 //留言
)

type WordsType uint8

type WordsInfo struct {
	BaseInfo
	Type   WordsType
	Owner  string //场景
	Words  string
	Target string //
	Device string
	Weight int32
	Quote  string
	Count  uint32
	Remark string
	States []uint8
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
	db.Count = 0
	db.States = make([]uint8, 0, 1)
	db.States = append(db.States, 0)
	if db.Assets == nil {
		db.Assets = make([]string, 0, 1)
	}
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

func (mine *cacheContext) GetWordsByPage(owner string, tp WordsType, page, num uint32) []*WordsInfo {
	if page < 1 {
		page = 1
	}
	if num < 1 {
		num = 10
	}
	start := (page - 1) * num
	array, err := nosql.GetWordsByPage(owner, uint8(tp), int64(start), int64(num))
	if err == nil {
		list := make([]*WordsInfo, 0, len(array))
		for _, db := range array {
			info := new(WordsInfo)
			info.initInfo(db)
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
	now, er := time.ParseInLocation("2006-01-02", date, time.Local)
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
	first, er := time.ParseInLocation("2006-01-02", from, time.Local)
	if er != nil {
		logger.Error("GetWordsByDate .... from error that " + er.Error())
		return nil
	}
	second, er := time.ParseInLocation("2006-01-02", to, time.Local)
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

func (mine *cacheContext) GetWordsByTarget(scene, uid, from, to string) []*WordsInfo {
	first, er := time.ParseInLocation("2006-01-02", from, time.Local)
	if er != nil {
		logger.Error("GetWordsByTarget .... from error that " + er.Error())
		return nil
	}
	second, er := time.ParseInLocation("2006-01-02", to, time.Local)
	if er != nil {
		logger.Error("GetWordsByTarget ....to error that " + er.Error())
		return nil
	}
	array, err := nosql.GetWordsByTarget(scene, uid)
	if err == nil {
		list := make([]*WordsInfo, 0, len(array))
		for _, item := range array {
			created := item.CreatedTime.Unix()
			if created > first.Unix() && created < second.Unix() {
				info := new(WordsInfo)
				info.initInfo(item)
				list = append(list, info)
			}
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetWordsByTarget2(owner, uid string, page, num uint32) (uint32, uint32, []*WordsInfo) {
	array, err := nosql.GetWordsByTarget(owner, uid)
	if err == nil {
		list := make([]*WordsInfo, 0, len(array))
		for _, item := range array {
			info := new(WordsInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return CheckPage(page, num, list)
	}
	return 0, 0, nil
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

func (mine *cacheContext) GetWordsByUserType(uid string, tp uint32) []*WordsInfo {
	array, err := nosql.GetWordsByUserType(uid, uint8(tp))
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

func (mine *cacheContext) UpdateWordsStars(owner, data string) error {
	type Star struct {
		UID   string `json:"uid"`
		User  string `json:"user"`
		Count uint32 `json:"count"`
	}
	list := make([]Star, 0, 40)
	err := json.Unmarshal([]byte(data), list)
	if err != nil {
		return err
	}
	for _, star := range list {
		er := nosql.UpdateWordsCount(star.UID, star.User, star.Count)
		if err != nil {
			return er
		}
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
	mine.Quote = db.Quote
	mine.Device = db.Device
	mine.Target = db.Target
	mine.Assets = db.Assets
	mine.Weight = db.Weight
	mine.Count = db.Count
	mine.States = db.States
	mine.Remark = db.Remark
	mine.Type = WordsType(db.Type)
}

func (mine *WordsInfo) UpdateWeight(weight int32, operator string) error {
	err := nosql.UpdateWordsWeight(mine.UID, operator, weight)
	if err == nil {
		mine.Weight = weight
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *WordsInfo) UpdateCount(count uint32, operator string) error {
	if count < mine.Count {
		return nil
	}
	err := nosql.UpdateWordsCount(mine.UID, operator, count)
	if err == nil {
		mine.Count = count
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *WordsInfo) UpdateContent(words, operator string) error {
	err := nosql.UpdateWordsContent(mine.UID, words, operator)
	if err == nil {
		mine.Words = words
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *WordsInfo) UpdateBase(words, target, quote, operator string) error {
	err := nosql.UpdateWordsBase(mine.UID, words, target, quote, operator)
	if err == nil {
		mine.Words = words
		mine.Target = target
		mine.Quote = quote
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *WordsInfo) UpdateStates(arr []uint8, remark, operator string) error {
	err := nosql.UpdateWordsStates(mine.UID, remark, operator, arr)
	if err == nil {
		mine.States = arr
		mine.Remark = remark
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
