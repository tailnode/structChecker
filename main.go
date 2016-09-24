package main

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type t struct {
	Name  string                 `json:"name" check:"notZero"`
	Regex string                 `check:"r:^[a-z]+$"`
	Map   map[string]interface{} `check:"notZero"`
	Slice []string               `check:"notZero"`
}

const checkKey = "check"

func main() {
	case1 := t{
		"",
		"abc0",
		nil,
		[]string{"abc"},
	}

	checkStruct(case1)
}
func checkStruct(t interface{}) error {
	st := reflect.TypeOf(t)
	sv := reflect.ValueOf(t)
	for i := 0; i < st.NumField(); i++ {
		if policy := st.Field(i).Tag.Get(checkKey); policy != "" {
			ok, err := checkValue(sv.Field(i), policy)
			if err != nil {
				fmt.Println(i, err)
				continue
			}
			if !ok {
				fmt.Println(i, "check failed")
				continue
			}
			fmt.Println(i, "check ok")
		}
	}
	return nil
}

func checkValue(value reflect.Value, policy string) (bool, error) {
	switch value.Kind() {
	case reflect.String:
		return checkString(value.String(), policy)
	}
	return true, nil
}

func checkString(str string, policy string) (bool, error) {
	switch {
	case policy == "notZero":
		if str != "" {
			return true, nil
		}
	case strings.HasPrefix(policy, "r:"):
		pattern := strings.TrimLeft(policy, "r:")
		matched, err := regexp.MatchString(pattern, str)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	default:
		return false, errors.New("invalid policy")

	}
	return false, nil
}
