package actions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/TomascpMarques/dynamic-querys-go/funcendpoints"
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
	"ReverseString": funcendpoints.ReverseString,
}

// Handler - Handles all of the requests coming into the server
func Handler(rw http.ResponseWriter, r *http.Request) {
	DQGLogger.Println("New action recived, checking and parssing...")
	// Gets the body from the request
	requestBody, _ := ioutil.ReadAll(r.Body)
	action := strings.TrimSpace(string(requestBody))

	if err := CheckRequestIsAction(action); err != nil {
		DQGLogger.Println("Error:", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Error: The request sent was not a valid action`))
		return
	}

	actionContents, err := ParseActionContents(action)
	if err != nil {
		DQGLogger.Println("Error:", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`Error: Not able to parse one or more content-parts of the action`))
		return
	}
	DQGLogger.Println("Endpoints called: ", actionContents.FuncCalls)

	functionCalMap, err := ParseActionBody(`"\w+":$`, actionContents)
	if err != nil {
		DQGLogger.Println(err)
		return
	}

	results, err := RunFunctionsGetReturns(functionCalMap)
	if err != nil {
		DQGLogger.Println("Error: ", err)
		return
	}
	send, err := json.Marshal(results)
	if err != nil {
		DQGLogger.Println(err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(send))

	DQGLogger.Println("No errors, all good")
	return
}
