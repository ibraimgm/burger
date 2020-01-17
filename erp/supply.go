package erp

import (
	"fmt"
	"github.com/ibraimgm/burger/app"
	"sync"
	"time"
)

/*
This represents an hypothetical scenario where the 'supply' logic
could be implemented by integrating with the existing ERP system of the company.

Of course, since this is a demo, we don't actually integrate it with
anything, since it is made up; the idea here is to demonstrate package design,
not actual implementation.
*/

type supply struct {
	stock map[app.Ingredient]uint
	mu    sync.Mutex
}

func NewSupply(initialStock map[app.Ingredient]uint) app.SupplyService {
	return &supply{
		stock: initialStock,
	}
}

func (s *supply) Consume(ingredients []app.Ingredient) bool {
	fmt.Printf("[SUPPLY] Checking if the needed itens are available\n")

	s.mu.Lock()
	defer s.mu.Unlock()

	// we either consume eveything we need, or consume nothing
	updated := make(map[app.Ingredient]uint)

	for _, need := range ingredients {
		qty, ok := updated[need]
		if !ok {
			qty = s.stock[need]
		}

		if qty == 0 {
			fmt.Printf("[SUPPLY] Not enough items in stock\n")
			return false
		}

		updated[need] = qty - 1
	}

	// no error? apply the changes!
	for k, v := range updated {
		s.stock[k] = v
	}

	fmt.Printf("[SUPPLY] Requested itens delivered to the kitchen\n")
	return true
}

func (s *supply) Buy(ingredients []app.Ingredient) {
	time.Sleep(3 * time.Second)
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, need := range ingredients {
		if s.stock[need] >= 5 {
			continue
		}

		s.stock[need] += 10
		fmt.Printf("[SUPPLY] Bougth '%s'\n", need)
	}
}
