package registry

import (
	"fmt"
	"sync"

	"github.com/initialed85/cameranator/pkg/persistence/graphql"
)

type Registry struct {
	mu          sync.Mutex
	modelByName map[string]*Model
}

func NewRegistry() *Registry {
	r := Registry{
		modelByName: make(map[string]*Model),
	}

	return &r
}

func (r *Registry) Register(model *Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.modelByName[model.name]
	if ok {
		return fmt.Errorf("model already exists for %#v", model.name)
	}

	r.modelByName[model.name] = model

	return nil
}

func (r *Registry) Unregister(model *Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.modelByName[model.name]
	if !ok {
		return fmt.Errorf("model does not exist for %#v", model.name)
	}

	delete(r.modelByName, model.name)

	return nil
}

func (r *Registry) GetModelAndClient(name string, client *graphql.Client) (*ModelAndClient, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	model, ok := r.modelByName[name]
	if !ok {
		return nil, fmt.Errorf("model does not exist for %#v", name)
	}

	return NewModelAndClient(model, client), nil
}
