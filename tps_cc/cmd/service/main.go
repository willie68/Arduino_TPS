package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/willie68/tps_cc/internal/api"
	"github.com/willie68/tps_cc/internal/apiv1"
	"github.com/willie68/tps_cc/internal/auth"
	"github.com/willie68/tps_cc/internal/health"
	"github.com/willie68/tps_cc/internal/serror"
	"github.com/willie68/tps_cc/internal/utils/httputils"
	"github.com/willie68/tps_cc/pkg/web"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	config "github.com/willie68/tps_cc/internal/config"

	jaegercfg "github.com/uber/jaeger-client-go/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httptracer"
	"github.com/go-chi/render"
	"github.com/willie68/tps_cc/internal/crypt"
	log "github.com/willie68/tps_cc/internal/logging"

	flag "github.com/spf13/pflag"
)

/*
apVersion implementing api version for this service
*/
const servicename = "tpscc"

var port int
var sslport int
var serviceURL string
var apikey string
var ssl bool
var configFile string
var serviceConfig config.Config
var Tracer opentracing.Tracer
var sslsrv *http.Server
var srv *http.Server

func init() {
	// variables for parameter override
	ssl = false
	log.Logger.Info("init service")
	flag.IntVarP(&port, "port", "p", 0, "port of the http server.")
	flag.IntVarP(&sslport, "sslport", "t", 0, "port of the https server.")
	flag.StringVarP(&configFile, "config", "c", config.File, "this is the path and filename to the config file")
	flag.StringVarP(&serviceURL, "serviceURL", "u", "", "service url from outside")
}

func apiRoutes() (*chi.Mux, error) {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		//middleware.DefaultCompress,
		middleware.Recoverer,
		cors.Handler(cors.Options{
			// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-mcs-username", "X-mcs-password", "X-mcs-profile"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
		httptracer.Tracer(Tracer, httptracer.Config{
			ServiceName:    servicename,
			ServiceVersion: "V" + apiv1.ApiVersion,
			SampleRate:     1,
			SkipFunc: func(r *http.Request) bool {
				return false
				//return r.URL.Path == "/livez"
			},
			Tags: map[string]interface{}{
				"_dd.measured": 1, // datadog, turn on metrics for http.request stats
				// "_dd1.sr.eausr": 1, // datadog, event sample rate
			},
		}),
	)
	if serviceConfig.Metrics.Enable {
		router.Use(
			api.MetricsHandler(api.MetricsConfig{
				SkipFunc: func(r *http.Request) bool {
					return false
				},
			}),
		)
	}
	if serviceConfig.Apikey {
		router.Use(
			api.SysAPIHandler(api.SysAPIConfig{
				Apikey: apikey,
				SkipFunc: func(r *http.Request) bool {
					path := strings.TrimSuffix(r.URL.Path, "/")
					if strings.HasSuffix(path, "/livez") {
						return true
					}
					if strings.HasSuffix(path, "/readyz") {
						return true
					}
					if strings.HasSuffix(path, "/metrics") {
						return true
					}
					if strings.HasPrefix(path, "/client") {
						return true
					}
					return false
				},
			}),
		)
	}
	// jwt is activated, register the Authenticator and Validator
	if strings.EqualFold(serviceConfig.Auth.Type, "jwt") {
		jwtConfig, err := auth.ParseJWTConfig(serviceConfig.Auth)
		if err != nil {
			return router, err
		}
		log.Logger.Infof("jwt config: %v", jwtConfig)
		jwtAuth := auth.JWTAuth{
			Config: jwtConfig,
		}
		router.Use(
			auth.Verifier(&jwtAuth),
			auth.Authenticator,
		)
	}

	// building the routes
	router.Route("/", func(r chi.Router) {
		r.Mount(apiv1.GenerateRoutes())
		r.Mount("/", health.Routes())
		if serviceConfig.Metrics.Enable {
			r.Mount("/metrics", promhttp.Handler())
		}
	})
	// adding a file server with web client asserts
	httputils.FileServer(router, "/client", http.FS(web.WebClientAssets))
	return router, nil
}

func healthRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		//middleware.DefaultCompress,
		middleware.Recoverer,
		httptracer.Tracer(Tracer, httptracer.Config{
			ServiceName:    servicename,
			ServiceVersion: "V" + apiv1.ApiVersion,
			SampleRate:     1,
			SkipFunc: func(r *http.Request) bool {
				return false
			},
			Tags: map[string]interface{}{
				"_dd.measured": 1, // datadog, turn on metrics for http.request stats
				// "_dd1.sr.eausr": 1, // datadog, event sample rate
			},
		}),
	)
	if serviceConfig.Metrics.Enable {
		router.Use(
			api.MetricsHandler(api.MetricsConfig{
				SkipFunc: func(r *http.Request) bool {
					return false
				},
			}),
		)
	}

	router.Route("/",
		func(r chi.Router) {
			r.Mount("/", health.Routes())
			if serviceConfig.Metrics.Enable {
				r.Mount("/metrics", promhttp.Handler())
			}
		})
	return router
}

