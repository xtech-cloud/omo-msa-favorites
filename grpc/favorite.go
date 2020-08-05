package grpc

import (
	"context"
	"errors"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	"omo.msa.favorite/cache"
	"omo.msa.favorite/proxy"
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
	array := info.GetEntities()
	tmp.Entities = make([]*pb.EntityInfo, 0, len(array))
	for _, value := range array {
		tmp.Entities = append(tmp.Entities, &pb.EntityInfo{Uid: value.UID, Name: value.Name})
	}
	return tmp
}

func (mine *FavoriteService)AddOne(ctx context.Context, in *pb.ReqFavoriteAdd, out *pb.ReplyFavoriteOne) error {
	inLog("favorite.add", in)
	if len(in.Owner) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the scene is empty")
	}
	owner := cache.GetOwner(in.Owner)
	info := new(cache.FavoriteInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Cover = in.Cover
	info.Creator = in.Operator
	err := owner.CreateFavorite(info)
	if err != nil {
		out.ErrorCode = pb.ResultStatus_DBException
	}else{
		out.Info = switchFavorite(in.Owner, info)
	}
	return err
}

func (mine *FavoriteService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteOne) error {
	inLog("favorite.get", in)
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the favorite is empty")
	}
	_,info := cache.GetFavorite(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the favorite not found")
	}
	out.Info = switchFavorite(in.Owner, info)
	return nil
}

func (mine *FavoriteService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	inLog("favorite.remove", in)
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the favorite is empty")
	}
	owner,info := cache.GetFavorite(in.Uid)
	if info == nil {
		return errors.New("the favorite not found")
	}
	err := owner.RemoveFavorite(in.Uid, in.Operator)
	out.Uid = in.Uid
	out.Owner = in.Owner
	if err != nil {
		out.ErrorCode = pb.ResultStatus_DBException
	}
	return err
}

func (mine *FavoriteService)GetList(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteList) error {
	inLog("favorite.list", in)
	if len(in.Owner) < 1 {
		return errors.New("the scene is empty")
	}
	scene := cache.GetOwner(in.Owner)
	array := scene.Favorites()
	out.List = make([]*pb.FavoriteInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchFavorite(in.Owner, val))
	}
	return nil
}

func (mine *FavoriteService)UpdateBase(ctx context.Context, in *pb.ReqFavoriteUpdate, out *pb.ReplyFavoriteOne) error {
	inLog("favorite.update", in)
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the favorite is empty")
	}
	_,info := cache.GetFavorite(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the favorite not found")
	}
	var err error
	if len(in.Cover) > 0 {
		err = info.UpdateCover(in.Cover, in.Operator)
	}
	if len(in.Name) > 0 || len(in.Remark) > 0 {
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}

	if err != nil {
		out.ErrorCode = pb.ResultStatus_DBException
	}
	out.Info = switchFavorite(info.Owner, info)
	return err
}

func (mine *FavoriteService)UpdateEntities(ctx context.Context, in *pb.ReqFavoriteEntities, out *pb.ReplyFavoriteEntities) error {
	inLog("favorite.entities", in)
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the favorite is empty")
	}
	_,info := cache.GetFavorite(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the favorite not found")
	}
	var err error
	list := make([]proxy.EntityInfo, 0, len(in.Entities))
	for _, val := range in.Entities {
		list = append(list, proxy.EntityInfo{UID: val.Uid, Name: val.Name})
	}

	err = info.UpdateEntities(in.Operator, list)
	if err != nil {
		out.ErrorCode = pb.ResultStatus_DBException
	}else{
		out.Entities = make([]*pb.EntityInfo, 0, len(info.GetEntities()))
		for _, value := range info.GetEntities() {
			out.Entities = append(out.Entities, &pb.EntityInfo{Uid: value.UID, Name: value.Name})
		}
	}
	return err
}

func (mine *FavoriteService)SubtractEntity(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteEntities) error {
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the favorite is empty")
	}
	_,info := cache.GetFavorite(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the favorite not found")
	}
	//err := info.SubtractEntity(in.Entity)
	//if err != nil {
	//	out.ErrorCode = pb.ResultStatus_DBException
	//}else{
	//	out.Entities = make([]*pb.EntityInfo, 0, len(info.GetEntities()))
	//	for _, value := range info.GetEntities() {
	//		out.Entities = append(out.Entities, &pb.EntityInfo{Uid: value.UID, Name: value.Name})
	//	}
	//}

	return nil
}

