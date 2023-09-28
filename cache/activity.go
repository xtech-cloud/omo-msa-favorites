package cache

import (
	"errors"
	"github.com/micro/go-micro/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
	"omo.msa.favorite/tool"
	"time"
)

const (
	ActivityStatusDraft   uint8 = 0 //草稿
	ActivityStatusCheck   uint8 = 1 // 审核中
	ActivityStatusPending uint8 = 2 // 审核通过，待发布或者释放
	ActivityStatusRelease uint8 = 3 // 释放成功
	ActivityStatusPublish uint8 = 4 // 发布成功
	ActivityStatusAbandon uint8 = 5 // 活动废弃
)

const (
	ActivityTypeNormal   uint8 = 0 //
	ActivityTypeTemplate uint8 = 1
	ActivityTypeOpen     uint8 = 2
)

const (
	OptionAgree  OptionType = 1 //审核同意
	OptionRefuse OptionType = 2 //审核拒绝
	OptionSwitch OptionType = 3 //切换关联
)

const (
	AccessFree = 0 //自由访问
	AccessHide = 1 //子场景不可见
	AccessMust = 2 //子场景可见且必做
)

type OptionType uint8

type ActivityInfo struct {
	Type        uint8
	Status      uint8
	Access      uint8
	SubmitLimit uint8 // 参与人提交资源的数量限制
	ShowResult  uint8 // 是否展示结果
	BaseInfo

	Participant uint32 //参与者数量
	Owner       string
	Cover       string //
	Remark      string // 活动介绍
	Organizer   string // 组织者
	Require     string // 活动要求

	Template string //引用的活动模板

	Duration proxy.DurationInfo //持续时间
	Place    proxy.PlaceInfo
	Prize    *proxy.PrizeInfo

	Quotes  []string //引用的实体对象
	Assets  []string
	Tags    []string
	Targets []string         //具体的班级，场景等
	Opuses  []proxy.OpusInfo //获奖作品
}

func (mine *cacheContext) GetActivity(uid string) (*ActivityInfo, error) {
	db, err := nosql.GetActivity(uid)
	if err == nil {
		info := new(ActivityInfo)
		info.initInfo(db)
		return info, nil
	}
	return nil, err
}

func (mine *cacheContext) CreateActivity(info *ActivityInfo) error {
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
	db.Duration = info.Duration
	db.Place = info.Place
	db.Limit = info.SubmitLimit
	db.Status = info.Status
	db.Show = info.ShowResult
	db.Prize = info.Prize
	db.Access = 0
	db.Template = info.Template
	db.Opuses = make([]proxy.OpusInfo, 0, 1)
	db.Participant = 0
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
	//db.Persons = info.Persons
	//if db.Persons == nil {
	//	db.Persons = make([]proxy.PersonInfo, 0, 1)
	//}

	err := nosql.CreateActivity(db)
	if err == nil {
		info.initInfo(db)
	}
	return err
}

