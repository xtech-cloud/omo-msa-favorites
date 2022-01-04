package cache

import (
	"errors"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

type FavoriteInfo struct {
	BaseInfo
	Status uint8
	Type   uint8 //
	Owner  string //该展览所属组织机构，scene, class等
	Cover  string
	Remark string
	Origin string //布展数据来源，比如活动
	Meta string //
	table string
	Tags   []string
	Keys   []string
	Targets []*proxy.ShowingInfo //目标设备
}

func getFavoriteTable(person bool) string {
	if person {
		return nosql.TableFavorite+ "_person"
	}else{
		return nosql.TableFavorite+"_scene"
	}
}


func (mine *cacheContext)CreateFavorite(info *FavoriteInfo, person bool) error {
	table := getFavoriteTable(person)
	db := new(nosql.Favorite)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetFavoriteNextID(table)
	db.CreatedTime = time.Now()
	db.Cover = info.Cover
	db.Name = info.Name
	db.Remark = info.Remark
	db.Owner = info.Owner
	db.Type = info.Type
	db.Origin = info.Origin
	db.State = info.Status
	db.Creator = info.Creator
	db.Operator = info.Operator
	db.Tags = info.Tags
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	db.Keys = info.Keys
	if db.Keys == nil {
		db.Keys = make([]string, 0, 1)
	}
	db.Targets = info.Targets
	if db.Targets == nil {
		db.Targets = make([]*proxy.ShowingInfo, 0, 1)
	}

	err := nosql.CreateFavorite(table, db)
	if err == nil {
		info.UID = db.UID.Hex()
		info.CreateTime = db.CreatedTime
		info.ID = db.ID
		info.UpdateTime = db.UpdatedTime
	}
	return err
}

func (mine *cacheContext)HadFavoriteByName(owner, name string, tp uint8, person bool) bool {
	table := getFavoriteTable(person)
	fav, err := nosql.GetFavoriteByName(table, owner, name, tp)
	if err != nil {
		return false
	}
	if fav != nil {
		return true
	}else{
		return true
	}
}

func (mine *cacheContext)RemoveFavorite(uid, operator string, person bool) error {
	table := getFavoriteTable(person)
	err := nosql.RemoveFavorite(table, uid, operator)
	return err
}

func (mine *cacheContext)GetFavorite(uid string, person bool) *FavoriteInfo {
	table := getFavoriteTable(person)
	db,err := nosql.GetFavorite(table, uid)
	if err == nil{
		info:= new(FavoriteInfo)
		info.initInfo(db, table)
		return info
	}
	return nil
}

func (mine *cacheContext)GetFavoriteByOrigin(user, uid string, person bool) *FavoriteInfo {
	table := getFavoriteTable(person)
	db,err := nosql.GetFavoriteByOrigin(table, user, uid)
	if err == nil{
		info:= new(FavoriteInfo)
		info.initInfo(db, table)
		return info
	}
	return nil
}

func (mine *cacheContext)GetFavoritesByOwner(uid string, person bool) []*FavoriteInfo {
	table := getFavoriteTable(person)
	array,err := nosql.GetFavoritesByOwner(table, uid)
	if err == nil{
		list := make([]*FavoriteInfo, 0, len(array))
		for _, item := range array {
			info := new(FavoriteInfo)
			info.initInfo(item, table)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext)GetFavoritesByType(owner string, kind uint8, person bool) []*FavoriteInfo {
	table := getFavoriteTable(person)
	var array []*nosql.Favorite
	var err error
	if kind == 1 {
		array,err = nosql.GetFavoritesByType(table, kind)
	}else{
		array,err = nosql.GetFavoritesByOwnerTP(table, owner, kind)
	}
	if err == nil{
		list := make([]*FavoriteInfo, 0, len(array))
		for _, item := range array {
			info := new(FavoriteInfo)
			info.initInfo(item, table)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext)GetFavoritesByList(person bool, array []string) []*FavoriteInfo {
	if array == nil || len(array) < 1 {
		return make([]*FavoriteInfo, 0, 1)
	}
	table := getFavoriteTable(person)
	list := make([]*FavoriteInfo, 0, 1)
	for _, s := range array {
		db,_ := nosql.GetFavorite(table, s)
		if db != nil {
			info := new(FavoriteInfo)
			info.initInfo(db, table)
			list = append(list, info)
		}
	}

	return list
}

func (mine *cacheContext) GetFavoritesByTargets(array []string, page, num uint32) (uint32, uint32, []*FavoriteInfo) {
	if array == nil || len(array) < 1 {
		return 0, 0, make([]*FavoriteInfo, 0, 1)
	}
	all := make([]*FavoriteInfo, 0, 10)
	table := getFavoriteTable(false)
	for _, s := range array {
		db, _ := nosql.GetFavoritesByTarget(table, s)
		if db != nil {
			for _, item := range db {
				info := new(FavoriteInfo)
				info.initInfo(item, table)
				all = append(all, info)
			}
		}
	}

	if num < 1 {
		num = 10
	}
	if len(all) < 1 {
		return 0, 0, make([]*FavoriteInfo, 0, 1)
	}
	max, pages, list := checkPage(page, num, all)
	return max, pages, list.([]*FavoriteInfo)
}

func (mine *FavoriteInfo)initInfo(db *nosql.Favorite, table string)  {
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
	mine.table = table
	mine.Status = db.State
	mine.Targets = db.Targets
	if mine.Keys == nil {
		mine.Keys = make([]string, 0 ,1)
	}
	if mine.Targets == nil {
		mine.Targets = make([]*proxy.ShowingInfo, 0, 1)
	}
}

func (mine *FavoriteInfo) GetKeys() []string {
	return mine.Keys
}

func (mine *FavoriteInfo) GetTargets() []*pb.ShowInfo {
	list := make([]*pb.ShowInfo, 0, len(mine.Targets))
	for _, item := range mine.Targets {
		list = append(list, &pb.ShowInfo{Target: item.Target, Effect: item.Effect, Skin: item.Skin, Slots: item.Slots} )
	}
	return list
}

func (mine *FavoriteInfo)UpdateBase(name, remark,operator string) error {
	if len(name) <1 {
		name = mine.Name
	}
	if len(remark) <1 {
		remark = mine.Remark
	}
	err := nosql.UpdateFavoriteBase(mine.table, mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo)UpdateMeta(operator, meta string) error {
	return nil
}

func (mine *FavoriteInfo)UpdateTags(operator string, tags []string) error {
	if tags == nil {
		return errors.New("the list of target is nil")
	}
	err := nosql.UpdateFavoriteTags(mine.table, mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo)UpdateCover(cover, operator string) error {
	err := nosql.UpdateFavoriteCover(mine.table, mine.UID, cover, operator)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo)UpdateStatus(st uint8, operator string) error {
	err := nosql.UpdateFavoriteState(mine.table, mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo) UpdateEntities(operator string, list []string) error {
	if list == nil {
		return errors.New("the list of target is nil")
	}
	err := nosql.UpdateFavoriteEntity(mine.table, mine.UID, operator, list)
	if err == nil {
		mine.Keys = list
		mine.Operator = operator
	}
	return err
}

func (mine *FavoriteInfo) HadKey(uid string) bool {
	for _, item := range mine.Keys {
		if item == uid {
			return true
		}
	}
	return false
}

func (mine *FavoriteInfo) HadTarget(uid string) bool {
	for _, item := range mine.Targets {
		if item.Target == uid {
			return true
		}
	}
	return false
}

func (mine *FavoriteInfo) AppendKey(uid string) error {
	if mine.HadKey(uid) {
		return nil
	}
	er := nosql.AppendFavoriteKey(mine.table, mine.UID, uid)
	if er == nil {
		mine.Keys = append(mine.Keys, uid)
	}
	return er
}

func (mine *FavoriteInfo) SubtractKey(uid string) error {
	if !mine.HadKey(uid) {
		return nil
	}
	er := nosql.SubtractFavoriteKey(mine.table, mine.UID, uid)
	if er == nil {
		for i := 0;i < len(mine.Keys);i += 1 {
			if mine.Keys[i] == uid {
				if i == len(mine.Keys) - 1{
					mine.Keys = append(mine.Keys[:i])
				}else{
					mine.Keys = append(mine.Keys[:i], mine.Keys[i+1:]...)
				}
				break
			}
		}
	}
	return er
}

func (mine *FavoriteInfo)UpdateTarget(uid, effect, skin, operator string, slots []string) error {
	if !mine.HadTarget(uid){
		return nil
	}
	if slots == nil {
		slots = make([]string, 0, 1)
	}
	array := make([]*proxy.ShowingInfo, 0, len(mine.Targets))
	for _, info := range mine.Targets {
		if info.Target == uid {
			info.Effect = effect
			info.Skin = skin
			info.Slots = slots
		}
		array = append(array, info)
	}
	err := nosql.UpdateFavoriteTarget(mine.table, mine.UID, operator, array)
	if err == nil {
		mine.Targets = array
	}
	return err
}

func (mine *FavoriteInfo)AppendTarget(show *proxy.ShowingInfo) error {
	if mine.HadTarget(show.Target) {
		_ = mine.SubtractTarget(show.Target)
	}
	er := nosql.AppendFavoriteTarget(mine.table, mine.UID, show)
	if er == nil {
		mine.Targets = append(mine.Targets, show)
	}
	return er
}

func (mine *FavoriteInfo)SubtractTarget(uid string) error {
	if !mine.HadTarget(uid) {
		return nil
	}
	er := nosql.SubtractFavoriteTarget(mine.table, mine.UID, uid)
	if er == nil {
		for i := 0;i < len(mine.Targets);i += 1 {
			if mine.Targets[i].Target == uid {
				if i == len(mine.Targets) - 1{
					mine.Targets = append(mine.Targets[:i])
				}else{
					mine.Targets = append(mine.Targets[:i], mine.Targets[i+1:]...)
				}
				break
			}
		}
	}
	return er
}

