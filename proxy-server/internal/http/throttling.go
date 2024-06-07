package http

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

const (
	randomNumberEndpointMaxRequestsPerTimePeriod        = 3
	internalServerErrorEndpointMaxRequestsPerTimePeriod = 3
)

const (
	randomNumberEndpointMaxRequestsTimePeriod        = 5
	internalServerErrorEndpointMaxRequestsTimePeriod = 5
)

const (
	throttlingEngineVisitorsCleanUpMinutes   = 3
	throttlingEngineCleanUpTimePeriodMinutes = 2
)

type ThrottlingEngine struct {
	Visitors      map[string]EndpointThrottlingLimiter
	VisitorsMutex *sync.Mutex
}

type EndpointThrottlingConfig struct {
	Endpoint      string
	Method        string
	MaxRequests   int
	TimePeriodSec int
}

type EndpointThrottlingLimiter struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
	Config   EndpointThrottlingConfig
}

func NewThrottlingEngine() ThrottlingEngine {
	return ThrottlingEngine{
		Visitors:      map[string]EndpointThrottlingLimiter{},
		VisitorsMutex: &sync.Mutex{},
	}
}

func NewEndpointThrottlingLimiter(req http.Request) EndpointThrottlingLimiter {
	etc := requestThrottlingConfig(req)

	r := rate.Every(time.Second * time.Duration(etc.TimePeriodSec))
	b := etc.MaxRequests

	return EndpointThrottlingLimiter{
		Limiter:  rate.NewLimiter(r, b),
		Config:   etc,
		LastSeen: time.Now(),
	}
}

func (te *ThrottlingEngine) StartThrottlingEngineCleanUp() {
	go te.cleanupThrottlingEngine()
}

func (te *ThrottlingEngine) CanAllowRequest(req http.Request) bool {
	id := requestId(req)

	etl, ok := te.Visitors[id]

	if !ok {
		te.VisitorsMutex.Lock()
		defer te.VisitorsMutex.Unlock()

		etl = NewEndpointThrottlingLimiter(req)
	}

	etl.LastSeen = time.Now()

	te.Visitors[id] = etl

	return etl.Limiter.Allow()
}

func (te *ThrottlingEngine) cleanupThrottlingEngine() {
	for {
		time.Sleep(throttlingEngineCleanUpTimePeriodMinutes * time.Minute)

		te.VisitorsMutex.Lock()

		for ip, v := range te.Visitors {
			if time.Since(v.LastSeen) > throttlingEngineVisitorsCleanUpMinutes*time.Minute {
				delete(te.Visitors, ip)
			}
		}

		te.VisitorsMutex.Unlock()
	}
}

func requestId(req http.Request) string {
	ipextr := echo.ExtractIPDirect()
	ip := ipextr(&req)

	id := fmt.Sprintf("%s%s", ip, req.Method)

	return id
}

func requestThrottlingConfig(req http.Request) EndpointThrottlingConfig {
	var etc EndpointThrottlingConfig

	switch req.URL.Path {
	case internalServerErrorRoute:
		etc = createInternalServerErrorThrottlingConfig()
	case randomNumberRoute:
		etc = createRandomNumberEndpointThrottlingConfig()
	}

	return etc
}

func createRandomNumberEndpointThrottlingConfig() EndpointThrottlingConfig {
	return EndpointThrottlingConfig{
		Endpoint:      randomNumberRoute,
		Method:        http.MethodGet,
		MaxRequests:   randomNumberEndpointMaxRequestsPerTimePeriod,
		TimePeriodSec: randomNumberEndpointMaxRequestsTimePeriod,
	}
}

func createInternalServerErrorThrottlingConfig() EndpointThrottlingConfig {
	return EndpointThrottlingConfig{
		Endpoint:      internalServerErrorRoute,
		Method:        http.MethodGet,
		MaxRequests:   internalServerErrorEndpointMaxRequestsPerTimePeriod,
		TimePeriodSec: internalServerErrorEndpointMaxRequestsTimePeriod,
	}
}
