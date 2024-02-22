package main

import (
	"bytes"
	"context"
	"elastic_test/model"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	cert, err := os.ReadFile("./http_ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	var (
		r   model.Response
		buf bytes.Buffer
	)

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"https://127.0.0.1:9200",
		},
		Username: "leyban",
		Password: "secret",
		CACert:   cert,
	})

	if err != nil {
		log.Fatal(err)
	}

	format := "02/01/2006 15:04"

	dateFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(format)
	dateTo := time.Date(2024, 1, 1, 23, 0, 0, 0, time.UTC).Format(format)

	fmt.Println("dateFrom:", dateFrom)
	fmt.Println("dateTo:", dateTo)

	hitSize := 3
	aggSize := 3

	depositStatusSuccess := 1

	queryParams := map[string]any{
		"size": hitSize,
		"query": map[string]any{
			"range": map[string]any{
				"created_at": map[string]any{
					"gte": dateFrom,
					"lte": dateTo,
				},
			},
		},

		"sort": []map[string]any{
			{
				"gross_amount": map[string]any{
					"order": "desc",
				},
			},
		},

		"aggs": map[string]any{
			"filter_by_status": map[string]any{
				"filter": map[string]any{
					"term": map[string]any{
						"status": depositStatusSuccess,
					},
				},

				"aggs": map[string]any{
					"group_by_member_id": map[string]any{
						"terms": map[string]any{
							"field": "member_id",
							"size":  aggSize,
							"order": map[string]any{
								"total_gross_amount": "desc",
							},
						},

						"aggs": map[string]any{
							"total_gross_amount": map[string]any{
								"sum": map[string]any{
									"field": "gross_amount",
								},
							},
						},
					},
				},
			},
		},
	}

	err = json.NewEncoder(&buf).Encode(queryParams)
	if err != nil {
		log.Fatal(err)
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("deposit"),
		es.Search.WithBody(&buf),
	)
	if err != nil {
		log.Println("An Error Occured")
		log.Fatal(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	log.Println(r)

	log.Println("Biggest Deposit", r.Hits.Hits[0].Source.LoginName)
	log.Println("Biggest Depositor", r.Aggregations.FilterStatus.GroupByMemberID.Buckets[0].Key)
}
