// Package elasticsearch permit to push or delete all datas from db to elasticsearch
package elasticsearch

import (
	"context"
	"net/http"

	"github.com/Lord-Y/synker/commons"
	"github.com/olivere/elastic/v7"
	"github.com/rs/zerolog/log"
)

// Ping permit to get elasticsearch status
func Ping() (b bool) {
	var (
		code   int
		client *elastic.Client
		err    error
	)
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(commons.GetElasticsearchURI()), elastic.SetGzip(true))
	if err != nil {
		log.Error().Err(err).Msg("Error occured while creating ES client")
		return
	}
	defer client.Stop()
	_, code, err = client.Ping(commons.GetElasticsearchURI()).HttpHeadOnly(true).Do(context.TODO())
	if code != http.StatusOK || err != nil {
		log.Error().Err(err).Msgf("Error occured while pinging ES http status %d", code)
		return
	}
	return true
}

// Client permit to create client connection to elasticsearch
func Client() (client *elastic.Client, err error) {
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(commons.GetElasticsearchURI()), elastic.SetGzip(true))
	if err != nil {
		return
	}
	return
}

// IndexAlreadyExist permit to check if elasticsearch index already exist
func IndexAlreadyExist(client *elastic.Client, index string) (b bool, err error) {
	defer client.Stop()

	ctx := context.Background()
	b, err = client.IndexExists(index).Do(ctx)
	if err != nil {
		return
	}
	return
}

// createIndex permit to create elasticsearch index with mapping provided
func createIndex(client *elastic.Client, index string, mapping string) (b bool, err error) {
	defer client.Stop()

	ctx := context.Background()
	b, err = client.IndexExists(index).Do(ctx)
	if err != nil {
		return
	}
	if !b {
		create, err := client.CreateIndex(index).BodyString(mapping).Do(ctx)
		if err != nil {
			return false, err
		}
		if !create.Acknowledged {
			return false, err
		}
	}
	return true, nil
}

// DeleteIndex permit to delete elasticsearch index provided
func DeleteIndex(client *elastic.Client, index string) (b bool, err error) {
	defer client.Stop()
	ctx := context.Background()

	resp, err := client.DeleteIndex(index).Do(ctx)
	if err != nil {
		return
	}
	if !resp.Acknowledged {
		return false, err
	}
	return true, nil
}
