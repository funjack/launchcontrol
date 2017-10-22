package kiiroo

import (
	"sort"
	"testing"
	"time"
)

func TestEventMarshalText(t *testing.T) {
	want := "1.23:2"
	e := Event{
		Time:  time.Millisecond * 1234,
		Value: 2,
	}

	v, err := e.MarshalText()
	if err != nil {
		t.Error(err)
	}

	if want != string(v) {
		t.Errorf("want: %q, got %q", want, v)
	}
}

func TestEventUnmarshalText(t *testing.T) {
	want := Event{
		Time:  time.Millisecond * 1230,
		Value: 2,
	}

	e := new(Event)
	err := e.UnmarshalText([]byte("1.23:2"))
	if err != nil {
		t.Error(err)
	}

	if want.Time != e.Time {
		t.Errorf("time does not match, want: %s, got: %s", want.Time, e.Time)
	}
	if want.Value != e.Value {
		t.Errorf("value does not match, want: %d, got: %d", want.Value, e.Value)
	}
}

func TestEventsMarshalText(t *testing.T) {
	want := "{1.23:2,1.50:4,3.00:0}"
	es := Events{
		{
			Time:  time.Millisecond * 1230,
			Value: 2,
		},
		{
			Time:  time.Millisecond * 1500,
			Value: 4,
		},
		{
			Time:  time.Millisecond * 3000,
			Value: 0,
		},
	}

	v, err := es.MarshalText()
	if err != nil {
		t.Error(err)
	}

	if want != string(v) {
		t.Errorf("want: %q, got %q", want, v)
	}
}

func TestEventsUnmarshalText(t *testing.T) {
	want := Events{
		{
			Time:  time.Millisecond * 1230,
			Value: 2,
		},
		{
			Time:  time.Millisecond * 1500,
			Value: 4,
		},
		{
			Time:  time.Millisecond * 3000,
			Value: 0,
		},
	}

	var es Events
	err := es.UnmarshalText([]byte("{1.23:2,1.50:4,3.00:0}"))
	if err != nil {
		t.Error(err)
	}

	if len(want) != len(es) {
		t.Errorf("length does not match, want: %d, got %d", len(want), len(es))
	}
	for i := range want {
		if want[i].Time != es[i].Time {
			t.Errorf("time does not match, want: %s, got: %s", want[i].Time, es[i].Time)
		}
		if want[i].Value != es[i].Value {
			t.Errorf("value does not match, want: %d, got: %d", want[i].Value, es[i].Value)
		}
	}
}

func TestEventsSeort(t *testing.T) {
	want := Events{
		{
			Time:  time.Millisecond * 1230,
			Value: 2,
		},
		{
			Time:  time.Millisecond * 1500,
			Value: 4,
		},
		{
			Time:  time.Millisecond * 3000,
			Value: 0,
		},
	}

	var es Events
	err := es.UnmarshalText([]byte("{1.50:4,1.23:2,3.00:0}"))
	if err != nil {
		t.Error(err)
	}
	sort.Sort(es)

	if len(want) != len(es) {
		t.Errorf("length does not match, want: %d, got %d", len(want), len(es))
	}
	for i := range want {
		if want[i].Time != es[i].Time {
			t.Errorf("time does not match, want: %s, got: %s", want[i].Time, es[i].Time)
		}
		if want[i].Value != es[i].Value {
			t.Errorf("value does not match, want: %d, got: %d", want[i].Value, es[i].Value)
		}
	}
}
