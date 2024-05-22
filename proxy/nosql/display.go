package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"omo.msa.favorite/proxy"
	"time"
)

type Display struct {
	UID         primitive.ObjectID     `bson:"_id"`
	ID          uint64                 `json:"id" bson:"id"`
	Name        string                 `json:"name" bson:"name"`
	CreatedTime time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time              `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time              `json:"deleteAt" bson:"deleteAt"`
	Creator     string                 `json:"creator" bson:"creator"`
	Operator    string                 `json:"operator" bson:"operator"`
	Cover       string                 `json:"cover" bson:"cover"`
	Remark      string                 `json:"remark" bson:"remark"`
	Owner       string                 `json:"owner" bson:"owner"`
	Status      uint8                  `json:"status" bson:"status"`
	Type        uint8                  `json:"type" bson:"type"`
	Access      uint8                  `json:"access" bson:"access"`
	Origin      string                 `json:"origin" bson:"origin"` //数据来源，可能是某次活动,或者是标准榜样
	Meta        string                 `json:"meta" bson:"meta"`     //源数据
	Banner      string                 `json:"banner" bson:"banner"`
	Poster      string                 `json:"poster" bson:"poster"`
	Tags        []string               `json:"tags" bsonL:"tags"`
	Keys        []string               `json:"keys" bson:"keys"`
	Scenes      []uint32               `json:"scenes" bson:"scenes"`
	Contents    []proxy.DisplayContent `json:"contents" bson:"contents"`
	Pending     []proxy.DisplayContent `json:"pending" bson:"pending"`
}

func CreateDisplay(info *Display) error {
	_, err := insertOne(TableDisplay, &info)
	return err
}

func GetDisplayNextID() uint64 {
	num, _ := getSequenceNext(TableDisplay)
	return num
}

func GetDisplays() ([]*Display, error) {
	cursor, err1 := findAll(TableDisplay, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Display, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Display)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDisplayFile(uid string) (*FileInfo, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Display.files uid is empty ")
	}
	result, err := findOne(TableDisplay, uid)
	if err != nil {
		return nil, err
	}
	info := new(FileInfo)
	err1 := result.Decode(&info)
	return info, err1
}

func RemoveDisplay(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Display uid is empty ")
	}
	_, err := removeOne(TableDisplay, uid, operator)
	return err
}

func GetDisplay(uid string) (*Display, error) {
	if len(uid) < 2 {
		return nil, errors.New("db Display uid is empty of GetDisplay")
	}
	result, err := findOne(TableDisplay, uid)
	if err != nil {
		return nil, err
	}
	model := new(Display)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetDisplayCount() int64 {
	num, _ := getTotalCount(TableDisplay)
	return num
}

func GetDisplaysCount(st uint8) int64 {
	def := new(time.Time)
	filter := bson.M{"status": st, "deleteAt": def}
	num, err1 := getCount(TableDisplay, filter)
	if err1 != nil {
		return num
	}

	return num
}

func GetDisplayByOrigin(owner, origin string) (*Display, error) {
	if len(origin) < 2 || len(owner) < 2 {
		return nil, errors.New("db origin uid is empty of GetDisplay")
	}
	filter := bson.M{"owner": owner, "origin": origin, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableDisplay, filter)
	if err != nil {
		return nil, err
	}
	model := new(Display)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetDisplayByName(owner, name string, tp uint8) (*Display, error) {
	if len(owner) < 2 || len(name) < 2 {
		return nil, errors.New("db owner or name is empty of GetDisplayByName")
	}
	filter := bson.M{"owner": owner, "name": name, "type": tp, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableDisplay, filter)
	if err != nil {
		return nil, err
	}
	model := new(Display)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetDisplaysByOwner(owner string) ([]*Display, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	cursor, err1 := findMany(TableDisplay, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Display, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Display)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDisplaysByContent(owner, key string) ([]*Display, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "contents.uid": key, "deleteAt": def}
	cursor, err1 := findMany(TableDisplay, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Display, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Display)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDisplaysByContent2(key string) ([]*Display, error) {
	def := new(time.Time)
	filter := bson.M{"contents.uid": key, "deleteAt": def}
	cursor, err1 := findMany(TableDisplay, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Display, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Display)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDisplaysByStatus(owner string, st uint8) ([]*Display, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "status": st, "deleteAt": def}
	cursor, err1 := findMany(TableDisplay, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Display, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Display)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDisplaysByOwnerTP(owner string, kind uint8) ([]*Display, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "type": kind, "deleteAt": def}
	cursor, err1 := findMany(TableDisplay, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Display, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Display)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDisplaysByType(kind uint8) ([]*Display, error) {
	def := new(time.Time)
	filter := bson.M{"type": kind, "deleteAt": def}
	cursor, err1 := findMany(TableDisplay, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Display, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Display)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDisplaysByTarget(owner, target string) ([]*Display, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "targets.target": target, "deleteAt": def}
	cursor, err1 := findMany(TableDisplay, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Display, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Display)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDisplaysByPage(st uint8, start, num int64) ([]*Display, error) {
	def := new(time.Time)
	filter := bson.M{"status": st, "deleteAt": def}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}}).SetLimit(num).SetSkip(start)
	cursor, err1 := findManyByOpts(TableDisplay, filter, opts)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Display, 0, 20)
	for cursor.Next(context.TODO()) {
		var node = new(Display)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

//func GetDisplaysByTargets(st uint8, targets []string) ([]*Display, error) {
//	def := new(time.Time)
//	in := bson.A{}
//	for _, target := range targets {
//		in = append(in, target)
//	}
//	filter := bson.M{"status": st, "$or": bson.A{bson.M{"targets": bson.M{"$in": in}}}, "deleteAt": def}
//	cursor, err1 := findMany(TableDisplay, filter, 0)
//	if err1 != nil {
//		return nil, err1
//	}
//	var items = make([]*Display, 0, 20)
//	for cursor.Next(context.Background()) {
//		var node = new(Display)
//		if err := cursor.Decode(&node); err != nil {
//			return nil, err
//		} else {
//			items = append(items, node)
//		}
//	}
//	return items, nil
//}

func UpdateDisplayBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDisplay, uid, msg)
	return err
}

func UpdateDisplayCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDisplay, uid, msg)
	return err
}

func UpdateDisplayBanner(uid, cover, operator string) error {
	msg := bson.M{"banner": cover, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDisplay, uid, msg)
	return err
}

func UpdateDisplayPoster(uid, cover, operator string) error {
	msg := bson.M{"poster": cover, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDisplay, uid, msg)
	return err
}

func UpdateDisplayState(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDisplay, uid, msg)
	return err
}

func UpdateDisplayAccess(uid, operator string, st uint8) error {
	msg := bson.M{"access": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDisplay, uid, msg)
	return err
}

func UpdateDisplayOwner(uid, owner, operator string) error {
	msg := bson.M{"owner": owner, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDisplay, uid, msg)
	return err
}

func UpdateDisplayType(uid, operator string, st uint8) error {
	msg := bson.M{"type": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDisplay, uid, msg)
	return err
}

func UpdateDisplayTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDisplay, uid, msg)
	return err
}

func UpdateDisplayContents(uid, operator string, arr []proxy.DisplayContent) error {
	msg := bson.M{"contents": arr, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDisplay, uid, msg)
	return err
}

func UpdateDisplayPending(uid, operator string, arr []proxy.DisplayContent) error {
	msg := bson.M{"pending": arr, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableDisplay, uid, msg)
	return err
}

func AppendDisplayContent(uid, operator string, content proxy.DisplayContent) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"contents": content}
	_, err := appendElement(TableDisplay, uid, operator, msg)
	return err
}

func SubtractDisplayContent(uid, key, operator string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"contents": bson.M{"uid": key}}
	_, err := removeElement(TableDisplay, uid, operator, msg)
	return err
}

func AppendDisplayPending(uid, operator string, content proxy.DisplayContent) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"pending": content}
	_, err := appendElement(TableDisplay, uid, operator, msg)
	return err
}

func SubtractDisplayPending(uid, key, operator string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"pending": bson.M{"uid": key}}
	_, err := removeElement(TableDisplay, uid, operator, msg)
	return err
}

//func UpdateDisplayTargets(uid, operator string, list []*proxy.ShowingInfo) error {
//	msg := bson.M{"targets": list, "operator": operator, "updatedAt": time.Now()}
//	_, err := updateOne(TableDisplay, uid, msg)
//	return err
//}

//func AppendDisplayTarget(uid string, target *proxy.ShowingInfo) error {
//	if len(uid) < 1 {
//		return errors.New("the uid is empty")
//	}
//	msg := bson.M{"targets": target}
//	_, err := appendElement(TableDisplay, uid, msg)
//	return err
//}

//func SubtractDisplayTarget(uid, target string) error {
//	if len(uid) < 1 {
//		return errors.New("the uid is empty")
//	}
//	msg := bson.M{"targets": bson.M{"target": target}}
//	_, err := removeElement(TableDisplay, uid, msg)
//	return err
//}
