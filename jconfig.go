package jconfig

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
)

type JConfig struct {
	Path       string
	Filename   string
	configType reflect.Type
	Data       interface{} // pointer of structure
}

func (c *JConfig) FilePath() string {
	return c.Path + "/" + c.Filename
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
		if _, err = file.Write([]byte(defContent)); err != nil {
			return err
		}
		file.Close()
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

func New(path, filename string, t interface{}) *JConfig {
	c := &JConfig{Path: path, Filename: filename}
	c.configType = reflect.TypeOf(t)

	return c
}

// return type is pointer to the structure.
func (c *JConfig) Load(defContent string) (interface{}, error) {
	if c.Filename == "" || c.Path == "" {
		return nil, errors.New("jconfig: invalid path")
	}

	if err := checkPath(c.Path); err != nil {
		return nil, err
	}
	name := c.FilePath()
	if err := checkFile(name, defContent); err != nil {
		return nil, err
	}

	v := reflect.New(c.configType)
	initializeStruct(c.configType, v.Elem())
	c.Data = v.Interface()
	//fmt.Println("type is", reflect.TypeOf(c.Data))

	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(b, c.Data); err != nil {
		return nil, err
	}

	return c.Data, nil
}

func (c *JConfig) Save() error {
	b, err := json.MarshalIndent(c.Data, "  ", "  ")
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
