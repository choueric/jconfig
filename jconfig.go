// Package jconfig provides struct JConfig to handle with configurations in
// JSON format.
package jconfig

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"reflect"
)

// JConfig structure is an entity representing configurations.
type JConfig struct {
	dir        string
	filename   string
	configType reflect.Type
	data       interface{} // pointer of structure
}

// FilePath returns the full path of configuration file.
func (c *JConfig) FilePath() string {
	return c.dir + "/" + c.filename
}

// Dir returns the path of directory containing configuration file.
func (c *JConfig) Dir() string {
	return c.dir
}

// Filename returns just the file name of configuration file.
func (c *JConfig) Filename() string {
	return c.filename
}

// New returns a pointer of JConfig that contains information of configuration
// file path and user-defined configuration type t.
func New(filepath string, t interface{}) *JConfig {
	var dir, file string
	if path.IsAbs(filepath) {
		dir = path.Dir(filepath)
		file = path.Base(filepath)
	} else {
		dir = "."
		file = filepath
	}
	c := &JConfig{dir: dir, filename: file}
	c.configType = reflect.TypeOf(t)

	return c
}

// Data retruns the user-defined configuration data.
func (c *JConfig) Data() interface{} {
	return c.data
}

// Load loads and parses the configuration file, then allocats a user-defined
// configuration variable which is stored in JConfig and is returned.
// If directory or file does not exist, Load will create with defContent.
func (c *JConfig) Load(defContent string) (interface{}, error) {
	if c.filename == "" || c.dir == "" {
		return nil, errors.New("jconfig: invalid path")
	}

	if err := checkPath(c.dir); err != nil {
		return nil, err
	}
	name := c.FilePath()
	if err := checkFile(name, defContent); err != nil {
		return nil, err
	}

	v := reflect.New(c.configType)
	initializeStruct(c.configType, v.Elem())
	c.data = v.Interface()
	//fmt.Println("type is", reflect.TypeOf(c.data))

	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, c.data); err != nil {
		return nil, err
	}

	return c.data, nil
}

// Save writes the user-defined configruations into JSON file.
func (c *JConfig) Save() error {
	b, err := json.MarshalIndent(c.data, "  ", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(c.FilePath())
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(b)

	return nil
}

func checkPath(path string) error {
	if err := os.MkdirAll(path, os.ModeDir|0777); err != nil {
		return err
	}
	return nil
}

func checkFile(filepath string, defContent string) error {
	_, err := os.Stat(filepath)
	if err != nil && os.IsNotExist(err) {
		file, err := os.Create(filepath)
		defer file.Close()
		if _, err = file.Write([]byte(defContent)); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			initializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
}
