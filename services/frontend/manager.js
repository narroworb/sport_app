// Manager page functionality
let managerId;
let isFavorite = false;
let currentStatsData = null;
let currentStatsType = 'aggregated';

// Переменные для сортировки
let currentTeamsSort = { column: 'team', direction: 'asc' };
let currentSeasonsSort = { column: 'season', direction: 'desc' };
let currentTeamsData = [];
let currentSeasonsData = [];

// Вспомогательная функция для экранирования HTML
function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function getManagerIdFromUrl() {
    const params = new URLSearchParams(window.location.search);
    return params.get('id');
}

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();

    managerId = getManagerIdFromUrl();
    
    if (!managerId) {
        const path = window.location.pathname;
        const pathParts = path.split('/');
        if (pathParts.length > 2 && pathParts[1] === 'manager') {
            managerId = pathParts[2];
        }
    }

    if (managerId) {
        await loadManagerData();
        if (TokenManager.hasToken()) {
            document.getElementById('fav-btn').style.display = 'inline-block';
        }
    }

    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
});

async function loadManagerStatsWithParams(params = {}) {
    let url = `/api/manager/${managerId}/stats`;
    const queryParams = [];
    
    if (params.by_team) {
        queryParams.push('by_team=true');
    }
    if (params.by_season) {
        queryParams.push('by_season=true');
    }
    
    if (queryParams.length > 0) {
        url += `?${queryParams.join('&')}`;
    }
    
    console.log('Loading manager stats from:', url);
    try {
        const response = await fetch(url);
        if (response.ok) {
            return await response.json();
        }
    } catch (error) {
        console.error('Error loading manager stats:', error);
    }
    return null;
}

// Общая статистика - выводим в виде сетки
function renderAggregatedStats(statsData) {
    const statsGrid = document.getElementById('manager-stats-grid');
    
    if (!statsData || !statsData.full) {
        statsGrid.innerHTML = '<p>No stats available</p>';
        return;
    }
    
    const stats = statsData.full;
    
    // Группируем статистику по категориям
    const generalStats = [
        { key: 'total_matches', label: 'Matches' },
        { key: 'win_percentage', label: 'Win %', multiply: 100, suffix: '%' },
        { key: 'avg_points', label: 'Avg Points' }
    ];
    
    const attackingStats = [
        { key: 'goals', label: 'Goals' },
        { key: 'goals_per_90', label: 'Goals/90' },
        { key: 'total_shots', label: 'Total Shots' },
        { key: 'shots_on_goal', label: 'Shots on Goal' }
    ];
    
    const defensiveStats = [
        { key: 'goals_conceded', label: 'Goals Conceded' },
        { key: 'goals_conceded_per_90', label: 'Conceded/90' },
        { key: 'yellow_cards', label: 'Yellow Cards' },
        { key: 'red_cards', label: 'Red Cards' }
    ];
    
    const possessionStats = [
        { key: 'average_ball_possession', label: 'Possession %', suffix: '%' },
        { key: 'average_pass_accuracy', label: 'Pass Accuracy %', multiply: 100, suffix: '%' },
        { key: 'total_passes', label: 'Total Passes' },
        { key: 'complete_passes', label: 'Complete Passes' }
    ];
    
    function formatValue(value, stat) {
        if (value === undefined || value === null) return '0';
        let displayValue = value;
        if (stat.multiply) {
            displayValue = value * stat.multiply;
        }
        if (typeof displayValue === 'number') {
            if (stat.suffix === '%' || (stat.key.includes('percentage') || stat.key.includes('accuracy'))) {
                return displayValue.toFixed(1) + '%';
            } else if (stat.key.includes('_per_90') && displayValue < 10) {
                return displayValue.toFixed(2);
            } else if (Number.isInteger(displayValue)) {
                return displayValue.toString();
            } else {
                return displayValue.toFixed(2);
            }
        }
        return displayValue;
    }
    
    let html = '<div style="display: flex; flex-direction: column; gap: 1.5rem;">';
    
    // Общая статистика
    html += `
        <div>
            <h3 style="margin-bottom: 0.75rem; color: #2c3e50;">📊 General</h3>
            <div class="stats-grid" style="grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));">
                ${generalStats.map(stat => {
                    const value = stats[stat.key];
                    return `
                        <div class="stat-box">
                            <div class="stat-value">${formatValue(value, stat)}</div>
                            <div class="stat-label">${stat.label}</div>
                        </div>
                    `;
                }).join('')}
            </div>
        </div>
    `;
    
    // Атакующая статистика
    html += `
        <div>
            <h3 style="margin-bottom: 0.75rem; color: #2c3e50;">⚽ Attacking</h3>
            <div class="stats-grid" style="grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));">
                ${attackingStats.map(stat => {
                    const value = stats[stat.key];
                    return `
                        <div class="stat-box">
                            <div class="stat-value">${formatValue(value, stat)}</div>
                            <div class="stat-label">${stat.label}</div>
                        </div>
                    `;
                }).join('')}
            </div>
        </div>
    `;
    
    // Оборонительная статистика
    html += `
        <div>
            <h3 style="margin-bottom: 0.75rem; color: #2c3e50;">🛡️ Defensive</h3>
            <div class="stats-grid" style="grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));">
                ${defensiveStats.map(stat => {
                    const value = stats[stat.key];
                    return `
                        <div class="stat-box">
                            <div class="stat-value">${formatValue(value, stat)}</div>
                            <div class="stat-label">${stat.label}</div>
                        </div>
                    `;
                }).join('')}
            </div>
        </div>
    `;
    
    // Владение и пасы
    html += `
        <div>
            <h3 style="margin-bottom: 0.75rem; color: #2c3e50;">🎯 Possession & Passing</h3>
            <div class="stats-grid" style="grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));">
                ${possessionStats.map(stat => {
                    const value = stats[stat.key];
                    return `
                        <div class="stat-box">
                            <div class="stat-value">${formatValue(value, stat)}</div>
                            <div class="stat-label">${stat.label}</div>
                        </div>
                    `;
                }).join('')}
            </div>
        </div>
    `;
    
    html += '</div>';
    statsGrid.innerHTML = html;
}

