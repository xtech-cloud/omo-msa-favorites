package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
	"omo.msa.favorite/proxy"
	"omo.msa.favorite/proxy/nosql"
	"omo.msa.favorite/tool"
	"strconv"
	"strings"
)

type ActivityService struct{}

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
	tmp.Poster = info.Poster
	tmp.Show = uint32(info.ShowResult)
	tmp.Type = uint32(info.Type)
	tmp.Status = uint32(info.Status)
	tmp.Date = &pb.DateInfo{Start: info.Duration.Begin(), Stop: info.Duration.End()}
	tmp.Place = &pb.PlaceInfo{Name: info.Place.Name, Location: info.Place.Location}
	tmp.Organizer = info.Organizer
	tmp.Template = info.Template
	tmp.Assets = info.Assets
	tmp.Participant = info.Participant
	tmp.Limit = uint32(info.SubmitLimit)
	tmp.Targets = info.Targets
	tmp.Access = uint32(info.Access)
	tmp.Quotes = info.Quotes
	tmp.Prize = switchPrize(info.Prize)
	tmp.Opuses = switchOpuses(info.Opuses)
	tmp.Records = switchHistories(info.GetHistories())
	return tmp
}

func switchPrize(info *proxy.PrizeInfo) *pb.PrizeInfo {
	if info == nil {
		t := new(pb.PrizeInfo)
		t.Ranks = make([]*pb.RankInfo, 0, 1)
		return t
	}
	tmp := new(pb.PrizeInfo)
	tmp.Name = info.Name
	tmp.Desc = info.Desc
	tmp.Ranks = make([]*pb.RankInfo, 0, len(info.Ranks))
	for _, rank := range info.Ranks {
		tmp.Ranks = append(tmp.Ranks, &pb.RankInfo{Index: rank.Index, Name: rank.Name, Count: rank.Count})
	}
	return tmp
}

func switchOpuses(list []proxy.OpusInfo) []*pb.OpusInfo {
	if list == nil {
		return make([]*pb.OpusInfo, 0, 1)
	}
	arr := make([]*pb.OpusInfo, 0, len(list))
	for _, info := range list {
		tmp := new(pb.OpusInfo)
		tmp.Rank = info.Rank
		tmp.Asset = info.Asset
		tmp.Remark = info.Remark
		arr = append(arr, tmp)
	}
	return arr
}

func switchHistories(dbs []*nosql.History) []*pb.RecordInfo {
	arr := make([]*pb.RecordInfo, 0, len(dbs))
	for _, item := range dbs {
		tmp := new(pb.RecordInfo)
		tmp.Operator = item.Creator
		tmp.From = item.From
		tmp.To = item.To
		tmp.Remark = item.Remark
		tmp.Stamp = item.CreatedTime.Unix()
		tmp.Content = item.Content
		tmp.Option = uint32(item.Option)
		arr = append(arr, tmp)
	}
	return arr
}

