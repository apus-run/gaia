package config

import (
	"testing"

	"github.com/apus-run/gaia/config/file"
)

func TestConfig(t *testing.T) {
	c := New(WithSource(
		file.NewFile("../../config/"),
		// file.NewFile("../../config/sea.yaml"),
		// file.NewFile("../../"),
	))

	testConfig(t, c)
}

func testConfig(t *testing.T, c Config) {
	if err := c.Load(); err != nil {
		t.Error(err)
	}
	d, _ := Sub("database")
	dsn := d.GetString("dsn")
	t.Logf("mode: %s", dsn)
	mode := Get("app.mode")
	t.Logf("mode: %s", mode)
	mode2 := File("sea").Get("app.mode")
	t.Logf("mode2: %s", mode2)

	addr := File("test").GetInt("http.addr")
	t.Logf("http: %d", addr)

	// AppConfig app config
	type AppConfig struct {
		App struct {
			Mode  string
			Grace bool
			Host  string
			Port  int
		}
	}

	var appConfig AppConfig
	sea := File("sea")
	err := sea.viper.Unmarshal(&appConfig)
	if err != nil {
		t.Errorf("error: %d", err)
	}

	t.Logf("AppConfig: %v", appConfig)
}
