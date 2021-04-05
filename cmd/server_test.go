package main

import (
	"github.com/go-kit/kit/log"
	"github.com/sir-hassan/pipi/backend"
	"net/http"
	"net/http/httptest"
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
			url:        "/movie/amazon/B07NJ7X55C",
			want:       "{\"title\":\"The Prodigy [dt./OV]\",\"release_year\":2019,\"actors\":[\"Taylor Schilling\",\"Jackson Robert Scott\",\"Colm Feore\"],\"poster\":\"https://images-na.ssl-images-amazon.com/images/I/81KFZJNk3eL._SX300_.jpg\",\"similar_ids\":[\"B08FMQTK65\",\"B01MUWNDPR\",\"B07ZY6PXX2\",\"B07RQ89RP8\",\"B08YWJ96M6\",\"B08SLHB8SB\",\"B08X2T9292\",\"B08QH2XW8M\",\"B08Y81994Y\",\"B08G597V3H\",\"B08KRZD96B\",\"B08KFDFKZ9\",\"B08LSWMDLN\",\"B08HLBMDQV\",\"B08KPLD3WP\",\"B08MZ3B1KC\",\"B07R7TZ2XZ\",\"B08R32M7DQ\",\"B08Q3T93JZ\",\"B08XTVDXJR\"]}",
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tc.url, nil)
			responseRecorder := httptest.NewRecorder()

			logger := log.NewNopLogger()
			client := backend.NewHttpClient(&http.Client{})
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
