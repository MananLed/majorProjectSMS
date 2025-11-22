package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/dto"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
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
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func NewNoticeRepository(ddbClient *dynamodb.Client, tableName string) *NoticeRepository {
	return &NoticeRepository{DynamoDbClient: ddbClient, TableName: tableName}
}

func (r *NoticeRepository) SaveNotice(notice model.Notice) error {
	
	notice.ID = utils.GenerateUUID()
	const customLayout = "2006-01-02 15:04:05.999999"
	formattedDate := notice.DateIssued.Format(customLayout)

	statement := "INSERT INTO " + r.TableName + " VALUE {'PK': ?, 'SK': ?, 'date_issued': ?, 'id': ?, 'content': ?, 'month': ?, 'year': ?}"
	
	_, err := r.DynamoDbClient.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
        Statement: &statement,
        Parameters: []types.AttributeValue{
            &types.AttributeValueMemberS{Value: "NOTICES"},
            &types.AttributeValueMemberS{Value: (strconv.Itoa(notice.Year) + "#" + fmt.Sprintf("%d", int(notice.Month)) + "#" + notice.ID.String())},
            &types.AttributeValueMemberS{Value: formattedDate},
            &types.AttributeValueMemberS{Value: notice.ID.String()},
            &types.AttributeValueMemberS{Value: notice.Content},
            &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", int(notice.Month))},
            &types.AttributeValueMemberN{Value: strconv.Itoa(notice.Year)},
        },
    })

	if err != nil {
		return err 
	}

	return nil
}

func (r *NoticeRepository) GetAllNotices() ([]model.Notice, error) {

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

	var noticeDetails dto.Notice

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

	var noticeDetails dto.Notice

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

	var noticeDetails dto.Notice

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
