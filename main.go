package main

import (
	"context"
	"github.com/ibraimgm/burger/http"
	"github.com/ibraimgm/burger/metrics"

	"github.com/ibraimgm/burger/app"
	"github.com/ibraimgm/burger/erp"
	"github.com/ibraimgm/burger/rabbitmq"
	"github.com/ibraimgm/burger/telnet"
)

func main() {
	// server will never 'end'
	ctx := context.Background()

	// wire up each service.
	// notice how we 'inject' the metrics collector
	supply := erp.NewSupply(make(map[app.Ingredient]uint))
	waiter := &waiterService{}
	host := metrics.NewCollector(rabbitmq.NewHostService(5, waiter))
	kitchen := rabbitmq.NewKitchenService(supply)
	waiter.host = host
	waiter.kitchen = kitchen

	// set up each server
	// remember - our 'host' variable is both a host and a metrics collector
	servers := []app.Server{
		telnet.New(":5555", host, host),
		http.New(":8080", host),
	}

	for _, server := range servers {
		server.Start(ctx)
	}

	println("Server started.")

	<-ctx.Done()
	println("Bye!")
}
