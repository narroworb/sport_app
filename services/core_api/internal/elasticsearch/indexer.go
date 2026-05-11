package elasticsearch

import (
	"context"
	"log"
	"time"
)

func (e *Elasticsearch) RunIndexing(ctx context.Context) error {
	tournaments, err := e.db.GetUnindexedTournaments(ctx)
	if err != nil {
		return err
	}
	if len(tournaments) > 0 {
		log.Printf("Indexing %d tournaments\n", len(tournaments))
		updatedIDs := make([]uint32, 0, len(tournaments))
		for _, tournament := range tournaments {
			if err := e.indexTournament(ctx, tournament); err != nil {
				log.Printf("Failed to index tournament %d: %v\n", tournament.ID, err)
			} else {
				updatedIDs = append(updatedIDs, tournament.ID)
			}
		}

		if len(updatedIDs) > 0 {
			if err := e.db.UpdateBatchTournamentsIndexedStatus(ctx, updatedIDs); err != nil {
				log.Printf("Failed to update indexed status for tournaments: %v\n", err)
			}
			log.Printf("Indexed %d tournaments\n", len(updatedIDs))
		}
	}

	teams, err := e.db.GetUnindexedTeams(ctx)
	if err != nil {
		return err
	}
	if len(teams) > 0 {
		log.Printf("Indexing %d teams\n", len(teams))
		updatedIDs := make([]uint32, 0, len(teams))
		for _, team := range teams {
			if err := e.indexTeam(ctx, team); err != nil {
				log.Printf("Failed to index team %d: %v\n", team.ID, err)
			} else {
				updatedIDs = append(updatedIDs, team.ID)
			}
		}
		if len(updatedIDs) > 0 {
			if err := e.db.UpdateBatchTeamsIndexedStatus(ctx, updatedIDs); err != nil {
				log.Printf("Failed to update indexed status for teams: %v\n", err)
			}
			log.Printf("Indexed %d teams\n", len(updatedIDs))
		}
	}

	mgrs, err := e.db.GetUnindexedManagers(ctx)
	if err != nil {
		return err
	}
	if len(mgrs) > 0 {
		log.Printf("Indexing %d managers\n", len(mgrs))
		updatedIDs := make([]uint32, 0, len(mgrs))
		for _, manager := range mgrs {
			if err := e.indexManager(ctx, manager); err != nil {
				log.Printf("Failed to index manager %d: %v\n", manager.ID, err)
			} else {
				updatedIDs = append(updatedIDs, manager.ID)
			}
		}
		if len(updatedIDs) > 0 {
			if err := e.db.UpdateBatchManagersIndexedStatus(ctx, updatedIDs); err != nil {
				log.Printf("Failed to update indexed status for managers: %v\n", err)
			} else {
			log.Printf("Indexed %d managers\n", len(updatedIDs))
		}
	}

	players, err := e.db.GetUnindexedPlayers(ctx)
	if err != nil {
		return err
	}
	if len(players) > 0 {
		log.Printf("Indexing %d players\n", len(players))
		updatedIDs := make([]uint32, 0, len(players))
		for _, player := range players {
			if err := e.indexPlayer(ctx, player); err != nil {
				log.Printf("Failed to index player %d: %v\n", player.ID, err)
			} else {
				updatedIDs = append(updatedIDs, player.ID)
			}
		}
		if len(updatedIDs) > 0 {
			if err := e.db.UpdateBatchPlayersIndexedStatus(ctx, updatedIDs); err != nil {
				log.Printf("Failed to update indexed status for players: %v\n", err)
			} else {
			log.Printf("Indexed %d players\n", len(updatedIDs))
			}
		}
	}

	return nil
}

func (es *Elasticsearch) endlessIndexing() {
	for {
		if err := es.RunIndexing(context.Background()); err != nil {
			log.Printf("Indexing error: %v\n", err)
		}
		time.Sleep(6 * time.Hour)
	}
}
