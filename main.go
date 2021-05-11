package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/jenkins-x/jx-logging/v3/pkg/log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/manifoldco/promptui"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var (
	options struct {
		projectID string
		filter    string
		dryRun    bool
	}
)

func init() {
	flag.StringVar(&options.projectID, "project", "", "GCP project with which to communicate")
	flag.StringVar(&options.filter, "filter", "pr-", "match secret names that contain this string, defaults to 'pr-'")
	flag.BoolVar(&options.dryRun, "dry-run", false, "dry run will print which secrets would be deleted but does not perform the actual delete action")

}

func main() {
	flag.Parse()

	if options.projectID == "" {
		log.Logger().Fatal("missing --project flag")
	}

	if options.filter == "" {
		log.Logger().Warnf("You are about to delete ALL secrets, from GCP project '%s'", options.projectID)
	} else {
		log.Logger().Warnf("You are about to delete secrets that contain '%s' in the name, from GCP project '%s'", options.filter, options.projectID)
	}

	if !yesNo() {
		return
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

		if options.filter != "" && !strings.Contains(resp.Name, options.filter) {
			// lets only GC secrets for version stream pull requests
			continue
		}
		count++

		if options.dryRun {
			log.Logger().Infof("dry run: found secret to delete: %s", resp.Name)
			continue

		} else {
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

	}
	log.Logger().Infof("deleted %v secrets", count)
}

func yesNo() bool {
	prompt := promptui.Select{
		Label: "Are you sure [Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Logger().Fatal("Prompt failed %v", err)
	}
	return result == "Yes"
}
