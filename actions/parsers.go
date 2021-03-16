package actions

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ParseActionBody - Parsses the action body text into usable data
func ParseActionBody(regex string, actionContents BodyContents) ([]Endpoint, error) {
	// `"\w+":$`
	// Defines the regex representing function calls
	functionCallRegexp := regexp.MustCompile(regex)
	functionCalMap := make([]Endpoint, 0)

	// Gets the list of function available on the function map,
	// that we can call
	ableToCall := reflect.ValueOf(FuncsStorage).MapKeys()
	call := ""
	for _, v := range ableToCall {
		call += v.String() + " "
	}

	// Iterates through the funcs: part of the action, and extracts it's functions allong with it's parameters
	for k, v := range actionContents.FuncsContent {

		check, err := CheckTypeAndConvert(v)
		if err != nil {
			DQGLogger.Println("Unable to convert to golang data type")
			return nil, err
		}
		// Skips not existing function calls
		if !strings.Contains(call, reflect.ValueOf(check).String()) {
			continue
		}
		// Checks if the current line is not a function call
		// if its not a function call it will search and get the following lines
		// as the previous function call parameters.
		if len(functionCallRegexp.FindAllString(v, -1)) != 0 {
			// Sets up the parameters list
			params := make([]interface{}, 0)
			for _, j := range actionContents.FuncsContent[k+1:] {

				// Skips iteration if curret line is a function call
				// after skiping, all the paremeters found
				// are attributted to the previous function call
				if string(j[len(j)-1]) == ":" {
					break
				}
				// Converts the current line content, into its appropriate go data type
				res, err := CheckTypeAndConvert(j)
				if err != nil {
					return nil, err
				}
				// Appends the converted value to the parameter array
				params = append(params, res)
			}
			// Checks if the number of params insside params is equal to the requeired number to call the function
			pnum, err := GetFunctionParamsNum(reflect.ValueOf(FuncsStorage[v[1:len(v)-2]]))
			if err != nil || len(params) != pnum {
				return nil, errors.New("bad parameters")
			}
			// If all went well, a function call and its params will be appended in the functionCalMap, and returned
			functionCalMap = append(functionCalMap, Endpoint{
				FuncName: v[1 : len(v)-2],
				Params:   params,
			})
		}
	}
	return functionCalMap, nil
}

/*
	* Regexp Patterns:
		? 1. "[a-zA-Z0-9_ ]+",+ -->  strings insside quotations (include spaces), check multiple times.
		? 2. \d+\.\d+, 		    -->  floats, check multiple times.
		? 3. \d+,+ 			    -->  integers, check multiple times.
		? 4. true,+|false,+     -->  booleans, check multiple times.
		? 5. \{.\},+		    -->  json-like strings ( "{\"name\":\"Golang\"}", or "{"age":43}", ), check multiple times.
		? 6. \[.\],+	  	    -->  arrays of interfaces, check multiple times.
		? 7. ^b".+"$		    -->  strings that will be translated to byte arrays
		? 8. ^b\{.+\}$			-->  json values that will be translated to byte arrays
*/

