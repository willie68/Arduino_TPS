package model

type SrcCommand struct {
	Line    int
	Cmd     int
	Data    int
	Comment string
}

type TemplateCmds struct {
	EqualsA []string          `json:"equalsA"`
	AEquals []string          `json:"aEquals"`
	ACalc   []string          `json:"aCalc"`
	SkipIf  []string          `json:"skipIf"`
	ABytes  []string          `json:"aBytes"`
	Debug   string            `json:"debug"`
	Boards  map[string]string `json:"boards"`
}
