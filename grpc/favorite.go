package grpc

import (
	"context"
	"errors"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	"omo.msa.favorite/cache"
)

type FavoriteService struct {}

func switchFavorite(scene string, info *cache.FavoriteInfo) *pb.FavoriteInfo {
	tmp := new(pb.FavoriteInfo)
	tmp.Scene = scene
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	array := info.GetEntities()
	tmp.Entities = make([]*pb.EntityInfo, 0, len(array))
	for _, value := range array {
		tmp.Entities = append(tmp.Entities, &pb.EntityInfo{Uid: value.UID, Name: value.Name})
	}
	return tmp
}

func (mine *FavoriteService)AddOne(ctx context.Context, in *pb.ReqFavoriteAdd, out *pb.ReplyFavoriteOne) error {
	if len(in.Scene) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the scene is empty")
	}
	scene := cache.GetScene(in.Scene)
	info := new(cache.FavoriteInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	err := scene.CreateFavorite(info)
	if err != nil {
		out.ErrorCode = pb.ResultStatus_DBException
	}else{
		out.Info = switchFavorite(in.Scene, info)
	}
	return err
}

func (mine *FavoriteService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteOne) error {
	if len(in.Scene) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the scene is empty")
	}
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the favorite is empty")
	}
	scene := cache.GetScene(in.Scene)
	info := scene.GetFavorite(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the favorite not found")
	}
	out.Info = switchFavorite(in.Scene, info)
	return nil
}

func (mine *FavoriteService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	if len(in.Scene) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the scene is empty")
	}
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the favorite is empty")
	}
	scene := cache.GetScene(in.Scene)
	err := scene.RemoveFavorite(in.Uid)
	out.Uid = in.Uid
	if err != nil {
		out.ErrorCode = pb.ResultStatus_DBException
	}
	return err
}

func (mine *FavoriteService)GetListByScene(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyFavoriteList) error {
	if len(in.Scene) < 1 {
		return errors.New("the scene is empty")
	}
	scene := cache.GetScene(in.Scene)
	array := scene.Favorites()
	out.List = make([]*pb.FavoriteInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchFavorite(in.Scene, val))
	}
	return nil
}

func (mine *FavoriteService)UpdateBase(ctx context.Context, in *pb.ReqFavoriteUpdate, out *pb.ReplyInfo) error {
	if len(in.Scene) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the scene is empty")
	}
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the favorite is empty")
	}
	scene := cache.GetScene(in.Scene)
	info := scene.GetFavorite(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the favorite not found")
	}
	var err error
	if len(in.Cover) > 0 {
		err = info.UpdateCover(in.Cover)
	}else{
		err = info.UpdateBase(in.Name, in.Remark)
	}
	if err != nil {
		out.ErrorCode = pb.ResultStatus_DBException
	}
	return err
}

func (mine *FavoriteService)AppendEntity(ctx context.Context, in *pb.ReqFavoriteEntity, out *pb.ReplyFavoriteEntities) error {
	if len(in.Scene) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the scene is empty")
	}
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the favorite is empty")
	}
	scene := cache.GetScene(in.Scene)
	info := scene.GetFavorite(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the favorite not found")
	}
	var err error
	for _, val := range in.Entities {
		err = info.AppendEntity(val.Uid, val.Name)
	}
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
	if len(in.Scene) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the scene is empty")
	}
	if len(in.Uid) < 1 {
		out.ErrorCode = pb.ResultStatus_Empty
		return errors.New("the favorite is empty")
	}
	scene := cache.GetScene(in.Scene)
	info := scene.GetFavorite(in.Uid)
	if info == nil {
		out.ErrorCode = pb.ResultStatus_NotExisted
		return errors.New("the favorite not found")
	}
	err := info.SubtractEntity(in.Entity)
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

