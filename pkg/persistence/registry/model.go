package registry

import (
	"fmt"

	"github.com/initialed85/cameranator/pkg/persistence/graphql"
	"github.com/initialed85/cameranator/pkg/utils"
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
		"",
		nil,
		"id",  // TODO: tied to database schema
		"asc", // TODO: tied to database schema
	)
	if err != nil {
		return fmt.Errorf("failed to invoke GetManyQuery: %v", err)
	}

	err = c.QueryAndExtract(query, m.name, &item)
	if err != nil {
		return fmt.Errorf("failed to invoke QueryAndExtract: %v", err)
	}

	return nil
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
		return fmt.Errorf("failed to invoke GetOneQuery: %v", err)
	}

	err = c.QueryAndExtract(query, m.name, &item)
	if err != nil {
		return fmt.Errorf("failed to invoke QueryAndExtract: %v", err)
	}

	return nil
}

func (m *Model) GetMany(
	c *graphql.Client,
	item interface{},
	conditionKey string,
	conditionValue interface{},
) error {
	query, err := graphql.GetManyQuery(
		m.name,
		m.reference,
		conditionKey,
		conditionValue,
		"id",  // TODO: tied to database schema
		"asc", // TODO: tied to database schema
	)
	if err != nil {
		return fmt.Errorf("failed to invoke GetManyQuery: %v", err)
	}

	err = c.QueryAndExtract(query, m.name, &item)
	if err != nil {
		return fmt.Errorf("failed to invoke QueryAndExtract: %v", err)
	}

	return nil
}

func (m *Model) Add(
	c *graphql.Client,
	item interface{},
	items interface{},
) error {
	query, err := graphql.InsertQuery(
		m.name,
		utils.Dereference(item),
	)
	if err != nil {
		return fmt.Errorf("failed to invoke InsertQuery: %v", err)
	}

	err = c.QueryAndExtract(
		query,
		fmt.Sprintf("insert_%v_one", m.name),
		&items,
	)
	if err != nil {
		return fmt.Errorf("failed to invoke QueryAndExtract: %v", err)
	}

	return nil
}

func (m *Model) Remove(
	c *graphql.Client,
	item interface{},
	items interface{},
) error {
	query, err := graphql.DeleteQuery(
		m.name,
		utils.Dereference(item),
	)
	if err != nil {
		return fmt.Errorf("failed to invoke DeleteQuery: %v", err)
	}

	err = c.QueryAndExtract(
		query,
		fmt.Sprintf("delete_%v", m.name),
		&items,
	)
	if err != nil {
		return fmt.Errorf("failed to invoke QueryAndExtract: %v", err)
	}

	return nil
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

func (m *ModelAndClient) Model() *Model {
	return m.model
}

func (m *ModelAndClient) Client() *graphql.Client {
	return m.client
}

func (m *ModelAndClient) GetAll(
	items interface{},
) error {
	err := m.model.GetAll(m.client, &items)
	if err != nil {
		return fmt.Errorf("failed to invoke GetAll: %v", err)
	}

	return nil
}

func (m *ModelAndClient) GetOne(
	items interface{},
	conditionKey string,
	conditionValue interface{},
) error {
	err := m.model.GetOne(m.client, &items, conditionKey, conditionValue)
	if err != nil {
		return fmt.Errorf("failed to invoke GetOne: %v", err)
	}

	return nil
}

func (m *ModelAndClient) GetMany(
	items interface{},
	conditionKey string,
	conditionValue interface{},
) error {
	err := m.model.GetMany(m.client, &items, conditionKey, conditionValue)
	if err != nil {
		return fmt.Errorf("failed to invoke GetMany: %v", err)
	}

	return nil
}

func (m *ModelAndClient) Add(
	item interface{},
	items interface{},
) error {
	err := m.model.Add(m.client, &item, &items)
	if err != nil {
		return fmt.Errorf("failed to invoke Add: %v", err)
	}

	return nil
}

func (m *ModelAndClient) Remove(
	item interface{},
	items interface{},
) error {
	err := m.model.Remove(m.client, &item, &items)
	if err != nil {
		return fmt.Errorf("failed to invoke Remove: %v", err)
	}

	return nil
}
