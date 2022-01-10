package tests

import (
	"errors"
	"testing"
)

func TestAssertErrorEquals(t *testing.T) {
	AssertErrorEquals(t, "test", errors.New("test"))
}

func TestAssertIntEquals(t *testing.T) {
	AssertIntEquals(t, 1, 1, "Testing int")
}

func TestAssertNil(t *testing.T) {
	AssertNil(t, nil, "Testing nil")
}

func TestAssertStringEquals(t *testing.T) {
	AssertStringEquals(t, "test", "test", "Testing test")
}

func TestAssertBooleanEquals(t *testing.T) {
	AssertBooleanEquals(t, true, true, "Testing bool")
}
