package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Snakder/Mon_go/internal/utils"
	"github.com/labstack/echo/v4"
)

type APIServer struct {
	config *utils.Config
	e      *echo.Echo
}

func New(config *utils.Config) *APIServer {
	updater := NewUpdateHandler()
	e := echo.New()
	e.GET("/", updater.getAllMetrics)
	e.GET("/value/:type/:name", updater.getMetric)
	e.POST("/update/:type/:name/:value", updater.postMetric)
	e.POST("/update/", updater.updateJSON)
	e.POST("/value/", updater.getJSON)
	return &APIServer{config: config, e: e}
}

type UpdateHandler struct {
	DB Storage
}

func NewUpdateHandler() *UpdateHandler {
	updater := new(UpdateHandler)
	updater.DB = NewDB()
	return updater
}

func (handler *UpdateHandler) getAllMetrics(c echo.Context) error {
	metrics := handler.DB.GetAll()
	answer := ""
	resp := c.Response()
	resp.Header().Set("Content-Type", "text/plain")
	for _, m := range metrics {
		answer += m.ID
		switch m.MType {
		case "counter":
			answer += fmt.Sprintf(" - %d\n", *m.Delta)
		case "gauge":
			answer += fmt.Sprintf(" - %.3f\n", *m.Value)
		}

	}
	return c.HTML(http.StatusOK, answer)
}

func (handler *UpdateHandler) getMetric(c echo.Context) error {
	var v string
	t := c.Param("type")
	name := c.Param("name")
	log.Printf("Get Metric type:%s name:%s", t, name)
	val, err := handler.DB.Get(t, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	resp := c.Response()
	resp.Header().Set("Content-Type", "text/plain")
	if val.Delta != nil {
		v = strconv.FormatInt(*val.Delta, 10)
	} else {
		v = strconv.FormatFloat(*val.Value, 'f', 3, 64)
	}
	return c.HTML(http.StatusOK, v)
}
func (handler *UpdateHandler) postMetric(c echo.Context) error {

	t := c.Param("type")
	name := c.Param("name")
	val := c.Param("value")
	log.Printf("Post Metric type:%s name:%s value:%s", t, name, val)
	err := handler.DB.Set(t, name, val)
	switch {
	case err == nil:
		return nil
	case err.Error() == "invalid type":
		return echo.NewHTTPError(http.StatusNotImplemented, err.Error())
	default:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}

}

func (handler *UpdateHandler) updateJSON(c echo.Context) error {
	m := new(utils.Metrics)
	m.Delta = new(int64)
	m.Value = new(float64)
	err := c.Bind(m)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	switch {
	}
	err = handler.DB.SetMetrica(m)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
func (handler *UpdateHandler) getJSON(c echo.Context) error {
	m := new(utils.Metrics)
	m.Delta = new(int64)
	m.Value = new(float64)
	err := c.Bind(m)
	if err != nil {
		return err
	}
	log.Print(m)
	metrics, err := handler.DB.Get(m.MType, m.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	resp := c.Response()
	resp.Header().Set("Content-Type", "application/json")

	return c.JSONPretty(http.StatusOK, metrics, "   ")
}

func (s *APIServer) Start(ctx context.Context) error {
	return s.e.Start(s.config.Address)
}
