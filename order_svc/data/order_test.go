package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zycgary/mxshop-go/order_svc/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_orderRepo_ListCartItems(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId int32
	}

	type expects func(sqlmock.Sqlmock)

	type asserts func(*testing.T, []*model.CartItem, error)

	tests := []struct {
		name    string
		args    args
		expects expects
		asserts asserts
	}{
		{
			name: "get cart items",
			args: args{
				ctx:    context.Background(),
				userId: 1,
			},
			expects: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "product_name", "product_price", "product_num", "checked"})
				for i := 1; i <= 10; i++ {
					rows.AddRow(i, 1, i, fmt.Sprintf("test product %d", i), 100, 1, 1)
				}

				m.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `cart_items` WHERE user_id = ? AND `cart_items`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(rows)
			},
			asserts: func(t *testing.T, items []*model.CartItem, err error) {
				assert.NoError(t, err)
				assert.Len(t, items, 10)
			},
		},
		{
			name: "iternal error",
			args: args{
				ctx:    context.Background(),
				userId: 1,
			},
			expects: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "product_name", "product_price", "product_num", "checked"})
				for i := 1; i <= 10; i++ {
					rows.AddRow(i, 1, i, fmt.Sprintf("test product %d", i), 100, 1, 1)
				}

				m.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `cart_items` WHERE user_id = ? AND `cart_items`.`deleted_at` IS NULL")).
					WithArgs(1).
					WillReturnError(errors.New("internal error"))
			},
			asserts: func(t *testing.T, items []*model.CartItem, err error) {
				assert.Error(t, err)
				assert.Nil(t, items)
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

			repo := NewOrderRepo(gormDB)

			tt.expects(mock)

			got, err := repo.ListCartItems(tt.args.ctx, tt.args.userId)

			tt.asserts(t, got, err)
		})
	}
}
