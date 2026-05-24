// Функционал страницы турнира
let tournamentId;
let positionHistoryData = null;
let positionChart = null;

// Пагинация для статистики игроков
let playersStatsCurrentPage = 1;
let playersStatsLimit = 20;
let playersStatsTotalCount = 0;
let isLoadingPlayersStats = false;

// Пагинация для статистики команд
let teamsStatsCurrentPage = 1;
let teamsStatsLimit = 20;
let teamsStatsTotalCount = 0;
let isLoadingTeamsStats = false;

// Текущие данные для сортировки
let currentTeamsStatsData = [];
let currentPlayersStatsData = [];
let currentTeamsSort = { column: 'team_name', direction: 'asc' };
let currentPlayersSort = { column: 'goals', direction: 'desc' };

// Пагинация для матчей по турам
let fixturesCurrentRound = 1;
let fixturesRoundsData = {};
let isLoadingFixtures = false;
let allRoundsList = [];

// Кэш для данных игроков
let playersDetailsCache = new Map();

function getTournamentIdFromUrl() {
    const params = new URLSearchParams(window.location.search);
    return params.get('id');
}

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();

    tournamentId = getTournamentIdFromUrl();
    
    if (!tournamentId) {
        const path = window.location.pathname;
        const pathParts = path.split('/');
        if (pathParts.length > 2 && pathParts[1] === 'tournament') {
            tournamentId = pathParts[2];
        }
    }

    if (tournamentId) {
        await loadTournamentData();
        if (TokenManager.hasToken()) {
            document.getElementById('fav-btn').style.display = 'inline-block';
        }
    }
    
    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
});

// Функция загрузки деталей игрока
async function loadPlayerDetails(athleteId) {
    if (playersDetailsCache.has(athleteId)) {
        return playersDetailsCache.get(athleteId);
    }
    
    try {
        const response = await fetch(`/api/player/${athleteId}/details`);
        if (response.ok) {
            const details = await response.json();
            playersDetailsCache.set(athleteId, details);
            return details;
        }
    } catch (error) {
        console.error(`Ошибка загрузки деталей игрока ${athleteId}:`, error);
    }
    return null;
}

