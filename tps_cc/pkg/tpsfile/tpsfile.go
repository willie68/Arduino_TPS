package tpsfile

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/willie68/tps_cc/pkg/model"
)

func ParseFile(file string) ([]model.SrcCommand, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return Parse(lines)
}

func Parse(tpsSrc []string) ([]model.SrcCommand, error) {
	commandSrc := make([]model.SrcCommand, 0)

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
				cmd := model.SrcCommand{
					Line:    len(commandSrc),
					Cmd:     0,
					Data:    0,
					Comment: "n.n.",
				}
				commandSrc = append(commandSrc, cmd)
			}
		}
		cmd := model.SrcCommand{
			Line:    n,
			Cmd:     c,
			Data:    d,
			Comment: cStr,
		}
		commandSrc = append(commandSrc, cmd)
	}
	return commandSrc, nil
}

func GetBuildFlags(src []model.SrcCommand) []string {
	rcReceiver := false
	enhancement := false
	servo := false
	tone := false

	for _, cmd := range src {
		if cmd.Cmd == 5 {
			if cmd.Data == 11 || cmd.Data == 12 {
				servo = true
			}
			if cmd.Data == 0 || cmd.Data >= 11 {
				enhancement = true
			}
		}
		if cmd.Cmd == 6 {
			if cmd.Data >= 11 {
				enhancement = true
			}
			if cmd.Data == 11 || cmd.Data == 12 {
				rcReceiver = true
			}
		}
		if cmd.Cmd == 7 {
			if cmd.Data >= 11 {
				enhancement = true
			}
		}
		if cmd.Cmd == 8 {
			if cmd.Data >= 8 {
				enhancement = true
			}
		}
		if cmd.Cmd == 12 {
			if cmd.Data == 0 {
				enhancement = true
			}
		}
		if cmd.Cmd == 14 {
			if cmd.Data != 0 {
				enhancement = true
			}
		}
		if cmd.Cmd == 15 {
			enhancement = true
			if cmd.Data == 6 || cmd.Data == 7 {
				servo = true
			}
			if cmd.Data == 8 {
				tone = true
			}
			if cmd.Data == 2 || cmd.Data == 3 {
				rcReceiver = true
			}
		}
	}

	flags := make([]string, 0)
	if rcReceiver {
		flags = append(flags, "TPS_RCRECEIVER")
	}
	if enhancement {
		flags = append(flags, "TPS_ENHANCEMENT")
	}
	if servo {
		flags = append(flags, "TPS_SERVO")
	}
	if tone {
		flags = append(flags, "TPS_TONE")
	}
	return flags
}
