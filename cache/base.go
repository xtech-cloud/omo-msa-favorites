package cache

import (
	"errors"
	"github.com/micro/go-micro/v2/logger"
	"omo.msa.favorite/config"
	"omo.msa.favorite/proxy/nosql"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	ObserveActivity = 1
	ObserveNotice = 2
	ObserveFav = 3
	ObserveArticle = 4
)

type BaseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator string
	Operator string
	CreateTime time.Time
	UpdateTime time.Time
}

type cacheContext struct {

}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if err == nil {
		num := nosql.GetDisplayCount()
		count := nosql.GetFavoriteCount()
		logger.Infof("the person favorite count = %d and the display count = %d", count, num)

		db,_ := nosql.GetActivity("616e7f56e1fd51b21c857b26")
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

func CheckPage( page, number uint32, all interface{}) (uint32, uint32, interface{}) {
	if number < 1 {
		number = 10
	}
	array := reflect.ValueOf(all)
	total := uint32(array.Len())
	maxPage := total/number + 1
	if page < 1 {
		return total, maxPage, all
	}

	if page > maxPage {
		page = maxPage
	}
	var start = (page - 1) * number
	var end = start + number
	if end > total {
		end = total
	}

	list := array.Slice(int(start), int(end))
	return total, maxPage, list.Interface()
}

func ParseDate(msg string) (year int, month time.Month, day int, err error) {
	if len(msg) < 1 {
		return 0,0,0, errors.New("the date is empty")
	}
	array := strings.Split(msg, "/")
	if array != nil && len(array) > 2 {
		y,_ := strconv.ParseUint(array[0], 10, 32)
		year = int(y)
		m,_ := strconv.ParseUint(array[1], 10, 32)
		month = time.Month(m)
		d,_ := strconv.ParseUint(array[2], 10, 32)
		day = int(d)
		return year,month, day, nil
	}else{
		return 0,0,0,errors.New("the split array is nil")
	}
}

func ParseTime(msg string) int64 {
	y,m,d,e := ParseDate(msg)
	if e != nil {
		return 0
	}
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix()
}
