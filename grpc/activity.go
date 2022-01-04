package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	"omo.msa.favorite/cache"
	"omo.msa.favorite/proxy"
)

type ActivityService struct {}

func switchActivity(owner string, info *cache.ActivityInfo) *pb.ActivityInfo {
	tmp := new(pb.ActivityInfo)
	tmp.Owner = owner
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Cover = info.Cover
	tmp.Require = info.Require
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Tags = info.Tags
	tmp.Owner = info.Owner
	tmp.Date = &pb.DateInfo{Start: info.Date.Start, Stop: info.Date.Stop}
	tmp.Place = &pb.PlaceInfo{Name: info.Place.Name, Location: info.Place.Location}
	tmp.Organizer = info.Organizer
	tmp.Assets = info.Assets
	tmp.Limit = uint32(info.AssetLimit)
	tmp.Participants = make([]*pb.PairInfo, 0, len(info.Persons))
	for _, person := range info.Persons {
		tmp.Participants = append(tmp.Participants, &pb.PairInfo{Key: person.Entity, Value: person.Event})
	}
	tmp.Targets = info.Targets
	return tmp
}

func (mine *ActivityService)AddOne(ctx context.Context, in *pb.ReqActivityAdd, out *pb.ReplyActivityInfo) error {
	path := "activity.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path,"the owner is empty", pb.ResultStatus_Empty)
		return nil
	}

	info := new(cache.ActivityInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Cover = in.Cover
	info.Require = in.Require
	info.Creator = in.Operator
	info.Tags = in.Tags
	info.Assets = in.Assets
	info.Owner = in.Owner
	info.Organizer = in.Organizer
	info.Targets = in.Targets
	info.Date = proxy.DateInfo{Start: in.Date.Start, Stop: in.Date.Stop}
	info.Place = proxy.PlaceInfo{Name: in.Place.Name, Location: in.Place.Location}
	info.Type = uint8(in.Type)
	info.AssetLimit = uint8(in.Limit)
	info.Persons = make([]proxy.PersonInfo, 0 ,1)
	err := cache.Context().CreateActivity(info)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchActivity(in.Owner, info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyActivityInfo) error {
	path := "activity.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the activity uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetActivity(in.Uid)
	if info == nil {
		out.Status = outError(path,"the activity not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchActivity(in.Owner, info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "activity.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the activity uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveActivity(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Parent = in.Owner
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService)GetList(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyActivityList) error {
	path := "activity.getList"
	inLog(path, in)
	if len(in.Owner) > 1 {
		array := cache.Context().GetActivityByOrganizer(in.Owner)
		out.List = make([]*pb.ActivityInfo, 0, len(array))
		for _, val := range array {
			out.List = append(out.List, switchActivity(in.Owner, val))
		}
	} else{
		out.List = make([]*pb.ActivityInfo, 0, 1)
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ActivityService)GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyActivityList) error {
	path := "activity.getByFilter"
	inLog(path, in)
	var array []*cache.ActivityInfo
	var max uint32 = 0
	var pages uint32 = 0
	if in.Key == "target" {

	}else if in.Key == "targets" {
		max, pages, array = cache.Context().GetActivitiesByTargets(in.List, in.Page, in.Number)
	}
	out.List = make([]*pb.ActivityInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchActivity(val.Owner, val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ActivityService)UpdateBase(ctx context.Context, in *pb.ReqActivityUpdate, out *pb.ReplyActivityInfo) error {
	path := "activity.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the activity uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetActivity(in.Uid)
	if info == nil {
		out.Status = outError(path,"the activity not found", pb.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if len(in.Cover) > 0 {
		err = info.UpdateCover(in.Cover, in.Operator)
	}
	if len(in.Name) > 0 || len(in.Remark) > 0 {
		err = info.UpdateBase(in.Name, in.Remark, in.Require, in.Operator, proxy.DateInfo{Start: in.Date.Start, Stop: in.Date.Stop},
			proxy.PlaceInfo{Name: in.Place.Name, Location: in.Place.Location})
	}
	if uint8(in.Limit) != info.AssetLimit {
		err = info.UpdateAssetLimit(in.Operator, uint8(in.Limit))
	}

	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchActivity(info.Owner, info)
	out.Status = outLog(path, out)
	return nil
}


func (mine *ActivityService)UpdateTags(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "activity.updateTags"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the activity uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	if in.List == nil || len(in.List) < 1 {
		out.Status = outError(path,"the activity tags is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetActivity(in.Uid)
	if info == nil {
		out.Status = outError(path,"the activity not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateTags(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = info.Tags
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService)UpdateAssets(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "activity.updateAssets"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the activity uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetActivity(in.Uid)
	if info == nil {
		out.Status = outError(path,"the activity not found", pb.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateAssets(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = info.Assets
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService)UpdateTargets(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "activity.updateTargets"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the activity uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetActivity(in.Uid)
	if info == nil {
		out.Status = outError(path,"the activity not found", pb.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateTargets(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = info.Assets
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService)AppendOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPairList) error {
	path := "activity.appendEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the activity uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetActivity(in.Uid)
	if info == nil {
		out.Status = outError(path,"the activity not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.AppendPerson(in.Owner, in.Flag)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = info.GetEntities()
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService)SubtractOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPairList) error {
	path := "activity.subtractEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the activity uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetActivity(in.Uid)
	if info == nil {
		out.Status = outError(path,"the activity not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.SubtractPerson(in.Owner)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.List = info.GetEntities()
	out.Status = outLog(path, out)
	return nil
}

