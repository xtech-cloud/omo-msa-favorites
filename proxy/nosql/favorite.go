package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Origin      string `json:"origin" bson:"origin"`
	Tags        []string `json:"tags" bsonL:"tags"`
	Entities    []string 		   `json:"entities" bson:"entities"`
}

func CreateFavorite(info *Favorite) error {
	_, err := insertOne(TableFavorite, &info)
	return err
}

func GetFavoriteNextID() uint64 {
	num, _ := getSequenceNext(TableFavorite)
	return num
}

func GetFavorites() ([]*Favorite, error) {
	cursor, err1 := findAll(TableFavorite, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Favorite, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Favorite)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
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

func GetFavoriteCount() int64 {
	num, _ := getCount(TableFavorite)
	return num
}

func GetFavoriteByOrigin(user, origin string) (*Favorite, error) {
	if len(origin) < 2 || len(user) < 2{
		return nil, errors.New("db origin uid is empty of GetFavorite")
	}
	filter := bson.M{"owner":user, "origin": origin, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableFavorite, filter)
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

func UpdateFavoriteBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func UpdateFavoriteCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func UpdateFavoriteTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}


func UpdateFavoriteEntity(uid, operator string, list []string) error {
	msg := bson.M{"entities": list, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func AppendFavoriteEntity(uid string, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"entities": key}
	_, err := appendElement(TableFavorite, uid, msg)
	return err
}

func SubtractFavoriteEntity(uid, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"entities": key}
	_, err := removeElement(TableFavorite, uid, msg)
	return err
}
