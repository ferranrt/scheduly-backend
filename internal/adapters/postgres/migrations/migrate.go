package migrations

import (
	_ "github.com/lib/pq"
	"gorm.io/gorm"
	"scheduly.io/core/internal/adapters/postgres/dbmodels"
)

func Migrate(db *gorm.DB) error {
	// Models to migrate
	models := []interface{}{
		&dbmodels.User{},
		&dbmodels.Session{},
	}

	// Auto-migrate all models
	return db.AutoMigrate(models...)
}
