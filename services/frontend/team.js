// Team page functionality
let teamId;
let currentTeamData = {};
let currentStats = null;
let availableSeasons = [];
let currentPlayersSeason = null;

// Пагинация для fixtures
let fixturesCurrentPage = 0;
let fixturesLimit = 10;
let fixturesTotalCount = 0;
let isLoadingFixtures = false;

// Фильтры для статистики
let currentStatsFilters = {
    season: null,
    dateFrom: null,
    dateTo: null
};

// Фильтр для таблицы
let currentTableSeason = null;
let allStandings = null;

function getTeamIdFromUrl() {
    const params = new URLSearchParams(window.location.search);
    return params.get('id');
}

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();

    teamId = getTeamIdFromUrl();
    
    if (!teamId) {
        const path = window.location.pathname;
        const pathParts = path.split('/');
        if (pathParts.length > 2 && pathParts[1] === 'team') {
            teamId = pathParts[2];
        }
    }

    if (teamId) {
        await loadTeamData();
        if (TokenManager.hasToken()) {
            document.getElementById('fav-btn').style.display = 'inline-block';
        }
    }

    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
    
    // Навешиваем обработчики на фильтры
    const applyFiltersBtn = document.getElementById('apply-filters');
    const resetFiltersBtn = document.getElementById('reset-filters');
    const applyTableSeasonBtn = document.getElementById('apply-table-season');
    
    if (applyFiltersBtn) applyFiltersBtn.addEventListener('click', applyStatsFilters);
    if (resetFiltersBtn) resetFiltersBtn.addEventListener('click', resetStatsFilters);
    if (applyTableSeasonBtn) applyTableSeasonBtn.addEventListener('click', applyTableSeasonFilter);
});

// Генерация списка сезонов
function generateSeasonsList() {
    const seasons = [];
    for (let year = 2015; year <= 2025; year++) {
        seasons.push(`${year}/${year + 1}`);
    }
    return seasons;
}

// Функция загрузки fixtures с пагинацией
async function loadFixturesWithPagination(page = 1, append = false) {
    if (!teamId || isLoadingFixtures) return;
    
    isLoadingFixtures = true;
    const offset = (page - 1) * fixturesLimit + 1;
    
    try {
        // Используем teamAPI.getFixtures с параметрами
        const fixtures = await teamAPI.getFixtures(teamId, fixturesLimit, offset);
        console.log('Fixtures loaded:', fixtures);
        
        if (!fixtures || !Array.isArray(fixtures)) {
            throw new Error('Invalid fixtures data');
        }
        
        if (fixtures.length < fixturesLimit) {
            fixturesTotalCount = offset + fixtures.length;
        } else {
            fixturesTotalCount = offset + fixturesLimit + 1;
        }
        
        renderFixturesList(fixtures, append);
        
    } catch (error) {
        console.error('Error loading fixtures:', error);
        const list = document.getElementById('fixtures-list');
        if (!append && list) {
            list.innerHTML = '<p>No fixtures available</p>';
        }
    } finally {
        isLoadingFixtures = false;
    }
}

