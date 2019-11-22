package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/FlagField/FlagField-Server/internal/pkg/validator"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	Session  SessionConfig  `json:"session"`
	Resource ResourceConfig `json:"resource"`
}

var (
	DefaultServerConfig = ServerConfig{
		Host:      "0.0.0.0",
		Port:      8080,
		DebugMode: false,
	}
	DefaultDatabaseConfig = DatabaseConfig{
		Type:      "",
		Parameter: "",
	}
	DefaultRedisConfig = RedisConfig{
		URI:                  "localhost:6379",
		MaxIdleConnections:   10,
		MaxActiveConnections: 500,
		MaxIdleTimeout:       10,
		Password:             "",
		Db:                   0,
		ConnectionTimeout:    3,
		ReadTimeout:          3,
		WriteTimeout:         3,
	}
	DefaultSessionConfig = SessionConfig{
		Name:     "session_id",
		MaxAge:   7200,
		Path:     "",
		Domain:   "",
		Secure:   false,
		HttpOnly: false,
	}
	DefaultResourceConfig = ResourceConfig{
		BaseDir: os.Getenv("FLAGFIELD_HOME") + "uploads/",
	}
	DefaultConfig = Config{DefaultServerConfig, DefaultDatabaseConfig, DefaultRedisConfig, DefaultSessionConfig, DefaultResourceConfig}
)

type ServerConfig struct {
	Host      string `json:"host" validate:"ip"`
	Port      uint16 `json:"port"`
	DebugMode bool   `json:"debug_mode"`
}

type DatabaseConfig struct {
	Type      string `json:"type" validate:"eq=mysql|eq=postgres|eq=sqlite3|eq=mssql,required"`
	Parameter string `json:"parameter" validate:"required"`
}

type RedisConfig struct {
	URI                  string `json:"uri" validate:"uri"`
	MaxIdleConnections   int    `json:"max_idle_connections"`
	MaxActiveConnections int    `json:"max_active_connections"`
	MaxIdleTimeout       int    `json:"max_idle_timeout"`
	Password             string `json:"password"`
	Db                   int    `json:"db"`
	ConnectionTimeout    int    `json:"connection_timeout"`
	ReadTimeout          int    `json:"read_timeout"`
	WriteTimeout         int    `json:"write_timeout"`
}

type SessionConfig struct {
	Name     string `json:"name" validate:"required"`
	MaxAge   int    `json:"max_age"`
	Path     string `json:"path"`
	Domain   string `json:"domain"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"http_only"`
}

type ResourceConfig struct {
	BaseDir string `json:"base_dir" validate:"required"`
}

func FromFile(path string) *Config {
	config := DefaultConfig
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}
	err = validator.Validate(config)
	if err != nil {
		panic(err)
	}
	return &config
}

func (c *Config) ToFile(path string) error {
	str, _ := json.Marshal(c)
	return ioutil.WriteFile(path, str, 0777)
}

//func FromRedis(rp *redis.Pool) (*Config, error) {
//	conn := rp.Get()
//	if conn.Err() != nil {
//		panic(conn.Err())
//	}
//	var cfg Config
//	err := importFromRedis(&cfg, &conn, "flagfield")
//	return &cfg, err
//}
//
//func importFromRedis(obj interface{}, r *redis.Conn, path string) error {
//	t := reflect.TypeOf(obj)
//	v := reflect.ValueOf(obj)
//	for t.Kind() == reflect.Ptr {
//		t = t.Elem()
//	}
//	for v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}
//	switch v.Kind() {
//	case reflect.Struct:
//		for i := 0; i < v.NumField(); i++ {
//			err := importFromRedis(v.Field(i).Addr().Interface(), r, fmt.Sprintf("%s.%s", path, strings.ToLower(t.Field(i).Name)))
//			if err != nil {
//				return err
//			}
//		}
//	default:
//		str, err := redis.Bytes((*r).Do("GET", path))
//		if err != nil {
//			panic(err)
//		}
//		err = json.Unmarshal(str, obj)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func exportToRedis(obj interface{}, r *redis.Conn, path string) {
//	t := reflect.TypeOf(obj)
//	v := reflect.ValueOf(obj)
//	for t.Kind() == reflect.Ptr {
//		t = t.Elem()
//	}
//	for v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}
//	switch v.Kind() {
//	case reflect.Struct:
//		for i := 0; i < v.NumField(); i++ {
//			exportToRedis(v.Field(i).Interface(), r, fmt.Sprintf("%s.%s", path, strings.ToLower(t.Field(i).Name)))
//		}
//	default:
//		str, err := json.Marshal(v.Interface())
//		if err != nil {
//			panic(err)
//		}
//		_, err = (*r).Do("SET", path, str)
//		if err != nil {
//			panic(err)
//		}
//	}
//}
//
//func (c *Config) Save(rp *redis.Pool) {
//	conn := rp.Get()
//	if conn.Err() != nil {
//		panic(conn.Err())
//	}
//	exportToRedis(c, &conn, "flagfield")
//}
//
//func Get(rp *redis.Pool, key string, value interface{}) error {
//	conn := rp.Get()
//	if conn.Err() != nil {
//		panic(conn.Err())
//	}
//	str, err := redis.Bytes(conn.Do("GET", key))
//	if err != nil {
//		panic(err)
//	}
//	err = json.Unmarshal(str, value)
//	return err
//}
//
//func Set(rp *redis.Pool, key string, value interface{}) error {
//	conn := rp.Get()
//	if conn.Err() != nil {
//		panic(conn.Err())
//	}
//	str, err := json.Marshal(value)
//	if err != nil {
//		return err
//	}
//	_, err = conn.Do("SET", key, str)
//	if err != nil {
//		panic(err)
//	}
//	return nil
//}
