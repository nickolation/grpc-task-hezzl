package configs

import (
	"strings"

	"github.com/spf13/viper"
)

type GrpcGonfigManager struct {
	Engine *viper.Viper

	sourceConfigPath string
	settingsPool     []string
}

const _defaultSourceConfigPath = "configs/settings"
const (
	_grpcServiceSettings = "grpc_service"
	_kafkaSettings       = "kafka"
	_storageSettings     = "storage"
)

func (cfg *GrpcGonfigManager) riseConfigs() error {
	for _, f := range cfg.settingsPool {
		cfg.Engine.AddConfigPath(cfg.sourceConfigPath)
		cfg.Engine.SetConfigName(f)
		err := cfg.Engine.MergeInConfig()
		if err != nil {
			return err
		}
	}

	return nil
}

func NewGrpcConfigManager() (*GrpcGonfigManager, error) {
	cfg := &GrpcGonfigManager{
		Engine:           viper.New(),
		sourceConfigPath: _defaultSourceConfigPath,
		settingsPool: []string{
			_grpcServiceSettings,
			_kafkaSettings,
			_storageSettings,
		},
	}

	err := cfg.riseConfigs()
	return cfg, err
}

// Postgres Db
func (cfg *GrpcGonfigManager) GetDbPort() string {
	const _dbPortTemplate = "db.port"
	return cfg.Engine.GetString(_dbPortTemplate)
}

func (cfg *GrpcGonfigManager) GetDbHost() string {
	const _dbHostTemplate = "db.host"
	return cfg.Engine.GetString(_dbHostTemplate)
}

func (cfg *GrpcGonfigManager) GetDbUsername() string {
	const _dbUsernameTemplate = "db.username"
	return cfg.Engine.GetString(_dbUsernameTemplate)
}

func (cfg *GrpcGonfigManager) GetDbName() string {
	const _dbNameTemplate = "db.dbname"
	return cfg.Engine.GetString(_dbNameTemplate)
}

// Redis
func (cfg *GrpcGonfigManager) GetRedisAddr() string {
	const _redisAddrTemplate = "redis.addr"
	return cfg.Engine.GetString(_redisAddrTemplate)
}

func (cfg *GrpcGonfigManager) GetRedisCacheTTL() int {
	const _redisCacheTTLTemplate = "redis.cache_ttl"
	return cfg.Engine.GetInt(_redisCacheTTLTemplate)
}

// Kafka
func (cfg *GrpcGonfigManager) GetKafkaClusterTopic() string {
	const _kafkaClusterTopicTemplate = "kafka.cluster_topic"
	return cfg.Engine.GetString(_kafkaClusterTopicTemplate)
}

func (cfg *GrpcGonfigManager) GetKafkaCliendId() string {
	const _kafkaClientIdTemplate = "kafka.client_id"
	return cfg.Engine.GetString(_kafkaClientIdTemplate)
}

func (cfg *GrpcGonfigManager) GetKafkaAddrs() []string {
	const (
		_kafkaAddrsTemplate = "kafka.addrs"
		_brockerStringSep   = ","
	)
	return strings.Split(cfg.Engine.GetString(_kafkaAddrsTemplate), _brockerStringSep)
}

// Grpc
func (cfg *GrpcGonfigManager) GetGrpcServerAddr() string {
	const _grpcServerAddr = "grpc.server.addr"
	return cfg.Engine.GetString(_grpcServerAddr)
}

func (cfg *GrpcGonfigManager) GetGrpcServerProtocol() string {
	const _grpcServerProtocol = "grpc.server.protocol"
	return cfg.Engine.GetString(_grpcServerProtocol)
}

func (cfg *GrpcGonfigManager) GetGrpcClientAddr() string {
	const _grpcClientAddr = "grpc.client_api.addr"
	return cfg.Engine.GetString(_grpcClientAddr)
}
