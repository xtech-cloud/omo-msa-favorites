package cache

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"omo.msa.favorite/tool"
	"sort"
	"time"
)

const (
	NoticeToFamily   = 0 //发送到小程序
	NoticeToDevice   = 1 //发送到设备
	NoticeToWebsite  = 2 //发送到网站网站
	NoticeToResident = 3 //对格桑码居民
	NoticeToRoute    = 4 //对应
)

type NoticeInfo struct {
	BaseInfo
	Type     uint8
	Status   MessageStatus
	Owner    string //该展览所属组织机构，scene
	Subtitle string
	Body     string
	Tags     []string
	Targets  []string //class, area等虚拟空间引用
}

func (mine *cacheContext) CreateNotice(info *NoticeInfo) error {
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
	db.Type = info.Type
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

func (mine *cacheContext) RemoveNotice(uid, operator string) error {
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

func (mine *cacheContext) GetNoticesByType(owner string, tp uint32) []*NoticeInfo {
	if owner == "" {
		return make([]*NoticeInfo, 0, 1)
	}
	array, err := nosql.GetNoticesByType(owner, tp)
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

func (mine *cacheContext) GetLatestNotice(owner string, tp uint32) *NoticeInfo {
	if owner == "" {
		return nil
	}
	array, err := nosql.GetNoticesByType(owner, tp)
	if err == nil && len(array) > 0 {
		sort.Slice(array, func(i, j int) bool {
			if array[i].CreatedTime.Unix() > array[j].CreatedTime.Unix() {
				return true
			} else {
				return false
			}
		})
		info := new(NoticeInfo)
		info.initInfo(array[0])
		return info
	}
	return nil
}

func (mine *cacheContext) GetNoticesByStatus(owner string, tp uint8, st MessageStatus) []*NoticeInfo {
	if owner == "" {
		return make([]*NoticeInfo, 0, 1)
	}
	array, err := nosql.GetNoticesByStatus(owner, tp, uint8(st))
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

func (mine *cacheContext) GetNoticesByTargets(owner string, array []string, st MessageStatus, tp uint8, page, num uint32) (uint32, uint32, []*NoticeInfo) {
	if array == nil || len(array) < 1 {
		return 0, 0, make([]*NoticeInfo, 0, 1)
	}
	all := make([]*NoticeInfo, 0, 10)
	var dbs []*nosql.Notice
	var er error
	if len(owner) < 1 {
		dbs, er = nosql.GetNoticesByTargets(uint8(st), tp, array)
	} else {
		dbs, er = nosql.GetNoticesByOTargets(owner, uint8(st), tp, array)
	}
	if er == nil {
		for _, db := range dbs {
			info := new(NoticeInfo)
			info.initInfo(db)
			all = append(all, info)
		}
	}
	if num < 1 {
		num = 10
	}
	if len(all) < 1 {
		return 0, 0, make([]*NoticeInfo, 0, 1)
	}
	return CheckPage(page, num, all)
}

func (mine *cacheContext) GetAllNoticesByTargets(owner string, st MessageStatus, tm uint64, array []string) []*NoticeInfo {
	if array == nil || len(array) < 1 {
		return make([]*NoticeInfo, 0, 1)
	}
	all := make([]*NoticeInfo, 0, 20)
	var dbs []*nosql.Notice
	var er error
	if len(owner) < 1 {
		dbs, er = nosql.GetNoticesByTargets(uint8(st), NoticeToFamily, array)
	} else {
		dbs, er = nosql.GetNoticesByOTargets(owner, uint8(st), NoticeToFamily, array)
	}
	if er == nil {
		var secs int64 = -3600 * 24 * 7
		for _, db := range dbs {
			if db.CreatedTime.Unix()-int64(tm) > secs {
				info := new(NoticeInfo)
				info.initInfo(db)
				all = append(all, info)
			}
		}
	}
	return all
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
	mine.Type = db.Type
	mine.Status = MessageStatus(db.Status)
	mine.Tags = db.Tags
	mine.Targets = db.Targets
	if mine.Targets == nil || len(mine.Targets) < 1 {
		mine.Targets = make([]string, 0, 15)
		for i := 0; i < 15; i += 1 {
			mine.Targets = append(mine.Targets, fmt.Sprintf("%d", i+1))
		}
		_ = nosql.UpdateNoticeTargets(mine.UID, mine.Operator, mine.Targets)
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
		return errors.New("the list of tags is nil")
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
		if st == MessageStatusAgree {
			_ = cacheCtx.updateRecord(mine.Owner, ObserveNotice, 1)
		}
	}
	return err
}

func (mine *NoticeInfo) UpdateTargets(operator string, list []string) error {
	if list == nil {
		return errors.New("the list of targets is nil")
	}
	err := nosql.UpdateNoticeTargets(mine.UID, operator, list)
	if err == nil {
		mine.Targets = list
		mine.Operator = operator
	}
	return err
}

func (mine *NoticeInfo) HadTargets(arr []string) bool {
	if mine.Targets == nil || len(mine.Targets) < 1 {
		return true
	}
	if tool.HasItem(mine.Targets, mine.Owner) {
		return true
	}
	if arr == nil || len(arr) < 1 {
		return false
	}
	for _, item := range arr {
		if tool.HasItem(mine.Targets, item) {
			return true
		}
	}
	return false
}
