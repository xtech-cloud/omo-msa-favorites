package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Repertory struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Remark      string             `json:"remark" bson:"remark"`
	Owner       string             `json:"owner" bson:"owner"`
	Type        uint8              `json:"type" bson:"type"`
	Bags      []string `json:"bags" bson:"bags"`
}

func CreateRepertory(info *Repertory) error {
	_, err := insertOne(TableRepertory, &info)
	return err
}

func GetRepertoryNextID() uint64 {
	num, _ := getSequenceNext(TableRepertory)
	return num
}

func GetRepertoryCount() int64 {
	num, _ := getTotalCount(TableRepertory)
	return num
}

func GetRepertories() ([]*Repertory, error) {
	cursor, err1 := findAll(TableRepertory, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Repertory, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Repertory)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveRepertory(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Repertory uid is empty ")
	}
	_, err := removeOne(TableRepertory, uid, operator)
	return err
}

func GetRepertory(uid string) (*Repertory, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Repertory uid is empty of GetRepertory")
	}

	result, err := findOne(TableRepertory, uid)
	if err != nil {
		return nil, err
	}
	model := new(Repertory)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetRepertoryByOwner(owner string) (*Repertory, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	result, err := findOneBy(TableRepertory, filter)
	if err != nil {
		return nil, err
	}
	model := new(Repertory)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateRepertoryBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableRepertory, uid, msg)
	return err
}

func UpdateRepertoryBags(uid, operator string, list []string) error {
	msg := bson.M{"bags": list, "operator":operator, "updatedAt": time.Now()}
	_, err := updateOne(TableRepertory, uid, msg)
	return err
}

func AppendRepertoryBag(uid, prop string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"bags": prop}
	_, err := appendElement(TableRepertory, uid, msg)
	return err
}

func SubtractRepertoryBag(uid, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"bags": key}
	_, err := removeElement(TableRepertory, uid, msg)
	return err
}
