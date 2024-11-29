package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
	"strconv"
	"time"
)

const (
	FavStatusDraft   uint8 = 0 //草稿
	FavStatusPending uint8 = 1 // 审核中，待发布或者释放
	FavStatusPublish uint8 = 2 // 发布成功
)

const (
	DisplayAccessRead  = 0
	DisplayAccessClose = 1
	DisplayAccessRW    = 2
	DisplayAccessWrite = 3
)

const (
	ContentOptionStable   ContentOption = 0
	ContentOptionAppend   ContentOption = 1
	ContentOptionSubtract ContentOption = 2
)

const DefaultOwner = "root"

type ContentOption uint8

type DisplayInfo struct {
	BaseInfo
	Status   uint8
	Type     uint8  //
	Access   uint8  //=0只可访问，=1不可访问也不可操作,=2可访问可操作，=3只可操作
	Owner    string //该展览所属组织机构，scene
	Cover    string
	Remark   string
	Origin   string //布展数据来源，比如活动,荣誉榜
	Banner   string //标语
	Poster   string //海报
	Meta     string //
	Tags     []string
	Targets  []uint32               //目标场景类型
	Contents []proxy.DisplayContent //正式的内容
	Pending  []proxy.DisplayContent //待审的内容
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
	db.Access = DisplayAccessClose
	db.Scenes = make([]uint32, 0, 1)
	if db.Owner == "" {
		db.Owner = DefaultOwner
	}
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	db.Contents = info.Contents
	if db.Contents == nil {
		db.Contents = make([]proxy.DisplayContent, 0, 1)
	}
	db.Pending = info.Pending
	if db.Pending == nil {
		db.Pending = make([]proxy.DisplayContent, 0, 1)
	}

	err := nosql.CreateDisplay(db)
	if err == nil {
		info.UID = db.UID.Hex()
		info.CreateTime = db.CreatedTime
		info.ID = db.ID
		info.Access = DisplayAccessClose
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
	if uid == "" {
		return make([]*DisplayInfo, 0, 1)
	}
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

func (mine *cacheContext) GetDisplaysByContent(owner, uid string) []*DisplayInfo {
	if uid == "" {
		return make([]*DisplayInfo, 0, 1)
	}
	array := make([]*nosql.Display, 0, 10)
	var err error
	if owner == "" {
		array, err = nosql.GetDisplaysByContent2(uid)
	} else {
		array, err = nosql.GetDisplaysByContent(owner, uid)
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

func hadTarget(arr []uint32, tp uint32) bool {
	if len(arr) == 0 {
		return true
	}
	for _, u := range arr {
		if u == tp {
			return true
		}
	}
	return false
}

func (mine *cacheContext) GetAvailableDisplays(arr []string, tp uint32) []*DisplayInfo {
	list := make([]*DisplayInfo, 0, 100)
	for _, uid := range arr {
		dbs, err := nosql.GetDisplaysByOwner(uid)
		if err == nil {
			for _, item := range dbs {
				if item.Status == FavStatusPublish && hadTarget(item.Scenes, tp) && (item.Access == DisplayAccessRW || item.Access == DisplayAccessWrite) {
					info := new(DisplayInfo)
					info.initInfo(item)
					list = append(list, info)
				}
			}
		}
	}

	return list
}

func (mine *cacheContext) GetVisibleDisplays(arr []string, tp uint32) []*DisplayInfo {
	list := make([]*DisplayInfo, 0, 100)
	for _, uid := range arr {
		dbs, err := nosql.GetDisplaysByOwner(uid)
		if err == nil {
			for _, item := range dbs {
				if hadTarget(item.Scenes, tp) && (item.Access == DisplayAccessRW || item.Access == DisplayAccessRead) {
					info := new(DisplayInfo)
					info.initInfo(item)
					list = append(list, info)
				}
			}
		}
	}

	return list
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
	var check = false
	if owner == "" {
		owner = DefaultOwner
		check = true
	}
	array, err = nosql.GetDisplaysByOwnerTP(owner, kind)
	if err == nil {
		list := make([]*DisplayInfo, 0, len(array))
		for _, item := range array {
			if check {
				if item.Access != DisplayAccessClose {
					info := new(DisplayInfo)
					info.initInfo(item)
					list = append(list, info)
				}
			} else {
				info := new(DisplayInfo)
				info.initInfo(item)
				list = append(list, info)
			}

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

func (mine *cacheContext) GetDisplayPages(page, number uint32) (uint32, uint32, []*DisplayInfo) {
	if page < 1 {
		page = 1
	}
	if number < 1 {
		number = 10
	}
	start := (page - 1) * number
	array, err := nosql.GetDisplaysByPage(FavStatusPublish, int64(start), int64(number))
	total := nosql.GetDisplaysCount(FavStatusPublish)
	pages := math.Ceil(float64(total) / float64(number))
	if err == nil {
		list := make([]*DisplayInfo, 0, len(array))
		for _, item := range array {
			info := new(DisplayInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return uint32(total), uint32(pages), list
	}
	return 0, 0, make([]*DisplayInfo, 0, 1)
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
	return CheckPage(page, num, all)
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

	mine.Status = db.Status
	mine.Banner = db.Banner
	mine.Poster = db.Poster
	mine.Access = db.Access
	mine.Targets = db.Scenes
	mine.Contents = db.Contents
	if mine.Contents == nil {
		mine.Contents = make([]proxy.DisplayContent, 0, 1)
	}
	mine.Pending = db.Pending
	if mine.Pending == nil {
		mine.Pending = make([]proxy.DisplayContent, 0, 1)
	}
	//if len(mine.Contents) < 1 && len(db.Keys) > 0 {
	//	arr := make([]proxy.DisplayContent, 0, len(db.Keys))
	//	for _, key := range db.Keys {
	//		arr = append(arr, proxy.DisplayContent{UID: key, Events: make([]string, 0, 1), Assets: make([]string, 0, 1)})
	//	}
	//	_ = nosql.UpdateDisplayContents(mine.UID, "", arr)
	//	mine.Contents = arr
	//}
	if mine.Owner == "" {
		_ = nosql.UpdateDisplayOwner(mine.UID, DefaultOwner, mine.Operator)
		mine.Owner = DefaultOwner
	}
	//if len(mine.Contents) > 0 {
	//	mine.Status = FavStatusPublish
	//	_ = nosql.UpdateDisplayState(mine.UID, mine.Operator, FavStatusPublish)
	//}
}

func (mine *DisplayInfo) GetContents() []*pb.DisplayContent {
	arr := make([]*pb.DisplayContent, 0, len(mine.Contents))
	for _, content := range mine.Contents {
		arr = append(arr, &pb.DisplayContent{Uid: content.UID,
			Remark: content.Remark,
			Events: content.Events,
			Assets: content.Assets})
	}
	return arr
}

func (mine *DisplayInfo) GetPending() []*pb.DisplayContent {
	arr := make([]*pb.DisplayContent, 0, len(mine.Pending))
	for _, content := range mine.Pending {
		arr = append(arr, &pb.DisplayContent{Uid: content.UID,
			Remark:    content.Remark,
			Option:    content.Option,
			Submitter: content.Submitter,
			Reviewer:  content.Reviewer,
			Events:    content.Events,
			Assets:    content.Assets})
	}
	return arr
}

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

func (mine *DisplayInfo) UpdateBanner(cover, operator string) error {
	err := nosql.UpdateDisplayBanner(mine.UID, cover, operator)
	if err == nil {
		mine.Banner = cover
		mine.Operator = operator
	}
	return err
}

func (mine *DisplayInfo) UpdatePoster(cover, operator string) error {
	err := nosql.UpdateDisplayPoster(mine.UID, cover, operator)
	if err == nil {
		mine.Poster = cover
		mine.Operator = operator
	}
	return err
}

func (mine *DisplayInfo) UpdateAccess(acc uint8, operator string) error {
	err := nosql.UpdateDisplayAccess(mine.UID, operator, acc)
	if err == nil {
		mine.Access = acc
		mine.Operator = operator
	}
	return err
}

func (mine *DisplayInfo) UpdateType(tp uint8, operator string) error {
	err := nosql.UpdateDisplayType(mine.UID, operator, tp)
	if err == nil {
		mine.Type = tp
		mine.Operator = operator
	}
	return err
}

func (mine *DisplayInfo) UpdateStatus(st uint8, operator, remark string) error {
	err := nosql.UpdateDisplayState(mine.UID, operator, st)
	if err == nil {
		opt := LogOptNull
		if mine.Status == FavStatusPending && st == FavStatusDraft {
			opt = LogOptRefuse
			for _, content := range mine.Contents {
				if !mine.HadPending(content.UID) {
					mine.Pending = append(mine.Pending, content)
				}
			}
			mine.Contents = make([]proxy.DisplayContent, 0, 1)
		} else if mine.Status == FavStatusPending && st == FavStatusPublish {
			opt = LogOptAgree
			if len(mine.Contents) < 1 {
				return errors.New("the stable content is empty")
			}
		} else if mine.Status == FavStatusDraft && st == FavStatusPending {
			opt = LogOptPend
			//if len(mine.Pending) < 1 {
			//	return errors.New("the pending content is empty")
			//}
		}
		mine.createHistory(operator, remark, "", mine.Status, st, opt)
		mine.Status = st
		mine.Operator = operator
		if st == FavStatusPublish {
			_ = cacheCtx.updateRecord(mine.Owner, RecodeFav, 1)
		}
	}
	return err
}

func (mine *DisplayInfo) GetHistories() []*nosql.History {
	dbs, _ := nosql.GetHistories(mine.UID)
	list := make([]*nosql.History, 0, len(dbs))
	for _, db := range dbs {
		list = append(list, db)
	}
	return list
}

func (mine *DisplayInfo) createHistory(operator, remark, content string, from, to uint8, opt uint32) {
	_ = cacheCtx.insertHistory(mine.UID, operator, remark, content, strconv.Itoa(int(from)), strconv.Itoa(int(to)), opt, HistoryDisplay)
}

func (mine *DisplayInfo) UpdateContents(operator string, list []*pb.DisplayContent) error {
	var err error
	if list == nil || len(list) < 1 {
		err = nosql.UpdateDisplayContents(mine.UID, operator, make([]proxy.DisplayContent, 0, 1))
		if err == nil {
			mine.Contents = make([]proxy.DisplayContent, 0, 1)
			mine.Operator = operator
		}
	} else {
		utc := time.Now().UTC().Unix()
		arr := make([]proxy.DisplayContent, 0, len(list))
		for _, s := range list {
			arr = append(arr, proxy.DisplayContent{UID: s.Uid, Remark: s.Remark, Submitter: operator,
				Stamp: utc, Events: s.Events, Assets: s.Assets})
		}
		err = nosql.UpdateDisplayContents(mine.UID, operator, arr)
		if err == nil {
			mine.Contents = arr
			mine.Operator = operator
		}
	}
	return err
}

func (mine *DisplayInfo) UpdatePending(operator string, list []*pb.DisplayContent) error {
	var err error
	if list == nil || len(list) < 1 {
		err = nosql.UpdateDisplayPending(mine.UID, operator, make([]proxy.DisplayContent, 0, 1))
		if err == nil {
			mine.Pending = make([]proxy.DisplayContent, 0, 1)
			mine.Operator = operator
		}
	} else {
		utc := time.Now().UTC().Unix()
		arr := make([]proxy.DisplayContent, 0, len(mine.Pending))
		for _, content := range mine.Pending {
			arr = append(arr, content)
		}
		for _, s := range list {
			if s.Option == uint32(ContentOptionAppend) && !mine.HadStableContent(s.Uid) && !mine.HadPending(s.Uid) {
				arr = append(arr, proxy.DisplayContent{UID: s.Uid, Option: s.Option, Submitter: operator,
					Stamp: utc, Remark: s.Remark, Events: s.Events, Assets: s.Assets})
			}
		}
		err = nosql.UpdateDisplayPending(mine.UID, operator, arr)
		if err == nil {
			mine.Pending = arr
			mine.Operator = operator
		}
	}
	return err
}

func (mine *DisplayInfo) HadStableContent(uid string) bool {
	for _, item := range mine.Contents {
		if item.UID == uid {
			return true
		}
	}

	return false
}

func (mine *DisplayInfo) HadPending(uid string) bool {
	for _, item := range mine.Pending {
		if item.UID == uid {
			return true
		}
	}

	return false
}

func (mine *DisplayInfo) GetPendContent(uid string) *proxy.DisplayContent {
	for _, item := range mine.Pending {
		if item.UID == uid {
			return &item
		}
	}
	return nil
}

func (mine *DisplayInfo) removePending(uid string) {
	for i := 0; i < len(mine.Pending); i += 1 {
		if mine.Pending[i].UID == uid {
			if i == len(mine.Pending)-1 {
				mine.Pending = append(mine.Pending[:i])
			} else {
				mine.Pending = append(mine.Pending[:i], mine.Pending[i+1:]...)
			}
			break
		}
	}
}

func (mine *DisplayInfo) AgreePending(uid, remark, operator string) error {
	item := mine.GetPendContent(uid)
	if item == nil {
		return errors.New("the item not found")
	}
	if item.Option == uint32(ContentOptionAppend) {
		return mine.AppendContent(uid, remark, operator, ContentOptionAppend, true)
	} else {
		er := mine.SubtractContent(uid, remark, operator, true, true)
		if er == nil {
			er = mine.SubtractContent(uid, remark, operator, false, false)
		}
		return er
	}
}

func (mine *DisplayInfo) AppendContent(uid, remark, operator string, opt ContentOption, stable bool) error {
	tmp := proxy.DisplayContent{
		UID:       uid,
		Submitter: operator,
		Option:    uint32(opt),
		Events:    make([]string, 0, 1),
		Assets:    make([]string, 0, 1),
	}
	var err error
	if stable {
		if mine.HadStableContent(uid) {
			return nil
		}
		err = nosql.AppendDisplayContent(mine.UID, operator, tmp)
		if err == nil {
			_ = mine.SubtractContent(uid, remark, operator, false, false)
			mine.createHistory(operator, remark, uid, FavStatusPending, FavStatusPublish, LogOptAgreeAdd)
			mine.Contents = append(mine.Contents, tmp)
		}
	} else {
		if opt == ContentOptionAppend && mine.HadStableContent(uid) {
			return nil
		}
		item := mine.GetPendContent(uid)
		if item != nil {
			return nil
		}
		err = nosql.AppendDisplayPending(mine.UID, operator, tmp)
		if err == nil {
			mine.Pending = append(mine.Pending, tmp)
			op := LogOptNull
			if tmp.Option == uint32(ContentOptionAppend) {
				op = LogOptRequestAdd
			} else {
				op = LogOptRequestDel
			}
			mine.createHistory(operator, remark, uid, FavStatusDraft, FavStatusPending, op)
		}
	}
	return err
}

func (mine *DisplayInfo) SubtractContent(uid, remark, operator string, stable, log bool) error {
	var er error
	if stable {
		if !mine.HadStableContent(uid) {
			return nil
		}
		er = nosql.SubtractDisplayContent(mine.UID, uid, operator)
		if er == nil {
			if log {
				mine.createHistory(operator, remark, uid, FavStatusPublish, FavStatusDraft, LogOptAgreeDel)
			}

			for i := 0; i < len(mine.Contents); i += 1 {
				if mine.Contents[i].UID == uid {
					if i == len(mine.Contents)-1 {
						mine.Contents = append(mine.Contents[:i])
					} else {
						mine.Contents = append(mine.Contents[:i], mine.Contents[i+1:]...)
					}
					break
				}
			}
		}
	} else {
		pend := mine.GetPendContent(uid)
		if pend == nil {
			return nil
		}
		er = nosql.SubtractDisplayPending(mine.UID, uid, operator)
		if er == nil {
			if log {
				op := LogOptNull
				if pend.Option == uint32(ContentOptionAppend) {
					op = LogOptRefuseAdd
				} else {
					op = LogOptRefuseDel
				}
				mine.createHistory(operator, remark, uid, FavStatusPending, FavStatusDraft, op)
			}
			mine.removePending(uid)
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
