package handler_test

import (
	"dns/internal/location/dto"
	"dns/internal/location/handler"
	"dns/internal/location/service"
	errUtil "dns/internal/util/err"
	"fmt"
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetLocationHandler_ServeHTTP(t *testing.T) {
	sectorID := 2
	calculator := service.NewDefaultCalculator(sectorID)
	validReq := "x=1.1&y=1.1&z=1.1&vel=2.0"
	validReqExpectedRes, _ := calculator.Calculate(dto.GetLocationRequest{
		X:   decimal.NewFromFloat(1.1),
		Y:   decimal.NewFromFloat(1.1),
		Z:   decimal.NewFromFloat(1.1),
		Vel: decimal.NewFromFloat(2.0),
	})

	tt := []struct {
		description string
		request     string
		expected    string
		statusCode  int
	}{
		{
			description: "all data is valid success",
			request:     "/location?" + validReq,
			expected:    fmt.Sprintf(`{"loc":"%s"}`, validReqExpectedRes),
			statusCode:  http.StatusOK,
		},
		{
			description: "x is empty validation error",
			request:     "/location?y=1.1&z=1.1&vel=2.0",
			expected:    fmt.Sprintf(`{"errors":{"x":"%s"}}`, errUtil.ValueIsNotFloat),
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "y is empty validation error",
			request:     "/location?x=1.1&z=1.1&vel=2.0",
			expected:    fmt.Sprintf(`{"errors":{"y":"%s"}}`, errUtil.ValueIsNotFloat),
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "z is empty validation error",
			request:     "/location?x=1.1&y=1.1&vel=2.0",
			expected:    fmt.Sprintf(`{"errors":{"z":"%s"}}`, errUtil.ValueIsNotFloat),
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "vel is empty validation error",
			request:     "/location?x=1.1&y=1.1&z=2.0",
			expected:    fmt.Sprintf(`{"errors":{"vel":"%s"}}`, errUtil.ValueIsNotFloat),
			statusCode:  http.StatusBadRequest,
		},
		{
			description: "all request is empty validation error",
			request:     "/location",
			expected: fmt.Sprintf(
				`{"errors":{"vel":"%s","x":"%s","y":"%s","z":"%s"}}`,
				errUtil.ValueIsNotFloat,
				errUtil.ValueIsNotFloat,
				errUtil.ValueIsNotFloat,
				errUtil.ValueIsNotFloat,
			),
			statusCode: http.StatusBadRequest,
		},
		{
			description: "vel is greater than light speed validation error",
			request:     "/location?x=1.1&y=1.1&z=1.1&vel=299793",
			expected: fmt.Sprintf(
				`{"errors":{"Vel":"%s"}}`,
				fmt.Sprintf(errUtil.ValueExceedsMax, 299792),
			),
			statusCode: http.StatusBadRequest,
		},
	}

	h := handler.NewGetLocationHandler(calculator)

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tc.request, nil)
			responseRecorder := httptest.NewRecorder()
			h.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("expected status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, responseRecorder.Body)
			}
		})
	}
}
