// Функционал страницы добавления данных

// Вспомогательные функции
function getText(id) {
    return document.getElementById(id)?.value.trim() || '';
}

function getBoolean(id) {
    return document.getElementById(id)?.value === 'true';
}

function getNumber(id) {
    const val = document.getElementById(id)?.value;
    return val ? Number(val) : 0;
}

function showResult(elementId, message, isError = false) {
    const el = document.getElementById(elementId);
    if (el) {
        el.textContent = message;
        el.className = `helper-text ${isError ? 'error' : 'success'}`;
        setTimeout(() => {
            if (el) {
                el.className = 'helper-text';
                el.textContent = '';
            }
        }, 5000);
    }
}

async function apiCall(url, options) {
    const token = TokenManager.getToken();
    const headers = {
        'Content-Type': 'application/json',
    };
    if (token) {
        headers['Authorization'] = token;
    }
    
    const response = await fetch(`/api${url}`, {
        ...options,
        headers,
    });
    
    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || `Ошибка ${response.status}`);
    }
    
    return await response.json();
}

// Создание игрока
async function createPlayer() {
    const payload = {
        first_name: getText('player-first'),
        last_name: getText('player-last'),
        position: getText('player-pos'),
        date_of_birth: getText('player-dob') || null,
        height: getNumber('player-height') || null,
        preffered_foot: getText('player-foot') || null,
        current_status: getText('player-status') || null,
        url_photo: getText('player-photo') || null,
        nation: {
            name: getText('player-nation') || null,
            url_flag: getText('player-flag') || null
        }
    };
    
    try {
        const res = await apiCall('/player', { method: 'POST', body: JSON.stringify(payload) });
        showResult('player-res', `✅ Игрок создан! ID: ${res.id || res.athlete_id || 'успешно'}`, false);
        document.getElementById('player-first').value = '';
        document.getElementById('player-last').value = '';
    } catch (e) {
        showResult('player-res', `❌ Ошибка: ${e.message}`, true);
    }
}

// Создание команды
async function createTeam() {
    const payload = {
        name: getText('team-name'),
        url_logo: getText('team-logo') || null
    };
    
    try {
        const res = await apiCall('/team', { method: 'POST', body: JSON.stringify(payload) });
        showResult('team-res', `✅ Команда создана! ID: ${res.id || res.team_id || 'успешно'}`, false);
        document.getElementById('team-name').value = '';
    } catch (e) {
        showResult('team-res', `❌ Ошибка: ${e.message}`, true);
    }
}

// Создание турнира
async function createTournament() {
    const payload = {
        name: getText('t-name'),
        country: {
            name: getText('t-country') || null,
            url_flag: getText('t-flag') || null
        },
        season: getText('t-season') || null,
        url_logo: getText('t-logo') || null
    };
    
    try {
        const res = await apiCall('/tournament', { method: 'POST', body: JSON.stringify(payload) });
        showResult('t-res', `✅ Турнир создан! ID: ${res.id || res.tournament_id || 'успешно'}`, false);
        document.getElementById('t-name').value = '';
    } catch (e) {
        showResult('t-res', `❌ Ошибка: ${e.message}`, true);
    }
}

// Создание тренера
async function createManager() {
    const payload = {
        first_name: getText('m-first'),
        last_name: getText('m-last'),
        nation: {
            name: getText('m-nation') || null,
            url_flag: getText('m-flag') || null
        },
        url_photo: getText('m-photo') || null,
        total_yellow_cards: getNumber('m-yellow'),
        total_red_cards: getNumber('m-red')
    };
    
    try {
        const res = await apiCall('/manager', { method: 'POST', body: JSON.stringify(payload) });
        showResult('m-res', `✅ Тренер создан! ID: ${res.id || res.manager_id || 'успешно'}`, false);
        document.getElementById('m-first').value = '';
        document.getElementById('m-last').value = '';
    } catch (e) {
        showResult('m-res', `❌ Ошибка: ${e.message}`, true);
    }
}

// Создание матча
async function createFixture() {
    const payload = {
        date: getText('f-date'),
        home_team: { team_id: getNumber('f-home') },
        away_team: { team_id: getNumber('f-away') },
        tournament: { tournament_id: getNumber('f-tid') },
        home_team_manager: { manager_id: getNumber('f-home-manager') },
        away_team_manager: { manager_id: getNumber('f-away-manager') },
        home_team_score: getNumber('f-home-score'),
        away_team_score: getNumber('f-away-score'),
        round: getNumber('f-round'),
        status: getText('f-status') || null
    };
    
    try {
        const res = await apiCall('/fixture', { method: 'POST', body: JSON.stringify(payload) });
        showResult('f-res', `✅ Матч создан! ID: ${res.id || res.match_id || 'успешно'}`, false);
    } catch (e) {
        showResult('f-res', `❌ Ошибка: ${e.message}`, true);
    }
}

