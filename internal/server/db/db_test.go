package db

import (
	"errors"
	"sync"
	"testing"

	"github.com/Snakder/Mon_go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestDB_Set(t *testing.T) {
	type args struct {
		t    string
		name string
		val  string
	}
	type want struct {
		db  *DB
		err error
	}
	tests := []struct {
		name string
		db   *DB
		args args
		want want
	}{
		{
			name: "TestGauge",
			db:   New(),
			args: args{t: "gauge", name: "Mymetric", val: "1.329184"},
			want: want{
				err: nil,
			},
		},
		{
			name: "TestInvalidGauge",
			db:   New(),
			args: args{t: "gauge", name: "Mymetric", val: "1,FSD29184"},
			want: want{
				err: errors.New("error"),
			},
		},
		{
			name: "TestCounter",
			db:   New(),
			args: args{t: "counter", name: "Mymetric", val: "14563"},
			want: want{
				err: nil,
			},
		},
		{
			name: "TestInvalidCounter",
			db:   New(),
			args: args{t: "counter", name: "Mymetric", val: "1.329184"},
			want: want{
				err: errors.New("error"),
			},
		},
		{
			name: "TestInvalidType",
			db:   New(),
			args: args{t: "Mytype", name: "Mymetric", val: "1sdfgsd4"},
			want: want{
				err: errors.New("unknown metric"),
			},
		}}
	for i, test := range tests {
		db := utils.NewMetricsStorage()
		m, _ := utils.NewMetrics(test.args.name, test.args.t, test.args.val)
		db[test.args.name] = m
		tests[i].want.db = &DB{Metrics: db, mut: &sync.Mutex{}}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("%s %s %s", tt.args.name, tt.args.t, tt.args.val)
			err := tt.db.Set(tt.args.name, tt.args.t, tt.args.val)
			if err != nil {
				if tt.want.err != nil {
					return
				}
				assert.Equal(t, err.Error(), tt.want.err.Error())
			}
			assert.Equal(t, tt.want.db, tt.db)
		})
	}
}

func TestDB_Get(t *testing.T) {
	initDB := New()
	initDB.Set("TestGauge", "gauge", "123.123")
	initDB.Set("TestCounter", "counter", "123")
	initDB.Set("Mygauge", "gauge", "14.563")
	initDB.Set("Mycounter", "counter", "14563")
	type args struct {
		t    string
		name string
	}
	type want struct {
		res utils.SysGather
		err error
	}
	tests := []struct {
		name string
		db   *DB
		args args
		want want
	}{
		{
			name: "TestGauge",
			db:   initDB,
			args: args{t: "gauge", name: "Mygauge"},
			want: want{res: initDB.Metrics["Mygauge"],
				err: nil,
			},
		},
		{
			name: "TestInvalidGauge",
			db:   initDB,
			args: args{t: "gauge", name: "Unknown"},
			want: want{res: nil,
				err: errors.New("unknown metric"),
			},
		},
		{
			name: "TestCounter",
			db:   initDB,
			args: args{t: "counter", name: "Mycounter"},
			want: want{res: initDB.Metrics["Mycounter"],
				err: nil,
			},
		},
		{
			name: "TestInvalidCounter",
			db:   initDB,
			args: args{t: "counter", name: "Unknown"},
			want: want{res: nil,
				err: errors.New("unknown metric"),
			},
		},
		{
			name: "TestInvalidType&Name",
			db:   initDB,
			args: args{t: "untype", name: "Mymetric"},
			want: want{res: nil,
				err: errors.New("unknown metric"),
			},
		},
		{
			name: "TestInvalidType",
			db:   initDB,
			args: args{t: "untype", name: "Mygauge"},
			want: want{res: nil,
				err: errors.New("invalid type"),
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.db.Get(tt.args.t, tt.args.name)
			if err != nil {
				assert.Equal(t, err.Error(), tt.want.err.Error())
				return
			}
			assert.Equal(t, tt.want.res, got)
		})
	}
}
