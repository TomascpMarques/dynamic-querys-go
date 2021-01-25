package actions

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ParseActionBody -
func ParseActionBody(regex string, actionContents BodyContents) ([]Endpoint, error) {
	// `"\w+":$`
	functionCallRegexp := regexp.MustCompile(regex)
	functionCalMap := make([]Endpoint, 0)

	for k, v := range actionContents.FuncsContent {
		if len(functionCallRegexp.FindAllString(v, -1)) != 0 {
			params := make([]interface{}, 0)
			for _, j := range actionContents.FuncsContent[k+1:] {
				if string(j[len(j)-1]) == ":" {
					break
				}

				if len(regexp.MustCompile(`^\d+$`).FindAllString(j[:len(j)-1], -1)) != 0 {
					cnvt, err := strconv.Atoi(j[:len(j)-1])
					if err != nil {
						return nil, errors.New("Error Converting to int")
					}
					params = append(params, cnvt)
					continue
				}
				if len(regexp.MustCompile(`^\d+\.\d+$`).FindAllString(j[:len(j)-1], -1)) != 0 {
					cnvt, err := strconv.ParseFloat(j[:len(j)-1], 64)
					if err != nil {
						return nil, errors.New("Error Converting to float64")

					}
					params = append(params, cnvt)
				}
				if len(regexp.MustCompile(`"\{\S+\}"`).FindAllString(j[:len(j)-1], -1)) != 0 {
					test := regexp.MustCompile(`\\`).ReplaceAllLiteralString(j[1:len(j)-2], "")
					var jsonBlob = []byte(test)
					var cnvt map[string]interface{}
					err := json.Unmarshal(jsonBlob, &cnvt)
					if err != nil {
						return nil, errors.New("Error Converting to map[string]interface{}")
					}
					params = append(params, cnvt)
				}
				if len(regexp.MustCompile(`"[a-zA-Z0-9_ ]+"`).FindAllString(j[:len(j)-1], -1)) != 0 {
					cnvt := j[1 : len(j)-2]
					params = append(params, cnvt)
				}
				if len(regexp.MustCompile(`^true$|^false$`).FindAllString(j[:len(j)-1], -1)) != 0 {
					cnvt, err := strconv.ParseBool(j[:len(j)-1])
					if err != nil {
						return nil, errors.New("Error Converting to boolean")
					}
					params = append(params, cnvt)
				}
			}
			pnum, err := GetFunctionParamsNum(reflect.ValueOf(FuncsStorage[v[1:len(v)-2]]))
			if err != nil || len(params) != pnum {
				return nil, errors.New("Bad Parameters")
			}
			functionCalMap = append(functionCalMap, Endpoint{
				funcName: v[1 : len(v)-2],
				params:   params,
			})
		}
	}
	return functionCalMap, nil
}

// ParseActionContents -
func ParseActionContents(request string) (BodyContents, error) {
	var result BodyContents
	result.ActionBody = strings.TrimSpace(request[7:])
	if len(result.ActionBody) == 0 {
		return BodyContents{}, errors.New("Error extracting the action body")
	}
	result.Authentication = regexp.MustCompile(`auth: ".+"|auth: ".+" +`).FindAllStringSubmatch(result.ActionBody, -1)
	if len(result.Authentication) == 0 {
		return BodyContents{}, errors.New("Error extracting authentication token")
	}
	result.FuncCalls = regexp.MustCompile(`"\w+":|"\w+":\s+`).FindAllStringSubmatch(result.ActionBody, -1)
	if len(result.FuncCalls) == 0 {
		return BodyContents{}, errors.New("Error extracting function calls")
	}
	result.FuncArgs = regexp.MustCompile(`"\w+",\n|"\w+",|\[.+\]\n|\[.+\]\s+\n|\[.+\]\s+|\[.+\]|\d+,|\d+.\d+|"[a-zA-Z0-9_ ]+",`).FindAllStringSubmatch(result.ActionBody, -1)
	if len(result.FuncArgs) == 0 {
		return BodyContents{}, errors.New("Error extracting the functions arguments")
	}
	result.FuncsContent = regexp.MustCompile(`\S+,|".+",|".+":`).FindAllString(result.ActionBody, -1)
	if len(result.FuncsContent) == 0 {
		return BodyContents{}, errors.New("Error extracting the <funcs:> content")
	}

	return result, nil
}
