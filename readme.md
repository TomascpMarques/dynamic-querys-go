# DynamicQuerysGo

## __Intro__

_DynamicQuerysGo_ is a simple project written in go, whith the intent to simulate a GraphQL like experience, in terms of querys and mutations. Here a request is called an action, which represents both an mutation and an query, the main purpose of these _actions_ is to resolve actions/requests to function calls, in a quick and simple way.

## __Examples__

__Action of Read/Get type:__

    {
        "action": {
            "auth": "mh354kh2vhvqyeoavq5yq8q5yyjuqqy5",
            "func": [
                {
                    "call": "GetRegisto",
                    "params": {
                        "id": "Registo123",
                        "fields": ["name", "id"]
                    }
                }
            ],
            "returns": [
                "success",
                "id",
                "nome"
            ] // or nil
        }
    }