function renderFixturesList(fixtures, append = false) {
    const list = document.getElementById('fixtures-list');
    if (!list) return;
    
    let html = '';
    if (append && list.innerHTML !== '<p>No fixtures available</p>') {
        html = list.innerHTML;
    }
    
    if (!fixtures || fixtures.length === 0) {
        if (!append) {
            list.innerHTML = '<p>No fixtures available</p>';
        }
        return;
    }
    
    html += fixtures.map(f => {
        const fixtureData = f.data || f;
        const homeTeam = fixtureData.home_team?.name || fixtureData.home_team_name || 'Home';
        const awayTeam = fixtureData.away_team?.name || fixtureData.away_team_name || 'Away';
        const homeLogo = fixtureData.home_team?.url_logo || '';
        const awayLogo = fixtureData.away_team?.url_logo || '';
        const homeScore = fixtureData.home_team_score ?? fixtureData.home_score ?? '-';
        const awayScore = fixtureData.away_team_score ?? fixtureData.away_score ?? '-';
        const matchId = fixtureData.match_id || fixtureData.id;
        const date = fixtureData.date ? new Date(fixtureData.date) : new Date();
        const status = fixtureData.status || 'Not started';
        const tournament = fixtureData.tournament?.name || '';
        
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
                    <span>${date.toLocaleDateString()}</span>
                    <span>${tournament}</span>
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
    
    const currentCount = fixturesCurrentPage * fixturesLimit;
    if (currentCount + fixturesLimit < fixturesTotalCount && !append) {
        html += `
            <div style="text-align: center; margin-top: 1rem;">
                <button id="load-more-fixtures" class="btn-secondary" onclick="loadMoreFixtures()">Load More ↓</button>
            </div>
        `;
    }
    
    list.innerHTML = html;
}

async function loadMoreFixtures() {
    fixturesCurrentPage++;
    await loadFixturesWithPagination(fixturesCurrentPage + 1, true);
}

// Функция загрузки статистики с фильтрами
async function loadTeamStatsWithFilters(season, dateFrom, dateTo) {
    let url = `/api/team/${teamId}/stats`;
    const params = [];
    
    if (season && season !== '') {
        params.push(`season=${encodeURIComponent(season)}`);
    }
    if (dateFrom && dateFrom !== '') {
        params.push(`dateFrom=${dateFrom}`);
    }
    if (dateTo && dateTo !== '') {
        params.push(`dateTo=${dateTo}`);
    }
    
    if (params.length > 0) {
        url += `?${params.join('&')}`;
    }
    
    console.log('Fetching team stats from:', url);
    try {
        const response = await fetch(url);
        if (response.ok) {
            return await response.json();
        } else {
            console.error('Stats response not OK:', response.status);
            return null;
        }
    } catch (error) {
        console.error('Error fetching stats:', error);
        return null;
    }
}

// Функция загрузки статистики игроков
async function loadTeamPlayersStatsWithSeason(season) {
    let url = `/api/team/${teamId}/players_stats`;
    if (season && season !== '') {
        url += `?season=${encodeURIComponent(season)}`;
    }
    
    console.log('Fetching team players stats from:', url);
    try {
        const response = await fetch(url);
        if (response.ok) {
            return await response.json();
        } else {
            console.error('Players stats response not OK:', response.status);
            return null;
        }
    } catch (error) {
        console.error('Error fetching players stats:', error);
        return null;
    }
}

// Функция обновления статистики
async function updateTeamStats() {
    const statsGrid = document.getElementById('team-stats-grid');
    statsGrid.innerHTML = '<div class="loading-spinner">Loading statistics...</div>';
    
    const stats = await loadTeamStatsWithFilters(
        currentStatsFilters.season, 
        currentStatsFilters.dateFrom, 
        currentStatsFilters.dateTo
    );
    
    if (stats) {
        currentStats = stats;
        renderTeamStats(stats);
        showFilterMessage('Statistics updated');
    } else {
        statsGrid.innerHTML = '<p>No stats available for selected period</p>';
        showFilterMessage('No data found for selected period', 'error');
    }
}

async function updateTeamPlayersStats(season) {
    const playersStats = await loadTeamPlayersStatsWithSeason(season);
    if (playersStats && typeof playersStats === 'object') {
        originalPlayersStats = playersStats;
        renderTeamPlayersStats(playersStats);
        setTimeout(addSortingHandlers, 100);
    } else {
        document.getElementById('players-tbody').innerHTML = '<tr><td colspan="8">No players stats available</td>';
    }
}

// Функция отрисовки статистики команды
function renderTeamStats(stats) {
    const statsGrid = document.getElementById('team-stats-grid');
    const statsData = stats.data || stats;
    const statsToShow = Object.entries(statsData).slice(0, 12);
    
    if (statsToShow.length > 0) {
        statsGrid.innerHTML = statsToShow.map(([key, value]) => {
            let displayValue = value;
            if (typeof value === 'number') {
                if (Number.isInteger(value)) {
                    displayValue = value;
                } else {
                    displayValue = value.toFixed(2);
                }
            }
            return `
                <div class="stat-box">
                    <div class="stat-value">${displayValue ?? 'N/A'}</div>
                    <div class="stat-label">${formatLabel(key)}</div>
                </div>
            `;
        }).join('');
    } else {
        statsGrid.innerHTML = '<p>No stats available</p>';
    }
}

// Переменная для хранения текущего состояния сортировки
let currentSort = {
    column: 'matches',
    direction: 'desc'
};

let currentPlayersStats = [];
let originalPlayersStats = null;

// Функция отрисовки игроков с их статистикой (с сортировкой)
function renderTeamPlayersStats(playersStats, sortColumn = null, sortDirection = null) {
    const tbody = document.getElementById('players-tbody');
    if (!tbody) return;
    
    if (!playersStats) {
        tbody.innerHTML = '<tr><td colspan="8">No players stats available</td>';
        return;
    }
    
    // Обновляем состояние сортировки
    if (sortColumn) {
        if (currentSort.column === sortColumn) {
            currentSort.direction = currentSort.direction === 'asc' ? 'desc' : 'asc';
        } else {
            currentSort.column = sortColumn;
            currentSort.direction = 'desc';
        }
    }
    
    // Обновляем иконки сортировки
    updateSortIcons();
    
    const allPlayers = [];
    const positions = ['G', 'D', 'M', 'F'];
    const positionNames = { 'G': 'Goalkeeper', 'D': 'Defender', 'M': 'Midfielder', 'F': 'Forward' };
    
    positions.forEach(pos => {
        if (playersStats[pos] && Array.isArray(playersStats[pos])) {
            playersStats[pos].forEach(player => {
                const minutes = player.minutes_played || 0;
                const goals = player.goals || 0;
                const goalsPer90 = minutes > 0 ? (goals / minutes) * 90 : 0;
                
                allPlayers.push({
                    ...player,
                    position_display: positionNames[pos],
                    position_code: pos,
                    matches_played: player.matches_played || 0,
                    goals: goals,
                    assists: player.assists || 0,
                    avg_rating: player.avg_rating || 0,
                    minutes_played: minutes,
                    goals_per_90: goalsPer90,
                    full_name: `${player.first_name || ''} ${player.last_name || ''}`.trim() || 'Unknown'
                });
            });
        }
    });
    
    // Сортировка
    allPlayers.sort((a, b) => {
        let aVal, bVal;
        
        switch (currentSort.column) {
            case 'name':
                aVal = a.full_name;
                bVal = b.full_name;
                return currentSort.direction === 'asc' 
                    ? aVal.localeCompare(bVal) 
                    : bVal.localeCompare(aVal);
            case 'position':
                aVal = a.position_display;
                bVal = b.position_display;
                return currentSort.direction === 'asc' 
                    ? aVal.localeCompare(bVal) 
                    : bVal.localeCompare(aVal);
            case 'matches':
                aVal = a.matches_played;
                bVal = b.matches_played;
                break;
            case 'goals':
                aVal = a.goals;
                bVal = b.goals;
                break;
            case 'assists':
                aVal = a.assists;
                bVal = b.assists;
                break;
            case 'rating':
                aVal = a.avg_rating;
                bVal = b.avg_rating;
                break;
            case 'goals_per_90':
                aVal = a.goals_per_90;
                bVal = b.goals_per_90;
                break;
            case 'minutes':
                aVal = a.minutes_played;
                bVal = b.minutes_played;
                break;
            default:
                aVal = a.matches_played;
                bVal = b.matches_played;
        }
        
        if (typeof aVal === 'string') {
            return currentSort.direction === 'asc' 
                ? aVal.localeCompare(bVal) 
                : bVal.localeCompare(aVal);
        } else {
            if (currentSort.direction === 'asc') {
                return aVal > bVal ? 1 : -1;
            } else {
                return aVal < bVal ? 1 : -1;
            }
        }
    });
    
    currentPlayersStats = allPlayers;
    originalPlayersStats = playersStats;
    
    if (allPlayers.length > 0) {
        tbody.innerHTML = allPlayers.map(player => {
            const name = player.full_name;
            const position = player.position_display;
            const playerId = player.athlete_id || player.id;
            const photo = player.url_photo || '';
            const nation = player.nation?.name || '';
            const flag = player.nation?.url_flag || '';
            const matches = player.matches_played;
            const goals = player.goals;
            const assists = player.assists;
            const rating = player.avg_rating ? player.avg_rating.toFixed(1) : '-';
            const goalsPer90 = player.goals_per_90.toFixed(2);
            const minutes = player.minutes_played;
            
            return `
                <tr onclick="goToPlayer(${playerId})" style="cursor:pointer;">
                    <td style="display: flex; align-items: center; gap: 10px; min-width: 200px;">
                        ${photo ? `<img src="${photo}" alt="${name}" style="width: 40px; height: 40px; border-radius: 50%; object-fit: cover;">` : 
                                  `<div style="width: 40px; height: 40px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center;">📷</div>`}
                        <div>
                            <strong>${name}</strong>
                            ${nation ? `<div style="font-size: 0.7rem; color: #666;">${flag ? `<img src="${flag}" style="width: 16px; vertical-align: middle;">` : ''} ${nation}</div>` : ''}
                        </div>
                       </td>
                      <td>${position}</td>
                      <td>${matches}</td>
                      <td>${goals}</td>
                      <td>${assists}</td>
                      <td>${rating}</td>
                      <td>${goalsPer90}</td>
                      <td>${minutes}</td>
                    </tr>
            `;
        }).join('');
    } else {
        tbody.innerHTML = '<table><td colspan="8">No players stats available</td>';
    }
}

// Функция обновления иконок сортировки
function updateSortIcons() {
    const headers = document.querySelectorAll('#team-players .sortable-table th');
    headers.forEach(header => {
        const sortColumn = header.getAttribute('data-sort');
        const iconSpan = header.querySelector('.sort-icon');
        if (iconSpan) {
            if (currentSort.column === sortColumn) {
                iconSpan.textContent = currentSort.direction === 'asc' ? '🔼' : '🔽';
            } else {
                iconSpan.textContent = '↕️';
            }
        }
    });
}

// Функция добавления обработчиков сортировки
function addSortingHandlers() {
    const headers = document.querySelectorAll('#team-players .sortable-table th');
    headers.forEach(header => {
        header.style.cursor = 'pointer';
        // Удаляем старый обработчик, чтобы не дублировать
        header.removeEventListener('click', header._sortHandler);
        const handler = () => {
            const sortColumn = header.getAttribute('data-sort');
            if (sortColumn && originalPlayersStats) {
                renderTeamPlayersStats(originalPlayersStats, sortColumn);
            }
        };
        header._sortHandler = handler;
        header.addEventListener('click', handler);
    });
}

async function applyStatsFilters() {
    const seasonFilter = document.getElementById('season-filter').value;
    const dateFrom = document.getElementById('date-from').value;
    const dateTo = document.getElementById('date-to').value;
    
    if (!seasonFilter && !dateFrom && !dateTo) {
        alert('Please select either a season or date range');
        return;
    }
    
    currentStatsFilters = {
        season: seasonFilter || null,
        dateFrom: dateFrom || null,
        dateTo: dateTo || null
    };
    
    await updateTeamStats();
}

function resetStatsFilters() {
    document.getElementById('season-filter').value = '';
    document.getElementById('date-from').value = '';
    document.getElementById('date-to').value = '';
    
    currentStatsFilters = {
        season: null,
        dateFrom: null,
        dateTo: null
    };
    
    updateTeamStats();
}

async function applyTableSeasonFilter() {
    const seasonSelect = document.getElementById('table-season-filter');
    const season = seasonSelect?.value || null;
    currentTableSeason = season;
    
    const standings = await loadStandingsWithSeason(season);
    if (standings && Array.isArray(standings)) {
        renderStandings(standings);
        showFilterMessage(`Standings for season ${season || 'current'}`, 'success');
    }
}

async function loadStandingsWithSeason(season) {
    let url = `/api/team/${teamId}/standings`;
    if (season && season !== '') {
        url += `?season=${encodeURIComponent(season)}`;
    }
    
    try {
        const response = await fetch(url);
        if (response.ok) {
            return await response.json();
        }
    } catch (error) {
        console.error('Error loading standings:', error);
    }
    return null;
}

function renderStandings(standings) {
    const tbody = document.getElementById('standing-tbody');
    if (!tbody) return;
    
    tbody.innerHTML = standings.map(standing => {
        const teamData = standing.team || standing;
        const teamName = teamData.name || teamData.team_name || 'Team';
        const teamIdStanding = teamData.team_id || teamData.id;
        const teamLogo = teamData.url_logo || '';
        const isCurrentTeam = teamIdStanding == teamId;
        
        return `
            <tr onclick="goToTeam(${teamIdStanding})" style="cursor:pointer; ${isCurrentTeam ? 'background: linear-gradient(90deg, #e3f2fd, #bbdef5); font-weight: bold; border-left: 4px solid #3498db;' : ''}">
                <td><strong>${standing.position || standing.pos || 'N/A'}</strong></td>
                <td style="display: flex; align-items: center; gap: 10px;">
                    ${teamLogo ? `<img src="${teamLogo}" alt="${teamName}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                    ${teamName} ${isCurrentTeam ? '🏠' : ''}
                   </td>
                <td>${standing.points || 0}</td>
                <td>${standing.matches_played || standing.played || 0}</td>
                <td>${standing.wins || 0}/${standing.draws || 0}/${standing.losses || 0}</td>
                </tr>
        `;
    }).join('');
}

function populateSeasonSelectors() {
    const seasons = generateSeasonsList();
    
    const seasonFilter = document.getElementById('season-filter');
    if (seasonFilter) {
        seasonFilter.innerHTML = '<option value="">Select Season</option>';
        seasons.forEach(season => {
            const option = document.createElement('option');
            option.value = season;
            option.textContent = season;
            seasonFilter.appendChild(option);
        });
    }
    
    const tableSeasonFilter = document.getElementById('table-season-filter');
    if (tableSeasonFilter) {
        tableSeasonFilter.innerHTML = '<option value="">Current Season</option>';
        seasons.forEach(season => {
            const option = document.createElement('option');
            option.value = season;
            option.textContent = season;
            tableSeasonFilter.appendChild(option);
        });
    }
}

function showFilterMessage(message, type = 'success') {
    const existingMsg = document.querySelector('.filter-message');
    if (existingMsg) existingMsg.remove();
    
    const msgDiv = document.createElement('div');
    msgDiv.className = `filter-message ${type}`;
    msgDiv.style.cssText = `
        position: fixed;
        top: 80px;
        right: 20px;
        background: ${type === 'success' ? '#2ecc71' : '#e74c3c'};
        color: white;
        padding: 12px 20px;
        border-radius: 8px;
        z-index: 1000;
        animation: slideIn 0.3s ease;
        box-shadow: 0 2px 10px rgba(0,0,0,0.2);
    `;
    msgDiv.textContent = message;
    document.body.appendChild(msgDiv);
    
    setTimeout(() => {
        msgDiv.style.opacity = '0';
        msgDiv.style.transition = 'opacity 0.3s';
        setTimeout(() => msgDiv.remove(), 300);
    }, 3000);
}

async function loadTeamData() {
    try {
        console.log('Loading team data for ID:', teamId);
        
        const details = await teamAPI.getDetails(teamId).catch(e => {
            console.error('Details error:', e);
            return null;
        });
        
        const stats = await teamAPI.getStats(teamId).catch(e => {
            console.error('Stats error:', e);
            return null;
        });
        
        const nextGame = await teamAPI.getNextGame(teamId).catch(e => {
            console.error('Next game error:', e);
            return null;
        });
        
        const standings = await teamAPI.getStandings(teamId).catch(e => {
            console.error('Standings error:', e);
            return null;
        });
        
        const manager = await teamAPI.getManager(teamId).catch(e => {
            console.error('Manager error:', e);
            return null;
        });
        
        const fixtures = await teamAPI.getFixtures(teamId).catch(e => {
            console.error('Fixtures error:', e);
            return null;
        });
        
        // Заполняем селекторы сезонов
        populateSeasonSelectors();
        
        allStandings = standings;
        
        // Определяем текущий сезон
        let season = null;
        if (details && details.tournament && details.tournament.season) {
            season = details.tournament.season;
        }
        if (!season) {
            const currentYear = new Date().getFullYear();
            const currentMonth = new Date().getMonth();
            season = currentMonth >= 6 ? `${currentYear}/${currentYear + 1}` : `${currentYear - 1}/${currentYear}`;
        }
        
        // Загружаем аналитику формы команды
        let teamForm = null;
        if (teamId && season) {
            try {
                const formUrl = `/api/analytics/team_form?team_id=${teamId}&season=${encodeURIComponent(season)}&matches_back=10&half_life_matches=5`;
                const formResponse = await fetch(formUrl);
                if (formResponse.ok) {
                    teamForm = await formResponse.json();
                }
            } catch (e) {
                console.error('Error loading team form:', e);
            }
        }

        console.log('Team details:', details);
        console.log('Team stats:', stats);
        console.log('Next game:', nextGame);
        console.log('Standings:', standings);
        console.log('Manager:', manager);
        console.log('Fixtures:', fixtures);
        console.log('Team form:', teamForm);

        if (!details) {
            document.getElementById('team-detail').innerHTML = '<p class="loading">Team not found</p>';
            return;
        }

        currentTeamData = details;

        // Update header with logo
        const teamName = details.name || details.team_name || 'Team';
        document.getElementById('team-name').innerHTML = `
            ${details.url_logo ? `<img src="${details.url_logo}" alt="${teamName}" style="height: 40px; vertical-align: middle; margin-right: 10px;">` : ''}
            ${teamName}
        `;

        // Отрисовываем статистику команды
        if (stats) {
            currentStats = stats;
            renderTeamStats(stats);
        }

        // Загружаем и отрисовываем статистику игроков (с текущим сезоном)
        await updateTeamPlayersStats(season);
        
        // Добавляем селектор сезона для игроков
        addPlayersSeasonSelector();

        // Блок аналитики формы
        if (teamForm && teamForm.form_index !== undefined) {
            const formIndex = (teamForm.form_index * 100).toFixed(1);
            const attackIndex = (teamForm.attack_index * 100).toFixed(1);
            const defenseIndex = (teamForm.defense_index * 100).toFixed(1);
            const confidence = (teamForm.confidence * 100).toFixed(1);
            const trend = teamForm.trend || 0;
            
            let trendIcon = '➡️';
            let trendText = 'Stable';
            let trendColor = '#3498db';
            if (trend > 0.05) {
                trendIcon = '📈';
                trendText = 'Improving';
                trendColor = '#2ecc71';
            } else if (trend < -0.05) {
                trendIcon = '📉';
                trendText = 'Declining';
                trendColor = '#e74c3c';
            }
            
            let formRating = '';
            let formColor = '';
            if (formIndex >= 70) {
                formRating = 'Excellent 🔥';
                formColor = '#2ecc71';
            } else if (formIndex >= 55) {
                formRating = 'Good 👍';
                formColor = '#3498db';
            } else if (formIndex >= 40) {
                formRating = 'Average 📊';
                formColor = '#f39c12';
            } else {
                formRating = 'Poor 👎';
                formColor = '#e74c3c';
            }
            
            const formHtml = `
                <div id="team-form-analytics" class="form-analytics-section" style="margin-top: 2rem; padding: 1.5rem; background: linear-gradient(135deg, #1e3799 0%, #2c3e50 100%); border-radius: 12px; color: white;">
                    <h3 style="text-align: center; margin-bottom: 1rem;">📊 Team Form Analytics (${season})</h3>
                    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem;">
                        <div class="form-card" style="text-align: center; background: rgba(255,255,255,0.15); border-radius: 8px; padding: 1rem;">
                            <div class="form-value" style="font-size: 2rem; font-weight: bold; color: ${formColor};">${formIndex}%</div>
                            <div class="form-label">Form Index</div>
                            <div class="form-rating" style="font-size: 0.8rem; margin-top: 0.5rem;">${formRating}</div>
                        </div>
                        <div class="form-card" style="text-align: center; background: rgba(255,255,255,0.15); border-radius: 8px; padding: 1rem;">
                            <div class="form-value" style="font-size: 2rem; font-weight: bold;">${attackIndex}%</div>
                            <div class="form-label">⚽ Attack Index</div>
                        </div>
                        <div class="form-card" style="text-align: center; background: rgba(255,255,255,0.15); border-radius: 8px; padding: 1rem;">
                            <div class="form-value" style="font-size: 2rem; font-weight: bold;">${defenseIndex}%</div>
                            <div class="form-label">🛡️ Defense Index</div>
                        </div>
                        <div class="form-card" style="text-align: center; background: rgba(255,255,255,0.15); border-radius: 8px; padding: 1rem;">
                            <div class="form-value" style="font-size: 2rem; font-weight: bold; color: ${trendColor};">${trendIcon} ${trend > 0 ? '+' : ''}${(trend * 100).toFixed(1)}%</div>
                            <div class="form-label">Trend</div>
                            <div class="form-trend" style="font-size: 0.8rem;">${trendText}</div>
                        </div>
                    </div>
                    <div style="margin-top: 1rem; font-size: 0.8rem; text-align: center; opacity: 0.8;">
                        📊 Based on last ${teamForm.matches_used || 0} matches | Confidence: ${confidence}%
                    </div>
                    ${teamForm.details ? `<div style="margin-top: 0.5rem; font-size: 0.7rem; text-align: center; opacity: 0.6;">${teamForm.details}</div>` : ''}
                </div>
            `;
            
            const statsGrid = document.getElementById('team-stats-grid');
            if (statsGrid && !document.getElementById('team-form-analytics')) {
                statsGrid.insertAdjacentHTML('afterend', formHtml);
            }
        }

        // Next game with team logos and tournament info
        if (nextGame && Object.keys(nextGame).length > 0) {
            const game = Array.isArray(nextGame) ? nextGame[0] : nextGame;
            const gameData = game.data || game;
            const homeTeam = gameData.home_team?.name || 'Home';
            const awayTeam = gameData.away_team?.name || 'Away';
            const homeLogo = gameData.home_team?.url_logo || '';
            const awayLogo = gameData.away_team?.url_logo || '';
            const tournamentName = gameData.tournament?.name || '';
            const tournamentLogo = gameData.tournament?.url_logo || '';
            const matchId = gameData.match_id || gameData.id;
            const date = gameData.date ? new Date(gameData.date) : new Date();
            
            document.getElementById('next-game').innerHTML = `
                <div class="match-card" style="cursor:pointer;" onclick="goToMatch(${matchId})">
                    <div class="match-header">
                        ${tournamentLogo ? `<img src="${tournamentLogo}" alt="${tournamentName}" style="height: 20px; vertical-align: middle; margin-right: 5px;">` : '🏆'}
                        ${tournamentName} • Round ${gameData.round || '?'}
                    </div>
                    <div class="match-header">${date.toLocaleDateString()} ${date.toLocaleTimeString([], {hour:'2-digit', minute:'2-digit'})}</div>
                    <div class="match-score">
                        <div class="team-container">
                            ${homeLogo ? `<img src="${homeLogo}" alt="${homeTeam}" style="height: 50px; width: 50px; object-fit: contain;">` : '<div style="width: 50px;"></div>'}
                            <span class="team-name">${homeTeam}</span>
                        </div>
                        <span class="score">VS</span>
                        <div class="team-container">
                            ${awayLogo ? `<img src="${awayLogo}" alt="${awayTeam}" style="height: 50px; width: 50px; object-fit: contain;">` : '<div style="width: 50px;"></div>'}
                            <span class="team-name">${awayTeam}</span>
                        </div>
                    </div>
                    <div class="match-status">${gameData.status || 'Upcoming'}</div>
                </div>
            `;
        } else {
            document.getElementById('next-game').innerHTML = '<p>No upcoming games</p>';
        }

        // Standings with current team highlighted
        if (standings && Array.isArray(standings)) {
            renderStandings(standings);
        } else {
            document.getElementById('standing-tbody').innerHTML = '<tr><td colspan="5">No standings available</tr>';
        }

        // Manager
        if (manager) {
            const managerData = manager;
            const firstName = managerData.first_name || '';
            const lastName = managerData.last_name || '';
            const name = `${firstName} ${lastName}`.trim() || 'Manager';
            const managerId = managerData.manager_id;
            const nation = managerData.nation?.name || '';
            const flag = managerData.nation?.url_flag || '';
            const photo = managerData.url_photo || '';
            
            document.getElementById('manager-card').innerHTML = `
                <div class="card" style="cursor:pointer; padding: 1rem;" onclick="goToManager(${managerId})">
                    <div style="display: flex; align-items: center; gap: 15px;">
                        ${photo ? `<img src="${photo}" alt="${name}" style="width: 80px; height: 80px; border-radius: 50%; object-fit: cover;">` : 
                                  `<div style="width: 80px; height: 80px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 2rem;">👨‍✈️</div>`}
                        <div>
                            <h3 style="margin: 0;">${name}</h3>
                            <p>Manager</p>
                            ${nation ? `<p class="result-info">${flag ? `<img src="${flag}" style="width: 20px; vertical-align: middle;">` : '🌍'} ${nation}</p>` : ''}
                        </div>
                    </div>
                </div>
            `;
        } else {
            document.getElementById('manager-card').innerHTML = '<p>No manager information available</p>';
        }

        // Fixtures с пагинацией
        fixturesCurrentPage = 0;
        await loadFixturesWithPagination(1, false);
        
        // Tournament info
        if (details.tournament) {
            const tournamentInfo = document.getElementById('team-info');
            if (tournamentInfo) {
                tournamentInfo.innerHTML = `${details.tournament.name || 'Tournament'} ${details.tournament.season || ''}`;
            }
        }
        
        // Проверяем избранное
        if (TokenManager.hasToken()) {
            try {
                const favorites = await teamAPI.getFavorites();
                if (favorites && Array.isArray(favorites)) {
                    const isFavorite = favorites.some(f => f.team_id === parseInt(teamId));
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
        console.error('Error loading team data:', error);
        document.getElementById('team-detail').innerHTML = '<p class="loading">Error loading team data. Please try again.</p>';
    }
}

// Функция добавления селектора сезона для игроков
function addPlayersSeasonSelector() {
    const playersTab = document.getElementById('team-players');
    if (!playersTab) return;
    
    if (document.getElementById('players-season-filter')) return;
    
    const seasonSelectorHtml = `
        <div style="display: flex; justify-content: flex-end; align-items: center; margin-bottom: 1rem; gap: 0.5rem;">
            <label>Season:</label>
            <select id="players-season-filter" style="padding: 0.5rem; border-radius: 4px; border: 1px solid #ddd;">
                <option value="">Current Season</option>
            </select>
            <button id="apply-players-season" class="btn-primary" style="padding: 0.5rem 1rem;">Apply</button>
        </div>
    `;
    
    const table = playersTab.querySelector('.table');
    if (table) {
        table.insertAdjacentHTML('beforebegin', seasonSelectorHtml);
        
        const seasons = generateSeasonsList();
        const seasonSelect = document.getElementById('players-season-filter');
        seasons.forEach(season => {
            const option = document.createElement('option');
            option.value = season;
            option.textContent = season;
            seasonSelect.appendChild(option);
        });
        
        document.getElementById('apply-players-season').addEventListener('click', async () => {
            const season = document.getElementById('players-season-filter').value;
            await updateTeamPlayersStats(season || null);
            showFilterMessage(`Players stats for season ${season || 'current'}`, 'success');
        });
    }
}

function switchTeamTab(tabName) {
    const tabs = document.querySelectorAll('.tab');
    const panes = document.querySelectorAll('.tab-pane');
    
    tabs.forEach(tab => tab.classList.remove('active'));
    panes.forEach(pane => pane.classList.remove('active'));
    
    if (event && event.target) {
        event.target.classList.add('active');
    } else {
        tabs.forEach(tab => {
            if (tab.textContent.toLowerCase().includes(tabName)) {
                tab.classList.add('active');
            }
        });
    }
    
    const activePane = document.getElementById(`team-${tabName}`);
    if (activePane) {
        activePane.classList.add('active');
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
            await teamAPI.removeFavorite(teamId);
            btn.textContent = '★ Add to Favorites';
        } else {
            await teamAPI.addFavorite(teamId);
            btn.textContent = '★ Remove from Favorites';
        }
    } catch (error) {
        console.error('Error toggling favorite:', error);
    }
}

function formatLabel(key) {
    return key.replace(/_/g, ' ')
              .split(' ')
              .map(w => w.charAt(0).toUpperCase() + w.slice(1))
              .join(' ');
}

function goToMatch(id) {
    window.location.href = `/match?id=${id}`;
}

function goToPlayer(id) {
    window.location.href = `/player?id=${id}`;
}

function goToTeam(id) {
    window.location.href = `/team?id=${id}`;
}

function goToManager(id) {
    window.location.href = `/manager?id=${id}`;
}

window.loadMoreFixtures = loadMoreFixtures;