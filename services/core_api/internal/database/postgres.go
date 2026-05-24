package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	conn *sqlx.DB
}

func NewPostgresDB() (*PostgresDB, error) {
	dsn := os.Getenv("POSTGRES_URL")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresDB{db}, nil
}

func (p *PostgresDB) GetFavoritePlayersIDs(ctx context.Context, userID int64) ([]uint32, error) {
	rows, err := p.conn.QueryContext(ctx, `SELECT athlete_id FROM user_favorite_athletes WHERE user_id=$1;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	athleteIDs := make([]uint32, 0, 10)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		athleteIDs = append(athleteIDs, uint32(id))
	}

	return athleteIDs, rows.Err()
}

func (p *PostgresDB) GetFavoriteManagersIDs(ctx context.Context, userID int64) ([]uint32, error) {
	rows, err := p.conn.QueryContext(ctx, `SELECT manager_id FROM user_favorite_managers WHERE user_id=$1;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	managerIDs := make([]uint32, 0, 10)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		managerIDs = append(managerIDs, uint32(id))
	}

	return managerIDs, rows.Err()
}

func (p *PostgresDB) GetFavoriteTeamsIDs(ctx context.Context, userID int64) ([]uint32, error) {
	rows, err := p.conn.QueryContext(ctx, `SELECT team_id FROM user_favorite_teams WHERE user_id=$1;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teamIDs := make([]uint32, 0, 10)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		teamIDs = append(teamIDs, uint32(id))
	}

	return teamIDs, rows.Err()
}

func (p *PostgresDB) GetFavoriteTournamentIDs(ctx context.Context, userID int64) ([]uint32, error) {
	rows, err := p.conn.QueryContext(ctx, `SELECT tournament_id FROM user_favorite_tournaments WHERE user_id=$1;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tournamentIDs := make([]uint32, 0, 10)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		tournamentIDs = append(tournamentIDs, uint32(id))
	}

	return tournamentIDs, rows.Err()
}

func (p *PostgresDB) SetFavoritePlayerByID(ctx context.Context, userID int64, playerID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM user_favorite_athletes WHERE user_id=$1 AND athlete_id=$2);`, userID, playerID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("player with ID=%d already favorite for user with ID=%d", playerID, userID)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_favorite_athletes(user_id, athlete_id) VALUES($1, $2);`, userID, playerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresDB) SetFavoriteManagerByID(ctx context.Context, userID int64, managerID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM user_favorite_managers WHERE user_id=$1 AND manager_id=$2);`, userID, managerID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("manager with ID=%d already favorite for user with ID=%d", managerID, userID)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_favorite_managers(user_id, manager_id) VALUES($1, $2);`, userID, managerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresDB) SetFavoriteTeamByID(ctx context.Context, userID int64, teamID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM user_favorite_teams WHERE user_id=$1 AND team_id=$2);`, userID, teamID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("team with ID=%d already favorite for user with ID=%d", teamID, userID)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_favorite_teams(user_id, team_id) VALUES($1, $2);`, userID, teamID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresDB) SetFavoriteTournamentByID(ctx context.Context, userID int64, tournamentID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM user_favorite_tournaments WHERE user_id=$1 AND tournament_id=$2);`, userID, tournamentID).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("tournament with ID=%d already favorite for user with ID=%d", tournamentID, userID)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_favorite_tournaments(user_id, tournament_id) VALUES($1, $2);`, userID, tournamentID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresDB) DeleteFavoritePlayerByID(ctx context.Context, userID int64, playerID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM user_favorite_athletes WHERE user_id=$1 AND athlete_id=$2);`, userID, playerID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("player with ID=%d already not favorite for user with ID=%d", playerID, userID)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM user_favorite_athletes WHERE user_id=$1 AND athlete_id=$2;`, userID, playerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresDB) DeleteFavoriteManagerByID(ctx context.Context, userID int64, managerID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM user_favorite_managers WHERE user_id=$1 AND manager_id=$2);`, userID, managerID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("manager with ID=%d already not favorite for user with ID=%d", managerID, userID)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM user_favorite_managers WHERE user_id=$1 AND manager_id=$2;`, userID, managerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresDB) DeleteFavoriteTeamByID(ctx context.Context, userID int64, teamID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM user_favorite_teams WHERE user_id=$1 AND team_id=$2);`, userID, teamID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("team with ID=%d already not favorite for user with ID=%d", teamID, userID)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM user_favorite_teams WHERE user_id=$1 AND team_id=$2;`, userID, teamID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresDB) DeleteFavoriteTournamentByID(ctx context.Context, userID int64, tournamentID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM user_favorite_tournaments WHERE user_id=$1 AND tournament_id=$2);`, userID, tournamentID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("tournament with ID=%d already not favorite for user with ID=%d", tournamentID, userID)
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM user_favorite_tournaments WHERE user_id=$1 AND tournament_id=$2;`, userID, tournamentID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
