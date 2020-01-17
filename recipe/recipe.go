package recipe

import (
	"github.com/ibraimgm/burger/app"
)

var Burger = app.Recipe{
	Name: "burger",
	Ingredients: []app.Ingredient{
		"bread",
		"burger",
		"cheese",
		"ketchup",
		"lettuce",
		"tomato",
		"pea",
		"corn",
	},
}

var DoubleBurger = app.Recipe{
	Name: "Double Burger",
	Ingredients: []app.Ingredient{
		"bread",
		"burger",
		"burger",
		"cheese",
		"barbecue",
		"bacon",
		"cheddar",
		"onion",
	},
}

var HotDog = app.Recipe{
	Name: "HotDog",
	Ingredients: []app.Ingredient{
		"bread",
		"sausage",
		"sausage",
		"cheese",
		"ketchup",
		"mayonnaise",
		"potatoes",
		"potatoes",
		"pea",
		"corn",
	},
}
