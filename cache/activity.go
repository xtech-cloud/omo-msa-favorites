package cache

import (
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
	Assets   []string
	Tags     []string
	Participants []string
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
	db.Tags = info.Tags
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	db.Assets = info.Assets
	if db.Assets == nil {
		db.Assets = make([]string, 0, 1)
	}
	db.Participants = info.Participants
	if db.Participants == nil {
		db.Participants = make([]string, 0, 1)
	}

	return nosql.CreateActivity(db)
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
	return nil
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
	return nil
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
	mine.Participants = db.Participants
}

func (mine *ActivityInfo)GetEntities() []string {
	return mine.Participants
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

func (mine *ActivityInfo)UpdateTags(operator string, tags []string) error {
	err := nosql.UpdateActivityTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
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
	err := nosql.UpdateActivityAssets(mine.UID, operator, list)
	if err == nil {
		mine.Assets = list
	}
	return err
}

func (mine *ActivityInfo)HadParticipant(uid string) bool {
	for _, item := range mine.Participants {
		if item == uid {
			return true
		}
	}
	return false
}

func (mine *ActivityInfo)AppendParticipant(uid string) error {
	if mine.HadParticipant(uid) {
		return nil
	}
	er := nosql.AppendActivityParticipant(mine.UID, uid)
	if er == nil {
		mine.Participants = append(mine.Participants, uid)
	}
	return er
}

func (mine *ActivityInfo) SubtractParticipant(uid string) error {
	if !mine.HadParticipant(uid) {
		return nil
	}
	er := nosql.SubtractActivityParticipant(mine.UID, uid)
	if er == nil {
		for i := 0;i < len(mine.Participants);i += 1 {
			if mine.Participants[i] == uid {
				mine.Participants = append(mine.Participants[:i], mine.Participants[i+1:]...)
				break
			}
		}
	}
	return er
}

