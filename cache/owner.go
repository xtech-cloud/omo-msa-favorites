package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

const (
	OwnerTypePerson = 1
	OwnerTypeUnit = 0
)

type OwnerInfo struct {
	UID string
	favorites []*FavoriteInfo
}

func AllOwners() []*OwnerInfo {
	return cacheCtx.boxes
}

func GetOwner(uid string) *OwnerInfo {
	for i := 0;i < len(cacheCtx.boxes);i += 1{
		if cacheCtx.boxes[i].UID == uid {
			return cacheCtx.boxes[i]
		}
	}
	scene := new(OwnerInfo)
	scene.initInfo(uid)
	cacheCtx.boxes = append(cacheCtx.boxes, scene)
	return scene
}

func (mine *OwnerInfo)initInfo(owner string)  {
	mine.UID = owner
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

func (mine *OwnerInfo)Favorites() []*FavoriteInfo {
	return mine.favorites
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
	db.Entities = make([]proxy.EntityInfo, 0, 1)
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