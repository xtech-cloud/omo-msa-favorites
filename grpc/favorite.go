package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	"omo.msa.favorite/cache"
)

type FavoriteService struct {}

func switchFavorite(owner string, info *cache.FavoriteInfo) *pb.FavoriteInfo {
	tmp := new(pb.FavoriteInfo)
	tmp.Owner = owner
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Owner = info.Owner
	tmp.Tags = info.Tags
	tmp.Origin = info.Origin
	tmp.Entities = info.GetEntities()
	return tmp
}

func (mine *FavoriteService)AddOne(ctx context.Context, in *pb.ReqFavoriteAdd, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path,"the scene is empty", pb.ResultStatus_Empty)
		return nil
	}
	if len(in.Origin) > 0 {
		tmp := cache.Context().GetFavoriteByOrigin(in.Owner, in.Origin)
		if tmp != nil {
			out.Info = switchFavorite(in.Owner, tmp)
			out.Status = outLog(path, out)
			return nil
		}
	}
	info := new(cache.FavoriteInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Cover = in.Cover
	info.Creator = in.Operator
	info.Tags = in.Tags
	info.Entities = in.Entities
	info.Origin = in.Origin
	info.Owner = in.Owner
	info.Type = uint8(in.Type)
	err := cache.Context().CreateFavorite(info)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchFavorite(in.Owner, info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchFavorite(in.Owner, info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)GetByOrigin(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.getByOrigin"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavoriteByOrigin(in.Owner, in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchFavorite(info.Owner, info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "favorite.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveFavorite(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Parent = in.Owner
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)GetList(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteList) error {
	path := "favorite.getList"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.List = make([]*pb.FavoriteInfo, 0, 1)
	}else{
		array := cache.Context().GetFavoritesByOwner(in.Owner)
		out.List = make([]*pb.FavoriteInfo, 0, len(array))
		for _, val := range array {
			out.List = append(out.List, switchFavorite(in.Owner, val))
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *FavoriteService)UpdateBase(ctx context.Context, in *pb.ReqFavoriteUpdate, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pb.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if len(in.Cover) > 0 {
		err = info.UpdateCover(in.Cover, in.Operator)
	}
	if len(in.Name) > 0 || len(in.Remark) > 0 {
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}

	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchFavorite(info.Owner, info)
	out.Status = outLog(path, out)
	return nil
}


func (mine *FavoriteService)UpdateTags(ctx context.Context, in *pb.ReqFavoriteTags, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.updateTags"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	if in.Tags == nil || len(in.Tags) < 1 {
		out.Status = outError(path,"the favorite tags is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateTags(in.Operator, in.Tags)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchFavorite(info.Owner, info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)UpdateEntities(ctx context.Context, in *pb.ReqFavoriteEntities, out *pb.ReplyFavoriteEntities) error {
	path := "favorite.updateEntities"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pb.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateEntities(in.Operator, in.Entities)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Entities = info.GetEntities()
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)AppendEntity(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteEntities) error {
	path := "favorite.appendEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.AppendEntity(in.Entity)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Entities = info.GetEntities()
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)SubtractEntity(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteEntities) error {
	path := "favorite.subtractEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.SubtractEntity(in.Entity)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Entities = info.GetEntities()
	out.Status = outLog(path, out)
	return nil
}

