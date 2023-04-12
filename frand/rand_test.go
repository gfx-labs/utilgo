package frand

// Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the
// terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.

import (
	"log"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	maxIntRange = 999999
)

var (
	stringChoices []string
	intChoices    []int
)

// Test that AlphaStringRange produces a string within specified min/max length
// parameters.  The actual randonimity of the string is not tested.
// nolint:gosec // insecure rand is fine here
func TestAlphaStringRange(t *testing.T) {
	min := rand.Intn(100)
	max := min + 1 + rand.Intn(100)
	s, err := AlphaStringRange(min, max)
	if err != nil {
		t.Error(err)
	}
	switch {
	case len(s) < min:
		t.Error("Random string is too short")
	case len(s) > max:
		t.Error("Random string is too short")
	}
}

// Test that IntRange produces an integer between min and max
func TestIntRange(t *testing.T) {
	min := 567
	max := 890
	i := IntRange(min, max)
	if i > max || i < min {
		t.Error("IntRange returned an out-of-range integer")
	}
}

// Test that the strings we produce are actually random.  This is done by
// comparing two 50,000 character generated random strings and checking that
// they differ.  It is quite unlikely, but not strictly impossible, that two
// truly random strings will be identical.
func TestRandonimity(t *testing.T) {
	l := 50000
	s1, err := AlphaString(l)
	if err != nil {
		t.Error(err)
	}
	s2, err := AlphaString(l)
	if err != nil {
		t.Error(err)
	}
	if s1 == s2 {
		msg := "Generated two identical 'random' strings - this is probably an error"
		t.Error(msg)
	}
}

// TestChoice tests that over the course of 1,000,000 calls on the same 100
// possible choices, the Choice() function returns every possible choice at
// least once.  Note, there is a VERY small chance this test could fail by
// random chance even when the code is working correctly.
func TestChoice(t *testing.T) {
	// Create a map associating each possible choice with a bool.
	chosen := make(map[int]bool)
	for _, v := range intChoices {
		chosen[v] = false
	}
	// Run Choice() a million times, and record which of the possible choices it returns.
	for i := 0; i < 1000000; i++ {
		c := ChoiceInt(intChoices)
		chosen[c] = true
	}
	// Fail if any of the choices was not chosen even once.
	for _, v := range chosen {
		if v == false {
			msg := "In 1,000,000 passes Choice() failed to return all 100 possible choices.  Something is probably wrong."
			t.Error(msg)
		}
	}
}

// TestWeightedChoice assembles a list of choices, weighted 0-9, and tests that
// over the course of 1,000,000 calls to WeightedChoice() each choice is
// returned more often than choices with a lower weight.
func TestWeightedChoice(t *testing.T) {
	// Make weighted choices
	var choices []Choice
	chosenCount := make(map[Choice]int)
	for i := 0; i < 10; i++ {
		c := Choice{
			Weight: i,
			Item:   i,
		}
		choices = append(choices, c)
		chosenCount[c] = 0
	}
	// Run WeightedChoice() a million times, and record how often it returns each
	// of the possible choices.
	for i := 0; i < 1000000; i++ {
		c := WeightedChoice(choices)
		chosenCount[c] += 1
	}
	// Test that higher weighted choices were chosen more often than their lower
	// weighted peers.
	for i, c := range choices[0 : len(choices)-1] {
		next := choices[i+1]
		expr := chosenCount[c] < chosenCount[next]
		assert.True(t, expr)
	}
}

// TestWeightedChoice assembles a list of choices, weighted 0-9, and tests that
// over the course of 1,000,000 calls to WeightedChoice() each choice is
// returned more often than choices with a lower weight.
func TestInverseWeightedChoice(t *testing.T) {
	// Make weighted choices
	var choices []Choice
	chosenCount := make(map[int]int)
	for i := 1; i < 10; i++ {
		c := Choice{
			Weight: i,
			Item:   i,
		}
		choices = append(choices, c)
		chosenCount[c.Weight] = 0
	}
	InvertChoicesT(choices)
	for z := 0; z < 4; z++ {
		// Run WeightedChoice() a million times, and record how often it returns each
		// of the possible choices.
		for i := 0; i < 1000000; i++ {
			c := WeightedChoice(choices)
			chosenCount[c.Weight] += 1
		}
		// Test that lower weighted choices were chosen more often than their higher
		// weighted peers.
		for i, c := range choices[0 : len(choices)-1] {
			next := choices[i+1]
			expr := chosenCount[c.Weight] > chosenCount[next.Weight]
			assert.True(t, expr)
		}
	}
}

// BenchmarkChoiceInt runs a benchmark on the ChoiceInt function.
func BenchmarkChoiceInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ChoiceInt(intChoices)
	}
}

// BenchmarkChoiceString runs a benchmark on the ChoiceString function.
func BenchmarkChoiceString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ChoiceString(stringChoices)
	}
}

// BenchmarkIntRange runs a benchmark on the IntRange function.
func BenchmarkIntRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntRange(0, math.MaxInt32)
	}
}

// BenchmarkIntRange runs a benchmark on the WeightedChoice function.
func BenchmarkWeightedChoice(b *testing.B) {
	// Create some random choices and weights before we start
	b.StopTimer()
	choices := []Choice{}
	for i := 0; i < 100; i++ {
		s, err := AlphaString(64)
		if err != nil {
			b.Error(err)
		}
		w := IntRange(1, 10)
		c := Choice{
			Item:   s,
			Weight: w,
		}
		choices = append(choices, c)
	}
	// Run the benchmark
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		WeightedChoice(choices)
	}
}

// init populates two arrays of random choices, intChoices and stringChoices,
// which will be used by various test and benchmark functions.
func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	// Random integers
	for i := 0; i < 100; i++ {
		randint := IntRange(0, maxIntRange)
		intChoices = append(intChoices, randint)
	}
	// Random strings
	for i := 0; i < 100; i++ {
		randstr, err := AlphaString(32)
		if err != nil {
			log.Panicln(err)
		}
		stringChoices = append(stringChoices, randstr)
	}
}
