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

// FuncMap - Map with the available functions to call
type FuncMap map[string]interface{}

// FuncsStorage -
var FuncsStorage = FuncMap{
	//*
	//* !!! STATE YOUR FUNCS HERE !!! **//
	//* You can implement your functions where-ever,
	//* as long the program can get through the path to them
	//* You can implement them in the folder funcendpoints, if you want everythin in the same place
	//*
	"ReverseString":        funcendpoints.ReverseString,
	"ReverseStringBool":    funcendpoints.ReverseStringBool,
	"TakeAnInterfaceArray": funcendpoints.TakeAInterfaceArray,
	"TakeAMap":             funcendpoints.TakeAMap,
}

// Handler - Handles all of the requests coming into the server
func Handler(rw http.ResponseWriter, r *http.Request) {
	DQGLogger.Println("New action recived, checking and parssing...")
	// Gets the body from the request
	requestBody, _ := ioutil.ReadAll(r.Body)
	action := strings.TrimSpace(string(requestBody))

	// Checks is the request sent contains a action
	if err := CheckRequestIsAction(action); err != nil {
		DQGLogger.Println("Error:", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"error":"The request sent was not a valid action"}`))
		return
	}

	// Gets the contents of the action, such as funcs:(function calls and its parameters) and auth:
	actionContents, err := ParseActionContents(action)
	if err != nil {
		DQGLogger.Println("Error:", err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"error":"Not able to parse one or more content-parts of the action"}`))
		return
	}

	DQGLogger.Println("Endpoints called: ", actionContents.FuncCalls)

	// Parsses the sent action values to go usable data types
	functionCalMap, err := ParseActionBody(`"\w+":$`, actionContents)
	if err != nil {
		DQGLogger.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"error":"Not able to convert to primitive"}`))
		return
	}

	// Runs the functions specified in the action
	results, err := RunFunctionsGetReturns(functionCalMap)
	if err != nil {
		DQGLogger.Println("Error: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"error":"Error runing the actions functions"}`))
		return
	}

	// Json encodes the functions results
	send, err := json.Marshal(results)
	if err != nil {
		DQGLogger.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(`{"error":"Unable to marshal the actions results"}`))
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write([]byte(send))

	DQGLogger.Println("No errors, all good")
	return
}
