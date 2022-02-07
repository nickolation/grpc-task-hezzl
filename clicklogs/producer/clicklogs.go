package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nickolation/grpc-task-hezzl/model"
)

const _formatClickhouseDateLayout = "2006-01-02 15:04:05"

type ClickHouseLog struct {
	PostgresId int `json:"postgres_id" form:"postgres_id"`
	Age        int `json:"age" form:"age"`

	Username    string `json:"username" form:"username"`
	Gender      string `json:"gender" form:"gender"`
	Description string `json:"description" form:"description"`
	Time        string `json:"time" form:"time"`
}

func logUserInfo(user *model.User) ([]byte, error) {
	u := &ClickHouseLog{
		PostgresId:  user.Id,
		Username:    user.Username,
		Age:         user.Age,
		Gender:      user.Gender,
		Description: user.Description,
		Time:        time.Now().Format(_formatClickhouseDateLayout),
	}

	return json.Marshal(u)
}

func (khp *KafkaHouseProducer) PushClickhouseUserLog(ctx context.Context, user *model.User) error {
	b, err := logUserInfo(user)
	if err != nil {
		return err
	}

	return khp.PushMessage(ctx, b)
}
