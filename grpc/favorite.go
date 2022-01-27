package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
	"omo.msa.favorite/proxy"
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
	tmp.Origin = info.Origin
	tmp.Status = uint32(info.Status)
	tmp.Keys = info.GetKeys()
	tmp.Targets = info.GetTargets()
	return tmp
}

func (mine *FavoriteService)AddOne(ctx context.Context, in *pb.ReqFavoriteAdd, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path,"the owner is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	if len(in.Origin) > 0 {
		tmp := cache.Context().GetFavoriteByOrigin(in.Owner, in.Origin, in.Person)
		if tmp != nil {
			out.Info = switchFavorite(tmp)
			out.Status = outLog(path, out)
			return nil
		}
	}
	if cache.Context().HadFavoriteByName(in.Owner, in.Name, uint8(in.Type), in.Person) {
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
	info.Origin = in.Origin
	info.Owner = in.Owner
	info.Status = uint8(in.Status)
	info.Type = uint8(in.Type)
	err := cache.Context().CreateFavorite(info, in.Person)
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
	info := cache.Context().GetFavorite(in.Uid, in.Person)
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

func (mine *FavoriteService)GetByOrigin(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.getByOrigin"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavoriteByOrigin(in.Owner, in.Uid, in.Person)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchFavorite(info)
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
	err := cache.Context().RemoveFavorite(in.Uid, in.Operator, in.Person)
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
			array = cache.Context().GetFavoritesByType(in.Owner, uint8(in.Type), in.Person)
		}else{
			array = cache.Context().GetFavoritesByOwner(in.Owner, in.Person)
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
	array := cache.Context().GetFavoritesByList(in.Person, in.List)
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
	if in.Key == "target" {

	}else if in.Key == "targets" {
		max, pages, array = cache.Context().GetFavoritesByTargets(in.Owner, in.List, in.Page, in.Number)
	}else if in.Key == "status" {
		st, er := strconv.ParseUint(in.Value, 10, 32)
		if er != nil {
			out.Status = outError(path,er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
		array = cache.Context().GetFavoritesByStatus(in.Owner, uint8(st), in.Person)
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

func (mine *FavoriteService)UpdateBase(ctx context.Context, in *pb.ReqFavoriteUpdate, out *pb.ReplyFavoriteInfo) error {
	path := "favorite.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid, in.Person)
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
	info := cache.Context().GetFavorite(in.Uid, in.Person)
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
	info := cache.Context().GetFavorite(in.Uid, in.Person)
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
	info := cache.Context().GetFavorite(in.Uid, in.Person)
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
	info := cache.Context().GetFavorite(in.Uid, in.Person)
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
	info := cache.Context().GetFavorite(in.Uid, in.Person)
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
	info := cache.Context().GetFavorite(in.Uid, in.Person)
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

func (mine *FavoriteService)UpdateTargets(ctx context.Context, in *pb.ReqFavoriteTargets, out *pb.ReplyInfo) error {
	path := "favorite.updateTargets"
	inLog(path, in)
	if in.List == nil || len(in.List) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	for _, uid := range in.List {
		info := cache.Context().GetFavorite(uid, false)
		if info == nil {
			out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		err := info.UpdateTargets(in.Operator, in.Targets)
		if err != nil {
			out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)UpdateTarget(ctx context.Context, in *pb.ReqFavoriteTarget, out *pb.ReplyFavoriteTargets) error {
	path := "favorite.updateTarget"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid, in.Person)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateTarget(in.Target, in.Effect, in.Skin, in.Operator, in.Slots)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Targets = info.GetTargets()
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)AppendTarget(ctx context.Context, in *pb.ReqFavoriteTarget, out *pb.ReplyFavoriteTargets) error {
	path := "favorite.appendTarget"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid, in.Person)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.AppendTarget(&proxy.ShowingInfo{Target: in.Target, Effect: in.Effect, Skin: in.Skin, Slots: in.Slots})
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Targets = info.GetTargets()
	out.Status = outLog(path, out)
	return nil
}

func (mine *FavoriteService)SubtractTarget(ctx context.Context, in *pb.ReqFavoriteTarget, out *pb.ReplyFavoriteTargets) error {
	path := "favorite.subtractTarget"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the favorite uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetFavorite(in.Uid, in.Person)
	if info == nil {
		out.Status = outError(path,"the favorite not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.SubtractTarget(in.Target)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Targets = info.GetTargets()
	out.Status = outLog(path, out)
	return nil
}

