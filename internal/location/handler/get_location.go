package handler

import (
	"dns/internal/location/dto"
	"dns/internal/location/service"
	errUtil "dns/internal/util/err"
	httpUtil "dns/internal/util/http"
	"fmt"
	"net/http"
)

type GetLocationHandler struct {
	locationCalculator service.LocationCalculator
}

func NewGetLocationHandler(c service.LocationCalculator) GetLocationHandler {
	return GetLocationHandler{locationCalculator: c}
}

func (h GetLocationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	req, err := dto.NewGetLocationRequest(r)
	if err != nil {
		httpUtil.HandleValidationError(w, err)
		return
	}
	if err = req.Validate(); err != nil {
		httpUtil.HandleValidationError(w, err)
		return
	}
	loc := h.locationCalculator.Calculate(req)
	res := dto.GetLocationResponse{Location: loc}
	if err := httpUtil.WriteJson(w, res); err != nil {
		fmt.Printf(errUtil.HttpResponseError+"\n", err)
	}
}
