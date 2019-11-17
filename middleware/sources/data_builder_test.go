package sources

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHttpDataBuilder_Header(t *testing.T) {
	t.Run("returns new header every time", func(t *testing.T) {
		builder := httpDataBuilder{}
		h := builder.Header()
		if len(h) > 0 {
			t.Errorf("Headers len is: %v", len(h))
		}
		h.Set("123", "123")
		h = builder.Header()
		if len(h) > 0 {
			t.Errorf("Headers len is: %v", len(h))
		}
	})

	t.Run("write header no panic", func(t *testing.T) {
		builder := httpDataBuilder{}
		builder.WriteHeader(123)
	})

	t.Run("write data succeeds and returns size of data", func(t *testing.T) {
		builder := httpDataBuilder{}
		d1 := []byte("test")
		d2 := []byte("test_123")
		l, err := builder.Write(d1)
		cmperr(t, err, nil)
		if l != len(d1) {
			t.Errorf("Expected to have equal size. Got: %v vs %v", len(d1), l)
		}
		l, err = builder.Write(d2)
		cmperr(t, err, nil)
		if l != len(d2) {
			t.Errorf("Expected to have equal size. Got: %v vs %v", len(d2), l)
		}
	})

	t.Run("build returns nil if nothing was written", func(t *testing.T) {
		builder := httpDataBuilder{}
		res, err := builder.Build()
		cmperr(t, err, nil)
		if res != nil {
			t.Errorf("Expected to nil,  got: %v", res)
		}
		res, err = builder.Build()
		cmperr(t, err, nil)
		if res != nil {
			t.Errorf("Expected to nil,  got: %v", res)
		}

	})

	t.Run("build returns last data written", func(t *testing.T) {
		builder := httpDataBuilder{}
		d1 := []byte("test")
		d2 := []byte("test_123")
		l, err := builder.Write(d1)
		cmperr(t, err, nil)
		if l != len(d1) {
			t.Errorf("Expected to have equal size. Got: %v vs %v", len(d1), l)
		}
		r1, err := builder.Build()
		cmperr(t, err, nil)
		if !cmp.Equal(r1, d1) {
			t.Errorf("Expected %v. Got: %v", d1, r1)
		}

		l, err = builder.Write(d2)
		cmperr(t, err, nil)
		if l != len(d2) {
			t.Errorf("Expected to have equal size. Got: %v vs %v", len(d2), l)
		}
		r2, err := builder.Build()
		cmperr(t, err, nil)
		if !cmp.Equal(r2, d2) {
			t.Errorf("Expected %v. Got: %v", d2, r2)
		}

	})
}
