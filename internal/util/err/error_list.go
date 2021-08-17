package err

const (
	ValueIsNotFloat        = "value is not float"
	HttpResponseError      = "error while writing to response: %s"
	SectorIDIsLessThanZero = "sectorID must be >= 0"
	ValueExceedsMax        = "must be no greater than %v"
	DecimalExpected        = "value is not a decimal"
)
