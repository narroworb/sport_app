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
	InsertFootballTeamTournamentPerformanceBatch(ctx context.Context, performanceBatch []models.TableRow, tournamentID uint32) error
	UpdateFootballTeamTournamentPerformanceBatch(ctx context.Context, performanceBatch map[uint32]models.TableRow, tournamentID uint32) error
	InsertFootballManagerWithID(ctx context.Context, manager *models.Manager) error
	UpdateFootballMatch(ctx context.Context, match *models.Match) error
	InsertFootballMatchWithID(ctx context.Context, match *models.Match, tournamentID uint32) error
	InsertFootballPlayerWithID(ctx context.Context, player *models.Player) error
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
			log.Printf("bad goalie match(id=%d) stats: %v\n", matchID, stats)
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
			log.Printf("bad players match(id=%d) stats: %v\n", matchID, stats)
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
	case "UpdateFootballTeamTournamentPerformanceBatch":
		if len(keyParts) != 2 {
			return fmt.Errorf("received bad key message for ufttpb")
		}
		seasonID, err := strconv.Atoi(keyParts[1])
		if err != nil {
			return fmt.Errorf("received bad seasonID in key: %s", keyParts[1])
		}

		table := make(map[uint32]models.TableRow)
		if err := json.Unmarshal([]byte(valueOfMessage), &table); err != nil {
			return fmt.Errorf("cannot unmarshal bad value message: %v", err)
		}

		err = c.db.UpdateFootballTeamTournamentPerformanceBatch(ctx, table, uint32(seasonID))
		if err != nil {
			return fmt.Errorf("error in update football team tournament performance batch: %v", err)
		}
	case "InsertFootballTeamTournamentPerformanceBatch":
		if len(keyParts) != 2 {
			return fmt.Errorf("received bad key message for ifttpb")
		}
		seasonID, err := strconv.Atoi(keyParts[1])
		if err != nil {
			return fmt.Errorf("received bad seasonID in key: %s", keyParts[1])
		}

		table := make([]models.TableRow, 0, 20)
		if err := json.Unmarshal([]byte(valueOfMessage), &table); err != nil {
			return fmt.Errorf("cannot unmarshal bad value message: %v", err)
		}

		err = c.db.InsertFootballTeamTournamentPerformanceBatch(ctx, table, uint32(seasonID))
		if err != nil {
			return fmt.Errorf("error in insert football team tournament performance batch: %v", err)
		}
	case "InsertFootballManager":
		var manager models.Manager
		if err := json.Unmarshal([]byte(valueOfMessage), &manager); err != nil {
			return fmt.Errorf("cannot unmarshal bad value message: %v", err)
		}

		if len(keyParts) != 2 {
			return fmt.Errorf("received bad key message for infmanager")
		}
		managerID, err := strconv.Atoi(keyParts[1])
		if err != nil || managerID != int(manager.ID) {
			return fmt.Errorf("received bad managerID in key: %s", keyParts[1])
		}

		err = c.db.InsertFootballManagerWithID(ctx, &manager)
		if err != nil {
			return fmt.Errorf("error in inserting manager: %v", err)
		}
	case "UpdateFootballMatch":
		var match models.Match
		if err := json.Unmarshal([]byte(valueOfMessage), &match); err != nil {
			return fmt.Errorf("cannot unmarshal bad value message: %v", err)
		}

		if len(keyParts) != 2 {
			return fmt.Errorf("received bad key message for upfmatch")
		}
		matchID, err := strconv.Atoi(keyParts[1])
		if err != nil || matchID != int(match.IDAppDB) {
			return fmt.Errorf("received bad matchID in key: %s", keyParts[1])
		}

		err = c.db.UpdateFootballMatch(ctx, &match)
		if err != nil {
			return fmt.Errorf("error in updating match: %v", err)
		}
	case "InsertFootballMatch":
		var match models.Match
		if err := json.Unmarshal([]byte(valueOfMessage), &match); err != nil {
			return fmt.Errorf("cannot unmarshal bad value message: %v", err)
		}

		if len(keyParts) != 3 {
			return fmt.Errorf("received bad key message for infmatch")
		}
		seasonID, err := strconv.Atoi(keyParts[2])
		if err != nil {
			return fmt.Errorf("received bad seasonID in key: %s", keyParts[2])
		}

		err = c.db.InsertFootballMatchWithID(ctx, &match, uint32(seasonID))
		if err != nil {
			return fmt.Errorf("error in inserting match: %v", err)
		}
	case "InsertFootballPlayer":
		var player models.Player
		if err := json.Unmarshal([]byte(valueOfMessage), &player); err != nil {
			return fmt.Errorf("cannot unmarshal bad value message: %v", err)
		}

		if len(keyParts) != 2 {
			return fmt.Errorf("received bad key message for infplayer")
		}
		playerID, err := strconv.Atoi(keyParts[1])
		if err != nil || playerID != int(player.ID) {
			return fmt.Errorf("received bad playerID in key: %s", keyParts[1])
		}

		err = c.db.InsertFootballPlayerWithID(ctx, &player)
		if err != nil {
			return fmt.Errorf("error in inserting player: %v", err)
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