// Функция загрузки статистики игроков с пагинацией
async function loadPlayersStatsWithPagination(page = 1, append = false) {
    if (!tournamentId || isLoadingPlayersStats) return;
    
    isLoadingPlayersStats = true;
    const offset = (page - 1) * playersStatsLimit + 1;
    
    try {
        const url = `/api/tournament/${tournamentId}/stats/players?limit=${playersStatsLimit}&offset=${offset}`;
        console.log('Загрузка статистики игроков из:', url);
        
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Не удалось загрузить статистику игроков: ${response.status}`);
        }
        
        const playersStats = await response.json();
        console.log('Статистика игроков загружена:', playersStats.length);
        
        if (!playersStats || !Array.isArray(playersStats)) {
            throw new Error('Некорректные данные статистики игроков');
        }
        
        if (playersStats.length < playersStatsLimit) {
            playersStatsTotalCount = offset + playersStats.length;
        } else {
            playersStatsTotalCount = offset + playersStatsLimit + 1;
        }
        
        // Загружаем детали для всех игроков
        const playersWithDetails = [];
        for (const stat of playersStats) {
            const details = await loadPlayerDetails(stat.athlete_id);
            playersWithDetails.push({
                ...stat,
                player_details: details
            });
        }
        
        if (append) {
            currentPlayersStatsData = [...currentPlayersStatsData, ...playersWithDetails];
        } else {
            currentPlayersStatsData = playersWithDetails;
        }
        
        await renderPlayersStats(currentPlayersStatsData, append);
        
    } catch (error) {
        console.error('Ошибка загрузки статистики игроков:', error);
        const tbody = document.getElementById('players-stats-tbody');
        if (!append && tbody) {
            tbody.innerHTML = '<tr><td colspan="10">Статистика игроков недоступна<\/td><\/tr>';
        }
    } finally {
        isLoadingPlayersStats = false;
    }
}

// Функция загрузки матчей турнира с пагинацией по турам
async function loadTournamentFixturesWithPagination(startRound = 1, roundsToLoad = 5) {
    if (!tournamentId || isLoadingFixtures) return;
    
    isLoadingFixtures = true;
    
    try {
        // Загружаем все матчи турнира
        const allFixtures = await tournamentAPI.getFixtures(tournamentId);
        
        if (!allFixtures || typeof allFixtures !== 'object') {
            throw new Error('Некорректные данные матчей');
        }
        
        // Собираем все матчи с информацией о туре
        const matchesByRound = {};
        let maxRound = 0;
        
        for (const round in allFixtures) {
            if (Array.isArray(allFixtures[round])) {
                const roundNum = parseInt(round);
                matchesByRound[roundNum] = allFixtures[round].map(match => ({
                    ...match,
                    round: roundNum
                }));
                maxRound = Math.max(maxRound, roundNum);
            }
        }
        
        // Получаем отсортированный список туров
        allRoundsList = Object.keys(matchesByRound).map(Number).sort((a, b) => a - b);
        fixturesRoundsData = matchesByRound;
        
        // Загружаем начальные туры
        const endRound = Math.min(startRound + roundsToLoad - 1, maxRound);
        const roundsToShow = allRoundsList.filter(r => r >= startRound && r <= endRound);
        
        renderFixturesByRounds(roundsToShow);
        
        // Добавляем кнопку для загрузки следующих туров
        if (endRound < maxRound) {
            addLoadMoreRoundsButton(endRound + 1, maxRound);
        } else {
            removeLoadMoreRoundsButton();
        }
        
    } catch (error) {
        console.error('Ошибка загрузки матчей турнира:', error);
        const list = document.getElementById('fixtures-list');
        if (list) {
            list.innerHTML = '<p>Матчи недоступны</p>';
        }
    } finally {
        isLoadingFixtures = false;
    }
}

function renderFixturesByRounds(roundsToShow) {
    const list = document.getElementById('fixtures-list');
    if (!list) return;
    
    let html = '';
    
    for (const round of roundsToShow) {
        const matches = fixturesRoundsData[round];
        if (!matches || matches.length === 0) continue;
        
        // Сортируем матчи по дате
        const sortedMatches = [...matches].sort((a, b) => new Date(a.date) - new Date(b.date));
        
        html += `
            <div class="round-group" style="margin-bottom: 2rem;">
                <div class="round-header" style="background: linear-gradient(135deg, #2c3e50, #34495e); color: white; padding: 0.75rem 1rem; border-radius: 8px 8px 0 0; margin-bottom: 0.5rem;">
                    <h3 style="margin: 0; font-size: 1.1rem;">🏆 Тур ${round}</h3>
                </div>
                <div class="round-matches">
                    ${sortedMatches.map(match => renderMatchCard(match)).join('')}
                </div>
            </div>
        `;
    }
    
    list.innerHTML = html;
}

function renderMatchCard(match) {
    const homeTeam = match.home_team?.name || 'Хозяева';
    const awayTeam = match.away_team?.name || 'Гости';
    const homeLogo = match.home_team?.url_logo || '';
    const awayLogo = match.away_team?.url_logo || '';
    const homeScore = match.home_team_score ?? '-';
    const awayScore = match.away_team_score ?? '-';
    const matchId = match.match_id;
    const date = match.date ? new Date(match.date) : new Date();
    const status = match.status || 'Not started';
    const tournament = match.tournament?.name || '';
    
    let statusClass = '';
    let statusText = '';
    if (status === 'Ended') {
        statusClass = 'status-ended';
        statusText = '✓ Завершен';
    } else if (status === 'Not started') {
        statusClass = 'status-scheduled';
        statusText = '⏱ Запланирован';
    } else if (status === 'In Progress' || status === 'Live') {
        statusClass = 'status-live';
        statusText = '🟢 В прямом эфире';
    }
    
    return `
        <div class="card match-card" style="cursor:pointer; margin-bottom: 0.75rem;" onclick="goToMatch(${matchId})">
            <div class="match-header" style="display: flex; justify-content: space-between; padding-bottom: 0.5rem; border-bottom: 1px solid #eee;">
                <span style="font-size: 0.8rem; color: #666;">${date.toLocaleDateString('ru-RU')}</span>
                <span style="font-size: 0.8rem; color: #666;">${tournament}</span>
            </div>
            <div class="match-score" style="display: flex; justify-content: space-between; align-items: center; padding: 0.75rem 0;">
                <div class="team-container" style="display: flex; align-items: center; gap: 10px; min-width: 140px;">
                    ${homeLogo ? `<img src="${homeLogo}" alt="${escapeHtml(homeTeam)}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                    <span class="team-name" style="font-weight: 600;">${escapeHtml(homeTeam)}</span>
                </div>
                <span class="score" style="font-size: 1.2rem; font-weight: bold; color: #2c3e50;">${homeScore} : ${awayScore}</span>
                <div class="team-container" style="display: flex; align-items: center; gap: 10px; min-width: 140px; justify-content: flex-end;">
                    <span class="team-name" style="font-weight: 600;">${escapeHtml(awayTeam)}</span>
                    ${awayLogo ? `<img src="${awayLogo}" alt="${escapeHtml(awayTeam)}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                </div>
            </div>
            <div class="match-status ${statusClass}" style="margin-top: 0.5rem; padding-top: 0.5rem; border-top: 1px solid #eee; text-align: center; font-size: 0.8rem;">
                ${statusText}
            </div>
        </div>
    `;
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function addLoadMoreRoundsButton(nextRound, maxRound) {
    const list = document.getElementById('fixtures-list');
    if (!list) return;
    
    // Удаляем существующую кнопку, если есть
    removeLoadMoreRoundsButton();
    
    const buttonHtml = `
        <div id="load-more-rounds-container" style="text-align: center; margin-top: 1rem;">
            <button id="load-more-rounds-btn" class="btn-primary" onclick="loadMoreRounds()" style="padding: 0.75rem 2rem;">
                📋 Показать следующие туры (${nextRound}-${Math.min(nextRound + 4, maxRound)})
            </button>
        </div>
    `;
    
    list.insertAdjacentHTML('beforeend', buttonHtml);
}

function removeLoadMoreRoundsButton() {
    const existingContainer = document.getElementById('load-more-rounds-container');
    if (existingContainer) {
        existingContainer.remove();
    }
}

async function loadMoreRounds() {
    if (isLoadingFixtures) return;
    
    // Находим следующий тур для загрузки
    const currentRoundsDisplayed = document.querySelectorAll('.round-group');
    let lastRoundDisplayed = 0;
    currentRoundsDisplayed.forEach(group => {
        const header = group.querySelector('.round-header h3');
        if (header) {
            const match = header.textContent.match(/Тур (\d+)/);
            if (match) {
                lastRoundDisplayed = Math.max(lastRoundDisplayed, parseInt(match[1]));
            }
        }
    });
    
    const nextStartRound = lastRoundDisplayed + 1;
    const endRound = Math.min(nextStartRound + 4, Math.max(...allRoundsList));
    
    if (nextStartRound > endRound) {
        removeLoadMoreRoundsButton();
        return;
    }
    
    isLoadingFixtures = true;
    
    // Показываем индикатор загрузки
    const loadMoreBtn = document.getElementById('load-more-rounds-btn');
    if (loadMoreBtn) {
        const originalText = loadMoreBtn.textContent;
        loadMoreBtn.textContent = '⏳ Загрузка...';
        loadMoreBtn.disabled = true;
        
        try {
            // Загружаем следующие туры
            const roundsToShow = allRoundsList.filter(r => r >= nextStartRound && r <= endRound);
            renderAdditionalRounds(roundsToShow);
            
            // Обновляем кнопку или удаляем если больше нет туров
            if (endRound >= Math.max(...allRoundsList)) {
                removeLoadMoreRoundsButton();
            } else {
                loadMoreBtn.textContent = originalText;
                loadMoreBtn.disabled = false;
            }
        } finally {
            isLoadingFixtures = false;
        }
    }
}

function renderAdditionalRounds(roundsToShow) {
    const list = document.getElementById('fixtures-list');
    if (!list) return;
    
    // Удаляем существующую кнопку загрузки
    removeLoadMoreRoundsButton();
    
    let html = '';
    
    for (const round of roundsToShow) {
        const matches = fixturesRoundsData[round];
        if (!matches || matches.length === 0) continue;
        
        const sortedMatches = [...matches].sort((a, b) => new Date(a.date) - new Date(b.date));
        
        html += `
            <div class="round-group" style="margin-bottom: 2rem;">
                <div class="round-header" style="background: linear-gradient(135deg, #2c3e50, #34495e); color: white; padding: 0.75rem 1rem; border-radius: 8px 8px 0 0; margin-bottom: 0.5rem;">
                    <h3 style="margin: 0; font-size: 1.1rem;">🏆 Тур ${round}</h3>
                </div>
                <div class="round-matches">
                    ${sortedMatches.map(match => renderMatchCard(match)).join('')}
                </div>
            </div>
        `;
    }
    
    list.insertAdjacentHTML('beforeend', html);
    
    // Добавляем кнопку для следующих туров, если нужно
    const maxRound = Math.max(...allRoundsList);
    const lastRound = roundsToShow[roundsToShow.length - 1];
    
    if (lastRound < maxRound) {
        addLoadMoreRoundsButton(lastRound + 1, maxRound);
    }
}

// Функция загрузки статистики команд с пагинацией
async function loadTeamsStatsWithPagination(page = 1, append = false) {
    if (!tournamentId || isLoadingTeamsStats) return;
    
    isLoadingTeamsStats = true;
    const offset = (page - 1) * teamsStatsLimit + 1;
    
    try {
        const url = `/api/tournament/${tournamentId}/stats/teams?limit=${teamsStatsLimit}&offset=${offset}`;
        console.log('Загрузка статистики команд из:', url);
        
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Не удалось загрузить статистику команд: ${response.status}`);
        }
        
        const teamsStats = await response.json();
        console.log('Статистика команд загружена:', teamsStats.length);
        
        if (!teamsStats || !Array.isArray(teamsStats)) {
            throw new Error('Некорректные данные статистики команд');
        }
        
        if (teamsStats.length < teamsStatsLimit) {
            teamsStatsTotalCount = offset + teamsStats.length;
        } else {
            teamsStatsTotalCount = offset + teamsStatsLimit + 1;
        }
        
        if (append) {
            currentTeamsStatsData = [...currentTeamsStatsData, ...teamsStats];
        } else {
            currentTeamsStatsData = teamsStats;
        }
        
        await renderTeamsStats(currentTeamsStatsData, append);
        
    } catch (error) {
        console.error('Ошибка загрузки статистики команд:', error);
        const tbody = document.getElementById('teams-stats-tbody');
        if (!append && tbody) {
            tbody.innerHTML = '<tr><td colspan="8">Статистика команд недоступна<\/td><\/tr>';
        }
    } finally {
        isLoadingTeamsStats = false;
    }
}