/*
CheckTypeAndConvert - Takes a string and compares it to a number of regex patterns, each representing golang data types,
	if the pattern and string match, the string will be converted to the represented data type.
	The order of regex comparisson is important in some cases, such as in strings, json-like strings ("{\"name\":\"Golang\"}").
*/
func CheckTypeAndConvert(word string) (interface{}, error) {
	wordNoCotations := word[1 : len(word)-1]
	wordNoSemicolon := word[:len(word)-1]
	// Checks string to see if its possible to convert into an array of interfaces
	if len(regexp.MustCompile(`^\[.+\],$`).FindAllString(word, -1)) != 0 {
		//temp := make([]interface{}, 0)
		// Extracts all available data types, by their corresponding pattern
		allTypesRegex := `".+",+|\d+\.\d+,+|\d+,+|true,+|false,+|\{.+\},+|\[.+\],+`
		extractedValues := regexp.MustCompile(allTypesRegex).FindAllStringSubmatch(wordNoCotations, -1)

		// Creates list to store the converted strings
		convertedList := make([]interface{}, 0)
		for _, v := range extractedValues {
			// Initiates the converssion process
			// Cals its parent function to make use of the already implemented regexp check/convert
			converted, err := CheckTypeAndConvert(v[0])
			if err != nil {
				return nil, errors.New("error converting")
			}

			// Adds the param converted from a string into a list to be returned
			convertedList = append(convertedList, converted)
		}
		return convertedList, nil
	}
	// CHecks and converts a string to byte array
	if len(regexp.MustCompile(`^b".+"$`).FindAllString(wordNoSemicolon, -1)) != 0 {
		cnvt := []byte(wordNoSemicolon[2 : len(wordNoSemicolon)-1])
		return cnvt, nil
	}
	// CHecks and converts a json string to a byte array
	if len(regexp.MustCompile(`^b\{.+\}$`).FindAllString(wordNoSemicolon, -1)) != 0 {
		cnvt := []byte(wordNoSemicolon[1:])
		return cnvt, nil
	}
	// Checks and converts a string to a integer
	if len(regexp.MustCompile(`^\d+$`).FindAllString(wordNoSemicolon, -1)) != 0 {
		cnvt, err := strconv.Atoi(wordNoSemicolon)
		if err != nil {
			return nil, errors.New("error converting to int")
		}
		return cnvt, nil
	}
	// Checks and converts a string to a float64
	if len(regexp.MustCompile(`^\d+\.\d+$`).FindAllString(wordNoSemicolon, -1)) != 0 {
		cnvt, err := strconv.ParseFloat(wordNoSemicolon, 64)
		if err != nil {
			return nil, errors.New("error converting to float64")

		}
		return cnvt, nil
	}
	// Checks and converts a string to an interface array
	if len(regexp.MustCompile(`\{.+\}`).FindAllString(wordNoSemicolon, -1)) != 0 {
		test := regexp.MustCompile(`\\`).ReplaceAllLiteralString(word[0:len(word)-1], "")
		var jsonBlob = []byte(test)

		var cnvt map[string]interface{}
		err := json.Unmarshal(jsonBlob, &cnvt)
		if err != nil {
			return nil, errors.New("error converting to map[string]interface{}")
		}
		return cnvt, nil
	}
	// Checks and properly formatts the string (pops the <""> and the <,>)
	if len(regexp.MustCompile(`^".+"$`).FindAllString(wordNoSemicolon, -1)) != 0 {
		cnvt := word[1 : len(word)-2]
		return cnvt, nil
	}
	// Checks and converts a string to a
	if len(regexp.MustCompile(`^true$|^false$|true|false`).FindAllString(wordNoSemicolon, -1)) != 0 {
		cnvt, err := strconv.ParseBool(wordNoSemicolon)
		if err != nil {
			return nil, errors.New("error converting to boolean")
		}
		return cnvt, nil
	}
	DQGLogger.Println("No regexp was hit!")
	return nil, errors.New("not able to translate into go primitive (bad parameter or no regex hit)")
}

//b\{.+\},+

// ParseActionContents - Separates the action into its dedicated content parts
func ParseActionContents(request string) (BodyContents, error) {
	var result BodyContents
	result.ActionBody = strings.TrimSpace(request[7:])
	if len(result.ActionBody) == 0 {
		return BodyContents{}, errors.New("error extracting the action body")
	}
	result.Authentication = regexp.MustCompile(`auth: ".+"|auth: ".+" +`).FindAllStringSubmatch(result.ActionBody, -1)
	if len(result.Authentication) == 0 {
		DQGLogger.Println("! Attention no auth token was sent !")
		result.Authentication = make([][]string, 0)
	}
	result.FuncCalls = regexp.MustCompile(`"\w+":\n|"\w+":\s+\n`).FindAllStringSubmatch(result.ActionBody, -1)
	if len(result.FuncCalls) == 0 {
		DQGLogger.Println(errors.New("error extracting function calls"))
	}
	result.FuncArgs = regexp.MustCompile(`"[a-zA-Z0-9_ ]+",+|\d+\.\d+,+|\d+,+|true,+|false,+|\{.+\},+|\[.+\],+|b".+",|b\{.+\},`).FindAllStringSubmatch(result.ActionBody, -1)
	if len(result.FuncArgs) == 0 {
		DQGLogger.Println(errors.New("error extracting the functions arguments"))
	}
	result.FuncsContent = regexp.MustCompile(`"[a-zA-Z0-9_ ]+",+|\d+\.\d+,+|\d+,+|true,+|false,+|\{.+\},+|\[.+\],+|".+",|".+":|b\{.+\},+|b".+",|b\{.+\},`).FindAllString(result.ActionBody, -1)
	if len(result.FuncsContent) == 0 {
		return BodyContents{}, errors.New("error extracting the <funcs:> content")
	}

	return result, nil
}
