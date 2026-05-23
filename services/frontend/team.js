// Функционал страницы команды
let teamId;
let currentTeamData = {};
let currentStats = null;
let availableSeasons = [];
let currentPlayersSeason = null;

// Пагинация для матчей
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

// Функция загрузки матчей с пагинацией
async function loadFixturesWithPagination(page = 1, append = false) {
    if (!teamId || isLoadingFixtures) return;
    
    isLoadingFixtures = true;
    const offset = (page - 1) * fixturesLimit + 1;
    
    try {
        // Используем teamAPI.getFixtures с параметрами
        const fixtures = await teamAPI.getFixtures(teamId, fixturesLimit, offset);
        console.log('Матчи загружены:', fixtures);
        
        if (!fixtures || !Array.isArray(fixtures)) {
            throw new Error('Некорректные данные матчей');
        }
        
        if (fixtures.length < fixturesLimit) {
            fixturesTotalCount = offset + fixtures.length;
        } else {
            fixturesTotalCount = offset + fixturesLimit + 1;
        }
        
        renderFixturesList(fixtures, append);
        
    } catch (error) {
        console.error('Ошибка загрузки матчей:', error);
        const list = document.getElementById('fixtures-list');
        if (!append && list) {
            list.innerHTML = '<p>Матчи недоступны</p>';
        }
    } finally {
        isLoadingFixtures = false;
    }
}

