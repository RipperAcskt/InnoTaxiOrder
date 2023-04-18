package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Elastic struct {
	Client *elasticsearch.Client
	cfg    *config.Config
}

func New(cfg *config.Config) (*Elastic, error) {
	cfgEs := elasticsearch.Config{
		Addresses: []string{cfg.ELASTIC_DB_HOST},
		Username:  cfg.ELASTIC_DB_USERNAME,
		Password:  cfg.ELASTIC_DB_PASSWORD,
	}
	es, err := elasticsearch.NewClient(cfgEs)
	if err != nil {
		return nil, fmt.Errorf("new client failed: %w", err)
	}

	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("info failed: %s", err)
	}

	defer res.Body.Close()

	response, err := es.Indices.Exists([]string{cfg.ELASTIC_DB_NAME})
	if err != nil {
		return nil, fmt.Errorf("exists failed: %w", err)
	}

	if response.StatusCode != 404 {
		return &Elastic{es, cfg}, nil
	}

	response, err = es.Indices.Create(cfg.ELASTIC_DB_NAME)
	if err != nil {
		return nil, fmt.Errorf("create failed: %w", err)
	}

	if response.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("decode err failed: %w", err)
		} else {
			return nil, fmt.Errorf("error: %v", e)
		}
	}

	return &Elastic{es, cfg}, nil
}

func (es *Elastic) CreateOrder(ctx context.Context, order model.Order) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	order.Status = model.StatusWaiting
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshal failed: %s", err)
	}
	req := esapi.IndexRequest{
		Index:   es.cfg.ELASTIC_DB_NAME,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(queryCtx, es.Client)
	if err != nil {
		return fmt.Errorf("req do failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("res error: %w", err)
	}

	return nil
}

func (es *Elastic) GetOrders(ctx context.Context, indexes []string) ([]*model.Order, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	fmt.Println(indexes)
	if len(indexes) == 0 {
		res, err := es.Client.Search(
			es.Client.Search.WithContext(queryCtx),
			es.Client.Search.WithIndex(es.cfg.ELASTIC_DB_NAME),
		)
		if err != nil {
			return nil, fmt.Errorf("search failed: %w", err)
		}

		return es.parseInfo(res)
	}

	var body bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"_id": indexes,
			},
		},
	}
	if err := json.NewEncoder(&body).Encode(query); err != nil {
		return nil, fmt.Errorf("encode failed: %w", err)
	}

	res, err := es.Client.Search(
		es.Client.Search.WithContext(queryCtx),
		es.Client.Search.WithIndex(es.cfg.ELASTIC_DB_NAME),
		es.Client.Search.WithBody(&body),
	)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer res.Body.Close()
	return es.parseInfo(res)
}

func (es *Elastic) parseInfo(res *esapi.Response) ([]*model.Order, error) {
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("decode err failed: %w", err)
		} else {
			return nil, fmt.Errorf("error: %v", e)
		}
	}
	// b, _ := io.ReadAll(res.Body)
	// fmt.Println(string(b))
	var info model.ElasticModel
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	var orders []*model.Order
	for _, el := range info.Hits.Hits {
		element := el
		element.Source.ID = element.ID
		orders = append(orders, &element.Source)
	}
	return orders, nil
}

func (es *Elastic) GetStatus(ctx context.Context, taxiType, status string) ([]*model.Order, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var body bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"Status": status,
						},
					},
					{
						"term": map[string]interface{}{
							"TaxiType": taxiType,
						},
					},
				},
			},
		},

		// "sort": map[string]interface{}{
		// 	"Date": "desc",
		// },
	}

	if err := json.NewEncoder(&body).Encode(query); err != nil {
		return nil, fmt.Errorf("encode failed: %w", err)
	}

	res, err := es.Client.Search(
		es.Client.Search.WithContext(queryCtx),
		es.Client.Search.WithIndex(es.cfg.ELASTIC_DB_NAME),
		es.Client.Search.WithBody(&body),
	)

	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return es.parseInfo(res)

}

func (es *Elastic) UpdateOrder(ctx context.Context, order *model.Order) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	body, err := json.Marshal(&order)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      es.cfg.ELASTIC_DB_NAME,
		DocumentID: order.ID,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
	}

	res, err := req.Do(queryCtx, es.Client)
	if err != nil {
		return fmt.Errorf("req do failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return fmt.Errorf("decode err failed: %w", err)
		} else {
			return fmt.Errorf("error: %v", e)
		}
	}

	return nil
}

func (es *Elastic) GetOrdersByUserID(ctx context.Context, index string, status string) ([]*model.Order, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var body bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"term": map[string]interface{}{
							"Status": status,
						},
					},
					{
						"term": map[string]interface{}{
							"UserID": index,
						},
					},
				},
			},
		},
	}
	if err := json.NewEncoder(&body).Encode(query); err != nil {
		return nil, fmt.Errorf("encode failed: %w", err)
	}

	res, err := es.Client.Search(
		es.Client.Search.WithContext(queryCtx),
		es.Client.Search.WithIndex(es.cfg.ELASTIC_DB_NAME),
		es.Client.Search.WithBody(&body),
	)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer res.Body.Close()
	return es.parseInfo(res)
}
