package actions

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// CheckRequestIsAction -
func CheckRequestIsAction(action string) error {
	if len(regexp.MustCompile(`^action:\n|^action:\s+\n`).FindAllString(action, -1)) == 0 {
		return errors.New("Request sent is not an action")
	}
	return nil
}

// CheckGivenParams - checks if the number of parmeters is the correct amount.
func CheckGivenParams(params []interface{}, numParams int) error {
	if len(params) != numParams {
		return errors.New("The number of params is insufficient")
	}
	return nil
}

// GetFunctionParamsNum - Gets the given function number of parameters, and returns it or an error.
func GetFunctionParamsNum(function reflect.Value) (int, error) {
	if function.Type().NumIn() == 0 {
		return 0, errors.New("The given function takes zero parameters")
	}
	return function.Type().NumIn(), nil
}

// ParseParamsIntoRVArray - Gets values from the params array, and puts their reflect value insside a []reflect.Value, and returns it.
func ParseParamsIntoRVArray(params []interface{}) ([]reflect.Value, error) {
	// Gets params from an interface array, and puts them in a []reflect.Value,
	// to be used in the call function, of the reflect package
	funcParams := make([]reflect.Value, len(params))
	for k, param := range params {
		funcParams[k] = reflect.ValueOf(param)
	}
	if len(funcParams) == 0 {
		return nil, errors.New("Error parssing the functions parameters, from given array")
	}

	return funcParams, nil
}

// CallFunc - Calls the function with by the name specified in funcName
func CallFunc(funcName string, params []interface{}) (interface{}, error) {
	// Gets function as reflect.Value to perform reflection,
	// to know things such as number of parameters
	function := reflect.ValueOf(FuncsStorage[funcName])

	numParams, err := GetFunctionParamsNum(function)
	if err != nil {
		return nil, err
	}

	// Checks if the passed parameters are more or less than the ones required
	err = CheckGivenParams(params, numParams)
	if err != nil {
		return nil, errors.New("The number of params is insufficient")
	}

	// Gets al the parameters passed in params
	// to be used in reflect.Call, as the called functions parameters
	funcParams, err := ParseParamsIntoRVArray(params)
	if err != nil {
		return nil, err
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

// RunFunctionsGetReturns - Runs the given functions in []FunctionPath, and returns the resulting function-call values.
func RunFunctionsGetReturns(functionCalMap []Endpoint) (map[string]interface{}, error) {
	results := make(map[string]interface{}, 1)
	for k, v := range functionCalMap {
		res, err := CallFunc(v.funcName, v.params)
		if err != nil {
			return nil, err
		}
		if _, ok := results[v.funcName]; ok {
			DQGLogger.Printf("Function call repeat, sending second function call as -> %v <-\n", v.funcName+"_Redo"+fmt.Sprint(k))
			results[v.funcName+"_Redo"+fmt.Sprint(k)] = res
		}
		results[v.funcName] = res
	}
	return results, nil
}