function renderFixturesList(fixtures, append = false) {
    const list = document.getElementById('fixtures-list');
    if (!list) return;
    
    let html = '';
    if (append && list.innerHTML !== '<p>Матчи недоступны</p>') {
        html = list.innerHTML;
    }
    
    if (!fixtures || fixtures.length === 0) {
        if (!append) {
            list.innerHTML = '<p>Матчи недоступны</p>';
        }
        return;
    }
    
    html += fixtures.map(f => {
        const fixtureData = f.data || f;
        const homeTeam = fixtureData.home_team?.name || fixtureData.home_team_name || 'Хозяева';
        const awayTeam = fixtureData.away_team?.name || fixtureData.away_team_name || 'Гости';
        const homeLogo = fixtureData.home_team?.url_logo || '';
        const awayLogo = fixtureData.away_team?.url_logo || '';
        const homeScore = fixtureData.home_team_score ?? fixtureData.home_score ?? '-';
        const awayScore = fixtureData.away_team_score ?? fixtureData.away_score ?? '-';
        const matchId = fixtureData.match_id || fixtureData.id;
        const date = fixtureData.date ? new Date(fixtureData.date) : new Date();
        const status = fixtureData.status || 'Not started';
        const tournament = fixtureData.tournament?.name || '';
        
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
            <div class="card match-card" style="cursor:pointer; margin-bottom: 1rem;" onclick="goToMatch(${matchId})">
                <div class="match-header" style="display: flex; justify-content: space-between;">
                    <span>${date.toLocaleDateString('ru-RU')}</span>
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
                <button id="load-more-fixtures" class="btn-secondary" onclick="loadMoreFixtures()">Показать ещё ↓</button>
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
    
    console.log('Загрузка статистики команды из:', url);
    try {
        const response = await fetch(url);
        if (response.ok) {
            return await response.json();
        } else {
            console.error('Статистика: ответ не OK:', response.status);
            return null;
        }
    } catch (error) {
        console.error('Ошибка загрузки статистики:', error);
        return null;
    }
}

// Функция загрузки статистики игроков
async function loadTeamPlayersStatsWithSeason(season) {
    let url = `/api/team/${teamId}/players_stats`;
    if (season && season !== '') {
        url += `?season=${encodeURIComponent(season)}`;
    }
    
    console.log('Загрузка статистики игроков команды из:', url);
    try {
        const response = await fetch(url);
        if (response.ok) {
            return await response.json();
        } else {
            console.error('Статистика игроков: ответ не OK:', response.status);
            return null;
        }
    } catch (error) {
        console.error('Ошибка загрузки статистики игроков:', error);
        return null;
    }
}

// Функция обновления статистики
async function updateTeamStats() {
    const statsGrid = document.getElementById('team-stats-grid');
    statsGrid.innerHTML = '<div class="loading-spinner">Загрузка статистики...</div>';
    
    const stats = await loadTeamStatsWithFilters(
        currentStatsFilters.season, 
        currentStatsFilters.dateFrom, 
        currentStatsFilters.dateTo
    );
    
    if (stats) {
        currentStats = stats;
        renderTeamStats(stats);
        showFilterMessage('Статистика обновлена');
    } else {
        statsGrid.innerHTML = '<p>Статистика за выбранный период недоступна</p>';
        showFilterMessage('Данные за выбранный период не найдены', 'error');
    }
}

async function updateTeamPlayersStats(season) {
    const playersStats = await loadTeamPlayersStatsWithSeason(season);
    if (playersStats && typeof playersStats === 'object') {
        originalPlayersStats = playersStats;
        renderTeamPlayersStats(playersStats);
        setTimeout(addSortingHandlers, 100);
    } else {
        document.getElementById('players-tbody').innerHTML = '<tr><td colspan="8">Статистика игроков недоступна</td><\/tr>';
    }
}

// Функция отрисовки статистики команды
function renderTeamStats(stats) {
    const statsGrid = document.getElementById('team-stats-grid');
    const statsData = stats.data || stats;
    const statsToShow = Object.entries(statsData).slice(0, 12);
    
    const labelMap = {
        total_matches: 'Матчи',
        goals: 'Голы',
        goals_conceded: 'Пропущено',
        wins: 'Победы',
        draws: 'Ничьи',
        losses: 'Поражения',
        win_percentage: 'Побед %',
        avg_goals_per_match: 'Голов за матч',
        avg_possession: 'Владение %'
    };
    
    function formatLabel(key) {
        return labelMap[key] || key.replace(/_/g, ' ')
              .split(' ')
              .map(w => w.charAt(0).toUpperCase() + w.slice(1))
              .join(' ');
    }
    
    if (statsToShow.length > 0) {
        statsGrid.innerHTML = statsToShow.map(([key, value]) => {
            let displayValue = value;
            if (typeof value === 'number') {
                if (key.includes('percentage') || key.includes('accuracy')) {
                    displayValue = (value * 100).toFixed(1) + '%';
                } else if (Number.isInteger(value)) {
                    displayValue = value;
                } else {
                    displayValue = value.toFixed(2);
                }
            }
            return `
                <div class="stat-box">
                    <div class="stat-value">${displayValue ?? 'Н/Д'}</div>
                    <div class="stat-label">${formatLabel(key)}</div>
                </div>
            `;
        }).join('');
    } else {
        statsGrid.innerHTML = '<p>Статистика недоступна</p>';
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
        tbody.innerHTML = '<tr><td colspan="8">Статистика игроков недоступна</td><\/tr>';
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
    const positionNames = { 'G': 'Вратарь', 'D': 'Защитник', 'M': 'Полузащитник', 'F': 'Нападающий' };
    
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
                    full_name: `${player.first_name || ''} ${player.last_name || ''}`.trim() || 'Неизвестно'
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
        tbody.innerHTML = '<tr><td colspan="8">Статистика игроков недоступна</td><\/tr>';
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
        alert('Пожалуйста, выберите сезон или диапазон дат');
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
        showFilterMessage(`Таблица за сезон ${season || 'текущий'}`, 'success');
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
        console.error('Ошибка загрузки таблицы:', error);
    }
    return null;
}

function renderStandings(standings) {
    const tbody = document.getElementById('standing-tbody');
    if (!tbody) return;
    
    tbody.innerHTML = standings.map(standing => {
        const teamData = standing.team || standing;
        const teamName = teamData.name || teamData.team_name || 'Команда';
        const teamIdStanding = teamData.team_id || teamData.id;
        const teamLogo = teamData.url_logo || '';
        const isCurrentTeam = teamIdStanding == teamId;
        
        return `
            <tr onclick="goToTeam(${teamIdStanding})" style="cursor:pointer; ${isCurrentTeam ? 'background: linear-gradient(90deg, #e3f2fd, #bbdef5); font-weight: bold; border-left: 4px solid #3498db;' : ''}">
                <td><strong>${standing.position || standing.pos || 'Н/Д'}</strong></td>
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
        seasonFilter.innerHTML = '<option value="">Выберите сезон</option>';
        seasons.forEach(season => {
            const option = document.createElement('option');
            option.value = season;
            option.textContent = season;
            seasonFilter.appendChild(option);
        });
    }
    
    const tableSeasonFilter = document.getElementById('table-season-filter');
    if (tableSeasonFilter) {
        tableSeasonFilter.innerHTML = '<option value="">Текущий сезон</option>';
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
        console.log('Загрузка данных команды для ID:', teamId);
        
        const details = await teamAPI.getDetails(teamId).catch(e => {
            console.error('Ошибка деталей:', e);
            return null;
        });
        
        const stats = await teamAPI.getStats(teamId).catch(e => {
            console.error('Ошибка статистики:', e);
            return null;
        });
        
        const nextGame = await teamAPI.getNextGame(teamId).catch(e => {
            console.error('Ошибка следующего матча:', e);
            return null;
        });
        
        const standings = await teamAPI.getStandings(teamId).catch(e => {
            console.error('Ошибка таблицы:', e);
            return null;
        });
        
        const manager = await teamAPI.getManager(teamId).catch(e => {
            console.error('Ошибка тренера:', e);
            return null;
        });
        
        const fixtures = await teamAPI.getFixtures(teamId).catch(e => {
            console.error('Ошибка матчей:', e);
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
                console.error('Ошибка загрузки формы команды:', e);
            }
        }

        console.log('Детали команды:', details);
        console.log('Статистика команды:', stats);
        console.log('Следующий матч:', nextGame);
        console.log('Таблица:', standings);
        console.log('Тренер:', manager);
        console.log('Матчи:', fixtures);
        console.log('Форма команды:', teamForm);

        if (!details) {
            document.getElementById('team-detail').innerHTML = '<p class="loading">Команда не найдена</p>';
            return;
        }

        currentTeamData = details;

        // Обновляем заголовок с логотипом
        const teamName = details.name || details.team_name || 'Команда';
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
            let trendText = 'Стабильна';
            let trendColor = '#3498db';
            if (trend > 0.05) {
                trendIcon = '📈';
                trendText = 'Улучшается';
                trendColor = '#2ecc71';
            } else if (trend < -0.05) {
                trendIcon = '📉';
                trendText = 'Ухудшается';
                trendColor = '#e74c3c';
            }
            
            let formRating = '';
            let formColor = '';
            if (formIndex >= 70) {
                formRating = 'Отлично 🔥';
                formColor = '#2ecc71';
            } else if (formIndex >= 55) {
                formRating = 'Хорошо 👍';
                formColor = '#3498db';
            } else if (formIndex >= 40) {
                formRating = 'Средне 📊';
                formColor = '#f39c12';
            } else {
                formRating = 'Плохо 👎';
                formColor = '#e74c3c';
            }
            
            const formHtml = `
                <div id="team-form-analytics" class="form-analytics-section" style="margin-top: 2rem; padding: 1.5rem; background: linear-gradient(135deg, #1e3799 0%, #2c3e50 100%); border-radius: 12px; color: white;">
                    <h3 style="text-align: center; margin-bottom: 1rem;">📊 Аналитика формы команды (${season})</h3>
                    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem;">
                        <div class="form-card" style="text-align: center; background: rgba(255,255,255,0.15); border-radius: 8px; padding: 1rem;">
                            <div class="form-value" style="font-size: 2rem; font-weight: bold; color: ${formColor};">${formIndex}%</div>
                            <div class="form-label">Индекс формы</div>
                            <div class="form-rating" style="font-size: 0.8rem; margin-top: 0.5rem;">${formRating}</div>
                        </div>
                        <div class="form-card" style="text-align: center; background: rgba(255,255,255,0.15); border-radius: 8px; padding: 1rem;">
                            <div class="form-value" style="font-size: 2rem; font-weight: bold;">${attackIndex}%</div>
                            <div class="form-label">⚽ Индекс атаки</div>
                        </div>
                        <div class="form-card" style="text-align: center; background: rgba(255,255,255,0.15); border-radius: 8px; padding: 1rem;">
                            <div class="form-value" style="font-size: 2rem; font-weight: bold;">${defenseIndex}%</div>
                            <div class="form-label">🛡️ Индекс защиты</div>
                        </div>
                        <div class="form-card" style="text-align: center; background: rgba(255,255,255,0.15); border-radius: 8px; padding: 1rem;">
                            <div class="form-value" style="font-size: 2rem; font-weight: bold; color: ${trendColor};">${trendIcon} ${trend > 0 ? '+' : ''}${(trend * 100).toFixed(1)}%</div>
                            <div class="form-label">Тренд</div>
                            <div class="form-trend" style="font-size: 0.8rem;">${trendText}</div>
                        </div>
                    </div>
                    <div style="margin-top: 1rem; font-size: 0.8rem; text-align: center; opacity: 0.8;">
                        📊 На основе ${teamForm.matches_used || 0} последних матчей | Достоверность: ${confidence}%
                    </div>
                    ${teamForm.details ? `<div style="margin-top: 0.5rem; font-size: 0.7rem; text-align: center; opacity: 0.6;">${teamForm.details}</div>` : ''}
                </div>
            `;
            
            const statsGrid = document.getElementById('team-stats-grid');
            if (statsGrid && !document.getElementById('team-form-analytics')) {
                statsGrid.insertAdjacentHTML('afterend', formHtml);
            }
        }

        // Следующий матч с логотипами команд и информацией о турнире
        if (nextGame && Object.keys(nextGame).length > 0) {
            const game = Array.isArray(nextGame) ? nextGame[0] : nextGame;
            const gameData = game.data || game;
            const homeTeam = gameData.home_team?.name || 'Хозяева';
            const awayTeam = gameData.away_team?.name || 'Гости';
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
                        ${tournamentName} • Тур ${gameData.round || '?'}
                    </div>
                    <div class="match-header">${date.toLocaleDateString('ru-RU')} ${date.toLocaleTimeString('ru-RU', {hour:'2-digit', minute:'2-digit'})}</div>
                    <div class="match-score">
                        <div class="team-container">
                            ${homeLogo ? `<img src="${homeLogo}" alt="${homeTeam}" style="height: 50px; width: 50px; object-fit: contain;">` : '<div style="width: 50px;"></div>'}
                            <span class="team-name">${homeTeam}</span>
                        </div>
                        <span class="score">ПРОТИВ</span>
                        <div class="team-container">
                            ${awayLogo ? `<img src="${awayLogo}" alt="${awayTeam}" style="height: 50px; width: 50px; object-fit: contain;">` : '<div style="width: 50px;"></div>'}
                            <span class="team-name">${awayTeam}</span>
                        </div>
                    </div>
                    <div class="match-status">${gameData.status || 'Предстоящий'}</div>
                </div>
            `;
        } else {
            document.getElementById('next-game').innerHTML = '<p>Нет предстоящих матчей</p>';
        }

        // Турнирная таблица с подсветкой текущей команды
        if (standings && Array.isArray(standings)) {
            renderStandings(standings);
        } else {
            document.getElementById('standing-tbody').innerHTML = '<tr><td colspan="5">Таблица недоступна</td><\/tr>';
        }

        // Тренер
        if (manager) {
            const managerData = manager;
            const firstName = managerData.first_name || '';
            const lastName = managerData.last_name || '';
            const name = `${firstName} ${lastName}`.trim() || 'Тренер';
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
                            <p>Тренер</p>
                            ${nation ? `<p class="result-info">${flag ? `<img src="${flag}" style="width: 20px; vertical-align: middle;">` : '🌍'} ${nation}</p>` : ''}
                        </div>
                    </div>
                </div>
            `;
        } else {
            document.getElementById('manager-card').innerHTML = '<p>Информация о тренере недоступна</p>';
        }

        // Матчи с пагинацией
        fixturesCurrentPage = 0;
        await loadFixturesWithPagination(1, false);
        
        // Информация о турнире
        if (details.tournament) {
            const tournamentInfo = document.getElementById('team-info');
            if (tournamentInfo) {
                tournamentInfo.innerHTML = `${details.tournament.name || 'Турнир'} ${details.tournament.season || ''}`;
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
                        btn.textContent = isFavorite ? '★ Удалить из избранного' : '★ Добавить в избранное';
                    }
                }
            } catch (e) {
                console.error('Ошибка проверки избранного:', e);
            }
        }
        
    } catch (error) {
        console.error('Ошибка загрузки данных команды:', error);
        document.getElementById('team-detail').innerHTML = '<p class="loading">Ошибка загрузки данных команды. Пожалуйста, попробуйте снова.</p>';
    }
}

// Функция добавления селектора сезона для игроков
function addPlayersSeasonSelector() {
    const playersTab = document.getElementById('team-players');
    if (!playersTab) return;
    
    if (document.getElementById('players-season-filter')) return;
    
    const seasonSelectorHtml = `
        <div style="display: flex; justify-content: flex-end; align-items: center; margin-bottom: 1rem; gap: 0.5rem;">
            <label>Сезон:</label>
            <select id="players-season-filter" style="padding: 0.5rem; border-radius: 4px; border: 1px solid #ddd;">
                <option value="">Текущий сезон</option>
            </select>
            <button id="apply-players-season" class="btn-primary" style="padding: 0.5rem 1rem;">Применить</button>
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
            showFilterMessage(`Статистика игроков за сезон ${season || 'текущий'}`, 'success');
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
        if (btn.textContent.includes('Удалить')) {
            await teamAPI.removeFavorite(teamId);
            btn.textContent = '★ Добавить в избранное';
        } else {
            await teamAPI.addFavorite(teamId);
            btn.textContent = '★ Удалить из избранного';
        }
    } catch (error) {
        console.error('Ошибка переключения избранного:', error);
    }
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