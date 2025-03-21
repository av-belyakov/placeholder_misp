package datamodels

/*
Пример ошибок от MISP
{
    "saved": false,
    "name": "Could not add Attribute",
    "message": "Could not add Attribute",
    "url": "\/attributes\/add",
    "errors": {
        "Event":{
            "uuid":[
                "Please provide a valid RFC 4122 UUID",
                ],
            },
        },
        "type": [
            "Options depend on the selected category."
        ],
        "value": [
            "Value not in the right type\/format. Please double check the value or select type \"other\"."
        ]
    }
}*/

type MispFormatError struct {
	Errors  map[string]interface{} `json:"errors"`
	Name    string                 `json:"name"`
	Message string                 `json:"message"`
	URL     string                 `json:"url"`
	Saved   bool                   `json:"saved"`
}
