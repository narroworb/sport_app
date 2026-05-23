// Match page functionality
let matchId;

// Переменные для сортировки игроков
let homePlayersData = [];
let awayPlayersData = [];
let homeGoaliesData = [];
let awayGoaliesData = [];
let currentHomePlayersSort = { column: 'rating', direction: 'desc' };
let currentAwayPlayersSort = { column: 'rating', direction: 'desc' };
let currentHomeGoaliesSort = { column: 'rating', direction: 'desc' };
let currentAwayGoaliesSort = { column: 'rating', direction: 'desc' };

function getMatchIdFromUrl() {
    const params = new URLSearchParams(window.location.search);
    return params.get('id');
}

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();

    matchId = getMatchIdFromUrl();
    
    if (!matchId) {
        const path = window.location.pathname;
        const pathParts = path.split('/');
        if (pathParts.length > 2 && pathParts[1] === 'match') {
            matchId = pathParts[2];
        }
    }

    if (matchId) {
        await loadMatchData();
    }

    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
});

// Функции сортировки для игроков
function sortHomePlayers(column) {
    if (currentHomePlayersSort.column === column) {
        currentHomePlayersSort.direction = currentHomePlayersSort.direction === 'asc' ? 'desc' : 'asc';
    } else {
        currentHomePlayersSort.column = column;
        currentHomePlayersSort.direction = 'desc';
    }
    renderHomePlayersTable();
}

function sortAwayPlayers(column) {
    if (currentAwayPlayersSort.column === column) {
        currentAwayPlayersSort.direction = currentAwayPlayersSort.direction === 'asc' ? 'desc' : 'asc';
    } else {
        currentAwayPlayersSort.column = column;
        currentAwayPlayersSort.direction = 'desc';
    }
    renderAwayPlayersTable();
}

function sortHomeGoalies(column) {
    if (currentHomeGoaliesSort.column === column) {
        currentHomeGoaliesSort.direction = currentHomeGoaliesSort.direction === 'asc' ? 'desc' : 'asc';
    } else {
        currentHomeGoaliesSort.column = column;
        currentHomeGoaliesSort.direction = 'desc';
    }
    renderHomeGoaliesTable();
}

function sortAwayGoalies(column) {
    if (currentAwayGoaliesSort.column === column) {
        currentAwayGoaliesSort.direction = currentAwayGoaliesSort.direction === 'asc' ? 'desc' : 'asc';
    } else {
        currentAwayGoaliesSort.column = column;
        currentAwayGoaliesSort.direction = 'desc';
    }
    renderAwayGoaliesTable();
}

// Функции рендеринга с сортировкой
function renderHomePlayersTable() {
    const sortedData = [...homePlayersData].sort((a, b) => {
        let aVal, bVal;
        switch (currentHomePlayersSort.column) {
            case 'name': aVal = a.name.toLowerCase(); bVal = b.name.toLowerCase(); break;
            case 'goals': aVal = a.goals || 0; bVal = b.goals || 0; break;
            case 'assists': aVal = a.assists || 0; bVal = b.assists || 0; break;
            case 'shots': aVal = a.shots || 0; bVal = b.shots || 0; break;
            case 'shots_on_target': aVal = a.shots_on_target || 0; bVal = b.shots_on_target || 0; break;
            case 'key_passes': aVal = a.key_passes || 0; bVal = b.key_passes || 0; break;
            case 'tackles': aVal = a.tackles || 0; bVal = b.tackles || 0; break;
            case 'interceptions': aVal = a.interceptions || 0; bVal = b.interceptions || 0; break;
            case 'fouls': aVal = a.fouls || 0; bVal = b.fouls || 0; break;
            case 'yellow_cards': aVal = a.yellow_cards || 0; bVal = b.yellow_cards || 0; break;
            case 'rating': aVal = a.rating || 0; bVal = b.rating || 0; break;
            case 'minutes': aVal = a.minutes || 0; bVal = b.minutes || 0; break;
            default: aVal = a.rating || 0; bVal = b.rating || 0;
        }
        if (currentHomePlayersSort.direction === 'asc') return aVal > bVal ? 1 : -1;
        return aVal < bVal ? 1 : -1;
    });
    
    const tbody = document.getElementById('home-players-tbody');
    if (tbody) {
        tbody.innerHTML = sortedData.map(p => `
            <tr onclick="goToPlayer(${p.playerId})" style="cursor:pointer;">
                <td style="display: flex; align-items: center; gap: 10px;">
                    ${p.photo ? `<img src="${p.photo}" alt="${p.name}" style="width: 35px; height: 35px; border-radius: 50%; object-fit: cover;">` : 
                              `<div style="width: 35px; height: 35px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center;">👤</div>`}
                    <div>
                        <strong>${p.name}</strong>
                        <div style="font-size: 0.7rem; color: #666;">${p.position} • ${p.minutes}' ${p.isStarter}</div>
                    </div>
                    </td>
                   <td style="text-align: center;">${p.goals}</td>
                   <td style="text-align: center;">${p.assists}</td>
                   <td style="text-align: center;">${p.shots}</td>
                   <td style="text-align: center;">${p.shots_on_target}</td>
                   <td style="text-align: center;">${p.key_passes}</td>
                   <td style="text-align: center;">${p.tackles}</td>
                   <td style="text-align: center;">${p.interceptions}</td>
                   <td style="text-align: center;">${p.fouls}</td>
                   <td style="text-align: center;">${p.yellow_cards}</td>
                   <td style="text-align: center;"><strong>${p.rating}</strong></td>
                  </tr>
        `).join('');
    }
    updateHomePlayersSortIcons();
}

