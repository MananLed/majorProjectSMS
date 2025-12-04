package repository

import (
	"context"
	"errors"
	"log"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserRepository struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type UserRepositoryInterface interface {
	AddUser(user model.User) error
	GetUserByIDAndPassword(id string, password string) (*model.User, error)
	UpdateUser(user model.User, previousEmail string) error
	ChangePassword(id string, role model.UserRole, email string, newHashedPassword string) error
	DeleteUserByID(id string, role model.UserRole, email string) error
	GetUserByID(id string) (*model.User, error)
}

func NewUserRepository(ddbClient *dynamodb.Client, tableName string) *UserRepository {
	return &UserRepository{DynamoDbClient: ddbClient, TableName: tableName}
}

func (r *UserRepository) AddUser(newUser model.User) error {

	statement := "INSERT INTO " + r.TableName + " VALUE {'PK': ?, 'SK': ?, 'email': ?, 'first_name': ?, 'flat': ?, 'id': ?, 'middle_name': ?, 'last_name': ?, 'mobile_number': ?, 'password': ?, 'role': ?}"

	input := &dynamodb.ExecuteTransactionInput{
		TransactStatements: []types.ParameterizedStatement{
			{
				Statement: aws.String(statement),
				Parameters: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "USERS"},
					&types.AttributeValueMemberS{Value: (newUser.Email + "#" + newUser.ID)},
					&types.AttributeValueMemberS{Value: newUser.Email},
					&types.AttributeValueMemberS{Value: newUser.FirstName},
					&types.AttributeValueMemberS{Value: newUser.Flat},
					&types.AttributeValueMemberS{Value: newUser.ID},
					&types.AttributeValueMemberS{Value: newUser.MiddleName},
					&types.AttributeValueMemberS{Value: newUser.LastName},
					&types.AttributeValueMemberS{Value: newUser.MobileNumber},
					&types.AttributeValueMemberS{Value: newUser.Password},
					&types.AttributeValueMemberS{Value: string(newUser.Role)},
				},
			},
			{
				Statement: aws.String(statement),
				Parameters: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: ("ROLE#" + string(newUser.Role))},
					&types.AttributeValueMemberS{Value: newUser.ID},
					&types.AttributeValueMemberS{Value: newUser.Email},
					&types.AttributeValueMemberS{Value: newUser.FirstName},
					&types.AttributeValueMemberS{Value: newUser.Flat},
					&types.AttributeValueMemberS{Value: newUser.ID},
					&types.AttributeValueMemberS{Value: newUser.MiddleName},
					&types.AttributeValueMemberS{Value: newUser.LastName},
					&types.AttributeValueMemberS{Value: newUser.MobileNumber},
					&types.AttributeValueMemberS{Value: newUser.Password},
					&types.AttributeValueMemberS{Value: string(newUser.Role)},
				},
			},
		},
	}

	_ , err := r.DynamoDbClient.ExecuteTransaction(context.TODO(), input)

	return err
}

func (r *UserRepository) GetUserByIDAndPassword(email string, password string) (*model.User, error) {
	var user model.User
	loginCredentials := dto.LoginRequestDTO{Email: email, Password: password}

	primaryKey, err := attributevalue.Marshal("USERS")
	if err != nil {
		return nil, errors.New("internal server error")
	}

	sortKey, err := attributevalue.Marshal(loginCredentials.Email)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	keyMap := map[string]types.AttributeValue{"PK": primaryKey, "SK": sortKey}

	log.Print(keyMap)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "USERS"},
			":skPrefix": &types.AttributeValueMemberS{Value: email},
		},
	}

	ctx := context.TODO()
	response, err := r.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var userDetails dto.User

	if err != nil {
		return nil, err
	} else {
		err = attributevalue.UnmarshalMap(response.Items[0], &userDetails)
		if err != nil {
			return nil, err
		}
	}

	user.ID = userDetails.ID
	user.Email = userDetails.Email
	user.FirstName = userDetails.FirstName
	user.LastName = userDetails.LastName
	user.MiddleName = userDetails.MiddleName
	user.MobileNumber = userDetails.MobileNumber
	user.Password = userDetails.Password
	user.Role = model.UserRole(userDetails.Role)
	user.Flat = userDetails.Flat

	return &user, err
}

