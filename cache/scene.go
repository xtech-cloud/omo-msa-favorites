package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
	"time"
)

type SceneInfo struct {
	UID string
	favorites []*FavoriteInfo
}

func AllScenes() []*SceneInfo {
	return cacheCtx.scenes
}

func GetScene(uid string) *SceneInfo {
	for i := 0;i < len(cacheCtx.scenes);i += 1{
		if cacheCtx.scenes[i].UID == uid {
			return cacheCtx.scenes[i]
		}
	}
	scene := new(SceneInfo)
	scene.initInfo(uid)
	cacheCtx.scenes = append(cacheCtx.scenes, scene)
	return scene
}

func (mine *SceneInfo)initInfo(scene string)  {
	mine.UID = scene
	array,err := nosql.GetFavoritesByScene(scene)
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

func (mine *SceneInfo)Favorites() []*FavoriteInfo {
	return mine.favorites
}

func (mine *SceneInfo)CreateFavorite(info *FavoriteInfo) error {
	db := new(nosql.Favorite)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetFavoriteNextID()
	db.CreatedTime = time.Now()
	db.Cover = info.Cover
	db.Name = info.Name
	db.Remark = info.Remark
	db.Scene = mine.UID
	db.Entities = make([]proxy.EntityInfo, 0, 1)
	err := nosql.CreateFavorite(db)
	if err == nil {
		info.initInfo(db)
		mine.favorites = append(mine.favorites, info)
	}
	return err
}

func (mine *SceneInfo)GetFavorite(uid string) *FavoriteInfo {
	for i := 0;i < len(mine.favorites);i += 1{
		if mine.favorites[i].UID == uid {
			return mine.favorites[i]
		}
	}
	return nil
}

func (mine *SceneInfo)RemoveFavorite(uid string) error {
	err := nosql.RemoveFavorite(uid)
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