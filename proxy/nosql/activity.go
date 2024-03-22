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

type Activity struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Cover   string `json:"cover" bson:"cover"`     //
	Remark  string `json:"remark" bson:"remark"`   //活动简介
	Require string `json:"require" bson:"require"` //活动要求
	Owner   string `json:"owner" bson:"owner"`
	Poster  string `json:"poster" bson:"poster"` //海报
	Type    uint8  `json:"type" bson:"type"`
	Limit   uint8  `json:"limit" bson:"limit"` //限制用户上传图片数量
	Status  uint8  `json:"status" bson:"status"`
	Access  uint8  `json:"access" bson:"access"` //访问权限
	Show    uint8  `json:"show" bson:"show"`     //显示结果

	Participant uint32 `json:"participant" bson:"participant"` //活动参与人数

	Organizer    string             `json:"organizer" bson:"organizer"`
	Template     string             `json:"template" bson:"template"`
	Certify      proxy.CertifyInfo  `json:"certify" bson:"certify"`
	Place        proxy.PlaceInfo    `json:"place" bson:"place"`
	Date         proxy.DateInfo     `json:"date" bson:"date"`
	Duration     proxy.DurationInfo `json:"duration" bson:"duration"`
	Prize        *proxy.PrizeInfo   `json:"prize" bson:"prize"` //奖项设置
	Tags         []string           `json:"tags" bsonL:"tags"`
	Assets       []string           `json:"assets" bson:"assets"`
	Targets      []string           `json:"targets" bson:"targets"` //引用的实体对象
	Quotes       []string           `json:"quotes" bson:"quotes"`
	Participants []string           `json:"participants" bson:"participants"` //弃用
	Persons      []proxy.PersonInfo `json:"persons" bson:"persons"`           //记录参与人信息
	Opuses       []proxy.OpusInfo   `json:"opuses" bson:"opuses"`             //获奖作品
}

func CreateActivity(info *Activity) error {
	_, err := insertOne(TableActivity, &info)
	return err
}

func GetActivityNextID() uint64 {
	num, _ := getSequenceNext(TableActivity)
	return num
}

