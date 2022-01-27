package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
	"strconv"
)

type ArticleService struct {}

func switchArticle(info *cache.ArticleInfo) *pb.ArticleInfo {
	tmp := new(pb.ArticleInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Body = info.Body
	tmp.Subtitle = info.Subtitle
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Owner = info.Owner
	tmp.Type = uint32(info.Type)
	tmp.Status = uint32(info.Status)
	tmp.Assets = info.Assets
	tmp.Targets = info.Targets
	return tmp
}

func (mine *ArticleService)AddOne(ctx context.Context, in *pb.ReqArticleAdd, out *pb.ReplyArticleInfo) error {
	path := "article.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path,"the owner is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	info := new(cache.ArticleInfo)
	info.Name = in.Name
	info.Subtitle = in.Subtitle
	info.Body = in.Body
	info.Creator = in.Operator
	info.Tags = in.Tags
	info.Targets = in.Targets
	info.Assets = in.Assets
	info.Owner = in.Owner
	info.Status = cache.MessageStatusDraft
	info.Type = cache.ArticleTypeDef
	err := cache.Context().CreateArticle(info)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchArticle(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ArticleService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyArticleInfo) error {
	path := "article.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the article uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetArticle(in.Uid)
	if info == nil {
		out.Status = outError(path,"the article not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchArticle(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ArticleService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "article.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the article uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveArticle(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Parent = in.Owner
	out.Status = outLog(path, out)
	return nil
}

func (mine *ArticleService)GetList(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyArticleList) error {
	path := "article.getList"
	inLog(path, in)
	var array []*cache.ArticleInfo
	var max uint32 = 0
	var pages uint32 = 0
	if in.Key == "type" {
		st, er := strconv.ParseUint(in.Value, 10, 32)
		if er == nil {
			array = cache.Context().GetArticlesByTP(in.Owner, uint8(st))
		}else{
			out.Status = outError(path,er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}else if in.Key == "status" {
		st, er := strconv.ParseUint(in.Value, 10, 32)
		if er == nil {
			array = cache.Context().GetArticlesByStatus(in.Owner, cache.MessageStatus(st))
		}else{
			out.Status = outError(path,er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}else if in.Key == "targets" {
		max, pages, array = cache.Context().GetArticlesByTargets(in.Owner, in.List, cache.MessageStatusAgree, in.Page, in.Number)
	}else if in.Key == "array" {
		array = cache.Context().GetArticlesByList(in.List)
	}else{
		array = cache.Context().GetArticlesByOwner(in.Owner)
	}
	out.List = make([]*pb.ArticleInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchArticle(val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ArticleService)GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "article.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *ArticleService)UpdateBase(ctx context.Context, in *pb.ReqArticleUpdate, out *pb.ReplyArticleInfo) error {
	path := "article.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the article uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetArticle(in.Uid)
	if info == nil {
		out.Status = outError(path,"the article not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Name != info.Name || in.Subtitle != info.Subtitle || in.Body != info.Body {
		err = info.UpdateBase(in.Name, in.Subtitle, in.Body, in.Operator)
		if err != nil {
			out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}
	if len(in.Tags) > 1 {
		err = info.UpdateTags(in.Operator, in.Targets)
		if err != nil {
			out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}

	if len(in.Targets) > 1{
		err = info.UpdateTargets(in.Operator, in.Targets)
		if err != nil {
			out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}

	out.Info = switchArticle(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ArticleService)UpdateStatus(ctx context.Context, in *pb.RequestState, out *pb.ReplyInfo) error {
	path := "article.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the article uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetArticle(in.Uid)
	if info == nil {
		out.Status = outError(path,"the article not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateStatus(cache.MessageStatus(in.Status), in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ArticleService)UpdateAssets(ctx context.Context, in *pb.RequestList, out *pb.ReplyInfo) error {
	path := "article.updateAssets"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the article uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetArticle(in.Uid)
	if info == nil {
		out.Status = outError(path,"the article not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateAssets(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ArticleService)UpdateTags(ctx context.Context, in *pb.RequestList, out *pb.ReplyInfo) error {
	path := "article.updateTags"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the article uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetArticle(in.Uid)
	if info == nil {
		out.Status = outError(path,"the article not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateTags(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ArticleService)UpdateTargets(ctx context.Context, in *pb.RequestList, out *pb.ReplyInfo) error {
	path := "article.updateTargets"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the article uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetArticle(in.Uid)
	if info == nil {
		out.Status = outError(path,"the article not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateTargets(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
