package main

import (
	"DynamicQuerysGo/actions"
	"DynamicQuerysGo/functions"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/gorilla/mux"
)

// FuncMap -
type FuncMap map[string]interface{}

// FuncsStorage -
var FuncsStorage = FuncMap{}

// CallFunc - Calls the function with by the name specified in funcName
func CallFunc(funcName string, params []interface{}) (interface{}, error) {
	// Gets function as reflect.Value to perform reflection,
	// to know things such as number of parameters
	function := reflect.ValueOf(FuncsStorage[funcName])

	// Checks if the passed parameters are more or less than the ones required
	if len(params) != function.Type().NumIn() {
		return nil, errors.New("The number of params is insufficient")
	}

	// Gets al the parameters passed in params
	// to be used in reflect.Call, as the called
	funcParams := make([]reflect.Value, len(params))
	for k, param := range params {
		funcParams[k] = reflect.ValueOf(param)
	}

	// Call calls the function v with the input arguments in.
	// For example, if len(funcParams) == 3, v.Call(funcParams),
	// represents the Go call v(funcParams[0], funcParams[1], funcParams[2]).
	calledFunction := function.Call(funcParams)

	// Gets the return values as interfaces, allocated in a interface array
	returned := make([]interface{}, len(calledFunction))
	for key, value := range calledFunction {
		returned[key] = value.Interface()
	}
	return returned, nil
}

// ActionsHandler -
func ActionsHandler(rw http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)

	var i actions.Action
	err := json.Unmarshal(requestBody, &i)
	if err != nil {
		log.Println(err)
	}

	funcParams, err := convertToPrimitives(i.ActionBody.Functions[0].FunctionParams)
	if err != nil {
		log.Panicf("Error: %s", err)
	}

	FuncsStorage = map[string]interface{}{
		"Test4": functions.Test4,
		"Test3": functions.Test3,
	}

	res, _ := CallFunc(i.ActionBody.Functions[0].FunctionCall, funcParams)
	fmt.Println(res)

	result, _ := json.MarshalIndent(&funcParams, "", " ")
	rw.Write(result)
}

func convertToPrimitives(x []interface{}) ([]interface{}, error) {
	converted := make([]interface{}, len(x))
	for k, v := range x {
		varType := reflect.TypeOf(v.(interface{}))
		if varType.String() == "[]interface {}" {
			arrayType := reflect.TypeOf(v.([]interface{})[0])

			switch arrayType.String() {
			case "string":
				newArray := make([]string, len(v.([]interface{})))
				for j, u := range v.([]interface{}) {
					newArray[j] = u.(string)
				}
				converted[k] = newArray
				break
			case "float64":
				newArray := make([]float64, len(v.([]interface{})))
				for j, u := range v.([]interface{}) {
					newArray[j] = u.(float64)
				}
				converted[k] = newArray
				break
			default:
				converted[k] = v
				break
			}

			if converted[k] == nil {
				return nil, errors.New("Error Converting")
			}
		} else {
			converted[k] = reflect.ValueOf(v).Convert(varType).Interface()
			if converted[k] == nil {
				return nil, errors.New("Error Converting")
			}
		}
	}
	return converted, nil
}

// DQGPORT - the port for where DynamicQuerysGo is located
var DQGPORT = os.Getenv("ENV_GOACTIONS_PORT")

// DEFAULTDQGPORT - default port for DynamicQuerysGo
const DEFAULTDQGPORT = "8000"

// DQGLogger - Logger for DynamicQuerysGo processes (can be stdout or a io.writer to a file)
var DQGLogger = log.New(os.Stdout, "GoActions [*] ", log.LstdFlags)

func main() {
	// Checks for port configuration for the service
	if DQGPORT == "" {
		DQGPORT = DEFAULTDQGPORT
	}

	// flag setup fo graceful-shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	routerMux := mux.NewRouter()
	routerMux.HandleFunc("/actions", ActionsHandler)

	// server setup
	srv := &http.Server{
		Handler:      routerMux,
		Addr:         "localhost:" + DQGPORT,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	// prevent server blocking
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			DQGLogger.Println(err)
		}
	}()

	// Graceful-Shutdown implementation
	// Credit - https://github.com/gorilla/mux#graceful-shutdown
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGKILL,
	// SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)

	DQGLogger.Println("shutting down")
	os.Exit(0)
}
