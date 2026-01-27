package roll

import (
	"crypto/rand"
	"math/big"
)

// DiceRoll represents a dice roll with its configuration and results.
type DiceRoll struct {
	count       int
	sides       int
	results     []int
	threshold   *int
	keepHighest bool
}

// NewDiceRoll creates a new DiceRoll instance.
func NewDiceRoll(count, sides int, threshold *int, keepHighest bool) *DiceRoll {
	return &DiceRoll{
		count:       count,
		sides:       sides,
		results:     []int{},
		threshold:   threshold,
		keepHighest: keepHighest,
	}
}

// Roll performs the dice roll and stores the results.
func (d *DiceRoll) Roll() []int {
	d.results = make([]int, d.count)
	for i := range d.count {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(d.sides)))
		d.results[i] = int(n.Int64()) + 1
	}

	return d.results
}

// GetResults returns the roll results.
func (d *DiceRoll) GetResults() []int {
	return d.results
}

// Sum returns the sum of all dice results.
func (d *DiceRoll) Sum() int {
	sum := 0
	for _, result := range d.results {
		sum += result
	}

	return sum
}

// GetHighest returns the highest die result.
func (d *DiceRoll) GetHighest() int {
	if len(d.results) == 0 {
		return 0
	}

	highest := d.results[0]
	for _, result := range d.results {
		if result > highest {
			highest = result
		}
	}

	return highest
}

// IsCriticalHit checks if all dice show maximum value.
func (d *DiceRoll) IsCriticalHit() bool {
	if len(d.results) == 0 {
		return false
	}

	for _, result := range d.results {
		if result != d.sides {
			return false
		}
	}

	return true
}

// IsCriticalFailure checks if all dice show minimum value (1).
func (d *DiceRoll) IsCriticalFailure() bool {
	if len(d.results) == 0 {
		return false
	}

	for _, result := range d.results {
		if result != 1 {
			return false
		}
	}

	return true
}

// IsSuccess determines if the roll was successful based on threshold
// Returns nil if no threshold is set.
func (d *DiceRoll) IsSuccess() *bool {
	if d.threshold == nil {
		return nil
	}

	var success bool
	if d.keepHighest {
		// Success if ANY die meets threshold
		success = d.GetHighest() >= *d.threshold
	} else {
		// Success if SUM meets threshold
		success = d.Sum() >= *d.threshold
	}

	return &success
}

// GetCount returns the number of dice.
func (d *DiceRoll) GetCount() int {
	return d.count
}

// GetSides returns the number of sides per die.
func (d *DiceRoll) GetSides() int {
	return d.sides
}

// GetThreshold returns the threshold if set.
func (d *DiceRoll) GetThreshold() *int {
	return d.threshold
}

// GetKeepHighest returns whether keep highest mode is enabled.
func (d *DiceRoll) GetKeepHighest() bool {
	return d.keepHighest
}
