package zen

import "gorm.io/gorm"

type Scope func(db *gorm.DB) *gorm.DB
