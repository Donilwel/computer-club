package repository_test

import (
	"computer-club/internal/repository"
	"computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	return gormDB, mock
}

func TestGetComputers(t *testing.T) {
	gormDB, mock := setupDB(t)
	repo := repository.NewComputerRepository(gormDB)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT \* FROM "computers"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "pc_number", "status"}).
			AddRow(1, 101, models.Free).
			AddRow(2, 102, models.Busy))

	computers, err := repo.GetComputers(ctx)
	assert.NoError(t, err)
	assert.Len(t, computers, 2)
	assert.Equal(t, models.Free, computers[0].Status)
}

func TestIsComputerAvailable(t *testing.T) {
	gormDB, mock := setupDB(t)
	repo := repository.NewComputerRepository(gormDB)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT \* FROM "computers" WHERE pc_number = \$1 ORDER BY "computers"."id" LIMIT \$2`).
		WithArgs(101, 1). // <- добавляем второй аргумент
		WillReturnRows(sqlmock.NewRows([]string{"id", "pc_number", "status"}).AddRow(1, 101, models.Busy))

	available, err := repo.IsComputerAvailable(ctx, nil, 101)
	assert.NoError(t, err)
	assert.True(t, available)
}

func TestIsComputerAvailable_NotFound(t *testing.T) {
	gormDB, mock := setupDB(t)
	repo := repository.NewComputerRepository(gormDB)
	ctx := context.Background()

	mock.ExpectQuery(`SELECT \* FROM "computers" WHERE pc_number = \$1 ORDER BY "computers"."id" LIMIT 1`).
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	available, err := repo.IsComputerAvailable(ctx, nil, 999)
	assert.ErrorIs(t, err, errors.ErrComputerNotFound)
	assert.False(t, available)
}

func TestMarkComputerFree_Error(t *testing.T) {
	gormDB, mock := setupDB(t)
	repo := repository.NewComputerRepository(gormDB)
	ctx := context.Background()

	mock.ExpectExec(`UPDATE "computers" SET status = \$1 WHERE pc_number = \$2`).
		WithArgs(models.Free, 101).
		WillReturnError(sql.ErrConnDone)

	err := repo.MarkComputerFree(ctx, nil, 101)
	assert.ErrorIs(t, err, errors.ErrUpdateComputerStatus)
}
