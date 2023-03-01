package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
	"omo.msa.favorite/proxy"
	"strconv"
)

type DisplayService struct{}

func switchDisplay(info *cache.DisplayInfo) *pb.DisplayInfo {
	tmp := new(pb.DisplayInfo)
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
	tmp.Banner = info.Banner
	tmp.Poster = info.Poster
	tmp.Status = uint32(info.Status)
	tmp.Contents = info.GetContents()
	return tmp
}

func (mine *DisplayService) AddOne(ctx context.Context, in *pb.ReqDisplayAdd, out *pb.ReplyDisplayInfo) error {
	path := "display.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path, "the owner is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	if len(in.Origin) > 0 {
		tmp := cache.Context().GetDisplayByOrigin(in.Owner, in.Origin)
		if tmp != nil {
			out.Info = switchDisplay(tmp)
			out.Status = outLog(path, out)
			return nil
		}
	}
	if cache.Context().HadDisplayByName(in.Owner, in.Name, uint8(in.Type)) {
		out.Status = outError(path, "the name is repeated", pbstatus.ResultStatus_Repeated)
		return nil
	}
	info := new(cache.DisplayInfo)
	info.Name = in.Name
	info.Remark = in.Remark
	info.Cover = in.Cover
	info.Creator = in.Operator
	info.Tags = in.Tags
	info.Contents = make([]proxy.DisplayContent, 0, len(in.Keys))
	for _, key := range in.Keys {
		info.Contents = append(info.Contents, proxy.DisplayContent{UID: key, Events: make([]string, 0, 1), Assets: make([]string, 0, 1)})
	}
	info.Origin = in.Origin
	info.Owner = in.Owner
	info.Status = uint8(in.Status)
	info.Type = uint8(in.Type)
	err := cache.Context().CreateDisplay(info)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchDisplay(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyDisplayInfo) error {
	path := "display.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetDisplay(in.Uid)
	if info == nil {
		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchDisplay(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "display.getStatistic"
	inLog(path, in)

	if in.Key == "type" {
		tp := parseStringToInt(in.Value)
		array := cache.Context().GetDisplaysByType(in.Owner, uint8(tp))
		out.Count = uint32(len(array))
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) GetByOrigin(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyDisplayInfo) error {
	path := "display.getByOrigin"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetDisplayByOrigin(in.Owner, in.Uid)
	if info == nil {
		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchDisplay(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "display.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveDisplay(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Parent = in.Owner
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) GetList(ctx context.Context, in *pb.ReqDisplayList, out *pb.ReplyDisplayList) error {
	path := "display.getList"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.List = make([]*pb.DisplayInfo, 0, 1)
	} else {
		var array []*cache.DisplayInfo
		if in.Type > 0 {
			array = cache.Context().GetDisplaysByType(in.Owner, uint8(in.Type))
		} else {
			array = cache.Context().GetDisplaysByOwner(in.Owner)
		}

		out.List = make([]*pb.DisplayInfo, 0, len(array))
		for _, val := range array {
			out.List = append(out.List, switchDisplay(val))
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *DisplayService) GetByList(ctx context.Context, in *pb.RequestList, out *pb.ReplyDisplayList) error {
	path := "display.getByList"
	inLog(path, in)
	array := cache.Context().GetDisplaysByList(in.List)
	out.List = make([]*pb.DisplayInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchDisplay(val))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *DisplayService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyDisplayList) error {
	path := "display.getByFilter"
	inLog(path, in)
	var array []*cache.DisplayInfo
	var max uint32 = 0
	var pages uint32 = 0
	if in.Key == "target" {

	} else if in.Key == "targets" {
		max, pages, array = cache.Context().GetDisplaysByTargets(in.Owner, in.List, in.Page, in.Number)
	} else if in.Key == "status" {
		if in.List != nil && len(in.List) > 1 {
			for _, val := range in.List {
				st, er := strconv.ParseUint(val, 10, 32)
				if er != nil {
					out.Status = outError(path, er.Error(), pbstatus.ResultStatus_FormatError)
					return nil
				}
				arr := cache.Context().GetDisplaysByStatus(in.Owner, uint8(st))
				if len(arr) > 0 {
					array = append(array, arr...)
				}
			}
		} else {
			st, er := strconv.ParseUint(in.Value, 10, 32)
			if er != nil {
				out.Status = outError(path, er.Error(), pbstatus.ResultStatus_FormatError)
				return nil
			}
			array = cache.Context().GetDisplaysByStatus(in.Owner, uint8(st))
		}
	}
	out.List = make([]*pb.DisplayInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchDisplay(val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *DisplayService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "display.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetDisplay(in.Uid)
	if info == nil {
		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Key == "banner" {
		err = info.UpdateBanner(in.Value, in.Operator)
	} else if in.Key == "poster" {
		err = info.UpdatePoster(in.Value, in.Operator)
	} else {
		err = errors.New("the key not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) UpdateBase(ctx context.Context, in *pb.ReqDisplayUpdate, out *pb.ReplyDisplayInfo) error {
	path := "display.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetDisplay(in.Uid)
	if info == nil {
		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
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
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchDisplay(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) UpdateMeta(ctx context.Context, in *pb.ReqDisplayMeta, out *pb.ReplyDisplayInfo) error {
	path := "display.updateMeta"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetDisplay(in.Uid)
	if info == nil {
		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateMeta(in.Operator, in.Meta)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchDisplay(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) UpdateStatus(ctx context.Context, in *pb.RequestState, out *pb.ReplyInfo) error {
	path := "display.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetDisplay(in.Uid)
	if info == nil {
		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateStatus(uint8(in.Status), in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) UpdateTags(ctx context.Context, in *pb.ReqDisplayTags, out *pb.ReplyDisplayInfo) error {
	path := "display.updateTags"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	if in.Tags == nil || len(in.Tags) < 1 {
		out.Status = outError(path, "the display tags is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetDisplay(in.Uid)
	if info == nil {
		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateTags(in.Operator, in.Tags)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchDisplay(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) UpdateContents(ctx context.Context, in *pb.ReqDisplayContents, out *pb.ReplyDisplayContents) error {
	path := "display.updateContents"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetDisplay(in.Uid)
	if info == nil {
		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateContents(in.Operator, in.Contents)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Contents = info.GetContents()
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) AppendContent(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyDisplayContents) error {
	path := "display.appendContent"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetDisplay(in.Uid)
	if info == nil {
		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.AppendContent(in.Flag)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Contents = info.GetContents()
	out.Status = outLog(path, out)
	return nil
}

func (mine *DisplayService) SubtractContent(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyDisplayContents) error {
	path := "display.subtractContent"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetDisplay(in.Uid)
	if info == nil {
		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.SubtractContent(in.Flag)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Contents = info.GetContents()
	out.Status = outLog(path, out)
	return nil
}

//func (mine *DisplayService) UpdateTargets(ctx context.Context, in *pb.ReqDisplayTargets, out *pb.ReplyInfo) error {
//	path := "display.updateTargets"
//	inLog(path, in)
//	if in.Owner == "" {
//		out.Status = outError(path, "the owner is empty", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	if in.List == nil || len(in.List) < 1 {
//		out.Status = outError(path, "the display list is empty", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	empty := false
//	if in.Targets == nil || len(in.Targets) < 1 {
//		empty = true
//	}
//	arr := cache.Context().GetDisplaysByOwner(in.Owner)
//	for _, info := range arr {
//		if tool.HasItem(in.List, info.UID) {
//			if empty {
//				_ = info.UpdateTargets(in.Operator, nil)
//			} else {
//				for _, target := range in.Targets {
//					_ = info.AppendSimpleTarget(target)
//				}
//			}
//		} else {
//			if empty {
//
//			} else {
//				for _, target := range in.Targets {
//					_ = info.SubtractTarget(target)
//				}
//			}
//		}
//	}
//
//	out.Status = outLog(path, out)
//	return nil
//}

//func (mine *DisplayService) UpdateTarget(ctx context.Context, in *pb.ReqDisplayTarget, out *pb.ReplyDisplayTargets) error {
//	path := "display.updateTarget"
//	inLog(path, in)
//	if len(in.Uid) < 1 {
//		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info := cache.Context().GetDisplay(in.Uid)
//	if info == nil {
//		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//	err := info.UpdateTarget(in.Target, in.Effect, in.Skin, in.Menu, in.Operator, in.Slots)
//	if err != nil {
//		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Targets = info.GetTargets()
//	out.Status = outLog(path, out)
//	return nil
//}

//func (mine *DisplayService) AppendTarget(ctx context.Context, in *pb.ReqDisplayTarget, out *pb.ReplyDisplayTargets) error {
//	path := "display.appendTarget"
//	inLog(path, in)
//	if len(in.Uid) < 1 {
//		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info := cache.Context().GetDisplay(in.Uid)
//	if info == nil {
//		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//	err := info.AppendTarget(&proxy.ShowingInfo{Target: in.Target, Effect: in.Effect, Alignment: in.Skin, Menu: in.Menu, Slots: in.Slots, UpdatedAt: time.Now()})
//	if err != nil {
//		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Targets = info.GetTargets()
//	out.Status = outLog(path, out)
//	return nil
//}

//func (mine *DisplayService) SubtractTarget(ctx context.Context, in *pb.ReqDisplayTarget, out *pb.ReplyDisplayTargets) error {
//	path := "display.subtractTarget"
//	inLog(path, in)
//	if len(in.Uid) < 1 {
//		out.Status = outError(path, "the display uid is empty", pbstatus.ResultStatus_Empty)
//		return nil
//	}
//	info := cache.Context().GetDisplay(in.Uid)
//	if info == nil {
//		out.Status = outError(path, "the display not found", pbstatus.ResultStatus_NotExisted)
//		return nil
//	}
//	err := info.SubtractTarget(in.Target)
//	if err != nil {
//		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
//		return nil
//	}
//	out.Targets = info.GetTargets()
//	out.Status = outLog(path, out)
//	return nil
//}
