package actions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// DQGLogger - Logger for DynamicQuerysGo processes (can be stdout or a io.writer to a file)
var DQGLogger = log.New(os.Stdout, "DynamicQuerysGo [*] ", log.LstdFlags)

// FuncMap -
type FuncMap map[string]interface{}

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

	// Sets up the action, populate it with the r.body content
	var i Action
	err := json.Unmarshal(requestBody, &i)
	if err != nil {
		DQGLogger.Println(err)
	}

	DQGLogger.Printf("New action to handle | auth: %s | Called functions: %v|", i.ActionBody.Authentication, GetCalledFuncs(i.ActionBody.Functions))

	returns, err := RunFunctionsGetReturns(i.ActionBody.Functions)
	if err != nil {
		DQGLogger.Println(err)
		return
	}

	result, _ := json.MarshalIndent(&returns, "", "\t")
	rw.Write(result)
}
