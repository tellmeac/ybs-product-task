package integration

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"os"
	"testing"
	"time"
)

type ServerTestSuite struct {
	suite.Suite
	ctx         context.Context
	cleanUpTest func()
	client      *resty.Client
}

func (s *ServerTestSuite) TestOrderPipeline() {
	r, err := s.client.R().
		SetBody(map[string]interface{}{
			"orders": nil,
		}).Post("/orders")
	require.NoError(s.T(), err)
	require.Equalf(s.T(), http.StatusBadRequest, r.StatusCode(),
		"For POST /orders with empty order list should be 400")

	r, err = s.client.R().
		SetBody(map[string]interface{}{
			"orders": []map[string]interface{}{
				{
					"weight":         25.5,
					"delivery_hours": nil,
					"cost":           500,
				},
			},
		}).Post("/orders")
	require.NoError(s.T(), err)
	require.Equalf(s.T(), http.StatusBadRequest, r.StatusCode(),
		"For POST /orders with empty delivery_hours field should be 400")

	var orders []struct {
		ID            int64    `json:"order_id"`
		Weight        float64  `json:"weight"`
		Region        int32    `json:"regions"`
		DeliveryHours []string `json:"delivery_hours"`
		Cost          int32    `json:"cost"`
	}

	r, err = s.client.R().
		SetBody(map[string]interface{}{
			"orders": []map[string]interface{}{
				{
					"weight":         25.5,
					"delivery_hours": []string{"12:00-14:00", "16:00-20:00"},
					"cost":           500,
				},
			},
		}).SetResult(&orders).Post("/orders")
	require.NoError(s.T(), err)
	require.Equalf(s.T(), http.StatusOK, r.StatusCode(),
		"Valid POST /orders")
	require.Len(s.T(), orders, 1)
	require.NotEmpty(s.T(), orders[0].ID)

	orders = nil
	r, err = s.client.R().SetResult(&orders).Get("/orders")
	require.NoError(s.T(), err)
	require.Equalf(s.T(), http.StatusOK, r.StatusCode(),
		"Valid GET /orders with some values")
	require.Len(s.T(), orders, 1)
	require.NotEmpty(s.T(), orders[0].ID)
}

func (s *ServerTestSuite) TestCourierPipeline() {
	r, err := s.client.R().
		SetBody(map[string]interface{}{
			"couriers": []map[string]interface{}{
				{
					"regions":       []int32{70, 77},
					"courier_type":  "BIKE",
					"working_hours": []string{"8:00-18:00"},
				},
			},
		}).Post("/couriers")
	require.NoError(s.T(), err)
	require.Equalf(s.T(), http.StatusOK, r.StatusCode(),
		"Valid POST /couriers request")

	var couriersResponse struct {
		Couriers []struct {
			ID           int64    `json:"courier_id"`
			Type         string   `json:"courier_type"`
			Regions      []int32  `json:"regions"`
			WorkingHours []string `json:"working_hours"`
		} `json:"couriers"`
	}

	r, err = s.client.R().SetResult(&couriersResponse).Get("/couriers")
	require.NoError(s.T(), err)
	require.Equalf(s.T(), http.StatusOK, r.StatusCode(),
		"Valid GET /couriers with some values")
	require.Len(s.T(), couriersResponse.Couriers, 1)
	require.NotEmpty(s.T(), couriersResponse.Couriers[0].ID)
}

func (s *ServerTestSuite) TestRateLimitMiddleware() {
	for i := 0; i < 10; i++ {
		r, err := s.client.R().Get("/couriers")
		require.NoError(s.T(), err)
		require.Equal(s.T(), http.StatusOK, r.StatusCode())
	}

	r, err := s.client.R().Get("/couriers")
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusTooManyRequests, r.StatusCode())

	r, err = s.client.R().Get("/orders")
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, r.StatusCode())
}

func (s *ServerTestSuite) SetupTest() {
	time.Sleep(time.Second)

	s.ctx, s.cleanUpTest = context.WithTimeout(context.Background(), time.Second)

	c := resty.New()
	c.SetBaseURL(os.Getenv("LAVKA_BASE_URL"))

	s.client = c
}

func (s *ServerTestSuite) TearDownTest() {
	s.cleanUpTest()
}

func TestBase(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}
