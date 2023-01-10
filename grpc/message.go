package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
	"sort"
	"strconv"
	"time"
)

type MessageService struct{}

func switchActivityMessage(info *cache.ActivityInfo, entity string) *pb.MessageInfo {
	tmp := new(pb.MessageInfo)
	tmp.Uid = info.UID
	tmp.Type = cache.ObserveActivity
	tmp.Name = info.Name
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Remark = info.Remark
	tmp.Owner = info.Owner
	tmp.Date = &pb.DateInfo{Start: info.Date.Start, Stop: info.Date.Stop}
	tmp.Organizer = info.Organizer
	tmp.Targets = info.Targets
	tmp.Tags = info.Tags
	tmp.Entity = entity
	return tmp
}

func switchNoticeMessage(info *cache.NoticeInfo, entity string) *pb.MessageInfo {
	tmp := new(pb.MessageInfo)
	tmp.Uid = info.UID
	tmp.Type = cache.ObserveNotice
	tmp.Name = info.Name
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Remark = info.Body
	tmp.Owner = info.Owner
	tmp.Date = &pb.DateInfo{}
	tmp.Organizer = info.Owner
	tmp.Entity = entity
	tmp.Targets = info.Targets
	tmp.Tags = info.Tags
	return tmp
}

func (mine *MessageService) GetByFilter(ctx context.Context, in *pb.ReqMessageFilter, out *pb.ReplyMessages) error {
	path := "message.getByFilter"
	inLog(path, in)
	all := make([]*pb.MessageInfo, 0, 200)
	var max uint32 = 0
	var pages uint32 = 0
	if in.Key == "targets" {
		all = getObserves(in.List)
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].Created > all[j].Created {
			return true
		} else {
			return false
		}
	})
	max, pages, list := cache.CheckPage(in.Page, in.Number, all)
	out.List = list.([]*pb.MessageInfo)
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d, max = %d", len(out.List), max))
	return nil
}

func (mine *MessageService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "message.getStatistic"
	inLog(path, in)
	out.Owner = in.Owner
	if in.Key == "type" {
		tp, er := strconv.ParseUint(in.Value, 10, 32)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_FormatError)
			return nil
		}
		out.Count = cache.Context().GetRecordCount(in.Owner, uint8(tp))
	}

	out.Status = outLog(path, out)
	return nil
}

func getObserves(list []*pb.TargetInfo) []*pb.MessageInfo {
	all := make([]*pb.MessageInfo, 0, 200)
	for _, item := range list {
		list1 := cache.Context().GetAllActivitiesByStatus(item.Owner, cache.ActivityStatusPublish)
		for _, info := range list1 {
			if info.IsAlive() && info.HadTargets(item.Targets) {
				tmp := getObserveMessage(info.UID, all)
				if tmp == nil {
					all = append(all, switchActivityMessage(info, item.Entity))
				} else {
					tmp.Entity = tmp.Entity + ";" + item.Entity
				}
			}
		}
		list2 := cache.Context().GetNoticesByStatus(item.Owner, cache.NoticeToFamily, cache.MessageStatusAgree)
		var secs int64 = -3600 * 24 * 7
		var from int64 = int64(item.Time)
		if from < 1 {
			from = time.Now().Unix()
		}
		for _, info := range list2 {
			dif := info.CreateTime.Unix() - from
			if dif > secs && info.HadTargets(item.Targets) {
				tmp := getObserveMessage(info.UID, all)
				if tmp == nil {
					all = append(all, switchNoticeMessage(info, item.Entity))
				} else {
					tmp.Entity = tmp.Entity + ";" + item.Entity
				}
			}

		}
		//logger.Info(fmt.Sprintf("the entity of %s activity count = %d; notice count = %d", item.Entity, len(list1), len(list2)))
	}
	return all
}

func getObserveMessage(uid string, list []*pb.MessageInfo) *pb.MessageInfo {
	for _, info := range list {
		if info.Uid == uid {
			return info
		}
	}
	return nil
}
