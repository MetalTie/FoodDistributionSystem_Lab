// The software suite consists of main.go,
// food, stats, gatekeeper,
// producer, circularque, and sales packages.
//
// The GateKeeper accepts FoodPacks from inter-galactic transports
// and repacks them into FoodPacks suitable for Sales to
// distribute to the planets in the system.

package gatekeeper

import (
    "fmt"
    "time"
    
    "FoodDistributionSystem_Lab/food"
    "FoodDistributionSystem_Lab/storage" // updated import
)

type GateKeeper struct {
    storage        *storage.SortedList // Changed from circularque
    acceptChan     chan food.FoodPack
    retrieveChan   chan retrieveRequest
    rejected       int
    startTime      time.Time
    endTime        time.Time
}

type retrieveRequest struct {
    response  chan food.FoodPack
    available chan bool
}

func NewGateKeeper(capacity int) *GateKeeper {
    gk := &GateKeeper{
        storage:        storage.NewSortedList(capacity), // use sorted list
        acceptChan:   make(chan food.FoodPack, 100),
        retrieveChan: make(chan retrieveRequest, 100),
        rejected:     0,
        startTime:    time.Now(),
        endTime:      time.Now().Add(40 * time.Second), // 40 seconds = 40 hours simulation
    }
    go gk.run()
    return gk
}

func (gk *GateKeeper) AcceptMessage(foodPack food.FoodPack) {
    gk.acceptChan <- foodPack
}

func (gk *GateKeeper) RetrieveMessage() (food.FoodPack, bool) {
    respChan := make(chan food.FoodPack)
    availChan := make(chan bool)
    
    gk.retrieveChan <- retrieveRequest{response: respChan, available: availChan}
    
    available := <-availChan
    if available {
        return <-respChan, true
    }
    return food.FoodPack{}, false
}

func (gk *GateKeeper) run() {
    time.Sleep(500 * time.Millisecond) // Allow initialization (0.5 seconds = 0.5 hours)
    
    for gk.rejected < 5 && time.Now().Before(gk.endTime) {
        select {
        case newFood := <-gk.acceptChan:
            if !gk.storage.IsFull() {
                gk.storage.AcceptMessage(newFood)
                fmt.Printf("GateKeeper insert accepted %s %c\n", newFood.FoodType, newFood.FoodShipment)
            } else {
                gk.rejected++
                fmt.Printf("Rejected by GateKeeper:\n")
                fmt.Printf("%s %c\n", newFood.FoodType, newFood.FoodShipment)
                fmt.Printf("Rejected = %d. Sent to another distribution facility!\n\n", gk.rejected)
            }
            
        case req := <-gk.retrieveChan:
            if !gk.storage.IsEmpty() {
                // foodItem, _ := gk.storage.RetrieveMessage()
                
                // Generate desired food type for management reporting
                mgtDesiredType := food.RandomFoodType()
                fmt.Printf("Mgt Desired Food Type To Sell is: %s\n", mgtDesiredType)
                // fmt.Printf("Actual type sold is: %s\n", foodItem.FoodType)
                fmt.Printf("Food pack removed by GateKeeper for shipment.\n\n")
                
                req.available <- true
                // req.response <- foodItem
            } else {
                req.available <- false
            }
        }
        
        time.Sleep(1100 * time.Millisecond) // Processing overhead (1.1 seconds = 1.1 hours simulation)
    }
    
    elapsed := time.Since(gk.startTime)
    fmt.Printf("\n\nHours of operation prior to closing: %.3f\n", elapsed.Seconds())
}