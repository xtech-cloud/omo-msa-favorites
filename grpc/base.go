package grpc

import (
	"encoding/json"
	"github.com/micro/go-micro/v2/logger"
	pb "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"strconv"
	"strings"
)

func stringToUints(source, split string) ([]uint8,error) {
	arr := make([]uint8, 0, 3)
	if strings.Contains(source, split) {
		ss := strings.Split(source, split)
		for _, s := range ss {
			st, er := strconv.ParseUint(s, 10, 32)
			if er != nil {
				return nil, er
			}
			arr = append(arr, uint8(st))
		}
	}else{
		st, er := strconv.ParseUint(source, 10, 32)
		if er != nil {
			return nil, er
		}
		arr = append(arr, uint8(st))
	}
	return arr, nil
}

func stringToArray(source, split string) ([]string,error) {
	arr := make([]string, 0, 3)
	if strings.Contains(source, split) {
		ss := strings.Split(source, split)
		for _, s := range ss {
			arr = append(arr, s)
		}
	}else{
		arr = append(arr, source)
	}
	return arr, nil
}

func inLog(name, data interface{})  {
	bytes, _ := json.Marshal(data)
	msg := ByteString(bytes)
	logger.Infof("[in.%s]:data = %s", name, msg)
}

func outError(name, msg string, code pbstatus.ResultStatus) *pb.ReplyStatus {
	logger.Warnf("[error.%s]:code = %d, msg = %s", name, code, msg)
	tmp := &pb.ReplyStatus{
		Code: uint32(code),
		Error: msg,
	}
	return tmp
}

func outLog(name, data interface{}) *pb.ReplyStatus {
	bytes, _ := json.Marshal(data)
	msg := ByteString(bytes)
	logger.Infof("[out.%s]:data = %s", name, msg)
	tmp := &pb.ReplyStatus{
		Code: 0,
		Error: "",
	}
	return tmp
}

func ByteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}
