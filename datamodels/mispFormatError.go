package datamodels

/*{
    "saved": false,
    "name": "Could not add Attribute",
    "message": "Could not add Attribute",
    "url": "\/attributes\/add",
    "errors": {
        "type": [
            "Options depend on the selected category."
        ],
        "value": [
            "Value not in the right type\/format. Please double check the value or select type \"other\"."
        ]
    }
}*/

type MispFormatError struct {
	Saved   bool          `json:"saved"`
	Name    string        `json:"name"`
	Message string        `json:"message"`
	URL     string        `json:"url"`
	Errors  ErrorsPattern `json:"errors"`
}

type ErrorsPattern struct {
	Type  []string `json:"type"`
	Value []string `json:"value"`
}
