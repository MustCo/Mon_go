package agent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MustCo/Mon_go/internal/utils"
	"github.com/go-resty/resty/v2"
)

type APIAgent struct {
	config *utils.Config
	client *resty.Client
}

func New(config *utils.Config) *APIAgent {
	return &APIAgent{config: config}
}

func (c *APIAgent) Report(ms utils.MetricsStorage) error {
	for _, v := range ms {
		err := c.sendJSON(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *APIAgent) sendJSON(m *utils.Metrics) error {
	resp, err := c.client.R().
		SetBody(*m).
		SetPathParams(map[string]string{
			"host": c.config.Address,
		}).
		Post("http://{host}/update/")
	log.Printf("Send JSON:\n%v", m)
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		fmt.Println("  Status Code:", resp.StatusCode())
		fmt.Println("  Status     :", resp.Status())
		fmt.Println("  Proto      :", resp.Proto())
		fmt.Println("  Time       :", resp.Time())
		fmt.Println("  Received At:", resp.ReceivedAt())
		fmt.Println("  Body       :\n", resp)
		return errors.New("invalid status code")
	}
	return nil

}

func (c *APIAgent) Start(ctx context.Context) error {
	m := utils.NewMetricsStorage()
	c.client = resty.New()
	c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	reports := time.NewTicker(c.config.ReportInterval)
	polls := time.NewTicker(c.config.PollInterval)
	for {
		select {
		case <-reports.C:
			err := c.Report(m)
			if err != nil {
				log.Println("Error", err)
			}
		case <-polls.C:
			m.Poll()
		case <-ctx.Done():
			log.Println("Exit by context")
			return nil
		}
	}

}
