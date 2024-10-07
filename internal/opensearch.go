package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type OpensearchAPI struct {
	client *opensearch.Client
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

	return &OpensearchAPI{
		client: opClient,
	}, nil
}

func (oapi *OpensearchAPI) CreateIndex(indexName string) (*opensearchapi.Response, error) {
	settings := strings.NewReader(`{
		'settings': {
			'index': {
				'number_of_shards': 1,
				'number_of_replicas': 0
				}
			}
		}`)

	req := opensearchapi.IndicesCreateRequest{
		Index: indexName,
		Body:  settings,
	}

	response, err := req.Do(context.Background(), oapi.client)

	return response, err
}

func (oapi *OpensearchAPI) DeleteIndex(indexName string) (*opensearchapi.Response, error) {

	req := opensearchapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	response, err := req.Do(context.Background(), oapi.client)

	return response, err
}

func (oapi *OpensearchAPI) AddFakeDocuments(indexName string, count int) (*opensearchapi.Response, error) {

	var builder strings.Builder

	for i := 0; i < count; i++ {
		builder.WriteString(fmt.Sprintf(`{ "index" : { "_index" : "go-test-index1", "_id" : "%s" } }`, gofakeit.UUID()))
		builder.WriteString("\n")

		// Add the actual document
		builder.WriteString(fmt.Sprintf(`{ "title" : "Movie %s", "director" : "Director %s", "year" : "%d" }`, gofakeit.MovieName(), gofakeit.Name(), gofakeit.Year()))
		builder.WriteString("\n")
	}

	// Convert the built string to an io.Reader
	reader := strings.NewReader(builder.String())

	blk, err := oapi.client.Bulk(reader)

	return blk, err
}
