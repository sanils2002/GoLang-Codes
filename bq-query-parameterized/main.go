package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// getFinalTargetUsers queries BigQuery to find identifiers that have already received campaigns
// and removes them from the identifiers map
func getFinalTargetUsers(ctx context.Context, client *bigquery.Client, identifiers map[string]string, label string, campaignName string, hours int64, considerAllCampaigns bool, campaignSource string) error {

	if hours == 0 {
		hours = 24
	}

	startTimestamp := time.Now().Unix() - (hours * 3600)
	query := "SELECT identifier FROM `datos_pocos.campaign_sent` WHERE event_timestamp BETWEEN TIMESTAMP_SECONDS(@startTs) AND CURRENT_TIMESTAMP() "
	params := []bigquery.QueryParameter{
		{Name: "startTs", Value: startTimestamp},
	}

	if campaignSource != "" {
		query += " AND eventsource = @campaignSource "
		params = append(params, bigquery.QueryParameter{Name: "campaignSource", Value: strings.ToLower(campaignSource)})
	}
	if !considerAllCampaigns && label != "" && campaignName != "" {
		query += " AND eventcampaign = @eventCampaign "
		params = append(params, bigquery.QueryParameter{Name: "eventCampaign", Value: label + "_" + campaignName})
	}
	query += " GROUP BY identifier "

	// Create and configure the query
	q := client.Query(query)
	q.Parameters = params

	// Run the query
	job, err := q.Run(ctx)
	if err != nil {
		log.Printf("Error running query: %v", err)
		return err
	}

	// Wait for query to complete
	status, err := job.Wait(ctx)
	if err != nil {
		log.Printf("Error waiting for job: %v", err)
		return err
	}
	if err := status.Err(); err != nil {
		log.Printf("Job completed with error: %v", err)
		return err
	}

	// Read results
	it, err := job.Read(ctx)
	if err != nil {
		log.Printf("Error reading results: %v", err)
		return err
	}

	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return err
		}
		if len(values) == 0 {
			log.Println("No data found in row")
			continue
		}
		if values[0] == nil {
			log.Println("Row has null value, skipping")
			continue
		}
		// Remove identifier from the map (already sent)
		delete(identifiers, strings.ToLower(values[0].(string)))
	}

	return nil
}

func main() {
	ctx := context.Background()

	// Initialize BigQuery client - replace with your project ID
	projectID := "datapipelineproduction"
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create BigQuery client: %v", err)
	}
	defer client.Close()

	// Sample identifiers map for testing
	identifiers := map[string]string{
		"9326517348": "9326517348",
	}

	fmt.Println("Identifiers before filtering:", identifiers)

	// Call the function
	err = getFinalTargetUsers(
		ctx,
		client,
		identifiers,
		"PWA",               // label
		"checkwhatsappFix3", // campaignName
		24,                  // hours (look back 24 hours)
		false,               // considerAllCampaigns
		"whatsapp",          // campaignSource
	)

	if err != nil {
		log.Printf("Error: %v", err)
	}

	fmt.Println("Identifiers after filtering:", identifiers)
}
