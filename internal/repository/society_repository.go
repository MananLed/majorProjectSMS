package repository

import (
	"context"

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

	// query := `
	// 	SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role, flat_no
	// 	FROM users WHERE role = $1
	// `

	// rows, err := s.db.Query(query, string(model.RoleResident))

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return nil, fmt.Errorf("failed to fetch residents: %w", err)
	// }
	// defer rows.Close()

	// var residents []model.User
	// for rows.Next() {
	// 	var user model.User
	// 	if err := rows.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber, &user.Email, &user.Password, &user.Role, &user.Flat); err != nil {
	// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 		return nil, fmt.Errorf("failed to scan resident: %w", err)
	// 	}
	// 	residents = append(residents, user)
	// }

	// if len(residents) > 0 {
	// 	fmt.Print(color.YellowString("Total Residents: "), len(residents))
	// } else {
	// 	fmt.Print("There are no residents currently.")
	// }
	// fmt.Println()

	// return residents, nil

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

	// query := `
	// 	SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role 
	// 	FROM users WHERE role = $1
	// `

	// rows, err := s.db.Query(query, string(model.RoleOfficer))

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return nil, fmt.Errorf("failed to fetch officers: %v", err)
	// }
	// defer rows.Close()

	// var officers []model.User
	// for rows.Next() {
	// 	var user model.User
	// 	if err := rows.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber, &user.Email, &user.Password, &user.Role); err != nil {
	// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 		return nil, fmt.Errorf("failed to scan resident: %w", err)
	// 	}
	// 	officers = append(officers, user)
	// }

	// if len(officers) > 0 {
	// 	fmt.Print(color.YellowString("Total Officers: "), len(officers))
	// } else {
	// 	fmt.Print("There are no officers currently.")
	// }
	// fmt.Println()

	// return officers, nil

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
