package repository

import (
	"context"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SocietyRepository struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type SocietyRepositoryInterface interface {
	GetAllResidents() ([]model.User, error)
	GetAllOfficers() ([]model.User, error)
}

func NewSocietyRepository(ddbClient *dynamodb.Client, tableName string) *SocietyRepository {
	return &SocietyRepository{DynamoDbClient: ddbClient, TableName: tableName}
}

func (s *SocietyRepository) GetAllResidents() ([]model.User, error) {

	var user model.User
	var users []model.User

	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "ROLE#resident"},
		},
	}

	ctx := context.TODO()
	response, err := s.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var userDetails dto.User

	if err != nil {
		return nil, err
	} else {
		for _, u := range response.Items {
			err = attributevalue.UnmarshalMap(u, &userDetails)
			if err != nil {
				return nil, err
			}
			user.ID = userDetails.ID
			user.Email = userDetails.Email
			user.FirstName = userDetails.FirstName
			user.MiddleName = userDetails.MiddleName
			user.LastName = userDetails.LastName
			user.MobileNumber = userDetails.MobileNumber
			user.Password = userDetails.Password
			user.Flat = userDetails.Flat
			user.Role = model.UserRole(userDetails.Role)
			users = append(users, user)
		}
	}

	return users, nil
}

func (s *SocietyRepository) GetAllOfficers() ([]model.User, error) {

	var user model.User
	var users []model.User

	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "ROLE#officer"},
		},
	}

	ctx := context.TODO()
	response, err := s.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var userDetails dto.User

	if err != nil {
		return nil, err
	} else {
		for _, u := range response.Items {
			err = attributevalue.UnmarshalMap(u, &userDetails)
			if err != nil {
				return nil, err
			}
			user.ID = userDetails.ID
			user.Email = userDetails.Email
			user.FirstName = userDetails.FirstName
			user.MiddleName = userDetails.MiddleName
			user.LastName = userDetails.LastName
			user.MobileNumber = userDetails.MobileNumber
			user.Password = userDetails.Password
			user.Flat = userDetails.Flat
			user.Role = model.UserRole(userDetails.Role)
			users = append(users, user)
		}
	}

	return users, nil
}
