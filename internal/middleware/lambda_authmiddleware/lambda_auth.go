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
			return response.ErrorResponse(http.StatusUnauthorized, "Unauthorized", 1003), nil
		}

		tokenString := authHeader[len("Bearer "):]
		token, err := utils.VerifyJwt(tokenString)
		if err != nil{
			return response.ErrorResponse(http.StatusUnauthorized, "Unauthorized", 1003), nil
		}


		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return response.ErrorResponse(http.StatusUnauthorized, "Unauthorized", 1003), nil
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return response.ErrorResponse(http.StatusUnauthorized, "Unauthorized", 1003), nil
		}

		role, ok := claims["role"].(string)
		if !ok {
			return response.ErrorResponse(http.StatusUnauthorized, "Unauthorized", 1003), nil
		}

		email, ok := claims["email"].(string)
		if !ok {
			return response.ErrorResponse(http.StatusUnauthorized, "Unauthorized", 1003), nil
		}

		flat, ok := claims["flat"].(string)
		if !ok {
			return response.ErrorResponse(http.StatusUnauthorized, "Unauthorized", 1003), nil
		}

		ctx = context.WithValue(ctx, utils.UserIDKey, userID)
		ctx = context.WithValue(ctx, utils.UserEmailKey, email)
		ctx = context.WithValue(ctx, utils.UserFlatKey, flat)
		ctx = context.WithValue(ctx, utils.UserRoleKey, role)

		return fn(ctx, req)
	}
}