// Функция отрисовки статистики команд с сортировкой
async function renderTeamsStats(teamsStats, append = false, sortColumn = null, sortDirection = null) {
    const tbody = document.getElementById('teams-stats-tbody');
    if (!tbody) return;
    
    // Обновляем состояние сортировки
    if (sortColumn) {
        if (currentTeamsSort.column === sortColumn) {
            currentTeamsSort.direction = currentTeamsSort.direction === 'asc' ? 'desc' : 'asc';
        } else {
            currentTeamsSort.column = sortColumn;
            currentTeamsSort.direction = 'desc';
        }
    }
    
    // Сортировка
    const sortedStats = [...teamsStats].sort((a, b) => {
        let aVal, bVal;
        
        switch (currentTeamsSort.column) {
            case 'team_name':
                const teamA = a.team?.name || a.team_name || '';
                const teamB = b.team?.name || b.team_name || '';
                return currentTeamsSort.direction === 'asc' 
                    ? teamA.localeCompare(teamB) 
                    : teamB.localeCompare(teamA);
            case 'matches':
                aVal = a.total_matches || 0;
                bVal = b.total_matches || 0;
                break;
            case 'win_rate':
                aVal = a.avg_win_rate || 0;
                bVal = b.avg_win_rate || 0;
                break;
            case 'goals':
                aVal = a.goals || 0;
                bVal = b.goals || 0;
                break;
            case 'goals_per_90':
                aVal = a.goals_per_90 || 0;
                bVal = b.goals_per_90 || 0;
                break;
            case 'goals_conceded':
                aVal = a.goals_conceded || 0;
                bVal = b.goals_conceded || 0;
                break;
            case 'goals_conceded_per_90':
                aVal = a.goals_conceded_per_90 || 0;
                bVal = b.goals_conceded_per_90 || 0;
                break;
            case 'possession':
                aVal = a.average_ball_possession || 0;
                bVal = b.average_ball_possession || 0;
                break;
            case 'total_shots':
                aVal = a.total_shots || 0;
                bVal = b.total_shots || 0;
                break;
            case 'shots_on_goal':
                aVal = a.shots_on_goal || 0;
                bVal = b.shots_on_goal || 0;
                break;
            case 'pass_accuracy':
                aVal = a.average_pass_accuracy || 0;
                bVal = b.average_pass_accuracy || 0;
                break;
            default:
                aVal = a.goals || 0;
                bVal = b.goals || 0;
        }
        
        if (currentTeamsSort.direction === 'asc') {
            return aVal > bVal ? 1 : -1;
        } else {
            return aVal < bVal ? 1 : -1;
        }
    });
    
    let html = '';
    if (!append) {
        html = '';
    } else {
        html = tbody.innerHTML;
    }
    
    // Загружаем детали команд для логотипов
    const teamDetailsMap = new Map();
    for (const stat of sortedStats) {
        const teamId = stat.team_id;
        if (!teamDetailsMap.has(teamId)) {
            const teamDetail = await teamAPI.getDetails(teamId).catch(() => null);
            if (teamDetail) {
                teamDetailsMap.set(teamId, teamDetail);
            }
        }
    }
    
    html += sortedStats.map(stat => {
        const teamDetail = teamDetailsMap.get(stat.team_id);
        const teamName = teamDetail?.name || stat.team?.name || `Команда ${stat.team_id}`;
        const teamLogo = teamDetail?.url_logo || stat.team?.url_logo || '';
        const winRate = stat.avg_win_rate ? (stat.avg_win_rate * 100).toFixed(1) : 0;
        const possession = stat.average_ball_possession ? stat.average_ball_possession.toFixed(1) : 0;
        const goalsPer90 = stat.goals_per_90 ? stat.goals_per_90.toFixed(2) : 0;
        const goalsConcededPer90 = stat.goals_conceded_per_90 ? stat.goals_conceded_per_90.toFixed(2) : 0;
        const passAccuracy = stat.average_pass_accuracy ? (stat.average_pass_accuracy * 100).toFixed(1) : 0;
        
        return `
            <tr onclick="goToTeam(${stat.team_id})" style="cursor:pointer;">
                <td style="display: flex; align-items: center; gap: 10px; min-width: 180px;">
                    ${teamLogo ? `<img src="${teamLogo}" alt="${teamName}" style="height: 25px; width: 25px; object-fit: contain;">` : ''}
                    ${teamName}
                  </td>
                  <td>${stat.total_matches || 0}</td>
                  <td>${winRate}%</td>
                  <td>${stat.goals || 0}</td>
                  <td>${goalsPer90}</td>
                  <td>${stat.goals_conceded || 0}</td>
                  <td>${goalsConcededPer90}</td>
                  <td>${possession}%</td>
                  <td>${stat.total_shots || 0}</td>
                  <td>${stat.shots_on_goal || 0}</td>
                  <td>${passAccuracy}%</td>
              </tr>
        `;
    }).join('');
    
    // Добавляем кнопку "Загрузить ещё" если есть еще данные
    if (teamsStatsCurrentPage * teamsStatsLimit < teamsStatsTotalCount && !append) {
        html += `
            <tr>
                <td colspan="11" style="text-align: center;">
                    <button class="btn-secondary" onclick="loadMoreTeamsStats()">Загрузить ещё ↓</button>
                <\/td>
            </tr>
        `;
    }
    
    tbody.innerHTML = html;
    updateTeamsSortIcons();
}

