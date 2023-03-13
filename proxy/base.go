package proxy

type EntityInfo struct {
	UID  string `json:"uid" bson:"uid"`
	Name string `json:"name" bson:"name"`
}

type DateInfo struct {
	Start string `json:"start" bson:"start"`
	Stop  string `json:"stop" bson:"stop"`
}

type PlaceInfo struct {
	Name     string `json:"name" bson:"name"`
	Location string `json:"location" bson:"location"`
}

type PersonInfo struct {
	Entity string `json:"entity" bson:"entity"`
	Event  string `json:"event" bson:"event"`
}

// 奖项评选要求
type PrizeInfo struct {
	Name  string     `json:"name" bson:"name"`   // 奖项名称
	Desc  string     `json:"desc" bson:"desc"`   // 描述
	Ranks []RankInfo `json:"ranks" bson:"ranks"` //奖项名称
}

// 奖项名称
type RankInfo struct {
	Index uint32 `json:"index" bson:"index"` // 第几名
	Name  string `json:"name" bson:"name"`   //名次的名称
	Count uint32 `json:"count" bson:"count"` //评选数量
}

// 作品信息
type OpusInfo struct {
	Rank   uint32 `json:"rank" bson:"rank"`
	Asset  string `json:"asset" bson:"asset"`   //
	Remark string `json:"remark" bson:"remark"` //评语
}

//type ShowingInfo struct {
//	//场所
//	string    `json:"target" bson:"target"` //场所
//	Effect    string                        `json:"effect" bson:"effect"` //展览的板式
//	Menu      string                        `json:"menu" bson:"menu"`     //所属目录
//	Alignment string                        `json:"align" bson:"align"`   //目录方向
//	Slots     []string                      `json:"slots" bson:"slots"`
//	UpdatedAt time.Time                     `json:"updatedAt" bson:"updatedAt"`
//}

type ShowContent struct {
	UID       string `json:"uid" bson:"uid"`       //展览UID
	Weight    uint32 `json:"weight" bson:"weight"` //排序权重
	Effect    string `json:"effect" bson:"effect"` //效果
	Menu      string `json:"menu" bson:"menu"`     //所属目录
	Alignment string `json:"align" bson:"align"`   //目录方向
}

type DisplayContent struct {
	UID string `json:"uid" bson:"uid"` //
	//Entity string   `json:"entity" bson:"entity"` //实体UID
	Events []string `json:"events" bson:"events"` //
	Assets []string `json:"assets" bson:"assets"` //
}

type ProductEffect struct {
	Pattern string
	Min     uint32
	Max     uint32
}
