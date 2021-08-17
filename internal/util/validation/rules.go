package validation

import (
	"dns/internal/util/err"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/shopspring/decimal"
)

func DecimalMaxAbs(max decimal.Decimal) validation.RuleFunc {
	return func(value interface{}) error {
		s, ok := value.(decimal.Decimal)
		if !ok {
			return errors.New(err.DecimalExpected)
		}
		if s.Abs().GreaterThan(max.Abs()) {
			return fmt.Errorf(err.ValueExceedsMax, max.String())
		}
		return nil
	}
}
