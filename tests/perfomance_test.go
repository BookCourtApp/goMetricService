package tests

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/storage"
)

const (
	host = "localhost:5078"
)

var (
	testJson = storage.Metric{
		TimeStamp:   time.Now().Format(time.DateTime),
		IsApp:       222,
		IsAuth:      222,
		IsNew:       5,
		ResWidth:    1488,
		ResHeight:   1337,
		UserAgent:   "testing agent",
		UserID:      "12321",
		SessionID:   "12345324",
		DeviceType:  "Andoird",
		Reffer:      "https://google.com/",
		Stage:       "Statistic",
		Action:      "Close",
		ExtraKeys:   []string{"first", "second"},
		ExtraValues: []string{"first", "second"},
	}
)

func TestPerfomance(t *testing.T) {
	urlQuery := url.URL{
		Scheme: "http",
		Host:   host,
	}

	client := httpexpect.Default(t, urlQuery.String())

	t.Run("Connection", func(t *testing.T) {
		if !checkConnection(host) {
			t.Fatal("can't dial a connection with host, is connection up ?")
		}
	})
	requests := 1
	for i := 0; i < requests; i++ {
		t.Run(fmt.Sprintf("Request: %d/%d", i+1, requests), func(t *testing.T) {
			_ = client.POST("/test").
				WithJSON(testJson).
				Expect().Status(http.StatusCreated)

		})
	}
}

func checkConnection(host string) bool {
	timeout := 1 * time.Second
	_, err := net.DialTimeout("tcp", host, timeout)
	return err == nil
}
