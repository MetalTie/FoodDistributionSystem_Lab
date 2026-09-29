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
    // "fmt"
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

// Sequential search (C Option requirement)
func (sl *SortedList) FindSequential(desiredType food.FoodType) (food.FoodPack, int, error) {
    sl.mutex.Lock()
    defer sl.mutex.Unlock()
    
    for i := 0; i < sl.size; i++ {
        if sl.items[i].FoodType == desiredType {
            found := sl.items[i]
            sl.removeAtIndex(i)
            return found, i, nil
        }
    }
    
    // Not found, return last item if available
    if sl.size > 0 {
        last := sl.items[sl.size-1]
        sl.removeAtIndex(sl.size-1)
        return last, sl.size-1, errors.New("desired type not available")
    }
    
    return food.FoodPack{}, -1, errors.New("storage is empty")
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

// func (q *CircularQue[T]) AcceptMessage(msg T) error {
//     if q.mesnum >= q.capacity {
//         return errors.New("ERROR - Message rejected - queue is full!")
//     }
//     q.rear = (q.rear + 1) % q.capacity
//     q.box[q.rear] = msg
//     q.mesnum++
//     return nil
// }

// func (q *CircularQue[T]) RetrieveMessage() (T, error) {
//     var zero T
//     if q.mesnum <= 0 {
//         return zero, errors.New("ERROR - No message in the queue to retrieve!")
//     }
//     q.front = (q.front + 1) % q.capacity
//     msg := q.box[q.front]
//     q.mesnum--
//     return msg, nil
// }