// Функция отрисовки статистики игроков с сортировкой
async function renderPlayersStats(playersStats, append = false, sortColumn = null, sortDirection = null) {
    const tbody = document.getElementById('players-stats-tbody');
    if (!tbody) return;
    
    // Обновляем состояние сортировки
    if (sortColumn) {
        if (currentPlayersSort.column === sortColumn) {
            currentPlayersSort.direction = currentPlayersSort.direction === 'asc' ? 'desc' : 'asc';
        } else {
            currentPlayersSort.column = sortColumn;
            currentPlayersSort.direction = 'desc';
        }
    }
    
    // Сортировка
    const sortedStats = [...playersStats].sort((a, b) => {
        let aVal, bVal;
        
        switch (currentPlayersSort.column) {
            case 'player_name':
                const playerA = a.player_details?.first_name && a.player_details?.last_name 
                    ? `${a.player_details.first_name} ${a.player_details.last_name}`
                    : `Игрок ${a.athlete_id}`;
                const playerB = b.player_details?.first_name && b.player_details?.last_name 
                    ? `${b.player_details.first_name} ${b.player_details.last_name}`
                    : `Игрок ${b.athlete_id}`;
                return currentPlayersSort.direction === 'asc' 
                    ? playerA.localeCompare(playerB) 
                    : playerB.localeCompare(playerA);
            case 'team_name':
                aVal = a.team_name || '';
                bVal = b.team_name || '';
                return currentPlayersSort.direction === 'asc' 
                    ? aVal.localeCompare(bVal) 
                    : bVal.localeCompare(aVal);
            case 'matches':
                aVal = a.matches_played || 0;
                bVal = b.matches_played || 0;
                break;
            case 'goals':
                aVal = a.goals || 0;
                bVal = b.goals || 0;
                break;
            case 'assists':
                aVal = a.assists || 0;
                bVal = b.assists || 0;
                break;
            case 'rating':
                aVal = a.avg_rating || 0;
                bVal = b.avg_rating || 0;
                break;
            case 'goals_per_90':
                aVal = a.goals_per_90 || 0;
                bVal = b.goals_per_90 || 0;
                break;
            case 'assists_per_90':
                aVal = a.assists_per_90 || 0;
                bVal = b.assists_per_90 || 0;
                break;
            case 'shots':
                aVal = a.total_shots || 0;
                bVal = b.total_shots || 0;
                break;
            case 'shots_on_target':
                aVal = a.shots_on_target || 0;
                bVal = b.shots_on_target || 0;
                break;
            case 'shots_on_target_per_90':
                aVal = a.shots_on_target_per_90 || 0;
                bVal = b.shots_on_target_per_90 || 0;
                break;
            case 'key_passes':
                aVal = a.key_passes || 0;
                bVal = b.key_passes || 0;
                break;
            case 'pass_attempts':
                aVal = a.pass_attempts || 0;
                bVal = b.pass_attempts || 0;
                break;
            case 'complete_passes':
                aVal = a.complete_passes || 0;
                bVal = b.complete_passes || 0;
                break;
            case 'pass_accuracy':
                const accA = a.pass_attempts > 0 ? (a.complete_passes / a.pass_attempts) * 100 : 0;
                const accB = b.pass_attempts > 0 ? (b.complete_passes / b.pass_attempts) * 100 : 0;
                aVal = accA;
                bVal = accB;
                break;
            case 'tackles':
                aVal = a.total_tackles || 0;
                bVal = b.total_tackles || 0;
                break;
            case 'interceptions':
                aVal = a.interceptions || 0;
                bVal = b.interceptions || 0;
                break;
            case 'fouls':
                aVal = a.fouls || 0;
                bVal = b.fouls || 0;
                break;
            case 'yellow_cards':
                aVal = a.yellow_cards || 0;
                bVal = b.yellow_cards || 0;
                break;
            default:
                aVal = a.goals || 0;
                bVal = b.goals || 0;
        }
        
        if (currentPlayersSort.direction === 'asc') {
            return aVal > bVal ? 1 : -1;
        } else {
            return aVal < bVal ? 1 : -1;
        }
    });
    
    let html = '';
    if (!append) {
        html = '';
    } else {
        html = tbody.innerHTML;
    }
    
    html += sortedStats.map(stat => {
        const playerName = stat.player_details?.first_name && stat.player_details?.last_name 
            ? `${stat.player_details.first_name} ${stat.player_details.last_name}`
            : `Игрок ${stat.athlete_id}`;
        const playerPhoto = stat.player_details?.url_photo || '';
        const goalsPer90 = stat.goals_per_90 ? stat.goals_per_90.toFixed(2) : '0.00';
        const assistsPer90 = stat.assists_per_90 ? stat.assists_per_90.toFixed(2) : '0.00';
        const passAccuracy = stat.pass_attempts > 0 ? ((stat.complete_passes / stat.pass_attempts) * 100).toFixed(1) : '0';
        
        return `
            <tr onclick="goToPlayer(${stat.athlete_id})" style="cursor:pointer;">
                <td style="display: flex; align-items: center; gap: 10px; min-width: 200px;">
                    ${playerPhoto ? `<img src="${playerPhoto}" alt="${playerName}" style="width: 35px; height: 35px; border-radius: 50%; object-fit: cover;">` : '<div style="width: 35px; height: 35px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center;">👤</div>'}
                    <div>
                        <strong>${playerName}</strong>
                        ${stat.team_name ? `<div style="font-size: 0.7rem; color: #666;">${stat.team_name}</div>` : ''}
                    </div>
                  </td>
                  <td>${stat.matches_played || 0}</td>
                  <td>${stat.goals || 0}</td>
                  <td>${goalsPer90}</td>
                  <td>${stat.assists || 0}</td>
                  <td>${assistsPer90}</td>
                  <td>${stat.avg_rating ? stat.avg_rating.toFixed(1) : 0}</td>
                  <td>${stat.total_shots || 0}</td>
                  <td>${stat.shots_on_target || 0}</td>
                  <td>${stat.shots_on_target_per_90 ? stat.shots_on_target_per_90.toFixed(2) : '0.00'}</td>
                  <td>${stat.key_passes || 0}</td>
                  <td>${stat.pass_attempts || 0}</td>
                  <td>${stat.complete_passes || 0}</td>
                  <td>${passAccuracy}%</td>
                  <td>${stat.total_tackles || 0}</td>
                  <td>${stat.interceptions || 0}</td>
                  <td>${stat.fouls || 0}</td>
                  <td>${stat.yellow_cards || 0}</td>
                  <td>${stat.minutes_played || 0}</td>
              </td>
        `;
    }).join('');
    
    // Добавляем кнопку "Загрузить ещё" если есть еще данные
    if (playersStatsCurrentPage * playersStatsLimit < playersStatsTotalCount && !append) {
        html += `
            <tr>
                <td colspan="19" style="text-align: center;">
                    <button class="btn-secondary" onclick="loadMorePlayersStats()">Загрузить ещё ↓</button>
                 <\/td>
            </tr>
        `;
    }
    
    tbody.innerHTML = html;
    updatePlayersSortIcons();
}

