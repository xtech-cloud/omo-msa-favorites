package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Notice struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Status  uint8 `json:"status" bson:"status"`
	Subtitle string `json:"subtitle" bson:"subtitle"`
	Body     string `json:"body" bson:"body"`
	Owner    string `json:"owner" bson:"owner"`

	Targets []string `json:"targets" bson:"targets"`
	Tags    []string `json:"tags" bson:"tags"`
}

func CreateNotice(info *Notice) error {
	_, err := insertOne(TableNotice, &info)
	return err
}

func GetNoticeNextID() uint64 {
	num, _ := getSequenceNext(TableNotice)
	return num
}

func RemoveNotice(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Notice uid is empty ")
	}
	_, err := removeOne(TableNotice, uid, operator)
	return err
}

func GetNotice(uid string) (*Notice, error) {
	if len(uid) < 2 {
		return nil, errors.New("db activity uid is empty of GetNotice")
	}
	result, err := findOne(TableNotice, uid)
	if err != nil {
		return nil, err
	}
	model := new(Notice)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetNoticeCount() int64 {
	num, _ := getTotalCount(TableNotice)
	return num
}

func GetNoticesByTargets(targets []string) ([]*Notice, error) {
	def := new(time.Time)
	in := bson.A{}
	for _, target := range targets {
		in = append(in, target)
	}
	filter := bson.M{ "$or":bson.A{bson.M{"targets": bson.M{"$in":in}}, bson.M{"targets":bson.M{"$ne":nil}}} , "deleteAt": def}
	cursor, err1 := findMany(TableNotice, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Notice, 0, 20)
	for cursor.Next(context.TODO()) {
		var node = new(Notice)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetNoticesByOTargets(owner string, targets []string) ([]*Notice, error) {
	def := new(time.Time)
	in := bson.A{}
	for _, target := range targets {
		in = append(in, target)
	}
	filter := bson.M{"owner":owner, "$or":bson.A{bson.M{"targets": bson.M{"$in":in}},bson.M{"targets":bson.M{"$ne":nil}}} , "deleteAt": def}
	cursor, err1 := findMany(TableNotice, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Notice, 0, 20)
	for cursor.Next(context.TODO()) {
		var node = new(Notice)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetNoticesByOwner(owner string) ([]*Notice, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	cursor, err1 := findMany(TableNotice, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Notice, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Notice)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateNoticeBase(uid, name, sub, body, operator string) error {
	msg := bson.M{"name": name, "body": body, "subtitle": sub, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableNotice, uid, msg)
	return err
}

func UpdateNoticeStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableNotice, uid, msg)
	return err
}

func UpdateNoticeTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableNotice, uid, msg)
	return err
}

func UpdateNoticeTargets(uid, operator string, list []string) error {
	msg := bson.M{"targets": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableNotice, uid, msg)
	return err
}
