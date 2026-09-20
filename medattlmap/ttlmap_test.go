package medattlmap

import (
	"testing"
	"time"
)

// TestTouch extends a live item's expiration (idle-timeout semantics).
func TestTouch(t *testing.T) {
	m := NewTTLMap(50*time.Millisecond, 10*time.Millisecond)
	defer m.Stop()

	m.Put("k", 0, "v")
	time.Sleep(30 * time.Millisecond) // still within the original 50ms TTL

	// Touch must extend the life past the original TTL.
	if !m.Touch("k") {
		t.Fatal("Touch on a live item must report true")
	}
	time.Sleep(30 * time.Millisecond) // past original expiry, inside touched window
	if _, ok := m.Get("k"); !ok {
		t.Fatal("touched item expired before the idle window")
	}

	// Stop touching → item expires on its own.
	time.Sleep(40 * time.Millisecond)
	if _, ok := m.Get("k"); ok {
		t.Fatal("item should have expired after the idle window elapsed")
	}
}

// TestTouchExpired deletes an already-expired item instead of resurrecting it.
func TestTouchExpired(t *testing.T) {
	m := NewTTLMap(20*time.Millisecond, 10*time.Millisecond)
	defer m.Stop()

	m.Put("k", 0, "v")
	time.Sleep(40 * time.Millisecond) // let it expire

	if m.Touch("k") {
		t.Fatal("Touch on an expired item must report false")
	}
	if _, ok := m.Get("k"); ok {
		t.Fatal("expired item must not be resurrected by Touch")
	}
}

// TestTouchMissing is a no-op for unknown keys.
func TestTouchMissing(t *testing.T) {
	m := NewTTLMap(time.Minute, 10*time.Millisecond)
	defer m.Stop()

	if m.Touch("nope") {
		t.Fatal("Touch on a missing key must report false")
	}
}
