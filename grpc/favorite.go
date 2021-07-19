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
	tmp.Type = uint32(info.Type)
	tmp.Tags = info.Tags
	tmp.Origin = info.Origin
	tmp.Keys = info.GetKeys()
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
		tmp := cache.Context().GetFavoriteByOrigin(in.Owner, in.Origin, in.Person)
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
	info.Keys = in.Keys
	info.Origin = in.Origin
	info.Owner = in.Owner
	info.Type = uint8(in.Type)
	err := cache.Context().CreateFavorite(info, in.Person)
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
	info := cache.Context().GetFavorite(in.Uid, in.Person)
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
	info := cache.Context().GetFavoriteByOrigin(in.Owner, in.Uid, in.Person)
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
	err := cache.Context().RemoveFavorite(in.Uid, in.Operator, in.Person)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Parent = in.Owner
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)GetList(ctx context.Context, in *pb.ReqFavoriteList, out *pb.ReplyFavoriteList) error {
	path := "favorite.getList"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.List = make([]*pb.FavoriteInfo, 0, 1)
	}else{
		var array []*cache.FavoriteInfo
		if in.Type > 0 {
			array = cache.Context().GetFavoritesByType(in.Owner, uint8(in.Type), in.Person)
		}else{
			array = cache.Context().GetFavoritesByOwner(in.Owner, in.Person)
		}

		out.List = make([]*pb.FavoriteInfo, 0, len(array))
		for _, val := range array {
			out.List = append(out.List, switchFavorite(in.Owner, val))
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *FavoriteService)GetByList(ctx context.Context, in *pb.RequestList, out *pb.ReplyFavoriteList) error {
	path := "favorite.getByList"
	inLog(path, in)
	array := cache.Context().GetFavoritesByList(in.Person, in.List)
	out.List = make([]*pb.FavoriteInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchFavorite(val.Owner, val))
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
	info := cache.Context().GetFavorite(in.Uid, in.Person)
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
	info := cache.Context().GetFavorite(in.Uid, in.Person)
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

func (mine *FavoriteService)UpdateKeys(ctx context.Context, in *pb.ReqFavoriteKeys, out *pb.ReplyFavoriteKeys) error {
	path := "favorite.updateKeys"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid, in.Person)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pb.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateEntities(in.Operator, in.Keys)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Keys = info.GetKeys()
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)AppendEntity(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteKeys) error {
	path := "favorite.appendEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid, in.Person)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.AppendEntity(in.Entity)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Keys = info.GetKeys()
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)SubtractEntity(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteKeys) error {
	path := "favorite.subtractEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid, in.Person)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.SubtractEntity(in.Entity)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Keys = info.GetKeys()
	out.Status = outLog(path, out)
	return nil
}

