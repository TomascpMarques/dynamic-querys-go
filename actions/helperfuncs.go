package actions

import (
	"errors"
	"reflect"
)

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

// ConvertToPrimitives -
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

// GetCalledFuncs -
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

// RunFunctionsGetReturns -
func RunFunctionsGetReturns(functions []FunctionPath) ([]interface{}, error) {
	returns := make([]interface{}, len(functions))
	// Iterates through all the queryed functions in the request
	for k := range functions {
		funcParams, err := ConvertToPrimitives(functions[k].FunctionParams)
		if err != nil {
			DQGLogger.Panicf("Error: %s", err)
		}
		res, err := CallFunc(functions[k].FunctionCall, funcParams)
		if err != nil {
			DQGLogger.Println("Error: Either bad params or called function error")
			continue
		}
		returns[k] = res
	}
	if len(returns) == 0 {
		return nil, errors.New("Error: functions returned nothing")
	}

	return returns, nil
}
