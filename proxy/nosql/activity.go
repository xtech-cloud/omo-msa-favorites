package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	Cover        string    `json:"cover" bson:"cover"`
	Remark       string    `json:"remark" bson:"remark"`
	Require      string `json:"require" bson:"require"`
	Owner        string    `json:"owner" bson:"owner"`
	Type         uint8     `json:"type" bson:"type"`
	Limit		 uint8 		`json:"limit" bson:"limit"`
	Organizer    string    `json:"organizer" bson:"organizer"`
	Place        proxy.PlaceInfo `json:"place" bson:"place"`
	Date         proxy.DateInfo  `json:"date" bson:"date"`
	Tags         []string  `json:"tags" bsonL:"tags"`
	Assets       []string  `json:"assets" bson:"assets"`
	Targets      []string `json:"targets" bson:"targets"`
	Participants []string  `json:"participants" bson:"participants"`
	Persons []proxy.PersonInfo  `json:"persons" bson:"persons"`
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
	num, _ := getCount(TableActivity)
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

func UpdateActivityBase(uid, name, remark, require, operator string, date proxy.DateInfo, place proxy.PlaceInfo) error {
	msg := bson.M{"name": name, "remark": remark, "require": require, "operator": operator, "date":date, "place":place, "updatedAt": time.Now()}
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

func AppendActivityParticipant(uid string, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"participants": key}
	_, err := appendElement(TableActivity, uid, msg)
	return err
}

func SubtractActivityParticipant(uid, key string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"participants": key}
	_, err := removeElement(TableActivity, uid, msg)
	return err
}

func AppendActivityPerson(uid string, person proxy.PersonInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"persons": person}
	_, err := appendElement(TableActivity, uid, msg)
	return err
}

func SubtractActivityPerson(uid, entity string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"persons": bson.M{"entity":entity}}
	_, err := removeElement(TableActivity, uid, msg)
	return err
}
