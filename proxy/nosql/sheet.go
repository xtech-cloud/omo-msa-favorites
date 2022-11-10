package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Sheet struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Remark  string   `json:"remark" bson:"remark"`
	Owner   string   `json:"owner" bson:"owner"`
	Quote   string   `json:"quote" bson:"quote"`
	Status  uint8    `json:"status" bson:"status"`
	Product uint8    `json:"product" bson:"product"`
	Keys    []string `json:"keys" bson:"keys"`
}

func CreateSheet(info *Sheet) error {
	_, err := insertOne(TableSheet, &info)
	return err
}

func GetSheetNextID() uint64 {
	num, _ := getSequenceNext(TableSheet)
	return num
}

func GetSheets() ([]*Sheet, error) {
	cursor, err1 := findAll(TableSheet, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Sheet, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Sheet)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveSheet(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db sheet uid is empty ")
	}
	_, err := removeOne(TableSheet, uid, operator)
	return err
}

func GetSheet(uid string) (*Sheet, error) {
	if len(uid) < 2 {
		return nil, errors.New("db sheet uid is empty of GetSheet")
	}
	result, err := findOne(TableSheet, uid)
	if err != nil {
		return nil, err
	}
	model := new(Sheet)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSheetCount() int64 {
	num, _ := getTotalCount(TableSheet)
	return num
}

func GetSheetByName(owner, name string) (*Sheet, error) {
	if len(owner) < 2 || len(name) < 2 {
		return nil, errors.New("db owner or name is empty of GetSheetByName")
	}
	filter := bson.M{"owner": owner, "name": name, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableSheet, filter)
	if err != nil {
		return nil, err
	}
	model := new(Sheet)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSheetByQuote(owner, quote string, tp uint8) (*Sheet, error) {
	if len(owner) < 2 || len(quote) < 2 {
		return nil, errors.New("db owner or quote is empty of GetSheetByQuote")
	}
	filter := bson.M{"owner": owner, "quote": quote, "product": tp, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableSheet, filter)
	if err != nil {
		return nil, err
	}
	model := new(Sheet)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSheetByQuote3(owner, quote string) (*Sheet, error) {
	if len(owner) < 2 || len(quote) < 2 {
		return nil, errors.New("db owner or quote is empty of GetSheetByQuote")
	}
	filter := bson.M{"owner": owner, "quote": quote, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableSheet, filter)
	if err != nil {
		return nil, err
	}
	model := new(Sheet)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSheetByQuote4(quote string) (*Sheet, error) {
	if len(quote) < 2 {
		return nil, errors.New("db owner or quote is empty of GetSheetByQuote")
	}
	filter := bson.M{"quote": quote, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableSheet, filter)
	if err != nil {
		return nil, err
	}
	model := new(Sheet)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSheetByQuote2(quote string, tp uint8) (*Sheet, error) {
	if len(quote) < 2 {
		return nil, errors.New("db quote is empty of GetSheetByQuote")
	}
	filter := bson.M{"quote": quote, "product": tp, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableSheet, filter)
	if err != nil {
		return nil, err
	}
	model := new(Sheet)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSheetsByQuote(quote string) ([]*Sheet, error) {
	if len(quote) < 2 {
		return nil, errors.New("db quote is empty of GetSheetByQuote")
	}
	def := new(time.Time)
	filter := bson.M{"quote": quote, "deleteAt": def}
	cursor, err1 := findMany(TableSheet, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Sheet, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Sheet)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetSheetsByOwner(owner string) ([]*Sheet, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	cursor, err1 := findMany(TableSheet, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Sheet, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Sheet)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetSheetsByOwnerTP(owner string, tp uint8) ([]*Sheet, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "product": tp, "deleteAt": def}
	cursor, err1 := findMany(TableSheet, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Sheet, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Sheet)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetSheetsByStatus(owner string, st uint8) ([]*Sheet, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "status": st, "deleteAt": def}
	cursor, err1 := findMany(TableSheet, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Sheet, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Sheet)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateSheetBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSheet, uid, msg)
	return err
}

func UpdateSheetState(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSheet, uid, msg)
	return err
}

func UpdateSheetQuote(uid, quote, operator string) error {
	msg := bson.M{"quote": quote, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSheet, uid, msg)
	return err
}

func UpdateSheetKeys(uid, operator string, list []string) error {
	msg := bson.M{"keys": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSheet, uid, msg)
	return err
}

func AppendSheetKey(uid string, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"keys": key}
	_, err := appendElement(TableSheet, uid, msg)
	return err
}

func SubtractSheetKey(uid, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"keys": key}
	_, err := removeElement(TableSheet, uid, msg)
	return err
}
