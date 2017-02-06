package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// Applications is the structure of the file containing all the templates variables
type Applications struct {
	Default map[string]interface{}            `yaml:"default" json:"default"`
	Apps    map[string]map[string]interface{} `yaml:"apps" json:"apps"`
	alters  map[string]VariableAlteration
}

type VariableAlteration func(interface{}) (interface{}, error)

func NewApplications(file string) (*Applications, error) {
	var ret Applications

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read applications file : %s", err)
	}

	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall applications file : %s", err)
	}

	return &ret, nil
}

// Names returns all the applications names availables
func (app *Applications) Names() []string {
	var ret []string

	for name := range app.Apps {
		ret = append(ret, name)
	}
	return ret
}

// AddAlteration adds the possibility to alter an application variable when getting the variables
func (app *Applications) AddAlteration(key string, ptr VariableAlteration) {
	if app.alters == nil {
		app.alters = make(map[string]VariableAlteration)
	}
	app.alters[key] = ptr
}

// Variables returns the default variables overwritten by the app variables
func (app *Applications) Variables(name string) (map[string]interface{}, error) {
	variables, exists := app.Apps[name]
	if !exists {
		return nil, fmt.Errorf("application %s does not exists", name)
	}

	var ret map[string]interface{} = make(map[string]interface{})
	for key, defaultValue := range app.Default {
		ret[key] = defaultValue
	}

	var err error
	for key, value := range variables {
		// Check if the variables contains '...' for prepend the default variables
		strValue, isString := value.(string)
		if isString && strings.HasPrefix(strValue, "...") {
			defaultValue, exists := app.Default[key]
			if !exists {
				return nil, fmt.Errorf("cannot preprend value to key %s : default value does not exists", key)
			}
			defaultStrValue, isString := defaultValue.(string)
			if !isString {
				return nil, fmt.Errorf("cannot preprend value to key %s : default value is not type of string", key)
			}
			value = strings.Replace(strValue, "...", defaultStrValue, 1)
		}

		// Check if we need to alter the value
		if ptr, exists := app.alters[key]; exists {
			value, err = ptr(value)
			if err != nil {
				return nil, fmt.Errorf("failed to apply alter on key %s : %s", key, err)
			}
		}
		ret[key] = value
	}
	ret["id"] = name

	return ret, nil
}
