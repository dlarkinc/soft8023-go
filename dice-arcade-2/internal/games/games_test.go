package games

import "testing"

func TestFactory(t *testing.T) {
	if _, err := New("highlow"); err != nil {
		t.Fatal(err)
	}
	if _, err := New("pig"); err != nil {
		t.Fatal(err)
	}
	if _, err := New("nope"); err == nil {
		t.Fatal("expected error")
	}
}
