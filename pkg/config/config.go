package config

import (
	"fmt"
	"strings"

	log "github.com/mgutz/logxi/v1"

	"github.com/mecode4food/cr-clan-bot/pkg/environment"
	"github.com/spf13/viper"
)

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	viper.SetConfigName(Environment())
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Warn("couldn't read configuration file", "environment", Environment(), "error", err)
	} else {
		log.Info("configuration loaded", "environment", Environment())
	}

	v := secretsViper()
	if err := viper.MergeConfigMap(v.AllSettings()); err != nil {
		log.Warn("couldn't read secrets file", "environment", Environment(), "error", err)
	} else {
		log.Info("secrets config merged", "environment", Environment())
	}
}

func secretsViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName(fmt.Sprintf("secrets/%s", Environment()))
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		log.Warn("couldn't read secret file", "environment", Environment(), "error", err)
	} else {
		log.Info("secrets loaded", "environment", Environment())
	}

	return v
}

func Environment() string {
	e := viper.GetString("APPLICATION_ENVIRONMENT")
	if e == "" {
		e = environment.PlatformDevelopment
	}
	return e
}

func IsDevelopment() bool {
	return Environment() == environment.PlatformDevelopment
}

func IsProduction() bool {
	return Environment() == environment.PlatformProduction
}

func IsTesting() bool {
	return Environment() == environment.PlatformTesting
}

func Viper() *viper.Viper {
	return viper.GetViper()
}
