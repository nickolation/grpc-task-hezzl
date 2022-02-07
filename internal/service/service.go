package service

import (
	"context"

	"github.com/nickolation/grpc-task-hezzl/cache"
	"github.com/nickolation/grpc-task-hezzl/clicklogs/producer"
	"github.com/sirupsen/logrus"

	repo "github.com/nickolation/grpc-task-hezzl/internal/repository"
	"github.com/nickolation/grpc-task-hezzl/model"
)

type GrpcUsersService interface {
	NewUser(ctx context.Context, u *model.User) (int, error)
	DeleteUser(username string) error
	GetUserList(ctx context.Context) ([]model.User, error)
	GetUserStringedList(ctx context.Context) (string, error)
}

type Service struct {
	Repo repo.GrpcUsersRepository
	Cache cache.GrpcCacheStorage
	Clicklog producer.KafkaHouseClusterProducer
}

func NewGrpcUsersService(r repo.GrpcUsersRepository, c cache.GrpcCacheStorage, l producer.KafkaHouseClusterProducer) *Service {
	return &Service{
		Repo: r,
		Cache: c,
		Clicklog: l,
	}
}

func (svice *Service) NewUser(ctx context.Context, u *model.User) (int, error) {
	u.SetHashUserValues()
	id, err := svice.Repo.NewUser(u)
	if err != nil {
		logrus.Errorf("new user postgres inserting - [%v]", err.Error())
		return 0, err
	}

	u.Id = id
	if err := svice.Clicklog.PushClickhouseUserLog(ctx, u); err != nil {
		return 0, err
	}

	return id, nil
}

func (svice *Service) DeleteUser(username string) error {
	return svice.Repo.DeleteUser(username)
}

func (svice *Service) GetUserList(ctx context.Context) ([]model.User, error) {
	l, err := svice.Repo.GetUserList()
	if err != nil {
		logrus.Errorf("get userList service - repo err - [%v]", err.Error())
		return l, err
	}

	return l, svice.Cache.SetToCasheUserList(ctx, l)
} 

func (svice *Service) GetUserStringedList(ctx context.Context) (string, error) {
	l, err := svice.Repo.GetUserList()
	if err != nil {
		logrus.Errorf("get userList service - repo err - [%v]", err.Error())
		return "", err
	}

	if err := svice.Cache.SetToCasheUserList(ctx, l); err != nil {
		return "", err  
	}

	return svice.Cache.GetCachedUserList(ctx)
}

