package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	"omo.msa.favorite/cache"
	"sort"
)

const (
	ObserveActivity = 1
	ObserveNotice = 2
)

type MessageService struct {}

func switchActivityMessage(info *cache.ActivityInfo, entity string) *pb.MessageInfo {
	tmp := new(pb.MessageInfo)
	tmp.Uid = info.UID
	tmp.Type = ObserveActivity
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
	tmp.Type = ObserveNotice
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

func (mine *MessageService)GetByFilter(ctx context.Context, in *pb.ReqMessageFilter, out *pb.ReplyMessages) error {
	path := "message.getByFilter"
	inLog(path, in)
	all := make([]*pb.MessageInfo, 0, 200)
	var max uint32 = 0
	var pages uint32 = 0
	if in.Key == "targets" {
		for _, item := range in.List {
			list1 := cache.Context().GetAllActivitiesByTargets(item.Owner,cache.ActivityStatusPublish, item.Time, item.Targets)
			for _, info := range list1 {
				tmp := getObserveMessage(info.UID, all)
				if tmp == nil {
					all = append(all, switchActivityMessage(info, item.Entity))
				}else{
					tmp.Entity = tmp.Entity+";"+item.Entity
				}
			}
			list2 := cache.Context().GetAllNoticesByTargets(item.Owner, cache.MessageStatusAgree, item.Time, item.Targets)
			for _, info := range list2 {
				tmp := getObserveMessage(info.UID, all)
				if tmp == nil {
					all = append(all, switchNoticeMessage(info, item.Entity))
				}else{
					tmp.Entity = tmp.Entity+";"+item.Entity
				}
			}
			//logger.Info(fmt.Sprintf("the entity of %s activity count = %d; notice count = %d", item.Entity, len(list1), len(list2)))
		}
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].Created > all[j].Created {
			return true
		}else{
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

func getObserveMessage(uid string, list []*pb.MessageInfo) *pb.MessageInfo {
	for _, info := range list {
		if info.Uid == uid {
			return info
		}
	}
	return nil
}
