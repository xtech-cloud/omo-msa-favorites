package cache

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.favorite/proxy/nosql"
	"omo.msa.favorite/tool"
	"time"
)

const ArticleTypeDef = 0

const (
	MessageStatusDraft MessageStatus = 0
	MessageStatusCheck MessageStatus = 1
	MessageStatusRefuse  MessageStatus = 2
	MessageStatusAgree  MessageStatus = 3
)

type MessageStatus uint8

type ArticleInfo struct {
	BaseInfo
	Status   MessageStatus
	Type     uint8  //
	Owner    string //该展览所属组织机构，scene, class等
	Subtitle string
	Body     string
	Tags     []string
	Targets  []string //班级，场景
	Assets   []string
}

func (mine *cacheContext)CreateArticle(info *ArticleInfo) error {
	db := new(nosql.Article)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetArticleNextID()
	db.CreatedTime = time.Now()
	db.Subtitle = info.Subtitle
	db.Name = info.Name
	db.Body = info.Body
	db.Owner = info.Owner
	db.Type = info.Type
	db.Status = uint8(info.Status)
	db.Creator = info.Creator
	db.Operator = info.Operator
	db.Tags = info.Tags
	if db.Tags == nil {
		db.Tags = make([]string, 0, 1)
	}
	db.Assets = info.Assets
	if db.Assets == nil {
		db.Assets = make([]string, 0, 1)
	}
	db.Targets = info.Targets
	if db.Targets == nil {
		db.Targets = make([]string, 0, 1)
	}

	err := nosql.CreateArticle(db)
	if err == nil {
		info.UID = db.UID.Hex()
		info.CreateTime = db.CreatedTime
		info.ID = db.ID
		info.UpdateTime = db.UpdatedTime
	}
	return err
}

func (mine *cacheContext) GetArticle(uid string) *ArticleInfo {
	db, err := nosql.GetArticle(uid)
	if err == nil {
		info := new(ArticleInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext)RemoveArticle(uid, operator string) error {
	err := nosql.RemoveArticle(uid, operator)
	return err
}

func (mine *cacheContext) GetArticlesByOwner(uid string) []*ArticleInfo {
	if uid == ""{
		return make([]*ArticleInfo, 0, 1)
	}
	array, err := nosql.GetArticlesByOwner(uid)
	if err == nil {
		list := make([]*ArticleInfo, 0, len(array))
		for _, item := range array {
			info := new(ArticleInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return make([]*ArticleInfo, 0, 1)
}

func (mine *cacheContext) GetArticlesByTP(owner string, st uint8) []*ArticleInfo {
	var array []*nosql.Article
	var err error
	array, err = nosql.GetArticlesByOwnerTP(owner, st)
	if err == nil {
		list := make([]*ArticleInfo, 0, len(array))
		for _, item := range array {
			info := new(ArticleInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return make([]*ArticleInfo, 0, 1)
}

func (mine *cacheContext) GetArticlesByStatus(owner string, st MessageStatus) []*ArticleInfo {
	var array []*nosql.Article
	var err error
	array, err = nosql.GetArticlesByOwnerStatus(owner, uint8(st))
	if err == nil {
		list := make([]*ArticleInfo, 0, len(array))
		for _, item := range array {
			info := new(ArticleInfo)
			info.initInfo(item)
			list = append(list, info)
		}
		return list
	}
	return make([]*ArticleInfo, 0, 1)
}

func (mine *cacheContext) GetArticlesByList(array []string) []*ArticleInfo {
	if array == nil || len(array) < 1 {
		return make([]*ArticleInfo, 0, 1)
	}
	list := make([]*ArticleInfo, 0, 1)
	for _, s := range array {
		db, _ := nosql.GetArticle(s)
		if db != nil {
			info := new(ArticleInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}

	return list
}

func (mine *cacheContext) GetArticlesByTargets(owner string, array []string, page, num uint32) (uint32, uint32, []*ArticleInfo) {
	if array == nil || len(array) < 1 {
		return 0, 0, make([]*ArticleInfo, 0, 1)
	}
	all := make([]*ArticleInfo, 0, 10)
	var dbs []*nosql.Article
	var er error
	if len(owner) < 1{
		dbs,er = nosql.GetArticlesByTargets(array)
	}else{
		dbs,er = nosql.GetArticlesByOTargets(owner, array)
	}
	if er == nil {
		for _, db := range dbs {
			info := new(ArticleInfo)
			info.initInfo(db)
			all = append(all, info)
		}
	}
	if num < 1 {
		num = 10
	}
	if len(all) < 1 {
		return 0, 0, make([]*ArticleInfo, 0, 1)
	}
	max, pages, list := checkPage(page, num, all)
	return max, pages, list.([]*ArticleInfo)
}

func (mine *ArticleInfo) initInfo(db *nosql.Article) {
	mine.UID = db.UID.Hex()
	mine.Subtitle = db.Subtitle
	mine.ID = db.ID
	mine.Name = db.Name
	mine.CreateTime = db.CreatedTime
	mine.UpdateTime = db.UpdatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Body = db.Body
	mine.Type = db.Type
	mine.Owner = db.Owner
	mine.Tags = db.Tags
	mine.Targets = db.Targets
	mine.Status = MessageStatus(db.Status)
	mine.Assets = db.Assets
	if mine.Targets == nil {
		mine.Targets = make([]string, 0, 1)
		_ = mine.UpdateTargets(mine.Operator, mine.Targets)
	}
	if mine.Assets == nil {
		mine.Assets = make([]string, 0, 1)
	}
	if mine.Tags == nil {
		mine.Tags = make([]string, 0, 1)
	}
}

func (mine *ArticleInfo) UpdateBase(name, sub, body, operator string) error {
	if len(name) < 1 {
		name = mine.Name
	}
	if len(sub) < 1 {
		sub = mine.Subtitle
	}
	if len(body) < 1 {
		body = mine.Body
	}
	err := nosql.UpdateArticleBase(mine.UID, name, sub, body, operator)
	if err == nil {
		mine.Name = name
		mine.Subtitle = sub
		mine.Body = body
		mine.Operator = operator
	}
	return err
}

func (mine *ArticleInfo) UpdateTags(operator string, tags []string) error {
	if tags == nil {
		return errors.New("the list of tags is nil")
	}
	err := nosql.UpdateArticleTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *ArticleInfo) UpdateStatus(st MessageStatus, operator string) error {
	err := nosql.UpdateArticleStatus(mine.UID, operator, uint8(st))
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *ArticleInfo) UpdateTargets(operator string, list []string) error {
	if list == nil {
		return errors.New("the list of targets is nil")
	}
	err := nosql.UpdateArticleTargets(mine.UID, operator, list)
	if err == nil {
		mine.Targets = list
		mine.Operator = operator
	}
	return err
}

func (mine *ArticleInfo) UpdateAssets(operator string, list []string) error {
	if list == nil {
		return errors.New("the list of assets is nil")
	}
	err := nosql.UpdateArticleAssets(mine.UID, operator, list)
	if err == nil {
		mine.Assets = list
		mine.Operator = operator
	}
	return err
}

func (mine *ArticleInfo)HadTargets(arr []string) bool {
	if mine.Targets == nil || len(mine.Targets) < 1 {
		return true
	}
	if arr == nil || len(arr) < 1 {
		return false
	}
	for _, item := range arr {
		if tool.HasItem(mine.Targets, item) {
			return true
		}
	}
	return false
}