// @title GoMicro service API
// @version 1.0
// @description The GoMicro service is a template for microservices written in go.
// @BasePath /api/v1
// @securityDefinitions.apikey api_key
// @in header
// @name apikey
func main() {
	configFolder, err := config.GetDefaultConfigFolder()
	if err != nil {
		panic("can't get config folder")
	}

	flag.Parse()

	log.Logger.Infof("starting server, config folder: %s", configFolder)
	defer log.Logger.Close()
	serror.Service = servicename
	if configFile == "" {
		configFolder, err := config.GetDefaultConfigFolder()
		if err != nil {
			log.Logger.Alertf("can't load config file: %s", err.Error())
			os.Exit(1)
		}
		configFolder = fmt.Sprintf("%s/service/", configFolder)
		err = os.MkdirAll(configFolder, os.ModePerm)
		if err != nil {
			log.Logger.Alertf("can't load config file: %s", err.Error())
			os.Exit(1)
		}
		configFile = configFolder + "/service.yaml"
	}
	config.File = configFile
	// autorestart starts here...
	if err := config.Load(); err != nil {
		log.Logger.Alertf("can't load config file: %s", err.Error())
		os.Exit(1)
	}

	serviceConfig = config.Get()
	initConfig()
	initLogging()

	log.Logger.Info("service is starting")

	var closer io.Closer
	Tracer, closer = initJaeger(servicename, serviceConfig.OpenTracing)
	opentracing.SetGlobalTracer(Tracer)
	defer closer.Close()

	healthCheckConfig := health.CheckConfig(serviceConfig.HealthCheck)

	health.InitHealthSystem(healthCheckConfig, Tracer)

	if serviceConfig.Sslport > 0 {
		ssl = true
		log.Logger.Info("ssl active")
	}

	apikey = getApikey()
	apiv1.APIKey = apikey
	if config.Get().Apikey {
		log.Logger.Infof("apikey: %s", apikey)
	}

	log.Logger.Infof("ssl: %t", ssl)
	log.Logger.Infof("serviceURL: %s", serviceConfig.ServiceURL)
	log.Logger.Infof("%s api routes", servicename)
	router, err := apiRoutes()
	if err != nil {
		log.Logger.Alertf("could not create api routes. %s", err.Error())
	}
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Logger.Infof("%s %s", method, route)
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Logger.Alertf("could not walk api routes. %s", err.Error())
	}
	log.Logger.Info("health api routes")
	healthRouter := healthRoutes()
	if err := chi.Walk(healthRouter, walkFunc); err != nil {
		log.Logger.Alertf("could not walk health routes. %s", err.Error())
	}

	if ssl {
		gc := crypt.GenerateCertificate{
			Organization: "MCS",
			Host:         "127.0.0.1",
			ValidFor:     10 * 365 * 24 * time.Hour,
			IsCA:         false,
			EcdsaCurve:   "P384",
			Ed25519Key:   false,
		}
		tlsConfig, err := gc.GenerateTLSConfig()
		if err != nil {
			log.Logger.Alertf("could not create tls config. %s", err.Error())
		}
		sslsrv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(serviceConfig.Sslport),
			WriteTimeout: time.Second * 30,
			ReadTimeout:  time.Second * 30,
			IdleTimeout:  time.Second * 120,
			Handler:      router,
			TLSConfig:    tlsConfig,
		}
		go func() {
			log.Logger.Infof("starting https server on address: %s", sslsrv.Addr)
			if err := sslsrv.ListenAndServeTLS("", ""); err != nil {
				log.Logger.Alertf("error starting server: %s", err.Error())
			}
		}()
		srv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(serviceConfig.Port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      healthRouter,
		}
		go func() {
			log.Logger.Infof("starting http server on address: %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				log.Logger.Alertf("error starting server: %s", err.Error())
			}
		}()
	} else {
		// own http server for the healthchecks
		srv = &http.Server{
			Addr:         "0.0.0.0:" + strconv.Itoa(serviceConfig.Port),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      router,
		}
		go func() {
			log.Logger.Infof("starting http server on address: %s", srv.Addr)
			if err := srv.ListenAndServe(); err != nil {
				log.Logger.Alertf("error starting server: %s", err.Error())
			}
		}()
	}

	log.Logger.Info("waiting for clients")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	if ssl {
		sslsrv.Shutdown(ctx)
	}

	log.Logger.Info("finished")

	os.Exit(0)
}

func initLogging() {
	log.Logger.SetLevel(serviceConfig.Logging.Level)
	var err error
	serviceConfig.Logging.Filename, err = config.ReplaceConfigdir(serviceConfig.Logging.Filename)
	if err != nil {
		log.Logger.Errorf("error on config dir: %v", err)
	}
	log.Logger.GelfURL = serviceConfig.Logging.Gelfurl
	log.Logger.GelfPort = serviceConfig.Logging.Gelfport
	log.Logger.Init()
}

func initConfig() {
	if port > 0 {
		serviceConfig.Port = port
	}
	if sslport > 0 {
		serviceConfig.Sslport = sslport
	}
	if serviceURL != "" {
		serviceConfig.ServiceURL = serviceURL
	}

	if serviceConfig.TempFolder != "" {
		apiv1.SrcDir = serviceConfig.TempFolder
		if _, err := os.Stat(apiv1.SrcDir); err != nil {
			os.MkdirAll(apiv1.SrcDir, os.ModePerm)
		}
	}
}

func initJaeger(servicename string, config config.OpenTracing) (opentracing.Tracer, io.Closer) {

	cfg := jaegercfg.Configuration{
		ServiceName: servicename,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: config.Host,
			CollectorEndpoint:  config.Endpoint,
		},
	}
	if (config.Endpoint == "") && (config.Host == "") {
		cfg.Disabled = true
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func getApikey() string {
	value := fmt.Sprintf("%s_%s", servicename, "default")
	apikey := fmt.Sprintf("%x", md5.Sum([]byte(value)))
	return strings.ToLower(apikey)
}
