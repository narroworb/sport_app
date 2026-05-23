// Tournament page functionality
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
        console.error(`Error loading player ${athleteId} details:`, error);
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
        console.log('Loading players stats from:', url);
        
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Failed to load players stats: ${response.status}`);
        }
        
        const playersStats = await response.json();
        console.log('Players stats loaded:', playersStats.length);
        
        if (!playersStats || !Array.isArray(playersStats)) {
            throw new Error('Invalid players stats data');
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
        console.error('Error loading players stats:', error);
        const tbody = document.getElementById('players-stats-tbody');
        if (!append && tbody) {
            tbody.innerHTML = '<tr><td colspan="10">No players stats available</td>' + '</tr>';
        }
    } finally {
        isLoadingPlayersStats = false;
    }
}

// Функция загрузки статистики команд с пагинацией
async function loadTeamsStatsWithPagination(page = 1, append = false) {
    if (!tournamentId || isLoadingTeamsStats) return;
    
    isLoadingTeamsStats = true;
    const offset = (page - 1) * teamsStatsLimit + 1;
    
    try {
        const url = `/api/tournament/${tournamentId}/stats/teams?limit=${teamsStatsLimit}&offset=${offset}`;
        console.log('Loading teams stats from:', url);
        
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Failed to load teams stats: ${response.status}`);
        }
        
        const teamsStats = await response.json();
        console.log('Teams stats loaded:', teamsStats.length);
        
        if (!teamsStats || !Array.isArray(teamsStats)) {
            throw new Error('Invalid teams stats data');
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
        console.error('Error loading teams stats:', error);
        const tbody = document.getElementById('teams-stats-tbody');
        if (!append && tbody) {
            tbody.innerHTML = '<tr><td colspan="8">No teams stats available</td>' + '</tr>';
        }
    } finally {
        isLoadingTeamsStats = false;
    }
}

// Функция отрисовки статистики команд с сортировкой (расширенная)
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
            case 'possession':
                aVal = a.average_ball_possession || 0;
                bVal = b.average_ball_possession || 0;
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
        const teamName = teamDetail?.name || stat.team?.name || `Team ${stat.team_id}`;
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
    
    // Добавляем кнопку "Load More" если есть еще данные
    if (teamsStatsCurrentPage * teamsStatsLimit < teamsStatsTotalCount && !append) {
        html += `
            <tr>
                <td colspan="11" style="text-align: center;">
                    <button class="btn-secondary" onclick="loadMoreTeamsStats()">Load More ↓</button>
                </td>
            </tr>
        `;
    }
    
    tbody.innerHTML = html;
    updateTeamsSortIcons();
}

