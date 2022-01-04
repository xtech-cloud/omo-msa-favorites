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

type PersonInfo struct {
	Entity string `json:"entity" bson:"entity"`
	Event string `json:"event" bson:"event"`
}

type ShowingInfo struct {
	//设备
	Target string `json:"target" bson:"target"`
	Effect string `json:"effect" bson:"effect"`
	Skin string `json:"skin" bson:"skin"`
	Slots []string `json:"slots" bson:"slots"`
}