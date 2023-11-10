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
	Status      uint8              `json:"status" bson:"status"`
	Type        uint8              `json:"type" bson:"type"`
	Meta        string             `json:"meta" bson:"meta"` //源数据
	Tags        []string           `json:"tags" bsonL:"tags"`
	Keys        []string           `json:"keys" bson:"keys"`
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
	num, _ := getTotalCount(TableFavorite)
	return num
}

func GetFavoriteByName(owner, name string, tp uint8) (*Favorite, error) {
	if len(owner) < 2 || len(name) < 2 {
		return nil, errors.New("db owner or name is empty of GetFavoriteByName")
	}
	filter := bson.M{"owner": owner, "name": name, "type": tp, "deleteAt": new(time.Time)}
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
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	cursor, err1 := findMany(TableFavorite, filter, 0)
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

func GetFavoritesByStatus(owner string, st uint8) ([]*Favorite, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "status": st, "deleteAt": def}
	cursor, err1 := findMany(TableFavorite, filter, 0)
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

func GetFavoritesByType(kind uint8) ([]*Favorite, error) {
	def := new(time.Time)
	filter := bson.M{"type": kind, "deleteAt": def}
	cursor, err1 := findMany(TableFavorite, filter, 0)
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

func GetFavoritesByOwnerTP(owner string, kind uint8) ([]*Favorite, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "type": kind, "deleteAt": def}
	cursor, err1 := findMany(TableFavorite, filter, 0)
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

func UpdateFavoriteBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func UpdateFavoriteCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func UpdateFavoriteState(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func UpdateFavoriteTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func UpdateFavoriteKeys(uid, operator string, list []string) error {
	msg := bson.M{"keys": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableFavorite, uid, msg)
	return err
}

func AppendFavoriteKey(uid string, key, operator string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"keys": key}
	_, err := appendElement(TableFavorite, uid, operator, msg)
	return err
}

func SubtractFavoriteKey(uid, key, operator string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"keys": key}
	_, err := removeElement(TableFavorite, uid, operator, msg)
	return err
}
