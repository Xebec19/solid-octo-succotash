package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	opensearch "github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

type OpensearchAPI struct {
	Client *opensearch.Client
}

func createOpensearchClient(host, port, username, password string) (*opensearch.Client, error) {

	address := fmt.Sprintf("%s:%s", host, port)

	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{address},
		Username:  username,
		Password:  password,
	})

	return client, err
}

func NewOpensearchClient(host, port, username, password string) (*OpensearchAPI, error) {

	opClient, err := createOpensearchClient(host, port, username, password)

	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprint(opClient.Info()))

	return &OpensearchAPI{
		Client: opClient,
	}, nil
}

func (oapi *OpensearchAPI) CreateIndex(indexName string) (*opensearchapi.Response, error) {
	settings := strings.NewReader(`{
		"settings": {
		  "index": {
			   "number_of_shards": 1,
			   "number_of_replicas": 2
			   }
			 }
		}`)

	// Create an index with non-default settings.
	res := opensearchapi.IndicesCreateRequest{
		Index: indexName,
		Body:  settings,
	}

	slog.Info(fmt.Sprint(res))

	response, err := res.Do(context.Background(), oapi.Client)

	return response, err
}

func (oapi *OpensearchAPI) DeleteIndex(indexName string) (*opensearchapi.Response, error) {

	req := opensearchapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	response, err := req.Do(context.Background(), oapi.Client)

	return response, err
}

func (oapi *OpensearchAPI) AddFakeDocuments(indexName string, count int) (*opensearchapi.Response, error) {

	var builder strings.Builder

	for i := 0; i < count; i++ {
		builder.WriteString(fmt.Sprintf(`{ "index" : { "_index" : "%s", "_id" : "%s" } }`, indexName, gofakeit.UUID()))
		builder.WriteString("\n")

		// Add the actual document
		builder.WriteString(fmt.Sprintf(`{ "title" : "%s", "director" : "%s", "year" : "%d" }`, gofakeit.MovieName(), gofakeit.Name(), gofakeit.Year()))
		builder.WriteString("\n")
	}

	// Convert the built string to an io.Reader
	reader := strings.NewReader(builder.String())

	blk, err := oapi.Client.Bulk(reader)

	return blk, err
}

func (opai *OpensearchAPI) SearchData(indexName string) (*opensearchapi.Response, error) {
	res, err := opai.Client.Search(
		opai.Client.Search.WithIndex(indexName),
		opai.Client.Search.WithSize(3),
		opai.Client.Search.WithScroll(time.Minute),
		opai.Client.Search.WithSort("_id"),
		// opai.Client.Search.WithSize(10),
	)

	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprint(res))

	return res, err

}

func (opai *OpensearchAPI) SearchDataWithScroll(indexName, scrollID string) (*opensearchapi.Response, error) {

	res, err := opai.Client.Scroll(
		// opai.Client.Search.WithIndex(indexName),
		// opai.Client.Search.WithScrollID(scrollID),
		// opai.Client.Search.WithSort("")
		// opai.Client.Search.WithSize(10),
		// opai.Client.Search.WithScroll(time.Minute),
		opai.Client.Scroll.WithScrollID(scrollID),
		opai.Client.Scroll.WithScroll(time.Minute))

	if err != nil {
		return nil, err
	}

	slog.Info(fmt.Sprint(res))

	return res, err

	// search := opensearchapi.SearchRequest{
	// 	Index: []string{indexName},
	// 	Body: strings.NewReader(fmt.Sprintf(`{
	// 		"scroll": "1m",
	// 		"scroll_id": "%s"
	// 	}`, scrollID)),
	// }

	// searchResponse, err := search.Do(context.Background(), opai.Client)

	// if err != nil {
	// 	return nil, err
	// }

	// return searchResponse, err
}
