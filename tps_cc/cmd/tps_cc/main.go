package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	flag "github.com/spf13/pflag"
	"github.com/willie68/tps_cc/pkg/tmpl"
)

var srcFile string

type SrcCommand struct {
	Line    int
	Cmd     int
	Data    int
	Comment string
}

type templateCmds struct {
	EqualsA []string `json:"equalsA"`
	AEquals []string `json:"aEquals"`
	ACalc   []string `json:"aCalc"`
	SkipIf  []string `json:"skipIf"`
	ABytes  []string `json:"aBytes"`
}

var commandSrc []SrcCommand
var name string
var path string
var TemplateCmds templateCmds

func init() {
	log.Println("init")
	flag.StringVarP(&srcFile, "source", "s", "", "this is the path and filename to the source file")

	commandSrc = make([]SrcCommand, 0)
	dat, err := tmpl.TemplateFS.ReadFile("files/template.json")
	if err != nil {
		log.Panicf("can't read template.json: %v", err)
	}

	err = json.Unmarshal(dat, &TemplateCmds)
	if err != nil {
		log.Panicf("can't unmarshall template.json: %v", err)
	}
}

func main() {
	log.Println("starting")

	flag.Parse()

	// reading tps file
	if srcFile == "" {
		log.Panicf("source file cant't be empty.")
	}

	name = strings.TrimSuffix(srcFile, ".tps")
	name = filepath.Base(name)
	path = filepath.Dir(srcFile)

	tpsSrc, err := readLines(srcFile)
	if err != nil {
		log.Panicf("can't read file: %v", err)
	}

	fmt.Printf("source: %v\r\n", tpsSrc)

	// generating structures
	commandSrc, err = parse(tpsSrc)
	if err != nil {
		log.Panicf("parsing error: %v", err)
	}

	// generating sources
	err = generate(commandSrc)
	if err != nil {
		log.Panicf("generating error: %v", err)
	}
	// compiling sources
}

func generate(commandSrc []SrcCommand) error {
	subPath := path + "/" + name
	os.Mkdir(subPath, os.ModePerm)
	dstFile := subPath + "/" + name + ".ino"
	copyNeededFiles(subPath)

	dat, err := tmpl.TemplateFS.ReadFile("files/template.ino")
	if err != nil {
		return fmt.Errorf("can't read template: %v", err)
	}

	tmpl, err := template.New(name).Parse(string(dat))
	if err != nil {
		return fmt.Errorf("can't parse template: %v", err)
	}

	main := generateMain(commandSrc)

	context := make(map[string]string)
	context["setup"] = ""
	context["main"] = main

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

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func parse(tpsSrc []string) ([]SrcCommand, error) {

	for _, l := range tpsSrc {
		if strings.HasPrefix(l, "#") {
			continue
		}
		s := strings.Split(l, ",")
		nStr := strings.Replace(s[0], "0X", "", -1)
		nStr = strings.Replace(nStr, "0x", "", -1)

		v, err := strconv.ParseInt(nStr, 16, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing line number: %v", err)
		}
		n := int(v)
		v, err = strconv.ParseInt(s[1], 16, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing command number: %v", err)
		}
		c := int(v)
		v, err = strconv.ParseInt(s[2], 16, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing data number: %v", err)
		}
		d := int(v)
		cStr := ""
		if len(s) > 3 {
			cStr = s[3]
		}
		if n > len(commandSrc) {
			for i := 0; i <= (n - len(commandSrc)); i++ {
				cmd := SrcCommand{
					Line:    len(commandSrc),
					Cmd:     0,
					Data:    0,
					Comment: "n.n.",
				}
				commandSrc = append(commandSrc, cmd)
			}
		}
		cmd := SrcCommand{
			Line:    n,
			Cmd:     c,
			Data:    d,
			Comment: cStr,
		}
		commandSrc = append(commandSrc, cmd)
	}
	return commandSrc, nil
}

func copyNeededFiles(subPath string) {

	// copy needed files
	files, err := tmpl.TemplateFS.ReadDir("files")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "template.ino") {
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

func generateMain(commandSrc []SrcCommand) string {
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
			fmt.Fprintf(w, "  %s;\r\n", TemplateCmds.EqualsA[c.Data])
		case 6:
			fmt.Fprintf(w, "  %s;\r\n", TemplateCmds.AEquals[c.Data])
		case 7:
			fmt.Fprintf(w, "  %s;\r\n", TemplateCmds.ACalc[c.Data])
		case 8:
			fmt.Fprintf(w, "  page=%d;\r\n", c.Data)
		case 9:
			fmt.Fprintf(w, "  tmpValue = (page * 16) + %d;goto *array[tmpValue];\r\n", c.Data)
		case 10:
			fmt.Fprintf(w, "  if (c>0) { c=c-1; tmpValue = (page * 16) + %d;goto *array[tmpValue];}\r\n", c.Data)
		case 11:
			fmt.Fprintf(w, "  if (d>0) { d=d-1; tmpValue = (page * 16) + %d;goto *array[tmpValue];}\r\n", c.Data)
		case 12:
			fmt.Fprintf(w, "  if (%s) { goto *array[%d];}\r\n", TemplateCmds.SkipIf[c.Data], x+2)
		case 13:
			fmt.Fprintf(w, "  saveaddr[saveCnt] = %d;\r\n  saveCnt++;\r\n  tmpValue = (page * 16) + %d;\r\n  goto *array[tmpValue];\r\n", x+1, c.Data)
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
			fmt.Fprintf(w, "  %s;\r\n", TemplateCmds.ABytes[c.Data])
		default:
			fmt.Fprintf(w, "  // unknown command 0x%x%x;\r\n", c.Cmd, c.Data)
			log.Printf("unknown command in line %d (0x%02x) : 0x%x%x", c.Line, c.Line, c.Cmd, c.Data)
		}
	}

	fmt.Fprint(w, "  ;")
	return w.String()
}
