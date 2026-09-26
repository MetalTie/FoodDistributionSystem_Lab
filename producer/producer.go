// The software suite consists of main.go,
// food, stats, gatekeeper,
// producer, circularque, and sales packages.
//
// This package simulates the arrival of food packets for the GateKeeper storage facility.
// FoodPacks are removed from interplanetary shipping crates by the GateKeeper
// and repacked into smaller FoodPacks more appropriate for handling at their final destination
// planet.
//
// This goroutine will be discarded in the final implementation as real cargo vessels will contact the
// GateKeeper directly to take their loads.

package producer

import (
    "fmt"
    "time"
    
    "FoodDistributionSystem_Lab/food"
    "FoodDistributionSystem_Lab/gatekeeper"
    "FoodDistributionSystem_Lab/stats"
)

type ProductGenerator struct {
    stats       *stats.Stats
    gatekeeper  *gatekeeper.GateKeeper
    id          int
}

func NewProductGenerator(id int, gk *gatekeeper.GateKeeper, st *stats.Stats) *ProductGenerator {
    return &ProductGenerator{
        stats:      st,
        gatekeeper: gk,
        id:         id,
    }
}

func (pg *ProductGenerator) Start() {
    go pg.run()
}

func (pg *ProductGenerator) run() {
    for {
        foodType := food.RandomFoodType()
        var shipment byte = 'M'
        if food.IsGrainVegetable(foodType) {
            shipment = 'B'
        }

        newFood := food.NewFoodPack(foodType, shipment)
        
        // Simulate preparation time
        if food.IsGrainVegetable(foodType) {
            time.Sleep(pg.stats.PrepareGrainVegetableFoodPackForSales())
        } else {
            time.Sleep(pg.stats.PrepareMeatFoodPackForSales())
        }

        fmt.Printf("%c delivered.\n", shipment)
        
        // Send to gatekeeper via channel
        pg.gatekeeper.AcceptMessage(newFood)
        
        // Schedule next arrival - exponentially distributed
        nextArrivalDelay := pg.stats.NextExponential() * 1.534
        fmt.Printf("Next grain shipment arrives %.5f Time units!\n", nextArrivalDelay)
        
        time.Sleep(time.Duration(nextArrivalDelay * float64(time.Second))) // 1 second = 1 hour simulation
    }
}