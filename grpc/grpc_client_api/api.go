package grpcclientapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nickolation/grpc-task-hezzl/configs"
	pb "github.com/nickolation/grpc-task-hezzl/grpc/proto"
	"github.com/nickolation/grpc-task-hezzl/model"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClientAPIServer struct {
	router *gin.Engine
	server *http.Server

	config *configs.GrpcGonfigManager

	client pb.UserActionsServiceClient
}

const (
	_badGrpcNewUserMethodTemplate = "bad grpc invoke - "
	_invalidMessageBodyTemplate   = "invalid message body for user struct"
)

func (gcs *grpcClientAPIServer) newUser(c *gin.Context) {
	user := &model.User{}

	if err := c.BindJSON(user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, _invalidMessageBodyTemplate)
		return
	}

	resp, err := gcs.client.NewUser(c, &pb.NewUserRequest{
		Username:    user.Username,
		Age:         int32(user.Age),
		Description: user.Description,
		Gender:      user.Gender,
		Password:    user.Password,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, _badGrpcNewUserMethodTemplate+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (gcs *grpcClientAPIServer) deleteUser(c *gin.Context) {
	const _usernameParamTemplate = "username"
	u := c.Param(_usernameParamTemplate)

	resp, err := gcs.client.DeleteUser(c, &pb.DeleteUserRequest{Username: u})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, _badGrpcNewUserMethodTemplate+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

const _defaultAllUsersStatus = true

func (gcs *grpcClientAPIServer) getUserStringedList(c *gin.Context) {
	resp, err := gcs.client.GetUserStringedList(c, &pb.GetUserListRequest{
		AllUsersStatus: _defaultAllUsersStatus,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, _badGrpcNewUserMethodTemplate+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (gcs *grpcClientAPIServer) getUserList(c *gin.Context) {
	resp, err := gcs.client.GetUserList(c, &pb.GetUserListRequest{
		AllUsersStatus: _defaultAllUsersStatus,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, _badGrpcNewUserMethodTemplate+err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

const (
	_testApiUrl      = "/api/test"
	_getListGroupUrl = "/list"

	_newUserApiUrl          = "/new"
	_deleteUserApiUrl       = "/delete/:username"
	_getStringedUserListUrl = "/string"
	_getUserListUrl         = "/list"
)

func (gcs *grpcClientAPIServer) initTestGroupRouter() {
	r := gin.Default()
	api := r.Group(_testApiUrl)
	{
		api.POST(_newUserApiUrl, gcs.newUser)
		api.DELETE(_deleteUserApiUrl, gcs.deleteUser)

		get := api.Group(_getListGroupUrl)
		{
			get.GET(_getStringedUserListUrl, gcs.getUserStringedList)
			get.GET(_getUserListUrl, gcs.getUserList)
		}

	}

	gcs.router = r
}

func (gcs *grpcClientAPIServer) initClientServer() {
	gcs.server = &http.Server{
		Addr:           gcs.config.GetGrpcClientAddr(),
		MaxHeaderBytes: 1 << 20,
		Handler:        gcs.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
}

func (gcs *grpcClientAPIServer) RunClientServer() error {
	return gcs.server.ListenAndServe()
} 

func (gcs *grpcClientAPIServer) ShutdownClientServer(ctx context.Context) error {
	return gcs.server.Shutdown(ctx)
}

func NewClientAPIServer(cfg *configs.GrpcGonfigManager) *grpcClientAPIServer {
	dial, err := grpc.Dial(cfg.GetGrpcServerAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logrus.Errorf("error with client dialer - [%v]", err.Error())
		return nil
	}

	c := &grpcClientAPIServer{
		client: pb.NewUserActionsServiceClient(dial),
		config: cfg,
	}
	c.initTestGroupRouter()
	c.initClientServer()

	return c
}
