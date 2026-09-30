// The GateKeeper utilizes a circular queue for storage and retrieval. The package has been
// designed to be generic to allow use in future systems with respect to the data type and queue size.
//
// The software suite consists of main.go,
// food, stats, gatekeeper,
// producer, sortedlist, and sales packages.
//
// This represents the software to manage an "embedded" planetary system
// food receiving and distribution system.

package storage // Changed the name to better fit the function of this package

import (
    "errors"
    "fmt"
    "sync"

    "FoodDistributionSystem_Lab/food" // Added this to grab from the food package
)

type SortedList struct {
    items     []food.FoodPack // Created an array to hold the food items
    capacity  int // changed to hold the amount
    size      int // Changed to keep track of the size
    mutex     sync.RWMutex // changed here 
}

func NewSortedList(capacity int) *SortedList {
    return &SortedList{
        items: make([]food.FoodPack, 0, capacity),
        capacity: capacity,
        size: 0,
    }
}

// Insert maintains sorted order with meat first, then grain/vegetable
func (sl *SortedList) Insert(item food.FoodPack) error {
    sl.mutex.Lock()
    defer sl.mutex.Unlock()
    
    if sl.size >= sl.capacity {
        return errors.New("storage is full")
    }
    
    // Find the correct insertion position
    position := sl.findInsertPosition(item)
    
    // Make space for new item
    sl.items = append(sl.items, food.FoodPack{})
    copy(sl.items[position+1:], sl.items[position:])
    
    // Insert the new item
    sl.items[position] = item
    sl.size++
    
    return nil
}

// Find the correct position to maintain sorted order
func (sl *SortedList) findInsertPosition(item food.FoodPack) int {
    for i := 0; i < sl.size; i++ {
        current := sl.items[i]
        
        // Meat comes before grain/vegetable
        if item.FoodShipment == 'M' && current.FoodShipment == 'B' {
            return i
        }
        
        // Same type, sort by food type
        if item.FoodShipment == current.FoodShipment {
            if int(item.FoodType) < int(current.FoodType) {
                return i
            }
        }
    }
    
    return sl.size  // Insert at end
}

// Binary search for efficient lookup (B Option)
func (sl *SortedList) FindBinary(desiredType food.FoodType) (food.FoodPack, int, error) {
    sl.mutex.Lock()
    defer sl.mutex.Unlock()
    
    if sl.size == 0 {
        return food.FoodPack{}, -1, errors.New("storage is empty")
    }
    
    low, high := 0, sl.size-1
    
    // Search for meat version first
    if found, index, err := sl.binarySearchForType(desiredType, 'M', low, high); err == nil {
        return found, index, nil
    }
    
    // Search for grain/vegetable version
    if found, index, err := sl.binarySearchForType(desiredType, 'B', low, high); err == nil {
        return found, index, nil
    }
    
    // Not found, return last item with apology
    lastIndex := sl.size - 1
    lastItem := sl.items[lastIndex]
    sl.removeAtIndex(lastIndex)
    
    fmt.Printf("Sorry, no food packets of the %s type are currently available\n", desiredType)
    return lastItem, lastIndex, errors.New("desired type not available")
}

// Helper function for binary search
func (sl *SortedList) binarySearchForType(desiredType food.FoodType, shipment byte, low, high int) (food.FoodPack, int, error) {
    target := food.FoodPack{
        FoodType:     desiredType,
        FoodShipment: shipment,
    }
    
    for low <= high {
        mid := low + (high-low)/2
        current := sl.items[mid]
        
        comparison := sl.compare(current, target)
        
        if comparison == 0 {
            // Found exact match
            found := current
            sl.removeAtIndex(mid)
            return found, mid, nil
        } else if comparison < 0 {
            low = mid + 1
        } else {
            high = mid - 1
        }
    }
    
    return food.FoodPack{}, -1, errors.New("not found")
}

// Comparison function for FoodPack items
func (sl *SortedList) compare(a, b food.FoodPack) int {
    // Meat comes before grain/vegetable
    if a.FoodShipment != b.FoodShipment {
        if a.FoodShipment == 'M' {
            return -1
        }
        return 1
    }
    
    // Same shipment type, compare by food type
    return int(a.FoodType) - int(b.FoodType)
}

// Remove element at index efficiently
func (sl *SortedList) removeAtIndex(index int) {
    if index < 0 || index >= sl.size {
        return
    }
    
    // Shift elements to fill the gap
    copy(sl.items[index:], sl.items[index+1:])
    
    // Truncate slice and update size
    sl.items = sl.items[:sl.size-1]
    sl.size--
    
    // Optional: Zero out the last element to help garbage collection
    if cap(sl.items) > sl.size {
        sl.items = sl.items[:cap(sl.items)]
        sl.items[sl.size] = food.FoodPack{}  // Zero value
        sl.items = sl.items[:sl.size]
    }
}

// Check if storage is empty
func (sl *SortedList) IsEmpty() bool {
    sl.mutex.RLock()
    defer sl.mutex.RUnlock()
    return sl.size == 0
}

// Check if storage is full
func (sl *SortedList) IsFull() bool {
    sl.mutex.RLock()
    defer sl.mutex.RUnlock()
    return sl.size >= sl.capacity
}

// Get current size
func (sl *SortedList) Size() int {
    sl.mutex.RLock()
    defer sl.mutex.RUnlock()
    return sl.size
}
// AcceptMessage handles incoming food packages and inserts them into storage.
func (sl *SortedList) AcceptMessage(item food.FoodPack) error {
	// Reuses your existing thread-safe Insert logic
	err := sl.Insert(item)
	if err != nil {
		return fmt.Errorf("AcceptMessage failed: %w", err)
	}
	return nil
}

// RetrieveMessage searches for a specific food type and retrieves it.
// If the desired type isn't available, it falls back to your FindBinary fallback behavior.
func (sl *SortedList) RetrieveMessage(desiredType food.FoodType) (food.FoodPack, error) {
	// Reuses your binary search logic
	item, _, err := sl.FindBinary(desiredType)
	if err != nil && err.Error() != "desired type not available" {
		return food.FoodPack{}, fmt.Errorf("RetrieveMessage failed: %w", err)
	}

	// Returns the found item or the fallback item (when desired type is missing)
	return item, nil
}