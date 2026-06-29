package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	client *http.Client
}

func NewHandler() *Handler {
	return &Handler{
		client: &http.Client{
			Timeout: 12 * time.Second,
		},
	}
}

type Place struct {
	Name      string  `json:"name,omitempty"`
	City      string  `json:"city,omitempty"`
	Town      string  `json:"town,omitempty"`
	Village   string  `json:"village,omitempty"`
	District  string  `json:"district,omitempty"`
	County    string  `json:"county,omitempty"`
	Region    string  `json:"region,omitempty"`
	State     string  `json:"state,omitempty"`
	Country   string  `json:"country,omitempty"`
	Lat       string  `json:"lat,omitempty"`
	Lon       string  `json:"lon,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type nominatimSearchItem struct {
	DisplayName string            `json:"display_name"`
	Lat         string            `json:"lat"`
	Lon         string            `json:"lon"`
	Address     map[string]string `json:"address"`
}

type nominatimReverseItem struct {
	DisplayName string            `json:"display_name"`
	Lat         string            `json:"lat"`
	Lon         string            `json:"lon"`
	Address     map[string]string `json:"address"`
}

func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	g := r.Group("/geo")
	g.GET("/search", h.Search)
	g.GET("/reverse", h.Reverse)
}

func (h *Handler) Search(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if len(q) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query must be at least 2 characters",
		})
		return
	}

	endpoint := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Set("q", q)
	params.Set("format", "jsonv2")
	params.Set("addressdetails", "1")
	params.Set("limit", "10")
	params.Set("countrycodes", "gh")

	req, err := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodGet,
		endpoint+"?"+params.Encode(),
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create geo request"})
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "BlomivaSchool/1.0 contact@blomiva.com")

	res, err := h.client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "geo provider unavailable"})
		return
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		c.JSON(http.StatusBadGateway, gin.H{"error": "geo provider failed"})
		return
	}

	var raw []nominatimSearchItem
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "invalid geo provider response"})
		return
	}

	places := make([]Place, 0, len(raw))
	for _, item := range raw {
		places = append(places, normalizePlace(item.DisplayName, item.Lat, item.Lon, item.Address))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": places,
	})
}

func (h *Handler) Reverse(c *gin.Context) {
	lat := strings.TrimSpace(c.Query("lat"))
	lon := strings.TrimSpace(c.Query("lon"))

	if lat == "" || lon == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "lat and lon are required",
		})
		return
	}

	endpoint := "https://nominatim.openstreetmap.org/reverse"
	params := url.Values{}
	params.Set("lat", lat)
	params.Set("lon", lon)
	params.Set("format", "jsonv2")
	params.Set("addressdetails", "1")
	params.Set("zoom", "12")

	req, err := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodGet,
		endpoint+"?"+params.Encode(),
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create reverse geo request"})
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "BlomivaSchool/1.0 contact@blomiva.com")

	res, err := h.client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "geo provider unavailable"})
		return
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		c.JSON(http.StatusBadGateway, gin.H{"error": "reverse geo provider failed"})
		return
	}

	var raw nominatimReverseItem
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "invalid reverse geo response"})
		return
	}

	place := normalizePlace(raw.DisplayName, raw.Lat, raw.Lon, raw.Address)

	c.JSON(http.StatusOK, gin.H{
		"data": place,
	})
}

func normalizePlace(displayName string, lat string, lon string, address map[string]string) Place {
	city := firstNonEmpty(
		address["city"],
		address["town"],
		address["village"],
		address["municipality"],
		address["suburb"],
		address["county"],
	)

	region := firstNonEmpty(
		address["state"],
		address["region"],
		address["province"],
	)

	district := firstNonEmpty(
		address["county"],
		address["district"],
		address["municipality"],
	)

	name := city
	if name == "" {
		name = firstDisplayPart(displayName)
	}

	return Place{
		Name:     name,
		City:     city,
		Town:     address["town"],
		Village:  address["village"],
		District: district,
		County:   address["county"],
		Region:   region,
		State:    address["state"],
		Country:  firstNonEmpty(address["country"], "Ghana"),
		Lat:      lat,
		Lon:      lon,
	}
}

func firstDisplayPart(value string) string {
	parts := strings.Split(value, ",")
	if len(parts) == 0 {
		return ""
	}

	return strings.TrimSpace(parts[0])
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		clean := strings.TrimSpace(value)
		if clean != "" {
			return clean
		}
	}
	return ""
}

func DebugURL(path string) {
	fmt.Println(path)
}
