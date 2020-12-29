package registry

import (
	"github.com/initialed85/cameranator/pkg/persistence/graphql"
)

type Model struct {
	name      string
	reference interface{}
}

func NewModel(
	name string,
	reference interface{},
) *Model {
	m := Model{
		name:      name,
		reference: reference,
	}

	return &m
}

func (m *Model) GetAll(
	c *graphql.Client,
	item interface{},
) error {
	query, err := graphql.GetManyQuery(
		m.name,
		m.reference,
		"id",  // TODO: tied to database schema
		"asc", // TODO: tied to database schema
	)
	if err != nil {
		return err
	}

	err = c.QueryAndExtract(query, m.name, &item)
	if err != nil {
		return err
	}

	return err
}

func (m *Model) GetOne(
	c *graphql.Client,
	item interface{},
	conditionKey string,
	conditionValue interface{},
) error {
	query, err := graphql.GetOneQuery(
		m.name,
		m.reference,
		conditionKey,
		conditionValue,
	)
	if err != nil {
		return err
	}

	err = c.QueryAndExtract(query, m.name, &item)
	if err != nil {
		return err
	}

	return err
}

func (m *Model) Add(
	c *graphql.Client,
	item interface{},
) error {
	query, err := graphql.InsertQuery(
		m.name,
		item,
	)
	if err != nil {
		return err
	}

	err = c.QueryAndExtract(query, m.name, &item)
	if err != nil {
		return err
	}

	return err
}

func (m *Model) Remove(
	c *graphql.Client,
	item interface{},
) error {
	query, err := graphql.DeleteQuery(
		m.name,
		item,
	)
	if err != nil {
		return err
	}

	err = c.QueryAndExtract(query, m.name, &item)
	if err != nil {
		return err
	}

	return err
}

type ModelAndClient struct {
	model  *Model
	client *graphql.Client
}

func NewModelAndClient(model *Model, client *graphql.Client) *ModelAndClient {
	m := ModelAndClient{
		model,
		client,
	}

	return &m
}

func (m *ModelAndClient) GetAll(
	item interface{},
) error {
	return m.model.GetAll(m.client, &item)
}

func (m *ModelAndClient) GetOne(
	item interface{},
	conditionKey string,
	conditionValue interface{},
) error {
	return m.model.GetOne(m.client, &item, conditionKey, conditionValue)
}

func (m *ModelAndClient) Add(
	item interface{},
) error {
	return m.model.Add(m.client, item)
}

func (m *ModelAndClient) Remove(
	item interface{},
) error {
	return m.model.Remove(m.client, item)
}
