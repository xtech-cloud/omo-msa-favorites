package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
)

type RepertoryService struct{}

func switchRepertory(info *cache.RepertoryInfo) *pb.RepertoryInfo {
	tmp := new(pb.RepertoryInfo)

	return tmp
}

func (mine *RepertoryService) AppendOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRepertoryInfo) error {
	path := "repertory.appendOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path, "the owner is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	owner, err := cache.Context().GetRepertory(in.Owner)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = owner.AppendAsset(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRepertory(owner)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RepertoryService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRepertoryInfo) error {
	path := "repertory.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the repertory uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetRepertory(in.Owner)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchRepertory(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RepertoryService) SubtractOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyRepertoryInfo) error {
	path := "repertory.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the repertory uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetRepertory(in.Owner)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = info.SubtractAsset(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRepertory(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *RepertoryService) UpdateList(ctx context.Context, in *pb.ReqRepertoryBags, out *pb.ReplyRepertoryInfo) error {
	path := "repertory.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the repertory uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetRepertory(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = info.UpdateBags(in.List, "")
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchRepertory(info)
	out.Status = outLog(path, out)
	return nil
}
