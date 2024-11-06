package shared

const MainComponentName = "Middleman Component"
const NetworkingComponentName = "Networker"

const ListeningPort = "3333"
const EndOfMessageToken = "<eom>"

// ---- Default Server Messages ----

// DefaultSuccessReponse The Default success response for the server.
// Don't forget to set the target component and time sent.
const DefaultSuccessReponse = `{
		"requester": {
			"title": "Middleman Component",
			"version": "1.0.0",
			"capabilities": []
		},
		"target": {
			"title": "",
			"version": "N/A",
			"capabilities": []
		},
		"time_sent": "2024-11-06T12:58:30Z",
		"payload": {
			"action": "Reply",
			"args": [
				{
					"type": "string",
					"value": "The server has received your request and will process it accordingly."
				}
			]
		}
	}`

// DefaultFailureResponse The Default failure response for the server.
// Don't forget to set the target component and time sent.
const DefaultFailureResponse = `{
		"requester": {
			"title": "Middleman Component",
			"version": "1.0.0",
			"capabilities": []
		},
		"target": {
			"title": "",
			"version": "N/A",
			"capabilities": []
		},
		"time_sent": "2024-11-06T12:58:30Z",
		"payload": {
			"action": "Reply",
			"args": [
				{
					"type": "string",
					"value": "Invalid Request, please verify your format."
				}
			]
		}
	}`
