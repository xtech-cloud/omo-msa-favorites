package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	"omo.msa.favorite/cache"
)

type RepertoryService struct {}

func switchRepertory(info *cache.OwnerInfo) *pb.RepertoryInfo {
	tmp := new(pb.RepertoryInfo)
	tmp.Owner = info.Owner
	tmp.Uid = info.UID
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Bags = info.Bags()
	return tmp
}

func (mine *RepertoryService)AppendOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRepertoryInfo) error {
	path := "repertory.appendOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path,"the owner is empty", pb.ResultStatus_Empty)
		return nil
	}
	owner := cache.GetOwner(in.Owner)
	if owner == nil {
		out.Status = outError(path,"the repertory not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := owner.AppendBag(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRepertory(owner)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RepertoryService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRepertoryInfo) error {
	path := "repertory.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the repertory uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.GetOwner(in.Uid)
	if info == nil {
		out.Status = outError(path,"the repertory not found", pb.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchRepertory(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RepertoryService)SubtractOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRepertoryInfo) error {
	path := "repertory.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the repertory uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.GetOwner(in.Owner)
	if info == nil {
		out.Status = outError(path,"the repertory not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.SubtractBag(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRepertory(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RepertoryService)UpdateList(ctx context.Context, in *pb.ReqRepertoryBags, out *pb.ReplyRepertoryInfo) error {
	path := "repertory.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the repertory uid is empty", pb.ResultStatus_Empty)
		return nil
	}
	info := cache.GetOwner(in.Uid)
	if info == nil {
		out.Status = outError(path,"the repertory not found", pb.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateBags(in.List, "")
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRepertory(info)
	out.Status = outLog(path, out)
	return nil
}
