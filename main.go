package main

import (
	"DynamicQuerysGo/actions"
	"encoding/json"
	"fmt"
)

func main() {
	jsonBlob := []byte(`{"action": {"auth": "mh354kh2vhvqyçyavq5yq8q5yyjuqqy5","func": [{"call": "GetRegisto","params": {"id": "Registo123","fields": ["name", "id"]}},{"call": "GetStuff","params": {"id": "123","fields": ["id"]}}],"returns": ["success","id","nome"] }}`)

	var i actions.Action
	err := json.Unmarshal(jsonBlob, &i)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(i)
	fmt.Println(i.ActionBody.Authentication)
	fmt.Println(i.ActionBody.Functions[0].FunctionCall)
	fmt.Println(i.ActionBody.Functions[0].FunctionParams)
	fmt.Println(i.ActionBody.Functions[0].FunctionParams["fields"].([]interface{})[0])
	fmt.Println(i.ActionBody.Functions[0].FunctionParams["id"])
	fmt.Println(i.ActionBody.Functions[1].FunctionCall)
	fmt.Println(i.ActionBody.Functions[1].FunctionParams)
	fmt.Println(i.ActionBody.Functions[1].FunctionParams["fields"].([]interface{})[0])
	fmt.Println(i.ActionBody.Functions[1].FunctionParams["id"])
	fmt.Println(i.ActionBody.Returns)
}
