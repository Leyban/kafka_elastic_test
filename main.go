package main

import (
	"context"
	"elastic_test/model"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	cert, err := os.ReadFile("./http_ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	var r model.Response

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

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("deposit"),
		es.Search.WithBody(strings.NewReader(`
        {
          "query": {
            "range": {
              " created_at": {
                "gte": "19/01/2024 00:00",
                "lte": "20/01/2024 23:00"
              }
            }
          },
          
          "size":3,
          "sort": [
            {
              " gross_amount": {
                "order": "desc"
              }
            }
          ], 
          
          "aggs": {
            "filter_status": {
              "filter": { 
                "term": { " status": "1" }
              },
              
              "aggs": {
                "groupby_member_id": {
                  "terms": {
                    "field": " member_id",
                    "size": 3,
                    "order": {
                      "total_gross_amount": "desc"
                    }
                  }, 
                  
                  "aggs": {
                    "total_gross_amount": {
                      "sum": {
                        "field": " gross_amount"
                      }
                    }
                  }
                }
              }
            }
          }
        }
        `)),
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
