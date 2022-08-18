package db

import (
	"errors"
	"strings"
	"sync"

	"github.com/Snakder/Mon_go/internal/utils"
)

func New() *DB {
	db := new(DB)
	db.Metrics = utils.NewMetricsStorage()
	db.mut = new(sync.Mutex)
	return db
}

type DB struct {
	mut     *sync.Mutex
	Metrics map[string]utils.SysGather
}

func (db *DB) Set(name, t, val string) error {
	m, err := utils.NewMetrics(name, t, val)
	if err != nil {
		return err
	}
	db.mut.Lock()
	db.Metrics[name] = m
	db.mut.Unlock()
	return nil

}

func (db *DB) Get(t, name string) (utils.SysGather, error) {
	db.mut.Lock()
	defer db.mut.Unlock()
	if m, ok := db.Metrics[name]; ok {
		_, mtype, _ := m.Areas()
		if mtype == strings.ToLower(t) {
			return m, nil
		} else {
			return nil, errors.New("invalid type")
		}
	}
	return nil, errors.New("unknown metric")

}

func (db *DB) GetAll() map[string]utils.SysGather {
	return db.Metrics
}
