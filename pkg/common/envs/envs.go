package envs

import (
	"github.com/spf13/viper"
	"log"
)

func Init() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error happened while config initialization: ", err)
	}
}
