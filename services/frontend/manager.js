// Функционал страницы тренера
let managerId;
let isFavorite = false;
let currentStatsData = null;
let currentStatsType = 'aggregated';

// Переменные для сортировки
let currentTeamsSort = { column: 'team', direction: 'asc' };
let currentSeasonsSort = { column: 'season', direction: 'desc' };
let currentTeamsData = [];
let currentSeasonsData = [];

// Пагинация для матчей тренера
let fixturesCurrentPage = 0;
let fixturesLimit = 10;
let fixturesTotalCount = 0;
let isLoadingFixtures = false;
let currentFixturesSeason = null;

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

    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    if (loginForm) loginForm.addEventListener('submit', handleLogin);
    if (registerForm) registerForm.addEventListener('submit', handleRegister);
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
    
    console.log('Загрузка статистики тренера из:', url);
    try {
        const response = await fetch(url);
        if (response.ok) {
            return await response.json();
        }
    } catch (error) {
        console.error('Ошибка загрузки статистики тренера:', error);
    }
    return null;
}

// Общая статистика - выводим в виде сетки
function renderAggregatedStats(statsData) {
    const statsGrid = document.getElementById('manager-stats-grid');
    
    if (!statsData || !statsData.full) {
        statsGrid.innerHTML = '<p>Статистика недоступна</p>';
        return;
    }
    
    const stats = statsData.full;
    
    // Группируем статистику по категориям
    const generalStats = [
        { key: 'total_matches', label: 'Матчи' },
        { key: 'win_percentage', label: 'Побед %', multiply: 100, suffix: '%' },
        { key: 'avg_points', label: 'Ср. очков' }
    ];
    
    const attackingStats = [
        { key: 'goals', label: 'Голы' },
        { key: 'goals_per_90', label: 'Голов/90' },
        { key: 'total_shots', label: 'Всего ударов' },
        { key: 'shots_on_goal', label: 'Ударов в створ' }
    ];
    
    const defensiveStats = [
        { key: 'goals_conceded', label: 'Пропущено голов' },
        { key: 'goals_conceded_per_90', label: 'Пропущено/90' },
        { key: 'yellow_cards', label: 'Желтые карточки' },
        { key: 'red_cards', label: 'Красные карточки' }
    ];
    
    const possessionStats = [
        { key: 'average_ball_possession', label: 'Владение %', suffix: '%' },
        { key: 'average_pass_accuracy', label: 'Точность пасов %', multiply: 100, suffix: '%' },
        { key: 'total_passes', label: 'Всего пасов' },
        { key: 'complete_passes', label: 'Точных пасов' }
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
            <h3 style="margin-bottom: 0.75rem; color: #2c3e50;">📊 Общая</h3>
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
            <h3 style="margin-bottom: 0.75rem; color: #2c3e50;">⚽ Атака</h3>
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
            <h3 style="margin-bottom: 0.75rem; color: #2c3e50;">🛡️ Защита</h3>
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
            <h3 style="margin-bottom: 0.75rem; color: #2c3e50;">🎯 Владение и пасы</h3>
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
                        <th data-sort="team" onclick="sortTeamsData('team')">Команда <span class="sort-icon">${currentTeamsSort.column === 'team' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="matches" onclick="sortTeamsData('matches')">Матчи <span class="sort-icon">${currentTeamsSort.column === 'matches' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="win_rate" onclick="sortTeamsData('win_rate')">Побед % <span class="sort-icon">${currentTeamsSort.column === 'win_rate' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="avg_points" onclick="sortTeamsData('avg_points')">Ср. очков <span class="sort-icon">${currentTeamsSort.column === 'avg_points' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="possession" onclick="sortTeamsData('possession')">Владение <span class="sort-icon">${currentTeamsSort.column === 'possession' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="goals" onclick="sortTeamsData('goals')">Голы <span class="sort-icon">${currentTeamsSort.column === 'goals' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="goals_per_90" onclick="sortTeamsData('goals_per_90')">Голов/90 <span class="sort-icon">${currentTeamsSort.column === 'goals_per_90' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="conceded" onclick="sortTeamsData('conceded')">Пропущено <span class="sort-icon">${currentTeamsSort.column === 'conceded' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="pass_acc" onclick="sortTeamsData('pass_acc')">Точность пасов <span class="sort-icon">${currentTeamsSort.column === 'pass_acc' ? (currentTeamsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
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
                        <th data-sort="season" onclick="sortSeasonsData('season')">Сезон <span class="sort-icon">${currentSeasonsSort.column === 'season' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="matches" onclick="sortSeasonsData('matches')">Матчи <span class="sort-icon">${currentSeasonsSort.column === 'matches' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="win_rate" onclick="sortSeasonsData('win_rate')">Побед % <span class="sort-icon">${currentSeasonsSort.column === 'win_rate' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="avg_points" onclick="sortSeasonsData('avg_points')">Ср. очков <span class="sort-icon">${currentSeasonsSort.column === 'avg_points' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="possession" onclick="sortSeasonsData('possession')">Владение <span class="sort-icon">${currentSeasonsSort.column === 'possession' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="goals" onclick="sortSeasonsData('goals')">Голы <span class="sort-icon">${currentSeasonsSort.column === 'goals' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="goals_per_90" onclick="sortSeasonsData('goals_per_90')">Голов/90 <span class="sort-icon">${currentSeasonsSort.column === 'goals_per_90' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="conceded" onclick="sortSeasonsData('conceded')">Пропущено <span class="sort-icon">${currentSeasonsSort.column === 'conceded' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
                        <th data-sort="pass_acc" onclick="sortSeasonsData('pass_acc')">Точность пасов <span class="sort-icon">${currentSeasonsSort.column === 'pass_acc' ? (currentSeasonsSort.direction === 'asc' ? '🔼' : '🔽') : '↕️'}</span></th>
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
        document.getElementById('manager-stats-grid').innerHTML = '<p>Статистика по командам недоступна</p>';
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
        document.getElementById('manager-stats-grid').innerHTML = '<p>Статистика по сезонам недоступна</p>';
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
    statsGrid.innerHTML = '<div class="loading-spinner">Загрузка статистики...</div>';
    
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
        statsGrid.innerHTML = '<p>Статистика для выбранного режима недоступна</p>';
    }
}

async function loadManagerData() {
    try {
        console.log('Загрузка данных тренера для ID:', managerId);
        
        const [details, stats, teams, fixtures] = await Promise.all([
            managerAPI.getDetails(managerId).catch(e => {
                console.error('Ошибка деталей:', e);
                return null;
            }),
            managerAPI.getStats(managerId).catch(e => {
                console.error('Ошибка статистики:', e);
                return null;
            }),
            managerAPI.getTeams(managerId).catch(e => {
                console.error('Ошибка команд:', e);
                return null;
            }),
            managerAPI.getFixtures(managerId).catch(e => {
                console.error('Ошибка матчей:', e);
                return null;
            }),
        ]);

        console.log('Детали тренера:', details);
        console.log('Статистика тренера:', stats);
        console.log('Команды тренера:', teams);
        console.log('Матчи тренера:', fixtures);

        if (!details) {
            document.getElementById('manager-detail').innerHTML = '<p class="loading">Тренер не найден</p>';
            return;
        }

        const managerName = `${details.first_name || ''} ${details.last_name || ''}`.trim() || 'Тренер';
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

        // Команды - ИСПРАВЛЕННЫЙ БЛОК
        if (teams && Array.isArray(teams) && teams.length > 0) {
            const list = document.getElementById('teams-list');
            
            // Группируем команды по названию для отображения всех сезонов
            const teamsMap = new Map();
            
            for (const item of teams) {
                const teamName = item.team?.name || 'Неизвестно';
                const teamId = item.team?.team_id;
                const teamLogo = item.team?.url_logo || '';
                const season = item.season || 'Текущий';
                
                if (!teamsMap.has(teamId)) {
                    teamsMap.set(teamId, {
                        id: teamId,
                        name: teamName,
                        logo: teamLogo,
                        seasons: []
                    });
                }
                teamsMap.get(teamId).seasons.push(season);
            }
            
            // Сортируем команды по названию
            const sortedTeams = Array.from(teamsMap.values()).sort((a, b) => a.name.localeCompare(b.name));
            
            list.innerHTML = sortedTeams.map(team => `
                <div class="card team-card" onclick="searchTeam('${encodeURIComponent(team.name)}')" style="display: flex; align-items: center; gap: 1rem; cursor: pointer;">
                    ${team.logo ? `<img src="${team.logo}" alt="${escapeHtml(team.name)}" style="width: 50px; height: 50px; object-fit: contain; border-radius: 8px;">` : '<div style="width: 50px; height: 50px; background: #f0f0f0; border-radius: 8px; display: flex; align-items: center; justify-content: center;">🏆</div>'}
                    <div style="flex: 1;">
                        <h3 style="margin-bottom: 0.25rem;">${escapeHtml(team.name)}</h3>
                        <p style="font-size: 0.8rem; color: #666;">Сезоны: ${team.seasons.sort().reverse().join(', ')}</p>
                    </div>
                </div>
            `).join('');
        } else {
            document.getElementById('teams-list').innerHTML = '<p>История команд недоступна</p>';
        }

        // Матчи с пагинацией и выбором сезона
        addManagerFixturesSeasonSelector();
        fixturesCurrentPage = 0;
        currentFixturesSeason = null;
        await loadManagerFixturesWithPagination(1, false, null);
        
        if (TokenManager.hasToken()) {
            try {
                const favorites = await managerAPI.getFavorites();
                if (favorites && Array.isArray(favorites)) {
                    isFavorite = favorites.some(f => f.manager_id === parseInt(managerId));
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
        console.error('Ошибка загрузки данных тренера:', error);
        document.getElementById('manager-detail').innerHTML = '<p class="loading">Ошибка загрузки данных тренера. Пожалуйста, попробуйте снова.</p>';
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

// Функция загрузки матчей тренера с пагинацией и сезоном
async function loadManagerFixturesWithPagination(page = 1, append = false, season = null) {
    if (!managerId || isLoadingFixtures) return;
    
    isLoadingFixtures = true;
    const offset = (page - 1) * fixturesLimit + 1;
    
    try {
        let url = `/api/manager/${managerId}/fixtures?limit=${fixturesLimit}&offset=${offset}`;
        
        if (season && season !== '') {
            url += `&season=${encodeURIComponent(season)}`;
        }
        
        console.log('Загрузка матчей тренера из:', url);
        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error('Ошибка загрузки матчей');
        }
        
        const fixtures = await response.json();
        console.log('Матчи тренера загружены:', fixtures);
        
        if (!fixtures || !Array.isArray(fixtures)) {
            throw new Error('Некорректные данные матчей');
        }
        
        if (fixtures.length < fixturesLimit) {
            fixturesTotalCount = offset + fixtures.length;
        } else {
            fixturesTotalCount = offset + fixturesLimit + 1;
        }
        
        renderManagerFixturesList(fixtures, append);
        
    } catch (error) {
        console.error('Ошибка загрузки матчей тренера:', error);
        const list = document.getElementById('fixtures-list');
        if (!append && list) {
            list.innerHTML = '<p>Матчи недоступны</p>';
        }
    } finally {
        isLoadingFixtures = false;
    }
}

function renderManagerFixturesList(fixtures, append = false) {
    const list = document.getElementById('fixtures-list');
    if (!list) return;
    
    let html = '';
    if (append && list.innerHTML !== '<p>Матчи недоступны</p>') {
        html = list.innerHTML;
        html = html.replace(/<div id="load-more-fixtures-container"[\s\S]*?<\/div>/, '');
    }
    
    if (!fixtures || fixtures.length === 0) {
        if (!append) {
            list.innerHTML = '<p>Матчи недоступны</p>';
        }
        return;
    }
    
    html += fixtures.map(f => {
        const homeTeam = f.home_team?.name || 'Хозяева';
        const awayTeam = f.away_team?.name || 'Гости';
        const homeLogo = f.home_team?.url_logo || '';
        const awayLogo = f.away_team?.url_logo || '';
        const homeScore = f.home_team_score ?? '-';
        const awayScore = f.away_team_score ?? '-';
        const matchId = f.match_id;
        const date = f.date ? new Date(f.date) : new Date();
        const status = f.status || 'Scheduled';
        const tournament = f.tournament?.name || '';
        
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
            <div class="card match-card" onclick="goToMatch(${matchId})">
                <div class="match-header">
                    <span>${date.toLocaleDateString('ru-RU')}</span>
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
    
    const totalLoaded = (fixturesCurrentPage + 1) * fixturesLimit;
    
    if (totalLoaded < fixturesTotalCount) {
        html += `
            <div id="load-more-fixtures-container" style="text-align: center; margin-top: 1rem;">
                <button id="load-more-fixtures" class="btn-secondary" onclick="loadMoreManagerFixtures()">Показать ещё ↓</button>
            </div>
        `;
    } else {
        const existingContainer = document.getElementById('load-more-fixtures-container');
        if (existingContainer) {
            existingContainer.remove();
        }
    }
    
    list.innerHTML = html;
}

async function loadMoreManagerFixtures() {
    if (isLoadingFixtures) return;
    
    fixturesCurrentPage++;
    await loadManagerFixturesWithPagination(fixturesCurrentPage + 1, true, currentFixturesSeason);
    
    setTimeout(() => {
        const loadMoreContainer = document.getElementById('load-more-fixtures-container');
        if (loadMoreContainer) {
            loadMoreContainer.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
        }
    }, 100);
}

function addManagerFixturesSeasonSelector() {
    const fixturesTab = document.getElementById('manager-fixtures');
    if (!fixturesTab) return;
    
    if (document.getElementById('manager-fixtures-season-filter')) return;
    
    const seasonSelectorHtml = `
        <div style="display: flex; justify-content: flex-end; align-items: center; margin-bottom: 1rem; gap: 0.5rem;">
            <label>Сезон матчей:</label>
            <select id="manager-fixtures-season-filter" style="padding: 0.5rem; border-radius: 4px; border: 1px solid #ddd;">
                <option value="">Текущий сезон</option>
            </select>
            <button id="apply-manager-fixtures-season" class="btn-primary" style="padding: 0.5rem 1rem;">Применить</button>
        </div>
    `;
    
    const fixturesList = document.getElementById('fixtures-list');
    if (fixturesList) {
        fixturesList.insertAdjacentHTML('beforebegin', seasonSelectorHtml);
        
        const seasons = generateSeasonsList();
        const seasonSelect = document.getElementById('manager-fixtures-season-filter');
        seasons.forEach(season => {
            const option = document.createElement('option');
            option.value = season;
            option.textContent = season;
            seasonSelect.appendChild(option);
        });
        
        document.getElementById('apply-manager-fixtures-season').addEventListener('click', async () => {
            currentFixturesSeason = document.getElementById('manager-fixtures-season-filter').value;
            fixturesCurrentPage = 0;
            
            const existingContainer = document.getElementById('load-more-fixtures-container');
            if (existingContainer) {
                existingContainer.remove();
            }
            
            await loadManagerFixturesWithPagination(1, false, currentFixturesSeason || null);
            showFilterMessageForManager(`Матчи за сезон ${currentFixturesSeason || 'текущий'}`, 'success');
        });
    }
}

function showFilterMessageForManager(message, type = 'success') {
    const existingMsg = document.querySelector('.manager-filter-message');
    if (existingMsg) existingMsg.remove();
    
    const msgDiv = document.createElement('div');
    msgDiv.className = `manager-filter-message ${type}`;
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

function generateSeasonsList() {
    const seasons = [];
    for (let year = 2015; year <= 2025; year++) {
        seasons.push(`${year}/${year + 1}`);
    }
    return seasons;
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
            btn.textContent = '★ Добавить в избранное';
        } else {
            await managerAPI.addFavorite(managerId);
            isFavorite = true;
            btn.textContent = '★ Удалить из избранного';
        }
    } catch (error) {
        console.error('Ошибка переключения избранного:', error);
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

window.sortTeamsData = sortTeamsData;
window.sortSeasonsData = sortSeasonsData;
window.updateStats = updateStats;
window.loadMoreManagerFixtures = loadMoreManagerFixtures;