function renderAwayPlayersTable() {
    const sortedData = [...awayPlayersData].sort((a, b) => {
        let aVal, bVal;
        switch (currentAwayPlayersSort.column) {
            case 'name': aVal = a.name.toLowerCase(); bVal = b.name.toLowerCase(); break;
            case 'goals': aVal = a.goals || 0; bVal = b.goals || 0; break;
            case 'assists': aVal = a.assists || 0; bVal = b.assists || 0; break;
            case 'shots': aVal = a.shots || 0; bVal = b.shots || 0; break;
            case 'shots_on_target': aVal = a.shots_on_target || 0; bVal = b.shots_on_target || 0; break;
            case 'key_passes': aVal = a.key_passes || 0; bVal = b.key_passes || 0; break;
            case 'tackles': aVal = a.tackles || 0; bVal = b.tackles || 0; break;
            case 'interceptions': aVal = a.interceptions || 0; bVal = b.interceptions || 0; break;
            case 'fouls': aVal = a.fouls || 0; bVal = b.fouls || 0; break;
            case 'yellow_cards': aVal = a.yellow_cards || 0; bVal = b.yellow_cards || 0; break;
            case 'rating': aVal = a.rating || 0; bVal = b.rating || 0; break;
            case 'minutes': aVal = a.minutes || 0; bVal = b.minutes || 0; break;
            default: aVal = a.rating || 0; bVal = b.rating || 0;
        }
        if (currentAwayPlayersSort.direction === 'asc') return aVal > bVal ? 1 : -1;
        return aVal < bVal ? 1 : -1;
    });
    
    const tbody = document.getElementById('away-players-tbody');
    if (tbody) {
        tbody.innerHTML = sortedData.map(p => `
            <tr onclick="goToPlayer(${p.playerId})" style="cursor:pointer;">
                <td style="display: flex; align-items: center; gap: 10px;">
                    ${p.photo ? `<img src="${p.photo}" alt="${p.name}" style="width: 35px; height: 35px; border-radius: 50%; object-fit: cover;">` : 
                              `<div style="width: 35px; height: 35px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center;">👤</div>`}
                    <div>
                        <strong>${p.name}</strong>
                        <div style="font-size: 0.7rem; color: #666;">${p.position} • ${p.minutes}' ${p.isStarter}</div>
                    </div>
                    </td>
                   <td style="text-align: center;">${p.goals}</td>
                   <td style="text-align: center;">${p.assists}</td>
                   <td style="text-align: center;">${p.shots}</td>
                   <td style="text-align: center;">${p.shots_on_target}</td>
                   <td style="text-align: center;">${p.key_passes}</td>
                   <td style="text-align: center;">${p.tackles}</td>
                   <td style="text-align: center;">${p.interceptions}</td>
                   <td style="text-align: center;">${p.fouls}</td>
                   <td style="text-align: center;">${p.yellow_cards}</td>
                   <td style="text-align: center;"><strong>${p.rating}</strong></td>
                  </tr>
        `).join('');
    }
    updateAwayPlayersSortIcons();
}

function renderHomeGoaliesTable() {
    const sortedData = [...homeGoaliesData].sort((a, b) => {
        let aVal, bVal;
        switch (currentHomeGoaliesSort.column) {
            case 'name': aVal = a.name.toLowerCase(); bVal = b.name.toLowerCase(); break;
            case 'saves': aVal = a.saves || 0; bVal = b.saves || 0; break;
            case 'conceded': aVal = a.conceded || 0; bVal = b.conceded || 0; break;
            case 'passes': aVal = a.passes || 0; bVal = b.passes || 0; break;
            case 'passes_accuracy': aVal = a.passes_accuracy || 0; bVal = b.passes_accuracy || 0; break;
            case 'penalty_saved': aVal = a.penalty_saved || 0; bVal = b.penalty_saved || 0; break;
            case 'rating': aVal = a.rating || 0; bVal = b.rating || 0; break;
            default: aVal = a.rating || 0; bVal = b.rating || 0;
        }
        if (currentHomeGoaliesSort.direction === 'asc') return aVal > bVal ? 1 : -1;
        return aVal < bVal ? 1 : -1;
    });
    
    const tbody = document.getElementById('home-goalies-tbody');
    if (tbody) {
        tbody.innerHTML = sortedData.map(g => `
            <tr onclick="goToPlayer(${g.playerId})" style="cursor:pointer;">
                <td style="display: flex; align-items: center; gap: 10px;">
                    ${g.photo ? `<img src="${g.photo}" alt="${g.name}" style="width: 35px; height: 35px; border-radius: 50%; object-fit: cover;">` : 
                              `<div style="width: 35px; height: 35px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center;">🧤</div>`}
                    <div>
                        <strong>${g.name}</strong>
                        <div style="font-size: 0.7rem; color: #666;">Goalkeeper • ${g.minutes}' ${g.isStarter}</div>
                    </div>
                    </td>
                   <td style="text-align: center;">${g.saves}</td>
                   <td style="text-align: center;">${g.conceded}</td>
                   <td style="text-align: center;">${g.passes}</td>
                   <td style="text-align: center;">${g.passes_accuracy}%</td>
                   <td style="text-align: center;">${g.penalty_saved}</td>
                   <td style="text-align: center;">${g.penalty_conceded}</td>
                   <td style="text-align: center;"><strong>${g.rating}</strong></td>
                  </tr>
        `).join('');
    }
    updateHomeGoaliesSortIcons();
}

function renderAwayGoaliesTable() {
    const sortedData = [...awayGoaliesData].sort((a, b) => {
        let aVal, bVal;
        switch (currentAwayGoaliesSort.column) {
            case 'name': aVal = a.name.toLowerCase(); bVal = b.name.toLowerCase(); break;
            case 'saves': aVal = a.saves || 0; bVal = b.saves || 0; break;
            case 'conceded': aVal = a.conceded || 0; bVal = b.conceded || 0; break;
            case 'passes': aVal = a.passes || 0; bVal = b.passes || 0; break;
            case 'passes_accuracy': aVal = a.passes_accuracy || 0; bVal = b.passes_accuracy || 0; break;
            case 'penalty_saved': aVal = a.penalty_saved || 0; bVal = b.penalty_saved || 0; break;
            case 'rating': aVal = a.rating || 0; bVal = b.rating || 0; break;
            default: aVal = a.rating || 0; bVal = b.rating || 0;
        }
        if (currentAwayGoaliesSort.direction === 'asc') return aVal > bVal ? 1 : -1;
        return aVal < bVal ? 1 : -1;
    });
    
    const tbody = document.getElementById('away-goalies-tbody');
    if (tbody) {
        tbody.innerHTML = sortedData.map(g => `
            <tr onclick="goToPlayer(${g.playerId})" style="cursor:pointer;">
                <td style="display: flex; align-items: center; gap: 10px;">
                    ${g.photo ? `<img src="${g.photo}" alt="${g.name}" style="width: 35px; height: 35px; border-radius: 50%; object-fit: cover;">` : 
                              `<div style="width: 35px; height: 35px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center;">🧤</div>`}
                    <div>
                        <strong>${g.name}</strong>
                        <div style="font-size: 0.7rem; color: #666;">Goalkeeper • ${g.minutes}' ${g.isStarter}</div>
                    </div>
                    </td>
                   <td style="text-align: center;">${g.saves}</td>
                   <td style="text-align: center;">${g.conceded}</td>
                   <td style="text-align: center;">${g.passes}</td>
                   <td style="text-align: center;">${g.passes_accuracy}%</td>
                   <td style="text-align: center;">${g.penalty_saved}</td>
                   <td style="text-align: center;">${g.penalty_conceded}</td>
                   <td style="text-align: center;"><strong>${g.rating}</strong></td>
                  </tr>
        `).join('');
    }
    updateAwayGoaliesSortIcons();
}