// Функция загрузки следующих страниц статистики команд
async function loadMoreTeamsStats() {
    teamsStatsCurrentPage++;
    await loadTeamsStatsWithPagination(teamsStatsCurrentPage, true);
}

// Функция загрузки следующих страниц статистики игроков
async function loadMorePlayersStats() {
    playersStatsCurrentPage++;
    await loadPlayersStatsWithPagination(playersStatsCurrentPage, true);
}

// Функция обновления иконок сортировки для команд
function updateTeamsSortIcons() {
    const headers = document.querySelectorAll('#tournament-teams-stats .sortable-table th');
    headers.forEach(header => {
        const sortColumn = header.getAttribute('data-sort');
        const iconSpan = header.querySelector('.sort-icon');
        if (iconSpan) {
            if (currentTeamsSort.column === sortColumn) {
                iconSpan.textContent = currentTeamsSort.direction === 'asc' ? '🔼' : '🔽';
            } else {
                iconSpan.textContent = '↕️';
            }
        }
    });
}

// Функция обновления иконок сортировки для игроков
function updatePlayersSortIcons() {
    const headers = document.querySelectorAll('#tournament-players-stats .sortable-table th');
    headers.forEach(header => {
        const sortColumn = header.getAttribute('data-sort');
        const iconSpan = header.querySelector('.sort-icon');
        if (iconSpan) {
            if (currentPlayersSort.column === sortColumn) {
                iconSpan.textContent = currentPlayersSort.direction === 'asc' ? '🔼' : '🔽';
            } else {
                iconSpan.textContent = '↕️';
            }
        }
    });
}

