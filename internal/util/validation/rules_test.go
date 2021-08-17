package validation_test

import (
	validation "dns/internal/util/validation"
	"github.com/shopspring/decimal"
	"testing"
)

func TestDecimalMaxAbs(t *testing.T) {
	tt := []struct {
		description string
		request     interface{}
		max         decimal.Decimal
		isError     bool
	}{
		{
			description: "valid input success, positive numbers, val < max",
			request:     decimal.NewFromFloat(10.49),
			max:         decimal.NewFromFloat(10.5),
			isError:     false,
		},
		{
			description: "valid input success, positive numbers, val == max",
			request:     decimal.NewFromFloat(10.5),
			max:         decimal.NewFromFloat(10.5),
			isError:     false,
		},
		{
			description: "valid input success, positive numbers, val > max",
			request:     decimal.NewFromFloat(10.51),
			max:         decimal.NewFromFloat(10.5),
			isError:     true,
		},
		{
			description: "valid input success negative numbers, abs comparison is done",
			request:     decimal.NewFromFloat(-10.49),
			max:         decimal.NewFromFloat(-10.5),
			isError:     false,
		},
		{
			description: "request value id not decimal fail 1",
			request:     11.345,
			max:         decimal.NewFromFloat(-10.45),
			isError:     true,
		},
		{
			description: "request value id not decimal fail 2",
			request:     "qwerty",
			max:         decimal.NewFromFloat(-10.45),
			isError:     true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			validator := validation.DecimalMaxAbs(tc.max)
			e := validator(tc.request)

			if tc.isError && e == nil {
				t.Error("expected error, but err is nil")
			}

			if !tc.isError && e != nil {
				t.Errorf("unexpected error: %s", e)
			}
		})
	}
}
