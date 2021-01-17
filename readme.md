# DynamicQuerysGo

## __Intro__

_DynamicQuerysGo_ is a go library with the objective to create a simple graphQL like solution to __easily implement CRUD functionalitty in REST APIs__.

So to achieve this goal a pattern to __simulate graph-like appearance and behavior__, and easy to use was created. Its sent by the API consumer as a POST request, but the request action is defined by the function called, __not the HTTP verb used__. The parameters required for the functions to work are given by the consumer within the body of the action.

<br><hr>

## __Concepts__

DynamicQuerysGo only does what you implement as functions, so words to define actions in the beginig of the request, such as query, mutation and subscription were droped, and the word __action__ is used to signify the begining of the request.

A dynamic request is based on three main parts, which are:
<div style="margin-left:25px">
    1. The actions body. <br>
    2. The functions called in the body. <br>
    3. And the return values, after action completion.
</div>

<br><hr>

## __Examples__

<br>

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
