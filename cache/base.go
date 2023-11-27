package cache

import (
	"errors"
	"github.com/micro/go-micro/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/config"
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
	"strconv"
	"strings"
	"time"
)

const (
	RecodeActivity = 1
	RecodeNotice   = 2
	RecodeFav      = 3
	RecodeArticle  = 4
	RecodeClick    = 5
)

const (
	MessageActivity = 1
	MessageNotice   = 99
)

const (
	HistoryActivity HistoryType = 1
	HistoryDisplay  HistoryType = 2
)

const (
	OptionAgree  OptionType = 1 //同意
	OptionRefuse OptionType = 2 //拒绝
	OptionSwitch OptionType = 3 //切换关联
)

const (
	LogOptNull       uint32 = 0
	LogOptRequestAdd uint32 = 1
	LogOptRequestDel uint32 = 2
	LogOptAgreeAdd   uint32 = 3
	LogOptAgreeDel   uint32 = 4
	LogOptRefuseAdd  uint32 = 5
	LogOptRefuseDel  uint32 = 6
	LogOptAgree      uint32 = 7
	LogOptRefuse     uint32 = 8
	LogOptPend       uint32 = 9 //
)

type OptionType uint8

type HistoryType uint8

type BaseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator    string
	Operator   string
	CreateTime time.Time
	UpdateTime time.Time
}

type cacheContext struct {
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}

	//var cstZone = time.FixedZone("CST", 8*3600) // 东八
	//time.Local = cstZone

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if err == nil {
		num := nosql.GetDisplayCount()
		count := nosql.GetFavoriteCount()
		logger.Infof("the person favorite count = %d and the display count = %d", count, num)

		//db, _ := nosql.GetActivity("616e7f56e1fd51b21c857b26")
		//if db != nil {
		//	info := new(ActivityInfo)
		//	info.initInfo(db)
		//}
		nosql.MoveTable()
	}
	return err
}

func Context() *cacheContext {
	return cacheCtx
}

//func CheckPage(page, number uint32, all interface{}) (uint32, uint32, interface{}) {
//	if number < 1 {
//		number = 10
//	}
//	array := reflect.ValueOf(all)
//	total := uint32(array.Len())
//	maxPage := total/number + 1
//	if page < 1 {
//		return total, maxPage, all
//	}
//
//	if page > maxPage {
//		page = maxPage
//	}
//	var start = (page - 1) * number
//	var end = start + number
//	if end > total {
//		end = total
//	}
//
//	list := array.Slice(int(start), int(end))
//	return total, maxPage, list.Interface()
//}

func CheckPage[T any](page, number uint32, all []T) (uint32, uint32, []T) {
	if len(all) < 1 {
		return 0, 0, make([]T, 0, 1)
	}
	if number < 1 {
		number = 10
	}
	total := uint32(len(all))
	if len(all) <= int(number) {
		return total, 1, all
	}
	maxPage := total/number + 1
	if page < 1 {
		return total, maxPage, all
	}

	var start = (page - 1) * number
	var end = start + number
	if end > total {
		end = total
	}
	list := make([]T, 0, number)
	list = append(all[start:end])
	return total, maxPage, list
}

func (mine *cacheContext) insertHistory(parent, operator, remark, content, from, to string, opt uint32, tp HistoryType) error {
	db := new(nosql.History)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetRecordNextID()
	db.Creator = operator
	db.CreatedTime = time.Now()
	db.Parent = parent
	db.From = from
	db.To = to
	db.Content = content
	db.Option = opt
	db.Type = uint8(tp)
	db.Remark = remark
	return nosql.CreateHistory(db)
}

func ParseDate(msg string) (year int, month time.Month, day int, err error) {
	if len(msg) < 1 {
		return 0, 0, 0, errors.New("the date is empty")
	}
	array := strings.Split(msg, "/")
	if array != nil && len(array) > 2 {
		y, _ := strconv.ParseUint(array[0], 10, 32)
		year = int(y)
		m, _ := strconv.ParseUint(array[1], 10, 32)
		month = time.Month(m)
		d, _ := strconv.ParseUint(array[2], 10, 32)
		day = int(d)
		return year, month, day, nil
	} else {
		return 0, 0, 0, errors.New("the split array is nil")
	}
}

func ParseDate2(msg string) (time.Time, error) {
	t, er := time.ParseInLocation("2006/01/02", msg, time.Local)
	if er == nil {
		return t, nil
	}
	y, m, d, e := ParseDate(msg)
	if e != nil {
		return time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local), e
	}
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC), nil
}

func ParseTime(msg string) int64 {
	y, m, d, e := ParseDate(msg)
	if e != nil {
		return 0
	}
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix()
}

func SwitchDate(start, stop string) (int64, int64) {
	//var cstZone = time.FixedZone("CST", 8*3600) // 东八
	//time.Local = cstZone
	var begin = proxy.DateToUTC(start, 0)
	var end int64 = 0
	var to time.Time
	now := time.Now()
	if begin < 2 { //没有正确的开始日期
		begin = now.Unix()
		to = now.AddDate(0, 0, 1)
		end = to.Unix()
	} else {
		end = proxy.DateToUTC(stop, 1)
		if end < 2 { //没有正确的结束日期
			from, _ := time.ParseInLocation("2006/01/02", start, time.Local)
			to = from.AddDate(0, 0, 1)
			end = to.Unix()
		}
	}
	return begin, end
}
