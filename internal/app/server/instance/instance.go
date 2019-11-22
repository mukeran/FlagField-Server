package instance

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"

	cfg "github.com/FlagField/FlagField-Server/internal/pkg/config"
)

type Instance struct {
	ginInstance  *gin.Engine
	gormInstance *gorm.DB
	redigoPool   *redis.Pool
	config       *cfg.Config
}

func (instance *Instance) GetGin() *gin.Engine {
	return instance.ginInstance
}

func (instance *Instance) GetDB() *gorm.DB {
	return instance.gormInstance
}

func (instance *Instance) GetRedis() *redis.Pool {
	return instance.redigoPool
}

func (instance *Instance) GetConfig() *cfg.Config {
	return instance.config
}

func New(ginInstance *gin.Engine, gormInstance *gorm.DB, redigoPool *redis.Pool, config *cfg.Config) *Instance {
	return &Instance{
		ginInstance:  ginInstance,
		gormInstance: gormInstance,
		redigoPool:   redigoPool,
		config:       config,
	}
}
