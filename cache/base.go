package cache

import (
	"github.com/micro/go-micro/v2/logger"
	"omo.msa.favorite/config"
	"omo.msa.favorite/proxy/nosql"
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
		num := nosql.GetFavoriteCount()
		count := nosql.GetRepertoryCount()
		logger.Infof("the favorite count = %d and the repertory count = %d", num, count)
	}
	return err
}

func Context() *cacheContext {
	return cacheCtx
}
