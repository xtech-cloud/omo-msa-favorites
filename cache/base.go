package cache

import (
	"errors"
	"github.com/micro/go-micro/v2/logger"
	"omo.msa.favorite/config"
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
	"strconv"
	"strings"
	"time"
)

const (
	ObserveActivity = 1
	ObserveNotice   = 2
	ObserveFav      = 3
	ObserveArticle  = 4
	ObserveClick    = 5
)

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

		db, _ := nosql.GetActivity("616e7f56e1fd51b21c857b26")
		if db != nil {
			info := new(ActivityInfo)
			info.initInfo(db)
		}
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
	t, er := time.ParseInLocation("2006/01/02", msg, time.UTC)
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
			from, _ := time.ParseInLocation("2006/01/02", start, time.UTC)
			to = from.AddDate(0, 0, 1)
			end = to.Unix()
		}
	}
	return begin, end
}
