package lib

type RequestBody struct {
	Platform   string                 `json:"platform"`
	CliVersion string                 `json:"cli_version"`
	Type       string                 `json:"type"`
	Data       map[string]interface{} `json:"data"`
}

type Document struct {
	Platform   string                 `json:"platform"`
	CliVersion string                 `json:"cli_version"`
	IpAddr     string                 `json:"ip"`
	Type       string                 `json:"type"`
	Data       map[string]interface{} `json:"data"`
}

func MapToDocument(r RequestBody, ip string) Document {
	return Document{
		Platform:   r.Platform,
		CliVersion: r.CliVersion,
		IpAddr:     ip,
		Type:       r.Type,
		Data:       r.Data,
	}
}

var BadRequest = map[string]interface{}{
	"status":  400,
	"title":   "Bad Request",
	"message": "Some fields might miss or malformed.",
}

var InternalServerError = map[string]interface{}{
	"status":  500,
	"title":   "Internal Server Error",
	"message": "Something went wrong.",
}

var MethodNotAllowed = map[string]interface{}{
	"status": 405,
	"title":  "Method not allowed",
}
