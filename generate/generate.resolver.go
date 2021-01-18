package generate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// HookList -
type HookList struct {
	Hooks   []string `json:"functions"`
	Package string   `json:"package"`
}

// CreateResolver -
func CreateResolver() error {
	defineFuncFileContent, err := ioutil.ReadFile("generate/define.funcs.jsonc")
	if err != nil {
		return err
	}
	fmt.Println(string(defineFuncFileContent))

	var cont HookList
	err = json.Unmarshal(defineFuncFileContent, &cont)
	if err != nil {
		return err
	}

	fmt.Println(cont.Hooks, cont.Package)

	// FuncCalls = map[string]interface{}{
	// 	"funcA": functions.Test,
	// }

	test, _ := ioutil.ReadFile("generate/genContents/head.txt")

	funcStoreGen := "\n\tFuncCalls = map[string]interface{} {\n"
	contGen := string(test) + funcStoreGen
	for _, v := range cont.Hooks {
		contGen += "\t\t\"" + v + "\": " + cont.Package + "." + v + ",\n"
	}
	contGen += "\t}"

	fmt.Println(contGen)

	return nil
}
