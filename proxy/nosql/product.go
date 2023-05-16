package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy"
	"time"
)

type Product struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name     string                 `json:"name" bson:"name"`
	Status   uint8                  `json:"status" bson:"status"`
	Type     uint8                  `json:"type" bson:"type"`
	Key      string                 `json:"key" bson:"key"`
	Entries  []string               `json:"entries" bson:"entries"`
	Menus    string                 `json:"menus" bson:"menus"`
	Remark   string                 `json:"remark" bson:"remark"`
	Templet  string                 `json:"templet" bson:"templet"`
	Catalogs string                 `json:"catalogs" bson:"catalogs"`
	Revises  []string               `json:"revises" bson:"revises"`
	Showings []string               `json:"shows" bson:"shows"`
	Effects  []*proxy.ProductEffect `json:"effects" bson:"effects"`
	Displays []*proxy.DisplayShow   `json:"displays" bson:"displays"`
}

func CreateProduct(info *Product) error {
	_, err := insertOne(TableProduct, &info)
	return err
}

func GetProductNextID() uint64 {
	num, _ := getSequenceNext(TableProduct)
	return num
}

func RemoveProduct(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Product uid is empty ")
	}
	_, err := removeOne(TableProduct, uid, operator)
	return err
}

func GetProduct(uid string) (*Product, error) {
	if len(uid) < 2 {
		return nil, errors.New("db activity uid is empty of GetProduct")
	}
	result, err := findOne(TableProduct, uid)
	if err != nil {
		return nil, err
	}
	model := new(Product)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetProducts() ([]*Product, error) {
	cursor, err1 := findAll(TableProduct, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Product, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Product)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetProductCount() int64 {
	num, _ := getTotalCount(TableProduct)
	return num
}

func GetProductByType(tp uint8) (*Product, error) {
	filter := bson.M{"type": tp, "deleteAt": new(time.Time)}
	result, err1 := findOneBy(TableProduct, filter)
	if err1 != nil {
		return nil, err1
	}
	model := new(Product)
	err2 := result.Decode(&model)
	if err2 != nil {
		return nil, err2
	}
	return model, nil
}

func GetProductsByDisplay(display string) ([]*Product, error) {
	filter := bson.M{"displays.uid": display, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableProduct, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Product, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Product)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateProductCatalog(uid, catalogs, operator string) error {
	msg := bson.M{"catalogs": catalogs, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableProduct, uid, msg)
	return err
}

func UpdateProductTemplet(uid, templet, operator string) error {
	msg := bson.M{"templet": templet, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableProduct, uid, msg)
	return err
}

func UpdateProductMenus(uid, menus, operator string) error {
	msg := bson.M{"menus": menus, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableProduct, uid, msg)
	return err
}

func UpdateProductRevises(uid, operator string, arr []string) error {
	msg := bson.M{"revises": arr, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableProduct, uid, msg)
	return err
}

func UpdateProductEntries(uid, operator string, arr []string) error {
	msg := bson.M{"entries": arr, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableProduct, uid, msg)
	return err
}

func UpdateProductEffects(uid, operator string, effects []*proxy.ProductEffect) error {
	msg := bson.M{"effects": effects, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableProduct, uid, msg)
	return err
}

func UpdateProductDisplays(uid, operator string, arr []*proxy.DisplayShow) error {
	msg := bson.M{"displays": arr, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableProduct, uid, msg)
	return err
}

func UpdateProductShows(uid, operator string, arr []string) error {
	msg := bson.M{"shows": arr, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableProduct, uid, msg)
	return err
}

func UpdateProductBase(uid, operator, name, key, remark string) error {
	msg := bson.M{"name": name, "key": key, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableProduct, uid, msg)
	return err
}
