package main

import (
	"DynamicQuerysGo/actions"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/mux"
)

// FuncMap -
type FuncMap map[string]interface{}

// FuncsStorage -
var FuncsStorage = FuncMap{}

// ActionHandler -
func ActionHandler(rw http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	var i actions.Action
	err := json.Unmarshal(body, &i)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(i)
	fmt.Println(i.ActionBody.Authentication)
	fmt.Println(i.ActionBody.Functions[0].FunctionCall)
	fmt.Println(i.ActionBody.Functions[0].FunctionParams)
	fmt.Println(i.ActionBody.Functions[0].FunctionParams["fields"].([]interface{})[0])
	fmt.Println("->", i.ActionBody.Functions[0].FunctionParams["id"])
	fmt.Println(i.ActionBody.Returns)

	rw.Write(body)
}

func main() {

	routerMux := mux.NewRouter()
	routerMux.HandleFunc("/actions", ActionHandler)
	http.Handle("/", routerMux)

	srv := &http.Server{
		Handler:      routerMux,
		Addr:         "localhost:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

	// jsonBlob := []byte(`{"action": {"auth": "mh354kh2vhvqyçyavq5yq8q5yyjuqqy5","func": [{"call": "GetRegisto","params": {"id": "Registo123","fields": ["name", "id"]}},{"call": "GetStuff","params": {"id": "123","fields": ["id"]}}],"returns": ["success","id","nome"] }}`)

	// var i actions.Action
	// err := json.Unmarshal(jsonBlob, &i)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(i)
	// fmt.Println(i.ActionBody.Authentication)
	// fmt.Println(i.ActionBody.Functions[0].FunctionCall)
	// fmt.Println(i.ActionBody.Functions[0].FunctionParams)
	// fmt.Println(i.ActionBody.Functions[0].FunctionParams["fields"].([]interface{})[0])
	// fmt.Println("->", i.ActionBody.Functions[0].FunctionParams["id"])
	// fmt.Println(i.ActionBody.Functions[1].FunctionCall)
	// fmt.Println(i.ActionBody.Functions[1].FunctionParams)
	// fmt.Println("-> >", i.ActionBody.Functions[1].FunctionParams["fields"].([]interface{})[0])
	// fmt.Println(i.ActionBody.Functions[1].FunctionParams["id"])
	// fmt.Println(i.ActionBody.Returns)

	// FuncsStorage = map[string]interface{}{
	// 	"funcA": functions.Test,
	// }

	// FuncsStorage["funcB"] = functions.Test2
	// FuncsStorage["funcC"] = functions.Test3

	// CallFunc("funcA", 1, 2, 3)
	// //prntFA := resFA.(int)
	// //fmt.Println(prntFA)
	// CallFunc("funcB", 2, 2, 3)

	// resFC, _ := CallFunc("funcC", 4, 5, 6)
	// fmt.Println(resFC)

	// type Pessoa struct {
	// 	Nome  string `json:"nome"`
	// 	Idade int    `json:"idade"`
	// }
	// x := Pessoa{
	// 	Nome:  "Tomás",
	// 	Idade: 18,
	// }

	// res, err := generate.CriarRegisto(&x)
	// if err != nil {
	// 	fmt.Println("Error")
	// }
	// fmt.Println(string(res))

	// fmt.Println("-------------------------------")
	// generate.CreateResolver()
}

// CallFunc - Calls the function with by the name specified in funcName
func CallFunc(funcName string, params ...interface{}) (interface{}, error) {
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
