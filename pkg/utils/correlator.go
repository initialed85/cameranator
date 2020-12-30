package utils

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

type Item struct {
	mu          sync.Mutex
	isComplete  bool
	correlation *Correlation
	name        string
	value       interface{}
}

func (i *Item) GetName() string {
	return i.name
}

func (i *Item) GetValue() interface{} {
	i.mu.Lock()
	defer i.mu.Unlock()

	return i.value
}

func (i *Item) SetValue(value interface{}) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.value = value
}

func (i *Item) Complete() {
	i.mu.Lock()

	if i.isComplete {
		i.mu.Unlock()
		return
	}

	i.isComplete = true

	i.mu.Unlock()

	log.Printf("item %#+v is complete", i.name)

	i.correlation.Complete()
}

type Correlation struct {
	mu            sync.Mutex
	items         []*Item
	isComplete    bool
	correlator    *Correlator
	correlationID uuid.UUID
	completeFn    func(*Correlation)
}

func (c *Correlation) NewItem(name string) *Item {
	item := Item{
		name:        name,
		isComplete:  false,
		correlation: c,
	}

	c.items = append(
		c.items,
		&item,
	)

	return &item
}

func (c *Correlation) GetCorrelationID() uuid.UUID {
	return c.correlationID
}

func (c *Correlation) GetItems() []*Item {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.items
}

func (c *Correlation) GetItem(name string) (*Item, error) {
	for _, item := range c.GetItems() {
		if item.GetName() == name {
			return item, nil
		}
	}

	return nil, fmt.Errorf("failed to find Item with name=%#+v", name)
}

func (c *Correlation) Complete() {
	c.mu.Lock()

	if c.isComplete {
		c.mu.Unlock()
		return
	}

	for _, item := range c.items {
		if !item.isComplete {
			c.mu.Unlock()
			return
		}
	}

	c.mu.Unlock()

	log.Printf("correlation %#+v is complete", c.correlationID.String())

	c.completeFn(c)

	c.correlator.Complete()
}

type Correlator struct {
	mu                         sync.Mutex
	correlationByCorrelationID map[uuid.UUID]*Correlation
}

func NewCorrelator() *Correlator {
	c := Correlator{
		correlationByCorrelationID: make(map[uuid.UUID]*Correlation),
	}

	return &c
}

func (c *Correlator) NewCorrelation(completeFn func(*Correlation)) *Correlation {
	c.mu.Lock()
	defer c.mu.Unlock()

	correlation := Correlation{
		items:         make([]*Item, 0),
		isComplete:    false,
		correlator:    c,
		correlationID: GetUUID(),
		completeFn:    completeFn,
	}

	c.correlationByCorrelationID[correlation.correlationID] = &correlation

	return &correlation
}

func (c *Correlator) Complete() {
	c.mu.Lock()
	defer c.mu.Unlock()

	toBeRemoved := make([]uuid.UUID, 0)
	for correlationID, correlation := range c.correlationByCorrelationID {
		if !correlation.isComplete {
			continue
		}

		toBeRemoved = append(toBeRemoved, correlationID)
	}

	for _, correlationID := range toBeRemoved {
		delete(c.correlationByCorrelationID, correlationID)
	}
}
