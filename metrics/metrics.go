package metrics

import (
	"github.com/ibraimgm/burger/app"
	"sync"
)

/*
This is a wrapper around the host service, to be able to collect
statistics about the restaurant operation.
*/

type Collector struct {
	Host app.HostService

	recipeCount   map[string]uint
	servedByTable map[uint]uint
	maxLineSize   int
	mu            sync.RWMutex
}

func NewCollector(host app.HostService) *Collector {
	return &Collector{
		Host: host,

		recipeCount:   make(map[string]uint),
		servedByTable: make(map[uint]uint),
	}
}

// metrics getters
func (c *Collector) RecipeCount() map[string]uint {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.recipeCount
}

func (c *Collector) ServedByTable() map[uint]uint {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.servedByTable
}

func (c *Collector) MaxLineSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.maxLineSize
}

// host methods - used to collect the data
func (c *Collector) Reserve(recipe app.Recipe) bool {
	ok := c.Host.Reserve(recipe)

	go func() {
		c.mu.Lock()
		if !ok {
			c.maxLineSize = len(c.Host.Queue())
		}

		c.recipeCount[recipe.Name]++
		c.mu.Unlock()
	}()

	return ok
}

func (c *Collector) Free(tableNo uint) {
	c.Host.Free(tableNo)

	go func() {
		c.mu.Lock()
		c.servedByTable[tableNo]++
		c.mu.Unlock()
	}()
}

func (c *Collector) Queue() []app.Recipe {
	return c.Host.Queue()
}
