package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.favorite/cache"
	"omo.msa.favorite/proxy"
)

type ProductService struct{}

func switchProduct(info *cache.ProductInfo) *pb.ProductInfo {
	tmp := new(pb.ProductInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Name = info.Name
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Created = info.CreateTime.Unix()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Type = uint32(info.Type)
	tmp.Status = uint32(info.Status)
	tmp.Remark = info.Remark
	tmp.Key = info.Key
	tmp.Templet = info.Templet
	tmp.Catalogs = info.Catalogs
	tmp.Entries = info.Entries
	tmp.Menus = info.Menus
	tmp.Revises = info.Revises
	tmp.Shows = info.Shows
	tmp.Effects = make([]*pb.ProductEffect, 0, len(info.Effects))
	for _, effect := range info.Effects {
		tmp.Effects = append(tmp.Effects, &pb.ProductEffect{Min: effect.Min, Max: effect.Max,
			Access: uint32(effect.Access), Pattern: effect.Pattern})
	}
	tmp.Displays = make([]*pb.DisplayShow, 0, len(info.Displays))
	for _, display := range info.Displays {
		tmp.Displays = append(tmp.Displays, &pb.DisplayShow{Uid: display.UID, Effect: display.Effect})
	}
	return tmp
}

func (mine *ProductService) AddOne(ctx context.Context, in *pb.ReqProductAdd, out *pb.ReplyProductInfo) error {
	path := "product.addOne"
	inLog(path, in)
	if len(in.Entries) < 1 {
		out.Status = outError(path, "the entries is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	effects := make([]*proxy.ProductEffect, 0, len(in.Effects))
	for _, effect := range in.Effects {
		effects = append(effects, &proxy.ProductEffect{Pattern: effect.Pattern, Min: effect.Min, Max: effect.Max})
	}
	info, err := cache.Context().CreateProduct(in.Name, in.Key, in.Menus, in.Templet, in.Remark, in.Operator, uint8(in.Type), in.Entries, in.Revises, effects)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchProduct(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ProductService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyProductInfo) error {
	path := "product.getOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the product uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err := cache.Context().GetProduct(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.Info = switchProduct(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ProductService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "product.removeOne"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the product uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	err := cache.Context().RemoveProduct(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Uid = in.Uid
	out.Parent = in.Owner
	out.Status = outLog(path, out)
	return nil
}

func (mine *ProductService) GetByFilter(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyProductList) error {
	path := "product.getByFilter"
	inLog(path, in)
	var array []*cache.ProductInfo
	var max uint32 = 0
	var pages uint32 = 0

	if in.Key == "status" {

	} else {
		array = cache.Context().GetAllProducts()
	}
	out.List = make([]*pb.ProductInfo, 0, len(array))
	for _, val := range array {
		out.List = append(out.List, switchProduct(val))
	}
	out.Total = max
	out.Pages = pages
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ProductService) GetStatistic(ctx context.Context, in *pb.RequestFilter, out *pb.ReplyStatistic) error {
	path := "product.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *ProductService) UpdateByFilter(ctx context.Context, in *pb.RequestUpdate, out *pb.ReplyInfo) error {
	path := "product.updateByFilter"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the product uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err1 := cache.Context().GetProduct(in.Uid)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Key == "templet" {
		err = info.UpdateTemplet(in.Value, in.Operator)
	} else if in.Key == "entries" {
		err = info.UpdateEntries(in.Operator, in.List)
	} else if in.Key == "menus" {
		err = info.UpdateMenus(in.Value, in.Operator)
	} else if in.Key == "revises" {
		err = info.UpdateRevises(in.Operator, in.List)
	} else if in.Key == "effects" {
		if len(in.Value) > 1 {
			array := make([]*pb.ProductEffect, 0, 10)
			err = json.Unmarshal([]byte(in.Value), &array)
			if err == nil {
				arr := make([]*proxy.ProductEffect, 0, 10)
				for _, effect := range array {
					arr = append(arr, &proxy.ProductEffect{Pattern: effect.Pattern, Min: effect.Min, Max: effect.Max, Access: uint8(effect.Access)})
				}
				err = info.UpdateEffects(in.Operator, arr)
			}
		}
	} else if in.Key == "catalogs" {
		err = info.UpdateCatalogs(in.Value, in.Operator)
	} else if in.Key == "shows" {
		err = info.UpdateShows(in.Operator, in.List)
	} else if in.Key == "displays" {
		if len(in.Value) > 1 {
			array := make([]*proxy.DisplayShow, 0, 10)
			err = json.Unmarshal([]byte(in.Value), &array)
			if err == nil {
				arr := make([]*proxy.DisplayShow, 0, 10)
				for _, item := range array {
					arr = append(arr, &proxy.DisplayShow{UID: item.UID, Effect: item.Effect})
				}
				err = info.UpdateDisplays(in.Operator, arr)
			}
		}
	} else if in.Key == "type" {
		tp := parseStringToInt(in.Value)
		err = info.UpdateType(in.Operator, uint32(tp))
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ProductService) UpdateBase(ctx context.Context, in *pb.ReqProductUpdate, out *pb.ReplyProductInfo) error {
	path := "product.updateBase"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the product uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err1 := cache.Context().GetProduct(in.Uid)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateBase(in.Name, in.Key, in.Remark, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchProduct(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ProductService) AppendDisplay(ctx context.Context, in *pb.ReqProductDisplay, out *pb.ReplyProductDisplays) error {
	path := "product.appendDisplay"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the product uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err1 := cache.Context().GetProduct(in.Uid)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.List = make([]*pb.DisplayShow, 0, len(info.Displays))
	for _, display := range info.Displays {
		out.List = append(out.List, &pb.DisplayShow{Uid: display.UID, Effect: display.Effect})
	}
	out.Status = outLog(path, out)
	return nil
}

func (mine *ProductService) SubtractDisplay(ctx context.Context, in *pb.ReqProductDisplay, out *pb.ReplyProductDisplays) error {
	path := "product.subtractDisplay"
	inLog(path, in)
	if len(in.Uid) < 1 {
		out.Status = outError(path, "the product uid is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	info, err1 := cache.Context().GetProduct(in.Uid)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.List = make([]*pb.DisplayShow, 0, len(info.Displays))
	for _, display := range info.Displays {
		out.List = append(out.List, &pb.DisplayShow{Uid: display.UID, Effect: display.Effect})
	}
	out.Status = outLog(path, out)
	return nil
}
