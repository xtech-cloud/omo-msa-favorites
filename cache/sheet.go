package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

type SheetInfo struct {
	BaseInfo
	Status  uint8
	ProductType  uint8
	Owner   string //该展览表所属用户等
	Remark  string
	Quote    string //
	Keys    []string
}

func (mine *cacheContext) CreateSheet(info *SheetInfo) error {
	db := new(nosql.Sheet)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetSheetNextID()
	db.CreatedTime = time.Now()
	db.Name = info.Name
	db.Remark = info.Remark
	db.Owner = info.Owner
	db.Status = info.Status
	db.Product = info.ProductType
	db.Creator = info.Creator
	db.Operator = info.Operator
	db.Quote = info.Quote
	db.Keys = info.Keys
	if db.Keys == nil {
		db.Keys = make([]string, 0, 1)
	}

	err := nosql.CreateSheet(db)
	if err == nil {
		info.UID = db.UID.Hex()
		info.CreateTime = db.CreatedTime
		info.ID = db.ID
		info.UpdateTime = db.UpdatedTime
	}
	return err
}


func (mine *cacheContext) HadSheetByName(owner, name string) bool {
	fav, err := nosql.GetSheetByName(owner, name)
	if err != nil {
		return false
	}
	if fav != nil {
		return true
	} else {
		return true
	}
}

func (mine *cacheContext) RemoveSheet(uid, operator string) error {
	err := nosql.RemoveSheet(uid, operator)
	return err
}

func (mine *cacheContext) GetSheet(uid string) *SheetInfo {
	db, err := nosql.GetSheet(uid)
	if err == nil {
		info := new(SheetInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetSheetBy(owner, quote string, tp uint32) *SheetInfo {
	var db *nosql.Sheet
	var err error
	if len(owner) > 1 {
		db, err = nosql.GetSheetByQuote(owner, quote, uint8(tp))
	}else{
		db, err = nosql.GetSheetByQuote2(quote, uint8(tp))
	}

	if err == nil {
		info := new(SheetInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetSheetsByOwner(uid string) []*SheetInfo {
	array, err := nosql.GetSheetsByOwner(uid)
	if err == nil {
		list := make([]*SheetInfo, 0, len(array))
		for _, item := range array {
			info := new(SheetInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *cacheContext) GetSheetsByQuote(quote string) []*SheetInfo {
	array, err := nosql.GetSheetsByQuote(quote)
	list := make([]*SheetInfo, 0, 4)
	if err == nil {
		for _, item := range array {
			info := new(SheetInfo)
			info.initInfo(item)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetSheetsByProduct(uid string, tp uint8) []*SheetInfo {
	array, err := nosql.GetSheetsByOwnerTP(uid, tp)
	if err == nil {
		list := make([]*SheetInfo, 0, len(array))
		for _, item := range array {
			info := new(SheetInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return nil
}

func (mine *SheetInfo) initInfo(db *nosql.Sheet) {
	mine.UID = db.UID.Hex()
	mine.Remark = db.Remark
	mine.ID = db.ID
	mine.Name = db.Name
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Owner = db.Owner
	mine.Status = db.Status
	mine.ProductType = db.Product
	mine.Quote = db.Quote
	mine.Keys = db.Keys
	if mine.Keys == nil {
		mine.Keys = make([]string, 0, 1)
	}
}

func (mine *SheetInfo) UpdateBase(name, remark, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(remark) < 1 {
		remark = mine.Remark
	}
	err := nosql.UpdateSheetBase(mine.UID, name, remark, operator)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
	}
	return err
}

func (mine *SheetInfo) UpdateQuote(operator, quote string) error {
	err := nosql.UpdateSheetQuote(mine.UID, quote, operator)
	if err == nil {
		mine.Quote = quote
		mine.Operator = operator
	}
	return err
}

func (mine *SheetInfo) UpdateStatus(st uint8, operator string) error {
	err := nosql.UpdateSheetState(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
		if st == FavStatusPublish {
			_ = cacheCtx.updateRecord(mine.Owner, ObserveFav, 1)
		}
	}
	return err
}

func (mine *SheetInfo) UpdateKeys(operator string, list []string) error {
	var err error
	if list == nil || len(list) < 1{
		err = nosql.UpdateSheetKeys(mine.UID, operator, make([]string, 0, 1))
		if err == nil {
			mine.Keys = make([]string, 0, 1)
			mine.Operator = operator
		}
	}else{
		err = nosql.UpdateSheetKeys(mine.UID, operator, list)
		if err == nil {
			mine.Keys = list
			mine.Operator = operator
		}
	}
	return err
}

func (mine *SheetInfo) HadKey(uid string) bool {
	for _, item := range mine.Keys {
		if item == uid {
			return true
		}
	}
	return false
}

func (mine *SheetInfo) AppendKey(uid string) error {
	if mine.HadKey(uid) {
		return nil
	}
	er := nosql.AppendSheetKey(mine.UID, uid)
	if er == nil {
		mine.Keys = append(mine.Keys, uid)
	}
	return er
}

func (mine *SheetInfo) SubtractKey(uid string) error {
	if !mine.HadKey(uid) {
		return nil
	}
	er := nosql.SubtractSheetKey(mine.UID, uid)
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
