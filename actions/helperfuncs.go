package actions

import (
	"errors"
	"reflect"
)

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

// ConvertToPrimitives - Converts the given values in x (extracted from JSON), to equivalent primitive golang data types.
func ConvertToPrimitives(x []interface{}) ([]interface{}, error) {
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

// GetCalledFuncs - Gets the functions called in the current Dynamic-Action.
func GetCalledFuncs(array []FunctionPath) string {
	calledFunctions := ""
	for _, v := range array {
		calledFunctions += v.FunctionCall + "; "
	}
	if calledFunctions == "" {
		return "Error: No functions called"
	}

	return calledFunctions
}

// RunFunctionsGetReturns - Runs the given functions in []FunctionPath, and returns the resulting function-call values.
func RunFunctionsGetReturns(functions []FunctionPath) ([]interface{}, error) {
	returns := make([]interface{}, len(functions))
	// Iterates through all the queryed functions in the request
	for k := range functions {
		// Conver the given JSON function params, into golang primitives
		funcParams, err := ConvertToPrimitives(functions[k].FunctionParams)
		if err != nil {
			DQGLogger.Panicf("Error: %s", err)
		}

		// Calls the function by the name specified insside the string
		res, err := CallFunc(functions[k].FunctionCall, funcParams)
		if err != nil {
			DQGLogger.Println("Error: Either bad params or called function error")
			continue
		}
		returns[k] = res
	}
	// checks for no return values
	if len(returns) == 0 {
		return nil, errors.New("Error: functions returned nothing")
	}

	return returns, nil
}
