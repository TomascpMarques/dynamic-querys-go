package actions

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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



	// funcNames[0][0][1:len(funcNames[0][0])-2]
	fmt.Println("| ->>", auth, "|", funcNames, funcArgs)
	rw.Write([]byte(actionBody))
}
