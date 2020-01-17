package app

// Ingredient is a single ingredient the can be used in a recipe.
// ex: bread, burger, sausage, etc.
type Ingredient string

// A recipe is nothing more than a name and a bunch of ingredients
type Recipe struct {
	Name        string
	Ingredients []Ingredient
}

// HostService represents the host(ess) that manager which client
// go in which table. If there is no tables available, clients go to a queue.
type HostService interface {
	Reserve(Recipe) bool // true = reserved; false = in the queue
	Free(tableNo uint)
	Queue() []Recipe
}

// Every order has a number, and every customer ill ask for only one item
type Order struct {
	OrderNo uint
	TableNo uint
	Recipe  Recipe
}

// WaiterService is the entry point from a client's order.
// This service manage the available tables and will put the orders into
// a queue, if no tables are available.
type WaiterService interface {
	// serve the client's order, i.e. send it to kitchen, deliver to table, etc.
	Serve(order Order)
}

// KitchenService represents the kitchen staff.
type KitchenService interface {
	// Prepare an order. Returns a channel that notifies when the order is done.
	Prepare(order Order) chan struct{}
}

// Supply is the kitchen supply.
type SupplyService interface {
	// Consume the needed ingredients. Returns false if there is no stock
	Consume(ingredients []Ingredient) bool

	// Buy the needed ingredients, in increased quantities for a better market discount
	Buy(ingredients []Ingredient)
}