func (mine *ActivityService) AddOne(ctx context.Context, in *pb.ReqActivityAdd, out *pb.ReplyActivityInfo) error {
	path := "activity.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path, "the owner is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	in.Name = strings.TrimSpace(in.Name)
	//if len(in.Template) > 1 {
	//	list := cache.Context().GetActivitiesByTemplate(in.Owner, in.Template)
	//	if len(list) > 0 {
	//		out.Status = outError(path, "the activity had clone by this owner", pbstatus.ResultStatus_Repeated)
	//		return nil
	//	}
	//}
	//if in.Targets == nil || len(in.Targets) < 1 {
	//	out.Status = outError(path, "the activity targets is not empty", pbstatus.ResultStatus_Empty)
	//	return nil
	//}

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
	info.ShowResult = uint8(in.Show)
	info.Status = uint8(in.Status)
	info.Template = in.Template

	info.Duration = proxy.DurationInfo{Start: proxy.DateToUTC(in.Date.Start, 0), Stop: proxy.DateToUTC(in.Date.Stop, 1)}
	info.Place = proxy.PlaceInfo{Name: in.Place.Name, Location: in.Place.Location}
	info.Type = uint8(in.Type)
	info.Opuses = make([]proxy.OpusInfo, 0, 1)
	info.SubmitLimit = uint8(in.Limit)
	info.Quotes = in.Quotes
	if in.Prize != nil {
		info.Prize = &proxy.PrizeInfo{
			Name:  in.Prize.Name,
			Desc:  in.Prize.Desc,
			Ranks: make([]proxy.RankInfo, 0, len(in.Prize.Ranks)),
		}
		for _, rank := range in.Prize.Ranks {
			info.Prize.Ranks = append(info.Prize.Ranks, proxy.RankInfo{Name: rank.Name, Index: rank.Index, Count: rank.Count})
		}
	} else {
		info.Prize = nil
	}
	err := cache.Context().CreateActivity(info)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchActivity(in.Owner, info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyActivityInfo) error {
	path := "activity.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetActivity(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchActivity(in.Owner, info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "activity.getStatistic"
	inLog(path, in)
	out.Key = in.Key
	if in.Key == "create" {
		out.Count = cache.Context().GetActivityCount(in.Owner)
	} else if in.Key == "clone" {
		out.Count = cache.Context().GetActivityCloneCount(in.Owner)
	} else if in.Key == "template_participant" {
		list := cache.Context().GetActivitiesByTemplate(in.Owner, in.Value)
		for _, info := range list {
			out.Count = out.Count + info.Participant
		}
	} else if in.Key == "ratio" {
		out.Count = cache.Context().GetActivityRatio(in.Value)
	} else if in.Key == "template" {
		out.Count = cache.Context().GetActivityTemplateCount(in.Value)
	} else if in.Key == "status" {
		st := parseStringToInt(in.Value)
		out.Count = cache.Context().GetActivityCountByStatus(in.Owner, uint32(st))
	} else if in.Key == "participant" {
		out.Count = cache.Context().GetActivityParticipant(in.Owner)
	} else if in.Key == "opuses" {
		out.Count = cache.Context().GetActivityOpusCount(in.Owner, in.Value)
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "activity.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveActivity(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Parent = in.Owner
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) GetList(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyActivityList) error {
	path := "activity.getList"
	inLog(path, in)
	if len(in.Owner) > 1 {
		used := false
		if len(in.Flag) > 0 {
			used = true
		}
		array := cache.Context().GetActivitiesByOwner(in.Owner, used)
		out.List = make([]*pb.ActivityInfo, 0, len(array))
		for _, val := range array {
			out.List = append(out.List, switchActivity(in.Owner, val))
		}
	} else {
		out.List = make([]*pb.ActivityInfo, 0, 1)
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ActivityService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyActivityList) error {
	path := "activity.getByFilter"
	inLog(path, in)
	var array []*cache.ActivityInfo
	var max uint32 = 0
	var pages uint32 = 0
	if in.Key == "targets" {
		max, pages, array = cache.Context().GetActivitiesByTargets(in.Owner, in.List, cache.ActivityStatusPublish, in.Page, in.Number)
	} else if in.Key == "status" {
		arr, er := stringToUints(in.Value, ";")
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
		max, pages, array = cache.Context().GetActivitiesByStatus(in.Owner, arr, in.Page, in.Number)
	} else if in.Key == "opuses" {

	} else if in.Key == "organizer" {
		array = cache.Context().GetActivitiesByOrganizer(in.Value)
	} else if in.Key == "template" {
		array = cache.Context().GetActivitiesByTemplate(in.Owner, in.Value)
	} else if in.Key == "show" {
		st, er := strconv.ParseUint(in.Value, 10, 32)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
		owners, er := stringToArray(in.Owner, ";")
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
		max, pages, array = cache.Context().GetActivitiesByShow(owners, uint8(st), in.Page, in.Number)
	} else if in.Key == "type" {
		st, er := strconv.ParseUint(in.Value, 10, 32)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
		array = cache.Context().GetAllActivitiesByType(in.Owner, uint8(st))
	} else if in.Key == "alive" {
		//当下时间未结束的活动数据
		array = cache.Context().GetAliveActivities(in.Owner)
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

func (mine *ActivityService) UpdateBase(ctx context.Context, in *pb.ReqActivityUpdate, out *pb.ReplyActivityInfo) error {
	path := "activity.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	in.Name = strings.TrimSpace(in.Name)
	info, err := cache.Context().GetActivity(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if len(in.Cover) > 0 && in.Cover != info.Cover {
		err = info.UpdateCover(in.Cover, in.Operator)
	}
	if len(in.Name) > 0 || len(in.Remark) > 0 {
		err = info.UpdateBase(in.Name, in.Remark, in.Require, in.Operator, proxy.DateInfo{Start: in.Date.Start, Stop: in.Date.Stop},
			proxy.PlaceInfo{Name: in.Place.Name, Location: in.Place.Location})
	}
	if uint8(in.Limit) != info.SubmitLimit {
		err = info.UpdateAssetLimit(in.Operator, uint8(in.Limit))
	}

	if !tool.EqualArray(info.Tags, in.Tags) {
		err = info.UpdateTags(in.Operator, in.Tags)
	}

	if !tool.EqualArray(info.Targets, in.Targets) {
		err = info.UpdateTargets(in.Operator, in.Targets)
	}

	if !tool.EqualArray(info.Assets, in.Assets) {
		err = info.UpdateAssets(in.Operator, in.Assets)
	}

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchActivity(info.Owner, info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "activity.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	activity, err := cache.Context().GetActivity(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_Empty)
		return nil
	}
	if in.Key == "participant" {
		if in.Value == "" {
			out.Status = outError(path, "the activity participant value is empty", pbstatus.ResultStatus_Empty)
			return nil
		}
		num, er := strconv.ParseUint(in.Value, 10, 32)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_FormatError)
			return nil
		}

		er = activity.UpdateParticipant(uint32(num))
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	} else if in.Key == "access" {
		num, er := strconv.ParseUint(in.Value, 10, 32)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_FormatError)
			return nil
		}
		er = activity.UpdateAccess(in.Operator, uint8(num))
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	} else if in.Key == "postpone" {
		//活动延期
		stop := proxy.DateToUTC(in.Value, 1)
		er := activity.UpdateStopDate(in.Operator, stop)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	} else if in.Key == "poster" {
		//修改活动海报
		er := activity.UpdatePoster(in.Operator, in.Value)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) UpdateTags(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "activity.updateTags"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	if in.List == nil {
		out.Status = outError(path, "the activity tags is nil", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetActivity(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = info.UpdateTags(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Tags
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) UpdateAssets(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "activity.updateAssets"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetActivity(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err = info.UpdateAssets(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Assets
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) UpdateTargets(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "activity.updateTargets"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	if in.List == nil {
		out.Status = outError(path, "the activity targets is nil", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetActivity(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = info.UpdateTargets(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Targets
	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) UpdateStatus(ctx context.Context, in *pb.RequestState, out *pb.ReplyInfo) error {
	path := "activity.updateTargets"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetActivity(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = info.UpdateStatus(in.Operator, in.Remark, uint8(in.Status))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) UpdateOpuses(ctx context.Context, in *pb.ReqActivityOpuses, out *pb.ReplyInfo) error {
	path := "activity.updateOpuses"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetActivity(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	list := make([]proxy.OpusInfo, 0, len(in.List))
	for _, opus := range in.List {
		list = append(list, proxy.OpusInfo{Rank: opus.Rank, Asset: opus.Asset, Remark: opus.Remark})
	}
	err = info.UpdateOpuses(in.Operator, list)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) UpdatePrize(ctx context.Context, in *pb.ReqActivityPrize, out *pb.ReplyInfo) error {
	path := "activity.updatePrize"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetActivity(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	list := make([]proxy.RankInfo, 0, len(in.Ranks))
	for _, rank := range in.Ranks {
		list = append(list, proxy.RankInfo{Index: rank.Index, Name: rank.Name, Count: rank.Count})
	}
	err = info.UpdatePrize(in.Operator, in.Name, in.Desc, list)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) UpdateShow(ctx context.Context, in *pb.RequestState, out *pb.ReplyInfo) error {
	path := "activity.updateShow"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetActivity(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = info.UpdateShowState(in.Operator, uint8(in.Status))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) AppendOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPairList) error {
	path := "activity.appendEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	out.Status = outError(path, "the fun not implement", pbstatus.ResultStatus_Empty)
	//info := cache.Context().GetActivity(in.Uid)
	//if info == nil {
	//	out.Status = outError(path,"the activity not found", pbstatus.ResultStatus_NotExisted)
	//	return nil
	//}
	//err := info.AppendPerson(in.Owner, in.Flag)
	//if err != nil {
	//	out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
	//	return nil
	//}
	//out.List = info.GetEntities()
	//out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) SubtractOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyPairList) error {
	path := "activity.subtractEntity"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the activity uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	out.Status = outError(path, "the fun not implement", pbstatus.ResultStatus_Empty)
	//info := cache.Context().GetActivity(in.Uid)
	//if info == nil {
	//	out.Status = outError(path,"the activity not found", pbstatus.ResultStatus_NotExisted)
	//	return nil
	//}
	//err := info.SubtractPerson(in.Owner)
	//if err != nil {
	//	out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
	//	return nil
	//}
	//out.List = info.GetEntities()
	//out.Status = outLog(path, out)
	return nil
}

func (mine *ActivityService) GetStrings(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyList) error {
	path := "activity.getStrings"
	inLog(path, in)
	if in.Key == "tags" {
		out.List = cache.Context().GetActivityTags(in.Value)
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}
