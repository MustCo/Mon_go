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

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', 3, 64)
}

func (c Counter) String() string {
	return strconv.FormatInt(int64(c), 10)
}

type SysGather interface {
	Areas() (id, mtype, value string)
	Update(value string) error
}
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type MetricsStorage map[string]SysGather

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

func (m *Metrics) Areas() (id, mtype, value string) {
	mtype = m.MType
	id = m.ID
	if m.Delta != nil {
		value = fmt.Sprintf("%d", *m.Delta)
	}
	if m.Value != nil {
		value = fmt.Sprintf("%.3f", *m.Value)
	}
	return

}

func (m *Metrics) Update(value string) error {
	switch m.MType {
	case "counter":
		log.Print("Create Counter", m)
		d, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		m.Delta = &d
		log.Print(m)
	case "gauge":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		m.Value = &v
	}
	return nil
}

func NewMetrics(id, mtype, value string) (*Metrics, error) {
	m := new(Metrics)
	m.ID = id
	m.MType = mtype
	log.Println("Metric")
	log.Print(m)
	m.Update(value)

	return m, nil
}

func NewMetricsStorage() MetricsStorage {
	return MetricsStorage{}
}

func Poll(poll_count string) map[string]SysGather {
	var v string
	m := NewMetricsStorage()
	m["PollCount"], _ = NewMetrics("PollCount", "counter", poll_count)
	t := reflect.ValueOf(1.1).Type()
	metrics := &runtime.MemStats{}
	runtime.ReadMemStats(metrics)
	mtrx := reflect.ValueOf(metrics).Elem()
	for i := 0; i < mtrx.NumField(); i++ {
		f := mtrx.Field(i)
		if f.CanConvert(t) {
			v = fmt.Sprintf("%.3f", f.Convert(t).Float())
			m[mtrx.Type().Field(i).Name], _ = NewMetrics(mtrx.Type().Field(i).Name, "gauge", v)
		}
	}
	seed := rand.NewSource(time.Now().UnixNano())
	r := fmt.Sprintf("%.3f", rand.New(seed).Float64())
	m["RandomValue"], _ = NewMetrics("RandomValue", "gauge", r)
	log.Println("Poll metrics")
	return m
}
