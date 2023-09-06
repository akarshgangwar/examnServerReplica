package migrations

import (
	"examn_go/infra/database"
	"examn_go/models"
)

// Migrate Add list of model add for migrations
// TODO later separate migration each models
func Migrate() {
	var migrationModels = []interface{}{&models.User{}, &models.Exam{}, &models.Question{}, &models.Option{}, &models.Attempt{}, &models.AttemptLog{}, &models.UserEvent{}, &models.SavedAnswer{} }
	db := database.DB
	err := db.AutoMigrate(migrationModels...)
	// err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.User{}, &models.Exam{}, &models.Question{}, &models.Option{})
	if err != nil {
		// throw error in go server
		panic(err)
	}
}
