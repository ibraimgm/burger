# Burger Shop sample

This is a sample project to show how to organize "service-like" packages in Go.
The code is intentionally convoluted to force certain situations (for example, de bidirectional dependency between `host`and `waiter`), so we can see how the design grows and adapts to these changes.

## Understanding the problem

In this demo, we're developing a *Burger Shop* application. Our restaurant has a variety of different services to support its day-to-day operation:

* The `HostService`, that represents the host, i. e., the person that checks if there are available tables and send the clients to the waiters;
* The `WaiterService`, that sends the client order to the kitchen and, when it's done, bring it back to the client;
* The `KitchenService`, that prepares the food;
* The `SupplyService`, that controls the consumed supllies and can buy more food, if needed.

All these services are defined as `interface`s, together with the domain "models" in the `app` package. This is done to be able to solve bidirectional dependencies, like the host/waiter relationship. The concrete implementations are scattered in different packages, to represent e. g. integration with different services or libraries:

* Package `rabbitmq` has concrete implementations for `HostService` and `KitchenService`. Of course, we don't actually use RabbitMQ; this is just a ay to help you visualize how an external service might be integrated on this architecture.
* Package `erp` implements a concrete `SupplyService`. As with the previous example, imagine this as a layer to talk to an existing legacy ERP system or something like that.
* The `WaiterService` is implemented directly in the `main` package (e. g. code that does not integrate with external libraries, etc.).

We also have other auxiliary packages:

* Both `http` and `telnet` provide concrete implementation for the `app.Server`. This is to demonstrate ho easy it is to have different kinds on servers using ashared business logic.
* The `recipe` package is just a bunch of global variables with our restaurant's menu, for easier access.
* Last, but not least, `metrics` contains a metrics collector that wraps around a `HostService` and collect additional data to display in our servers.

As you can probably guess, the `main.go` file does all the setup and dependency injection.

## Running the demo

Just build the project and run the main binary (there are no command-line parameters available):

```bash
make && ./burger
```

If everything is alright, you can go to [http://localhost:8080/](http://localhost:8080/) and see the statistics collected. To actually be able to see the services talking to each other connect via telnet (`telnet localhost 5555`) and use the `help` command to see the available options. The output of what each service is doing will be displayed on the server's `stdout`, so keep an eye there.
