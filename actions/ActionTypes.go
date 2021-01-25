package actions

// Endpoint - Represent the function called in the action schema, its name and parameters
type Endpoint struct {
	FuncName string        // Function name to call
	Params   []interface{} // Parameters to be passed into the function
}

// BodyContents - Represents the extende structure of the action body
type BodyContents struct {
	ActionBody     string     // Body of the action
	Authentication [][]string // JWT auth tokens
	FuncCalls      [][]string // Function names passed in the schema
	FuncArgs       [][]string // Passed params in the schema
	FuncsContent   []string   // Content of the " funcs: " section of the schema
}
