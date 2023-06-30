package proxy

import "time"

type EntityInfo struct {
	UID  string `json:"uid" bson:"uid"`
	Name string `json:"name" bson:"name"`
}

type DateInfo struct {
	Start string `json:"start" bson:"start"`
	Stop  string `json:"stop" bson:"stop"`
}

type DurationInfo struct {
	Start int64 `json:"start" bson:"start"`
	Stop  int64 `json:"stop" bson:"stop"`
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

type DisplayShow struct {
	UID    string `json:"uid" bson:"uid"`
	Effect string `json:"effect" bson:"effect"`
}

type ShowContent struct {
	UID       string `json:"uid" bson:"uid"`       //展览UID
	Weight    uint32 `json:"weight" bson:"weight"` //排序权重
	Effect    string `json:"effect" bson:"effect"` //效果
	Menu      string `json:"menu" bson:"menu"`     //所属目录
	Alignment string `json:"align" bson:"align"`   //目录方向
	Local     uint32 `json:"local" bson:"local"`   //是否本地展览
}

type DisplayContent struct {
	UID string `json:"uid" bson:"uid"` //实体UID
	//Entity string   `json:"entity" bson:"entity"` //实体UID
	Events []string `json:"events" bson:"events"` //
	Assets []string `json:"assets" bson:"assets"` //
}

type ProductEffect struct {
	Pattern string
	Min     uint32
	Max     uint32
}

func (mine *DateInfo) BeginUTC() int64 {
	return DateToUTC(mine.Start)
}

func (mine *DateInfo) EndUTC() int64 {
	return DateToUTC(mine.Stop)
}

func (mine *DurationInfo) Begin() string {
	return UTCToDate(mine.Start)
}

func (mine *DurationInfo) End() string {
	return UTCToDate(mine.Stop)
}

func DateToUTC(date string) int64 {
	if date == "" {
		return 0
	}
	t, e := time.ParseInLocation("2006/01/02", date, time.Local)
	if e != nil {
		return 0
	}
	return t.Unix()
}

func UTCToDate(utc int64) string {
	if utc < 1 {
		return ""
	}
	t := time.Unix(utc, 0)
	return t.Format("2006/01/02")
}
