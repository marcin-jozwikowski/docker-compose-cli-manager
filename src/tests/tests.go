package tests

import "testing"

func AssertErrorEquals(t *testing.T, expected string, err error) {
	if err == nil {
		t.Errorf("Expected error: %s, got nil", expected)
	}

	if err.Error() != expected {
		t.Errorf("Invalid error. Expected %s, got %s", expected, err)
	}
}

func AssertStringEquals(t *testing.T, expected, value, name string) {
	if value != expected {
		t.Errorf("Invalid string value on %s. Expected %s, got %s", name, expected, value)
	}
}

func AssertIntEquals(t *testing.T, expected, actual int, name string) {
	if expected != actual {
		t.Errorf("Invalid int value on %s. Expected %d, got %d", name, expected, actual)
	}
}

func AssertNil(t *testing.T, obj interface{}, description string) {
	if obj != nil {
		t.Errorf("Unexpected value in %s. Expected nil, got %+v", description, obj)
	}
}

func AssertBooleanEquals(t *testing.T, expected bool, value bool, name string) {
	if value != expected {
		t.Errorf("Invalid %s. Expected %t got %t", name, expected, value)
	}
}
