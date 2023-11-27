package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
	"strconv"
	"time"
)

type MessageService struct{}

func switchActivityMessage(info *cache.ActivityInfo, entity string) *pb.MessageInfo {
	tmp := new(pb.MessageInfo)
	tmp.Uid = info.UID
	tmp.Type = uint32(info.Type)
	tmp.Name = info.Name
	tmp.Created = uint64(info.CreateTime.Unix())
	if info.UpdateTime.Unix() > 0 {
		tmp.Created = uint64(info.UpdateTime.Unix())
	}
	tmp.Remark = info.Remark
	tmp.Owner = info.Owner
	tmp.Date = &pb.DateInfo{Start: info.Duration.Begin(), Stop: info.Duration.End()}
	tmp.Organizer = info.Organizer
	tmp.Targets = info.Targets
	tmp.Tags = info.Tags
	tmp.Entity = entity
	return tmp
}

func switchNoticeMessage(info *cache.NoticeInfo, entity string) *pb.MessageInfo {
	tmp := new(pb.MessageInfo)
	tmp.Uid = info.UID
	tmp.Type = cache.MessageNotice
	tmp.Name = info.Name
	tmp.Created = uint64(info.CreateTime.Unix())
	if info.UpdateTime.Unix() > 0 {
		tmp.Created = uint64(info.UpdateTime.Unix())
	}
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
		tp := parseStringToInt(in.Value)
		if tp < 0 {
			all = getObserves(in.List)
		} else {
			all = getObservesByType(in.List, tp)
		}
	}

	max, pages, list := cache.CheckPage(in.Page, in.Number, all)
	out.List = list
	out.Total = max
	out.Pages = pages
	//out.Status = outLog(path, fmt.Sprintf("the length = %d, max = %d", len(out.List), max))
	out.Status = outLog(path, out)
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
					if item.Entity != "" {
						tmp.Entity = tmp.Entity + ";" + item.Entity
					}
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
					if item.Entity != "" {
						tmp.Entity = tmp.Entity + ";" + item.Entity
					}
				}
			}
		}
		//logger.Info(fmt.Sprintf("the entity of %s activity count = %d; notice count = %d", item.Entity, len(list1), len(list2)))
	}
	return all
}

func getObservesByType(list []*pb.TargetInfo, tp int) []*pb.MessageInfo {
	all := make([]*pb.MessageInfo, 0, 200)
	if tp == cache.MessageNotice {
		for _, item := range list {
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
						if item.Entity != "" {
							tmp.Entity = tmp.Entity + ";" + item.Entity
						}
					}
				}
			}
		}
	} else {
		for _, item := range list {
			list1 := cache.Context().GetAllActivitiesByStatusTP(item.Owner, cache.ActivityStatusPublish, uint8(tp))
			for _, info := range list1 {
				if info.IsAlive() && info.HadTargets(item.Targets) {
					tmp := getObserveMessage(info.UID, all)
					if tmp == nil {
						all = append(all, switchActivityMessage(info, item.Entity))
					} else {
						if item.Entity != "" {
							tmp.Entity = tmp.Entity + ";" + item.Entity
						}
					}
				}
			}
		}
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