func (r *UserRepository) UpdateUser(updatedUser model.User, previousEmail string) error {

	if previousEmail == "" {
		updateRequestStatement := "UPDATE " + r.TableName + " SET first_name = ?, middle_name = ?, last_name = ?, mobile_number = ? WHERE PK = ? AND SK = ?"

		input := &dynamodb.ExecuteTransactionInput{
			TransactStatements: []types.ParameterizedStatement{
				{
					Statement: aws.String(updateRequestStatement),
					Parameters: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: updatedUser.FirstName},
						&types.AttributeValueMemberS{Value: updatedUser.MiddleName},
						&types.AttributeValueMemberS{Value: updatedUser.LastName},
						&types.AttributeValueMemberS{Value: updatedUser.MobileNumber},
						&types.AttributeValueMemberS{Value: "USERS"},
						&types.AttributeValueMemberS{Value: (updatedUser.Email + "#" + updatedUser.ID)},
					},
				},
				{
					Statement: aws.String(updateRequestStatement),
					Parameters: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: updatedUser.FirstName},
						&types.AttributeValueMemberS{Value: updatedUser.MiddleName},
						&types.AttributeValueMemberS{Value: updatedUser.LastName},
						&types.AttributeValueMemberS{Value: updatedUser.MobileNumber},
						&types.AttributeValueMemberS{Value: ("ROLE#" + string(updatedUser.Role))},
						&types.AttributeValueMemberS{Value: updatedUser.ID},
					},
				},
			},
		}

		_ , err := r.DynamoDbClient.ExecuteTransaction(context.TODO(), input)

		if err != nil{
			return err
		}
	} else {
		deleteRequestStatement := "DELETE FROM " + r.TableName + " WHERE PK = ? AND SK = ?"
		updateRequestStatement := "UPDATE " + r.TableName + " SET first_name = ?, middle_name = ?, last_name = ?, mobile_number = ?, email = ? WHERE PK = ? AND SK = ?"
		insertRequestStatement := "INSERT INTO " + r.TableName + " VALUE {'PK': ?, 'SK': ?, 'email': ?, 'first_name': ?, 'flat': ?, 'id': ?, 'middle_name': ?, 'last_name': ?, 'mobile_number': ?, 'password': ?, 'role': ?}"

		input := &dynamodb.ExecuteTransactionInput{
			TransactStatements: []types.ParameterizedStatement{
				{
					Statement: aws.String(deleteRequestStatement),
					Parameters: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: "USERS"},
						&types.AttributeValueMemberS{Value: previousEmail + "#" + updatedUser.ID},
					},
				},
				{
					Statement: aws.String(insertRequestStatement),
					Parameters: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: "USERS"},
						&types.AttributeValueMemberS{Value: (updatedUser.Email + "#" + updatedUser.ID)},
						&types.AttributeValueMemberS{Value: updatedUser.Email},
						&types.AttributeValueMemberS{Value: updatedUser.FirstName},
						&types.AttributeValueMemberS{Value: updatedUser.Flat},
						&types.AttributeValueMemberS{Value: updatedUser.ID},
						&types.AttributeValueMemberS{Value: updatedUser.MiddleName},
						&types.AttributeValueMemberS{Value: updatedUser.LastName},
						&types.AttributeValueMemberS{Value: updatedUser.MobileNumber},
						&types.AttributeValueMemberS{Value: updatedUser.Password},
						&types.AttributeValueMemberS{Value: string(updatedUser.Role)},
					},
				},
				{
					Statement: aws.String(updateRequestStatement),
					Parameters: []types.AttributeValue{
						&types.AttributeValueMemberS{Value: updatedUser.FirstName},
						&types.AttributeValueMemberS{Value: updatedUser.MiddleName},
						&types.AttributeValueMemberS{Value: updatedUser.LastName},
						&types.AttributeValueMemberS{Value: updatedUser.MobileNumber},
						&types.AttributeValueMemberS{Value: updatedUser.Email},
						&types.AttributeValueMemberS{Value: ("ROLE#" + string(updatedUser.Role))},
						&types.AttributeValueMemberS{Value: updatedUser.ID},
					},
				},
			},
		}

		_ , err := r.DynamoDbClient.ExecuteTransaction(context.TODO(), input)

		if err != nil{
			return err
		}
	}
	return nil
}

func (r *UserRepository) ChangePassword(id string, role model.UserRole, email string, newHashedPassword string) error {

	updateRequestStatement := "UPDATE " + r.TableName + " SET password = ? WHERE PK = ? AND SK = ?"

	input := &dynamodb.ExecuteTransactionInput{
		TransactStatements: []types.ParameterizedStatement{
			{
				Statement: aws.String(updateRequestStatement),
				Parameters: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: (newHashedPassword)},
					&types.AttributeValueMemberS{Value: "ROLE#" + string(role)},
					&types.AttributeValueMemberS{Value: id},
				},
			},
			{
				Statement: aws.String(updateRequestStatement),
				Parameters: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: (newHashedPassword)},
					&types.AttributeValueMemberS{Value: "USERS"},
					&types.AttributeValueMemberS{Value: email + "#" + id},
				},
			},
		},
	}

	_ , err := r.DynamoDbClient.ExecuteTransaction(context.TODO(), input)

	return err
}

func (r *UserRepository) DeleteUserByID(id string, role model.UserRole, email string) error {

	deleteUserStatement := "DELETE FROM " + r.TableName + " WHERE PK = ? AND SK = ?"

	input := &dynamodb.ExecuteTransactionInput{
		TransactStatements: []types.ParameterizedStatement{
			{
				Statement: aws.String(deleteUserStatement),
				Parameters: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "ROLE#" + string(role)},
					&types.AttributeValueMemberS{Value: id},
				},
			},
			{
				Statement: aws.String(deleteUserStatement),
				Parameters: []types.AttributeValue{
					&types.AttributeValueMemberS{Value: "USERS"},
					&types.AttributeValueMemberS{Value: (email + "#" + id)},
				},
			},
		},
	}

	_ , err := r.DynamoDbClient.ExecuteTransaction(context.TODO(), input)

	return err
}

func (r *UserRepository) GetUserByID(id string) (*model.User, error) {

	var user model.User

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "USERS"},
			":skPrefix": &types.AttributeValueMemberS{Value: id},
		},
	}

	ctx := context.TODO()
	response, err := r.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, nil
	}

	var userDetails dto.User

	if err != nil {
		return nil, err
	} else {
		err = attributevalue.UnmarshalMap(response.Items[0], &userDetails)
		if err != nil {
			return nil, err
		}
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
	return &user, err
}
