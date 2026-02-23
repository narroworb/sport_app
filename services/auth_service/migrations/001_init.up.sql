CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO users(username, password_hash) VALUES('admin', '');

CREATE TABLE IF NOT EXISTS user_favorite_teams (
    id SERIAL PRIMARY KEY,
    team_id INT NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(team_id, user_id)
);

CREATE TABLE IF NOT EXISTS user_favorite_athletes (
    id SERIAL PRIMARY KEY,  
    athlete_id INT NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(athlete_id, user_id)
);

CREATE TABLE IF NOT EXISTS user_favorite_tournaments (
    id SERIAL PRIMARY KEY,
    tournament_id INT NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(tournaments_id, user_id)
);

CREATE TABLE IF NOT EXISTS user_favorite_managers (
    id SERIAL PRIMARY KEY,
    manager_id INT NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(manager_id, user_id)
);
