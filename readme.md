# DynamicQuerysGo / GoActions

## __Intro__

_DynamicQuerysGo_ is a simple API project written in go, whith the intent to simulate a GraphQL like experience, in terms of querys and mutations. Here a request is called an action, which represents both an mutation and an query, the main purpose of these _actions_ is to resolve actions/requests to function calls, in a quick and simple way.

## __What I plan to add__ - (updated)

1. __DONE!__ Reduce the request/action foot-print, sent to the server. Make it more readable, like graph, or a docker compose file, type of readable.

1. ~~Add authentication middleware, in the format of JSON WebTokens.~~

1. __DONE!__ Make my own request parser to suport objective one.

## __Things to keep in mind__ - v.1.1-Release

Settings:

1. The port where the server is located defaults to 8000, but it also reads from the enviorment variable _ENV_GOACTIONS_PORT_

Since DynamicQuerysGo is based arround a custom request schema (like in graph), conssider the following:

1. Any form of rquest authentication, should be implemented by the user, the idea is to use JWT, but the auth fiedl only takes strings

1. The beginning of the request sholud have the `action:` keyword, other wise the no functions will be called, and the action body will not be parssed.

1. The function calls are deffined by the format `"funcName":` and the parameters to be passed are written in the lines bellow, one line per each parameter, and the paremeters are declared similarly as json values ending the line with a semicolon, ex:

    | Variable type      | Variable value    |
    | :----------------  | :---------------- |
      Integers           | 1234,
      Strings            | "Muck Nuck",
      Float64            | 142.6356,
      Booleans           | true / false
      map\[string\]interface{} (to decode json) | "{\\"name\\":\\"Muck Nuck\\"}", OR "{"name":123124}",
      \[ \]interface{} (multi/single type arrays)|  \[1231,14.13,true,"asdasd",\],

    
1. The action schema should be written in the following these steps:
    * First line is for the keyword `action:`.
    * Second line can ommit the `auth:` field, but shouldn't _(unlless you don't want any authentication through JWTs)_.
    * The function calls are specified by `"functionName":` and the next lines are the arguments, one line per each argument passed.
    * The only indentation to worry about is in the paremeter passing for functions, you could write the action with no indentation but the parameters need to be in their own singular lines.

## __Examples__

__Action calling function ReverseString:__ _(action schema oriented)_

    action:
        auth: "JWT EXAMPLE"
        funcs:
            "ReverseString":
                "Hello GO",

__OR__

    action:
    auth: "JWT EXAMPLE"
    funcs:
    "ReverseString":
    "Hello GO",

__Action calling multiple funcs:__ _(action schema oriented)_

    action:
        auth: "JWT EXAMPLE"
        funcs:
            "ReverseString":
                "Hello GO",
            "ReverseStringBool":
                true,
                "Hello GO", 
            "TakeAnInterfaceArray":
                [1231,14.13,true,"Muck","Nuck","{\\"name\\":123124}",],
            "TakeAMap":
                "{"age":43}",
__Result:__ _(json)_
    
    {
        "ReverseString": [
            {
                "reverssed": "OG olleH"
            }
        ],
        "ReverseStringBool": [
            {
                "reversse": true,
                "reverssed": "OG olleH"
            }
        ],
        "TakeAMap": [
            {
                "age": 43
            }
        ],
        "TakeAnInterfaceArray": [
            [
                1231,
                14.13,
                true,
                "asdasd",
                "asdasd",
                {
                    "name": 123124
                }
            ]
        ]
    }
    
