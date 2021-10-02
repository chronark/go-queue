package main

import (
	"log"
	"os"

	"github.com/chronark/go-queue/fauna"
	f "github.com/fauna/faunadb-go/v4/faunadb"
)

func main() {

	faunaKey := os.Getenv("FAUNA_KEY")
	if faunaKey == "" {
		panic("FAUNA_KEY not set")
	}

	client := f.NewFaunaClient(faunaKey)

	queries := []f.Expr{
		f.CreateCollection(f.Obj{
			"name":         fauna.COLLECTION_TODO,
			"history_days": 0,
		}),
		f.CreateCollection(f.Obj{
			"name":         fauna.COLLECTION_IN_PROGRESS,
			"history_days": 0,
		}),
		f.CreateCollection(f.Obj{
			"name":         fauna.COLLECTION_DONE,
			"history_days": 0,
			"ttl_days":     30,
		}),
		f.CreateIndex(f.Obj{
			"name":       fauna.INDEX_TODO_BY_TOPIC,
			"source":     f.Collection(fauna.COLLECTION_TODO),
			"serialized": true,
			"terms": []f.Obj{
				{
					"field": []string{"data", "header", "topic"},
				},
			},
			"values": []f.Obj{
				{"field": []string{"ts"}},
				{"field": []string{"ref"}},
				{"field": []string{"data"}},
			},
		}),

		f.CreateIndex(f.Obj{
			"name":       fauna.INDEX_IN_PROGRESS_BY_ID,
			"source":     f.Collection(fauna.COLLECTION_IN_PROGRESS),
			"unique":     true,
			"serialized": true,
			"terms": []f.Obj{
				{
					"field": []string{"data", "header", "id"},
				},
			},
		}),
	}

	for _, query := range queries {
		res, err := client.Query(query)
		if err != nil {
			panic(err)
		}
		log.Printf("%+v", res)
	}

}
