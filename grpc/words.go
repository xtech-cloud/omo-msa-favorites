package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
)

type WordsService struct{}

func switchWords(info *cache.WordsInfo) *pb.WordsInfo {
	tmp := new(pb.WordsInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Words = info.Words
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Owner = info.Owner
	tmp.Target = info.Target
	tmp.Weight = uint32(info.Weight)
	tmp.Assets = info.Assets
	tmp.Type = uint32(info.Type)
	tmp.Device = info.Device
	return tmp
}

func (mine *WordsService) AddOne(ctx context.Context, in *pb.ReqWordsAdd, out *pb.ReplyWordsInfo) error {
	path := "words.addOne"
	inLog(path, in)
	if len(in.Owner) < 1 {
		out.Status = outError(path, "the owner is empty", pbstatus.ResultStatus_Empty)
		return nil
	}

	if len(in.Target) < 1 {
		in.Target = in.Owner
	}

	info, err := cache.Context().CreateWords(in.Words, in.Owner, in.Target, in.Device, in.Operator, in.Quote, in.Assets, cache.WordsType(in.Type))
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchWords(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *WordsService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyWordsInfo) error {
	path := "words.getOne"
	inLog(path, in)

	var info *cache.WordsInfo
	if len(in.Flag) < 2 {
		info = cache.Context().GetWords(in.Uid)
	} else if in.Flag == "today" {
		info = cache.Context().GetWordsByToday(in.Owner, in.Operator, in.Uid)
	}

	if info == nil {
		out.Status = outError(path, "the words not found", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchWords(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *WordsService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "words.getStatistic"
	inLog(path, in)
	if in.Key == "template" {
		arr := cache.Context().GetWordsByQuote(in.Owner, in.Value)
		out.Count = uint32(len(arr))
	} else if in.Key == "device" {
		out.Count, _ = cache.Context().GetWordsCountByDevice(in.Value)
	} else if in.Key == "today" {
		out.Count, _ = cache.Context().GetWordsCountByToday(in.Value)
	} else {
		out.Count, _ = cache.Context().GetWordsCountByScene(in.Owner)
	}
	out.Owner = in.Owner
	out.Key = in.Key
	out.Status = outLog(path, out)
	return nil
}

func (mine *WordsService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "words.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the words uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveWords(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Parent = in.Owner
	out.Status = outLog(path, out)
	return nil
}

func (mine *WordsService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyWordsList) error {
	path := "words.getByFilter"
	inLog(path, in)
	var array []*cache.WordsInfo
	var max uint32 = 0
	var pages uint32 = 0
	if in.Key == "" {
		array = cache.Context().GetWordsByOwner(in.Owner)
	} else if in.Key == "target" {
		array = cache.Context().GetWordsByTarget(in.Value)
	} else if in.Key == "type" {
		tp := parseStringToInt(in.Value)
		array = cache.Context().GetWordsByOwnerTP(in.Owner, cache.WordsType(tp))
	} else if in.Key == "user" {
		array = cache.Context().GetWordsByUser(in.Value)
	} else if in.Key == "quote" {
		array = cache.Context().GetWordsByQuote(in.Owner, in.Value)
	} else {

	}
	out.List = make([]*pb.WordsInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchWords(val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *WordsService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "words.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the words uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	var info *cache.WordsInfo
	var err error
	info = cache.Context().GetWords(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the words", pbstatus.ResultStatus_Empty)
		return nil
	}
	if in.Key == "weight" {
		w := parseStringToInt(in.Value)
		err = info.UpdateWeight(int32(w), in.Operator)
	} else if in.Key == "assets" {
		err = info.UpdateAssets(in.List, in.Operator)
	} else if in.Key == "words" {
		err = info.UpdateBase(in.Value, in.Operator)
	} else {
		err = errors.New("not defined the field key")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_Empty)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
