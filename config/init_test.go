package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetCfgFile(t *testing.T) {
	SetCfgFile()

	appName := viper.GetString("app.name")
	assert.Equal(t, "MyApp", appName)
}
