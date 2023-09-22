package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
	"omo.msa.favorite/proxy"
	"strconv"
	"strings"
)

type NoticeService struct{}

func switchNotice(info *cache.NoticeInfo) *pb.NoticeInfo {
	tmp := new(pb.NoticeInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Name = info.Name
	tmp.Type = uint32(info.Type)
	tmp.Body = info.Body
	tmp.Subtitle = info.Subtitle
	tmp.Owner = info.Owner
	tmp.Status = uint32(info.Status)
	tmp.Targets = info.Targets
	tmp.Interval = info.Interval
	tmp.Showtime = info.Showtime
	tmp.Date = &pb.DateInfo{Start: info.Duration.Begin(), Stop: info.Duration.End()}
	return tmp
}

func (mine *NoticeService) AddOne(ctx context.Context, in *pb.ReqNoticeAdd, out *pb.ReplyNoticeInfo) error {
	path := "notice.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path, "the owner is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	in.Name = strings.TrimSpace(in.Name)
	info := new(cache.NoticeInfo)
	info.Name = in.Name
	info.Subtitle = in.Subtitle
	info.Body = in.Body
	info.Creator = in.Operator
	info.Tags = make([]string, 0, 1)
	info.Targets = in.Targets
	info.Owner = in.Owner
	info.Type = uint8(in.Type)
	info.Interval = in.Interval
	info.Showtime = in.Showtime

	begin, end := cache.SwitchDate(in.Date.Start, in.Date.Stop)
	info.Duration = proxy.DurationInfo{Start: begin, Stop: end}
	info.Status = cache.MessageStatusDraft
	err := cache.Context().CreateNotice(info)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchNotice(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *NoticeService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyNoticeInfo) error {
	path := "notice.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the notice uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetNotice(in.Uid)
	if info == nil {
		out.Status = outError(path, "the notice not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchNotice(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *NoticeService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "notice.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *NoticeService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "notice.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the notice uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveNotice(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Parent = in.Owner
	out.Status = outLog(path, out)
	return nil
}

func (mine *NoticeService) GetList(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyNoticeList) error {
	path := "notice.getList"
	inLog(path, in)
	var max uint32 = 0
	var pages uint32 = 0
	var array []*cache.NoticeInfo
	if in.Key == "status" {
		//st, er := strconv.ParseUint(in.Value, 10, 32)
		//if er == nil {
		//	array = cache.Context().getn("", cache.MessageStatus(st))
		//}else{
		//	out.Status = outError(path,er.Error(), pb.ResultStatus_DBException)
		//	return nil
		//}
	} else if in.Key == "targets" {
		max, pages, array = cache.Context().GetNoticesByTargets(in.Owner, in.List, cache.MessageStatusAgree, cache.NoticeToFamily, in.Page, in.Number)
	} else if in.Key == "array" {
		array = cache.Context().GetNoticesByList(in.List)
	} else if in.Key == "type" {
		tp, er := strconv.ParseUint(in.Value, 10, 32)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
		array = cache.Context().GetNoticesByType(in.Owner, uint32(tp))
	} else if in.Key == "latest" {
		tp, er := strconv.ParseUint(in.Value, 10, 32)
		if er == nil {
			array = cache.Context().GetLatestNotices(in.Owner, uint32(tp))
		} else {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_FormatError)
			return nil
		}
	} else if in.Key == "alive" {
		tp, er := strconv.ParseUint(in.Value, 10, 32)
		if er == nil {
			array = cache.Context().GetAliveNotices(in.Owner, uint32(tp))
		} else {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_FormatError)
			return nil
		}
	} else {
		array = cache.Context().GetNoticesByOwner(in.Owner)
	}
	out.List = make([]*pb.NoticeInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchNotice(val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *NoticeService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "notice.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the notice uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *NoticeService) UpdateBase(ctx context.Context, in *pb.ReqNoticeUpdate, out *pb.ReplyNoticeInfo) error {
	path := "notice.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the notice uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	in.Name = strings.TrimSpace(in.Name)
	info := cache.Context().GetNotice(in.Uid)
	if info == nil {
		out.Status = outError(path, "the notice not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateBase(in.Name, in.Subtitle, in.Body, in.Operator, in.Interval, in.Showtime, proxy.DateInfo{Start: in.Date.Start, Stop: in.Date.Stop})
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	if len(in.Targets) > 1 {
		err = info.UpdateTargets(in.Operator, in.Targets)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}

	out.Info = switchNotice(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *NoticeService) UpdateStatus(ctx context.Context, in *pb.RequestState, out *pb.ReplyInfo) error {
	path := "notice.updateStatus"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the notice uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetNotice(in.Uid)
	if info == nil {
		out.Status = outError(path, "the notice not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateStatus(cache.MessageStatus(in.Status), "")
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *NoticeService) UpdateTags(ctx context.Context, in *pb.RequestList, out *pb.ReplyInfo) error {
	path := "notice.updateTags"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the notice uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetNotice(in.Uid)
	if info == nil {
		out.Status = outError(path, "the notice not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateTags(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *NoticeService) UpdateTargets(ctx context.Context, in *pb.RequestList, out *pb.ReplyInfo) error {
	path := "notice.updateTargets"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the notice uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info := cache.Context().GetNotice(in.Uid)
	if info == nil {
		out.Status = outError(path, "the notice not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	err = info.UpdateTargets(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
