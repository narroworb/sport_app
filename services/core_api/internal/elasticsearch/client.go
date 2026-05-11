package elasticsearch

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/narroworb/core_api/internal/models"
)

type Elasticsearch struct {
	client *elasticsearch.Client
}

func NewElasticsearch() (*Elasticsearch, error) {
	addr := os.Getenv("ELASTICSEARCH_ADDR")
	if addr == "" {
		addr = "http://localhost:9200"
	}

	cfg := elasticsearch.Config{
		Addresses: []string{addr},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := client.Info(client.Info.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to elasticsearch: %w", err)
	}
	res.Body.Close()

	es := &Elasticsearch{client: client}
	if err := es.createIndices(ctx); err != nil {
		return nil, fmt.Errorf("failed to create indices: %w", err)
	}

	return es, nil
}

func (es *Elasticsearch) createIndices(ctx context.Context) error {
	indices := []struct {
		name   string
		schema string
	}{
		{
			name: "players",
			schema: `{
				"mappings": {
					"properties": {
						"id": {"type": "integer"},
						"first_name": {"type": "text", "analyzer": "standard"},
						"last_name": {"type": "text", "analyzer": "standard"},
						"position": {"type": "keyword"},
						"nation": {"type": "keyword"},
						"current_status": {"type": "keyword"},
						"height": {"type": "integer"},
						"url_photo": {"type": "keyword"}
					}
				}
			}`,
		},
		{
			name: "teams",
			schema: `{
				"mappings": {
					"properties": {
						"id": {"type": "integer"},
						"name": {"type": "text", "analyzer": "standard"},
						"url_logo": {"type": "keyword"}
					}
				}
			}`,
		},
		{
			name: "managers",
			schema: `{
				"mappings": {
					"properties": {
						"id": {"type": "integer"},
						"first_name": {"type": "text", "analyzer": "standard"},
						"last_name": {"type": "text", "analyzer": "standard"},
						"nation": {"type": "keyword"},
						"url_photo": {"type": "keyword"}
					}
				}
			}`,
		},
		{
			name: "tournaments",
			schema: `{
				"mappings": {
					"properties": {
						"id": {"type": "integer"},
						"name": {"type": "text", "analyzer": "standard"},
						"country": {"type": "keyword"},
						"season": {"type": "keyword"},
						"url_logo": {"type": "keyword"}
					}
				}
			}`,
		},
	}

	for _, idx := range indices {
		res, err := es.client.Indices.Exists([]string{idx.name})
		if err != nil {
			return fmt.Errorf("failed to check index %s: %w", idx.name, err)
		}

		if res.StatusCode != 200 {
			res, err = es.client.Indices.Create(
				idx.name,
				es.client.Indices.Create.WithBody(strings.NewReader(idx.schema)),
			)
			if err != nil {
				return fmt.Errorf("failed to create index %s: %w", idx.name, err)
			}
			res.Body.Close()
		} else {
			res.Body.Close()
		}
	}

	return nil
}

func (es *Elasticsearch) Search(ctx context.Context, query string, entityType string, filters models.SearchFilters) ([]models.SearchResult, error) {
	if strings.TrimSpace(query) == "" {
		return []models.SearchResult{}, nil
	}

	var indices []string
	switch entityType {
	case "player":
		indices = []string{"players"}
	case "team":
		indices = []string{"teams"}
	case "manager":
		indices = []string{"managers"}
	case "tournament":
		indices = []string{"tournaments"}
	case "all":
		indices = []string{"players", "teams", "managers", "tournaments"}
	default:
		indices = []string{"players", "teams", "managers", "tournaments"}
	}

	searchQuery := es.buildSearchQuery(query, entityType, filters)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("failed to encode search query: %w", err)
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex(indices...),
		es.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: status %d", res.StatusCode)
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Score  float32         `json:"_score"`
				Index  string          `json:"_index"`
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	var results []models.SearchResult
	for _, hit := range result.Hits.Hits {
		var entity interface{}
		var eType string

		switch hit.Index {
		case "players":
			var player models.Player
			if err := json.Unmarshal(hit.Source, &player); err != nil {
				continue
			}
			entity = player
			eType = "player"
		case "teams":
			var team models.Team
			if err := json.Unmarshal(hit.Source, &team); err != nil {
				continue
			}
			entity = team
			eType = "team"
		case "managers":
			var manager models.Manager
			if err := json.Unmarshal(hit.Source, &manager); err != nil {
				continue
			}
			entity = manager
			eType = "manager"
		case "tournaments":
			var tournament models.Tournament
			if err := json.Unmarshal(hit.Source, &tournament); err != nil {
				continue
			}
			entity = tournament
			eType = "tournament"
		}

		results = append(results, models.SearchResult{
			Type:  eType,
			Score: hit.Score,
			Data:  entity,
		})
	}

	return results, nil
}

func (es *Elasticsearch) buildSearchQuery(query string, entityType string, filters models.SearchFilters) map[string]interface{} {
	from := (filters.Page - 1) * filters.PageSize
	if filters.Page < 1 {
		from = 0
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 10
	}

	var must []map[string]interface{}

	if strings.TrimSpace(query) != "" {
		multiMatch := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     query,
				"fields":    es.getSearchFields(entityType),
				"fuzziness": "AUTO",
				"operator":  "or",
			},
		}
		must = append(must, multiMatch)
	}

	if filters.Position != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"position": filters.Position,
			},
		})
	}

	if filters.Nation != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"nation": filters.Nation,
			},
		})
	}

	if filters.Season != "" {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"season": filters.Season,
			},
		})
	}

	return map[string]interface{}{
		"from": from,
		"size": filters.PageSize,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"sort": []map[string]interface{}{
			{
				"_score": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}
}

func (es *Elasticsearch) getSearchFields(entityType string) []string {
	switch entityType {
	case "player":
		return []string{"first_name^2", "last_name^2", "position", "nation"}
	case "team":
		return []string{"name^3"}
	case "manager":
		return []string{"first_name^2", "last_name^2", "nation"}
	case "tournament":
		return []string{"name^3", "country"}
	default:
		return []string{"first_name", "last_name", "name", "position", "nation", "country"}
	}
}

func (es *Elasticsearch) IndexPlayer(ctx context.Context, player models.Player) error {
	return es.indexDocument(ctx, "players", player.ID, player)
}

func (es *Elasticsearch) IndexTeam(ctx context.Context, team models.Team) error {
	return es.indexDocument(ctx, "teams", team.ID, team)
}

func (es *Elasticsearch) IndexManager(ctx context.Context, manager models.Manager) error {
	return es.indexDocument(ctx, "managers", manager.ID, manager)
}

func (es *Elasticsearch) IndexTournament(ctx context.Context, tournament models.Tournament) error {
	return es.indexDocument(ctx, "tournaments", tournament.ID, tournament)
}

func (es *Elasticsearch) indexDocument(ctx context.Context, index string, id uint32, doc interface{}) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	res, err := es.client.Index(
		index,
		strings.NewReader(string(data)),
		es.client.Index.WithDocumentID(fmt.Sprintf("%d", id)),
		es.client.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to index document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index error: status %d", res.StatusCode)
	}

	return nil
}

func (es *Elasticsearch) DeleteIndex(ctx context.Context, index string) error {
	res, err := es.client.Indices.Delete([]string{index})
	if err != nil {
		return fmt.Errorf("failed to delete index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("delete index error: status %d", res.StatusCode)
	}

	return nil
}

func hashFilters(filters models.SearchFilters) string {
	data := fmt.Sprintf("%s:%s:%s:%s:%d:%d",
		filters.Query,
		filters.Position,
		filters.Nation,
		filters.Season,
		filters.Page,
		filters.PageSize,
	)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