// Функция отрисовки статистики игроков с сортировкой (полная версия)
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
                    : `Player ${a.athlete_id}`;
                const playerB = b.player_details?.first_name && b.player_details?.last_name 
                    ? `${b.player_details.first_name} ${b.player_details.last_name}`
                    : `Player ${b.athlete_id}`;
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
            case 'shots':
                aVal = a.total_shots || 0;
                bVal = b.total_shots || 0;
                break;
            case 'shots_on_target':
                aVal = a.shots_on_target || 0;
                bVal = b.shots_on_target || 0;
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
            : `Player ${stat.athlete_id}`;
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
              </tr>
        `;
    }).join('');
    
    // Добавляем кнопку "Load More" если есть еще данные
    if (playersStatsCurrentPage * playersStatsLimit < playersStatsTotalCount && !append) {
        html += `
            <tr>
                <td colspan="19" style="text-align: center;">
                    <button class="btn-secondary" onclick="loadMorePlayersStats()">Load More ↓</button>
                </td>
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
        console.log('Loading tournament data for ID:', tournamentId);
        
        const [details, table, fixtures, positionHistory] = await Promise.all([
            tournamentAPI.getDetails(tournamentId).catch(e => {
                console.error('Details error:', e);
                return null;
            }),
            tournamentAPI.getTable(tournamentId).catch(e => {
                console.error('Table error:', e);
                return null;
            }),
            tournamentAPI.getFixtures(tournamentId).catch(e => {
                console.error('Fixtures error:', e);
                return null;
            }),
            tournamentAPI.getTableGraph(tournamentId).catch(e => {
                console.error('Position history error:', e);
                return null;
            }),
        ]);

        console.log('Tournament details:', details);
        console.log('Table:', table);
        console.log('Fixtures:', fixtures);
        console.log('Position history:', positionHistory);

        if (!details) {
            document.getElementById('tournament-detail').innerHTML = '<p class="loading">Tournament not found</p>';
            return;
        }

        positionHistoryData = positionHistory;

        // Update header with logo
        const tournamentName = details.name || 'Tournament';
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

        // Table (турнирная таблица)
        if (table && Array.isArray(table)) {
            const tbody = document.getElementById('table-tbody');
            tbody.innerHTML = table.map(team => {
                const teamData = team.team || team;
                const teamName = teamData.name || 'Team';
                const teamId = teamData.team_id;
                const teamLogo = teamData.url_logo || '';
                
                return `
                    <tr onclick="goToTeam(${teamId})" style="cursor:pointer;">
                        <td><strong>${team.position || team.pos || 'N/A'}</strong></td>
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

        // Fixtures
        if (fixtures && typeof fixtures === 'object') {
            const list = document.getElementById('fixtures-list');
            const allMatches = [];
            for (const round in fixtures) {
                if (Array.isArray(fixtures[round])) {
                    fixtures[round].forEach(match => {
                        allMatches.push({ ...match, round: parseInt(round) });
                    });
                }
            }
            allMatches.sort((a, b) => new Date(a.date) - new Date(b.date));
            
            list.innerHTML = allMatches.slice(0, 50).map(match => {
                const homeTeam = match.home_team?.name || 'Home';
                const awayTeam = match.away_team?.name || 'Away';
                const homeLogo = match.home_team?.url_logo || '';
                const awayLogo = match.away_team?.url_logo || '';
                const homeScore = match.home_team_score ?? '-';
                const awayScore = match.away_team_score ?? '-';
                const matchId = match.match_id;
                const date = match.date ? new Date(match.date) : new Date();
                const status = match.status || 'Scheduled';
                const round = match.round || 'N/A';
                
                let statusClass = '';
                let statusText = status;
                if (status === 'Ended') {
                    statusClass = 'status-ended';
                    statusText = '✓ Finished';
                } else if (status === 'Not started') {
                    statusClass = 'status-scheduled';
                    statusText = '⏱ Scheduled';
                }
                
                return `
                    <div class="card match-card" style="cursor:pointer; margin-bottom: 1rem;" onclick="goToMatch(${matchId})">
                        <div class="match-header" style="display: flex; justify-content: space-between;">
                            <span>Round ${round}</span>
                            <span>${date.toLocaleDateString()}</span>
                        </div>
                        <div class="match-score" style="display: flex; justify-content: space-between; align-items: center;">
                            <div class="team-container" style="display: flex; align-items: center; gap: 10px; min-width: 120px;">
                                ${homeLogo ? `<img src="${homeLogo}" alt="${homeTeam}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                                <span class="team-name">${homeTeam}</span>
                            </div>
                            <span class="score" style="font-size: 1.2rem; font-weight: bold;">${homeScore} : ${awayScore}</span>
                            <div class="team-container" style="display: flex; align-items: center; gap: 10px; min-width: 120px;">
                                <span class="team-name">${awayTeam}</span>
                                ${awayLogo ? `<img src="${awayLogo}" alt="${awayTeam}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                            </div>
                        </div>
                        <div class="match-status ${statusClass}" style="margin-top: 10px; padding-top: 10px; border-top: 1px solid #eee; text-align: center;">
                            ${statusText}
                        </div>
                    </div>
                `;
            }).join('');
        }

        // Position History Chart
        if (positionHistoryData && Array.isArray(positionHistoryData) && positionHistoryData.length > 0) {
            initPositionChart(positionHistoryData);
        } else {
            const historyDiv = document.getElementById('tournament-history');
            if (historyDiv) {
                const noDataMsg = document.createElement('div');
                noDataMsg.style.textAlign = 'center';
                noDataMsg.style.padding = '2rem';
                noDataMsg.style.color = '#666';
                noDataMsg.innerHTML = '<p>No position history data available for this tournament.</p>';
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
                        btn.textContent = isFavorite ? '★ Remove from Favorites' : '★ Add to Favorites';
                    }
                }
            } catch (e) {
                console.error('Error checking favorites:', e);
            }
        }
        
    } catch (error) {
        console.error('Error loading tournament data:', error);
        document.getElementById('tournament-detail').innerHTML = '<p class="loading">Error loading tournament data. Please try again.</p>';
    }
}

function initPositionChart(historyData) {
    console.log('initPositionChart called with historyData length:', historyData.length);
    
    const rounds = historyData.map((_, index) => `Round ${index + 1}`);
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
        teamSelector.innerHTML = '<option value="all">📊 Show All Teams</option>';
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
                            return `${teamName}: ${position === null ? 'N/A' : position + ' place'}`;
                        }
                    }
                }
            },
            scales: {
                y: {
                    reverse: true,
                    title: { display: true, text: 'Position', font: { weight: 'bold' } },
                    min: 1,
                    max: maxTeams,
                    ticks: { stepSize: 1 }
                },
                x: {
                    title: { display: true, text: 'Round', font: { weight: 'bold' } }
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
        if (btn.textContent.includes('Remove')) {
            await tournamentAPI.removeFavorite(tournamentId);
            btn.textContent = '★ Add to Favorites';
        } else {
            await tournamentAPI.addFavorite(tournamentId);
            btn.textContent = '★ Remove from Favorites';
        }
    } catch (error) {
        console.error('Error toggling favorite:', error);
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