package shared

const MainComponentName = "Middleman Component"
const NetworkingComponentName = "Networker"
const RoutingComponentName = "Router"

const ListeningAddress = "127.0.0.1"
const ListeningPort = "3333"
const EndOfMessageToken = "<eom>"

// ---- Default Server Messages ----

// DefaultSuccessResponsePayload The Default success response for the server.
// Don't forget to set the target component and time sent.
const DefaultSuccessResponsePayload = `{
		"action": "response",
		"args": [
			{
				"type": "string",
				"value": "The server has received your request and will process it accordingly."
			},
			{
				"type": "int",
				"value": "1"
			}
		]
	}`

const RegisterSuccessResponsePayload = `{
		"action": "response",
		"args": [
			{
				"type": "string",
				"value": "The server has successfully registered you."
			},
			{
				"type": "int",
				"value": "0"
			}
		]
	}`

// InvalidRequestResponsePayload The Default failure response for the server.
// Don't forget to set the target component and time sent.
const InvalidRequestResponsePayload = `{
		"action": "error",
		"args": [
			{
				"type": "string",
				"value": "Invalid Request, please verify your format."
			},
			{
				"type": "int",
				"value": "355"
			}
		]
	}`

const NoMatchFoundResponsePayload = `{
		"action": "error",
		"args": [
			{
				"type": "string",
				"value": "There was no match for the request action.\nPlease verify that your action argument types and numbers are correct."
			},
			{
				"type": "int",
				"value": "354"
			}
		]
	}`