function updateHomePlayersSortIcons() {
    const headers = document.querySelectorAll('#home-players-table .sortable-th');
    headers.forEach(header => {
        const column = header.getAttribute('data-sort');
        const icon = header.querySelector('.sort-icon');
        if (icon) {
            if (currentHomePlayersSort.column === column) {
                icon.textContent = currentHomePlayersSort.direction === 'asc' ? '🔼' : '🔽';
            } else {
                icon.textContent = '↕️';
            }
        }
    });
}

function updateAwayPlayersSortIcons() {
    const headers = document.querySelectorAll('#away-players-table .sortable-th');
    headers.forEach(header => {
        const column = header.getAttribute('data-sort');
        const icon = header.querySelector('.sort-icon');
        if (icon) {
            if (currentAwayPlayersSort.column === column) {
                icon.textContent = currentAwayPlayersSort.direction === 'asc' ? '🔼' : '🔽';
            } else {
                icon.textContent = '↕️';
            }
        }
    });
}

function updateHomeGoaliesSortIcons() {
    const headers = document.querySelectorAll('#home-goalies-table .sortable-th');
    headers.forEach(header => {
        const column = header.getAttribute('data-sort');
        const icon = header.querySelector('.sort-icon');
        if (icon) {
            if (currentHomeGoaliesSort.column === column) {
                icon.textContent = currentHomeGoaliesSort.direction === 'asc' ? '🔼' : '🔽';
            } else {
                icon.textContent = '↕️';
            }
        }
    });
}

function updateAwayGoaliesSortIcons() {
    const headers = document.querySelectorAll('#away-goalies-table .sortable-th');
    headers.forEach(header => {
        const column = header.getAttribute('data-sort');
        const icon = header.querySelector('.sort-icon');
        if (icon) {
            if (currentAwayGoaliesSort.column === column) {
                icon.textContent = currentAwayGoaliesSort.direction === 'asc' ? '🔼' : '🔽';
            } else {
                icon.textContent = '↕️';
            }
        }
    });
}

// Функция для определения статуса матча
function getMatchStatusInfo(details) {
    const matchDate = details.date ? new Date(details.date) : new Date();
    const now = new Date();
    const status = details.status || 'Scheduled';
    
    if (status === 'Ended') {
        return { isFinished: true, isLive: false, isPending: false, message: null };
    }
    
    if (status === 'Live' || status === 'In Progress' || status === 'First Half' || status === 'Second Half') {
        return { isFinished: false, isLive: true, isPending: false, message: '⚽ Match is currently in progress! Live updates coming soon.' };
    }
    
    if (matchDate > now) {
        const timeUntil = Math.ceil((matchDate - now) / (1000 * 60));
        return { 
            isFinished: false, isLive: false, isPending: false, 
            message: `⏰ Match starts in approximately ${timeUntil} minutes` 
        };
    }
    
    if (matchDate < now && status !== 'Ended') {
        return { 
            isFinished: false, isLive: false, isPending: true, 
            message: '⏳ The match has ended. Statistics will be updated shortly. Please check back in a few minutes.' 
        };
    }
    
    return { isFinished: false, isLive: false, isPending: false, message: null };
}

// Функция загрузки деталей турнира
async function loadTournamentDetails(tournamentId) {
    if (!tournamentId) return null;
    try {
        const response = await fetch(`/api/tournament/${tournamentId}/details`);
        if (response.ok) {
            return await response.json();
        }
    } catch (error) {
        console.error('Error loading tournament details:', error);
    }
    return null;
}

