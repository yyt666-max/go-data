package auto

import (
	"context"
	"fmt"

	"github.com/eolinker/eosc"
)

var (
	services = eosc.BuildUntyped[string, CompleteService]()
)

func GetService(name string) (CompleteService, bool) {
	return services.Get(name)
}
func RegisterService(name string, service CompleteService) {
	services.Set(name, service)
}
func CompleteLabels[T any](ctx context.Context, vs ...T) {
	handlers := createLabelHandler(vs)
	if handlers == nil {
		return
	}
	for name, h := range handlers {
		s, has := services.Get(name)
		if has {
			labels := s.GetLabels(ctx, h.UUIDS()...)
			if labels == nil {
				labels = make(map[string]string)
			}
			h.Set(labels)
		} else {
			h.Set(make(map[string]string))
		}
	}
}

func SearchIDCheck[T any](ctx context.Context, vs ...T) error {
	checkMap := searchIDCheck(vs)
	if checkMap == nil {
		return nil
	}
	for name, checks := range checkMap {
		for _, h := range checks {
			if len(h.UUIDS()) == 0 {
				continue
			}
			s, has := services.Get(name)
			if has {
				labels := s.GetLabels(ctx, h.UUIDS()...)
				if labels == nil || len(labels) == 0 {
					return fmt.Errorf("%s(%v) not found", h.Name(), h.UUIDS())
				}
				if len(labels) != len(h.UUIDS()) {
					notFoundIds := make([]string, 0)
					for _, id := range h.UUIDS() {
						if _, has := labels[id]; !has {
							notFoundIds = append(notFoundIds, id)
						}
					}
					return fmt.Errorf("%s(%v) not found", h.Name(), notFoundIds)
				}
			}
		}
	}
	return nil
}

type LabelHandler interface {
	UUIDS() []string
	Set(labels map[string]string)
}
type CompleteService interface {
	GetLabels(ctx context.Context, ids ...string) map[string]string
}
