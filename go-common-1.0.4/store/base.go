package store

import (
	"context"
	"go/ast"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var _ IBaseStore[any] = (*Store[any])(nil)

type Store[T any] struct {
	imlTransaction
	//IDB        `autowired:""`
	UniqueList []string
	Model      *T
}

func (b *Store[T]) CountByGroup(ctx context.Context, keyword string, wm map[string]interface{}, groupBy string) (map[string]int64, error) {
	db := b.DB(ctx)

	if keyword != "" {
		ids := make([]interface{}, 0)
		err := db.Model(b.Model).Select("DISTINCT target").Where("label like ?", "%"+keyword+"%").Scan(&ids).Error
		if err != nil {
			return nil, err
		}
		wm["id"] = ids
	}

	rows, err := db.Model(b.Model).Select([]string{groupBy, "count(*)"}).Where(wm).Group(groupBy).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rs := map[string]int64{}
	for rows.Next() {
		var key string
		var count int64
		err = rows.Scan(&key, &count)
		if err != nil {
			return nil, err
		}
		rs[key] = count
	}
	return rs, err
}

func (b *Store[T]) UniqueIndex() {
	b.Model = new(T)
	modelType := reflect.TypeOf(new(T)).Elem()
	for i := 0; i < modelType.NumField(); i++ {
		if fieldStruct := modelType.Field(i); ast.IsExported(fieldStruct.Name) {
			tagSetting := schema.ParseTagSetting(fieldStruct.Tag.Get("gorm"), ";")
			if _, ok := tagSetting["UNIQUEINDEX"]; ok {
				b.UniqueList = append(b.UniqueList, tagSetting["COLUMN"])
			}
		}
	}
}

func (b *Store[T]) OnComplete() {
	ctx := context.Background()
	err := b.DB(ctx).AutoMigrate(new(T))
	if err != nil {
		panic(err)
	}
	b.UniqueIndex()
}

func (b *Store[T]) Get(ctx context.Context, id int64) (*T, error) {
	value := new(T)
	err := b.DB(ctx).First(value, id).Error
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (b *Store[T]) GetByUUID(ctx context.Context, uuid string) (*T, error) {
	value := new(T)
	err := b.DB(ctx).Model(b.Model).Where("uuid = ?", uuid).First(value).Error
	if err != nil {
		return nil, err
	}
	return value, nil
}
func (b *Store[T]) Insert(ctx context.Context, t ...*T) error {
	return b.DB(ctx).Create(t).Error
}

func (b *Store[T]) Save(ctx context.Context, t *T) error {

	var v interface{} = t
	if table, ok := v.(Table); ok {

		if table.IdValue() != 0 {
			return b.DB(ctx).Save(t).Error
		}
		//没查到主键ID的数据 看看有没有唯一索引 有唯一索引 用唯一索引更新所有字段
		if len(b.UniqueList) > 0 {
			return b.UpdateByUnique(ctx, t, b.UniqueList)
		}
	}
	return b.Insert(ctx, t)
}

func (b *Store[T]) UpdateByUnique(ctx context.Context, t *T, uniques []string) error {
	columns := make([]clause.Column, 0, len(uniques))
	for _, unique := range uniques {
		columns = append(columns, clause.Column{
			Name: unique,
		})
	}
	return b.DB(ctx).Clauses(clause.OnConflict{
		Columns:   columns,
		UpdateAll: true,
	}).Create(t).Error
}

func (b *Store[T]) Delete(ctx context.Context, id ...int64) (int, error) {
	if len(id) == 0 {
		return 0, nil
	}
	result := b.DB(ctx).Delete(b.Model, id)

	return int(result.RowsAffected), result.Error
}
func (b *Store[T]) DeleteUUID(ctx context.Context, uuid string) error {
	db := b.DB(ctx)
	return db.Model(b.Model).Where("uuid = ?", uuid).Delete(b.Model).Error
}
func (b *Store[M]) SoftDelete(ctx context.Context, m map[string]interface{}) error {
	db := b.DB(ctx)
	return db.Model(b.Model).Where(m).Update("is_delete", true).Error

}
func (b *Store[M]) SoftDeleteQuery(ctx context.Context, sql string, args ...interface{}) error {
	db := b.DB(ctx)
	return db.Model(b.Model).Where(sql, args).Update("is_delete", true).Error

}

func (b *Store[T]) UpdateWhere(ctx context.Context, w map[string]interface{}, m map[string]interface{}) (int64, error) {
	var t T
	result := b.DB(ctx).Model(&t).Where(w).Updates(m)

	return result.RowsAffected, result.Error
}

func (b *Store[T]) Update(ctx context.Context, t *T) (int, error) {

	result := b.DB(ctx).Updates(t)

	return int(result.RowsAffected), result.Error
}

func (b *Store[T]) UpdateField(ctx context.Context, field string, value interface{}, sql string, args ...interface{}) (int64, error) {

	result := b.DB(ctx).Model(b.Model).Where(sql, args...).UpdateColumn(field, value)

	return result.RowsAffected, result.Error
}

func (b *Store[T]) DeleteWhere(ctx context.Context, m map[string]interface{}) (int64, error) {
	if len(m) == 0 {
		return 0, gorm.ErrMissingWhereClause
	}
	result := b.DB(ctx).Where(m).Delete(b.Model)

	return result.RowsAffected, result.Error
}
func (b *Store[T]) DeleteQuery(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	result := b.DB(ctx).Model(b.Model).Delete(sql, args...)
	return result.RowsAffected, result.Error
}

func (b *Store[T]) CountWhere(ctx context.Context, m map[string]interface{}) (int64, error) {
	if len(m) == 0 {
		return 0, gorm.ErrMissingWhereClause
	}
	var count int64
	err := b.DB(ctx).Model(b.Model).Where(m).Count(&count).Error

	return count, err
}
func (b *Store[T]) CountQuery(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	var count int64
	err := b.DB(ctx).Model(b.Model).Where(sql, args...).Count(&count).Error
	return count, err
}

func (b *Store[T]) List(ctx context.Context, m map[string]interface{}, order ...string) ([]*T, error) {
	list := make([]*T, 0)
	db := b.DB(ctx).Where(m)
	for _, o := range order {
		db = db.Order(o)
	}
	err := db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
func (b *Store[T]) ListQuery(ctx context.Context, where string, args []interface{}, order string) ([]*T, error) {
	list := make([]*T, 0)
	db := b.DB(ctx)
	db = db.Where(where, args...)
	if order != "" {
		db = db.Order(order)
	}
	err := db.Find(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (b *Store[T]) First(ctx context.Context, m map[string]interface{}, order ...string) (*T, error) {
	value := new(T)
	db := b.DB(ctx)

	err := db.Where(m).First(value).Order(order).Error
	if err != nil {
		return nil, err
	}

	return value, nil
}
func (b *Store[T]) FirstQuery(ctx context.Context, where string, args []interface{}, order string) (*T, error) {
	value := new(T)
	db := b.DB(ctx)
	if order != "" {
		db = db.Order(order)
	}
	err := db.Where(where, args...).Take(value).Error
	if err != nil {
		return nil, err
	}

	return value, nil
}
func (b *Store[T]) ListPage(ctx context.Context, where string, pageNum, pageSize int, args []interface{}, order string) ([]*T, int64, error) {
	list := make([]*T, 0, pageSize)
	db := b.DB(ctx).Where(where, args...)
	if order != "" {
		db = db.Order(order)
	}
	count := int64(0)
	err := db.Model(list).Count(&count).Limit(pageSize).Offset(PageIndex(pageNum, pageSize)).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, count, nil
}

//
//// Transaction 执行事务
//func (b *Store[T]) Transaction(ctx context.Context, f func(context.Context) error) error {
//	if b.IsTxCtx(ctx) {
//		return f(ctx)
//	}
//
//	return b.DB(ctx).Transaction(func(tx *gorm.DB) error {
//		txCtx := context.WithValue(ctx, TxContextKey, tx)
//		return f(txCtx)
//	})
//}
