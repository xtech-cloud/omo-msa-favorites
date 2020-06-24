package proxy


type EntityInfo struct {
	UID string `json:"uid" bson:"uid"`
	Name string `json:"name" bson:"name"`
}