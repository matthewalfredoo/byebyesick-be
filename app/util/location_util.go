package util

import (
	"encoding/json"
	"fmt"
	"halodeksik-be/app/appconfig"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"io/ioutil"
	"net/http"
	"strings"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Bounds struct {
	Northeast Location `json:"northeast"`
	Southwest Location `json:"southwest"`
}

type Geometry struct {
	Bounds       Bounds   `json:"bounds"`
	Location     Location `json:"location"`
	LocationType string   `json:"location_type"`
	Viewport     Bounds   `json:"viewport"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Result struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormattedAddress  string             `json:"formatted_address"`
	Geometry          Geometry           `json:"geometry"`
	PlaceID           string             `json:"place_id"`
	Types             []string           `json:"types"`
}

type PlusCode struct {
	CompoundCode string `json:"compound_code"`
	GlobalCode   string `json:"global_code"`
}

type LocationJSONData struct {
	PlusCode PlusCode `json:"plus_code"`
	Results  []Result `json:"results"`
	Status   string   `json:"status"`
}

type LocationUtil interface {
	ValidateLatLong(city string, province string, lat string, long string) error
}

func NewLocationUtil(region string) LocationUtil {

	return &LocationUtilImpl{
		googleUrl:    appconfig.Config.GmapUrl,
		region:       region,
		responseType: "json",
		apiKey:       appconfig.Config.GmapKey,
		shortenedArea: map[string]string{
			"DKI Jakarta":   "Daerah Khusus Ibukota Jakarta",
			"DI Yogyakarta": "Daerah Istimewa Yogyakarta",
		},
	}
}

type LocationUtilImpl struct {
	googleUrl     string
	responseType  string
	region        string
	apiKey        string
	shortenedArea map[string]string
}

func (l *LocationUtilImpl) ValidateLatLong(city string, province string, lat string, long string) error {
	url := fmt.Sprintf("%s/%s?latlng=%s,%s&key=%s&result_type=%s&language=%s&region=%s", l.googleUrl, l.responseType, lat, long, l.apiKey, appconstant.AreaCityLevel, "id", "id")
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var data LocationJSONData

	err = json.Unmarshal(responseData, &data)
	if err != nil {
		return err
	}

	if data.Status != appconstant.GMAPStatusOK {
		return apperror.ErrInvalidLatLong
	}

	latLongCity, latLongProvince := l.getCityAndProvince(data)

	latLongCity = l.parseRegency(latLongCity)

	if !strings.Contains(city, latLongCity) {
		return apperror.ErrInvalidLatLong
	}

	province = l.getLongProvinceName(province)

	if !strings.Contains(province, latLongProvince) {
		return apperror.ErrInvalidLatLong
	}

	return nil

}

func (l *LocationUtilImpl) getLongProvinceName(province string) string {
	if val, ok := l.shortenedArea[province]; ok {
		return val
	}
	return province
}

func (l *LocationUtilImpl) getCityAndProvince(data LocationJSONData) (string, string) {
	var latLongCity string
	var latLongProvince string

	if len(data.Results) > 0 {
		loc := data.Results[0]
		for _, address := range loc.AddressComponents {
			for _, aType := range address.Types {
				if aType == appconstant.AreaCityLevel {
					latLongCity = address.LongName
					break
				} else if aType == appconstant.AreaProvinceLevel {
					latLongProvince = address.LongName
					break
				}
			}
		}
		return latLongCity, latLongProvince
	} else {
		return "", ""
	}
}

func (l *LocationUtilImpl) parseRegency(city string) string {
	splitCity := strings.Split(city, " ")
	if splitCity[0] == appconstant.RegencyId {
		rep := []string{appconstant.RegencyIdShort}
		rep = append(rep, splitCity[1])
		return strings.Join(rep, " ")
	}
	return city
}
