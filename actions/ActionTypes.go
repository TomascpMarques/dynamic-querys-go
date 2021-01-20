package actions

// Action - Defines the action as the body of the json action field,
// 		    which as the body of the recieved action.
type Action struct {
	ActionBody ActionBody `json:"action"`
}

// ActionBody - Body of the action in itself, has all the fields
//			   that are required for its decoding/functionality.
type ActionBody struct {
	Authentication string         `json:"auth,-"`
	Functions      []FunctionPath `json:"func"`
	Returns        []string       `json:"returns,omitempty"`
}

// FunctionPath - Defines the function to be called (it's path/name),
//				  and it's parameters.
type FunctionPath struct {
	FunctionCall   string                 `json:"call"`
	FunctionParams []interface{} `json:"params,omitempty"`
}
