package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

const (
	FavStatusDraft   uint8 = 0 //草稿
	FavStatusCheck   uint8 = 1 // 审核中
	FavStatusPending uint8 = 2 // 审核通过，待发布或者释放
	FavStatusPublish uint8 = 3 // 发布成功
)

type DisplayInfo struct {
	BaseInfo
	Status uint8
	Type   uint8  //
	Owner  string //该展览所属组织机构，scene
	Cover  string
	Remark string
	Origin string //布展数据来源，比如活动
	Banner string //标语
	Poster string //海报
	Meta   string //
	Tags   []string
	Keys   []string
	//Targets []*proxy.ShowingInfo //目标效果配置
}

func (mine *cacheContext) CreateDisplay(info *DisplayInfo) error {
	db := new(nosql.Display)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetDisplayNextID()
	db.CreatedTime = time.Now()
	db.Cover = info.Cover
	db.Name = info.Name
	db.Remark = info.Remark
	db.Owner = info.Owner
	db.Type = info.Type
	db.Origin = info.Origin
	db.Status = info.Status
	db.Creator = info.Creator
	db.Operator = info.Operator
	db.Banner = info.Banner
	db.Poster = info.Poster
	db.Tags = info.Tags
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	db.Keys = info.Keys
	if db.Keys == nil {
		db.Keys = make([]string, 0, 1)
	}
	//db.Targets = info.Targets
	//if db.Targets == nil {
	//	db.Targets = make([]*proxy.ShowingInfo, 0, 1)
	//}

	err := nosql.CreateDisplay(db)
	if err == nil {
		info.UID = db.UID.Hex()
		info.CreateTime = db.CreatedTime
		info.ID = db.ID
		info.UpdateTime = db.UpdatedTime
	}
	return err
}

func (mine *cacheContext) HadDisplayByName(owner, name string, tp uint8) bool {
	fav, err := nosql.GetDisplayByName(owner, name, tp)
	if err != nil {
		return false
	}
	if fav != nil {
		return true
	} else {
		return true
	}
}

func (mine *cacheContext) RemoveDisplay(uid, operator string) error {
	err := nosql.RemoveDisplay(uid, operator)
	return err
}

