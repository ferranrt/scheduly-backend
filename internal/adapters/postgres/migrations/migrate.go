package migrations

import (
	"ferranrt.com/scheduly-backend/internal/adapters/postgres/dbmodels"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
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