// Функция сортировки для таблицы команд
function sortTeamsData(column) {
    if (currentTeamsSort.column === column) {
        currentTeamsSort.direction = currentTeamsSort.direction === 'asc' ? 'desc' : 'asc';
    } else {
        currentTeamsSort.column = column;
        currentTeamsSort.direction = 'asc';
    }
    
    const sortedData = [...currentTeamsData].sort((a, b) => {
        let aVal, bVal;
        
        switch (column) {
            case 'team':
                aVal = a.name.toLowerCase();
                bVal = b.name.toLowerCase();
                return currentTeamsSort.direction === 'asc' 
                    ? aVal.localeCompare(bVal) 
                    : bVal.localeCompare(aVal);
            case 'matches':
                aVal = a.total_matches || 0;
                bVal = b.total_matches || 0;
                break;
            case 'win_rate':
                aVal = (a.win_percentage || 0) * 100;
                bVal = (b.win_percentage || 0) * 100;
                break;
            case 'avg_points':
                aVal = a.avg_points || 0;
                bVal = b.avg_points || 0;
                break;
            case 'possession':
                aVal = a.average_ball_possession || 0;
                bVal = b.average_ball_possession || 0;
                break;
            case 'goals':
                aVal = a.goals || 0;
                bVal = b.goals || 0;
                break;
            case 'goals_per_90':
                aVal = a.goals_per_90 || 0;
                bVal = b.goals_per_90 || 0;
                break;
            case 'conceded':
                aVal = a.goals_conceded || 0;
                bVal = b.goals_conceded || 0;
                break;
            case 'pass_acc':
                aVal = (a.average_pass_accuracy || 0) * 100;
                bVal = (b.average_pass_accuracy || 0) * 100;
                break;
            default:
                aVal = a.total_matches || 0;
                bVal = b.total_matches || 0;
        }
        
        if (currentTeamsSort.direction === 'asc') {
            return aVal > bVal ? 1 : -1;
        } else {
            return aVal < bVal ? 1 : -1;
        }
    });
    
    renderTeamsTable(sortedData);
}

