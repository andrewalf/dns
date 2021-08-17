package dto

import (
	"dns/internal/util/err"
	"dns/internal/util/validation"
	"errors"
	v "github.com/go-ozzo/ozzo-validation"
	"github.com/shopspring/decimal"
	"net/http"
)

type GetLocationResponse struct {
	Location decimal.Decimal `json:"loc"`
}

// we need decimal in order to save precision,
// I think it's important for coordinates calculation
type GetLocationRequest struct {
	X   decimal.Decimal
	Y   decimal.Decimal
	Z   decimal.Decimal
	Vel decimal.Decimal
}

func (r GetLocationRequest) Validate() error {
	// in terms of space i think that drones dto the future fly with speed measured in km/sec :)
	lightSpeedKmPerSec := decimal.NewFromInt(299792)

	// all other fields don't need validation, it was already done in dto while parsing request
	// and I guess theres no coordinates restrictions, space is endless (theoretically haha)
	return v.ValidateStruct(&r,
		v.Field(&r.Vel, v.By(validation.DecimalMaxAbs(lightSpeedKmPerSec))),
	)
}

func NewGetLocationRequest(r *http.Request) (GetLocationRequest, error) {
	const (
		x   = "x"
		y   = "y"
		z   = "z"
		vel = "vel"
	)

	errs := v.Errors{}
	values := make(map[string]decimal.Decimal, 4)

	for _, key := range [4]string{x, y, z, vel} {
		raw := r.URL.Query().Get(key)
		d, e := decimal.NewFromString(raw)
		if e != nil {
			errs[key] = errors.New(err.ValueIsNotFloat)
		}
		values[key] = d
	}

	if len(errs) != 0 {
		return GetLocationRequest{}, errs
	}

	return GetLocationRequest{
		X:   values[x],
		Y:   values[y],
		Z:   values[z],
		Vel: values[vel],
	}, nil
}