// Функция добавления обработчиков сортировки для команд
function addTeamsSortingHandlers() {
    const headers = document.querySelectorAll('#tournament-teams-stats .sortable-table th');
    headers.forEach(header => {
        header.style.cursor = 'pointer';
        header.removeEventListener('click', header._teamsSortHandler);
        const handler = () => {
            const sortColumn = header.getAttribute('data-sort');
            if (sortColumn && currentTeamsStatsData.length > 0) {
                renderTeamsStats(currentTeamsStatsData, false, sortColumn);
            }
        };
        header._teamsSortHandler = handler;
        header.addEventListener('click', handler);
    });
}

// Функция добавления обработчиков сортировки для игроков
function addPlayersSortingHandlers() {
    const headers = document.querySelectorAll('#tournament-players-stats .sortable-table th');
    headers.forEach(header => {
        header.style.cursor = 'pointer';
        header.removeEventListener('click', header._playersSortHandler);
        const handler = () => {
            const sortColumn = header.getAttribute('data-sort');
            if (sortColumn && currentPlayersStatsData.length > 0) {
                renderPlayersStats(currentPlayersStatsData, false, sortColumn);
            }
        };
        header._playersSortHandler = handler;
        header.addEventListener('click', handler);
    });
}

async function loadTournamentData() {
    try {
        console.log('Загрузка данных турнира для ID:', tournamentId);
        
        const [details, table, fixtures, positionHistory] = await Promise.all([
            tournamentAPI.getDetails(tournamentId).catch(e => {
                console.error('Ошибка деталей:', e);
                return null;
            }),
            tournamentAPI.getTable(tournamentId).catch(e => {
                console.error('Ошибка таблицы:', e);
                return null;
            }),
            tournamentAPI.getFixtures(tournamentId).catch(e => {
                console.error('Ошибка матчей:', e);
                return null;
            }),
            tournamentAPI.getTableGraph(tournamentId).catch(e => {
                console.error('Ошибка истории позиций:', e);
                return null;
            }),
        ]);

        console.log('Детали турнира:', details);
        console.log('Таблица:', table);
        console.log('Матчи:', fixtures);
        console.log('История позиций:', positionHistory);

        if (!details) {
            document.getElementById('tournament-detail').innerHTML = '<p class="loading">Турнир не найден</p>';
            return;
        }

        positionHistoryData = positionHistory;

        // Обновляем заголовок с логотипом
        const tournamentName = details.name || 'Турнир';
        const tournamentLogo = details.url_logo || '';
        const season = details.season || '';
        const countryFlag = details.country?.url_flag || '';
        const countryName = details.country?.name || '';
        
        document.getElementById('tournament-name').innerHTML = `
            ${tournamentLogo ? `<img src="${tournamentLogo}" alt="${tournamentName}" style="height: 40px; vertical-align: middle; margin-right: 10px;">` : ''}
            ${tournamentName}
        `;
        
        document.getElementById('tournament-season').innerHTML = `
            ${countryFlag ? `<img src="${countryFlag}" style="width: 20px; vertical-align: middle; margin-right: 5px;">` : '🏆'} 
            ${season} ${countryName ? `(${countryName})` : ''}
        `;

        // Турнирная таблица
        if (table && Array.isArray(table)) {
            const tbody = document.getElementById('table-tbody');
            tbody.innerHTML = table.map(team => {
                const teamData = team.team || team;
                const teamName = teamData.name || 'Команда';
                const teamId = teamData.team_id;
                const teamLogo = teamData.url_logo || '';
                
                return `
                    <tr onclick="goToTeam(${teamId})" style="cursor:pointer;">
                        <td><strong>${team.position || team.pos || 'Н/Д'}</strong></td>
                        <td style="display: flex; align-items: center; gap: 10px;">
                            ${teamLogo ? `<img src="${teamLogo}" alt="${teamName}" style="height: 25px; width: 25px; object-fit: contain;">` : ''}
                            ${teamName}
                         </td>
                         <td>${team.matches_played || team.played || 0}</td>
                         <td>${team.wins || 0}</td>
                         <td>${team.draws || 0}</td>
                         <td>${team.losses || 0}</td>
                         <td>${team.goals_scored || team.gf || 0}</td>
                         <td>${team.goals_conceded || team.ga || 0}</td>
                         <td><strong>${team.points || 0}</strong></td>
                     </tr>
                `;
            }).join('');
        }

        // Загружаем статистику команд и игроков с пагинацией
        teamsStatsCurrentPage = 1;
        playersStatsCurrentPage = 1;
        
        await loadTeamsStatsWithPagination(1, false);
        await loadPlayersStatsWithPagination(1, false);
        
        setTimeout(() => {
            addTeamsSortingHandlers();
            addPlayersSortingHandlers();
        }, 200);

        // Матчи с группировкой по турам и пагинацией
        await loadTournamentFixturesWithPagination(1, 5);

        // График истории позиций
        if (positionHistoryData && Array.isArray(positionHistoryData) && positionHistoryData.length > 0) {
            initPositionChart(positionHistoryData);
        } else {
            const historyDiv = document.getElementById('tournament-history');
            if (historyDiv) {
                const noDataMsg = document.createElement('div');
                noDataMsg.style.textAlign = 'center';
                noDataMsg.style.padding = '2rem';
                noDataMsg.style.color = '#666';
                noDataMsg.innerHTML = '<p>Данные об истории позиций для этого турнира отсутствуют.</p>';
                historyDiv.appendChild(noDataMsg);
            }
        }
        
        // Проверяем избранное
        if (TokenManager.hasToken()) {
            try {
                const favorites = await tournamentAPI.getFavorites();
                if (favorites && Array.isArray(favorites)) {
                    const isFavorite = favorites.some(f => f.tournament_id === parseInt(tournamentId));
                    const btn = document.getElementById('fav-btn');
                    if (btn) {
                        btn.textContent = isFavorite ? '★ Удалить из избранного' : '★ Добавить в избранное';
                    }
                }
            } catch (e) {
                console.error('Ошибка проверки избранного:', e);
            }
        }
        
    } catch (error) {
        console.error('Ошибка загрузки данных турнира:', error);
        document.getElementById('tournament-detail').innerHTML = '<p class="loading">Ошибка загрузки данных турнира. Пожалуйста, попробуйте снова.</p>';
    }
}

