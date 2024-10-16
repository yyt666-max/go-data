package store

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Table interface {
	schema.Tabler
	IdValue() int64
}

type IDB interface {
	DB(ctx context.Context) *gorm.DB
	IsTxCtx(ctx context.Context) bool
}

type IBaseStore[T any] interface {
	Get(ctx context.Context, id int64) (*T, error)
	GetByUUID(ctx context.Context, uuid string) (*T, error)
	Save(ctx context.Context, t *T) error
	UpdateByUnique(ctx context.Context, t *T, uniques []string) error
	Delete(ctx context.Context, id ...int64) (int, error)
	UpdateWhere(ctx context.Context, w map[string]interface{}, m map[string]interface{}) (int64, error)
	Update(ctx context.Context, t *T) (int, error)
	UpdateField(ctx context.Context, field string, value interface{}, sql string, args ...interface{}) (int64, error)
	DeleteWhere(ctx context.Context, m map[string]interface{}) (int64, error)
	DeleteUUID(ctx context.Context, uuid string) error
	DeleteQuery(ctx context.Context, sql string, args ...interface{}) (int64, error)
	CountWhere(ctx context.Context, m map[string]interface{}) (int64, error)
	CountQuery(ctx context.Context, sql string, args ...interface{}) (int64, error)
	CountByGroup(ctx context.Context, keyword string, m map[string]interface{}, group string) (map[string]int64, error)
	SoftDelete(ctx context.Context, where map[string]interface{}) error
	SoftDeleteQuery(ctx context.Context, sql string, args ...interface{}) error
	Insert(ctx context.Context, t ...*T) error
	List(ctx context.Context, m map[string]interface{}, order ...string) ([]*T, error)
	ListQuery(ctx context.Context, sql string, args []interface{}, order string) ([]*T, error)
	First(ctx context.Context, m map[string]interface{}, order ...string) (*T, error)
	FirstQuery(ctx context.Context, sql string, args []interface{}, order string) (*T, error)
	ListPage(ctx context.Context, sql string, pageNum, pageSize int, args []interface{}, order string) ([]*T, int64, error)
	ITransaction
}
