// The software suite consists of main.go,
// food, stats, gatekeeper,
// producer, circularque, and sales packages.
//
// This represents the software to manage an "embedded" planetary system
// food receiving and distribution system.

package stats

import (
    "math/rand"
    "time"
)

type Stats struct {
    rand *rand.Rand
}

func NewStats() *Stats {
    return &Stats{
        rand: rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

// Exponential distribution using interpolation
func (s *Stats) NextExponential() float64 {
    x := s.rand.Float64()
    
    if x == 0.0 {
        return 0.0
    } else if x <= 0.1 {
        return (x-0.0)*1.04 + 0.0
    } else if x <= 0.2 {
        return (x-0.1)*1.18 + 0.104
    } else if x <= 0.3 {
        return (x-0.2)*1.33 + 0.222
    } else if x <= 0.4 {
        return (x-0.3)*1.54 + 0.355
    } else if x <= 0.5 {
        return (x-0.4)*1.81 + 0.509
    } else if x <= 0.6 {
        return (x-0.5)*2.25 + 0.690
    } else if x <= 0.7 {
        return (x-0.6)*2.85 + 0.915
    } else if x <= 0.75 {
        return (x-0.70)*3.60 + 1.2
    } else if x <= 0.8 {
        return (x-0.75)*4.40 + 1.38
    } else if x <= 0.84 {
        return (x-0.8)*5.75 + 1.60
    } else if x <= 0.88 {
        return (x-0.84)*7.25 + 1.83
    } else if x <= 0.9 {
        return (x-0.88)*9.00 + 2.12
    } else if x <= 0.92 {
        return (x-0.90)*11.0 + 2.30
    } else if x <= 0.94 {
        return (x-0.92)*14.5 + 2.52
    } else if x <= 0.95 {
        return (x-0.94)*18.0 + 2.81
    } else if x <= 0.97 {
        return (x-0.95)*30.0 + 2.99
    } else if x <= 0.98 {
        return (x-0.97)*40.0 + 3.50
    } else if x <= 0.99 {
        return (x-0.98)*70.0 + 3.90
    } else if x <= 0.995 {
        return (x-0.99)*140.0 + 4.60
    } else if x <= 0.999 {
        return (x-0.998)*800.0 + 6.20
    } else {
        return (x-0.9997)*1000.0 + 8.0
    }
}

// Time required to arrange raw food packets for sale - Step function
func (s *Stats) PrepareGrainVegetableFoodPackForSales() time.Duration {
    randFloat := s.rand.Float64()
    if randFloat <= 0.333 {
        return time.Duration(0.25 * float64(time.Second))  // 0.25 hours
    } else if randFloat <= 0.9 {
        return time.Duration(0.66 * float64(time.Second))  // 0.66 hours
    } else {
        return time.Duration(0.75 * float64(time.Second))  // 0.75 hours
    }
}

func (s *Stats) PrepareMeatFoodPackForSales() time.Duration {
    randFloat := s.rand.Float64()
    if randFloat <= 0.20 {
        return time.Duration(0.15 * float64(time.Second))  // 20% - 0.15 hours
    } else if randFloat <= 0.50 {
        return time.Duration(0.35 * float64(time.Second))  // 30% - 0.35 hours
    } else if randFloat <= 0.90 {
        return time.Duration(0.40 * float64(time.Second))  // 40% - 0.40 hours
    } else {
        return time.Duration(1.20 * float64(time.Second))  // 10% - 1.20 hours
    }
}