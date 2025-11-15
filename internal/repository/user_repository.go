package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/crypto/bcrypt"
)

// type TableBasics struct {
// 	DynamoDbClient *dynamodb.Client
// 	TableName      string
// }

// func (basics *TableBasics) GetUserByIDAndPassword(email string, password string) (*model.User, error){
// 	loginCredentials := dto.LoginRequestDTO{Email: email, Password: password}

// 	emailDDB, err := attributevalue.Marshal(loginCredentials.Email)
// 	if err != nil {
// 		return nil, errors.New("Internal Server Error")
// 	}

// 	passwordDDB, err := attributevalue.Marshal(loginCredentials.Password)
// 	if err != nil {
// 		return nil, errors.New("Internal Server Error")
// 	}

// 	keyMap := map[string]types.AttributeValue{"email": emailDDB, "password":passwordDDB}

// 	ctx := context.TODO()

// 	response, err := basics.DynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{
// 		Key: keyMap, TableName: aws.String(basics.TableName),
// 	})
// }

type UserRepository struct {
	mu             sync.Mutex
	db             *sql.DB
	DynamoDbClient *dynamodb.Client
	TableName      string
}

type UserRepositoryInterface interface {
	AddUser(user model.User) error
	GetUserByIDAndPassword(id string, password string) (*model.User, error)
	UpdateUser(user model.User) error
	ChangePassword(id string, newHashedPassword string) error
	IsPasswordUnique(hashedPassword string) bool
	DeleteUserByID(id string) error
	GetUserByID(id string) (*model.User, error)
}

func NewUserRepository(ddbClient *dynamodb.Client, tableName string) *UserRepository {
	return &UserRepository{DynamoDbClient: ddbClient, TableName: tableName}
}

func (r *UserRepository) AddUser(newUser model.User) error {

	// query := `
	// 	INSERT INTO users (id, first_name, middle_name, last_name, mobile_number, email, password, role, flat_no)
	// 	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	// `

	// _, err := r.db.Exec(query, newUser.ID, newUser.FirstName, newUser.MiddleName, newUser.LastName, newUser.MobileNumber, newUser.Email, newUser.Password, newUser.Role, newUser.Flat)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return err
	// }

	// return nil
	statement := "INSERT INTO " + r.TableName + " VALUE {'PK': ?, 'SK': ?, 'email': ?, 'first_name': ?, 'flat': ?, 'id': ?, 'middle_name': ?, 'last_name': ?, 'mobile_number': ?, 'password': ?, 'role': ?}"

	_, err := r.DynamoDbClient.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
        Statement: &statement,
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
    })

	if err != nil {
		return err 
	}

	_, err = r.DynamoDbClient.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
        Statement: &statement,
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
    })

	return err
}

func (r *UserRepository) GetUserByIDAndPassword(email string, password string) (*model.User, error) {
	var user model.User
	loginCredentials := dto.LoginRequestDTO{Email: email, Password: password}

	primaryKey, err := attributevalue.Marshal("USERS")
	if err != nil {
		return nil, errors.New("Internal Server Error")
	}

	sortKey, err := attributevalue.Marshal(loginCredentials.Email)
	if err != nil {
		return nil, errors.New("Internal Server Error")
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

	type User struct {
		PK           string  `dynamodbav:"PK"`
		SK           string  `dynamodbav:"SK"`
		ID           string  `dynamodbav:"id"` 
		Email        string  `dynamodbav:"email"`
		FirstName    string  `dynamodbav:"first_name"`
		LastName     string  `dynamodbav:"last_name"`
		MiddleName   string `dynamodbav:"middle_name"`
		MobileNumber string `dynamodbav:"mobile_number"`
		Password     string  `dynamodbav:"password"`
		Role         string  `dynamodbav:"role"`
		Flat         string `dynamodbav:"flat"` 
	}

	var userDetails User

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

func (r *UserRepository) UpdateUser(updatedUser model.User) error {

	query := `
		UPDATE users
		SET first_name = $1, middle_name = $2, last_name = $3, mobile_number = $4, email = $5, password = $6, role = $7
		WHERE id = $8
	`

	result, err := r.db.Exec(query, updatedUser.FirstName, updatedUser.MiddleName, updatedUser.LastName, updatedUser.MobileNumber, updatedUser.Email, updatedUser.Password, updatedUser.Role, updatedUser.ID)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("UpdateUser error: %v", err))
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) ChangePassword(id string, newHashedPassword string) error {

	query := `
		UPDATE users
		SET password = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, newHashedPassword, id)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) IsPasswordUnique(password string) bool {

	rows, err := r.db.Query(`
		SELECT password FROM users
	`)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return false
	}

	defer rows.Close()

	for rows.Next() {
		var hashed string
		if err := rows.Scan(&hashed); err != nil {
			return false
		}
		if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil {
			return false
		}
	}
	return true
}

func (r *UserRepository) DeleteUserByID(id string) error {

	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) GetUserByID(id string) (*model.User, error) {
	// var user model.User

	// query := `
	// 	SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role, flat_no
	// 	FROM users
	// 	WHERE id = $1
	// `

	// err := r.db.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber, &user.Email, &user.Password, &user.Role, &user.Flat)

	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return nil, errors.New("user not found")
	// 	}
	// 	logger.LogToFile(fmt.Sprintf("error fetching user: %v", err))
	// 	return nil, err
	// }

	// return &user, nil
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

	type User struct {
		PK           string  `dynamodbav:"PK"`
		SK           string  `dynamodbav:"SK"`
		ID           string  `dynamodbav:"id"` 
		Email        string  `dynamodbav:"email"`
		FirstName    string  `dynamodbav:"first_name"`
		LastName     string  `dynamodbav:"last_name"`
		MiddleName   string `dynamodbav:"middle_name"`
		MobileNumber string `dynamodbav:"mobile_number"`
		Password     string  `dynamodbav:"password"`
		Role         string  `dynamodbav:"role"`
		Flat         string `dynamodbav:"flat"` 
	}

	var userDetails User

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