// Функция сортировки для таблицы сезонов
function sortSeasonsData(column) {
    if (currentSeasonsSort.column === column) {
        currentSeasonsSort.direction = currentSeasonsSort.direction === 'asc' ? 'desc' : 'asc';
    } else {
        currentSeasonsSort.column = column;
        currentSeasonsSort.direction = 'asc';
    }
    
    const sortedData = [...currentSeasonsData].sort((a, b) => {
        let aVal, bVal;
        
        switch (column) {
            case 'season':
                aVal = a.season;
                bVal = b.season;
                return currentSeasonsSort.direction === 'asc' 
                    ? aVal.localeCompare(bVal) 
                    : bVal.localeCompare(aVal);
            case 'matches':
                aVal = a.total_matches || 0;
                bVal = b.total_matches || 0;
                break;
            case 'win_rate':
                aVal = (a.win_percentage || 0) * 100;
                bVal = (b.win_percentage || 0) * 100;
                break;
            case 'avg_points':
                aVal = a.avg_points || 0;
                bVal = b.avg_points || 0;
                break;
            case 'possession':
                aVal = a.average_ball_possession || 0;
                bVal = b.average_ball_possession || 0;
                break;
            case 'goals':
                aVal = a.goals || 0;
                bVal = b.goals || 0;
                break;
            case 'goals_per_90':
                aVal = a.goals_per_90 || 0;
                bVal = b.goals_per_90 || 0;
                break;
            case 'conceded':
                aVal = a.goals_conceded || 0;
                bVal = b.goals_conceded || 0;
                break;
            case 'pass_acc':
                aVal = (a.average_pass_accuracy || 0) * 100;
                bVal = (b.average_pass_accuracy || 0) * 100;
                break;
            default:
                aVal = a.total_matches || 0;
                bVal = b.total_matches || 0;
        }
        
        if (currentSeasonsSort.direction === 'asc') {
            return aVal > bVal ? 1 : -1;
        } else {
            return aVal < bVal ? 1 : -1;
        }
    });
    
    renderSeasonsTable(sortedData);
}

// Отрисовка таблицы команд
function renderTeamsTable(data) {
    const statsGrid = document.getElementById('manager-stats-grid');
    
    let html = `
        <div class="manager-stats-wrapper">
            <table class="manager-stats-table sortable-table">
                <thead>
                    <tr>
                        <th data-sort="team" onclick="sortTeamsData('team')">Team <span class="sort-icon">${currentTeamsSort.column === 'team' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="matches" onclick="sortTeamsData('matches')">Matches <span class="sort-icon">${currentTeamsSort.column === 'matches' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="win_rate" onclick="sortTeamsData('win_rate')">Win % <span class="sort-icon">${currentTeamsSort.column === 'win_rate' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="avg_points" onclick="sortTeamsData('avg_points')">Avg Points <span class="sort-icon">${currentTeamsSort.column === 'avg_points' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="possession" onclick="sortTeamsData('possession')">Possession <span class="sort-icon">${currentTeamsSort.column === 'possession' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="goals" onclick="sortTeamsData('goals')">Goals <span class="sort-icon">${currentTeamsSort.column === 'goals' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="goals_per_90" onclick="sortTeamsData('goals_per_90')">Goals/90 <span class="sort-icon">${currentTeamsSort.column === 'goals_per_90' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="conceded" onclick="sortTeamsData('conceded')">Conceded <span class="sort-icon">${currentTeamsSort.column === 'conceded' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="pass_acc" onclick="sortTeamsData('pass_acc')">Pass Acc <span class="sort-icon">${currentTeamsSort.column === 'pass_acc' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                    </tr>
                </thead>
                <tbody>
    `;
    
    for (const team of data) {
        const winPercentage = team.win_percentage ? (team.win_percentage * 100).toFixed(1) : 0;
        const avgPoints = team.avg_points ? team.avg_points.toFixed(2) : 0;
        const possession = team.average_ball_possession ? team.average_ball_possession.toFixed(1) : 0;
        const goalsPer90 = team.goals_per_90 ? team.goals_per_90.toFixed(2) : 0;
        const passAccuracy = team.average_pass_accuracy ? (team.average_pass_accuracy * 100).toFixed(1) : 0;
        
        html += `
            <tr onclick="searchTeam('${encodeURIComponent(team.name)}')">
                <td><strong>${escapeHtml(team.name)}</strong></td>
                <td>${team.total_matches || 0}</td>
                <td>${winPercentage}%</td>
                <td>${avgPoints}</td>
                <td>${possession}%</td>
                <td>${team.goals || 0}</td>
                <td>${goalsPer90}</td>
                <td>${team.goals_conceded || 0}</td>
                <td>${passAccuracy}%</td>
            </tr>
        `;
    }
    
    html += `
                </tbody>
            </table>
        </div>
    `;
    
    statsGrid.innerHTML = html;
}

