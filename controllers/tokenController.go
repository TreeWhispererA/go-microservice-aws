package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"

	"encoding/csv"
	"encoding/json"
	"os"

	"blockparty.co/test/db"
	middlewares "blockparty.co/test/handlers"
	"blockparty.co/test/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var svc = db.Dbconnect()
var tableName = "MetaData"

func Test(c *gin.Context) {
	middlewares.SuccessMessageResponse("Congratulations... It's working.", c.Writer)
}

// Retrieves all metadata items from the DynamoDB table
var GetTokens = gin.HandlerFunc(func(c *gin.Context) {

	var tokens []*models.Metadata

	// Create a ScanInput object return all items
	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	// Scan API to retrieve all items from the table
	result, err := svc.Scan(params)
	if err != nil {
		color.Red("Scan API call failed: %s", err)
	}

	// Iterate over each item in the result and unmarshal it into a Metadata struct
	for _, i := range result.Items {
		item := models.Metadata{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			color.Red("Got error unmarshalling: %s", err)
		}
		tokens = append(tokens, &item)
	}

	middlewares.SuccessArrRespond(tokens, "Metadata", c.Writer)
})

// Retrieves a single metadata item from the DynamoDB table based on the cid parameter
var GetTokenByID = gin.HandlerFunc(func(c *gin.Context) {

	cid := c.Param("cid")

	// GetItem API to retrieve the item with the specified cid
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(cid),
			},
		},
	})
	if err != nil {
		color.Red("Got error calling GetItem: %s", err)
		middlewares.ErrorResponse("Database Connection Error...", c.Writer)
		return
	}

	// If no matching item is found, return an error response
	if result.Item == nil {
		msg := "Could not find '" + cid + "'"
		middlewares.ErrorResponse(msg, c.Writer)
		return
	}

	// Unmarshal the item into a Metadata struct
	item := models.Metadata{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		color.Red("Failed to unmarshal Record, %v", err)
		middlewares.ErrorResponse("Failed to unmarshal Record", c.Writer)
		return
	}

	middlewares.SuccessOneRespond(item, "Metadata", c.Writer)
})

// stores the scraped metadata in a DynamoDB table
func storeData(metadata []*models.Metadata) error {

	for _, item := range metadata {

		// Marshal the item into a DynamoDB attribute value map
		av_item, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			color.Red("Got error marshalling metadata item: %s", err)
			return err
		}

		// Create a PutItemInput object
		input := &dynamodb.PutItemInput{
			Item:      av_item,
			TableName: aws.String(tableName),
		}

		// Call the PutItem API to store the item in the table
		_, err = svc.PutItem(input)
		if err != nil {
			color.Red("Got error calling PutItem: %s", err)
			return err
		}

		color.Cyan("Successfully added '"+item.ID+" in ", tableName)
	}
	return nil
}

func scrapeMetadata(cid string) (*models.Metadata, error) {
	url := fmt.Sprintf("https://blockpartyplatform.mypinata.cloud/ipfs/%s", cid)
	response, err := http.Get(url)
	if err != nil {
		color.Red("Got error calling http.Get: %s", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			color.Red("Got error calling ioutil.ReadAll: %s", err)
			return nil, err
		}

		var metadata models.Metadata
		err = json.Unmarshal(body, &metadata)
		if err != nil {
			color.Red("Got error Unmarshal: %s", err)
			return nil, err
		}

		metadata.ID = cid

		return &metadata, nil
	}

	return nil, errors.New("failed to retrieve data")
}

func ScrapFunc(c *gin.Context) {

	file, err := os.Open("ipfs_cids.csv")
	if err != nil {
		color.Red("File Open Failed...")
		middlewares.ErrorResponse("File Open Failed...", c.Writer)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the lines from the CSV file
	lines, err := reader.ReadAll()
	if err != nil {
		color.Red("File Read Failed...")
		middlewares.ErrorResponse("File Read Failed...", c.Writer)
		return
	}

	// Convert the lines to an array of strings
	var result []string
	for _, line := range lines {
		result = append(result, line[0])
	}

	var tokens []*models.Metadata

	for _, cid := range result {
		// Scrape metadata for the current cid
		metadata, err := scrapeMetadata(cid)
		if err != nil {
			color.Red("Failed to retrieve metadata for CID %s: %v\n", cid, err)
			continue
		}
		tokens = append(tokens, metadata)
	}

	// Store the tokens slice in the database
	err = storeData(tokens)

	if err != nil {
		color.Red(err.Error())
		middlewares.ErrorResponse(err.Error(), c.Writer)
	}

	// Send a successful response containing the tokens
	middlewares.SuccessArrRespond(tokens, `Metadata`, c.Writer)
}