func GetActivities() ([]*Activity, error) {
	cursor, err1 := findAll(TableActivity, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Activity)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveActivity(uid, operator string) error {
	if len(uid) < 2 {
		return errors.New("db Activity uid is empty ")
	}
	_, err := removeOne(TableActivity, uid, operator)
	return err
}

func GetActivity(uid string) (*Activity, error) {
	if len(uid) < 2 {
		return nil, errors.New("db activity uid is empty of GetActivity")
	}
	result, err := findOne(TableActivity, uid)
	if err != nil {
		return nil, err
	}
	model := new(Activity)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetActivityCount() int64 {
	num, _ := getTotalCount(TableActivity)
	return num
}

func GetActivityCountByOwner(owner string) int64 {
	filter := bson.M{"owner": owner, "$or": bson.A{bson.M{"template": ""}, bson.M{"template": bson.M{"$exists": false}}}, "deleteAt": new(time.Time)}
	num, _ := getCount(TableActivity, filter)
	return num
}

func GetActivityCountByStatus(owner string, st uint8) int64 {
	filter := bson.M{"owner": owner, "status": st, "deleteAt": new(time.Time)}
	num, _ := getCount(TableActivity, filter)
	return num
}

func GetActivityCountByClone(owner string) int64 {
	filter := bson.M{"owner": owner, "template": bson.M{"$exists": true, "$ne": ""}, "deleteAt": new(time.Time)}
	num, _ := getCount(TableActivity, filter)
	return num
}

func GetActivityByOrganizer(organizer string) ([]*Activity, error) {
	def := new(time.Time)
	filter := bson.M{"organizer": organizer, "deleteAt": def}
	cursor, err1 := findMany(TableActivity, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesByOwner(owner string) ([]*Activity, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "deleteAt": def}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}}).SetLimit(10)
	cursor, err1 := findManyByOpts(TableActivity, filter, opts)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.TODO()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesByTargets(st uint8, targets []string) ([]*Activity, error) {
	def := new(time.Time)
	in := bson.A{}
	for _, target := range targets {
		in = append(in, target)
	}
	filter := bson.M{"status": st, "$or": bson.A{bson.M{"targets": bson.M{"$in": in}}}, "deleteAt": def}
	cursor, err1 := findMany(TableActivity, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.TODO()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesByOTargets(owner string, st uint8, targets []string) ([]*Activity, error) {
	def := new(time.Time)
	in := bson.A{}
	for _, target := range targets {
		in = append(in, target)
	}
	//filter := bson.M{"owner":owner, "status":st, "$or":bson.A{bson.M{"targets": bson.M{"$in":in}},bson.M{"targets":bson.M{"$ne":nil}}} , "deleteAt": def}
	filter := bson.M{"owner": owner, "status": st, "$or": bson.A{bson.M{"targets": bson.M{"$in": in}}}, "deleteAt": def}
	cursor, err1 := findMany(TableActivity, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.TODO()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesByStates(owner string, states []uint8) ([]*Activity, error) {
	def := new(time.Time)
	in := bson.A{}
	for _, st := range states {
		in = append(in, bson.M{"status": st})
	}
	filter := bson.M{"owner": owner, "$or": in, "deleteAt": def}
	cursor, err1 := findMany(TableActivity, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.TODO()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesByStatus(owner string, status uint8) ([]*Activity, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "status": status, "deleteAt": def}
	cursor, err1 := findMany(TableActivity, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesByStatusTP(owner string, status, tp uint8) ([]*Activity, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "type": tp, "status": status, "deleteAt": def}
	cursor, err1 := findMany(TableActivity, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesByType(owner string, tp uint8) ([]*Activity, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "type": tp, "deleteAt": def}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}}).SetLimit(100)
	cursor, err1 := findManyByOpts(TableActivity, filter, opts)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesByQuote(quote string) ([]*Activity, error) {
	def := new(time.Time)
	filter := bson.M{"quotes": quote, "deleteAt": def}
	cursor, err1 := findMany(TableActivity, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesBySceneType(scene string, tp uint8) ([]*Activity, error) {
	def := new(time.Time)
	filter := bson.M{"owner": scene, "type": tp, "deleteAt": def}
	cursor, err1 := findMany(TableActivity, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesByShow(owners []string, st uint8) ([]*Activity, error) {
	def := new(time.Time)
	in := bson.A{}
	for _, ss := range owners {
		in = append(in, bson.M{"owner": ss})
	}
	filter := bson.M{"$or": bson.A{in}, "show": st, "template": "", "deleteAt": def}
	cursor, err1 := findMany(TableActivity, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesByTemplate(template string) ([]*Activity, error) {
	def := new(time.Time)
	filter := bson.M{"template": template, "deleteAt": def}
	cursor, err1 := findMany(TableActivity, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetActivitiesByOwnTemplate(owner, template string) ([]*Activity, error) {
	def := new(time.Time)
	filter := bson.M{"owner": owner, "template": template, "deleteAt": def}
	cursor, err1 := findMany(TableActivity, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Activity, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Activity)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateActivityDuration(uid string, date proxy.DurationInfo) error {
	msg := bson.M{"duration": date, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityBase(uid, name, remark, require, operator string, date *proxy.DurationInfo, place proxy.PlaceInfo) error {
	msg := bson.M{"name": name, "remark": remark, "require": require, "operator": operator, "duration": date, "place": place, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityCover(uid, cover, operator string) error {
	msg := bson.M{"cover": cover, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityLimit(uid, operator string, num uint8) error {
	msg := bson.M{"limit": num, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityStatus(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityParticipant(uid string, count uint32) error {
	msg := bson.M{"participant": count, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityShowState(uid, operator string, st uint8) error {
	msg := bson.M{"show": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityAccess(uid, operator string, st uint8) error {
	msg := bson.M{"access": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityCertify(uid, operator string, certify proxy.CertifyInfo) error {
	msg := bson.M{"certify": certify, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityPoster(uid, operator, poster string) error {
	msg := bson.M{"poster": poster, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityStop(uid, operator string, stop int64) error {
	msg := bson.M{"duration.stop": stop, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityPrize(uid, operator string, prize *proxy.PrizeInfo) error {
	msg := bson.M{"prize": prize, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityOpuses(uid, operator string, list []proxy.OpusInfo) error {
	msg := bson.M{"opuses": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityAssets(uid, operator string, list []string) error {
	msg := bson.M{"assets": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}

func UpdateActivityTargets(uid, operator string, list []string) error {
	msg := bson.M{"targets": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableActivity, uid, msg)
	return err
}
