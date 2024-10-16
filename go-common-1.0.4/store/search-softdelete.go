package store

import (
	"context"
	"fmt"
)

var _ ISearchStore[any] = (*SearchStoreSoftDelete[any])(nil)

type SearchStoreSoftDelete[T any] struct {
	SearchStore[T]
}

func (s *SearchStoreSoftDelete[T]) Delete(ctx context.Context, id ...int64) (int, error) {

	r := s.Store.DB(ctx).Where(map[string]interface{}{
		"id": id,
	}).Update("is_delete", true)
	return int(r.RowsAffected), r.Error
}

func (s *SearchStoreSoftDelete[T]) DeleteWhere(ctx context.Context, m map[string]interface{}) (int64, error) {

	r := s.Store.DB(ctx).Where(m).Update("is_delete", true)
	return r.RowsAffected, r.Error
}

func (s *SearchStoreSoftDelete[T]) DeleteUUID(ctx context.Context, uuid string) error {
	r := s.Store.DB(ctx).Where(map[string]interface{}{
		"uuid": uuid,
	}).Update("is_delete", true)
	return r.Error
}

func (s *SearchStoreSoftDelete[T]) DeleteQuery(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	r := s.Store.DB(ctx).Where(sql, args...).Update("is_delete", true)
	return r.RowsAffected, r.Error
}

func (s *SearchStoreSoftDelete[T]) CountWhere(ctx context.Context, m map[string]interface{}) (int64, error) {
	vm := m
	if vm == nil {
		vm = map[string]interface{}{}
	}
	vm["is_delete"] = false
	return s.SearchStore.CountWhere(ctx, vm)
}

func (s *SearchStoreSoftDelete[T]) CountQuery(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	if sql == "" {
		return s.SearchStore.CountQuery(ctx, "is_delete = false", false)
	}
	return s.SearchStore.CountQuery(ctx, fmt.Sprintf("(%s) AND is_delete = false", sql), args...)
}

func (s *SearchStoreSoftDelete[T]) List(ctx context.Context, m map[string]interface{}, order ...string) ([]*T, error) {
	vm := m
	if vm == nil {
		vm = map[string]interface{}{}
	}
	return s.SearchStore.List(ctx, vm, order...)
}

func (s *SearchStoreSoftDelete[T]) ListQuery(ctx context.Context, sql string, args []interface{}, order string) ([]*T, error) {
	if sql != "" {
		sql = fmt.Sprintf("(%s) AND is_delete = false", sql)
	} else {
		sql = "is_delete = false"
	}
	return s.SearchStore.ListQuery(ctx, sql, args, order)
}

func (s *SearchStoreSoftDelete[T]) First(ctx context.Context, m map[string]interface{}, order ...string) (*T, error) {
	if m == nil {
		m = map[string]interface{}{}
	}
	m["is_delete"] = false
	return s.SearchStore.First(ctx, m, order...)
}

func (s *SearchStoreSoftDelete[T]) FirstQuery(ctx context.Context, sql string, args []interface{}, order string) (*T, error) {
	if sql != "" {
		sql = fmt.Sprintf("(%s) AND is_delete = false", sql)
	} else {
		sql = "is_delete = false"
	}
	return s.SearchStore.FirstQuery(ctx, sql, args, order)
}

func (s *SearchStoreSoftDelete[T]) ListPage(ctx context.Context, sql string, pageNum, pageSize int, args []interface{}, order string) ([]*T, int64, error) {
	if sql != "" {
		sql = fmt.Sprintf("(%s) AND is_delete = false", sql)
	} else {
		sql = "is_delete = false"
	}

	return s.SearchStore.ListPage(ctx, sql, pageNum, pageSize, args, order)

}

func (s *SearchStoreSoftDelete[T]) Search(ctx context.Context, keyword string, condition map[string]interface{}, sortRule ...string) ([]*T, error) {
	if condition == nil {
		condition = map[string]interface{}{}
	}
	condition["is_delete"] = false
	return s.SearchStore.Search(ctx, keyword, condition, sortRule...)
}

func (s *SearchStoreSoftDelete[T]) Count(ctx context.Context, keyword string, condition map[string]interface{}) (int64, error) {
	if condition == nil {
		condition = map[string]interface{}{}
	}
	condition["is_delete"] = false
	return s.SearchStore.Count(ctx, keyword, condition)
}

func (s *SearchStoreSoftDelete[T]) CountByGroup(ctx context.Context, keyword string, condition map[string]interface{}, groupBy string) (map[string]int64, error) {
	if condition == nil {
		condition = map[string]interface{}{}
	}
	condition["is_delete"] = false
	return s.SearchStore.CountByGroup(ctx, keyword, condition, groupBy)
}

func (s *SearchStoreSoftDelete[T]) SearchByPage(ctx context.Context, keyword string, condition map[string]interface{}, page int, pageSize int, sortRule ...string) ([]*T, int64, error) {
	if condition == nil {
		condition = map[string]interface{}{}
	}
	condition["is_delete"] = false
	return s.SearchStore.SearchByPage(ctx, keyword, condition, page, pageSize, sortRule...)
}
