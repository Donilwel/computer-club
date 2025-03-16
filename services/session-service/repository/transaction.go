package repository

import (
	"computer-club/pkg/errors"
	"gorm.io/gorm"
)

type Transaction interface {
	Commit() error
	Rollback() error
	DB() *gorm.DB
}

type GormTransaction struct {
	tx *gorm.DB
}

func (t *GormTransaction) Commit() error {
	if err := t.tx.Commit().Error; err != nil {
		return errors.ErrCommitData
	}
	return nil
}

func (t *GormTransaction) Rollback() error {
	return t.tx.Rollback().Error
}

func (t *GormTransaction) DB() *gorm.DB {
	return t.tx
}
