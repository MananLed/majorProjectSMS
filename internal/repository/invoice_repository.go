package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type InvoiceRepositoryInterface interface {
	SaveInvoice(model.Invoice) error
	GetInvoiceByMonthAndYear(time.Month, int) (*model.Invoice, error)
	GetInvoicesByYear(int) ([]model.Invoice, error)
}

type InvoiceRepository struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func NewInvoiceRepository(ddbClient *dynamodb.Client, tableName string) *InvoiceRepository {
	return &InvoiceRepository{DynamoDbClient: ddbClient, TableName: tableName}
}

func (r *InvoiceRepository) SaveInvoice(invoice model.Invoice) error {
	
	// query := `
	// 	INSERT INTO invoices (id, month, year, amount)
	// 	VALUES ($1, $2, $3, $4)
	// `
	
	// _, err := r.DB.Exec(query, invoice.ID, invoice.Month, invoice.Year, invoice.Amount)

	// if err != nil {
		// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
		// 	return err
	// }
	
	// return nil
	
	invoice.ID = utils.GenerateUUID()

	statement := "INSERT INTO " + r.TableName + " VALUE {'PK': ?, 'SK': ?, 'id': ?, 'amount': ?, 'month': ?, 'year': ?}"
	
	_, err := r.DynamoDbClient.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
        Statement: &statement,
        Parameters: []types.AttributeValue{
            &types.AttributeValueMemberS{Value: "INVOICES"},
            &types.AttributeValueMemberS{Value: (strconv.Itoa(invoice.Year) + "#" + fmt.Sprintf("%d", int(invoice.Month)) + "#" + invoice.ID.String())},
            &types.AttributeValueMemberS{Value: invoice.ID.String()},
            &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", invoice.Amount)},
            &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", int(invoice.Month))},
            &types.AttributeValueMemberN{Value: strconv.Itoa(invoice.Year)},
        },
    })

	if err != nil {
		return err 
	}

	return nil
}

func (r *InvoiceRepository) GetInvoiceByMonthAndYear(month time.Month, year int) (*model.Invoice, error) {
	// var invoice model.Invoice

	// query := `
	// 	SELECT id, month, year, amount
	// 	FROM invoices
	// 	WHERE month = $1 AND year = $2
	// `

	// err := r.DB.QueryRow(query, int(month), year).Scan(&invoice.ID, &invoice.Month, &invoice.Year, &invoice.Amount)

	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return nil, errors.New("invoice not found")
	// 	}
	// 	logger.LogToFile(fmt.Sprintf("error fetching invoice: %v", err))
	// 	return nil, err
	// }

	// return &invoice, nil

	var invoice model.Invoice

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "INVOICES"},
			":skPrefix": &types.AttributeValueMemberS{Value: strconv.Itoa(year) + "#" + fmt.Sprintf("%d", int(month))},
		},
	}

	ctx := context.TODO()
	response, err := r.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	type Invoice struct {
		PK     string  `dynamodbav:"PK"`
		SK     string  `dynamodbav:"SK"`
		ID     string  `dynamodbav:"id"`
		Amount float64 `dynamodbav:"amount"`
		Month  string  `dynamodbav:"month"`
		Year   int     `dynamodbav:"year"`
	}

	var invoiceDetails Invoice

	if err != nil {
		return nil, err
	} else {
		err = attributevalue.UnmarshalMap(response.Items[0], &invoiceDetails)
		if err != nil {
			return nil, err
		}

		invoice.ID, _ = uuid.Parse(invoiceDetails.ID)
		invoice.Amount = invoiceDetails.Amount
		monthInt, _ := strconv.Atoi(invoiceDetails.Month)
		invoice.Month = time.Month(monthInt)
		invoice.Year = invoiceDetails.Year
	}

	return &invoice, nil
}

func (r *InvoiceRepository) GetInvoicesByYear(year int) ([]model.Invoice, error) {
	// var query string

	// if year == 0 {
	// 	query = `
	// 	SELECT id, month, year, amount
	// 	FROM invoices
	// `
	// } else {
	// 	query = `
	// 	SELECT id, month, year, amount
	// 	FROM invoices
	// 	WHERE year = $1
	// `
	// }

	// var rows *sql.Rows
	// var err error
	// if year == 0 {
	// 	rows, err = r.DB.Query(query)
	// } else {
	// 	rows, err = r.DB.Query(query, year)
	// }

	// if err != nil {
	// 	logger.LogToFile(fmt.Sprintf("error: %v", err))
	// 	return nil, err
	// }
	// defer rows.Close()

	// var invoices []model.Invoice
	// for rows.Next() {
	// 	var inv model.Invoice
	// 	if err := rows.Scan(&inv.ID, &inv.Month, &inv.Year, &inv.Amount); err != nil {
	// 		logger.LogToFile(fmt.Sprintf("error scanning invoice row: %v", err))
	// 		return nil, err
	// 	}
	// 	invoices = append(invoices, inv)
	// }

	// return invoices, nil
	var invoice model.Invoice
	var invoices []model.Invoice

	if year == 0 {
		input := &dynamodb.QueryInput{
			TableName:              aws.String(r.TableName),
			KeyConditionExpression: aws.String("PK = :pkValue"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pkValue": &types.AttributeValueMemberS{Value: "INVOICES"},
			},
		}

		ctx := context.TODO()
		response, err := r.DynamoDbClient.Query(ctx, input)
		if err != nil {
			return nil, err
		}

		type Invoice struct {
			PK     string  `dynamodbav:"PK"`
			SK     string  `dynamodbav:"SK"`
			ID     string  `dynamodbav:"id"`
			Amount float64 `dynamodbav:"amount"`
			Month  string  `dynamodbav:"month"`
			Year   int     `dynamodbav:"year"`
		}

		var invoiceDetails Invoice

		if err != nil {
			return nil, err
		} else {
			for _, i := range response.Items {
				err = attributevalue.UnmarshalMap(i, &invoiceDetails)
				if err != nil {
					return nil, err
				}

				invoice.ID, _ = uuid.Parse(invoiceDetails.ID)
				invoice.Amount = invoiceDetails.Amount
				monthInt, _ := strconv.Atoi(invoiceDetails.Month)
				invoice.Month = time.Month(monthInt)
				invoice.Year = invoiceDetails.Year
				invoices = append(invoices, invoice)
			}
		}
		return invoices, nil
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String("PK = :pkValue AND begins_with(SK, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pkValue":  &types.AttributeValueMemberS{Value: "INVOICES"},
			":skPrefix": &types.AttributeValueMemberS{Value: strconv.Itoa(year)},
		},
	}

	ctx := context.TODO()
	response, err := r.DynamoDbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	type Invoice struct {
		PK     string  `dynamodbav:"PK"`
		SK     string  `dynamodbav:"SK"`
		ID     string  `dynamodbav:"id"`
		Amount float64 `dynamodbav:"amount"`
		Month  string  `dynamodbav:"month"`
		Year   int     `dynamodbav:"year"`
	}

	var invoiceDetails Invoice

	if err != nil {
		return nil, err
	} else {
		for _, i := range response.Items {
			err = attributevalue.UnmarshalMap(i, &invoiceDetails)
			if err != nil {
				return nil, err
			}

			invoice.ID, _ = uuid.Parse(invoiceDetails.ID)
			invoice.Amount = invoiceDetails.Amount
			monthInt, _ := strconv.Atoi(invoiceDetails.Month)
			invoice.Month = time.Month(monthInt)
			invoice.Year = invoiceDetails.Year
			invoices = append(invoices, invoice)
		}
	}
	return invoices, nil
}
