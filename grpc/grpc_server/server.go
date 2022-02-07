package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/nickolation/grpc-task-hezzl/cache"
	"github.com/nickolation/grpc-task-hezzl/clicklogs/producer"
	"github.com/nickolation/grpc-task-hezzl/configs"
	pb "github.com/nickolation/grpc-task-hezzl/grpc/proto"
	"github.com/nickolation/grpc-task-hezzl/internal/repository"
	"github.com/nickolation/grpc-task-hezzl/internal/service"
	"github.com/nickolation/grpc-task-hezzl/model"
	"github.com/nickolation/grpc-task-hezzl/storage/db"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type grpcUserActionsServer struct {
	service service.GrpcUsersService

	pb.UnimplementedUserActionsServiceServer
}

func newGrpcUserActionsServer(s service.GrpcUsersService) *grpcUserActionsServer {
	return &grpcUserActionsServer{
		service: s,
	}
}

func (srv *grpcUserActionsServer) NewUser(ctx context.Context,
	r *pb.NewUserRequest) (*pb.NewUserResponse, error) {
	u := &model.User{
		Username:    r.GetUsername(),
		Password:    r.GetPassword(),
		Gender:      r.GetGender(),
		Age:         int(r.GetAge()),
		Description: r.GetDescription(),
	}
	id, err := srv.service.NewUser(ctx, u)
	return &pb.NewUserResponse{
		InsertStatus: err == nil,
		PostgresId:   int32(id),
	}, err
}

func (srv *grpcUserActionsServer) DeleteUser(ctx context.Context,
	r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	n := r.GetUsername()
	err := srv.service.DeleteUser(n)
	return &pb.DeleteUserResponse{
		DeleteStatus:    err == nil,
		DeletedUsername: n,
	}, err
}

func (srv *grpcUserActionsServer) GetUserStringedList(ctx context.Context,
	r *pb.GetUserListRequest) (*pb.GetUserStringedListResponse, error) {
	var err error
	var s string

	if r.AllUsersStatus {
		s, err = srv.service.GetUserStringedList(ctx)
	}

	return &pb.GetUserStringedListResponse{
		SelectAllStatus:  err == nil,
		StringJSONResult: s,
	}, err
}

func convertUsersToPbUserList(l []model.User) []*pb.User {
	r := make([]*pb.User, 0)
	for _, u := range l {
		r = append(r, &pb.User{
			Id:          int32(u.Id),
			Username:    u.Username,
			Age:         int32(u.Age),
			Description: u.Description,
			Gender:      u.Gender,
			Password:    u.Password,
			Hash:        u.Hash,
		})
	}

	return r
}

func (srv *grpcUserActionsServer) GetUserList(ctx context.Context,
	r *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	var err error
	l := make([]model.User, 0)

	if r.AllUsersStatus {
		l, err = srv.service.GetUserList(ctx)
	}

	return &pb.GetUserListResponse{
		SelectAllStatus: err == nil,
		UserList:        convertUsersToPbUserList(l),
	}, err
}

func main() {
	errChan := make(chan error, 0)
	go func() {
		logrus.Fatalf("error with the setupping and setting the grpc-server - [%v]", <-errChan)
	}()

	if err := godotenv.Load(); err != nil {
		errChan <- err
	}

	cfg, err := configs.NewGrpcConfigManager()
	if err != nil {
		errChan <- err
	}

	d, err := db.ConnectToPostgresDb(
		db.WithPort(cfg.GetDbPort()),
		db.WithHost(cfg.GetDbHost()),
		db.WithDbName(cfg.GetDbName()),
		db.WithDisabledSSLMode(),
		db.WithPassword(os.Getenv("POSTGRES_PASSWORD")),
		db.WithUsername(cfg.GetDbUsername()),
	)

	if err != nil {
		errChan <- err
	}

	r := repository.NewGrpsUserRepository(d)
	c := cache.NewGrpcCacheStorage(cfg)
	l := producer.NewKafkaHouseClusterManager(cfg)
	s := service.NewGrpcUsersService(r, c, l)

	grpcUserServer := newGrpcUserActionsServer(s)
	srv := grpc.NewServer()
	pb.RegisterUserActionsServiceServer(srv, grpcUserServer)

	listener, err := net.Listen(cfg.GetGrpcServerProtocol(), cfg.GetGrpcServerAddr())
	if err != nil {
		errChan <- err
	}

	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Printf("err with serves by grpcServer - [%v]", err)
		}
	}() 

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit 

	srv.GracefulStop()
	if err := d.Close(); err != nil {
		logrus.Errorf("error occurred on db disconnect %s", err.Error())
		return
	}
}