func (mine *cacheContext) GetActivitiesByOrganizer(uid string) []*ActivityInfo {
	array, err := nosql.GetActivityByOrganizer(uid)
	if err == nil {
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

func (mine *cacheContext) RemoveActivity(uid, operator string) error {
	err := nosql.RemoveActivity(uid, operator)
	return err
}

func (mine *cacheContext) GetActivityCount(owner string) uint32 {
	num := nosql.GetActivityCountByOwner(owner)
	return uint32(num)
}

func (mine *cacheContext) GetActivityCountByStatus(owner string, st uint32) uint32 {
	num := nosql.GetActivityCountByStatus(owner, uint8(st))
	return uint32(num)
}

func (mine *cacheContext) GetActivityParticipant(owner string) uint32 {
	var num uint32 = 0
	list, _ := nosql.GetActivitiesByOwner(owner)
	for _, activity := range list {
		num = num + activity.Participant
	}
	return num
}

func (mine *cacheContext) GetActivityRatio(uid string) uint32 {
	activity, err := mine.GetActivity(uid)
	if err != nil {
		logger.Warn("the activity not found that err = " + err.Error())
		return 0
	}
	min, max := activity.GetRatio()
	return min / max
}

func (mine *cacheContext) GetActivityTemplateCount(uid string) uint32 {
	list, err := nosql.GetActivitiesByTemplate(uid)
	if err != nil {
		logger.Warn("the activity not found that err = " + err.Error())
		return 0
	}
	return uint32(len(list))
}

func (mine *cacheContext) GetActivityCloneCount(owner string) uint32 {
	num := nosql.GetActivityCountByClone(owner)
	return uint32(num)
}

func (mine *cacheContext) GetActivityOpusCount(owner, uid string) uint32 {
	if uid != "" {

	} else {

	}

	return 0
}

func (mine *cacheContext) GetActivitiesByOwner(uid string, usable bool) []*ActivityInfo {
	array, err := nosql.GetActivitiesByOwner(uid)
	if err == nil {
		list := make([]*ActivityInfo, 0, len(array))
		for _, item := range array {
			if usable {
				if item.Status > ActivityStatusPending {
					info := new(ActivityInfo)
					info.initInfo(item)
					list = append(list, info)
				}
			} else {
				info := new(ActivityInfo)
				info.initInfo(item)
				list = append(list, info)
			}
		}
		return list
	}
	return make([]*ActivityInfo, 0, 1)
}

func (mine *cacheContext) GetActivitiesByTargets(owner string, array []string, st uint8, page, num uint32) (uint32, uint32, []*ActivityInfo) {
	if array == nil || len(array) < 1 {
		return 0, 0, make([]*ActivityInfo, 0, 1)
	}
	all := make([]*ActivityInfo, 0, 10)
	var dbs []*nosql.Activity
	var er error
	if len(owner) < 1 {
		dbs, er = nosql.GetActivitiesByTargets(st, array)
	} else {
		dbs, er = nosql.GetActivitiesByOTargets(owner, st, array)
	}
	if er == nil {
		for _, db := range dbs {
			info := new(ActivityInfo)
			info.initInfo(db)
			all = append(all, info)
		}
	}
	if num < 1 {
		num = 10
	}
	if len(all) < 1 {
		return 0, 0, make([]*ActivityInfo, 0, 1)
	}

	return CheckPage(page, num, all)
}

func (mine *cacheContext) GetAllActivitiesByTargets(owner string, st uint8, tm uint64, array []string) []*ActivityInfo {
	if array == nil || len(array) < 1 {
		return make([]*ActivityInfo, 0, 1)
	}
	all := make([]*ActivityInfo, 0, 20)
	var dbs []*nosql.Activity
	var er error
	if len(owner) < 1 {
		dbs, er = nosql.GetActivitiesByTargets(st, array)
	} else {
		dbs, er = nosql.GetActivitiesByOTargets(owner, st, array)
	}
	if er == nil {
		var secs int64 = -3600 * 24 * 7
		for _, db := range dbs {
			start := ParseTime(db.Date.Start)
			if start-int64(tm) > secs {
				info := new(ActivityInfo)
				info.initInfo(db)
				all = append(all, info)
			}
		}
	}
	return all
}

func (mine *cacheContext) GetAllActivitiesByStatus(owner string, state uint8) []*ActivityInfo {
	if len(owner) < 1 {
		return make([]*ActivityInfo, 0, 1)
	}
	all := make([]*ActivityInfo, 0, 10)
	db, _ := nosql.GetActivitiesByStatus(owner, state)
	if db != nil {
		for _, item := range db {
			info := new(ActivityInfo)
			info.initInfo(item)
			all = append(all, info)
		}
	}

	return all
}

//获取所有当前时间可用的活动
func (mine *cacheContext) GetAliveActivities(owner string) []*ActivityInfo {
	if len(owner) < 1 {
		return make([]*ActivityInfo, 0, 1)
	}
	dbs, _ := nosql.GetActivitiesByStatus(owner, ActivityStatusPublish)
	all := make([]*ActivityInfo, 0, len(dbs))
	for _, db := range dbs {
		if db.Duration.Stop > time.Now().Unix() {
			info := new(ActivityInfo)
			info.initInfo(db)
			all = append(all, info)
		}
	}
	return all
}

func (mine *cacheContext) GetAllActivitiesByType(owner string, tp uint8) []*ActivityInfo {
	if len(owner) < 1 {
		return make([]*ActivityInfo, 0, 1)
	}
	all := make([]*ActivityInfo, 0, 10)
	db, _ := nosql.GetActivitiesByType(owner, tp)
	if db != nil {
		for _, item := range db {
			info := new(ActivityInfo)
			info.initInfo(item)
			all = append(all, info)
		}
	}

	return all
}

func (mine *cacheContext) GetActivityTags() []string {
	list := make([]string, 0, 30)
	dbs, _ := nosql.GetActivities()
	if dbs != nil {
		for _, db := range dbs {
			for _, tag := range db.Tags {
				if !tool.HasItem(list, tag) {
					list = append(list, tag)
				}
			}
		}
	}

	return list
}

func (mine *cacheContext) GetActivitiesByStatus(owner string, states []uint8, page, num uint32) (uint32, uint32, []*ActivityInfo) {
	if len(owner) < 1 {
		return 0, 0, make([]*ActivityInfo, 0, 1)
	}
	all := make([]*ActivityInfo, 0, 10)
	db, _ := nosql.GetActivitiesByStates(owner, states)
	if db != nil {
		for _, item := range db {
			info := new(ActivityInfo)
			info.initInfo(item)
			all = append(all, info)
		}
	}

	if num < 1 {
		num = 10
	}
	if len(all) < 1 {
		return 0, 0, make([]*ActivityInfo, 0, 1)
	}
	return CheckPage(page, num, all)
}

func (mine *cacheContext) GetActivitiesByShow(owners []string, st uint8, page, num uint32) (uint32, uint32, []*ActivityInfo) {
	if owners == nil || len(owners) < 1 {
		return 0, 0, make([]*ActivityInfo, 0, 1)
	}
	all := make([]*ActivityInfo, 0, 10)
	db, _ := nosql.GetActivitiesByShow(owners, st)
	if db != nil {
		for _, item := range db {
			info := new(ActivityInfo)
			info.initInfo(item)
			all = append(all, info)
		}
	}

	if num < 1 {
		num = 10
	}
	if len(all) < 1 {
		return 0, 0, make([]*ActivityInfo, 0, 1)
	}
	return CheckPage(page, num, all)
}

func (mine *cacheContext) GetActivitiesByTemplate(owner, template string) []*ActivityInfo {
	if len(template) < 1 {
		return make([]*ActivityInfo, 0, 1)
	}
	all := make([]*ActivityInfo, 0, 10)
	if len(owner) < 1 {
		dbs, _ := nosql.GetActivitiesByTemplate(template)
		if dbs != nil {
			for _, item := range dbs {
				info := new(ActivityInfo)
				info.initInfo(item)
				all = append(all, info)
			}
		}
	} else {
		dbs, _ := nosql.GetActivitiesByOwnTemplate(owner, template)
		if dbs != nil {
			for _, item := range dbs {
				info := new(ActivityInfo)
				info.initInfo(item)
				all = append(all, info)
			}
		}
	}
	return all
}

func (mine *ActivityInfo) initInfo(db *nosql.Activity) {
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
	mine.Template = db.Template
	mine.Duration = db.Duration
	mine.Place = db.Place
	mine.Prize = db.Prize
	mine.ShowResult = db.Show
	mine.Access = db.Access
	mine.Assets = db.Assets
	mine.Tags = db.Tags
	mine.Targets = db.Targets
	mine.SubmitLimit = db.Limit
	mine.Participant = db.Participant
	mine.Status = db.Status

	if mine.Targets == nil || len(mine.Targets) > 16 {
		mine.Targets = make([]string, 0, 1)
		mine.Targets = append(mine.Targets, mine.Owner)
		_ = nosql.UpdateActivityTargets(mine.UID, mine.Operator, mine.Targets)
	}
	if mine.Assets == nil {
		mine.Assets = make([]string, 0, 1)
		_ = nosql.UpdateActivityAssets(mine.UID, mine.Operator, mine.Assets)
	}
	mine.Opuses = db.Opuses
	if mine.Opuses == nil {
		mine.Opuses = make([]proxy.OpusInfo, 0, 1)
		_ = nosql.UpdateActivityOpuses(mine.UID, mine.Operator, mine.Opuses)
	}
	if db.Date.Start != "" && db.Duration.Start < 1 {
		duration := proxy.DurationInfo{
			Start: proxy.DateToUTC(db.Date.Start, 0),
			Stop:  proxy.DateToUTC(db.Date.Stop, 1),
		}
		mine.Duration = duration
		_ = nosql.UpdateActivityDuration(mine.UID, duration)
	}
}

func (mine *ActivityInfo) UpdateBase(name, remark, require, operator string, date proxy.DateInfo, place proxy.PlaceInfo) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	dur := proxy.DurationInfo{Start: 0, Stop: 0}
	if len(date.Start) < 1 {
		dur.Start = mine.Duration.Start
	} else {
		dur.Start = proxy.DateToUTC(date.Start, 0)
	}
	if len(date.Stop) < 1 {
		dur.Stop = mine.Duration.Stop
	} else {
		dur.Stop = proxy.DateToUTC(date.Stop, 1)
	}
	if len(place.Location) < 1 {
		place.Location = mine.Place.Location
	}
	err := nosql.UpdateActivityBase(mine.UID, name, remark, require, operator, &dur, place)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Duration = dur
		mine.Place = place
		mine.Require = require
	}
	return err
}