func (mine *cacheContext) GetDisplay(uid string) *DisplayInfo {
	db, err := nosql.GetDisplay(uid)
	if err == nil {
		info := new(DisplayInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetDisplayByOrigin(user, uid string) *DisplayInfo {
	db, err := nosql.GetDisplayByOrigin(user, uid)
	if err == nil {
		info := new(DisplayInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetDisplaysByOwner(uid string) []*DisplayInfo {
	array, err := nosql.GetDisplaysByOwner(uid)
	if err == nil {
		list := make([]*DisplayInfo, 0, len(array))
		for _, item := range array {
			info := new(DisplayInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetDisplaysByStatus(uid string, st uint8) []*DisplayInfo {
	array, err := nosql.GetDisplaysByStatus(uid, st)
	if err == nil {
		list := make([]*DisplayInfo, 0, len(array))
		for _, item := range array {
			info := new(DisplayInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetDisplaysByType(owner string, kind uint8) []*DisplayInfo {
	var array []*nosql.Display
	var err error
	if kind == 1 {
		array, err = nosql.GetDisplaysByType(kind)
	} else {
		array, err = nosql.GetDisplaysByOwnerTP(owner, kind)
	}
	if err == nil {
		list := make([]*DisplayInfo, 0, len(array))
		for _, item := range array {
			info := new(DisplayInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetDisplaysByList(array []string) []*DisplayInfo {
	if array == nil || len(array) < 1 {
		return make([]*DisplayInfo, 0, 1)
	}
	list := make([]*DisplayInfo, 0, 1)
	for _, s := range array {
		db, _ := nosql.GetDisplay(s)
		if db != nil {
			info := new(DisplayInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}

	return list
}

func (mine *cacheContext) GetDisplaysByTargets(owner string, array []string, page, num uint32) (uint32, uint32, []*DisplayInfo) {
	if array == nil || len(array) < 1 {
		return 0, 0, make([]*DisplayInfo, 0, 1)
	}
	all := make([]*DisplayInfo, 0, 10)
	for _, s := range array {
		db, _ := nosql.GetDisplaysByTarget(owner, s)
		if db != nil {
			for _, item := range db {
				info := new(DisplayInfo)
				info.initInfo(item)
				all = append(all, info)
			}
		}
	}

	if len(all) < 1 {
		return 0, 0, make([]*DisplayInfo, 0, 1)
	}
	max, pages, list := CheckPage(page, num, all)
	return max, pages, list.([]*DisplayInfo)
}

func (mine *DisplayInfo) initInfo(db *nosql.Display) {
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
	mine.Origin = db.Origin
	mine.Tags = db.Tags
	mine.Keys = db.Keys
	mine.Status = db.Status
	mine.Banner = db.Banner
	mine.Poster = db.Poster
	//mine.Targets = db.Targets
	if mine.Keys == nil {
		mine.Keys = make([]string, 0, 1)
	}
	//if mine.Targets == nil {
	//	mine.Targets = make([]*proxy.ShowingInfo, 0, 1)
	//	_ = nosql.UpdateDisplayTargets(mine.UID, mine.Operator, mine.Targets)
	//}
}

func (mine *DisplayInfo) GetKeys() []string {
	return mine.Keys
}

//func (mine *DisplayInfo) GetTargets() []*pb.ShowInfo {
//	list := make([]*pb.ShowInfo, 0, len(mine.Targets))
//	for _, item := range mine.Targets {
//		list = append(list, &pb.ShowInfo{Target: item.Target, Effect: item.Effect, Skin: item.Alignment, Slots: item.Slots})
//	}
//	return list
//}

func (mine *DisplayInfo) UpdateBase(name, remark, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	err := nosql.UpdateDisplayBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *DisplayInfo) UpdateMeta(operator, meta string) error {
	return nil
}

func (mine *DisplayInfo) UpdateTags(operator string, tags []string) error {
	if tags == nil {
		return errors.New("the list of tags is nil")
	}
	err := nosql.UpdateDisplayTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *DisplayInfo) UpdateCover(cover, operator string) error {
	err := nosql.UpdateDisplayCover(mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *DisplayInfo) UpdateStatus(st uint8, operator string) error {
	err := nosql.UpdateDisplayState(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
		if st == FavStatusPublish {
			_ = cacheCtx.updateRecord(mine.Owner, ObserveFav, 1)
		}
	}
	return err
}

func (mine *DisplayInfo) UpdateEntities(operator string, list []string) error {
	var err error
	if list == nil || len(list) < 1 {
		err = nosql.UpdateDisplayKeys(mine.UID, operator, make([]string, 0, 1))
		if err == nil {
			mine.Keys = make([]string, 0, 1)
			mine.Operator = operator
		}
	} else {
		err = nosql.UpdateDisplayKeys(mine.UID, operator, list)
		if err == nil {
			mine.Keys = list
			mine.Operator = operator
		}
	}
	return err
}

func (mine *DisplayInfo) HadKey(uid string) bool {
	for _, item := range mine.Keys {
		if item == uid {
			return true
		}
	}
	return false
}

//func (mine *DisplayInfo) HadTarget(uid string) bool {
//	for _, item := range mine.Targets {
//		if item.Target == uid {
//			return true
//		}
//	}
//	return false
//}

func (mine *DisplayInfo) AppendKey(uid string) error {
	if mine.HadKey(uid) {
		return nil
	}
	er := nosql.AppendDisplayKey(mine.UID, uid)
	if er == nil {
		mine.Keys = append(mine.Keys, uid)
	}
	return er
}

func (mine *DisplayInfo) SubtractKey(uid string) error {
	if !mine.HadKey(uid) {
		return nil
	}
	er := nosql.SubtractDisplayKey(mine.UID, uid)
	if er == nil {
		for i := 0; i < len(mine.Keys); i += 1 {
			if mine.Keys[i] == uid {
				if i == len(mine.Keys)-1 {
					mine.Keys = append(mine.Keys[:i])
				} else {
					mine.Keys = append(mine.Keys[:i], mine.Keys[i+1:]...)
				}
				break
			}
		}
	}
	return er
}

//func (mine *DisplayInfo) UpdateTarget(uid, effect, align, menu, operator string, slots []string) error {
//	if !mine.HadTarget(uid) {
//		return nil
//	}
//	if slots == nil {
//		slots = make([]string, 0, 1)
//	}
//	array := make([]*proxy.ShowingInfo, 0, len(mine.Targets))
//	for _, info := range mine.Targets {
//		if info.Target == uid {
//			info.Effect = effect
//			info.Alignment = align
//			info.Slots = slots
//			info.Menu = menu
//			info.UpdatedAt = time.Now()
//		}
//		array = append(array, info)
//	}
//	err := nosql.UpdateDisplayTargets(mine.UID, operator, array)
//	if err == nil {
//		mine.Targets = array
//	}
//	return err
//}

//func (mine *DisplayInfo) UpdateTargets(operator string, targets []string) error {
//	var array []*proxy.ShowingInfo
//	if targets != nil {
//		array = make([]*proxy.ShowingInfo, 0, len(mine.Targets))
//		for _, item := range targets {
//			info := new(proxy.ShowingInfo)
//			info.Target = item
//			info.Effect = ""
//			info.Alignment = ""
//			info.Slots = make([]string, 0, 1)
//			info.UpdatedAt = time.Now()
//			array = append(array, info)
//		}
//	} else {
//		array = make([]*proxy.ShowingInfo, 0, 1)
//	}
//	err := nosql.UpdateDisplayTargets(mine.UID, operator, array)
//	if err == nil {
//		mine.Targets = array
//	}
//	return err
//}

//func (mine *DisplayInfo) AppendTarget(show *proxy.ShowingInfo) error {
//	if mine.HadTarget(show.Target) {
//		_ = mine.SubtractTarget(show.Target)
//	}
//	er := nosql.AppendDisplayTarget(mine.UID, show)
//	if er == nil {
//		mine.Targets = append(mine.Targets, show)
//	}
//	return er
//}

//func (mine *DisplayInfo) AppendSimpleTarget(target string) error {
//	if target == "" {
//		return errors.New("the target is empty")
//	}
//	if mine.HadTarget(target) {
//		return nil
//	}
//	show := new(proxy.ShowingInfo)
//	show.Target = target
//	show.Effect = ""
//	show.Alignment = ""
//	show.Menu = ""
//	show.Slots = make([]string, 0, 1)
//	show.UpdatedAt = time.Now()
//	er := nosql.AppendDisplayTarget(mine.UID, show)
//	if er == nil {
//		mine.Targets = append(mine.Targets, show)
//	}
//	return er
//}

//func (mine *DisplayInfo) SubtractTarget(sn string) error {
//	if sn == "" {
//		return errors.New("the target is empty")
//	}
//	if !mine.HadTarget(sn) {
//		return nil
//	}
//	er := nosql.SubtractDisplayTarget(mine.UID, sn)
//	if er == nil {
//		for i := 0; i < len(mine.Targets); i += 1 {
//			if mine.Targets[i].Target == sn {
//				if i == len(mine.Targets)-1 {
//					mine.Targets = append(mine.Targets[:i])
//				} else {
//					mine.Targets = append(mine.Targets[:i], mine.Targets[i+1:]...)
//				}
//				break
//			}
//		}
//	}
//	return er
//}
