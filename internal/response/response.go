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

	// w.Header().Set("content-Type","application/json")
	// w.WriteHeader(code)
	// json.NewEncoder(w).Encode(response)

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

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(statusCode)
	// json.NewEncoder(w).Encode(response)

	body, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body: string(body),
	}
}

// func LambdaResponse(code int, data any, message string) events.APIGatewayProxyResponse{
// 	body, _ := json.Marshal(map[string]any{
// 		"message": message,
// 		"data": data,
// 	})
// 	return events.APIGatewayProxyResponse{
// 		StatusCode: code,
// 		Body: string(body),
// 	}
// }