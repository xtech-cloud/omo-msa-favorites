package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
)

type SheetService struct {}

func switchSheet(info *cache.SheetInfo) *pb.SheetInfo {
	tmp := new(pb.SheetInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Owner = info.Owner
	tmp.Quote = info.Quote
	tmp.Status = uint32(info.Status)
	tmp.Keys = info.Keys
	return tmp
}

func (mine *SheetService)AddOne(ctx context.Context, in *pb.ReqSheetAdd, out *pb.ReplySheetInfo) error {
	path := "sheet.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path,"the owner is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	if cache.Context().HadSheetByName(in.Owner, in.Name) {
		out.Status = outError(path,"the name is repeated", pbstatus.ResultStatus_Repeated)
		return nil
	}
	info := new(cache.SheetInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Creator = in.Operator
	info.Keys = in.Keys
	info.Owner = in.Owner
	info.Status = uint8(in.Status)
	info.Quote = in.Quote
	err := cache.Context().CreateSheet(info)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchSheet(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SheetService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySheetInfo) error {
	path := "sheet.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the sheet uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetSheet(in.Uid)
	if info == nil {
		out.Status = outError(path,"the sheet not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchSheet(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SheetService)GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "sheet.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *SheetService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "sheet.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the sheet uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveSheet(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Parent = in.Owner
	out.Status = outLog(path, out)
	return nil
}

func (mine *SheetService)GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplySheetList) error {
	path := "sheet.getByFilter"
	inLog(path, in)
	var array []*cache.SheetInfo
	var max uint32 = 0
	var pages uint32 = 0
	if in.Key == "" {
		array = cache.Context().GetSheetsByOwner(in.Owner)
	}else{

	}
	out.List = make([]*pb.SheetInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchSheet(val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *SheetService)UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "sheet.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the sheet uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *SheetService)UpdateBase(ctx context.Context, in *pb.ReqSheetUpdate, out *pb.ReplySheetInfo) error {
	path := "sheet.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the sheet uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetSheet(in.Uid)
	if info == nil {
		out.Status = outError(path,"the sheet not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if len(in.Name) > 0 || len(in.Remark) > 0 {
		err = info.UpdateBase(in.Name, in.Remark, in.Operator)
	}

	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchSheet(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SheetService)UpdateStatus(ctx context.Context, in *pb.RequestState, out *pb.ReplyInfo) error {
	path := "sheet.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the sheet uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetSheet(in.Uid)
	if info == nil {
		out.Status = outError(path,"the sheet not found", pbstatus.ResultStatus_NotExisted)
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


func (mine *SheetService)UpdateKeys(ctx context.Context, in *pb.ReqSheetKeys, out *pb.ReplySheetKeys) error {
	path := "sheet.updateKeys"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the sheet uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetSheet(in.Uid)
	if info == nil {
		out.Status = outError(path,"the sheet not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateKeys(in.Operator, in.Keys)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Keys = info.Keys
	out.Status = outLog(path, out)
	return nil
}

func (mine *SheetService)AppendKey(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySheetKeys) error {
	path := "sheet.appendEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the sheet uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetSheet(in.Uid)
	if info == nil {
		out.Status = outError(path,"the sheet not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.AppendKey(in.Flag)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Keys = info.Keys
	out.Status = outLog(path, out)
	return nil
}

func (mine *SheetService)SubtractKey(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySheetKeys) error {
	path := "sheet.subtractEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path,"the sheet uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetSheet(in.Uid)
	if info == nil {
		out.Status = outError(path,"the sheet not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.SubtractKey(in.Flag)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Keys = info.Keys
	out.Status = outLog(path, out)
	return nil
}
