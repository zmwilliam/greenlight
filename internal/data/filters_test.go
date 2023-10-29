package data_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/zmwilliam/greenlight/internal/data"
	"github.com/zmwilliam/greenlight/internal/validator"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		desc            string
		input           data.Filters
		expected_errors map[string]string
	}{
		{
			desc: "filters are valid",
			input: data.Filters{
				Page:         1,
				PageSize:     1,
				Sort:         "id",
				SortSafelist: []string{"id"},
			},
			expected_errors: map[string]string{},
		},
		{
			desc:  "page and page_size must be GT zero",
			input: data.Filters{Sort: "id", SortSafelist: []string{"id"}},
			expected_errors: map[string]string{
				"page":      "must be greater than zero",
				"page_size": "must be greater than zero",
			},
		},
		{
			desc: "page and page_size must respect max value",
			input: data.Filters{
				Page:         10_000_001,
				PageSize:     101,
				Sort:         "id",
				SortSafelist: []string{"id"},
			},
			expected_errors: map[string]string{
				"page":      "must be a maximum of 10 million",
				"page_size": "must be a maximum of 100",
			},
		},
		{
			desc: "sort value is invalid",
			input: data.Filters{
				Page:         1,
				PageSize:     1,
				Sort:         "name",
				SortSafelist: []string{"id"},
			},
			expected_errors: map[string]string{
				"sort": "invalid sort value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			v := validator.New()
			tt.input.Validate(v)

			if len(tt.expected_errors) > 0 && v.Valid() {
				t.Error("Filters should not be valid")
			}

			if diff := cmp.Diff(tt.expected_errors, v.Errors); diff != "" {
				t.Errorf("validation errors does not match (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestSortValue(t *testing.T) {
	t.Run("valid sort value without prefix", func(t *testing.T) {
		input := data.Filters{
			Sort:         "name",
			SortSafelist: []string{"id", "name"},
		}

		expected_value := "name"
		got := input.SortValue()

		if got != expected_value {
			t.Errorf("expected %s, got %s", expected_value, got)
		}
	})

	t.Run("valid sort value with prefix", func(t *testing.T) {
		input := data.Filters{
			Sort:         "-name",
			SortSafelist: []string{"id", "name"},
		}

		expected_value := "name"
		got := input.SortValue()

		if got != expected_value {
			t.Errorf("expected %s, got %s", expected_value, got)
		}
	})

	t.Run("panic when safe list is empty", func(t *testing.T) {
		defer func() {
			expected_reason := "invalid sort value"
			if reason := recover(); reason != expected_reason {
				t.Errorf("expected to panic with reason: %s, got %v", expected_reason, reason)
			}
		}()

		input := data.Filters{
			Sort:         "name",
			SortSafelist: []string{},
		}

		input.SortValue()
	})

	t.Run("panic when value not in safe list", func(t *testing.T) {
		defer func() {
			expected_reason := "invalid sort value"
			if reason := recover(); reason != expected_reason {
				t.Errorf("expected to panic with reason: %s, got %v", expected_reason, reason)
			}
		}()

		input := data.Filters{
			Sort:         "name",
			SortSafelist: []string{"id", "inserted_at"},
		}

		input.SortValue()
	})
}

func TestSortDirection(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "DESC when value has prefix",
			input:    "-name",
			expected: "DESC",
		},
		{
			desc:     "ASC when value has no prefix",
			input:    "name",
			expected: "ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			input := data.Filters{Sort: tt.input}

			expected := tt.expected
			got := input.SortDirection()

			if got != expected {
				t.Errorf("expected %s sort direction, got %s", expected, got)
			}
		})
	}
}