// Отрисовка таблицы сезонов
function renderSeasonsTable(data) {
    const statsGrid = document.getElementById('manager-stats-grid');
    
    let html = `
        <div class="manager-stats-wrapper">
            <table class="manager-stats-table sortable-table">
                <thead>
                    <tr>
                        <th data-sort="season" onclick="sortSeasonsData('season')">Season <span class="sort-icon">${currentSeasonsSort.column === 'season' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="matches" onclick="sortSeasonsData('matches')">Matches <span class="sort-icon">${currentSeasonsSort.column === 'matches' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="win_rate" onclick="sortSeasonsData('win_rate')">Win % <span class="sort-icon">${currentSeasonsSort.column === 'win_rate' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="avg_points" onclick="sortSeasonsData('avg_points')">Avg Points <span class="sort-icon">${currentSeasonsSort.column === 'avg_points' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="possession" onclick="sortSeasonsData('possession')">Possession <span class="sort-icon">${currentSeasonsSort.column === 'possession' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="goals" onclick="sortSeasonsData('goals')">Goals <span class="sort-icon">${currentSeasonsSort.column === 'goals' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="goals_per_90" onclick="sortSeasonsData('goals_per_90')">Goals/90 <span class="sort-icon">${currentSeasonsSort.column === 'goals_per_90' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="conceded" onclick="sortSeasonsData('conceded')">Conceded <span class="sort-icon">${currentSeasonsSort.column === 'conceded' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="pass_acc" onclick="sortSeasonsData('pass_acc')">Pass Acc <span class="sort-icon">${currentSeasonsSort.column === 'pass_acc' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                    </tr>
                </thead>
                <tbody>
    `;
    
    for (const season of data) {
        const winPercentage = season.win_percentage ? (season.win_percentage * 100).toFixed(1) : 0;
        const avgPoints = season.avg_points ? season.avg_points.toFixed(2) : 0;
        const possession = season.average_ball_possession ? season.average_ball_possession.toFixed(1) : 0;
        const goalsPer90 = season.goals_per_90 ? season.goals_per_90.toFixed(2) : 0;
        const passAccuracy = season.average_pass_accuracy ? (season.average_pass_accuracy * 100).toFixed(1) : 0;
        
        html += `
            <tr>
                <td><strong>${escapeHtml(season.season)}</strong></td>
                <td>${season.total_matches || 0}</td>
                <td>${winPercentage}%</td>
                <td>${avgPoints}</td>
                <td>${possession}%</td>
                <td>${season.goals || 0}</td>
                <td>${goalsPer90}</td>
                <td>${season.goals_conceded || 0}</td>
                <td>${passAccuracy}%</td>
            </tr>
        `;
    }
    
    html += `
                </tbody>
            </table>
        </div>
    `;
    
    statsGrid.innerHTML = html;
}

