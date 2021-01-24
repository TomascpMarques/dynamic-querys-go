package actions

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// DQGLogger - Logger for DynamicQuerysGo processes (can be stdout or a io.writer to a file)
var DQGLogger = log.New(os.Stdout, "DynamicQuerysGo [*] ", log.LstdFlags)

// FuncMap -
type FuncMap map[string]interface{}

type endpoint struct {
	funcName string
	params   interface{}
}

// FuncsStorage -
var FuncsStorage = FuncMap{
	//*
	//* !!! STATE YOUR FUNCS HERE !!! **//
	//*
	"test123": test123,
}

func test123(n float64) float64 {
	return n * n
}

// Handler - Handles all of the requests coming into the server
func Handler(rw http.ResponseWriter, r *http.Request) {
	// Gets the body from the request
	requestBody, _ := ioutil.ReadAll(r.Body)

	action := strings.TrimSpace(string(requestBody))

	fmt.Println("\n", action)
	if !strings.Contains(action[:7], "action:") {
		fmt.Println("Erro: request sent is not an action")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	actionBody := strings.TrimSpace(action[7:])
	auth := regexp.MustCompile(`auth: ".+"|auth: ".+" +`).FindAllStringSubmatch(actionBody, -1)
	funcNames := regexp.MustCompile(`"\w+":|"\w+":\s+`).FindAllStringSubmatch(actionBody, -1)
	funcArgs := regexp.MustCompile(`"\w+",\n|"\w+",|\[.+\]\n|\[.+\]\s+\n|\[.+\]\s+|\[.+\]|\d+,|\d+.\d+,`).FindAllStringSubmatch(actionBody, -1)
	functionsAndArgs := regexp.MustCompile(`\S+,|"\S+":`).FindAllString(actionBody, -1)
	// funcNames[0][0][1:len(funcNames[0][0])-2]
	fmt.Println("| ->>", auth, "|", len(funcNames), funcArgs, "|-> ", functionsAndArgs)

	pnum, err := GetFunctionParamsNum(reflect.ValueOf(FuncsStorage["test123"]))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("-> ", pnum)

	/*
		pnum, err := GetFunctionParamsNum(reflect.ValueOf(FuncsStorage[string(string(v[0])[1:len(v[0])-2])]))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("passed point 1")
		x.ActionBody.Functions = make([]FunctionPath, 1)
		x.ActionBody.Functions[0].FunctionCall = string(string(v[0])[1:len(v[0])-2])
		test := string(funcArgs[0][0:pnum][0])
		fmt.Println("->>", test[0:len(test)-1])
		x.ActionBody.Functions[0].FunctionParams = append(x.ActionBody.Functions[0].FunctionParams, test[0:len(test)-1])
	*/

	functionCallRegexp := regexp.MustCompile(`"\w+":`)
	functionCalMap := make([]endpoint, 0)

	for k, v := range functionsAndArgs {
		if len(functionCallRegexp.FindAllString(v, -1)) != 0 {
			params := make([]string, 0)
			for _, j := range functionsAndArgs[k+1:] {
				if string(j[len(j)-1]) == ":" {
					break
				}
				params = append(params, j)
			}
			fmt.Println(params)
			functionCalMap = append(functionCalMap, endpoint{
				funcName: v[1 : len(v)-2],
				params:   params,
			})
		}
	}
	fmt.Println(functionCalMap)
	

	rw.Write([]byte(actionBody))
}
