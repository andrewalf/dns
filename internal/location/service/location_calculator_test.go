package service_test

import (
	"dns/internal/location/dto"
	"dns/internal/location/service"
	"github.com/shopspring/decimal"
	"testing"
)

func TestDefaultCalculator_Calculate(t *testing.T) {
	tt := []struct {
		description   string
		request       dto.GetLocationRequest
		sectorID      int
		expected      decimal.Decimal
		errorExpected bool
	}{
		{
			description: "Test case valid success 1",
			request: dto.GetLocationRequest{
				X:   decimal.NewFromFloat(1.1),
				Y:   decimal.NewFromFloat(1.1),
				Z:   decimal.NewFromFloat(1.1),
				Vel: decimal.NewFromFloat(2),
			},
			sectorID:      1,
			expected:      decimal.NewFromFloat(1.1*3 + 2),
			errorExpected: false,
		},
		{
			description: "Test case valid success 2",
			request: dto.GetLocationRequest{
				X:   decimal.NewFromFloat(1.1),
				Y:   decimal.NewFromFloat(1.1),
				Z:   decimal.NewFromFloat(1.1),
				Vel: decimal.NewFromFloat(2),
			},
			sectorID:      2,
			expected:      decimal.NewFromFloat(1.1*2*3 + 2),
			errorExpected: false,
		},
		{
			description: "Test case valid success 3",
			request: dto.GetLocationRequest{
				X:   decimal.NewFromFloat(0),
				Y:   decimal.NewFromFloat(10),
				Z:   decimal.NewFromFloat(-3.3),
				Vel: decimal.NewFromFloat(2.34),
			},
			sectorID:      2,
			expected:      decimal.NewFromFloat(10*2 + (-3.3 * 2) + 2.34),
			errorExpected: false,
		},
		{
			description:   "Test case sectorID less then zero fail",
			request:       dto.GetLocationRequest{},
			sectorID:      -2,
			errorExpected: false,
		},
		{
			description:   "Test case sectorID is zero fail",
			request:       dto.GetLocationRequest{},
			sectorID:      0,
			errorExpected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			s := service.NewDefaultCalculator(tc.sectorID)
			r, e := s.Calculate(tc.request)

			if tc.errorExpected && e != nil {
				t.Error("expected error, but err is nil")
			}

			if !tc.expected.Equal(r) {
				t.Errorf("expected '%s', got '%s'", tc.expected.String(), r.String())
			}
		})
	}
}
