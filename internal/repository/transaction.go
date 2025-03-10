package repository

import "gorm.io/gorm"

type Transaction interface {
	Commit() error
	Rollback() error
	DB() *gorm.DB
}

type GormTransaction struct {
	tx *gorm.DB
}

func (t *GormTransaction) Commit() error {
	return t.tx.Commit().Error
}

func (t *GormTransaction) Rollback() error {
	return t.tx.Rollback().Error
}

func (t *GormTransaction) DB() *gorm.DB {
	return t.tx
}
