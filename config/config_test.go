package config

import (
	"testing"
)

// AppConfig app config
type AppConfig struct {
	App struct {
		Mode string
		Host string
		Port int
	}
}

func TestConfig(t *testing.T) {
	err := Load("../testdata/config/")
	// err := Load("../testdata/config/app.json")
	if err != nil {
		t.Error(err)
	}

	testConfig(t)
}

func testConfig(t *testing.T) {
	mode := Get("app.mode")
	t.Logf("mode: %s", mode)
	mode2 := File("gaia").Get("app.mode")
	t.Logf("mode2: %s", mode2)

	addr := File("app").GetInt("http.addr")
	t.Logf("http: %d", addr)

	var appConfig AppConfig
	gaia := File("gaia")
	err := gaia.Unmarshal(&appConfig)
	if err != nil {
		t.Errorf("error: %d", err)
	}

	t.Logf("AppConfig: %v", appConfig)
}
