package application

import (
	"time"

	"github.com/initialed85/cameranator/pkg/persistence/graphql"
	"github.com/initialed85/cameranator/pkg/persistence/model"
	"github.com/initialed85/cameranator/pkg/persistence/registry"
)

type Application struct {
	registry *registry.Registry
	client   *graphql.Client
}

func NewApplication(url string, timeout time.Duration) (*Application, error) {
	var err error

	r := registry.NewRegistry()

	err = r.Register(
		registry.NewModel("camera", model.Camera{}),
	)
	if err != nil {
		return nil, err
	}

	err = r.Register(
		registry.NewModel("video", model.Video{}),
	)
	if err != nil {
		return nil, err
	}

	err = r.Register(
		registry.NewModel("image", model.Image{}),
	)
	if err != nil {
		return nil, err
	}

	err = r.Register(
		registry.NewModel("event", model.Event{}),
	)
	if err != nil {
		return nil, err
	}

	a := Application{
		registry: r,
		client:   graphql.NewClient(url, timeout),
	}

	return &a, nil
}

func (a *Application) GetModelAndClient(name string) (*registry.ModelAndClient, error) {
	return a.registry.GetModelAndClient(name, a.client)
}
