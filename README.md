# jconfig

This is a Go package to parse a configuration file using JSON.

## Installation

```sh
go get github.com/choueric/jconfig

import (
		"github.com/choueric/jconfig"
	   )
```

## Usage
-----

```go
package main

import (
	"fmt"
	"github.com/choueric/jconfig"
)

const DefContent = `
{
	"server": "127.0.0.1:8088"
}
`

type Config struct {
	Server string `json:"server"`
}

func getConfig() *Config {
	jc := jconfig.New(".", "config.json", Config{})

	if _, err := jc.Load(DefContent); err != nil {
		fmt.Println("load config error:", err)
		return nil
	}

	config := jc.Data().(*Config)
	return config
}
```

Refer to `jconfig_test.go` for more details of how to use it.

Another usage is config.go of [kbdashboard](https://github.com/choueric/kernelBuildDashboard.git)
