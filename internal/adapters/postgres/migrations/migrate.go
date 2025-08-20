package migrations

import (
	"bifur.app/core/internal/adapters/postgres/dbmodels"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Models to migrate
	models := []interface{}{
		&dbmodels.User{},
		&dbmodels.Source{},
		&dbmodels.Center{},
	}

	// Auto-migrate all models
	return db.AutoMigrate(models...)
}