function renderStatsByTeam(statsData) {
    if (!statsData || typeof statsData !== 'object' || Object.keys(statsData).length === 0) {
        document.getElementById('manager-stats-grid').innerHTML = '<p>No team stats available</p>';
        return;
    }
    
    // Преобразуем объект в массив для сортировки
    currentTeamsData = Object.keys(statsData).map(teamName => ({
        name: teamName,
        ...statsData[teamName]
    }));
    
    sortTeamsData(currentTeamsSort.column);
}

function renderStatsBySeason(statsData) {
    if (!statsData || typeof statsData !== 'object' || Object.keys(statsData).length === 0) {
        document.getElementById('manager-stats-grid').innerHTML = '<p>No season stats available</p>';
        return;
    }
    
    // Преобразуем объект в массив для сортировки
    currentSeasonsData = Object.keys(statsData).map(season => ({
        season: season,
        ...statsData[season]
    }));
    
    sortSeasonsData(currentSeasonsSort.column);
}

async function updateStats(type) {
    currentStatsType = type;
    const statsGrid = document.getElementById('manager-stats-grid');
    statsGrid.innerHTML = '<div class="loading-spinner">Loading statistics...</div>';
    
    let params = {};
    if (type === 'by_team') {
        params = { by_team: true };
    } else if (type === 'by_season') {
        params = { by_season: true };
    }
    
    const stats = await loadManagerStatsWithParams(params);
    if (stats) {
        currentStatsData = stats;
        if (type === 'by_team') {
            renderStatsByTeam(stats);
        } else if (type === 'by_season') {
            renderStatsBySeason(stats);
        } else {
            renderAggregatedStats(stats);
        }
    } else {
        statsGrid.innerHTML = '<p>No stats available for selected view</p>';
    }
}

