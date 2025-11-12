package authenticationmiddleware

import (
	"context"
	"net/http"

	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

func AuthorizedInvoke(fn func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		authHeader := req.Headers["Authorization"]
		if authHeader == "" {
			return response.LambdaResponse(http.StatusUnauthorized, map[string]any{"errorCode": 1003}, "Unauthorized"),nil
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := utils.VerifyJwt(tokenString)
		if err != nil{
			return response.LambdaResponse(http.StatusUnauthorized, map[string]any{"errorCode": 1003}, "Unauthorized"),nil
		}


		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return response.LambdaResponse(http.StatusUnauthorized, map[string]any{"errorCode": 1003}, "Unauthorized"),nil
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return response.LambdaResponse(http.StatusUnauthorized, map[string]any{"errorCode": 1003}, "Unauthorized"),nil
		}

		role, ok := claims["role"].(string)
		if !ok {
			return response.LambdaResponse(http.StatusUnauthorized, map[string]any{"errorCode": 1003}, "Unauthorized"),nil
		}

		email, ok := claims["email"].(string)
		if !ok {
			return response.LambdaResponse(http.StatusUnauthorized, map[string]any{"errorCode": 1003}, "Unauthorized"),nil
		}

		flat, ok := claims["flat"].(string)
		if !ok {
			return response.LambdaResponse(http.StatusUnauthorized, map[string]any{"errorCode": 1003}, "Unauthorized"),nil
		}

		ctx = context.WithValue(ctx, utils.UserIDKey, userID)
		ctx = context.WithValue(ctx, utils.UserEmailKey, email)
		ctx = context.WithValue(ctx, utils.UserFlatKey, flat)
		ctx = context.WithValue(ctx, utils.UserRoleKey, role)

		return fn(ctx, req)
	}
}