package grpc

import (
	"context"
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
	tmp.Entry = info.Entry
	tmp.Menus = info.Menus
	tmp.Revises = info.Revises
	tmp.Effects = make([]*pb.ProductEffect, 0, len(info.Effects))
	for _, effect := range info.Effects {
		tmp.Effects = append(tmp.Effects, &pb.ProductEffect{Min: effect.Min, Max: effect.Max, Pattern: effect.Pattern})
	}
	return tmp
}

func (mine *ProductService) AddOne(ctx context.Context, in *pb.ReqProductAdd, out *pb.ReplyProductInfo) error {
	path := "product.addOne"
	inLog(path, in)
	if len(in.Entry) < 1 {
		out.Status = outError(path, "the owner is empty", pbstatus.ResultStatus_Empty)
		return nil
	}
	effects := make([]*proxy.ProductEffect, 0, len(in.Effects))
	for _, effect := range in.Effects {
		effects = append(effects, &proxy.ProductEffect{Pattern: effect.Pattern, Min: effect.Min, Max: effect.Max})
	}
	info, err := cache.Context().CreateProduct(in.Name, in.Key, in.Entry, in.Menus, in.Templet, in.Remark, in.Operator, uint8(in.Type), in.Revises, effects)
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
	} else if in.Key == "entry" {
		err = info.UpdateEntry(in.Value, in.Operator)
	} else if in.Key == "menus" {
		err = info.UpdateMenus(in.Value, in.Operator)
	} else if in.Key == "revises" {
		err = info.UpdateRevises(in.Operator, in.List)
	} else if in.Key == "effect" {
		err = info.UpdateCatalogs(in.Value, in.Operator)
	} else if in.Key == "catalogs" {
		err = info.UpdateCatalogs(in.Value, in.Operator)
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
