package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy"
	"time"
)

type Favorite struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Cover        string              `json:"cover" bson:"cover"`
	Remark       string `json:"remark" bson:"remark"`
	Scene        string `json:"scene" bson:"scene"`
	Entities    []proxy.EntityInfo `json:"entities" bson:"entities"`
}

func CreateFavorite(info *Favorite) error {
	_, err := insertOne(TableFavorite, &info)
	return err
}

func GetFavoriteNextID() uint64 {
	num, _ := getSequenceNext(TableFavorite)
	return num
}

func GetFavoriteFile(uid string) (*FileInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Favorite.files uid is empty ")
	}
	result, err := findOne(TableFavorite, uid)
	if err != nil {
		return nil, err
	}
	info := new(FileInfo)
	err1 := result.Decode(&info)
	return info, err1
}

func RemoveFavorite(uid string) error {
	if len(uid) < 2 {
		return errors.New("db Favorite uid is empty ")
	}
	_, err := removeOne(TableFavorite, uid)
	return err
}

func GetFavorite(uid string) (*Favorite, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Favorite uid is empty of GetFavorite")
	}

	result, err := findOne(TableFavorite, uid)
	if err != nil {
		return nil, err
	}
	model := new(Favorite)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetFavoritesByScene(scene string) ([]*Favorite, error) {
	var items = make([]*Favorite, 0, 20)
	def := new(time.Time)
	filter := bson.M{"scene": scene, "deleteAt": def}
	cursor, err1 := findMany(TableFavorite, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	for cursor.Next(context.Background()) {
		var node = new(Favorite)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateFavoriteBase(uid string, name, remark string) error {
	msg := bson.M{"name": name,"remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func UpdateFavoriteCover(uid string, cover string) error {
	msg := bson.M{"cover": cover, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func AppendFavoriteEntity(uid string, entity proxy.EntityInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"entities": entity}
	_, err := appendElement(TableFavorite, uid, msg)
	return err
}

func SubtractFavoriteEntity(uid string, entity string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"entities": bson.M{ "uid": entity }}
	_, err := removeElement(TableFavorite, uid, msg)
	return err
}