function initPositionChart(historyData) {
    console.log('initPositionChart вызван с historyData длиной:', historyData.length);
    
    const rounds = historyData.map((_, index) => `Тур ${index + 1}`);
    const teamsMap = new Map();
    
    historyData.forEach((roundData, roundIndex) => {
        roundData.forEach(teamEntry => {
            const teamId = teamEntry.team.team_id;
            const teamName = teamEntry.team.name;
            const teamLogo = teamEntry.team.url_logo;
            const position = teamEntry.position;
            
            if (!teamsMap.has(teamId)) {
                teamsMap.set(teamId, {
                    id: teamId,
                    name: teamName,
                    logo: teamLogo,
                    positions: new Array(historyData.length).fill(null)
                });
            }
            const team = teamsMap.get(teamId);
            team.positions[roundIndex] = position;
        });
    });
    
    const teams = Array.from(teamsMap.values()).sort((a, b) => a.name.localeCompare(b.name));
    
    const teamSelector = document.getElementById('team-selector');
    if (teamSelector) {
        teamSelector.innerHTML = '<option value="all">📊 Показать все команды</option>';
        teams.forEach(team => {
            const option = document.createElement('option');
            option.value = team.name;
            option.textContent = team.name;
            teamSelector.appendChild(option);
        });
    }
    
    window.positionChartData = {
        rounds: rounds,
        teams: teams
    };
    
    updatePositionChart();
}

