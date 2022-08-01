package agent

import (
	"testing"

	"github.com/Snakder/Mon_go/internal/utils"
)

func TestAPIAgent_SendMetric(t *testing.T) {
	type args struct {
		name string
		g    utils.Gauge
		c    utils.Counter
	}
	tests := []struct {
		name    string
		c       *APIAgent
		args    args
		wantErr bool
	}{{
		name:    "Test counter 1",
		c:       &APIAgent{config: &utils.Config{Address: "127.0.0.1:8080", ReportInterval: 10, PollInterval: 2}},
		args:    args{name: "Mymetric", c: 0, g: 125.453},
		wantErr: true,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
