package nosql

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Record struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Type  uint8  `json:"type" bson:"type"`
	Count uint32 `json:"count" bson:"count"`
	Owner string `json:"owner" bson:"owner"`
	Begin string `json:"begin" bson:"begin"`
}

func CreateRecord(info *Record) error {
	_, err := insertOne(TableRecord, &info)
	return err
}

func GetRecordNextID() uint64 {
	num, _ := getSequenceNext(TableRecord)
	return num
}

func GetRecord(uid string) (*Record, error) {
	if len(uid) < 2 {
		return nil, errors.New("db activity uid is empty of GetRecord")
	}
	result, err := findOne(TableRecord, uid)
	if err != nil {
		return nil, err
	}
	model := new(Record)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetRecordCount() int64 {
	num, _ := getTotalCount(TableRecord)
	return num
}

func GetRecordsByType(owner string, tp uint8) (*Record, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "type": tp, "deleteAt": def}
	result, err := findOneBy(TableRecord, filter)
	if err != nil {
		return nil, err
	}
	model := new(Record)
	err = result.Decode(&model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func GetRecordsByDate(owner, begin string, tp uint8) (*Record, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "begin": begin, "type": tp, "deleteAt": def}
	result, err := findOneBy(TableRecord, filter)
	if err != nil {
		return nil, err
	}
	model := new(Record)
	err = result.Decode(&model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func UpdateRecordCount(uid string, count uint32) error {
	msg := bson.M{"count": count, "updatedAt": time.Now()}
	_, err := updateOne(TableRecord, uid, msg)
	return err
}
