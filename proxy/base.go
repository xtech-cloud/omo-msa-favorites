package proxy


type EntityInfo struct {
	UID string `json:"uid" bson:"uid"`
	Name string `json:"name" bson:"name"`
}

type DateInfo struct {
	Start string `json:"start" bson:"start"`
	Stop string `json:"stop" bson:"stop"`
}

type PlaceInfo struct {
	Name string `json:"name" bson:"name"`
	Location string `json:"location" bson:"location"`
}