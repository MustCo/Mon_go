package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Snakder/Mon_go/internal/server/db"
	"github.com/Snakder/Mon_go/internal/utils"
	"github.com/labstack/echo/v4"
)

func TestUpdateHandler_postMetric(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	testMetrics := utils.NewMetricsStorage()
	testMetrics["TestGauge"] = utils.NewMetrics("TestGauge", "gauge")
	*testMetrics["TestGauge"].Value = 123.124
	testMetrics["TestCounter"] = utils.NewMetrics("TestCounter", "counter")
	*testMetrics["TestCounter"].Delta = 123

	type want struct {
		db  *db.DB
		sc  int
		err error
	}
	tests := []struct {
		name    string
		handler *UpdateHandler
		args    args
		want    want
	}{
		{
			name:    "Invalid Path",
			handler: NewUpdateHandler(),
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodGet, "/status", nil)},
			want:    want{db: db.New(), sc: http.StatusMethodNotAllowed, err: echo.NewHTTPError(http.StatusNotImplemented, "invalid type")}},
		{
			name:    "Post_gauge",
			handler: NewUpdateHandler(),
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodPost, "/update/gauge/TestGauge/123.124", nil)},
			want:    want{db: db.New(), sc: http.StatusOK, err: nil},
		},
		{
			name:    "Post_counter",
			handler: NewUpdateHandler(),
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest(http.MethodPost, "/update/counter/TestCounter/123", nil)},
			want:    want{db: db.New(), sc: http.StatusOK, err: nil},
		},
	}
	tests[1].want.db.Metrics["TestGauge"] = testMetrics["TestGauge"]
	tests[2].want.db.Metrics["TestCounter"] = testMetrics["TestCounter"]
	e := echo.New()
	updater := NewUpdateHandler()
	e.POST("/update/:type/:name/:value", updater.postMetric)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e.NewContext(tt.args.r, tt.args.w)
		})
	}
}
