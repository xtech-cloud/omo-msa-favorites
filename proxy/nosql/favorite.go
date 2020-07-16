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
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`
	Cover       string             `json:"cover" bson:"cover"`
	Remark      string             `json:"remark" bson:"remark"`
	Owner       string             `json:"owner" bson:"owner"`
	Type        uint8              `json:"type" bson:"type"`
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

func RemoveFavorite(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Favorite uid is empty ")
	}
	_, err := removeOne(TableFavorite, uid, operator)
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

func GetFavoritesByOwner(owner string) ([]*Favorite, error) {
	var items = make([]*Favorite, 0, 20)
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
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

func UpdateFavoriteBase(uid string, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func UpdateFavoriteCover(uid string, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func UpdateFavoriteEntity(uid, operator string, list []proxy.EntityInfo) error {
	msg := bson.M{"entities": list, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}
