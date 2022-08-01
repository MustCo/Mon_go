package utils

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

type Gauge float64
type Counter int64

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type SysGather interface {
	Areas() (t, ID, value string)
}

type MetricsStorage map[string]*Metrics

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', 3, 64)
}

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

func (m *Metrics) String() string {
	s := fmt.Sprintf("ID:%s\ntype:%s\n", m.ID, m.MType)
	if m.Delta != nil {
		s += fmt.Sprintf("delta:%v\n", *m.Delta)
	}
	if m.Delta != nil {
		s += fmt.Sprintf("value:%v\n", *m.Value)
	}
	return s
}

func (m *Metrics) Areas() (t, ID, value string) {
	t = m.MType
	ID = m.ID
	if m.Delta != nil {
		value = fmt.Sprintf("%v", *m.Delta)
	}
	if m.Value != nil {
		value = fmt.Sprintf("%v", *m.Value)
	}
	return

}

func NewMetrics(id, mtype string) *Metrics {
	switch mtype {
	case "counter":
		return &Metrics{
			ID:    id,
			MType: mtype,
			Delta: new(int64),
		}
	case "gauge":
		return &Metrics{
			ID:    id,
			MType: mtype,
			Value: new(float64),
		}
	}
	return &Metrics{
		ID:    id,
		MType: mtype,
	}
}

func NewMetricsStorage() MetricsStorage {
	return MetricsStorage{}
}

func (m MetricsStorage) Poll() {
	g := "gauge"
	if _, ok := m["PollCount"]; !ok {
		m["PollCount"] = NewMetrics("PollCount", "counter")
	}
	*m["PollCount"].Delta += 1
	metrics := &runtime.MemStats{}
	runtime.ReadMemStats(metrics)
	mtrx := reflect.ValueOf(metrics).Elem()
	for i := 0; i < mtrx.NumField(); i++ {
		f := mtrx.Field(i)
		m[mtrx.Type().Field(i).Name] = NewMetrics(mtrx.Type().Field(i).Name, g)
		switch f.Kind() {
		case reflect.Int32, reflect.Int64:
			*m[mtrx.Type().Field(i).Name].Value = float64(f.Int())

		case reflect.Uint64, reflect.Uint32:
			*m[mtrx.Type().Field(i).Name].Value = float64(f.Uint())

		case reflect.Float32, reflect.Float64:
			*m[mtrx.Type().Field(i).Name].Value = f.Float()
		}

	}
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	m["RandomValue"] = NewMetrics("RandomValue", g)
	*m["RandomValue"].Value = r.Float64()
	log.Println("Poll metrics")
}
