package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Article struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Status   uint8  `json:"status" bson:"status"`
	Subtitle string `json:"subtitle" bson:"subtitle"`
	Body     string `json:"body" bson:"body"`
	Owner    string `json:"owner" bson:"owner"`
	Type     uint8  `json:"type" bson:"type"`

	Assets  []string `json:"assets" bson:"assets"`
	Targets []string `json:"targets" bson:"targets"`
	Tags    []string `json:"tags" bson:"tags"`
}

func CreateArticle(info *Article) error {
	_, err := insertOne(TableArticle, &info)
	return err
}

func GetArticleNextID() uint64 {
	num, _ := getSequenceNext(TableArticle)
	return num
}

func RemoveArticle(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db article uid is empty ")
	}
	_, err := removeOne(TableArticle, uid, operator)
	return err
}

func GetArticle(uid string) (*Article, error) {
	if len(uid) < 2 {
		return nil, errors.New("db article uid is empty of GetArticle")
	}
	result, err := findOne(TableArticle, uid)
	if err != nil {
		return nil, err
	}
	model := new(Article)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetArticleCount() int64 {
	num, _ := getTotalCount(TableArticle)
	return num
}

func GetArticlesByOwnerTP(owner string, kind uint8) ([]*Article, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "type": kind, "deleteAt": def}
	cursor, err1 := findMany(TableArticle, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Article, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Article)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetArticlesByOwnerStatus(owner string, st uint8) ([]*Article, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "status": st, "deleteAt": def}
	cursor, err1 := findMany(TableArticle, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Article, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Article)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetArticlesByType(kind uint8) ([]*Article, error) {
	def := new(time.Time)
	filter := bson.M{"type": kind, "deleteAt": def}
	cursor, err1 := findMany(TableArticle, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Article, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Article)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetArticlesByOwner(owner string) ([]*Article, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	cursor, err1 := findMany(TableArticle, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Article, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Article)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetArticlesByTargets(st uint8, targets []string) ([]*Article, error) {
	def := new(time.Time)
	in := bson.A{}
	for _, target := range targets {
		in = append(in, target)
	}
	filter := bson.M{"status":st, "$or":bson.A{bson.M{"targets": bson.M{"$in":in}}, bson.M{"targets":bson.M{"$ne":nil}}} , "deleteAt": def}
	cursor, err1 := findMany(TableArticle, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Article, 0, 20)
	for cursor.Next(context.TODO()) {
		var node = new(Article)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetArticlesByOTargets(owner string, st uint8, targets []string) ([]*Article, error) {
	def := new(time.Time)
	in := bson.A{}
	for _, target := range targets {
		in = append(in, target)
	}
	filter := bson.M{"owner":owner,"status":st, "$or":bson.A{bson.M{"targets": bson.M{"$in":in}}} , "deleteAt": def}
	cursor, err1 := findMany(TableArticle, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Article, 0, 20)
	for cursor.Next(context.TODO()) {
		var node = new(Article)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateArticleBase(uid, name, sub, body, operator string) error {
	msg := bson.M{"name": name, "body": body, "subtitle": sub, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArticle, uid, msg)
	return err
}

func UpdateArticleStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArticle, uid, msg)
	return err
}

func UpdateArticleTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArticle, uid, msg)
	return err
}

func UpdateArticleAssets(uid, operator string, list []string) error {
	msg := bson.M{"assets": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArticle, uid, msg)
	return err
}

func UpdateArticleTargets(uid, operator string, list []string) error {
	msg := bson.M{"targets": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableArticle, uid, msg)
	return err
}
