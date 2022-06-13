package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/prometheus/common/log"
	"github.com/willie68/tps_asm/internal/asm"
)

type MyEvent struct {
	Asm          string `json:"asm"`
	Name         string `json:"name"`
	OutputFormat string `json:"outputformat"`
	Hardware     string `json:"hardware"`
}

type MyResponse struct {
	Binary      []byte `json:"binary"`
	Name        string `json:"name"`
	Format      string `json:"format"`
	Destination string `json:"destination"`
	Hardware    string `json:"hardware"`
}

func HandleLambdaEvent(p MyEvent) (MyResponse, error) {
	scanner := bufio.NewScanner(strings.NewReader(p.Asm))
	src := make([]string, 0)
	for scanner.Scan() {
		src = append(src, scanner.Text())
		// handle first encountered error while reading
		if err := scanner.Err(); err != nil {
			return MyResponse{}, fmt.Errorf("error scanning source: %v", err)
		}
	}

	tpsasm := asm.Assembler{
		Hardware:     asm.ParseHardware(p.Hardware),
		Source:       src,
		Includes:     "",
		Filename:     p.Name,
		Outputformat: strings.ToLower(p.OutputFormat),
	}

	errs := tpsasm.Parse()
	if len(errs) > 0 {
		errsl := make([]string, 0)
		for _, err := range errs {
			log.Errorf("%v", err)
			errsl = append(errsl, fmt.Sprintf("%+v", err))
		}
		return MyResponse{}, fmt.Errorf("compile error: %v", strings.Join(errsl, "; "))
	}

	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	err := asm.TpsFile(foo, tpsasm)
	if err != nil {
		return MyResponse{}, fmt.Errorf("error generating file: %v", err)
	}
	foo.Flush()

	return MyResponse{
		Binary:      tpsasm.Binary,
		Name:        tpsasm.Filename,
		Format:      tpsasm.Outputformat,
		Destination: b.String(),
		Hardware:    tpsasm.Hardware.String(),
	}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
