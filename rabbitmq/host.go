package rabbitmq

/*
This represents an hypothetical scenario where the 'host' logic
could be implemented using RabbitMQ as backend technology.

Of course, since this is a demo, we don't actually integrate it with
RabbitMQ; the idea here is to demonstrate package design, not actual
implementation.
*/

import (
	"fmt"
	"github.com/ibraimgm/burger/app"
	"sync"
)

type hostService struct {
	waiter app.WaiterService

	orderNo uint
	tables  map[uint]bool
	queued  []app.Recipe
	mu      sync.Mutex
}

func NewHostService(availableTables uint, waiter app.WaiterService) app.HostService {
	tables := make(map[uint]bool)

	for i := 1; i <= 5; i++ {
		tables[uint(i)] = true
	}

	return &hostService{
		waiter: waiter,
		tables: tables,
	}
}

func (h *hostService) Reserve(recipe app.Recipe) bool {
	fmt.Printf("[HOST] Received a new client asking for a '%s'\n", recipe.Name)

	h.mu.Lock()
	defer h.mu.Unlock()

	for k, v := range h.tables {
		if v {
			fmt.Printf("[HOST] Table #%d is available!\n", k)
			h.placeOrder(k, recipe)
			return true
		}
	}

	fmt.Printf("[HOST] No available tables; you will have to wait for your '%s', sir!\n", recipe.Name)
	h.queued = append(h.queued, recipe)
	return false
}

func (h *hostService) Free(tableNo uint) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.tables[tableNo] = true
	fmt.Printf("[HOST] Table #%d is now available.\n", tableNo)

	if len(h.queued) == 0 {
		return
	}

	next := h.queued[0]
	h.queued = h.queued[1:]
	fmt.Printf("[HOST] The '%s' is the first in line.\n", next.Name)
	h.placeOrder(tableNo, next)
}

func (h *hostService) Queue() []app.Recipe {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.queued
}

// this assumes the caller got a lock (yeah, right...)
func (h *hostService) placeOrder(tableNo uint, recipe app.Recipe) {
	h.tables[tableNo] = false
	h.orderNo++

	fmt.Printf("[HOST] Generating order #%d for table #%d\n", h.orderNo, tableNo)
	order := app.Order{
		OrderNo: h.orderNo,
		TableNo: tableNo,
		Recipe:  recipe,
	}

	fmt.Printf("[HOST] Order sent to the waiter\n")
	h.waiter.Serve(order)
}
