package repository

import (
	"context"

	"errors"
	"fmt"


	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CredentialRepository struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type CredentialRepositoryInterface interface {
	DeleteUserByIDAndRole(id string, role model.UserRole) error
}

func NewCredentialRepository(ddbClient *dynamodb.Client, tableName string) *CredentialRepository {
	return &CredentialRepository{DynamoDbClient: ddbClient, TableName: tableName}
}

func (r *CredentialRepository) DeleteUserByIDAndRole(id string, role model.UserRole) error {
	// var exists bool

	// query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND role = $2)`

	// err := r.db.QueryRow(query, id, string(role)).Scan(&exists)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return fmt.Errorf("failed to check user existence: %v", err)
	// }

	// if !exists {
	// 	return errors.New("user not found")
	// }

	// query = `DELETE FROM users WHERE id = $1 AND role = $2`

	// _, err = r.db.Exec(query, id, string(role))

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return fmt.Errorf("failed to delete user: %v", err)
	// }

	// return nil

	fetchUserStatement := "SELECT * FROM " + r.TableName + " WHERE PK = ? AND SK = ?"

	result, err := r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(fetchUserStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ROLE#" + string(role)},
			&types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return err
	}

	if len(result.Items) == 0 {
		return errors.New("user not found")
	}

	type User struct {
		PK           string `dynamodbav:"PK"`
		SK           string `dynamodbav:"SK"`
		ID           string `dynamodbav:"id"`
		Email        string `dynamodbav:"email"`
		FirstName    string `dynamodbav:"first_name"`
		LastName     string `dynamodbav:"last_name"`
		MiddleName   string `dynamodbav:"middle_name"`
		MobileNumber string `dynamodbav:"mobile_number"`
		Password     string `dynamodbav:"password"`
		Role         string `dynamodbav:"role"`
		Flat         string `dynamodbav:"flat"`
	}

	var userDetails User

	err = attributevalue.UnmarshalMap(result.Items[0], &userDetails)
	if err != nil {
		return err
	}

	deleteUserStatement := "DELETE FROM " + r.TableName + " WHERE PK = ? AND SK = ?"

	_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(deleteUserStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "ROLE#" + string(role)},
			&types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	_, err = r.DynamoDbClient.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(deleteUserStatement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: "USERS"},
			&types.AttributeValueMemberS{Value: (userDetails.Email + "#" + id)},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}
