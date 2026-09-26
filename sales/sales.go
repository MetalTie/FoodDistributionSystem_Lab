// The software suite consists of main.go,
// food, stats, gatekeeper,
// producer, circularque, and sales packages.
//
// The specification require the ability to create multiple points of sale.

package sales

import (
    "fmt"
    "time"
    
    //"FoodDistributionSystem_Lab/food"
    "FoodDistributionSystem_Lab/gatekeeper"
    "FoodDistributionSystem_Lab/stats"
)

type RetailSales struct {
    stats      *stats.Stats
    gatekeeper *gatekeeper.GateKeeper
    id         int
}

func NewRetailSales(id int, gk *gatekeeper.GateKeeper, st *stats.Stats) *RetailSales {
    return &RetailSales{
        stats:      st,
        gatekeeper: gk,
        id:         id,
    }
}

func (rs *RetailSales) Start() {
    go rs.run()
}

func (rs *RetailSales) run() {
    time.Sleep(1 * time.Second) // Allow for initialization activities (1 second = 1 hour simulation)
    
    for {
        foodItem, available := rs.gatekeeper.RetrieveMessage()
        if available {
            // The time to sell a product is exponentially distributed with mean 2.0 hours
            salesTime := rs.stats.NextExponential() * 2.0
            time.Sleep(time.Duration(salesTime * float64(time.Second)))
            
            fmt.Printf("Retail Sales successfully sold %s %c\n", foodItem.FoodType, foodItem.FoodShipment)
        } else {
            time.Sleep(100 * time.Millisecond) // Brief pause if no food available
        }
    }
}