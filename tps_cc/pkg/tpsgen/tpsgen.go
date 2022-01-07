package tpsgen

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/willie68/tps_cc/internal/config"
	"github.com/willie68/tps_cc/pkg/model"
	"github.com/willie68/tps_cc/pkg/tmpl"
)

var TemplateCmds model.TemplateCmds

func init() {
	dat, err := tmpl.TemplateFS.ReadFile("files/template.json")
	if err != nil {
		log.Panicf("can't read template.json: %v", err)
	}

	err = json.Unmarshal(dat, &TemplateCmds)
	if err != nil {
		log.Panicf("can't unmarshall template.json: %v", err)
	}
}

type TPSGen struct {
	Name         string
	Path         string
	Debug        bool
	TemplateCmds model.TemplateCmds
	Board        string
	BuildFlags   []string
}

func (t *TPSGen) Generate(commandSrc []model.SrcCommand) error {
	subPath := t.Path + "/" + t.Name
	os.Mkdir(subPath, os.ModePerm)
	dstFile := subPath + "/" + t.Name + ".ino"
	t.copyNeededFiles(subPath)

	dat, err := tmpl.TemplateFS.ReadFile("files/template.ino")
	if err != nil {
		return fmt.Errorf("can't read template: %v", err)
	}

	tmpl, err := template.New(t.Name).Parse(string(dat))
	if err != nil {
		return fmt.Errorf("can't parse template: %v", err)
	}

	main := t.generateMain(commandSrc)

	context := make(map[string]string)
	context["setup"] = ""
	context["main"] = main
	context["debug"] = ""
	if t.Debug {
		context["debug"] = t.TemplateCmds.Debug
	}
	context["flags"] = ""
	if t.BuildFlags != nil {
		s := new(strings.Builder)
		for _, flag := range t.BuildFlags {
			s.WriteString(fmt.Sprintf("#define %s \r\n", flag))
		}

		context["flags"] = s.String()
	}

	file, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("can't create dest file: %v", err)
	}
	defer file.Close()
	err = tmpl.Execute(file, context)
	if err != nil {
		panic(err)
	}

	return nil
}

func (t *TPSGen) copyNeededFiles(subPath string) {

	// copy needed files
	files, err := tmpl.TemplateFS.ReadDir("files")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "template") {
			continue
		}
		dst := subPath + "/" + f.Name()

		fin, err := tmpl.TemplateFS.Open("files/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer fin.Close()

		fout, err := os.Create(dst)
		if err != nil {
			log.Fatal(err)
		}
		defer fout.Close()

		_, err = io.Copy(fout, fin)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (t *TPSGen) generateMain(commandSrc []model.SrcCommand) string {
	w := new(strings.Builder)

	fmt.Fprint(w, "  static void *array[] = { ")
	for x, c := range commandSrc {
		if x > 0 {
			fmt.Fprint(w, ", ")
		}
		fmt.Fprintf(w, "&&label_%d", c.Line)
	}
	fmt.Fprint(w, "};\r\n")

	for x, c := range commandSrc {
		log.Printf("%d: %v\r\n", x, c)
		fmt.Fprintf(w, "label_%d: ", c.Line)
		if c.Comment != "" {
			fmt.Fprintf(w, "// %s", c.Comment)
		}
		fmt.Fprint(w, "\r\n")
		switch c.Cmd {
		case 0:
			fmt.Fprintf(w, " ; // noop \r\n")
		case 1:
			fmt.Fprintf(w, "  doPort(%d);\r\n", c.Data)
		case 2:
			fmt.Fprintf(w, "  doDelay(%d);\r\n", c.Data)
		case 3:
			fmt.Fprintf(w, "  goto label_%d;\r\n", c.Line-c.Data)
		case 4:
			fmt.Fprintf(w, "  a=%d;\r\n", c.Data)
		case 5:
			fmt.Fprintf(w, "  %s;\r\n", t.TemplateCmds.EqualsA[c.Data])
		case 6:
			fmt.Fprintf(w, "  %s;\r\n", t.TemplateCmds.AEquals[c.Data])
		case 7:
			fmt.Fprintf(w, "  %s;\r\n", t.TemplateCmds.ACalc[c.Data])
		case 8:
			fmt.Fprintf(w, "  page=%d;\r\n", c.Data)
		case 9:
			fmt.Fprintf(w, "  goto *array[(page * 16) + %d];\r\n", c.Data)
		case 10:
			fmt.Fprintf(w, "  if (c>0) { c=c-1; goto *array[(page * 16) + %d];}\r\n", c.Data)
		case 11:
			fmt.Fprintf(w, "  if (d>0) { d=d-1; goto *array[(page * 16) + %d];}\r\n", c.Data)
		case 12:
			fmt.Fprintf(w, "  if (%s) { goto *array[%d];}\r\n", t.TemplateCmds.SkipIf[c.Data], x+2)
		case 13:
			fmt.Fprintf(w, "  saveaddr[saveCnt] = %d;\r\n  saveCnt++;\r\n  goto *array[(page * 16) + %d];\r\n", x+1, c.Data)
		case 14:
			switch c.Data {
			case 0:
				fmt.Fprint(w, "  if (saveCnt < 0) {\r\n    doReset();\r\n    goto *array[0];\r\n  }\r\n  saveCnt -= 1;\r\n  tmpValue = saveaddr[saveCnt];\r\n  goto *array[tmpValue];\r\n")
			case 15:
				fmt.Fprint(w, "  doReset();\r\n  goto *array[0];\r\n")
			default:
				fmt.Fprint(w, "  ;\r\n")
			}
		case 15:
			fmt.Fprintf(w, "  %s;\r\n", t.TemplateCmds.ABytes[c.Data])
		default:
			fmt.Fprintf(w, "  // unknown command 0x%x%x;\r\n", c.Cmd, c.Data)
			log.Printf("unknown command in line %d (0x%02x) : 0x%x%x", c.Line, c.Line, c.Cmd, c.Data)
		}
	}

	fmt.Fprint(w, "  ;")
	return w.String()
}

//ZipFiles zip all generated files into a single zip file
func (t *TPSGen) ZipFiles() (string, error) {
	subPath := t.Path + "/" + t.Name

	dest := t.Path + "/" + t.Name + ".zip"
	file, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		inZipPath := t.Name + "/" + filepath.Base(path)

		f, err := w.Create(inZipPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(subPath, walker)
	if err != nil {
		return "", err
	}
	return dest, nil
}

func (t *TPSGen) Compile() (string, error) {
	s := new(strings.Builder)
	for _, flag := range t.BuildFlags {
		s.WriteString(fmt.Sprintf("-D%s ", flag))
	}

	cmd := exec.Command(
		config.Get().ArduinoCli,
		"compile",
		"--clean",
		"-e",
		"-b",
		t.TemplateCmds.Boards[t.Board],
		"--output-dir",
		t.Path+"/"+t.Name,
		t.Path+"/"+t.Name,
		fmt.Sprintf("--build-property=\"build.extra_flags=%s\"", s.String()),
	)

	log.Printf("compile: %s", cmd)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	//	log.Printf("out:%s", outStr)
	//	log.Printf("err:%s", errStr)

	if err != nil {
		return string(outStr), fmt.Errorf("compile with exit code: %d, \r\n%s", cmd.ProcessState.ExitCode(), string(errStr))
	}

	return string(outStr), nil
}