async function loadManagerData() {
    try {
        console.log('Loading manager data for ID:', managerId);
        
        const [details, stats, teams, fixtures] = await Promise.all([
            managerAPI.getDetails(managerId).catch(e => {
                console.error('Details error:', e);
                return null;
            }),
            managerAPI.getStats(managerId).catch(e => {
                console.error('Stats error:', e);
                return null;
            }),
            managerAPI.getTeams(managerId).catch(e => {
                console.error('Teams error:', e);
                return null;
            }),
            managerAPI.getFixtures(managerId).catch(e => {
                console.error('Fixtures error:', e);
                return null;
            }),
        ]);

        console.log('Manager details:', details);
        console.log('Manager stats:', stats);
        console.log('Manager teams:', teams);
        console.log('Manager fixtures:', fixtures);

        if (!details) {
            document.getElementById('manager-detail').innerHTML = '<p class="loading">Manager not found</p>';
            return;
        }

        const managerName = `${details.first_name || ''} ${details.last_name || ''}`.trim() || 'Manager';
        const photo = details.url_photo || '';
        const nationFlag = details.nation?.url_flag || '';
        const nationName = details.nation?.name || '';
        
        document.getElementById('manager-name').innerHTML = `
            ${photo ? `<img src="${photo}" alt="${managerName}" style="height: 60px; border-radius: 50%; vertical-align: middle; margin-right: 15px;">` : ''}
            ${managerName}
        `;
        
        const nationSpan = document.getElementById('manager-nation');
        if (nationName) {
            nationSpan.innerHTML = `${nationFlag ? `<img src="${nationFlag}" style="width: 20px; vertical-align: middle; margin-right: 5px;">` : '🌍'} ${nationName}`;
        }

        currentStatsData = stats;
        renderAggregatedStats(stats);

        // Teams
        if (teams && typeof teams === 'object' && Object.keys(teams).length > 0) {
            const list = document.getElementById('teams-list');
            const teamsArray = [];
            
            for (const [teamName, seasons] of Object.entries(teams)) {
                if (Array.isArray(seasons)) {
                    seasons.forEach(season => {
                        teamsArray.push({ name: teamName, season: season });
                    });
                } else {
                    teamsArray.push({ name: teamName, season: seasons || 'Current' });
                }
            }
            
            list.innerHTML = teamsArray.map(team => `
                <div class="card team-card" onclick="searchTeam('${encodeURIComponent(team.name)}')">
                    <h3>${escapeHtml(team.name)}</h3>
                    <p>${escapeHtml(team.season)}</p>
                </div>
            `).join('');
        } else {
            document.getElementById('teams-list').innerHTML = '<p>No team history available</p>';
        }

        // Fixtures
        if (fixtures && Array.isArray(fixtures)) {
            const list = document.getElementById('fixtures-list');
            list.innerHTML = fixtures.slice(0, 20).map(f => {
                const homeTeam = f.home_team?.name || 'Home';
                const awayTeam = f.away_team?.name || 'Away';
                const homeLogo = f.home_team?.url_logo || '';
                const awayLogo = f.away_team?.url_logo || '';
                const homeScore = f.home_team_score ?? '-';
                const awayScore = f.away_team_score ?? '-';
                const matchId = f.match_id;
                const date = f.date ? new Date(f.date) : new Date();
                const status = f.status || 'Scheduled';
                const tournament = f.tournament?.name || '';
                
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
                    <div class="card match-card" onclick="goToMatch(${matchId})">
                        <div class="match-header">
                            <span>${date.toLocaleDateString()}</span>
                            <span>${escapeHtml(tournament)}</span>
                        </div>
                        <div class="match-score">
                            <div class="team-container">
                                ${homeLogo ? `<img src="${homeLogo}" alt="${escapeHtml(homeTeam)}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                                <span class="team-name">${escapeHtml(homeTeam)}</span>
                            </div>
                            <span class="score">${homeScore} : ${awayScore}</span>
                            <div class="team-container">
                                <span class="team-name">${escapeHtml(awayTeam)}</span>
                                ${awayLogo ? `<img src="${awayLogo}" alt="${escapeHtml(awayTeam)}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                            </div>
                        </div>
                        <div class="match-status ${statusClass}">
                            ${statusText}
                        </div>
                    </div>
                `;
            }).join('');
        } else {
            document.getElementById('fixtures-list').innerHTML = '<p>No fixtures available</p>';
        }
        
        if (TokenManager.hasToken()) {
            try {
                const favorites = await managerAPI.getFavorites();
                if (favorites && Array.isArray(favorites)) {
                    isFavorite = favorites.some(f => f.manager_id === parseInt(managerId));
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
        console.error('Error loading manager data:', error);
        document.getElementById('manager-detail').innerHTML = '<p class="loading">Error loading manager data. Please try again.</p>';
    }
}

function switchManagerTab(tabName) {
    const tabs = document.querySelectorAll('.tab');
    const panes = document.querySelectorAll('.tab-pane');
    
    tabs.forEach(tab => tab.classList.remove('active'));
    panes.forEach(pane => pane.classList.remove('active'));
    
    if (event && event.target) {
        event.target.classList.add('active');
    }
    
    const activePane = document.getElementById(`manager-${tabName}`);
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
        if (isFavorite) {
            await managerAPI.removeFavorite(managerId);
            isFavorite = false;
            btn.textContent = '★ Add to Favorites';
        } else {
            await managerAPI.addFavorite(managerId);
            isFavorite = true;
            btn.textContent = '★ Remove from Favorites';
        }
    } catch (error) {
        console.error('Error toggling favorite:', error);
    }
}

function goToMatch(id) {
    if (id && id > 0) {
        window.location.href = `/match?id=${id}`;
    }
}

function searchTeam(teamName) {
    if (teamName) {
        window.location.href = `/search?q=${encodeURIComponent(teamName)}`;
    }
}

// Делаем функции сортировки глобальными
window.sortTeamsData = sortTeamsData;
window.sortSeasonsData = sortSeasonsData;
window.updateStats = updateStats;