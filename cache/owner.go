package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

const (
	OwnerTypePerson = 1
	OwnerTypeUnit = 0
)

type OwnerInfo struct {
	BaseInfo
	Owner string
	Remark string
	favorites []*FavoriteInfo
	bags []string
}

func AllOwners() []*OwnerInfo {
	return cacheCtx.boxes
}

func GetOwner(uid string) *OwnerInfo {
	for i := 0;i < len(cacheCtx.boxes);i += 1{
		if cacheCtx.boxes[i].Owner == uid {
			return cacheCtx.boxes[i]
		}
	}
	db,err := nosql.GetRepertoryByOwner(uid)
	if err == nil {
		tmp := new(OwnerInfo)
		tmp.initByDB(db)
		cacheCtx.boxes = append(cacheCtx.boxes, tmp)
		return tmp
	}
	info := new(OwnerInfo)
	info.initInfo(uid)
	cacheCtx.boxes = append(cacheCtx.boxes, info)
	return info
}

func (mine *OwnerInfo)initInfo(owner string)  {
	mine.Owner = owner
	mine.bags = make([]string, 0, 1)
	array,err := nosql.GetFavoritesByOwner(owner)
	if err == nil{
		mine.favorites = make([]*FavoriteInfo, 0, len(array))
		for _, value := range array {
			fav := new(FavoriteInfo)
			fav.initInfo(value)
			mine.favorites = append(mine.favorites, fav)
		}
	}else{
		mine.favorites = make([]*FavoriteInfo, 0, 1)
	}
}

func (mine *OwnerInfo)initByDB(db *nosql.Repertory)  {
	mine.initInfo(db.Owner)
	mine.UID = db.UID.Hex()
	mine.Remark = db.Remark
	mine.ID = db.ID
	mine.Name = db.Name
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
}

func (mine *OwnerInfo)Favorites() []*FavoriteInfo {
	return mine.favorites
}

func (mine *OwnerInfo)addFavorite(db *nosql.Favorite)  {
	info := new(FavoriteInfo)
	info.initInfo(db)
	mine.favorites = append(mine.favorites, info)
}

func (mine *OwnerInfo)CreateFavorite(info *FavoriteInfo) error {
	db := new(nosql.Favorite)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetFavoriteNextID()
	db.CreatedTime = time.Now()
	db.Cover = info.Cover
	db.Name = info.Name
	db.Remark = info.Remark
	db.Owner = mine.UID
	db.Type = info.Type
	db.Creator = info.Creator
	db.Operator = info.Operator
	db.Entities = make([]string, 0, 1)
	err := nosql.CreateFavorite(db)
	if err == nil {
		info.initInfo(db)
		mine.favorites = append(mine.favorites, info)
	}
	return err
}

func (mine *OwnerInfo)GetFavorite(uid string) *FavoriteInfo {
	for i := 0;i < len(mine.favorites);i += 1{
		if mine.favorites[i].UID == uid {
			return mine.favorites[i]
		}
	}
	return nil
}

func (mine *OwnerInfo)RemoveFavorite(uid, operator string) error {
	err := nosql.RemoveFavorite(uid, operator)
	if err == nil {
		for i := 0;i < len(mine.favorites);i += 1{
			if mine.favorites[i].UID == uid {
				mine.favorites = append(mine.favorites[:i], mine.favorites[i+1:]...)
				break
			}
		}
	}
	return err
}

func (mine *OwnerInfo)createRepertory(uid string) (*nosql.Repertory,error) {
	db := new(nosql.Repertory)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRepertoryNextID()
	db.CreatedTime = time.Now()
	db.Owner = uid
	err := nosql.CreateRepertory(db)
	if err != nil {
		return nil, err
	}else{
		return db, nil
	}
}

func (mine *OwnerInfo)Bags() []string {
	return mine.bags
}

func (mine *OwnerInfo)HadBag(uid string) bool {
	for _, bag := range mine.bags {
		if bag == uid {
			return true
		}
	}
	return false
}

func (mine *OwnerInfo)AppendBag(uid string) error {
	if mine.UID == "" {
		db,err := mine.createRepertory(mine.Owner)
		if err != nil {
			return err
		}
		mine.UID = db.UID.Hex()
		mine.Creator = db.Creator
		mine.CreateTime = db.CreatedTime
		mine.UpdateTime = db.UpdatedTime
	}
	if mine.HadBag(uid) {
		return nil
	}
	er := nosql.AppendRepertoryBag(mine.UID, uid)
	if er == nil {
		mine.bags = append(mine.bags, uid)
	}
	return er
}

func (mine *OwnerInfo)SubtractBag(uid string) error {
	if !mine.HadBag(uid) {
		return nil
	}
	er := nosql.SubtractRepertoryBag(mine.UID, uid)
	if er == nil {
		for i := 0;i < len(mine.bags);i += 1 {
			if mine.bags[i] == uid {
				mine.bags = append(mine.bags[:i], mine.bags[i+1:]...)
				break
			}
		}
	}
	return er
}

func (mine *OwnerInfo)UpdateBags(list []string, operator string) error {
	er := nosql.UpdateRepertoryBags(mine.UID, operator, list)
	if er == nil {
		mine.bags = list
	}
	return er
}