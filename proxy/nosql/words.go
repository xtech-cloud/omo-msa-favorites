package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Words struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Words  string   `json:"words" bson:"words"`
	Owner  string   `json:"owner" bson:"owner"`
	Target string   `json:"target" bson:"target"` //给目标留言
	Type   uint8    `json:"type" bson:"type"`
	Weight int32    `json:"weight" bson:"weight"`
	Quote  string   `json:"quote" bson:"quote"`
	Device string   `json:"device" bson:"device"`
	Assets []string `json:"assets" bson:"assets"`
}

func CreateWords(info *Words) error {
	_, err := insertOne(TableWords, &info)
	return err
}

func GetWordsNextID() uint64 {
	num, _ := getSequenceNext(TableWords)
	return num
}

func RemoveWords(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db sheet uid is empty ")
	}
	_, err := removeOne(TableWords, uid, operator)
	return err
}

func GetWords(uid string) (*Words, error) {
	if len(uid) < 2 {
		return nil, errors.New("db sheet uid is empty of GetWords")
	}
	result, err := findOne(TableWords, uid)
	if err != nil {
		return nil, err
	}
	model := new(Words)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetWordsCount() int64 {
	num, _ := getTotalCount(TableWords)
	return num
}

func GetWordsByOwnerType(owner string, tp uint8) ([]*Words, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "type": tp, "deleteAt": def}
	cursor, err1 := findMany(TableWords, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Words, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Words)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetWordsByDevice(device string) ([]*Words, error) {
	def := new(time.Time)
	filter := bson.M{"device": device, "deleteAt": def}
	cursor, err1 := findMany(TableWords, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Words, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Words)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetWordsCountByDevice(device string) (int64, error) {
	def := new(time.Time)
	filter := bson.M{"device": device, "deleteAt": def}
	return getCount(TableWords, filter)
}

func GetWordsCountByScene(owner string) (int64, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	return getCount(TableWords, filter)
}

func GetWordsCountByDate(device string) (int64, error) {
	def := new(time.Time)
	filter := bson.M{"device": device, "deleteAt": def}
	return getCount(TableWords, filter)
}

func GetWordsByCreator(owner, user, device string, tp uint8) ([]*Words, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "creator": user, "device": device, "type": tp, "deleteAt": def}
	cursor, err1 := findMany(TableWords, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Words, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Words)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetWordsByUser(user string) ([]*Words, error) {
	def := new(time.Time)
	filter := bson.M{"creator": user, "deleteAt": def}
	cursor, err1 := findMany(TableWords, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Words, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Words)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetWordsByOwner(owner string) ([]*Words, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	cursor, err1 := findMany(TableWords, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Words, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Words)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetWordsByQuote(owner, quote string) ([]*Words, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "quote": quote, "deleteAt": def}
	cursor, err1 := findMany(TableWords, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Words, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Words)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetWordsByTarget(target string) ([]*Words, error) {
	def := new(time.Time)
	filter := bson.M{"target": target, "deleteAt": def}
	cursor, err1 := findMany(TableWords, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Words, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Words)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateWordsBase(uid, words, operator string) error {
	msg := bson.M{"words": words, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableWords, uid, msg)
	return err
}

func UpdateWordsState(uid, operator string, st int32) error {
	msg := bson.M{"weight": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableWords, uid, msg)
	return err
}

func UpdateWordsQuote(uid, quote, operator string) error {
	msg := bson.M{"quote": quote, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableWords, uid, msg)
	return err
}

func UpdateWordsAssets(uid, operator string, list []string) error {
	msg := bson.M{"assets": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableWords, uid, msg)
	return err
}

func AppendWordsKey(uid string, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"keys": key}
	_, err := appendElement(TableWords, uid, msg)
	return err
}

func SubtractWordsKey(uid, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"keys": key}
	_, err := removeElement(TableWords, uid, msg)
	return err
}
