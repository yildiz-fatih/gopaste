package models

import (
	"errors"
	"strings"
	"testing"
)

func TestRandomSlug(t *testing.T) {
	t.Run("correct length", func(t *testing.T) {
		lengths := []int{1, 8, 16}

		for _, want := range lengths {
			slug, err := randomSlug(want)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := len(slug)

			if want != got {
				t.Errorf("want %d, got %d", want, got)
			}
		}
	})
	t.Run("error for non-positive length", func(t *testing.T) {
		lengths := []int{0, -1, -100}

		for _, length := range lengths {
			_, err := randomSlug(length)
			if !errors.Is(err, errInvalidLength) {
				t.Errorf("expected an error for length %d", length)
			}
		}
	})
	t.Run("only uses allowed characters", func(t *testing.T) {
		// intentionally duplicated so the test is independent of the implementation
		const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

		slug, err := randomSlug(50)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, char := range slug {
			if !strings.ContainsRune(allowedChars, char) {
				t.Errorf("contains disallowed character %q", char)
			}
		}
	})
}
