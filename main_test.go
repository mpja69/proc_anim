package main

import (
	"fmt"
	"math"
	"testing"
)

func TestDistance(t *testing.T) {
	p := Point{3, 0}
	q := Point{0, 4}
	actual := distance(p, q)
	expected := 5.0
	fmt.Printf("PASS - distance(): actual: %v,  expected: %v\n", actual, expected)
	if actual != expected {
		t.Errorf("ERR: actual: %v,  expected: %v", actual, expected)
	}
}

func TestMag(t *testing.T) {
	p := Point{3, 4}
	actual := p.Mag()
	expected := 5.0
	fmt.Printf("PASS - Mag(): actual: %v,  expected: %v\n", actual, expected)
	if actual != expected {
		t.Errorf("ERR: actual: %v,  expected: %v", actual, expected)
	}
}

func TestAdd(t *testing.T) {
	p := Point{3, 0}
	q := Point{0, 4}
	actual := p.Add(q)
	expected := Point{3, 4}
	fmt.Printf("PASS - Add(): actual: %v,  expected: %v\n", actual, expected)
	if actual != expected {
		t.Errorf("ERR: actual: %v,  expected: %v", actual, expected)
	}
}

func TestSub(t *testing.T) {
	p := Point{3, 0}
	q := Point{0, 4}
	actual := p.Sub(q)
	expected := Point{3, -4}
	fmt.Printf("PASS - Sub(): actual: %v,  expected: %v\n", actual, expected)
	if actual != expected {
		t.Errorf("ERR: actual: %v,  expected: %v", actual, expected)
	}
}

func TestSetMag(t *testing.T) {
	p := Point{6, 8}
	actual := p.SetMag(5)
	expected := Point{3, 4}

	actual.x = math.Round(actual.x)
	actual.y = math.Round(actual.y)
	fmt.Printf("PASS - SetMag(): actual: %v,  expected: %v\n", actual, expected)
	if actual != expected {
		t.Errorf("ERR: actual: %v,  expected: %v", actual, expected)
	}
}

func TestSetConstraint(t *testing.T) {
	p := Point{7, 9}
	q := Point{1, 1}
	actual := SetConstraint(p, q, 5)
	expected := Point{4, 5}

	actual.x = math.Round(actual.x)
	actual.y = math.Round(actual.y)
	fmt.Printf("PASS - SetConstraint(): actual: %v,  expected: %v\n", actual, expected)
	if actual != expected {
		t.Errorf("ERR: actual: %v,  expected: %v", actual, expected)
	}
}
