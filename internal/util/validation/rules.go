package validation

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/shopspring/decimal"
)

func DecimalMaxAbs(max decimal.Decimal) validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(decimal.Decimal)
		if !ok {
			return errors.New("input is not decimal")
		}
		if s.Abs().GreaterThan(max) {
			return fmt.Errorf("must be no greater than %v", max.String())
		}
		return nil
	}
}
