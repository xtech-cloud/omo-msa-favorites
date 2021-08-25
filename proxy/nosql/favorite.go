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
	UpdatedTime time.Time `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time `json:"deleteAt" bson:"deleteAt"`
	Creator     string    `json:"creator" bson:"creator"`
	Operator    string    `json:"operator" bson:"operator"`
	Cover       string    `json:"cover" bson:"cover"`
	Remark      string    `json:"remark" bson:"remark"`
	Owner       string    `json:"owner" bson:"owner"`
	Type        uint8     `json:"type" bson:"type"`
	Origin      string    `json:"origin" bson:"origin"`
	Tags        []string  `json:"tags" bsonL:"tags"`
	Keys        []string  `json:"keys" bson:"keys"`
}

func CreateFavorite(table string, info *Favorite) error {
	_, err := insertOne(table, &info)
	return err
}

func GetFavoriteNextID(table string) uint64 {
	num, _ := getSequenceNext(table)
	return num
}

func GetFavorites(table string) ([]*Favorite, error) {
	cursor, err1 := findAll(table, 0)
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

func GetFavoriteFile(table, uid string) (*FileInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Favorite.files uid is empty ")
	}
	result, err := findOne(table, uid)
	if err != nil {
		return nil, err
	}
	info := new(FileInfo)
	err1 := result.Decode(&info)
	return info, err1
}

func RemoveFavorite(table, uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Favorite uid is empty ")
	}
	_, err := removeOne(table, uid, operator)
	return err
}

func GetFavorite(table, uid string) (*Favorite, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Favorite uid is empty of GetFavorite")
	}
	result, err := findOne(table, uid)
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

func GetFavoriteCount(table string) int64 {
	num, _ := getCount(table)
	return num
}

func GetFavoriteByOrigin(table, user, origin string) (*Favorite, error) {
	if len(origin) < 2 || len(user) < 2{
		return nil, errors.New("db origin uid is empty of GetFavorite")
	}
	filter := bson.M{"owner":user, "origin": origin, "deleteAt": new(time.Time)}
	result, err := findOneBy(table, filter)
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

func GetFavoriteByName(table, owner, name string) (*Favorite, error) {
	if len(owner) < 2 || len(name) < 2{
		return nil, errors.New("db owner or name is empty of GetFavoriteByName")
	}
	filter := bson.M{"owner":owner, "name": name, "deleteAt": new(time.Time)}
	result, err := findOneBy(table, filter)
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

func GetFavoritesByOwner(table, owner string) ([]*Favorite, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	cursor, err1 := findMany(table, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Favorite, 0, 20)
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


func GetFavoritesByType(table, owner string, kind uint8) ([]*Favorite, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "type":kind, "deleteAt": def}
	cursor, err1 := findMany(table, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Favorite, 0, 20)
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

func UpdateFavoriteBase(table, uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(table, uid, msg)
	return err
}

func UpdateFavoriteCover(table, uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(table, uid, msg)
	return err
}

func UpdateFavoriteTags(table, uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(table, uid, msg)
	return err
}

func UpdateFavoriteEntity(table, uid, operator string, list []string) error {
	msg := bson.M{"keys": list, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(table, uid, msg)
	return err
}

func AppendFavoriteEntity(table, uid string, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"keys": key}
	_, err := appendElement(table, uid, msg)
	return err
}

func SubtractFavoriteEntity(table, uid, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"keys": key}
	_, err := removeElement(table, uid, msg)
	return err
}
