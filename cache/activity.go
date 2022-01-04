package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

type ActivityInfo struct {
	BaseInfo
	Type     uint8
	Owner    string
	Cover    string
	Remark   string
	Organizer string
	Require string
	Date proxy.DateInfo
	Place proxy.PlaceInfo
	AssetLimit uint8
	Assets   []string
	Tags     []string
	Targets []string //班级，场景等
	Persons []proxy.PersonInfo
}

func (mine *cacheContext)GetActivity(uid string) *ActivityInfo {
	db,err := nosql.GetActivity(uid)
	if err == nil{
		info:= new(ActivityInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext)CreateActivity(info *ActivityInfo) error {
	db := new(nosql.Activity)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetActivityNextID()
	db.CreatedTime = time.Now()
	db.Cover = info.Cover
	db.Name = info.Name
	db.Remark = info.Remark
	db.Require = info.Require
	db.Owner = info.Owner
	db.Type = info.Type
	db.Organizer = info.Organizer
	db.Creator = info.Creator
	db.Operator = info.Operator
	db.Date = info.Date
	db.Place = info.Place
	db.Limit = info.AssetLimit
	db.Participants = make([]string, 0, 1)
	db.Tags = info.Tags
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	db.Assets = info.Assets
	if db.Assets == nil {
		db.Assets = make([]string, 0, 1)
	}
	db.Targets = info.Targets
	if db.Targets == nil {
		db.Targets = make([]string, 0, 1)
	}
	db.Persons = info.Persons
	if db.Persons == nil {
		db.Persons = make([]proxy.PersonInfo, 0, 1)
	}

	err := nosql.CreateActivity(db)
	if err == nil {
		info.UID = db.UID.Hex()
		info.CreateTime = db.CreatedTime
		info.ID = db.ID
	}
	return err
}

func (mine *cacheContext)GetActivityByOrganizer(uid string) []*ActivityInfo {
	array,err := nosql.GetActivityByOrganizer(uid)
	if err == nil{
		list := make([]*ActivityInfo, 0, len(array))
		for _, item := range array {
			info := new(ActivityInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return make([]*ActivityInfo, 0, 1)
}

func (mine *cacheContext)RemoveActivity(uid, operator string) error {
	err := nosql.RemoveActivity(uid, operator)
	return err
}

func (mine *cacheContext)GetActivitiesByOwner(uid string) []*ActivityInfo {
	array,err := nosql.GetActivitiesByOwner(uid)
	if err == nil{
		list := make([]*ActivityInfo, 0, len(array))
		for _, item := range array {
			info := new(ActivityInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return make([]*ActivityInfo, 0, 1)
}

func (mine *cacheContext) GetActivitiesByTargets(array []string, page, num uint32) (uint32, uint32, []*ActivityInfo) {
	if array == nil || len(array) < 1 {
		return 0, 0, make([]*ActivityInfo, 0, 1)
	}
	all := make([]*ActivityInfo, 0, 10)
	for _, s := range array {
		db, _ := nosql.GetActivitiesByTarget(s)
		if db != nil {
			for _, item := range db {
				info := new(ActivityInfo)
				info.initInfo(item)
				all = append(all, info)
			}
		}
	}

	if num < 1 {
		num = 10
	}
	if len(all) < 1 {
		return 0, 0, make([]*ActivityInfo, 0, 1)
	}
	max, pages, list := checkPage(page, num, all)
	return max, pages, list.([]*ActivityInfo)
}

func (mine *ActivityInfo)initInfo(db *nosql.Activity)  {
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
	mine.Require = db.Require
	mine.Organizer = db.Organizer
	mine.Date = db.Date
	mine.Place = db.Place
	mine.Assets = db.Assets
	mine.Tags = db.Tags
	mine.Targets = db.Targets
	mine.AssetLimit = db.Limit
	mine.Persons = db.Persons
	if mine.Targets == nil {
		mine.Targets = make([]string, 0 ,1)
	}
	if mine.Persons == nil {
		mine.Persons = make([]proxy.PersonInfo, 0, 5)
	}
	if db.Participants != nil && len(db.Participants) > 0 {
		for _, item := range db.Participants {
			mine.Persons = append(mine.Persons, proxy.PersonInfo{Entity: "", Event: item})
		}
	}
}

func (mine *ActivityInfo)GetEntities() []*pb.PairInfo {
	list := make([]*pb.PairInfo, 0, len(mine.Persons))
	for _, person := range mine.Persons {
		list = append(list, &pb.PairInfo{Key: person.Entity, Value: person.Event})
	}
	return list
}

func (mine *ActivityInfo)UpdateBase(name, remark, require,operator string, date proxy.DateInfo, place proxy.PlaceInfo) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}

	if len(date.Start) < 1 {
		date.Start = mine.Date.Start
	}
	if len(date.Stop) < 1 {
		date.Stop = mine.Date.Stop
	}
	if len(place.Location) < 1 {
		place.Location = mine.Place.Location
	}
	err := nosql.UpdateActivityBase(mine.UID, name, remark, require, operator, date, place)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Date = date
		mine.Place = place
		mine.Require = require
	}
	return err
}

func (mine *ActivityInfo)UpdateTargets(operator string, list []string) error {
	err := nosql.UpdateActivityTargets(mine.UID, operator, list)
	if err == nil {
		mine.Targets = list
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo)UpdateTags(operator string, tags []string) error {
	if tags == nil {
		return errors.New("the list of target is nil")
	}
	err := nosql.UpdateActivityTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo)UpdateAssetLimit(operator string, num uint8) error {
	err := nosql.UpdateActivityLimit(mine.UID, operator, num)
	if err == nil {
		mine.AssetLimit = num
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo)UpdateCover(cover, operator string) error {
	err := nosql.UpdateActivityCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo) UpdateAssets(operator string, list []string) error {
	if list == nil {
		return errors.New("the list of target is nil")
	}
	err := nosql.UpdateActivityAssets(mine.UID, operator, list)
	if err == nil {
		mine.Assets = list
	}
	return err
}

//func (mine *ActivityInfo)HadParticipant(uid string) bool {
//	for _, item := range mine.Participants {
//		if item == uid {
//			return true
//		}
//	}
//	return false
//}

func (mine *ActivityInfo)HadPersonByEvent(uid string) bool {
	for _, item := range mine.Persons {
		if item.Event == uid {
			return true
		}
	}
	return false
}

func (mine *ActivityInfo)HadPerson(uid string) bool {
	for _, item := range mine.Persons {
		if item.Entity == uid {
			return true
		}
	}
	return false
}

//func (mine *ActivityInfo)AppendParticipant(uid string) error {
//	if mine.HadParticipant(uid) {
//		return nil
//	}
//	er := nosql.AppendActivityParticipant(mine.UID, uid)
//	if er == nil {
//		mine.Participants = append(mine.Participants, uid)
//	}
//	return er
//}
//
//func (mine *ActivityInfo) SubtractParticipant(uid string) error {
//	if !mine.HadParticipant(uid) {
//		return nil
//	}
//	er := nosql.SubtractActivityParticipant(mine.UID, uid)
//	if er == nil {
//		for i := 0;i < len(mine.Participants);i += 1 {
//			if mine.Participants[i] == uid {
//				mine.Participants = append(mine.Participants[:i], mine.Participants[i+1:]...)
//				break
//			}
//		}
//	}
//	return er
//}

func (mine *ActivityInfo)AppendPerson(uid, event string) error {
	if mine.HadPersonByEvent(event) {
		return nil
	}
	person := proxy.PersonInfo{Entity: uid, Event: event}
	er := nosql.AppendActivityPerson(mine.UID, person)
	if er == nil {
		mine.Persons = append(mine.Persons, person)
	}
	return er
}

func (mine *ActivityInfo) SubtractPerson(uid string) error {
	if !mine.HadPersonByEvent(uid) {
		return nil
	}
	er := nosql.SubtractActivityPerson(mine.UID, uid)
	if er == nil {
		for i := 0;i < len(mine.Persons);i += 1 {
			if mine.Persons[i].Entity == uid {
				if i == len(mine.Persons) - 1 {
					mine.Persons = append(mine.Persons[:i])
				}else{
					mine.Persons = append(mine.Persons[:i], mine.Persons[i+1:]...)
				}

				break
			}
		}
	}
	return er
}

