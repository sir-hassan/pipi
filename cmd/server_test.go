package main

import (
	"github.com/go-kit/kit/log"
	"github.com/sir-hassan/pipi/backend"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestPiPiHandler(t *testing.T) {
	tt := []struct {
		name       string
		url        string
		want       string
		statusCode int
	}{
		{
			name:       "valid movie url",
			url:        "/movie/amazon/B08MZ3B1KC",
			want:       "{\"title\":\"Mercy Black [dt./OV]\",\"release_year\":2020,\"actors\":[\"Daniella Pineda\",\"Austin Amelio\",\"Elle LaMont\"],\"poster\":\"https://images-na.ssl-images-amazon.com/images/S/sgp-catalog-images/region_DE/lighthouse-LHE28437196-Full-Image_GalleryBackground-de-DE-1613717563794._SX1080_.jpg\",\"similar_ids\":[\"tail/B08KP\",\"tail/B08FM\",\"tail/B08X2\",\"tail/B08JG\",\"tail/B08KF\",\"tail/B08MS\",\"tail/B08LQ\",\"tail/B08P9\",\"tail/B08KR\",\"tail/B08LD\",\"tail/B08LT\",\"tail/B08R3\",\"tail/B08XT\",\"tail/B08HL\",\"tail/B08MD\",\"tail/B08HH\",\"tail/B08MD\",\"tail/B08QH\",\"tail/B085X\",\"tail/B08QH\",\"tail/B00MX\",\"tail/B07DH\",\"tail/B083X\",\"tail/B00VQ\",\"tail/B08L7\",\"tail/B07YM\",\"tail/B07WC\",\"tail/B01M1\",\"tail/B08MK\",\"tail/B018U\",\"tail/B074M\"]}",
			statusCode: http.StatusOK,
		},
		{
			name:       "valid movie url",
			url:        "/movie/amazon/B00K19SD8Q",
			want:       "{\"title\":\"Um Jeden Preis [dt./OV]\",\"release_year\":2013,\"actors\":[\"Dennis Quaid\",\"Zac Efron\",\"Kim Dickens\"],\"poster\":\"https://images-na.ssl-images-amazon.com/images/S/sgp-catalog-images/region_DE/universum-00664000-Full-Image_GalleryBackground-de-DE-1617099345129._SX1080_.jpg\",\"similar_ids\":[\"tail/B00IX\",\"tail/B08TB\",\"tail/B00N1\",\"tail/B00IM\",\"tail/B00TP\",\"tail/B01N3\",\"tail/B0172\",\"tail/B08P3\",\"tail/B00FY\",\"tail/B00HD\",\"tail/B0742\",\"tail/B00IL\",\"tail/B01GR\",\"tail/B00FC\",\"tail/B00FA\",\"tail/B00HD\",\"tail/B014X\",\"tail/B00IK\",\"tail/B00IP\",\"tail/B08VZ\",\"tail/B01DT\",\"tail/B07FT\",\"tail/B07FT\",\"tail/B08WK\",\"tail/B00I8\",\"tail/B00FF\",\"tail/B00I0\",\"tail/B087W\",\"tail/B07M8\",\"tail/B00FZ\",\"tail/B081K\",\"tail/B01N3\",\"tail/B00IU\",\"tail/B0811\",\"tail/B081Z\",\"tail/B00JZ\",\"tail/B01GF\",\"tail/B0747\",\"tail/B00HD\",\"tail/B00ZG\",\"tail/B07YP\",\"tail/B014J\",\"tail/B01KU\",\"tail/B00N9\",\"tail/B07S6\",\"tail/B01I1\",\"tail/B00J2\",\"tail/B00H3\",\"tail/B00LI\",\"tail/B07TB\",\"tail/B00N1\",\"tail/B075X\",\"tail/B01G4\",\"tail/B087Z\",\"tail/B00G0\",\"tail/B07ZG\",\"tail/B00H3\",\"tail/B0149\",\"tail/B07ZG\",\"tail/B00IG\",\"tail/B00FZ\",\"tail/B07L6\",\"tail/B00H3\"]}",
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tc.url, nil)
			responseRecorder := httptest.NewRecorder()

			logger := log.NewNopLogger()
			path, _ := os.Getwd()
			client := backend.NewFilesClient(path + "/../test_pages")

			handler := createHandler(logger, client)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}
			if strings.TrimSpace(strings.TrimSpace(responseRecorder.Body.String())) != tc.want {
				t.Errorf("want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}
