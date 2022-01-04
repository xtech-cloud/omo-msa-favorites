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
	UpdatedTime time.Time `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time `json:"deleteAt" bson:"deleteAt"`
	Creator     string    `json:"creator" bson:"creator"`
	Operator    string    `json:"operator" bson:"operator"`
	Cover       string    `json:"cover" bson:"cover"`
	Remark      string    `json:"remark" bson:"remark"`
	Owner       string    `json:"owner" bson:"owner"`
	State       uint8 	  `json:"state" bson:"state"`
	Type        uint8     `json:"type" bson:"type"`
	Origin      string    `json:"origin" bson:"origin"` //数据来源，可能是某次活动
	Meta        string 	  `json:"meta" bson:"meta"` //源数据
	Tags        []string  `json:"tags" bsonL:"tags"`
	Keys        []string  `json:"keys" bson:"keys"`
	Targets     []*proxy.ShowingInfo `json:"targets" bson:"targets"`
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

func GetFavoriteByName(table, owner, name string, tp uint8) (*Favorite, error) {
	if len(owner) < 2 || len(name) < 2{
		return nil, errors.New("db owner or name is empty of GetFavoriteByName")
	}
	filter := bson.M{"owner":owner, "name": name, "type":tp, "deleteAt": new(time.Time)}
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

func GetFavoritesByOwnerTP(table, owner string, kind uint8) ([]*Favorite, error) {
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

func GetFavoritesByType(table string, kind uint8) ([]*Favorite, error) {
	def := new(time.Time)
	filter := bson.M{"type":kind, "deleteAt": def}
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

func GetFavoritesByTarget(table, target string) ([]*Favorite, error) {
	def := new(time.Time)
	filter := bson.M{"targets": bson.M{"target": target} , "deleteAt": def}
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

func UpdateFavoriteState(table, uid, operator string, st uint8) error {
	msg := bson.M{"state": st, "operator":operator, "updatedAt": time.Now()}
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

func AppendFavoriteKey(table, uid string, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"keys": key}
	_, err := appendElement(table, uid, msg)
	return err
}

func SubtractFavoriteKey(table, uid, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"keys": key}
	_, err := removeElement(table, uid, msg)
	return err
}

func UpdateFavoriteTarget(table, uid, operator string, list []*proxy.ShowingInfo) error {
	msg := bson.M{"targets": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(table, uid, msg)
	return err
}

func AppendFavoriteTarget(table, uid string, target *proxy.ShowingInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"targets": target}
	_, err := appendElement(table, uid, msg)
	return err
}

func SubtractFavoriteTarget(table, uid, target string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"targets": bson.M{"target":target}}
	_, err := removeElement(table, uid, msg)
	return err
}
