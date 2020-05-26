package check

import (
	"testing"
)

func TestNewMonitor(t *testing.T) {
	m, err := New("Alerty", "https://alerty.online", 10, 1)
	if err != nil {
		t.Errorf("%v", err)
	}
	if m.Port != 443 {
		t.Errorf("Expected 443 port in monitor, got %v", m.Port)
	}
}

func TestNewMonitor2(t *testing.T) {
	m, err := New("Alerty Socket", "tcp://192.0.0.1:8080", 10, 1)
	
	if err != nil {
		t.Errorf("%v", err)
	}
	if m.Port != 8080 {
		t.Errorf("Expected 8080 port in monitor, got %v", m.Port)
	}
}
