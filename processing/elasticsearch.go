// Package processing provide all requirements to process change data capture
package processing

import (
	"context"
	"net/http"

	"github.com/Lord-Y/synker/commons"
	"github.com/olivere/elastic/v7"
)

// ePing permit to get elasticsearch status
func (c *Validate) ePing() (b bool) {
	var (
		code   int
		client *elastic.Client
		err    error
	)
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(commons.GetElasticsearchURI()), elastic.SetGzip(true))
	if err != nil {
		c.Logger.Error().Err(err).Msg("Error occured while creating ES client")
		return
	}
	defer client.Stop()
	_, code, err = client.Ping(commons.GetElasticsearchURI()).HttpHeadOnly(true).Do(context.TODO())
	if code != http.StatusOK || err != nil {
		c.Logger.Error().Err(err).Msgf("Error occured while pinging ES http status %d", code)
		return
	}
	return true
}

// eClient permit to create client connection to elasticsearch
func (c *Validate) eClient() (client *elastic.Client, err error) {
	client, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(commons.GetElasticsearchURI()), elastic.SetGzip(true))
	if err != nil {
		return
	}
	return
}

// indexAlreadyExist permit to check if elasticsearch index already exist
func (c *Validate) indexAlreadyExist(client *elastic.Client, index string) (b bool, err error) {
	defer client.Stop()

	ctx := context.Background()
	b, err = client.IndexExists(index).Do(ctx)
	if err != nil {
		return
	}
	return
}

// createIndex permit to create elasticsearch index with mapping provided
func (c *Validate) createIndex(client *elastic.Client, index string, mapping string) (b bool, err error) {
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

// deleteIndex permit to delete elasticsearch index provided
func (c *Validate) deleteIndex(client *elastic.Client, index string) (b bool, err error) {
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
