package main

import (
	"context"
	"log"

	"cloud.google.com/go/bigquery"
	"example.com/hash/go-logger"
)

func QueryWithCount(client *bigquery.Client, ctx context.Context, Query string) (*bigquery.RowIterator, int64, error) {

	q := client.Query(Query)
	job, err := q.Run(ctx)
	if err != nil {
		logger.Error("Error while running query in BigQuery", err)
		return nil, 0, err
	}

	// Wait for job to complete
	status, err := job.Wait(ctx)
	if err != nil {
		logger.Error("Error while waiting for BigQuery job", err)
		return nil, 0, err
	}

	if err := status.Err(); err != nil {
		logger.Error("Error in BigQuery job status", err)
		return nil, 0, err
	}

	// Get RowIterator from job (this gives us access to schema)
	iterator, err := job.Read(ctx)
	if err != nil {
		logger.Error("Error while reading BigQuery job results", err)
		return nil, 0, err
	}

	// Get total rows from the iterator (TotalRows is populated after job completes)
	// TotalRows is uint64, convert to int64
	totalRows := int64(iterator.TotalRows)

	return iterator, totalRows, nil
}

func main() {
	ctx := context.Background()

	// Example query - replace with your actual query
	query := `SELECT DISTINCT phone AS email FROM (
SELECT cc.phone
FROM (
SELECT
pv.entityid,
COUNT(pv.identifier) AS counts
FROM datapipelineproduction.datos_pocos.page_view AS pv
WHERE pv.page_type IN ('product_detail','variant_detail')
AND TIMESTAMP_ADD(TIMESTAMP(pv.event_timestamp), INTERVAL 330 MINUTE) >= TIMESTAMP('2025-10-01')
AND TIMESTAMP_ADD(TIMESTAMP(pv.event_timestamp), INTERVAL 330 MINUTE) < TIMESTAMP('2025-12-31')
GROUP BY pv.entityid
) AS t0
JOIN datapipelineproduction.datos_deposito_banco.purplle_purplle2_contact_contact AS cc
ON CAST(cc.user_id AS STRING) = t0.entityid
LEFT JOIN datapipelineproduction.datos_pocos.app_open AS apo
ON CAST(cc.user_id AS STRING) = apo.entityid
AND TIMESTAMP_ADD(TIMESTAMP(apo.event_timestamp), INTERVAL 330 MINUTE) >= TIMESTAMP('2025-12-31')
INNER JOIN datapipelineproduction.datos_deposito_banco.purplle_purplle2_shop_order AS shop_order
ON shop_order.contact_id = cc.id
AND shop_order.status IN ('Verified','In Process','Shipped','Partly Shipped','Complete')
AND shop_order.order_type = 'b2c'
AND shop_order.tenant = 'PURPLLE_COM'
INNER JOIN datapipelineproduction.datos_deposito_banco.purplle_purplle2_shop_orderitem AS soi
ON shop_order.id = soi.order_id
AND CAST(soi.our_price AS FLOAT64) <> 0
WHERE apo.entityid IS NULL
AND t0.counts >= 3
AND CAST(cc.user_id AS STRING) NOT IN (
SELECT DISTINCT entityid
FROM datapipelineproduction.datos_pocos.campaign_sent
WHERE DATE(TIMESTAMP_ADD(TIMESTAMP(event_timestamp), INTERVAL 330 MINUTE)) >= CURRENT_DATE() - 30
AND LOWER(eventsource) = 'whatsapp'
)
)`

	// Initialize BigQuery client - replace with your project ID
	projectID := "datapipelineproduction"
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create BigQuery client: %v", err)
	}
	defer client.Close()

	iterator, count, err := QueryWithCount(client, ctx, query)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

	log.Printf("Total rows: %d", count)

	// Read first row
	var firstRow []bigquery.Value
	err = iterator.Next(&firstRow)
	if err != nil {
		if err.Error() == "bigquery: no rows in result set" || err.Error() == "EOF" {
			log.Println("No rows found")
		} else {
			log.Fatalf("Error reading first row: %v", err)
		}
	} else {
		log.Printf("First row: %v", firstRow)
	}
}
