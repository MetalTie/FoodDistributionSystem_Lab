// The software suite consists of main.go,
// food, stats, gatekeeper,
// producer, storage, and sales packages.
//
// This represents the software to manage an "embedded" planetary system
// food receiving and distribution system.

package food

import (
    "math/rand"
    "time"
)

type FoodType int

const (
    Wheat FoodType = iota
    Beans
    Corn
    Rice
    Potatoes
    Squash
    Tomato
    Steak
    Pork
    Fish
    Fowel
)

func (ft FoodType) String() string {
    return []string{"Wheat", "Beans", "Corn", "Rice", "Potatoes", "Squash", "Tomato", "Steak", "Pork", "Fish", "Fowel"}[ft]
}

type FoodPack struct {
    FoodType     FoodType
    FoodShipment byte
}

func NewFoodPack(foodType FoodType, shipment byte) FoodPack {
    return FoodPack{
        FoodType:     foodType,
        FoodShipment: shipment,
    }
}

func init() {
    rand.Seed(time.Now().UnixNano())
}

func RandomFoodType() FoodType {
    return FoodType(rand.Intn(int(Fowel) + 1))
}

func IsGrainVegetable(ft FoodType) bool {
    return ft >= Wheat && ft <= Tomato
}