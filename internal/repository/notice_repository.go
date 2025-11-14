package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type NoticeRepositoryInterface interface {
	SaveNotice(notice model.Notice) error
	GetAllNotices() ([]model.Notice, error)
	GetNoticesByMonthYear(month time.Month, year int) ([]model.Notice, error)
	GetNoticesByYear(year int) ([]model.Notice, error)
}

type NoticeRepository struct {
	DB             *sql.DB
	mu             sync.Mutex
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func NewNoticeRepository(ddbClient *dynamodb.Client, tableName string) *NoticeRepository {
	return &NoticeRepository{DynamoDbClient: ddbClient, TableName: tableName}
}

func (r *NoticeRepository) SaveNotice(notice model.Notice) error {

	notice.ID = utils.GenerateUUID()

	query := `
		INSERT INTO notices (id, date_issued, content, month, year)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.DB.Exec(query, notice.ID, notice.DateIssued, notice.Content, notice.Month, notice.Year)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}
	return nil
}

func (r *NoticeRepository) GetAllNotices() ([]model.Notice, error) {
	// query := `
	// SELECT id, date_issued, content, month, year
	// FROM notices
	// ORDER BY date_issued DESC
	// `

	// rows, err := r.DB.Query(query)
	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return nil, err
	// }
	// defer rows.Close()

	// var notices []model.Notice
	// for rows.Next() {
	// 	var n model.Notice
	// 	if err := rows.Scan(&n.ID, &n.DateIssued, &n.Content, &n.Month, &n.Year); err != nil {
	// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 		return nil, err
	// 	}
	// 	notices = append(notices, n)
	// }

	// return notices, nil

	var notice model.Notice
	var notices []model.Notice

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue": &types.AttributeValueMemberS{Value: "NOTICES"},
		},
	}

	ctx := context.TODO()
	response, err := r.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	type Notice struct {
		PK          string `dynamodbav:"PK"`
		SK          string `dynamodbav:"SK"`
		ID          string `dynamodbav:"id"`
		Content     string `dynamodbav:"content"`
		Date_Issued string `dynamodbav:"date_issued"`
		Month       string `dynamodbav:"month"`
		Year        int    `dynamodbav:"year"`
	}

	var noticeDetails Notice

	if err != nil {
		return nil, err
	} else {
		for _, n := range response.Items {
			err = attributevalue.UnmarshalMap(n, &noticeDetails)
			if err != nil {
				return nil, err
			}

			layout := "2006-01-02 15:04:05.000000"

			notice.ID, _ = uuid.Parse(noticeDetails.ID)
			notice.DateIssued, _ = time.Parse(layout, noticeDetails.Date_Issued)
			notice.Content = noticeDetails.Content
			monthInt, _ := strconv.Atoi(noticeDetails.Month)
			notice.Month = time.Month(monthInt)
			notice.Year = noticeDetails.Year

			notices = append(notices, notice)
		}
	}

	return notices, nil
}

func (r *NoticeRepository) GetNoticesByMonthYear(month time.Month, year int) ([]model.Notice, error) {
	// query := `
	// SELECT id, date_issued, content, month, year
	// FROM notices
	// WHERE month = $1 AND year = $2
	// ORDER BY date_issued DESC
	// `

	// rows, err := r.DB.Query(query, month, year)

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return nil, err
	// }
	// defer rows.Close()

	// var notices []model.Notice
	// for rows.Next() {
	// 	var n model.Notice
	// 	if err := rows.Scan(&n.ID, &n.DateIssued, &n.Content, &n.Month, &n.Year); err != nil {
	// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 		return nil, err
	// 	}
	// 	notices = append(notices, n)
	// }

	// return notices, nil
	var notice model.Notice
	var notices []model.Notice

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "NOTICES"},
			":skPrefix": &types.AttributeValueMemberS{Value: strconv.Itoa(year) + "#" + fmt.Sprintf("%d", int(month))},
		},
	}

	ctx := context.TODO()
	response, err := r.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	type Notice struct {
		PK          string     `dynamodbav:"PK"`
		SK          string     `dynamodbav:"SK"`
		ID          string  `dynamodbav:"id"`
		Content     string     `dynamodbav:"content"`
		Date_Issued string  `dynamodbav:"date_issued"`
		Month       string `dynamodbav:"month"`
		Year        int        `dynamodbav:"year"`
	}

	var noticeDetails Notice

	if err != nil {
		return nil, err
	} else {
		for _, n := range response.Items {
			err = attributevalue.UnmarshalMap(n, &noticeDetails)
			if err != nil {
				return nil, err
			}

			layout := "2006-01-02 15:04:05.000000"

			notice.ID, _ = uuid.Parse(noticeDetails.ID)
			notice.DateIssued, _ = time.Parse(layout, noticeDetails.Date_Issued)
			notice.Content = noticeDetails.Content
			monthInt, _ := strconv.Atoi(noticeDetails.Month)
			notice.Month = time.Month(monthInt)
			notice.Year = noticeDetails.Year

			notices = append(notices, notice)
		}
	}

	return notices, nil
}

func (r *NoticeRepository) GetNoticesByYear(year int) ([]model.Notice, error) {
	var notice model.Notice
	var notices []model.Notice

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "NOTICES"},
			":skPrefix": &types.AttributeValueMemberS{Value: strconv.Itoa(year)},
		},
	}

	ctx := context.TODO()
	response, err := r.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	type Notice struct {
		PK          string     `dynamodbav:"PK"`
		SK          string     `dynamodbav:"SK"`
		ID          string  `dynamodbav:"id"`
		Content     string     `dynamodbav:"content"`
		Date_Issued string  `dynamodbav:"date_issued"`
		Month       string `dynamodbav:"month"`
		Year        int        `dynamodbav:"year"`
	}

	var noticeDetails Notice

	if err != nil {
		return nil, err
	} else {
		for _, n := range response.Items {
			err = attributevalue.UnmarshalMap(n, &noticeDetails)
			if err != nil {
				return nil, err
			}

			layout := "2006-01-02 15:04:05.000000"

			notice.ID, _ = uuid.Parse(noticeDetails.ID)
			notice.DateIssued, _ = time.Parse(layout, noticeDetails.Date_Issued)
			notice.Content = noticeDetails.Content
			monthInt, _ := strconv.Atoi(noticeDetails.Month)
			notice.Month = time.Month(monthInt)
			notice.Year = noticeDetails.Year

			notices = append(notices, notice)
		}
	}

	return notices, nil
}
