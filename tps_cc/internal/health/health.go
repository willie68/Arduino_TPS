package health

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
	"github.com/willie68/tps_cc/internal/config"
	log "github.com/willie68/tps_cc/internal/logging"
)

var myhealthy bool

/*
This is the readiness check you will have to provide.
*/
func check(tracer opentracing.Tracer) (bool, string) {
	cmd := exec.Command(
		config.Get().ArduinoCli,
		"version",
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	message := fmt.Sprintf("arduino-cli with exit code: %d, \r\n%s\r\n%s", cmd.ProcessState.ExitCode(), outStr, errStr)
	log.Logger.Info(message)
	if err != nil {
		return false, message
	}
	return true, message
}

//##### template internal functions for processing the healthchecks #####
var healthmessage string
var healthy bool
var lastChecked time.Time
var period int

// CheckConfig configuration for the healthcheck system
type CheckConfig struct {
	Period int
}

// Msg a health message
type Msg struct {
	Message   string `json:"message"`
	LastCheck string `json:"lastCheck,omitempty"`
}

// InitHealthSystem initialise the complete health system
func InitHealthSystem(config CheckConfig, tracer opentracing.Tracer) {
	period = config.Period
	log.Logger.Infof("healthcheck starting with period: %d seconds", period)
	healthmessage = "service starting"
	healthy = false
	doCheck(tracer)
	go func() {
		background := time.NewTicker(time.Second * time.Duration(period))
		for _ = range background.C {
			doCheck(tracer)
		}
	}()
}

/*
internal function to process the health check
*/
func doCheck(tracer opentracing.Tracer) {
	var msg string
	healthy, msg = check(tracer)
	if !healthy {
		healthmessage = msg
	} else {
		healthmessage = ""
	}
	lastChecked = time.Now()
}

/*
Routes getting all routes for the health endpoint
*/
func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/livez", GetHealthyEndpoint)
	router.Get("/readyz", GetReadinessEndpoint)
	return router
}

/*
GetHealthyEndpoint liveness probe
*/
func GetHealthyEndpoint(response http.ResponseWriter, req *http.Request) {
	render.Status(req, http.StatusOK)
	render.JSON(response, req, Msg{
		Message: fmt.Sprintf("service started"),
	})
}

/*
GetReadinessEndpoint is this service ready for taking requests, e.g. formaly known as health checks
*/
func GetReadinessEndpoint(response http.ResponseWriter, req *http.Request) {
	t := time.Now()
	if t.Sub(lastChecked) > (time.Second * time.Duration(2*period)) {
		healthy = false
		healthmessage = "Healthcheck not running"
	}
	if healthy {
		render.Status(req, http.StatusOK)
		render.JSON(response, req, Msg{
			Message:   "service up and running",
			LastCheck: lastChecked.String(),
		})
	} else {
		render.Status(req, http.StatusServiceUnavailable)
		render.JSON(response, req, Msg{
			Message:   fmt.Sprintf("service is unavailable: %s", healthmessage),
			LastCheck: lastChecked.String(),
		})
	}
}

/*
sendMessage sending a span message to tracer
*/
func sendMessage(tracer opentracing.Tracer, message string) {
	span := tracer.StartSpan("say-hello")
	println(message)
	span.Finish()
}
