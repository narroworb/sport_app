package elasticsearch

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/narroworb/core_api/internal/handlers"
	"github.com/narroworb/core_api/internal/models"
)

type Elasticsearch struct {
	client *elasticsearch.Client
	db     handlers.AnalyticDatabaseInterface
}

// Структуры для парсинга результатов Elasticsearch
type esPlayer struct {
	ID            uint32 `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Position      string `json:"position"`
	Nation        string `json:"nation"`
	CurrentStatus string `json:"current_status"`
	Height        uint16 `json:"height"`
	URLPhoto      string `json:"url_photo"`
}

type esManager struct {
	ID        uint32 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nation    string `json:"nation"`
	URLPhoto  string `json:"url_photo"`
}

type esTournament struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	Season   string `json:"season"`
	URLLogo  string `json:"url_logo"`
}

func NewElasticsearch(db handlers.AnalyticDatabaseInterface) (*Elasticsearch, error) {
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

	es := &Elasticsearch{client: client, db: db}
	if err := es.createIndices(ctx); err != nil {
		return nil, fmt.Errorf("failed to create indices: %w", err)
	}

	go es.endlessIndexing()

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
			var player esPlayer
			if err := json.Unmarshal(hit.Source, &player); err != nil {
				log.Printf("error in parsing player from esearch %s\n", hit.Source)
				continue
			}
			entity = player
			eType = "player"
		case "teams":
			var team models.Team
			if err := json.Unmarshal(hit.Source, &team); err != nil {
				log.Printf("error in parsing team from esearch %s\n", hit.Source)
				continue
			}
			entity = team
			eType = "team"
		case "managers":
			var manager esManager
			if err := json.Unmarshal(hit.Source, &manager); err != nil {
				log.Printf("error in parsing manager from esearch %s\n", hit.Source)
				continue
			}
			entity = manager
			eType = "manager"
		case "tournaments":
			var tournament esTournament
			if err := json.Unmarshal(hit.Source, &tournament); err != nil {
				log.Printf("error in parsing tournament from esearch %s\n", hit.Source)
				continue
			}
			entity = tournament
			eType = "tournament"
		default:
			continue
		}

		if entity == nil {
			continue
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

func (es *Elasticsearch) indexPlayer(ctx context.Context, player models.Player) error {
	esPlayer := map[string]interface{}{
		"id":             player.ID,
		"first_name":     player.FirstName,
		"last_name":      player.LastName,
		"position":       player.Position,
		"nation":         player.Nation.Name,
		"current_status": player.CurrentStatus,
		"height":         player.Height,
		"url_photo":      player.URLPhoto,
	}

	return es.indexDocument(ctx, "players", player.ID, esPlayer)
}

func (es *Elasticsearch) indexTeam(ctx context.Context, team models.Team) error {
	return es.indexDocument(ctx, "teams", team.ID, team)
}

func (es *Elasticsearch) indexManager(ctx context.Context, manager models.Manager) error {
	esManager := map[string]interface{}{
		"id":         manager.ID,
		"first_name": manager.FirstName,
		"last_name":  manager.LastName,
		"nation":     manager.Nation.Name,
		"url_photo":  manager.URLPhoto,
	}

	return es.indexDocument(ctx, "managers", manager.ID, esManager)
}

func (es *Elasticsearch) indexTournament(ctx context.Context, tournament models.Tournament) error {
	esTournament := map[string]interface{}{
		"id":       tournament.ID,
		"name":     tournament.Name,
		"country":  tournament.Country.Name,
		"season":   tournament.Season,
		"url_logo": tournament.URLLogo,
	}

	return es.indexDocument(ctx, "tournaments", tournament.ID, esTournament)
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
		var errMsg struct {
			Error struct {
				Type      string `json:"type"`
				Reason    string `json:"reason"`
				RootCause []struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"root_cause"`
			} `json:"error"`
			Status int `json:"status"`
		}

		if err := json.NewDecoder(res.Body).Decode(&errMsg); err != nil {
			return fmt.Errorf("index error: status %d, body: %s", res.StatusCode, res.String())
		}

		return fmt.Errorf("index error: %s - %s (status %d)",
			errMsg.Error.Type, errMsg.Error.Reason, res.StatusCode)
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
