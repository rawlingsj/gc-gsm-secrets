package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var (
	options struct {
		projectID string
	}
)

func init() {
	flag.StringVar(&options.projectID, "project", "", "GCP project with which to communicate")

}

func main() {
	flag.Parse()

	if options.projectID == "" {
		log.Logger().Fatal("missing --project flag")
	}

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Logger().Fatalf("failed to create secretmanager client: %v", err)
	}

	parent := fmt.Sprintf("projects/%s", options.projectID)

	// Build the request.
	req := &secretmanagerpb.ListSecretsRequest{
		Parent: parent,
	}

	// Call the API.
	it := client.ListSecrets(ctx, req)
	count := 0
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Logger().Fatalf("failed to list secret versions: %v", err)
		}

		if !strings.Contains(resp.Name, "pr-") {
			// lets only GC secrets for version stream pull requests
			continue
		}
		count++
		// Build the request.
		req := &secretmanagerpb.DeleteSecretRequest{
			Name: resp.Name,
		}

		// Call the API.
		if err := client.DeleteSecret(ctx, req); err != nil {
			log.Logger().Fatalf("failed to delete secret: %v", err)
		}

		log.Logger().Infof("deleted secret %s\n", resp.Name)
	}
	log.Logger().Infof("deleted %v secrets", count)
}
