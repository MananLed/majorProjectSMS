package response

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type response struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	ErrorCode int         `json:"errorcode,omitempty"`
	Data      any         `json:"data,omitempty"`
}

func SuccessResponse(data any , message string, code int) events.APIGatewayProxyResponse{

	response:= response{
		Status: "Success",
		Message: message,
		Data: data,
	}

	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body: string(body),
	}
}

func ErrorResponse(statusCode int, errMessage string, code int) events.APIGatewayProxyResponse{
	response:=response{
		Status: "fail",
		Message: errMessage,
		ErrorCode: code,
	}

	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body: string(body),
	}
}