// Функция отрисовки статистики в Overview
function renderOverviewStats(teamStats, homeTeamName, awayTeamName, homeTeamLogo, awayTeamLogo) {
    const overviewDiv = document.getElementById('match-overview');
    if (!overviewDiv) return;
    
    if (!teamStats) return;
    
    // Очищаем overview от старых элементов
    const elementsToRemove = overviewDiv.querySelectorAll('.match-stats-container, .stats-grid, #home-team-stats, #away-team-stats');
    elementsToRemove.forEach(el => el.remove());
    
    const homePossession = teamStats.ball_possession_home_team ?? 0;
    const awayPossession = teamStats.ball_possession_away_team ?? 0;
    const homeShots = teamStats.total_shots_home_team ?? 0;
    const awayShots = teamStats.total_shots_away_team ?? 0;
    const homeShotsOnGoal = teamStats.shots_on_goal_home_team ?? 0;
    const awayShotsOnGoal = teamStats.shots_on_goal_away_team ?? 0;
    const homePassAccuracy = teamStats.total_passes_home_team > 0 ? ((teamStats.complete_passes_home_team / teamStats.total_passes_home_team) * 100).toFixed(1) : 0;
    const awayPassAccuracy = teamStats.total_passes_away_team > 0 ? ((teamStats.complete_passes_away_team / teamStats.total_passes_away_team) * 100).toFixed(1) : 0;
    const homeCorners = teamStats.corner_kicks_home_team ?? 0;
    const awayCorners = teamStats.corner_kicks_away_team ?? 0;
    const homeFouls = teamStats.fouls_home_team ?? 0;
    const awayFouls = teamStats.fouls_away_team ?? 0;
    const homeYellowCards = teamStats.yellow_cards_home_team ?? 0;
    const awayYellowCards = teamStats.yellow_cards_away_team ?? 0;
    const homeRedCards = teamStats.red_cards_home_team ?? 0;
    const awayRedCards = teamStats.red_cards_away_team ?? 0;
    
    const statsHtml = `
        <div class="match-stats-container" style="margin: 1.5rem 0;">
            <div style="display: flex; justify-content: space-between; align-items: center; gap: 1rem; margin-bottom: 2rem;">
                <div style="flex: 1; text-align: center;">
                    <div style="font-size: 2.5rem; font-weight: bold; color: #2c3e50;">${homeShots}</div>
                    <div style="color: #666; font-size: 0.9rem;">Shots</div>
                    <div style="font-size: 2rem; font-weight: bold; color: #2c3e50; margin-top: 1rem;">${homePossession}%</div>
                    <div style="color: #666; font-size: 0.9rem;">Possession</div>
                </div>
                <div style="font-size: 1.5rem; font-weight: bold; color: #3498db;">VS</div>
                <div style="flex: 1; text-align: center;">
                    <div style="font-size: 2.5rem; font-weight: bold; color: #2c3e50;">${awayShots}</div>
                    <div style="color: #666; font-size: 0.9rem;">Shots</div>
                    <div style="font-size: 2rem; font-weight: bold; color: #2c3e50; margin-top: 1rem;">${awayPossession}%</div>
                    <div style="color: #666; font-size: 0.9rem;">Possession</div>
                </div>
            </div>
            
            <div style="background: #f8f9fa; border-radius: 12px; padding: 1.5rem;">
                <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 2rem;">
                    <div>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 0.75rem;">
                            <span>Shots on Goal</span>
                            <span style="font-weight: bold;">${homeShotsOnGoal}</span>
                        </div>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 0.75rem;">
                            <span>Pass Accuracy</span>
                            <span style="font-weight: bold;">${homePassAccuracy}%</span>
                        </div>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 0.75rem;">
                            <span>Corners</span>
                            <span style="font-weight: bold;">${homeCorners}</span>
                        </div>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 0.75rem;">
                            <span>Fouls</span>
                            <span style="font-weight: bold;">${homeFouls}</span>
                        </div>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 0.75rem;">
                            <span>Yellow Cards</span>
                            <span style="font-weight: bold;">${homeYellowCards}</span>
                        </div>
                        <div style="display: flex; justify-content: space-between;">
                            <span>Red Cards</span>
                            <span style="font-weight: bold;">${homeRedCards}</span>
                        </div>
                    </div>
                    <div>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 0.75rem;">
                            <span>Shots on Goal</span>
                            <span style="font-weight: bold;">${awayShotsOnGoal}</span>
                        </div>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 0.75rem;">
                            <span>Pass Accuracy</span>
                            <span style="font-weight: bold;">${awayPassAccuracy}%</span>
                        </div>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 0.75rem;">
                            <span>Corners</span>
                            <span style="font-weight: bold;">${awayCorners}</span>
                        </div>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 0.75rem;">
                            <span>Fouls</span>
                            <span style="font-weight: bold;">${awayFouls}</span>
                        </div>
                        <div style="display: flex; justify-content: space-between; margin-bottom: 0.75rem;">
                            <span>Yellow Cards</span>
                            <span style="font-weight: bold;">${awayYellowCards}</span>
                        </div>
                        <div style="display: flex; justify-content: space-between;">
                            <span>Red Cards</span>
                            <span style="font-weight: bold;">${awayRedCards}</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `;
    
    const infoMessage = overviewDiv.querySelector('.info-message');
    if (infoMessage) {
        infoMessage.insertAdjacentHTML('afterend', statsHtml);
    } else {
        overviewDiv.insertAdjacentHTML('afterbegin', statsHtml);
    }
}

