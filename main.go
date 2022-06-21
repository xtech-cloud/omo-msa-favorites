package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/registry/etcdv3/v2"
	proto "github.com/xtech-cloud/omo-msp-favorites/proto/favorite"
	"io"
	"omo.msa.favorite/cache"
	"omo.msa.favorite/config"
	"omo.msa.favorite/grpc"
	"os"
	"path/filepath"
	"time"
)

var (
	BuildVersion string
	BuildTime    string
	CommitID     string
)

func main() {
	config.Setup()
	err := cache.InitData()
	if err != nil {
		panic(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("omo.msa.favorite"),
		micro.Version(BuildVersion),
		micro.RegisterTTL(time.Second*time.Duration(config.Schema.Service.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Schema.Service.Interval)),
		micro.Address(config.Schema.Service.Address),
	)
	// Initialise service
	service.Init()
	// Register Handler
	_ = proto.RegisterFavoriteServiceHandler(service.Server(), new(grpc.FavoriteService))
	_ = proto.RegisterActivityServiceHandler(service.Server(), new(grpc.ActivityService))
	_ = proto.RegisterArticleServiceHandler(service.Server(), new(grpc.ArticleService))
	_ = proto.RegisterNoticeServiceHandler(service.Server(), new(grpc.NoticeService))
	_ = proto.RegisterMessageServiceHandler(service.Server(), new(grpc.MessageService))
	_ = proto.RegisterDisplayServiceHandler(service.Server(), new(grpc.DisplayService))
	_ = proto.RegisterSheetServiceHandler(service.Server(), new(grpc.SheetService))
	_ = proto.RegisterWordsServiceHandler(service.Server(), new(grpc.WordsService))

	app, _ := filepath.Abs(os.Args[0])

	logger.Info("-------------------------------------------------------------")
	logger.Info("- Micro Service Agent -> Run")
	logger.Info("-------------------------------------------------------------")
	logger.Infof("- version      : %s", BuildVersion)
	logger.Infof("- application  : %s", app)
	logger.Infof("- md5          : %s", md5hex(app))
	logger.Infof("- build        : %s", BuildTime)
	logger.Infof("- commit       : %s", CommitID)
	logger.Info("-------------------------------------------------------------")
	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

func md5hex(_file string) string {
	h := md5.New()

	f, err := os.Open(_file)
	if err != nil {
		return ""
	}
	defer f.Close()

	io.Copy(h, f)

	return hex.EncodeToString(h.Sum(nil))
}
