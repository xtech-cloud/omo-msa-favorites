package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type History struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Type    uint8  `json:"type" bson:"type"`
	Option  uint32 `json:"option" bson:"option"`
	From    string `json:"from" json:"from"`
	To      string `json:"to" bson:"to"`
	Parent  string `json:"parent" bson:"parent"`
	Remark  string `json:"remark" bson:"remark"`
	Content string `json:"content" bson:"content"`
}

func CreateHistory(info *History) error {
	_, err := insertOne(TableHistory, info)
	if err != nil {
		return err
	}
	return nil
}

func GetHistoryNextID() uint64 {
	num, _ := getSequenceNext(TableHistory)
	return num
}

func GetHistory(uid string) (*History, error) {
	result, err := findOne(TableHistory, uid)
	if err != nil {
		return nil, err
	}
	model := new(History)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetHistories(uid string) ([]*History, error) {
	var items = make([]*History, 0, 20)
	filter := bson.M{"parent": uid}
	cursor, err1 := findMany(TableHistory, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(History)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetHistoriesBy(parent, to string, op uint8) ([]*History, error) {
	var items = make([]*History, 0, 20)
	filter := bson.M{"parent": parent, "option": op, "to": to}
	cursor, err1 := findMany(TableHistory, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(History)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}
