package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

// 展览表集合
type SheetInfo struct {
	BaseInfo
	Status      uint8
	ProductType uint8
	Owner       string //该展览表所属场景
	Remark      string
	Quote       string              //关联的场所UID， 班级等
	Contents    []proxy.ShowContent //展览集合
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
	db.Contents = info.Contents
	if db.Contents == nil {
		db.Contents = make([]proxy.ShowContent, 0, 1)
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

func (mine *cacheContext) IsUsed(display string) bool {
	list, err := nosql.GetSheetsByDisplay(display)
	if err != nil {
		return false
	}
	if len(list) > 0 {
		return true
	}
	arr, er := nosql.GetProductsByDisplay(display)
	if er != nil {
		return false
	}
	if len(arr) > 0 {
		return true
	}
	return false
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

func (mine *cacheContext) GetSheetBy(owner, quote string) *SheetInfo {
	var db *nosql.Sheet
	var err error
	if len(owner) > 1 {
		db, err = nosql.GetSheetByQuote3(owner, quote)
	} else {
		db, err = nosql.GetSheetByQuote4(quote)
	}
	//if tp < 1 {
	//
	//} else {
	//	if len(owner) > 1 {
	//		db, err = nosql.GetSheetByQuote(owner, quote, uint8(tp))
	//	} else {
	//		db, err = nosql.GetSheetByQuote2(quote, uint8(tp))
	//	}
	//}

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
	mine.Contents = db.Contents
	if mine.Contents == nil {
		mine.Contents = make([]proxy.ShowContent, 0, 1)
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

func (mine *SheetInfo) UpdateKeys(operator string, list []proxy.ShowContent) error {
	var err error
	if list == nil || len(list) < 1 {
		err = nosql.UpdateSheetDisplay(mine.UID, operator, make([]proxy.ShowContent, 0, 1))
		if err == nil {
			mine.Contents = make([]proxy.ShowContent, 0, 1)
			mine.Operator = operator
		}
	} else {
		err = nosql.UpdateSheetDisplay(mine.UID, operator, list)
		if err == nil {
			mine.Contents = list
			mine.Operator = operator
		}
	}
	return err
}

func (mine *SheetInfo) HadContent(uid string) bool {
	for _, item := range mine.Contents {
		if item.UID == uid {
			return true
		}
	}
	return false
}

func (mine *SheetInfo) AppendContent(uid, operator, effect, menu, align string, weight uint32) error {
	if mine.HadContent(uid) {
		return nil
	}
	tmp := proxy.ShowContent{
		UID:       uid,
		Weight:    weight,
		Effect:    effect,
		Menu:      menu,
		Alignment: align,
	}
	er := nosql.AppendSheetContent(mine.UID, operator, tmp)
	if er == nil {
		mine.Contents = append(mine.Contents, tmp)
	}
	return er
}

func (mine *SheetInfo) SubtractContent(uid, operator string) error {
	if !mine.HadContent(uid) {
		return nil
	}
	er := nosql.SubtractSheetContent(mine.UID, operator, uid)
	if er == nil {
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
	return er
}
