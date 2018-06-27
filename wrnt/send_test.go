package wrnt

import (
	"testing"
)

func TestSend_Pack_Uninited1(t *testing.T) {
	s := NewSend()
	_, err := s.Pack()
	if err != ErrNotInited {
		t.Error("empty send shall not Pack")
	}
}

func TestSend_Pack_Uninited2(t *testing.T) {
	s := NewSend()
	s.AddItems("a", "b")
	_, err := s.Pack()
	if err != ErrNotInited {
		t.Error("uninited send with added items shall not Pack")
	}
}

func TestSend_Pack_Empty_confidmed(t *testing.T) {
	s := NewSend()
	s.Confirm(100)
	testPackCall(t, s, []string{}, 101)
}

func TestSend_Pack_add_before_confirm(t *testing.T) {
	s := NewSend()
	slice := []string{"a", "b"}
	s.AddItems(slice...)
	s.Confirm(100)
	testPackCall(t, s, slice, 101)
}

func TestSend_Pack_add_after_confirm(t *testing.T) {
	s := NewSend()
	slice := []string{"a", "b"}
	s.Confirm(100)
	s.AddItems(slice...)
	testPackCall(t, s, slice, 101)
}

func TestSend_Pack_multiple_adds_with_partial_confirm(t *testing.T) {
	s := NewSend()
	s.Confirm(100)
	s.AddItems("a", "b", "c")
	s.Confirm(102)
	testPackCall(t, s, []string{"c"}, 103)
	s.AddItems("d", "e", "f")
	testPackCall(t, s, []string{"c", "d", "e", "f"}, 103)
	s.Confirm(106)
	testPackCall(t, s, []string{}, 107)
}

func testPackCall(t *testing.T, s *Send, wait []string, waitN int) {
	got, err := s.Pack()
	if err != nil {
		t.Error(err)
		return
	}
	if got.BaseN != waitN {
		t.Error("wait baseN ", waitN, " got ", got.BaseN)
	}
	if !eqSlices(got.Items, wait) {
		t.Error("wait Items ", wait, " got ", got.Items)
	}
}

func eqSlices(a, b []string) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	} else if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
