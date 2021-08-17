package service

import (
	"dns/internal/location/dto"
	"dns/internal/util/err"
	"errors"
	"github.com/shopspring/decimal"
)

type LocationCalculator interface {
	Calculate(req dto.GetLocationRequest) (decimal.Decimal, error)
}

type DefaultCalculator struct {
	sectorID int
}

func NewDefaultCalculator(sectorId int) *DefaultCalculator {
	return &DefaultCalculator{sectorID: sectorId}
}

func (s *DefaultCalculator) Calculate(req dto.GetLocationRequest) (decimal.Decimal, error) {
	if s.sectorID <= 0 {
		return decimal.Decimal{}, errors.New(err.SectorIDIsLessThanZero)
	}
	sectorID := decimal.NewFromInt(int64(s.sectorID))
	newX := req.X.Mul(sectorID)
	newY := req.Y.Mul(sectorID)
	newZ := req.Z.Mul(sectorID)
	return newX.Add(newY).Add(newZ).Add(req.Vel), nil
}
