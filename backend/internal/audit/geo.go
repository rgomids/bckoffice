package audit

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// GeoInfo contém informações geográficas básicas.
type GeoInfo struct {
	Country string  `json:"country"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

// GeoService consulta informações geográficas de um IP.
type GeoService interface {
	Lookup(ctx context.Context, ip string) (GeoInfo, error)
}

// HttpGeoService implementa GeoService via chamada HTTP pública.
type HttpGeoService struct {
	client *http.Client
}

// NewHttpGeoService cria um HttpGeoService com timeout configurado.
func NewHttpGeoService() *HttpGeoService {
	timeout := 5 * time.Second
	if v := os.Getenv("GEO_TIMEOUT_MS"); v != "" {
		if ms, err := time.ParseDuration(v + "ms"); err == nil {
			timeout = ms
		}
	}
	return &HttpGeoService{client: &http.Client{Timeout: timeout}}
}

func (h *HttpGeoService) Lookup(ctx context.Context, ip string) (GeoInfo, error) {
	url := "https://ipapi.co/" + ip + "/json/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return GeoInfo{}, err
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return GeoInfo{}, err
	}
	defer resp.Body.Close()
	var out struct {
		Country string  `json:"country_name"`
		City    string  `json:"city"`
		Lat     float64 `json:"latitude"`
		Lon     float64 `json:"longitude"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return GeoInfo{}, err
	}
	return GeoInfo{Country: out.Country, City: out.City, Lat: out.Lat, Lon: out.Lon}, nil
}
