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

// GetLocation godoc
// @Router /location [get]
// @Summary Retrieves storages location by drones coordinates and velocity
// @Produce json
//
// @Param x query number false "x coordinate"
// @Param y query number false "y coordinate"
// @Param z query number false "z coordinate"
// @Param vel query number false "drone velocity"
//
// @Success 200 {object} dto.GetLocationResponse
// @Failure 400 {object} httpUtil.ErrorResponse
// @Failure 500 {object} httpUtil.ErrorResponse
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
	loc, err := h.locationCalculator.Calculate(req)
	if err != nil {
		httpUtil.HandleServerError(w, err)
		return
	}
	res := dto.GetLocationResponse{Location: loc}
	if err := httpUtil.WriteJson(w, res); err != nil {
		fmt.Printf(errUtil.HttpResponseError+"\n", err)
	}
}
