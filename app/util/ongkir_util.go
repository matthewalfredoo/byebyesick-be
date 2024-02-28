package util

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/apperror"
	"io/ioutil"
	"net/http"
	"strings"
)

type RajaongkirResponse struct {
	Rajaongkir Rajaongkir `json:"rajaongkir"`
}

type Rajaongkir struct {
	Status  Status    `json:"status"`
	Results []Results `json:"results,omitempty"`
}

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type Results struct {
	Costs []Costs `json:"costs"`
}

type Costs struct {
	Cost []Cost `json:"cost"`
}

type Cost struct {
	Value int32 `json:"value"`
}

type OngkirUtil interface {
	GetCost(originCityId, destCityId int64, weight int32) (decimal.Decimal, error)
}

type RajaOngkirUtil struct {
	url         string
	contentType string
	apiKey      string
}

func NewRajaOngkirUtil() *RajaOngkirUtil {
	return &RajaOngkirUtil{
		url:         appconfig.Config.RajaongkirUrl,
		contentType: "application/x-www-form-urlencoded",
		apiKey:      appconfig.Config.RajaongkirKey,
	}
}

func (util *RajaOngkirUtil) GetCost(originCityId, destCityId int64, weight int32) (decimal.Decimal, error) {
	req, err := http.NewRequest(
		"POST",
		util.url,
		strings.NewReader(
			fmt.Sprintf(
				"origin=%d&destination=%d&weight=%d&courier=pos",
				originCityId,
				destCityId,
				weight,
			),
		),
	)
	if err != nil {
		return decimal.Zero, err
	}
	req.Header.Add("Key", util.apiKey)
	req.Header.Add("Content-Type", util.contentType)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return decimal.Zero, err
	}

	var rajaOngkirResponse RajaongkirResponse
	if err := json.Unmarshal(body, &rajaOngkirResponse); err != nil {
		return decimal.Zero, err
	}
	if rajaOngkirResponse.Rajaongkir.Status.Code != http.StatusOK {
		return decimal.Zero, apperror.ErrGetShipmentCost
	}

	return decimal.NewFromInt32(rajaOngkirResponse.Rajaongkir.Results[0].Costs[0].Cost[0].Value), nil
}
