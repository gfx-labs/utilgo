package generic

import "testing"

func TestMapN(t *testing.T) {
	var m Map11[string, int, bool, int, uint, string, int8, byte, bool, string, int64, string]
	m.Store("a", 1, false, 2, 54, "hello world", -6, 64, true, "hi", 584932, "wow it works")
	v, ok := m.Load("a", 1, false, 2, 54, "hello world", -6, 64, true, "hi", 584932)
	if !ok {
		t.Error("expected value to exist")
		return
	}
	if v != "wow it works" {
		t.Error("expected value to be 'wow it works'")
		return
	}
	v, ok = m.LoadAndDelete("a", 1, false, 2, 54, "hello world", -6, 64, true, "hi", 584932)
	if !ok {
		t.Error("expected value to exist")
		return
	}
	if v != "wow it works" {
		t.Error("expected value to be 'wow it works'")
		return
	}
	v, loaded := m.LoadOrStore("b", 434, true, 43, 320343, "abcdef", -69, 1, false, "woa", 43434343433, "a")
	if loaded {
		t.Error("expected value to not exist")
	}
	if v != "a" {
		t.Error("expected value to be 'a'")
	}
	m.Delete("b", 434, true, 43, 320343, "abcdef", -69, 1, false, "woa", 43434343433)
	_, ok = m.Load("b", 434, true, 43, 320343, "abcdef", -69, 1, false, "woa", 43434343433)
	if ok {
		t.Error("expected value to not exist")
	}
}