async function loadMatchData() {
    try {
        console.log('Loading match data for ID:', matchId);
        
        const details = await fixtureAPI.getDetails(matchId).catch(e => {
            console.error('Details error:', e);
            return null;
        });
        
        const playersStats = await fixtureAPI.getPlayersStats(matchId).catch(e => {
            console.error('Players stats error:', e);
            return null;
        });
        
        const goaliesStats = await fixtureAPI.getGoaliesStats(matchId).catch(e => {
            console.error('Goalies stats error:', e);
            return null;
        });
        
        const teamStats = await fixtureAPI.getTeamsStats(matchId).catch(e => {
            console.error('Team stats error:', e);
            return null;
        });

        console.log('Match details:', details);
        console.log('Players stats:', playersStats);
        console.log('Goalies stats:', goaliesStats);
        console.log('Team stats:', teamStats);

        if (!details) {
            const detailDiv = document.getElementById('match-detail');
            if (detailDiv) detailDiv.innerHTML = '<p class="loading">Match not found</p>';
            return;
        }

        const matchStatusInfo = getMatchStatusInfo(details);
        const isFinished = matchStatusInfo.isFinished;
        const isLive = matchStatusInfo.isLive;
        const isPending = matchStatusInfo.isPending;
        
        console.log('Match status info:', matchStatusInfo);

        const homeTeamId = details.home_team?.team_id || details.home_team_id;
        const awayTeamId = details.away_team?.team_id || details.away_team_id;
        
        let homeTeamDetails = null;
        let awayTeamDetails = null;
        let tournamentDetails = null;
        
        if (homeTeamId) {
            homeTeamDetails = await teamAPI.getDetails(homeTeamId).catch(e => {
                console.error('Home team details error:', e);
                return null;
            });
        }
        
        if (awayTeamId) {
            awayTeamDetails = await teamAPI.getDetails(awayTeamId).catch(e => {
                console.error('Away team details error:', e);
                return null;
            });
        }
        
        if (details.tournament?.tournament_id) {
            tournamentDetails = await loadTournamentDetails(details.tournament.tournament_id);
        }
        
        let matchProbabilities = null;
        let tournamentStandings = null;
        let tournamentId = null;
        
        if (!isFinished && !isLive && !isPending && details.tournament?.tournament_id) {
            tournamentId = details.tournament.tournament_id;
            const season = details.tournament?.season || '2024/2025';
            
            try {
                const probUrl = `/api/analytics/match_win_probabilities?home_team_id=${homeTeamId}&away_team_id=${awayTeamId}&season=${encodeURIComponent(season)}`;
                const probResponse = await fetch(probUrl);
                if (probResponse.ok) {
                    matchProbabilities = await probResponse.json();
                    console.log('Match probabilities:', matchProbabilities);
                }
            } catch (e) {
                console.error('Error loading probabilities:', e);
            }
            
            try {
                const standingsData = await tournamentAPI.getTable(tournamentId);
                if (standingsData && Array.isArray(standingsData)) {
                    tournamentStandings = standingsData;
                    console.log('Tournament standings:', tournamentStandings);
                }
            } catch (e) {
                console.error('Error loading standings:', e);
            }
        }
        
        // Загружаем детали всех игроков
        const allPlayers = [];
        if (playersStats) {
            if (playersStats.home_team) allPlayers.push(...playersStats.home_team);
            if (playersStats.away_team) allPlayers.push(...playersStats.away_team);
        }
        if (goaliesStats) {
            if (goaliesStats.home_team) allPlayers.push(...goaliesStats.home_team);
            if (goaliesStats.away_team) allPlayers.push(...goaliesStats.away_team);
        }
        
        const playerIds = [...new Set(allPlayers.map(p => p.athlete?.athlete_id).filter(id => id && id > 0))];
        console.log('Player IDs to fetch:', playerIds);
        
        const playerDetailsMap = new Map();
        for (const playerId of playerIds) {
            const playerDetails = await playerAPI.getDetails(playerId).catch(e => {
                console.error(`Player ${playerId} details error:`, e);
                return null;
            });
            if (playerDetails) {
                playerDetailsMap.set(playerId, playerDetails);
            }
        }

        // Update header
        const homeTeamName = homeTeamDetails?.name || details.home_team?.name || `Team ${homeTeamId}`;
        const awayTeamName = awayTeamDetails?.name || details.away_team?.name || `Team ${awayTeamId}`;
        const homeTeamLogo = homeTeamDetails?.url_logo || '';
        const awayTeamLogo = awayTeamDetails?.url_logo || '';
        const homeScore = details.home_team_score ?? details.home_score ?? 0;
        const awayScore = details.away_team_score ?? details.away_score ?? 0;
        const round = details.round || 'N/A';
        const tournamentName = tournamentDetails?.name || details.tournament?.name || '';
        const tournamentLogo = tournamentDetails?.url_logo || details.tournament?.url_logo || '';

        const homeTeamEl = document.getElementById('home-team');
        if (homeTeamEl) {
            homeTeamEl.innerHTML = `
                ${homeTeamLogo ? `<img src="${homeTeamLogo}" alt="${homeTeamName}" style="height: 40px; vertical-align: middle; margin-right: 10px;">` : ''}
                ${homeTeamName}
            `;
            homeTeamEl.style.cursor = 'pointer';
            homeTeamEl.dataset.id = homeTeamId;
            homeTeamEl.onclick = () => goToTeam(homeTeamId);
        }
        
        const awayTeamEl = document.getElementById('away-team');
        if (awayTeamEl) {
            awayTeamEl.innerHTML = `
                ${awayTeamLogo ? `<img src="${awayTeamLogo}" alt="${awayTeamName}" style="height: 40px; vertical-align: middle; margin-right: 10px;">` : ''}
                ${awayTeamName}
            `;
            awayTeamEl.style.cursor = 'pointer';
            awayTeamEl.dataset.id = awayTeamId;
            awayTeamEl.onclick = () => goToTeam(awayTeamId);
        }

        const matchScoreEl = document.getElementById('match-score');
        if (matchScoreEl) matchScoreEl.textContent = `${homeScore} : ${awayScore}`;
        
        const matchDate = details.date ? new Date(details.date) : new Date();
        const dateStr = matchDate.toLocaleDateString() + ' ' + matchDate.toLocaleTimeString([], {hour:'2-digit', minute:'2-digit'});
        
        const tournamentInfoHtml = `
            <div style="display: flex; align-items: center; justify-content: center; gap: 1rem; margin-top: 0.5rem;">
                ${tournamentLogo ? `<img src="${tournamentLogo}" alt="${tournamentName}" style="height: 25px;">` : '🏆'}
                <span>${tournamentName}</span>
                <span>•</span>
                <span>Round ${round}</span>
            </div>
        `;
        
        const dateElement = document.getElementById('match-date');
        if (dateElement) dateElement.innerHTML = `${dateStr}<br>${tournamentInfoHtml}`;
        
        let statusText = details.status || 'Scheduled';
        if (isLive) {
            statusText = '🔴 LIVE - Match in progress';
        } else if (isPending) {
            statusText = '⏳ Match ended - Awaiting statistics update';
        } else if (!isFinished && matchDate > new Date()) {
            statusText = '📅 Upcoming match';
        }
        const statusEl = document.getElementById('match-status');
        if (statusEl) statusEl.textContent = statusText;

        const homeTeamNameHeader = document.getElementById('home-team-name-header');
        if (homeTeamNameHeader) {
            homeTeamNameHeader.innerHTML = homeTeamLogo ? `<img src="${homeTeamLogo}" alt="${homeTeamName}" style="height: 30px; vertical-align: middle; margin-right: 10px;">${homeTeamName}` : homeTeamName;
        }
        
        const awayTeamNameHeader = document.getElementById('away-team-name-header');
        if (awayTeamNameHeader) {
            awayTeamNameHeader.innerHTML = awayTeamLogo ? `<img src="${awayTeamLogo}" alt="${awayTeamName}" style="height: 30px; vertical-align: middle; margin-right: 10px;">${awayTeamName}` : awayTeamName;
        }

        const overviewDiv = document.getElementById('match-overview');
        if (overviewDiv) {
            // Удаляем старые блоки статистики
            const oldStatsGrid = overviewDiv.querySelector('.stats-grid');
            if (oldStatsGrid) oldStatsGrid.remove();
            const oldMatchStats = overviewDiv.querySelector('.match-stats-container');
            if (oldMatchStats) oldMatchStats.remove();
            const oldHomeTeamStats = overviewDiv.querySelector('#home-team-stats');
            if (oldHomeTeamStats) oldHomeTeamStats.remove();
            const oldAwayTeamStats = overviewDiv.querySelector('#away-team-stats');
            if (oldAwayTeamStats) oldAwayTeamStats.remove();
        }
        
        // Отрисовка статистики в Overview
        if (teamStats && (isFinished || isLive) && overviewDiv) {
            renderOverviewStats(teamStats, homeTeamName, awayTeamName, homeTeamLogo, awayTeamLogo);
        }
        
        // Информационное сообщение о статусе
        let infoMessageHtml = '';
        if (isLive) {
            infoMessageHtml = `
                <div class="info-message live" style="background: #ff4757; color: white; padding: 1rem; border-radius: 8px; margin-bottom: 1rem; text-align: center; animation: pulse 2s infinite;">
                    🔴 LIVE! The match is currently in progress. Follow the action in real-time!
                </div>
            `;
        } else if (isPending) {
            infoMessageHtml = `
                <div class="info-message pending" style="background: #ffa502; color: white; padding: 1rem; border-radius: 8px; margin-bottom: 1rem; text-align: center;">
                    ⏳ The match has ended. Statistics will be updated shortly. Please check back in a few minutes.
                </div>
            `;
        } else if (!isFinished && matchDate > new Date()) {
            const timeUntil = Math.ceil((matchDate - new Date()) / (1000 * 60));
            infoMessageHtml = `
                <div class="info-message upcoming" style="background: #1e3799; color: white; padding: 1rem; border-radius: 8px; margin-bottom: 1rem; text-align: center;">
                    📅 Match starts in approximately ${timeUntil} minutes. Check back for live updates!
                </div>
            `;
        }
        
        if (infoMessageHtml && overviewDiv) {
            const existingInfo = overviewDiv.querySelector('.info-message');
            if (existingInfo) existingInfo.remove();
            overviewDiv.insertAdjacentHTML('afterbegin', infoMessageHtml);
        }

        // Предсказания для предстоящих матчей
        if (!isFinished && !isLive && !isPending && matchProbabilities && overviewDiv) {
            const homeProb = (matchProbabilities.p_home_win || 0) * 100;
            const drawProb = (matchProbabilities.p_draw || 0) * 100;
            const awayProb = (matchProbabilities.p_away_win || 0) * 100;
            
            const probabilitiesHtml = `
                <div class="probabilities-section" style="margin-top: 1rem; padding: 1.5rem; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 12px; color: white;">
                    <h3 style="text-align: center; margin-bottom: 1rem;">🎯 Match Prediction</h3>
                    <div style="display: flex; justify-content: space-around; gap: 1rem; flex-wrap: wrap; margin-bottom: 1rem;">
                        <div class="probability-card" style="text-align: center; flex: 1; min-width: 100px; background: rgba(255,255,255,0.2); border-radius: 8px; padding: 1rem;">
                            <div class="probability-value" style="font-size: 2rem; font-weight: bold;">${homeProb.toFixed(1)}%</div>
                            <div class="probability-label">🏠 Home Win</div>
                        </div>
                        <div class="probability-card" style="text-align: center; flex: 1; min-width: 100px; background: rgba(255,255,255,0.2); border-radius: 8px; padding: 1rem;">
                            <div class="probability-value" style="font-size: 2rem; font-weight: bold;">${drawProb.toFixed(1)}%</div>
                            <div class="probability-label">🤝 Draw</div>
                        </div>
                        <div class="probability-card" style="text-align: center; flex: 1; min-width: 100px; background: rgba(255,255,255,0.2); border-radius: 8px; padding: 1rem;">
                            <div class="probability-value" style="font-size: 2rem; font-weight: bold;">${awayProb.toFixed(1)}%</div>
                            <div class="probability-label">✈️ Away Win</div>
                        </div>
                    </div>
                    <div style="font-size: 0.8rem; text-align: center; margin-top: 1rem; opacity: 0.8;">
                        📊 Model: ${matchProbabilities.model || 'Poisson'} | Expected goals: 🏠 ${matchProbabilities.lambda_home?.toFixed(2) || '?'} - ${matchProbabilities.lambda_away?.toFixed(2) || '?'} ✈️
                    </div>
                </div>
            `;
            overviewDiv.insertAdjacentHTML('beforeend', probabilitiesHtml);
        }
        
        // Таблица турнира для предстоящих матчей
        if (!isFinished && !isLive && !isPending && tournamentStandings && tournamentStandings.length > 0 && overviewDiv) {
            const standingsHtml = `
                <div class="standings-preview" style="margin-top: 2rem;">
                    <h3>🏆 Tournament Standings</h3>
                    <div style="overflow-x: auto;">
                        <table class="table" style="min-width: 600px;">
                            <thead>
                                <tr>
                                    <th>Pos</th>
                                    <th>Team</th>
                                    <th>P</th>
                                    <th>W</th>
                                    <th>D</th>
                                    <th>L</th>
                                    <th>GF</th>
                                    <th>GA</th>
                                    <th>Pts</th>
                                </tr>
                            </thead>
                            <tbody id="standings-preview-tbody"></tbody>
                        </table>
                    </div>
                </div>
            `;
            overviewDiv.insertAdjacentHTML('beforeend', standingsHtml);
            
            const standingsTbody = document.getElementById('standings-preview-tbody');
            if (standingsTbody) {
                standingsTbody.innerHTML = tournamentStandings.map(team => {
                    const teamData = team.team || team;
                    const teamName = teamData.name || 'Team';
                    const teamId = teamData.team_id;
                    const teamLogo = teamData.url_logo || '';
                    const isHomeTeam = teamId === homeTeamId;
                    const isAwayTeam = teamId === awayTeamId;
                    
                    return `
                        <tr onclick="goToTeam(${teamId})" style="cursor:pointer; ${isHomeTeam || isAwayTeam ? 'background: #e3f2fd; font-weight: bold;' : ''}">
                            <td><strong>${team.position || team.pos || 'N/A'}</strong></td>
                            <td style="display: flex; align-items: center; gap: 8px;">
                                ${teamLogo ? `<img src="${teamLogo}" alt="${teamName}" style="height: 25px; width: 25px; object-fit: contain;">` : ''}
                                ${teamName} ${isHomeTeam ? '🏠' : ''} ${isAwayTeam ? '✈️' : ''}
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
        }

        // Team Stats таблица
        const tbody = document.getElementById('team-stats-tbody');
        if (teamStats && (isFinished || isLive) && tbody) {
            const statsToShow = [
                { name: 'Shots on Goal', home: teamStats.shots_on_goal_home_team ?? 0, away: teamStats.shots_on_goal_away_team ?? 0 },
                { name: 'Total Shots', home: teamStats.total_shots_home_team ?? 0, away: teamStats.total_shots_away_team ?? 0 },
                { name: 'Blocked Shots', home: teamStats.blocked_shots_home_team ?? 0, away: teamStats.blocked_shots_away_team ?? 0 },
                { name: 'Fouls', home: teamStats.fouls_home_team ?? 0, away: teamStats.fouls_away_team ?? 0 },
                { name: 'Corner Kicks', home: teamStats.corner_kicks_home_team ?? 0, away: teamStats.corner_kicks_away_team ?? 0 },
                { name: 'Yellow Cards', home: teamStats.yellow_cards_home_team ?? 0, away: teamStats.yellow_cards_away_team ?? 0 },
                { name: 'Red Cards', home: teamStats.red_cards_home_team ?? 0, away: teamStats.red_cards_away_team ?? 0 },
                { name: 'Total Passes', home: teamStats.total_passes_home_team ?? 0, away: teamStats.total_passes_away_team ?? 0 },
                { name: 'Complete Passes', home: teamStats.complete_passes_home_team ?? 0, away: teamStats.complete_passes_away_team ?? 0 },
                { name: 'Offsides', home: teamStats.offsides_home_team ?? 0, away: teamStats.offsides_away_team ?? 0 },
                { name: 'Shots Inside Box', home: teamStats.shots_inside_box_home_team ?? 0, away: teamStats.shots_inside_box_away_team ?? 0 }
            ];
            tbody.innerHTML = statsToShow.map(stat => `
                <tr>
                    <td><strong>${stat.name}</strong></td>
                    <td style="text-align: center;">${stat.home}</td>
                    <td style="text-align: center;">${stat.away}</td>
                </tr>
            `).join('');
        }

        // Players Stats - сбор данных для сортировки
        if (playersStats && (isFinished || isLive)) {
            homePlayersData = [];
            awayPlayersData = [];
            
            if (playersStats.home_team && Array.isArray(playersStats.home_team)) {
                homePlayersData = playersStats.home_team.map(p => {
                    const playerId = p.athlete?.athlete_id;
                    const playerDetails = playerDetailsMap.get(playerId);
                    const name = playerDetails?.first_name && playerDetails?.last_name 
                        ? `${playerDetails.first_name} ${playerDetails.last_name}`
                        : `${p.athlete?.first_name || ''} ${p.athlete?.last_name || ''}`.trim() || 'Unknown';
                    const photo = playerDetails?.url_photo || '';
                    const rating = p.rating ? p.rating.toFixed(1) : '0';
                    const position = playerDetails?.position || p.athlete?.position || '';
                    const positionNames = { 'G': 'Goalkeeper', 'D': 'Defender', 'M': 'Midfielder', 'F': 'Forward' };
                    const positionDisplay = positionNames[position] || position || 'N/A';
                    
                    return {
                        playerId, name, photo, position: positionDisplay,
                        minutes: p.minutes_played || 0,
                        isStarter: p.start_player ? '⭐' : '🔄',
                        goals: p.goals || 0,
                        assists: p.assists || 0,
                        shots: p.total_shots || 0,
                        shots_on_target: p.shots_on_target || 0,
                        key_passes: p.key_passes || 0,
                        tackles: p.total_tackles || 0,
                        interceptions: p.interceptions || 0,
                        fouls: p.fouls || 0,
                        yellow_cards: p.yellow_cards || 0,
                        rating: parseFloat(rating)
                    };
                });
            }
            
            if (playersStats.away_team && Array.isArray(playersStats.away_team)) {
                awayPlayersData = playersStats.away_team.map(p => {
                    const playerId = p.athlete?.athlete_id;
                    const playerDetails = playerDetailsMap.get(playerId);
                    const name = playerDetails?.first_name && playerDetails?.last_name 
                        ? `${playerDetails.first_name} ${playerDetails.last_name}`
                        : `${p.athlete?.first_name || ''} ${p.athlete?.last_name || ''}`.trim() || 'Unknown';
                    const photo = playerDetails?.url_photo || '';
                    const rating = p.rating ? p.rating.toFixed(1) : '0';
                    const position = playerDetails?.position || p.athlete?.position || '';
                    const positionNames = { 'G': 'Goalkeeper', 'D': 'Defender', 'M': 'Midfielder', 'F': 'Forward' };
                    const positionDisplay = positionNames[position] || position || 'N/A';
                    
                    return {
                        playerId, name, photo, position: positionDisplay,
                        minutes: p.minutes_played || 0,
                        isStarter: p.start_player ? '⭐' : '🔄',
                        goals: p.goals || 0,
                        assists: p.assists || 0,
                        shots: p.total_shots || 0,
                        shots_on_target: p.shots_on_target || 0,
                        key_passes: p.key_passes || 0,
                        tackles: p.total_tackles || 0,
                        interceptions: p.interceptions || 0,
                        fouls: p.fouls || 0,
                        yellow_cards: p.yellow_cards || 0,
                        rating: parseFloat(rating)
                    };
                });
            }
            
            renderHomePlayersTable();
            renderAwayPlayersTable();
        }

        // Goalies Stats - сбор данных для сортировки
        if (goaliesStats && (isFinished || isLive)) {
            homeGoaliesData = [];
            awayGoaliesData = [];
            
            if (goaliesStats.home_team && Array.isArray(goaliesStats.home_team)) {
                homeGoaliesData = goaliesStats.home_team.map(g => {
                    const playerId = g.athlete?.athlete_id;
                    const playerDetails = playerDetailsMap.get(playerId);
                    const name = playerDetails?.first_name && playerDetails?.last_name 
                        ? `${playerDetails.first_name} ${playerDetails.last_name}`
                        : `${g.athlete?.first_name || ''} ${g.athlete?.last_name || ''}`.trim() || 'Unknown';
                    const photo = playerDetails?.url_photo || '';
                    const rating = g.rating ? g.rating.toFixed(1) : '0';
                    const passesAccuracy = g.pass_attempts > 0 ? ((g.complete_passes / g.pass_attempts) * 100).toFixed(1) : 0;
                    
                    return {
                        playerId, name, photo,
                        minutes: g.minutes_played || 0,
                        isStarter: g.start_player ? '⭐' : '🔄',
                        saves: g.saves || 0,
                        conceded: g.goals_conceded || 0,
                        passes: g.pass_attempts || 0,
                        passes_accuracy: parseFloat(passesAccuracy),
                        penalty_saved: g.penalty_saved || 0,
                        penalty_conceded: g.penalty_conceded || 0,
                        rating: parseFloat(rating)
                    };
                });
            }
            
            if (goaliesStats.away_team && Array.isArray(goaliesStats.away_team)) {
                awayGoaliesData = goaliesStats.away_team.map(g => {
                    const playerId = g.athlete?.athlete_id;
                    const playerDetails = playerDetailsMap.get(playerId);
                    const name = playerDetails?.first_name && playerDetails?.last_name 
                        ? `${playerDetails.first_name} ${playerDetails.last_name}`
                        : `${g.athlete?.first_name || ''} ${g.athlete?.last_name || ''}`.trim() || 'Unknown';
                    const photo = playerDetails?.url_photo || '';
                    const rating = g.rating ? g.rating.toFixed(1) : '0';
                    const passesAccuracy = g.pass_attempts > 0 ? ((g.complete_passes / g.pass_attempts) * 100).toFixed(1) : 0;
                    
                    return {
                        playerId, name, photo,
                        minutes: g.minutes_played || 0,
                        isStarter: g.start_player ? '⭐' : '🔄',
                        saves: g.saves || 0,
                        conceded: g.goals_conceded || 0,
                        passes: g.pass_attempts || 0,
                        passes_accuracy: parseFloat(passesAccuracy),
                        penalty_saved: g.penalty_saved || 0,
                        penalty_conceded: g.penalty_conceded || 0,
                        rating: parseFloat(rating)
                    };
                });
            }
            
            renderHomeGoaliesTable();
            renderAwayGoaliesTable();
        } else if (isPending) {
            const homePlayersTbody = document.getElementById('home-players-tbody');
            const awayPlayersTbody = document.getElementById('away-players-tbody');
            const homeGoaliesTbody = document.getElementById('home-goalies-tbody');
            const awayGoaliesTbody = document.getElementById('away-goalies-tbody');
            
            if (homePlayersTbody) homePlayersTbody.innerHTML = '<tr><td colspan="12" style="text-align: center;">⏳ Player statistics are being processed...<\/td><\/tr>';
            if (awayPlayersTbody) awayPlayersTbody.innerHTML = '<tr><td colspan="12" style="text-align: center;">⏳ Player statistics are being processed...<\/td><\/tr>';
            if (homeGoaliesTbody) homeGoaliesTbody.innerHTML = '<tr><td colspan="8" style="text-align: center;">⏳ Goalkeeper statistics are being processed...<\/td><\/tr>';
            if (awayGoaliesTbody) awayGoaliesTbody.innerHTML = '<tr><td colspan="8" style="text-align: center;">⏳ Goalkeeper statistics are being processed...<\/td><\/tr>';
        } else if (!isFinished && !isLive) {
            const homePlayersTbody = document.getElementById('home-players-tbody');
            const awayPlayersTbody = document.getElementById('away-players-tbody');
            const homeGoaliesTbody = document.getElementById('home-goalies-tbody');
            const awayGoaliesTbody = document.getElementById('away-goalies-tbody');
            
            if (homePlayersTbody) homePlayersTbody.innerHTML = '<tr><td colspan="12" style="text-align: center;">📊 Player stats will be available after the match<\/td><\/tr>';
            if (awayPlayersTbody) awayPlayersTbody.innerHTML = '<tr><td colspan="12" style="text-align: center;">📊 Player stats will be available after the match<\/td><\/tr>';
            if (homeGoaliesTbody) homeGoaliesTbody.innerHTML = '<tr><td colspan="8" style="text-align: center;">📊 Goalkeeper stats will be available after the match<\/td><\/tr>';
            if (awayGoaliesTbody) awayGoaliesTbody.innerHTML = '<tr><td colspan="8" style="text-align: center;">📊 Goalkeeper stats will be available after the match<\/td><\/tr>';
        }
        
        // Обновляем иконки сортировки
        updateHomePlayersSortIcons();
        updateAwayPlayersSortIcons();
        updateHomeGoaliesSortIcons();
        updateAwayGoaliesSortIcons();
        
    } catch (error) {
        console.error('Error loading match data:', error);
        const detailDiv = document.getElementById('match-detail');
        if (detailDiv) detailDiv.innerHTML = '<p class="loading">Error loading match data. Please try again.</p>';
    }
}

function switchMatchTab(tabName) {
    const tabs = document.querySelectorAll('.tab');
    const panes = document.querySelectorAll('.tab-pane');
    
    tabs.forEach(tab => tab.classList.remove('active'));
    panes.forEach(pane => pane.classList.remove('active'));
    
    if (event && event.target) {
        event.target.classList.add('active');
    }
    
    const activePane = document.getElementById(`match-${tabName}`);
    if (activePane) {
        activePane.classList.add('active');
    }
}

function goToPlayer(id) {
    if (id && id > 0) {
        window.location.href = `/player?id=${id}`;
    }
}

function goToTeam(id) {
    if (id && id > 0) {
        window.location.href = `/team?id=${id}`;
    }
}

function goToMatch(id) {
    if (id && id > 0) {
        window.location.href = `/match?id=${id}`;
    }
}

// Делаем функции сортировки глобальными
window.sortHomePlayers = sortHomePlayers;
window.sortAwayPlayers = sortAwayPlayers;
window.sortHomeGoalies = sortHomeGoalies;
window.sortAwayGoalies = sortAwayGoalies;