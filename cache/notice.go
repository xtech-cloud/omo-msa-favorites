package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

type NoticeInfo struct {
	BaseInfo
	Status MessageStatus
	Owner    string //该展览所属组织机构，scene, class等
	Subtitle string
	Body     string
	Tags     []string
	Targets  []string
}

func (mine *cacheContext)CreateNotice(info *NoticeInfo) error {
	db := new(nosql.Notice)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetNoticeNextID()
	db.CreatedTime = time.Now()
	db.Subtitle = info.Subtitle
	db.Name = info.Name
	db.Body = info.Body
	db.Owner = info.Owner
	db.Creator = info.Creator
	db.Operator = info.Operator
	db.Tags = info.Tags
	db.Status = uint8(info.Status)
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}

	db.Targets = info.Targets
	if db.Targets == nil {
		db.Targets = make([]string, 0, 1)
	}

	err := nosql.CreateNotice(db)
	if err == nil {
		info.UID = db.UID.Hex()
		info.CreateTime = db.CreatedTime
		info.ID = db.ID
		info.UpdateTime = db.UpdatedTime
	}
	return err
}

func (mine *cacheContext) GetNotice(uid string) *NoticeInfo {
	db, err := nosql.GetNotice(uid)
	if err == nil {
		info := new(NoticeInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext)RemoveNotice(uid, operator string) error {
	err := nosql.RemoveNotice(uid, operator)
	return err
}

func (mine *cacheContext) GetNoticesByOwner(uid string) []*NoticeInfo {
	if uid == "" {
		return make([]*NoticeInfo, 0, 1)
	}
	array, err := nosql.GetNoticesByOwner(uid)
	if err == nil {
		list := make([]*NoticeInfo, 0, len(array))
		for _, item := range array {
			info := new(NoticeInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return make([]*NoticeInfo, 0, 1)
}

func (mine *cacheContext) GetNoticesByTargets(array []string) []*NoticeInfo {
	if array == nil || len(array) < 1 {
		return make([]*NoticeInfo, 0, 1)
	}
	list := make([]*NoticeInfo, 0, 10)
	for _, s := range array {
		db, _ := nosql.GetNoticesByTargets(s)
		if db != nil {
			for _, item := range db {
				info := new(NoticeInfo)
				info.initInfo(item)
				list = append(list, info)
			}
		}
	}

	return list
}

func (mine *cacheContext) GetNoticesByList(array []string) []*NoticeInfo {
	if array == nil || len(array) < 1 {
		return make([]*NoticeInfo, 0, 1)
	}
	list := make([]*NoticeInfo, 0, 1)
	for _, s := range array {
		db, _ := nosql.GetNotice(s)
		if db != nil {
			info := new(NoticeInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}

	return list
}

func (mine *NoticeInfo) initInfo(db *nosql.Notice) {
	mine.UID = db.UID.Hex()
	mine.Subtitle = db.Subtitle
	mine.ID = db.ID
	mine.Name = db.Name
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Body = db.Body
	mine.Owner = db.Owner
	mine.Status = MessageStatus(db.Status)
	mine.Tags = db.Tags
	mine.Targets = db.Targets
	if mine.Targets == nil {
		mine.Targets = make([]string, 0, 1)
	}

	if mine.Tags == nil {
		mine.Tags = make([]string, 0, 1)
	}
}

func (mine *NoticeInfo) UpdateBase(name, sub, body, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(sub) < 1 {
		sub = mine.Subtitle
	}
	if len(body) < 1 {
		body = mine.Body
	}
	err := nosql.UpdateNoticeBase(mine.UID, name, sub, body, operator)
	if err == nil {
		mine.Name = name
		mine.Subtitle = sub
		mine.Body = body
		mine.Operator = operator
	}
	return err
}

func (mine *NoticeInfo) UpdateTags(operator string, tags []string) error {
	if tags == nil {
		return errors.New("the list of target is nil")
	}
	err := nosql.UpdateNoticeTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *NoticeInfo) UpdateStatus(st MessageStatus, operator string) error {
	err := nosql.UpdateNoticeStatus(mine.UID, operator, uint8(st))
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *NoticeInfo) UpdateTargets(operator string, list []string) error {
	if list == nil {
		return errors.New("the list of target is nil")
	}
	err := nosql.UpdateNoticeTargets(mine.UID, operator, list)
	if err == nil {
		mine.Targets = list
		mine.Operator = operator
	}
	return err
}