function updatePositionChart() {
    if (!window.positionChartData) return;
    
    const teamSelector = document.getElementById('team-selector');
    if (!teamSelector) return;
    
    const selectedValues = Array.from(teamSelector.selectedOptions).map(opt => opt.value);
    const showAll = selectedValues.includes('all') || selectedValues.length === 0;
    
    let teamsToShow = window.positionChartData.teams;
    if (!showAll) {
        teamsToShow = window.positionChartData.teams.filter(team => 
            selectedValues.includes(team.name)
        );
    }
    
    const datasets = teamsToShow.map(team => {
        const hue = (team.name.length * 37) % 360;
        const color = `hsl(${hue}, 70%, 55%)`;
        const data = team.positions.map(pos => pos);
        
        return {
            label: team.name,
            data: data,
            borderColor: color,
            backgroundColor: color.replace('55%', '15%'),
            borderWidth: 2,
            pointRadius: 4,
            pointHoverRadius: 6,
            fill: false,
            tension: 0.1
        };
    });
    
    if (datasets.length === 0) return;
    
    if (typeof Chart === 'undefined') {
        const script = document.createElement('script');
        script.src = 'https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js';
        script.onload = () => createChart(datasets);
        document.head.appendChild(script);
    } else {
        createChart(datasets);
    }
}

function createChart(datasets) {
    const canvas = document.getElementById('position-chart');
    if (!canvas) return;
    
    const ctx = canvas.getContext('2d');
    
    if (positionChart) positionChart.destroy();
    
    const maxTeams = window.positionChartData.teams.length;
    
    positionChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: window.positionChartData.rounds,
            datasets: datasets
        },
        options: {
            responsive: true,
            maintainAspectRatio: true,
            plugins: {
                legend: {
                    position: 'right',
                    labels: { font: { size: 10 }, boxWidth: 12 }
                },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            const teamName = context.dataset.label;
                            const position = context.raw;
                            return `${teamName}: ${position === null ? 'Н/Д' : position + ' место'}`;
                        }
                    }
                }
            },
            scales: {
                y: {
                    reverse: true,
                    title: { display: true, text: 'Позиция', font: { weight: 'bold' } },
                    min: 1,
                    max: maxTeams,
                    ticks: { stepSize: 1 }
                },
                x: {
                    title: { display: true, text: 'Тур', font: { weight: 'bold' } }
                }
            }
        }
    });
}

function resetChartSelection() {
    const teamSelector = document.getElementById('team-selector');
    if (teamSelector) {
        Array.from(teamSelector.options).forEach(opt => opt.selected = false);
        const allOption = teamSelector.querySelector('option[value="all"]');
        if (allOption) allOption.selected = true;
    }
    updatePositionChart();
}

function switchTournamentTab(tabName) {
    const tabs = document.querySelectorAll('.tab');
    const panes = document.querySelectorAll('.tab-pane');
    
    tabs.forEach(tab => tab.classList.remove('active'));
    panes.forEach(pane => pane.classList.remove('active'));
    
    if (event && event.target) {
        event.target.classList.add('active');
    }
    
    const activePane = document.getElementById(`tournament-${tabName}`);
    if (activePane) activePane.classList.add('active');
    
    if (tabName === 'history' && positionChart) {
        setTimeout(() => positionChart?.resize(), 100);
    }
}

async function toggleFavorite() {
    if (!TokenManager.hasToken()) {
        toggleAuthPanel();
        return;
    }

    try {
        const btn = document.getElementById('fav-btn');
        if (btn.textContent.includes('Удалить')) {
            await tournamentAPI.removeFavorite(tournamentId);
            btn.textContent = '★ Добавить в избранное';
        } else {
            await tournamentAPI.addFavorite(tournamentId);
            btn.textContent = '★ Удалить из избранного';
        }
    } catch (error) {
        console.error('Ошибка переключения избранного:', error);
    }
}

function goToMatch(id) {
    if (id && id > 0) window.location.href = `/match?id=${id}`;
}

function goToTeam(id) {
    if (id && id > 0) window.location.href = `/team?id=${id}`;
}

function goToPlayer(id) {
    if (id && id > 0) window.location.href = `/player?id=${id}`;
}

window.loadMoreTeamsStats = loadMoreTeamsStats;
window.loadMorePlayersStats = loadMorePlayersStats;
window.loadMoreRounds = loadMoreRounds;