package aws

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/ihopethisisfine/helloworld/internal/pkg/storage"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var svc dynamodbiface.DynamoDBAPI
var userStorage UserStorage

var createdTableOutput *dynamodb.CreateTableOutput

// the TestMain functions runs before any test
func TestMain(m *testing.M) {
	// get a context for our request
	ctx := context.Background()
	// create a container request for DynamoDB
	req := testcontainers.ContainerRequest{
		Image: "amazon/dynamodb-local:latest",
		// in-memory version is good enough
		Cmd: []string{"-jar", "DynamoDBLocal.jar", "-inMemory"},
		// by default, DynamoDB runs on port 8000
		ExposedPorts: []string{"8000/tcp"},
		// block until the port is available
		WaitingFor: wait.NewHostPortStrategy("8000"),
	}

	// start the container
	d, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		panic(err)
	}

	// stop the container
	defer d.Terminate(ctx)

	// get the IP and port of DynamoDB instance to connect to the right endpoints
	ip, err := d.Host(ctx)

	if err != nil {
		panic(err)
	}

	port, err := d.MappedPort(ctx, "8000")

	if err != nil {
		panic(err)
	}

	// create a new session with custom endpoint
	// need to specify a region, otherwise there's a fatal error
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String(fmt.Sprintf("http://%s:%s", ip, port)),
		Region:   aws.String("eu-central-1"),
	}))

	svc = dynamodb.New(sess)

	userStorage = NewUserStorage(sess, time.Second*5)

	table := "users"

	if createdTableOutput == nil {
		createdTableOutput, err = svc.CreateTable(&dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("username"),
					AttributeType: aws.String("S"),
				},
			},
			BillingMode:            nil,
			GlobalSecondaryIndexes: nil,
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("username"),
					KeyType:       aws.String(dynamodb.KeyTypeHash),
				},
			},
			LocalSecondaryIndexes: nil,
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
			SSESpecification:    nil,
			StreamSpecification: nil,
			TableName:           aws.String(table),
			Tags:                nil,
		})

		if err != nil {
			panic(err)
		}
	}

	// run the tests
	os.Exit(m.Run())
}

func TestInsert(t *testing.T) {
	insertedUser := storage.User{Username: "Mike", DateOfBirth: "2020-12-02"}
	err := userStorage.Put(context.Background(), insertedUser)

	if err != nil {
		t.Fatalf("could not insert: %s", err.Error())
	}
	u, err := userStorage.Find(context.Background(), insertedUser.Username)

	if err != nil {
		t.Fatalf("could not find: %s", err.Error())
	}

	if insertedUser.DateOfBirth != u.DateOfBirth {
		t.Fatalf("output \"%s\" is wrong! Should be \"%s\" instead.", u, insertedUser)
	}
}

func TestUpdate(t *testing.T) {

	user := storage.User{Username: "Eva", DateOfBirth: "2020-03-02"}
	if err := userStorage.Put(context.Background(), user); err != nil {
		t.Fatalf("could not insert: %s", err.Error())
	}

	user = storage.User{Username: "Eva", DateOfBirth: "2020-07-03"}
	if err := userStorage.Put(context.Background(), user); err != nil {
		t.Fatalf("could not update: %s", err.Error())
	}

	expectedUser, err := userStorage.Find(context.Background(), user.Username)
	if err != nil {
		t.Fatalf("could not find: %s", err.Error())
	}

	if user.DateOfBirth != expectedUser.DateOfBirth {
		t.Fatalf("output \"%s\" is wrong! Should be \"%s\" instead.", expectedUser, user)
	}
}
