package cache

import (
	"encoding/json"

	"github.com/nickolation/grpc-task-hezzl/model"
)
 
func marshalUserList(list []model.User) ([]byte, error) {
	return json.Marshal(list)
}