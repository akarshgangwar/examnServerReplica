package main

import (
	"examn_go/config"
	"examn_go/infra/database"
	"examn_go/infra/logger"
	"examn_go/migrations"
	"examn_go/routers"
	"github.com/spf13/viper"
	"time"
)

func main() {

	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
	masterDSN, replicaDSN := config.DbConfiguration()

	if err := database.DbConnection(masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}
	//later separate migration
	migrations.Migrate()

	router := routers.SetupRoute()
	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}
