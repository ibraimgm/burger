package rabbitmq

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ibraimgm/burger/app"
)

/*
This represents an hypothetical scenario where the 'kitchen' logic
could be implemented using RabbitMQ as backend technology.

Of course, since this is a demo, we don't actually integrate it with
RabbitMQ; the idea here is to demonstrate package design, not actual
implementation.
*/

type pendingOrder struct {
	doneCh chan struct{}
	order  app.Order
}

type kitchenService struct {
	pending chan pendingOrder
	supply  app.SupplyService
}

func NewKitchenService(supply app.SupplyService) app.KitchenService {
	return &kitchenService{supply: supply}
}

func (k *kitchenService) Prepare(order app.Order) chan struct{} {
	if k.pending == nil {
		k.pending = make(chan pendingOrder)

		go k.executeOrders()
	}

	doneCh := make(chan struct{})

	go func() {
		fmt.Printf("[KITCHEN] Queued order #%d (%s)\n", order.OrderNo, order.Recipe.Name)
		k.pending <- pendingOrder{order: order, doneCh: doneCh}
	}()

	return doneCh
}

// infinite loop that reads
func (k *kitchenService) executeOrders() {
	for p := range k.pending {
		fmt.Printf("[KITCHEN] Checking order #%d (%s)\n", p.order.OrderNo, p.order.Recipe.Name)

		for !k.supply.Consume(p.order.Recipe.Ingredients) {
			fmt.Printf("[KITCHEN] Not enough ingredients to make a %s! Buying more...\n", p.order.Recipe.Name)
			k.supply.Buy(p.order.Recipe.Ingredients)
		}

		fmt.Printf("[KITCHEN] Cooking order #%d (%s)\n", p.order.OrderNo, p.order.Recipe.Name)
		t := rand.Intn(len(p.order.Recipe.Ingredients)) + 2
		time.Sleep(time.Duration(t) * time.Second)

		fmt.Printf("[KITCHEN] Order #%d (%s) for table #%d is done!\n", p.order.OrderNo, p.order.Recipe.Name, p.order.TableNo)
		close(p.doneCh)
	}
}
