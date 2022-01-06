package apiv1

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/willie68/tps_cc/internal/serror"
	"github.com/willie68/tps_cc/internal/utils"
	"github.com/willie68/tps_cc/pkg/tpsfile"
	"github.com/willie68/tps_cc/pkg/tpsgen"

	"github.com/willie68/tps_cc/internal/utils/httputils"
)

var (
	postGenerateCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tps_cc_post_gen_total",
		Help: "The total number of generating requests",
	})
)

var srcDir = "example"

/*
ConfigRoutes getting all routes for the config endpoint
*/
func GenerateRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/", PostGenerateEndpoint)
	return router
}

// PostGenerateEndpoint create a new store for a tenant
// because of the automatic store creation, this method will always return 201
// @Summary Create a new store for a tenant
// @Tags configs
// @Accept  json
// @Produce  json
// @Param payload body string true "Add store"
// @Success 201 {string} string "tenant"
// @Failure 400 {object} serror.Serr "client error information as json"
// @Failure 500 {object} serror.Serr "server error information as json"
// @Router /config [post]
func PostGenerateEndpoint(response http.ResponseWriter, request *http.Request) {
	mimeType := request.Header.Get("Content-Type")

	var cntLength int64
	var filename string
	var f io.Reader
	if strings.HasPrefix(mimeType, "multipart/form-data") {
		err := request.ParseMultipartForm(1024 * 1024 * 1024)
		if err != nil {
			httputils.Err(response, request, serror.InternalServerError(err))
			return
		}
		mpf, fileHeader, err := request.FormFile("file")
		if err != nil {
			httputils.Err(response, request, serror.InternalServerError(err))
			return
		}

		mimeType = fileHeader.Header.Get("Content-type")
		cntLength = fileHeader.Size
		filename = fileHeader.Filename
		f = mpf
		defer mpf.Close()
	} else {
		mpf := request.Body
		defer mpf.Close()
		cntLength = -1
		filename = request.Header.Get("tpsfilename")
		if filename == "" {
			filename = "tps_file.tps"
		}
		f = mpf
	}

	SrcFile := srcDir + "/" + filename
	file, err := os.Create(SrcFile)
	if err != nil {
		httputils.Err(response, request, serror.InternalServerError(err))
		return
	}
	defer file.Close()
	io.Copy(file, f)

	name := strings.TrimSuffix(filename, ".tps")
	path := utils.GenerateRamdomPath(srcDir)

	// generating structures
	commandSrc, err := tpsfile.ParseFile(SrcFile)
	if err != nil {
		log.Printf("parsing error: %v", err)
		httputils.Err(response, request, serror.InternalServerError(err))
		return
	}

	flags := tpsfile.GetBuildFlags(commandSrc)
	Board := "arduino_uno"

	// generating sources
	tpsgen := tpsgen.TPSGen{
		Name:         name,
		Path:         path,
		TemplateCmds: tpsgen.TemplateCmds,
		Debug:        true,
		Board:        Board,
		BuildFlags:   flags,
	}
	err = tpsgen.Generate(commandSrc)
	if err != nil {
		log.Printf("generating error: %v", err)
		httputils.Err(response, request, serror.InternalServerError(err))
		return
	}

	// compiling sources

	AutoCompile := true
	if AutoCompile {
		// arduino-cli compile --clean -e -b arduino:avr:uno --output-dir /home/arduinocli/Arduino_TPS/dest/ ./ --build-property="build.extra_flags=-DTPS_ENHANCEMENT -DTPS_SERVO -DTPS_SERIAL_PRG -DTPS_TONE" >>log.$TPS_VERSION.log 2>&1
		// cp /home/arduinocli/Arduino_TPS/dest/TPS.ino.hex /home/arduinocli/Arduino_TPS/dest/TPS.$TPS_VERSION.UNO.ENHANCEMENT.SERVO.TONE.SERIAL_PRG.hex
		log.Printf("starting compilation for %s", Board)
		r, err := tpsgen.Compile()
		if err != nil {
			log.Printf("compile error: %s,\r\n %v", r, err)
			httputils.Err(response, request, serror.InternalServerError(err))
			return
		}
		log.Printf("compilation finnished with: %s", r)
	}

	// building a zip
	zip, err := tpsgen.ZipFiles()
	if err != nil {
		log.Printf("generating zip error: %v", err)
		httputils.Err(response, request, serror.InternalServerError(err))
		return
	}
	log.Printf("building Zip File: %s", zip)

	/*
		r := make(map[string]interface{})
		r["name"] = filename
		r["contentLength"] = cntLength
		r["mimetype"] = mimeType
	*/

	response.Header().Set("Content-Type", "application/zip")
	response.Header().Set("TPS-Content-Length", fmt.Sprintf("%d", cntLength))
	response.WriteHeader(http.StatusCreated)
	r, err := os.Open(zip)
	if err != nil {
		httputils.Err(response, request, serror.InternalServerError(err))
		return
	}
	defer r.Close()
	_, err = io.Copy(response, r)
	if err != nil {
		httputils.Err(response, request, serror.InternalServerError(err))
		return
	}
	go func() {
		err := os.RemoveAll(path)
		if err != nil {
			log.Print(err)
		}
	}()
	return
}
