package apiv1

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/willie68/tps_asm/internal/asm"
	log "github.com/willie68/tps_asm/internal/logging"
	"github.com/willie68/tps_asm/internal/serror"
	"github.com/willie68/tps_asm/internal/utils/httputils"
)

var (
	postAsmCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tpsasm_post_asm_total",
		Help: "The total number of post asm requests",
	})
)

/*
AsmRoutes getting all routes for the asm endpoint
*/
func AsmRoutes() (string, *chi.Mux) {
	router := chi.NewRouter()
	router.Post("/generate", PostAsm)
	return baseURL + "/asm", router
}

// PostConfig create a new store for a tenant
// because of the automatic store creation, this method will always return 201
// @Summary Create a new store for a tenant
// @Tags configs
// @Accept  json
// @Produce  json
// @Security api_key
// @Param tenant header string true "Tenant"
// @Param payload body string true "Add store"
// @Success 201 {string} string "tenant"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Failure 500 {object} serror.Serr "server error information as json"
// @Router /config [post]
func PostAsm(response http.ResponseWriter, request *http.Request) {
	log.Info("postasm")
	postAsmCounter.Inc()
	request.ParseMultipartForm(10 << 20)

	file, handler, err := request.FormFile("file")
	if err != nil {
		msg := fmt.Sprintf("error retrieving file: %v", err)
		httputils.Err(response, request, serror.BadRequest(nil, "retrieving-file", msg))
		return
	}
	defer file.Close()

	hardware := request.FormValue("board")
	outputformat := request.FormValue("outputformat")
	log.Infof("Uploaded File: %+v", handler.Filename)
	log.Infof("File Size: %+v", handler.Size)
	log.Infof("MIME Header: %+v", handler.Header)
	log.Infof("Hardwareboard: %s", hardware)

	fileScanner := bufio.NewScanner(file)
	src := make([]string, 0)
	for fileScanner.Scan() {
		src = append(src, fileScanner.Text())
		// handle first encountered error while reading
		if err := fileScanner.Err(); err != nil {
			msg := fmt.Sprintf("error retrieving file: %v", err)
			httputils.Err(response, request, serror.BadRequest(nil, "retrieving-file", msg))
			return
		}
	}

	tpsasm := asm.Assembler{
		Hardware:     asm.ParseHardware(hardware),
		Source:       src,
		Includes:     "",
		Filename:     handler.Filename,
		Outputformat: strings.ToLower(outputformat),
	}
	errs := tpsasm.Parse()
	if len(errs) > 0 {
		errsl := make([]string, 0)
		for _, err := range errs {
			log.Errorf("%v", err)
			errsl = append(errsl, fmt.Sprintf("%+v", err))
		}
		serr := serror.New(http.StatusBadRequest, "compile-error", strings.Join(errsl, "; "))
		httputils.Err(response, request, serr)
		return
	}
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	err = asm.TpsFile(foo, tpsasm)
	if err != nil {
		msg := fmt.Sprintf("error generating file: %v", err)
		httputils.Err(response, request, serror.BadRequest(nil, "generating-file", msg))
		return
	}
	foo.Flush()

	asmresult := struct {
		Binary      []byte `json: binary`
		Name        string `json: name`
		Format      string `json: format`
		Destination string `json: destination`
	}{
		Binary:      tpsasm.Binary,
		Name:        tpsasm.Filename,
		Format:      tpsasm.Outputformat,
		Destination: b.String(),
	}
	render.Status(request, http.StatusOK)
	render.JSON(response, request, asmresult)
}
