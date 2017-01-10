package jconfig_test

import (
	"fmt"
	"testing"

	"github.com/choueric/jconfig"
)

const DefaultConfig = `{
	"editor": "vim",
	"current": 0,
	"profile": [
	{
		"name":"first",
		"src_dir":"/home/user/kernel",
		"arch":"arm",
		"cross_compile":"arm-eabi-",
		"target":"uImage",
		"output_dir":"./arm_build"
	},
	{
		"name":"second",
		"src_dir":"/home/user/kernel",
		"arch":"x86",
		"target":"zImage",
		"output_dir":"./x86_build"
	}
	]
}
`

type Profile struct {
	Name        string `json:"name"`
	SrcDir      string `json:"src_dir"`
	Arch        string `json:"arch"`
	Target      string `json:"target"`
	CrossComile string `json:"cross_compile"`
	OutputDir   string `json:"output_dir"`
}

type Config struct {
	Editor   string     `json:"editor"`
	Current  int        `json:"current"`
	Profiles []*Profile `json:"profile"`
}

func (p *Profile) String() string {
	return fmt.Sprintf(" [%s]: use '%s' to build '%s' on '%s' from '%s' to '%s'\n",
		p.Name, p.CrossComile, p.Target, p.Arch, p.SrcDir, p.OutputDir)
}

func (c *Config) String() string {
	return fmt.Sprintf("Editor: %s, Current Profile: %d\n%v\n",
		c.Editor, c.Current, c.Profiles)
}

func Test_New(t *testing.T) {
	config := jconfig.New(".", "config.json", Config{})

	if config.FilePath() != "./config.json" {
		t.Error("filepath is not ./config.json instead of", config.FilePath())
	}

	if config.Data() != nil {
		t.Error("config Data should be nil")
	}

}

func Test_Load(t *testing.T) {
	config := jconfig.New(".", "config.json", Config{})

	p, err := config.Load(DefaultConfig)
	if err != nil {
		t.Error("load config error:", err)
	}

	cc := config.Data().(*Config)
	cc.Current = cc.Current + 1

	pp := p.(*Config)
	if pp.Current != cc.Current {
		t.Error("data does not match.")
	}

	fmt.Println(pp)
}

func Test_Save(t *testing.T) {
	config := jconfig.New(".", "config.json", Config{})

	p, err := config.Load(DefaultConfig)
	if err != nil {
		t.Error("load config error:", err)
		return
	}

	pp := p.(*Config)
	pp.Current = 88

	if err := config.Save(); err != nil {
		t.Error("save config error:", err)
	}
}
