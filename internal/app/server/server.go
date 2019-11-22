package server

import (
	"fmt"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/captcha"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/notification"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/admin"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/config"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/contest"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/resource"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/session"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/statistic"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/submission"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/team"
	timeController "github.com/FlagField/FlagField-Server/internal/app/server/controllers/time"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/user"
	"github.com/FlagField/FlagField-Server/internal/pkg/validator"

	"github.com/FlagField/FlagField-Server/internal/app/server/instance"
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	cfg "github.com/FlagField/FlagField-Server/internal/pkg/config"
)

const (
	Version string = "0.0.1"
)

// @title FlagField-Server
// @version 0.0.1
// @description A CTF platform
// @contact.name mukeran
// @contact.email report@mukeran.com

func showMotd() {
	motd := `   ___ _               ___ _      _     _       __      by mukeran
  / __\ | __ _  __ _  / __(_) ___| | __| |     / _\ ___ _ ____   _____ _ __ 
 / _\ | |/ _` + "`" + ` |/ _` + "`" + ` |/ _\ | |/ _ \ |/ _` + "`" + ` |_____\ \ / _ \ '__\ \ / / _ \ '__|
/ /   | | (_| | (_| / /   | |  __/ | (_| |_____|\ \  __/ |   \ V /  __/ |   
\/    |_|\__,_|\__, \/    |_|\___|_|\__,_|     \__/\___|_|    \_/ \___|_|   
               |___/`
	fmt.Println(motd)
	fmt.Print("\n")
}

func showVersion(debug bool) {
	line := `FlagField-Server v` + Version
	if debug {
		line += " (Debug mode on)"
	}
	fmt.Println(line)
}

type Server struct {
	_gin    *gin.Engine
	_gorm   *gorm.DB
	_redigo *redis.Pool
	config  cfg.Config
	public  *instance.Instance
}

func (s *Server) checkDatabase() {
	var name string
	err := s._gorm.Raw("select database()").Row().Scan(&name)
	if err != nil {
		panic("No selected database")
	}
}

func (s *Server) SetConfig(config cfg.Config) {
	s.config = config
}

func (s *Server) initGorm() {
	var err error
	s._gorm, err = gorm.Open(s.config.Database.Type, s.config.Database.Parameter)
	s._gorm.SingularTable(true)
	if err != nil {
		panic(err)
	}
	s.checkDatabase()
}

func (s *Server) initGin() {
	if !s.config.Server.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	s._gin = gin.New()
	s._gin.Use(gin.Logger())
	s._gin.NoRoute(Handler404())
	s._gin.NoMethod(Handler404())
	validator.RegisterGin() // Register extra validation types
}

func (s *Server) initRedigo() {
	s._redigo = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", s.config.Redis.URI,
				redis.DialPassword(s.config.Redis.Password),
				redis.DialDatabase(s.config.Redis.Db),
				redis.DialConnectTimeout(time.Duration(s.config.Redis.ConnectionTimeout)*time.Second),
				redis.DialReadTimeout(time.Duration(s.config.Redis.ReadTimeout)*time.Second),
				redis.DialWriteTimeout(time.Duration(s.config.Redis.WriteTimeout)*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
		MaxIdle:     s.config.Redis.MaxIdleConnections,
		MaxActive:   s.config.Redis.MaxActiveConnections,
		IdleTimeout: time.Duration(s.config.Redis.MaxIdleTimeout) * time.Second,
		Wait:        true,
	}
	rds := s._redigo.Get()
	if rds.Err() != nil {
		panic(rds.Err())
	}
	err := rds.Close()
	if err != nil {
		panic(err)
	}
}

func (s *Server) RegisterRouter() {
	v1 := s._gin.Group("/v1")
	v1.Use(Recovery(), InitRequest(s.public), ResumeDefault(), ResumeSession(s.config.Session), Transaction())
	admin.RegisterRouter(v1)
	captcha.RegisterRouter(v1)
	config.RegisterRouter(v1)
	contest.RegisterRouter(v1)
	notification.RegisterRouter(v1)
	resource.RegisterRouter(v1)
	session.RegisterRouter(v1)
	statistic.RegisterRouter(v1)
	submission.RegisterRouter(v1)
	team.RegisterRouter(v1)
	user.RegisterRouter(v1)
	timeController.RegisterRouter(v1)
}

func (s *Server) Run() {
	var err error
	/* Show information */
	showMotd()
	showVersion(s.config.Server.DebugMode)
	/* Initialize gin, gorm and redigo */
	s.initGorm()
	s.initGin()
	s.initRedigo()
	/* RegisterRouter public */
	s.public = instance.New(s._gin, s._gorm, s._redigo, &s.config)
	/* Register routers */
	s.RegisterRouter()
	/* Run gin app */
	err = s._gin.Run(fmt.Sprintf("%s:%d", s.config.Server.Host, int(s.config.Server.Port)))
	if err != nil {
		panic(err)
	}
}

func New() *Server {
	server := new(Server)
	server.config = cfg.DefaultConfig
	return server
}
