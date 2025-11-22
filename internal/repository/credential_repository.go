package repository

import (
	"context"

	"errors"
	"fmt"

	"github.com/MananLed/majorProjectSMS/internal/dto"
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

	var userDetails dto.User

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
