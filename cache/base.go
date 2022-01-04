package cache

import (
	"github.com/micro/go-micro/v2/logger"
	"omo.msa.favorite/config"
	"omo.msa.favorite/proxy/nosql"
	"reflect"
	"time"
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
	//boxes []*OwnerInfo
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if err == nil {
		table := getFavoriteTable(true)
		num := nosql.GetFavoriteCount(table)
		count := nosql.GetRepertoryCount()
		logger.Infof("the person favorite count = %d and the repertory count = %d", num, count)
	}
	return err
}

func Context() *cacheContext {
	return cacheCtx
}

func checkPage( page, number uint32, all interface{}) (uint32, uint32, interface{}) {
	if number < 1 {
		number = 10
	}
	array := reflect.ValueOf(all)
	total := uint32(array.Len())
	maxPage := total/number + 1
	if page < 1 {
		return total, maxPage, all
	}

	var start = (page - 1) * number
	var end = start + number
	if end > total {
		end = total
	}

	list := array.Slice(int(start), int(end))
	return total, maxPage, list.Interface()
}
