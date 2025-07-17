package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Data map[string][]string

type Hpc struct {
	TimeStamp string  `json:"timestamp"`
	Gapower   float64 `json:"gapower"`
	Grpower   float64 `json:"grpower"`
	Voltage   float64 `json:"voltage"`
	Gintens   float64 `json:"gintens"`
	Sm1       float64 `json:"sm1"`
	Sm2       float64 `json:"sm2"`
	Sm3       float64 `json:"sm3"`
}

type Reals [7]float64

func (h Hpc) String() string {
	return fmt.Sprintf("TimeStamp: \"%s\", Gapower: %0.3f, Grpower: %0.3f, "+
		"Voltage: %0.3f, Gintens: %0.3f, Sm1: %0.3f, Sm2: %0.3f, Sm3: %0.3f", h.TimeStamp,
		h.Gapower, h.Grpower, h.Voltage, h.Gintens, h.Sm1, h.Sm2, h.Sm3)
}

func convert(values []string) Reals {
	var result Reals
	for i, str := range values[:7] {
		if fVal, err := strconv.ParseFloat(str, 64); err == nil {
			result[i] = fVal
		} else {
			fmt.Fprintln(os.Stderr, "Cannot convert string to float:", err)
			return Reals{}
		}
	}
	return result
}

func getHpc(now time.Time, data Data) (Hpc, error) {
	const format string = "2006-01-02 15:04:00" // 2025-05-08 15:04:00
	ts := now.Format(format)
	var result Hpc
	if values, ok := data[ts]; !ok {
		err := fmt.Errorf("TimeStamp %s not found", ts)
		result = Hpc{TimeStamp: ts} // {TimeStamp:"2006-01-02 15:04:00", Gapower:0, Grpower:0, Voltage:0, Gintens:0, Sm1:0, Sm2:0, Sm3:0}
		return result, err
	} else {
		fValue := convert(values)
		// fmt.Printf("%s: %v\n", ts, values)
		result = Hpc{
			TimeStamp: ts,
			Gapower:   fValue[0],
			Grpower:   fValue[1],
			Voltage:   fValue[2],
			Gintens:   fValue[3],
			Sm1:       fValue[4],
			Sm2:       fValue[5],
			Sm3:       fValue[6],
		}
	}
	return result, nil
}

func homePage(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Household power consumption REST API Server")
}

func getValue(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Processing request from %s: User-Agent: %s\n", r.RemoteAddr, r.UserAgent())
	now := time.Now()
	result, err := getHpc(now, data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot get current date:", err)
		http.Error(w, "Cannot get current date", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot write response:", err)
	}
}

var data Data

func main() {
	now := time.Now()
	data = make(Data) // "2025-05-08 15:04:00": [1.756 0.304 235.670 7.400 1.000 2.000 17.000]
	createData(now, data)

	http.HandleFunc("/", homePage) // http://localhost:4000/
	// Response: Household power consumption REST API Server
	http.HandleFunc("/api/v1/gethpc", getValue) // http://localhost:4000/api/v1/gethpc
	// Response: {"timestamp":"2025-05-12 12:12:00","gapower":0.426,"grpower":0.318,"voltage":236.16,"gintens":2.2,"sm1":0,"sm2":2,"smg3":0}

	fmt.Println("Starting API server on port 4000...")
	if err := http.ListenAndServe(":4000", nil); err != nil {
		// Значение nil во втором аргументе означает, что запросы будут
		// обрабатываться с использованием функций, заданных при помощи HandleFunc.
		fmt.Fprintln(os.Stderr, "Cannot start API server:", err)
		os.Exit(1)
	}
}