func (mine *ActivityInfo) UpdateStatus(operator string, st uint8) error {
	err := nosql.UpdateActivityStatus(mine.UID, operator, st)
	if err == nil {
		mine.createHistory(operator, "", mine.Status, st)
		mine.Status = st
		mine.Operator = operator
		mine.UpdateTime = time.Now()
		if st == ActivityStatusPublish || st == ActivityStatusRelease {
			_ = cacheCtx.updateRecord(mine.Owner, ObserveActivity, 1)
		}
	}
	return err
}

// UpdatePrize 更新活动评奖要求
func (mine *ActivityInfo) UpdatePrize(operator, name, desc string, ranks []proxy.RankInfo) error {
	prize := proxy.PrizeInfo{Name: name, Desc: desc, Ranks: ranks}
	err := nosql.UpdateActivityPrize(mine.UID, operator, &prize)
	if err == nil {
		mine.Prize = &prize
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo) UpdateShowState(operator string, st uint8) error {
	err := nosql.UpdateActivityShowState(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo) UpdateAccess(operator string, st uint8) error {
	err := nosql.UpdateActivityAccess(mine.UID, operator, st)
	if err == nil {
		mine.Access = st
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo) UpdateOpuses(operator string, opuses []proxy.OpusInfo) error {
	err := nosql.UpdateActivityOpuses(mine.UID, operator, opuses)
	if err == nil {
		mine.Opuses = opuses
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo) UpdateTargets(operator string, list []string) error {
	if list == nil {
		return errors.New("the list of targets is nil")
	}
	err := nosql.UpdateActivityTargets(mine.UID, operator, list)
	if err == nil {
		mine.Targets = list
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo) UpdateTags(operator string, tags []string) error {
	if tags == nil {
		return errors.New("the list of tag is nil")
	}
	err := nosql.UpdateActivityTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo) UpdateAssetLimit(operator string, num uint8) error {
	err := nosql.UpdateActivityLimit(mine.UID, operator, num)
	if err == nil {
		mine.SubmitLimit = num
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo) UpdateCover(cover, operator string) error {
	err := nosql.UpdateActivityCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *ActivityInfo) UpdateAssets(operator string, list []string) error {
	if list == nil {
		return errors.New("the list of assets is nil")
	}
	err := nosql.UpdateActivityAssets(mine.UID, operator, list)
	if err == nil {
		mine.Assets = list
	}
	return err
}

func (mine *ActivityInfo) IsAlive() bool {
	if mine.Duration.Stop > time.Now().Unix() {
		return true
	}
	return false
}

func (mine *ActivityInfo) HadTargets(arr []string) bool {
	//全场景
	if mine.Targets == nil || len(mine.Targets) < 1 {
		return true
	}
	//全场景
	if tool.HasItem(mine.Targets, mine.Owner) {
		return true
	}
	if arr == nil || len(arr) < 1 {
		return true
	}
	for _, item := range arr {
		if tool.HasItem(mine.Targets, item) {
			return true
		}
	}
	return false
}

func (mine *ActivityInfo) GetHistories() ([]*nosql.History, error) {
	dbs, err := nosql.GetHistories(mine.UID)
	if err != nil {
		return nil, err
	}
	list := make([]*nosql.History, 0, len(dbs))
	for _, db := range dbs {
		if db.Option == uint8(OptionAgree) || db.Option == uint8(OptionRefuse) {
			list = append(list, db)
		}
	}
	return list, nil
}

func (mine *ActivityInfo) createHistory(operator, remark string, from, to uint8) {
	opt := OptionAgree
	if to > from {
		opt = OptionAgree
	} else {
		opt = OptionRefuse
	}

	_ = mine.insertHistory(operator, remark, string(from), string(to), opt)
}

func (mine *ActivityInfo) insertHistory(operator, remark, from, to string, opt OptionType) error {
	db := new(nosql.History)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRecordNextID()
	db.Creator = operator
	db.CreatedTime = time.Now()
	db.Parent = mine.UID
	db.From = from
	db.To = to
	db.Option = uint8(opt)
	db.Remark = remark
	return nosql.CreateHistory(db)
}

func (mine *ActivityInfo) GetRatio() (min, max uint32) {

	return 0, 0
}

func (mine *ActivityInfo) UpdateParticipant(num uint32) error {
	err := nosql.UpdateActivityParticipant(mine.UID, num)
	if err == nil {
		mine.Participant = num
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

//func (mine *ActivityInfo)HadPersonByEvent(uid string) bool {
//	for _, item := range mine.Persons {
//		if item.Event == uid {
//			return true
//		}
//	}
//	return false
//}
//
//func (mine *ActivityInfo)HadPerson(uid string) bool {
//	for _, item := range mine.Persons {
//		if item.Entity == uid {
//			return true
//		}
//	}
//	return false
//}
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

//func (mine *ActivityInfo)AppendPerson(uid, event string) error {
//	if mine.HadPersonByEvent(event) {
//		return nil
//	}
//	person := proxy.PersonInfo{Entity: uid, Event: event}
//	er := nosql.AppendActivityPerson(mine.UID, person)
//	if er == nil {
//		mine.Persons = append(mine.Persons, person)
//	}
//	return er
//}
//
//func (mine *ActivityInfo) SubtractPerson(uid string) error {
//	if !mine.HadPersonByEvent(uid) {
//		return nil
//	}
//	er := nosql.SubtractActivityPerson(mine.UID, uid)
//	if er == nil {
//		for i := 0;i < len(mine.Persons);i += 1 {
//			if mine.Persons[i].Entity == uid {
//				if i == len(mine.Persons) - 1 {
//					mine.Persons = append(mine.Persons[:i])
//				}else{
//					mine.Persons = append(mine.Persons[:i], mine.Persons[i+1:]...)
//				}
//
//				break
//			}
//		}
//	}
//	return er
//}
