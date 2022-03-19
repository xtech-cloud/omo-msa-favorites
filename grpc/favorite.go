package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
	"strconv"
)

type FavoriteService struct {}

func switchFavorite(info *cache.FavoriteInfo) *pb.FavoriteInfo {
	tmp := new(pb.FavoriteInfo)
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
	tmp.Status = uint32(info.Status)
	tmp.Keys = info.GetKeys()
	return tmp
}

func (mine *FavoriteService)AddOne(ctx context.Context, in *pb.ReqFavoriteAdd, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path,"the owner is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	if cache.Context().HadFavoriteByName(in.Owner, in.Name, uint8(in.Type)) {
		out.Status = outError(path,"the name is repeated", pbstatus.ResultStatus_Repeated)
		return nil
	}
	info := new(cache.FavoriteInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Cover = in.Cover
	info.Creator = in.Operator
	info.Tags = in.Tags
	info.Keys = in.Keys
	info.Owner = in.Owner
	info.Status = uint8(in.Status)
	info.Type = uint8(in.Type)
	err := cache.Context().CreateFavorite(info)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchFavorite(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchFavorite(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "favorite.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "favorite.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveFavorite(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
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
			array = cache.Context().GetFavoritesByType(in.Owner, uint8(in.Type))
		}else{
			array = cache.Context().GetFavoritesByOwner(in.Owner)
		}

		out.List = make([]*pb.FavoriteInfo, 0, len(array))
		for _, val := range array {
			out.List = append(out.List, switchFavorite(val))
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *FavoriteService)GetByList(ctx context.Context, in *pb.RequestList, out *pb.ReplyFavoriteList) error {
	path := "favorite.getByList"
	inLog(path, in)
	array := cache.Context().GetFavoritesByList(in.List)
	out.List = make([]*pb.FavoriteInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchFavorite(val))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *FavoriteService)GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyFavoriteList) error {
	path := "favorite.getByFilter"
	inLog(path, in)
	var array []*cache.FavoriteInfo
	var max uint32 = 0
	var pages uint32 = 0
	if in.Key == "status" {
		if in.List != nil && len(in.List) > 1 {
			for _, val := range in.List {
				st, er := strconv.ParseUint(val, 10, 32)
				if er != nil {
					out.Status = outError(path,er.Error(), pbstatus.ResultStatus_FormatError)
					return nil
				}
				arr := cache.Context().GetFavoritesByStatus(in.Owner, uint8(st))
				if len(arr) > 0 {
					array = append(array, arr...)
				}
			}
		}else{
			st, er := strconv.ParseUint(in.Value, 10, 32)
			if er != nil {
				out.Status = outError(path,er.Error(), pbstatus.ResultStatus_FormatError)
				return nil
			}
			array = cache.Context().GetFavoritesByStatus(in.Owner, uint8(st))
		}
	}
	out.List = make([]*pb.FavoriteInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchFavorite(val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *FavoriteService)UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "favorite.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)UpdateBase(ctx context.Context, in *pb.ReqFavoriteUpdate, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
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
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchFavorite(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)UpdateMeta(ctx context.Context, in *pb.ReqFavoriteMeta, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.updateMeta"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateMeta(in.Operator, in.Meta)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchFavorite(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)UpdateStatus(ctx context.Context, in *pb.RequestState, out *pb.ReplyInfo) error {
	path := "favorite.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateStatus(uint8(in.Status), in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)UpdateTags(ctx context.Context, in *pb.ReqFavoriteTags, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.updateTags"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	if in.Tags == nil || len(in.Tags) < 1 {
		out.Status = outError(path,"the favorite tags is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateTags(in.Operator, in.Tags)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchFavorite(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)UpdateKeys(ctx context.Context, in *pb.ReqFavoriteKeys, out *pb.ReplyFavoriteKeys) error {
	path := "favorite.updateKeys"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateEntities(in.Operator, in.Keys)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Keys = info.GetKeys()
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)AppendKey(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteKeys) error {
	path := "favorite.appendEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.AppendKey(in.Flag)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Keys = info.GetKeys()
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)SubtractKey(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteKeys) error {
	path := "favorite.subtractEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.SubtractKey(in.Flag)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Keys = info.GetKeys()
	out.Status = outLog(path, out)
	return nil
}

