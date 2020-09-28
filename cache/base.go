package cache

import (
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
	boxes []*OwnerInfo
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}
	cacheCtx.boxes = make([]*OwnerInfo, 0, 100)

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if err == nil {
		list,er := nosql.GetFavorites()
		if er == nil {
			for _, item := range list {
				owner := GetOwner(item.Owner)
				owner.addFavorite(item)
			}
		}
	}
	return err
}
