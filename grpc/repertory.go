package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	"omo.msa.favorite/cache"
)

type RepertoryService struct {}

func switchRepertory(info *cache.RepertoryInfo) *pb.RepertoryInfo {
	tmp := new(pb.RepertoryInfo)

	return tmp
}

func (mine *RepertoryService)AppendOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRepertoryInfo) error {
	path := "repertory.appendOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path,"the owner is empty", pb.ResultStatus_Empty)
		return nil
	}
	owner,err := cache.Context().GetRepertory(in.Owner)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_NotExisted)
		return nil
	}
	err = owner.AppendAsset(in.Uid)
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
	info,err := cache.Context().GetRepertory(in.Owner)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_NotExisted)
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
	info,err := cache.Context().GetRepertory(in.Owner)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_NotExisted)
		return nil
	}
	err = info.SubtractAsset(in.Uid)
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
	info,err := cache.Context().GetRepertory(in.Uid)
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_NotExisted)
		return nil
	}
	err = info.UpdateBags(in.List, "")
	if err != nil {
		out.Status = outError(path,err.Error(), pb.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRepertory(info)
	out.Status = outLog(path, out)
	return nil
}
