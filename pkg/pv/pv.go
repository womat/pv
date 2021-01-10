package pv

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/womat/debug"
)

const httpRequestTimeout = 10 * time.Second

type Measurements struct {
	sync.RWMutex
	Timestamp time.Time
	Power     float64
	Energy    float64
	config    struct {
		meterURL string
	}
}

type meterURLBody struct {
	Timestamp time.Time `json:"Time"`
	Runtime   float64   `json:"Runtime"`
	Measurand struct {
		E float64 `json:"e"`
		P float64 `json:"p"`
	} `json:"Measurand"`
}

func New() *Measurements {
	return &Measurements{}
}

func (m *Measurements) SetMeterURL(url string) {
	m.config.meterURL = url
}

func (m *Measurements) Read() (err error) {
	start := time.Now()
	if err = m.readInverterMeter(); err != nil {
		return
	}
	debug.DebugLog.Printf("runtime to request data: %vs", time.Since(start).Seconds())

	return
}

func (m *Measurements) readInverterMeter() (err error) {
	var r meterURLBody

	if err = read(m.config.meterURL, &r); err != nil {
		return
	}

	m.Lock()
	defer m.Unlock()

	m.Power = r.Measurand.P
	m.Energy = r.Measurand.E
	m.Timestamp = r.Timestamp
	return
}

func read(url string, data interface{}) (err error) {
	done := make(chan bool, 1)
	go func() {
		// ensures that data is sent to the channel when the function is terminated
		defer func() {
			select {
			case done <- true:
			default:
			}
			close(done)
		}()

		debug.TraceLog.Printf("performing http get: %v\n", url)

		var resp *http.Response
		if resp, err = http.Get(url); err != nil {
			return
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if err = json.Unmarshal(bodyBytes, data); err != nil {
			return
		}
	}()

	// wait for API Data
	select {
	case <-done:
	case <-time.After(httpRequestTimeout):
		err = errors.New("timeout during receive data")
	}

	if err != nil {
		debug.ErrorLog.Println(err)
		return
	}
	return
}
