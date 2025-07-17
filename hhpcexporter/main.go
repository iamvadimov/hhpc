package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed rootPage.html
var rootPageHTML []byte

var errInvalidResponse = errors.New("unexpected response from the server")

type hpcValues struct {
	TimeStamp string  `json:"timestamp"`
	Gapower   float64 `json:"gapower"`
	Grpower   float64 `json:"grpower"`
	Voltage   float64 `json:"voltage"`
	Gintens   float64 `json:"gintens"`
	Sm1       float64 `json:"sm1"`
	Sm2       float64 `json:"sm2"`
	Sm3       float64 `json:"sm3"`
}

func (hv *hpcValues) getHpcValues(client *http.Client, url string) error {
	response, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf(
			"%w: invalid status code: %s",
			errInvalidResponse,
			response.Status,
		)
		log.Println(err)
		return err
	}

	if err := json.NewDecoder(response.Body).Decode(hv); err != nil {
		log.Println("*!*", err)
		return err
	}
	return nil
}

type metrics struct {
	results *hpcValues
	up      float64
	expire  time.Time
	sync.RWMutex
}

func (m *metrics) getMetrics(client *http.Client, url string) *metrics {
	m.Lock()
	defer m.Unlock()

	if time.Now().Before(m.expire) {
		return m
	}

	m.up = 1
	if err := m.results.getHpcValues(client, url); err != nil {
		m.up = 0
		m.results.Gapower = 0
		m.results.Grpower = 0
		m.results.Voltage = 0
		m.results.Gintens = 0
		m.results.Sm1 = 0
		m.results.Sm2 = 0
		m.results.Sm3 = 0
	}

	m.expire = time.Now().Add(20 * time.Second)

	return m
}

func (m *metrics) status() float64 {
	m.RLock()
	defer m.RUnlock()

	return m.up
}

func (m *metrics) gapower() float64 {
	m.RLock()
	defer m.RUnlock()

	return m.results.Gapower
}

func (m *metrics) grpower() float64 {
	m.RLock()
	defer m.RUnlock()

	return m.results.Grpower
}

func (m *metrics) voltage() float64 {
	m.RLock()
	defer m.RUnlock()

	return m.results.Voltage
}

func (m *metrics) gintens() float64 {
	m.RLock()
	defer m.RUnlock()

	return m.results.Gintens
}

func (m *metrics) sm1() float64 {
	m.RLock()
	defer m.RUnlock()

	return m.results.Sm1
}

func (m *metrics) sm2() float64 {
	m.RLock()
	defer m.RUnlock()

	return m.results.Sm2
}

func (m *metrics) sm3() float64 {
	m.RLock()
	defer m.RUnlock()

	return m.results.Sm3
}

func newMux(url string) http.Handler {
	mux := http.NewServeMux()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	m := &metrics{
		results: &hpcValues{},
	}

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "global_active_power",
			Help:        "Household global minute-averaged active power.",
			ConstLabels: prometheus.Labels{"unit": "kilowatt"},
		},
		func() float64 {
			return m.getMetrics(client, url).gapower()
		},
	)

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "global_reactive_power",
			Help:        "Household global minute-averaged reactive power.",
			ConstLabels: prometheus.Labels{"unit": "kilowatt"},
		},
		func() float64 {
			return m.getMetrics(client, url).grpower()
		},
	)

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "voltage",
			Help:        "Minute-averaged voltage.",
			ConstLabels: prometheus.Labels{"unit": "volt"},
		},
		func() float64 {
			return m.getMetrics(client, url).voltage()
		},
	)

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "global_intensity",
			Help:        "Household global minute-averaged current intensity.",
			ConstLabels: prometheus.Labels{"unit": "ampere"},
		},
		func() float64 {
			return m.getMetrics(client, url).gintens()
		},
	)

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "sub_metering_1",
			Help:        "Energy sub-metering No. 1.",
			ConstLabels: prometheus.Labels{"unit": "watt-hour of active energy"},
		},
		func() float64 {
			return m.getMetrics(client, url).sm1()
		},
	)

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "sub_metering_2",
			Help:        "Energy sub-metering No. 2.",
			ConstLabels: prometheus.Labels{"unit": "watt-hour of active energy"},
		},
		func() float64 {
			return m.getMetrics(client, url).sm2()
		},
	)

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name:        "sub_metering_3",
			Help:        "Energy sub-metering No. 3.",
			ConstLabels: prometheus.Labels{"unit": "watt-hour of active energy"},
		},
		func() float64 {
			return m.getMetrics(client, url).sm3()
		},
	)

	promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "hhpc_up",
			Help: "Hhpc Server Status.",
		},
		func() float64 {
			return m.getMetrics(client, url).status()
		},
	)

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write(rootPageHTML)
	})

	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func main() {
	hhpcURL := os.Getenv("HHPC_SERVER_URL")
	// hhpcURL := "http://localhost:4000/api/v1/gethpc"
	s := &http.Server{
		Addr:         ":3030",
		Handler:      newMux(hhpcURL),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := s.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
