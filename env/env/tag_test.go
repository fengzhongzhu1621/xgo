package env

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/fengzhongzhu1621/xgo/tests"
)

type config struct {
	Home         string        `env:"HOME"`
	Port         int           `env:"PORT"               envDefault:"3000"`
	Password     string        `env:"PASSWORD,unset"`
	IsProduction bool          `env:"PRODUCTION"`
	Hosts        []string      `env:"HOSTS"                                       envSeparator:":"`
	Duration     time.Duration `env:"DURATION"`
	TempFolder   string        `env:"TEMP_FOLDER,expand" envDefault:"${HOME}/tmp"`
}

func TestEnvParse(t *testing.T) {
	os.Setenv("HOME", "/Users/user")
	os.Setenv("PORT", "8080")
	os.Setenv("PASSWORD", "mysecretpassword")
	os.Setenv("PRODUCTION", "true")
	os.Setenv("HOSTS", "localhost:127.0.0.1")
	os.Setenv("DURATION", "30s")
	os.Setenv("TEMP_FOLDER", "${HOME}/tmp")

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	tests.PrintStruct(cfg)

	// {
	// 		"Home": "/Users/user",
	// 		"Port": 8080,
	// 		"Password": "mysecretpassword",
	// 		"IsProduction": true,
	// 		"Hosts": [
	// 			"localhost",
	// 			"127.0.0.1"
	// 		],
	// 		"Duration": 30000000000,
	// 		"TempFolder": "/Users/user/tmp"
	// }
}
