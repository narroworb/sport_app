package database

import (
	"context"
	"database/sql"
	"errors"
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
	tx, err := p.conn.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, `SELECT athlete_id FROM user_favorite_athletes WHERE user_id=$1;`, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	athleteIDs := make([]uint32, 0, 10)
	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return nil, err
		}

		athleteIDs = append(athleteIDs, uint32(id))
	}

	return athleteIDs, tx.Commit()
}

func (p *PostgresDB) GetFavoriteManagersIDs(ctx context.Context, userID int64) ([]uint32, error) {
	tx, err := p.conn.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, `SELECT manager_id FROM user_favorite_managers WHERE user_id=$1;`, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	managerIDs := make([]uint32, 0, 10)
	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return nil, err
		}

		managerIDs = append(managerIDs, uint32(id))
	}

	return managerIDs, tx.Commit()
}

func (p *PostgresDB) GetFavoriteTeamsIDs(ctx context.Context, userID int64) ([]uint32, error) {
	tx, err := p.conn.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, `SELECT team_id FROM user_favorite_teams WHERE user_id=$1;`, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	teamIDs := make([]uint32, 0, 10)
	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return nil, err
		}

		teamIDs = append(teamIDs, uint32(id))
	}

	return teamIDs, tx.Commit()
}

func (p *PostgresDB) GetFavoriteTournamentIDs(ctx context.Context, userID int64) ([]uint32, error) {
	tx, err := p.conn.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, `SELECT tournament_id FROM user_favorite_tournaments WHERE user_id=$1;`, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tournamentIDs := make([]uint32, 0, 10)
	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return nil, err
		}

		tournamentIDs = append(tournamentIDs, uint32(id))
	}

	return tournamentIDs, tx.Commit()
}

func (p *PostgresDB) SetFavoritePlayerByID(ctx context.Context, userID int64, playerID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.QueryContext(ctx, `SELECT id FROM user_favorite_athletes WHERE user_id=$1 AND athlete_id=$2;`, userID, playerID)
	if err == nil {
		tx.Rollback()
		return fmt.Errorf("player with ID=%d already favorite", playerID)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_favorite_athletes(user_id, athlete_id) VALUES($1, $2);`, userID, playerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *PostgresDB) SetFavoriteManagerByID(ctx context.Context, userID int64, managerID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.QueryContext(ctx, `SELECT id FROM user_favorite_managers WHERE user_id=$1 AND manager_id=$2;`, userID, managerID)
	if err == nil {
		tx.Rollback()
		return fmt.Errorf("manager with ID=%d already favorite", managerID)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_favorite_managers(user_id, manager_id) VALUES($1, $2);`, userID, managerID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *PostgresDB) SetFavoriteTeamByID(ctx context.Context, userID int64, teamID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.QueryContext(ctx, `SELECT id FROM user_favorite_teams WHERE user_id=$1 AND team_id=$2;`, userID, teamID)
	if err == nil {
		tx.Rollback()
		return fmt.Errorf("team with ID=%d already favorite", teamID)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_favorite_teams(user_id, team_id) VALUES($1, $2);`, userID, teamID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *PostgresDB) SetFavoriteTournamentByID(ctx context.Context, userID int64, tournamentID int64) error {
	tx, err := p.conn.Begin()
	if err != nil {
		return err
	}

	_, err = tx.QueryContext(ctx, `SELECT id FROM user_favorite_tournaments WHERE user_id=$1 AND tournament_id=$2;`, userID, tournamentID)
	if err == nil {
		tx.Rollback()
		return fmt.Errorf("tournament with ID=%d already favorite", tournamentID)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO user_favorite_tournaments(user_id, tournament_id) VALUES($1, $2);`, userID, tournamentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
