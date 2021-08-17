package service

import (
	"dns/internal/location/dto"
	"github.com/shopspring/decimal"
)

type LocationCalculator interface {
	Calculate(req dto.GetLocationRequest) decimal.Decimal
}

type DefaultCalculator struct {
	sectorID int
}

func NewDefaultCalculator(sectorId int) *DefaultCalculator {
	return &DefaultCalculator{sectorID: sectorId}
}

func (s *DefaultCalculator) Calculate(req dto.GetLocationRequest) decimal.Decimal {
	sectorID := decimal.NewFromInt(int64(s.sectorID))
	newX := req.X.Mul(sectorID)
	newY := req.Y.Mul(sectorID)
	newZ := req.Z.Mul(sectorID)
	return newX.Add(newY).Add(newZ).Add(req.Vel)
}
