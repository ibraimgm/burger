package main

import (
	"fmt"
	"github.com/ibraimgm/burger/app"
	"math/rand"
	"time"
)

/*
This represents an implementation done entirely by us (i. e., "our code").
*/

type waiterService struct {
	host    app.HostService
	kitchen app.KitchenService
}

func (w *waiterService) Serve(order app.Order) {
	fmt.Printf("[WAITER] Sending the order #%d, of table #%d to the kitchen.\n", order.OrderNo, order.TableNo)
	done := w.kitchen.Prepare(order)

	go func() {
		<-done
		fmt.Printf("[WAITER] The order #%d of table #%d is delivered! Bon appetit!\n", order.OrderNo, order.TableNo)
		t := rand.Intn(10) + 1
		time.Sleep(time.Duration(t) * time.Second)

		fmt.Printf("[WAITER] Client of table #%d paid the bill and left with a smile.\n", order.TableNo)
		w.host.Free(order.TableNo)
	}()
}
