package functions

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	logClient       *logging.Client
	ctx             context.Context
	firestoreClient *firestore.Client
	logger          *logging.Logger
	storageClient   *storage.Client
	pubSubClient    *pubsub.Client
	projectID       = "capstone-412502"
	err             error
)

// Initializing firestore connection and cloud logging.
func init() {
	ctx = context.Background()

	logClient, err = logging.NewClient(ctx, projectID, option.WithCredentialsFile("capstone-key.json"))

	if err != nil {
		log.Fatalf("Failed to create logging client: %v", err)
	}

	// Creates a logger for the specified log name.
	logger = logClient.Logger("my-log")

	firestoreClient, err = firestore.NewClientWithDatabase(ctx, projectID, "capstone", option.WithCredentialsFile("capstone-key.json"))
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)

		logger.Log(logging.Entry{
			Payload:  "Failed to create Firestore client",
			Severity: logging.Error,
		})

	}

	// Initialize Google Cloud Storage client
	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("capstone-key.json"))
	if err != nil {
		log.Fatalf("Failed to create storage client: %v", err)

		logger.Log(logging.Entry{
			Payload:  "Failed to create cloud storage client client",
			Severity: logging.Error,
		})
		return
	}

	pubSubClient, err = pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("capstone-key.json"))
	if err != nil {
		log.Fatalf("Failed to create pub/sub client: %v", err)

		logger.Log(logging.Entry{
			Payload:  "Failed to create cloud pub/sub client client",
			Severity: logging.Error,
		})
		return
	}

}
