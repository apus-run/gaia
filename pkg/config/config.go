package config

import (
	"errors"
	"log"
	"os"
	xpath "path"
	"path/filepath"
	"strings"
	"time"

	"github.com/apus-run/gaia/pkg/config/file"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	defaultFile = "gaia"
	files       map[string]*Config

	// Get 查询配置/环境变量
	Get func(string) string
)

type Config struct {
	*viper.Viper
}

func Load(path string) error {
	src := file.NewSource(path)
	fs, _ := src.Load()
	files = make(map[string]*Config, len(fs))
	for _, f := range fs {
		v, err := load(f)
		if err != nil {
			panic(err)
		}
		name := strings.TrimSuffix(xpath.Base(f), filepath.Ext(f))
		log.Printf("文件名: %s", name)
		files[name] = &Config{v}
	}
	Get = GetString

	return nil
}

func load(fileName string) (*viper.Viper, error) {
	fileType := file.Format(fileName)

	v := viper.New()
	v.SetConfigFile(fileName)
	v.SetConfigType(fileType)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("Using config file: %s [%s]\n", viper.ConfigFileUsed(), err)
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	// 读取匹配的环境变量
	v.AutomaticEnv()

	return v, nil
}

func Watch() {
	for _, v := range files {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("Config file changed: %s", e.Name)
		})
	}
}

// Set 设置配置，仅用于测试
func Set(key string, value interface{}) { File(defaultFile).Set(key, value) }

func GetBool(key string) bool              { return File(defaultFile).GetBool(key) }
func GetDuration(key string) time.Duration { return File(defaultFile).GetDuration(key) }
func GetFloat64(key string) float64        { return File(defaultFile).GetFloat64(key) }
func GetInt(key string) int                { return File(defaultFile).GetInt(key) }
func GetInt32(key string) int32            { return File(defaultFile).GetInt32(key) }
func GetInt64(key string) int64            { return File(defaultFile).GetInt64(key) }
func GetIntSlice(key string) []int         { return File(defaultFile).GetIntSlice(key) }
func GetSizeInBytes(key string) uint       { return File(defaultFile).GetSizeInBytes(key) }
func GetString(key string) string          { return File(defaultFile).GetString(key) }
func GetStringSlice(key string) []string   { return File(defaultFile).GetStringSlice(key) }
func GetTime(key string) time.Time         { return File(defaultFile).GetTime(key) }
func GetUint(key string) uint              { return File(defaultFile).GetUint(key) }
func GetUint32(key string) uint32          { return File(defaultFile).GetUint32(key) }
func GetUint64(key string) uint64          { return File(defaultFile).GetUint64(key) }

func GetStringMap(key string) map[string]interface{} { return File(defaultFile).GetStringMap(key) }
func GetStringMapString(key string) map[string]string {
	return File(defaultFile).GetStringMapString(key)
}
func GetStringMapStringSlice(key string) map[string][]string {
	return File(defaultFile).GetStringMapStringSlice(key)
}

// GetEnvString get value from env.
// application parameters take precedence over environment variables
// env := GetEnvString("APP_ENV", "")
func GetEnvString(key string, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

// File 根据文件名获取对应配置对象
// 如果要读取 foo.yaml 配置，可以 File("foo").Get("bar")
func File(name string) *Config {
	return files[name]
}
