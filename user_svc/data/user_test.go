package data

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/zycgary/mxshop-go/user_svc/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUserRepo_ListUser(t *testing.T) {
	type args struct {
		page     int32
		pageSize int32
	}

	type expects func(sqlmock.Sqlmock)

	type asserts func(*testing.T, int64, []*model.User, error)

	tests := []struct {
		name    string
		args    args
		expects expects
		asserts asserts
	}{
		{
			name: "get user list",
			args: args{
				page:     1,
				pageSize: 10,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`deleted_at` IS NULL")).
					WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(20))
				rows := sqlmock.NewRows([]string{"id", "nickname", "password", "mobile"})
				for i := 1; i <= 10; i++ {
					rows.AddRow(i, "test", "test", "123")
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL LIMIT ?")).
					WithArgs(10).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, total int64, users []*model.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(20), total)
				assert.Len(t, users, 10)
			},
		},
		{
			name: "get user list with page 2",
			args: args{
				page:     2,
				pageSize: 10,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`deleted_at` IS NULL")).
					WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(20))
				rows := sqlmock.NewRows([]string{"id", "nickname", "password", "mobile"})
				for i := 1; i <= 10; i++ {
					rows.AddRow(i, "test", "test", "123")
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL LIMIT ? OFFSET ?")).
					WithArgs(10, 10).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, total int64, users []*model.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(20), total)
				assert.Len(t, users, 10)
			},
		},
		{
			name: "empty user list",
			args: args{
				page:     1,
				pageSize: 10,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`deleted_at` IS NULL")).
					WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))
			},
			asserts: func(t *testing.T, total int64, users []*model.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, int64(0), total)
				assert.Empty(t, users)
			},
		},
		{
			name: "internal error when count",
			args: args{
				page:     1,
				pageSize: 10,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`deleted_at` IS NULL")).
					WillReturnError(errors.New("internal error"))
			},
			asserts: func(t *testing.T, total int64, users []*model.User, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), total)
				assert.Nil(t, users)
			},
		},
		{
			name: "internal error when get user list",
			args: args{
				page:     1,
				pageSize: 10,
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE `users`.`deleted_at` IS NULL")).
					WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(20))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL LIMIT ?")).
					WithArgs(10).
					WillReturnError(errors.New("internal error"))
			},
			asserts: func(t *testing.T, total int64, users []*model.User, err error) {
				assert.Error(t, err)
				assert.Equal(t, int64(0), total)
				assert.Nil(t, users)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("new sqlmock: %s", err)
			}
			defer func(db *sql.DB) {
				_ = db.Close()
			}(db)

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				SkipInitializeWithVersion: true,
				Conn:                      db,
			}))
			if err != nil {
				t.Fatalf("open gorm: %s", err)
			}

			repo := NewUserRepo(gormDB)

			tt.expects(mock)

			total, users, err := repo.ListUser(tt.args.page, tt.args.pageSize)

			tt.asserts(t, total, users, err)
		})
	}
}

func TestUserRepo_GetUserByMobile(t *testing.T) {
	type args struct {
		mobile string
	}

	type expects func(sqlmock.Sqlmock)

	type asserts func(*testing.T, *model.User, error)

	tests := []struct {
		name    string
		args    args
		expects expects
		asserts asserts
	}{
		{
			name: "Get user by mobile",
			args: args{
				mobile: "123",
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE mobile = ? AND `users`.`deleted_at` IS NULL LIMIT ?")).
					WithArgs("123", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "nickname", "password", "mobile"}).AddRow(1, "test", "test", "123"))
			},
			asserts: func(t *testing.T, user *model.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "123", user.Mobile)
			},
		},
		{
			name: "User not found",
			args: args{
				mobile: "123",
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE mobile = ? AND `users`.`deleted_at` IS NULL LIMIT ?")).
					WithArgs("123", 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "nickname", "password", "mobile"}))
			},
			asserts: func(t *testing.T, user *model.User, err error) {
				assert.Error(t, err)
				assert.Equal(t, "user not found", err.Error())
				assert.Nil(t, user)
			},
		},
		{
			name: "Internal error",
			args: args{
				mobile: "123",
			},
			expects: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE mobile = ? AND `users`.`deleted_at` IS NULL LIMIT ?")).
					WithArgs("123", 1).
					WillReturnError(errors.New("internal error"))
			},
			asserts: func(t *testing.T, user *model.User, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("new sqlmock: %s", err)
			}
			defer func(db *sql.DB) {
				_ = db.Close()
			}(db)

			gormDB, err := gorm.Open(mysql.New(mysql.Config{
				SkipInitializeWithVersion: true,
				Conn:                      db,
			}))
			if err != nil {
				t.Fatalf("open gorm: %s", err)
			}

			repo := NewUserRepo(gormDB)

			tt.expects(mock)

			user, err := repo.GetUserByMobile("123")

			tt.asserts(t, user, err)
		})
	}
}
