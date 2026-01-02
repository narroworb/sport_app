package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/narroworb/data_collector/internal/models"
	"github.com/segmentio/kafka-go"
)

type DatabaseInterface interface {
	IncrementYellowCardsManager(ctx context.Context, managerID uint32) error
	IncrementRedCardsManager(ctx context.Context, managerID uint32) error
	InsertFootballMatchStats(ctx context.Context, stats models.TeamMatchStats, matchID uint32) (uint32, error)
	InsertFootballGoalieMatchStatsBatchNotPointer(ctx context.Context, statsBatch map[uint32]models.GoalieStatsInMatch, matchID uint32) error
	InsertFootballPlayerMatchStatsBatchNotPointer(ctx context.Context, statsBatch map[uint32]models.PlayerStatsInMatch, matchID uint32) error
	InsertFootballTeamTournamentPerformance(ctx context.Context, rowTable *models.TableRow, tournamentID uint32) error
	UpdateFootballTeamTournamentPerformance(ctx context.Context, rowTable *models.TableRow, tournamentID uint32, statID uint32) error
}

type KafkaConsumer struct {
	reader *kafka.Reader
	db     DatabaseInterface
}

func NewKafkaConsumer(groupID string, db DatabaseInterface) *KafkaConsumer {
	broker := os.Getenv("KAFKA_ADDR")
	topic := os.Getenv("KAFKA_TOPIC")

	return &KafkaConsumer{
		db: db,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{broker},
			Topic:          topic,
			GroupID:        groupID,
			CommitInterval: 0,
		}),
	}
}

func (c *KafkaConsumer) Close() {
	c.reader.Close()
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		keyParts := strings.Split(string(m.Key), "|")

		if err := c.processMessage(ctx, keyParts, m.Value); err != nil {
			log.Printf("error processing message: %v", err)
			continue
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("commit failed: %v", err)
		}
	}
}

func (c *KafkaConsumer) processMessage(ctx context.Context, keyParts []string, valueOfMessage []byte) error {
	if len(keyParts) <= 1 {
		return fmt.Errorf("received bad key message")
	}

	switch keyParts[0] {
	case "InsertFootballMatchStats":
		var stats models.TeamMatchStats
		if err := json.Unmarshal([]byte(valueOfMessage), &stats); err != nil {
			return fmt.Errorf("cannot unmarshal bad value message: %v", err)
		}

		if len(keyParts) != 2 {
			return fmt.Errorf("received bad key message for infms")
		}
		matchID, err := strconv.Atoi(keyParts[1])
		if err != nil {
			return fmt.Errorf("received bad matchID in key: %s", keyParts[1])
		}

		_, err = c.db.InsertFootballMatchStats(ctx, stats, uint32(matchID))
		if err != nil {
			return fmt.Errorf("error in inserting team stats: %v", err)
		}

	case "InsertFootballGoalieMatchStatsBatch":
		var stats map[uint32]models.GoalieStatsInMatch
		if err := json.Unmarshal([]byte(valueOfMessage), &stats); err != nil {
			return fmt.Errorf("cannot unmarshal bad value message: %v", err)
		}

		if len(keyParts) != 2 {
			return fmt.Errorf("received bad key message for infgs")
		}
		matchID, err := strconv.Atoi(keyParts[1])
		if err != nil {
			return fmt.Errorf("received bad matchID in key: %s", keyParts[1])
		}

		err = c.db.InsertFootballGoalieMatchStatsBatchNotPointer(ctx, stats, uint32(matchID))
		if err != nil {
			log.Printf("bad goalie match(id=%d) stats: %v", matchID, stats)
			return fmt.Errorf("error in inserting goalies stats: %v", err)
		}

	case "InsertFootballPlayerMatchStatsBatch":
		var stats map[uint32]models.PlayerStatsInMatch
		if err := json.Unmarshal([]byte(valueOfMessage), &stats); err != nil {
			return fmt.Errorf("cannot unmarshal bad value message: %v", err)
		}

		if len(keyParts) != 2 {
			return fmt.Errorf("received bad key message for infps")
		}
		matchID, err := strconv.Atoi(keyParts[1])
		if err != nil {
			return fmt.Errorf("received bad matchID in key: %s", keyParts[1])
		}

		err = c.db.InsertFootballPlayerMatchStatsBatchNotPointer(ctx, stats, uint32(matchID))
		if err != nil {
			log.Printf("bad players match(id=%d) stats: %v", matchID, stats)
			return fmt.Errorf("error in inserting players stats: %v", err)
		}

	case "IncrementYellowCardsManager":
		if len(keyParts) != 2 {
			return fmt.Errorf("received bad key message for incycm")
		}
		managerID, err := strconv.Atoi(keyParts[1])
		if err != nil {
			return fmt.Errorf("received bad managerID in key: %s", keyParts[1])
		}

		err = c.db.IncrementYellowCardsManager(ctx, uint32(managerID))
		if err != nil {
			return fmt.Errorf("error in increment yellow cards manager: %v", err)
		}

	case "IncrementRedCardsManager":
		if len(keyParts) != 2 {
			return fmt.Errorf("received bad key message for incrcm")
		}
		managerID, err := strconv.Atoi(keyParts[1])
		if err != nil {
			return fmt.Errorf("received bad managerID in key: %s", keyParts[1])
		}

		err = c.db.IncrementRedCardsManager(ctx, uint32(managerID))
		if err != nil {
			return fmt.Errorf("error in increment red cards manager: %v", err)
		}
	case "InsertFootballTeamTournamentPerformance":
		if len(keyParts) != 3 {
			return fmt.Errorf("received bad key message for ifttp")
		}
		seasonID, err := strconv.Atoi(keyParts[1])
		if err != nil {
			return fmt.Errorf("received bad seasonID in key: %s", keyParts[1])
		}

		var row models.TableRow
		if err := json.Unmarshal([]byte(valueOfMessage), &row); err != nil {
			return fmt.Errorf("cannot unmarshal bad value message: %v", err)
		}

		err = c.db.InsertFootballTeamTournamentPerformance(ctx, &row, uint32(seasonID))
		if err != nil {
			return fmt.Errorf("error in insert football team tournament performance: %v", err)
		}
	case "UpdateFootballTeamTournamentPerformance":
		if len(keyParts) != 4 {
			return fmt.Errorf("received bad key message for ifttp")
		}
		seasonID, err := strconv.Atoi(keyParts[1])
		if err != nil {
			return fmt.Errorf("received bad seasonID in key: %s", keyParts[1])
		}
		statID, err := strconv.Atoi(keyParts[3])
		if err != nil {
			return fmt.Errorf("received bad statID in key: %s", keyParts[3])
		}

		var row models.TableRow
		if err := json.Unmarshal([]byte(valueOfMessage), &row); err != nil {
			return fmt.Errorf("cannot unmarshal bad value message: %v", err)
		}

		err = c.db.UpdateFootballTeamTournamentPerformance(ctx, &row, uint32(seasonID), uint32(statID))
		if err != nil {
			return fmt.Errorf("error in update football team tournament performance: %v", err)
		}

	default:
		return fmt.Errorf("undefined operation in kafka message: %v", keyParts[0])
	}
	return nil
}

func (c *KafkaConsumer) DrainAndProcess(ctx context.Context, timeout time.Duration) error {
	for {
		ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
		m, err := c.reader.ReadMessage(ctxTimeout)
		cancel()

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		}

		keyParts := strings.Split(string(m.Key), "|")

		if err := c.processMessage(ctx, keyParts, m.Value); err != nil {
			log.Printf("error processing message: %v", err)
			continue
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("commit failed: %v", err)
		}
	}
}