// Создание статистики матча
async function createStats() {
    const matchID = getNumber('s-match');
    if (!matchID) {
        showResult('s-res', '❌ Укажите ID матча', true);
        return;
    }
    
    const players = [];
    const goalies = [];
    const teams = {
        match_id: matchID,
        shots_on_goal_home_team: getNumber('t-shots-on-goal-home'),
        shots_on_goal_away_team: getNumber('t-shots-on-goal-away'),
        total_shots_home_team: getNumber('t-total-shots-home'),
        total_shots_away_team: getNumber('t-total-shots-away'),
        blocked_shots_home_team: getNumber('t-blocked-shots-home'),
        blocked_shots_away_team: getNumber('t-blocked-shots-away'),
        fouls_home_team: getNumber('t-fouls-home'),
        fouls_away_team: getNumber('t-fouls-away'),
        corner_kicks_home_team: getNumber('t-corner-kicks-home'),
        corner_kicks_away_team: getNumber('t-corner-kicks-away'),
        ball_possession_home_team: getNumber('t-ball-possession-home'),
        ball_possession_away_team: getNumber('t-ball-possession-away'),
        yellow_cards_home_team: getNumber('t-yellow-cards-home'),
        yellow_cards_away_team: getNumber('t-yellow-cards-away'),
        red_cards_home_team: getNumber('t-red-cards-home'),
        red_cards_away_team: getNumber('t-red-cards-away'),
        total_passes_home_team: getNumber('t-total-passes-home'),
        total_passes_away_team: getNumber('t-total-passes-away'),
        complete_passes_home_team: getNumber('t-complete-passes-home'),
        complete_passes_away_team: getNumber('t-complete-passes-away'),
        offsides_home_team: getNumber('t-offsides-home'),
        offsides_away_team: getNumber('t-offsides-away'),
        shots_inside_box_home_team: getNumber('t-shots-inside-box-home'),
        shots_inside_box_away_team: getNumber('t-shots-inside-box-away')
    };

    const playerId = getNumber('s-player');
    if (playerId) {
        players.push({
            match_id: matchID,
            athlete: { athlete_id: playerId },
            start_player: getBoolean('s-start'),
            rating: Number(getText('s-rating')) || 0,
            minutes_played: getNumber('s-min'),
            goals: getNumber('s-goals'),
            assists: getNumber('s-assists'),
            blocked_shots: getNumber('s-blocked'),
            interceptions: getNumber('s-interceptions'),
            total_tackles: getNumber('s-tackles'),
            dribbled_past: getNumber('s-dribbled'),
            duels: getNumber('s-duels'),
            duels_won: getNumber('s-duels-won'),
            fouls: getNumber('s-fouls'),
            was_fouled: getNumber('s-was-fouled'),
            pass_attempts: getNumber('s-pass-attempts'),
            complete_passes: getNumber('s-complete-passes'),
            key_passes: getNumber('s-key-passes'),
            shots_on_target: getNumber('s-shots-on-target'),
            total_shots: getNumber('s-total-shots'),
            dribble_attempts: getNumber('s-dribble-attempts'),
            complete_dribbles: getNumber('s-complete-dribbles'),
            penalty_scored: getNumber('s-penalty-scored'),
            penalty_missed: getNumber('s-penalty-missed'),
            yellow_cards: getNumber('s-yellow-cards'),
            red_cards: getNumber('s-red-cards'),
            captain: getBoolean('s-captain'),
            home_team_player: getBoolean('s-home-player')
        });
    }

    const goalieId = getNumber('g-player');
    if (goalieId) {
        goalies.push({
            match_id: matchID,
            athlete: { athlete_id: goalieId },
            start_player: getBoolean('g-start'),
            rating: Number(getText('g-rating')) || 0,
            minutes_played: getNumber('g-min'),
            goals: getNumber('g-goals'),
            assists: getNumber('g-assists'),
            goals_conceded: getNumber('g-goals-conceded'),
            saves: getNumber('g-saves'),
            pass_attempts: getNumber('g-pass-attempts'),
            complete_passes: getNumber('g-complete-passes'),
            key_passes: getNumber('g-key-passes'),
            penalty_saved: getNumber('g-penalty-saved'),
            penalty_conceded: getNumber('g-penalty-conceded'),
            fouls: getNumber('g-fouls'),
            was_fouled: getNumber('g-was-fouled'),
            yellow_cards: getNumber('g-yellow-cards'),
            red_cards: getNumber('g-red-cards'),
            captain: getBoolean('g-captain'),
            home_team_player: getBoolean('g-home-player')
        });
    }

    const payload = {};
    if (players.length > 0) payload.players = players;
    if (goalies.length > 0) payload.goalies = goalies;
    if (players.length > 0 || goalies.length > 0) payload.teams = teams;

    if (players.length === 0 && goalies.length === 0) {
        showResult('s-res', '❌ Укажите хотя бы одного игрока или вратаря', true);
        return;
    }

    try {
        const res = await apiCall('/match_stats', { method: 'POST', body: JSON.stringify(payload) });
        showResult('s-res', `✅ Статистика добавлена!`, false);
    } catch (e) {
        showResult('s-res', `❌ Ошибка: ${e.message}`, true);
    }
}

// Создание строки турнирной таблицы
async function createTable() {
    const tournamentId = getNumber('table-tournament-id');
    if (!tournamentId) {
        showResult('table-res', '❌ Укажите ID турнира', true);
        return;
    }
    
    const payload = {
        tournament_id: tournamentId,
        rows: [{
            team: { team_id: getNumber('table-team-id') },
            points: getNumber('table-points'),
            position: getNumber('table-position'),
            matches_played: getNumber('table-matches'),
            wins: getNumber('table-wins'),
            draws: getNumber('table-draws'),
            losses: getNumber('table-losses'),
            goals_scored: getNumber('table-goals-for'),
            goals_conceded: getNumber('table-goals-against')
        }]
    };
    
    try {
        const res = await apiCall('/tournament/table', { method: 'POST', body: JSON.stringify(payload) });
        showResult('table-res', `✅ Строка таблицы добавлена!`, false);
    } catch (e) {
        showResult('table-res', `❌ Ошибка: ${e.message}`, true);
    }
}