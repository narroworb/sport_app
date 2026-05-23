// Функционал страницы игрока
let playerId;
let isFavorite = false;
let currentStats = null;
let currentTeams = null;
let currentFixtures = null;
let availableSeasons = [];
let currentSimilarParams = {
    top_k: 5,
    min_minutes: 100,
    season: null
};

// Пагинация для матчей
let fixturesCurrentPage = 0;
let fixturesLimit = 10;
let fixturesTotalCount = 0;
let isLoadingFixtures = false;

function getPlayerIdFromUrl() {
    const params = new URLSearchParams(window.location.search);
    return params.get('id');
}

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();

    playerId = getPlayerIdFromUrl();

    if (!playerId) {
        const path = window.location.pathname;
        const pathParts = path.split('/');
        if (pathParts.length > 2 && pathParts[1] === 'player') {
            playerId = pathParts[2];
        }
    }

    if (playerId) {
        await loadPlayerData();
        if (TokenManager.hasToken()) {
            document.getElementById('fav-btn').style.display = 'inline-block';
        }
    }

    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
    
    // Навешиваем обработчики на фильтры статистики
    const applySeasonBtn = document.getElementById('apply-season');
    const applyDateBtn = document.getElementById('apply-date');
    const resetBtn = document.getElementById('reset-filters');
    
    if (applySeasonBtn) applySeasonBtn.addEventListener('click', applySeasonFilter);
    if (applyDateBtn) applyDateBtn.addEventListener('click', applyDateFilter);
    if (resetBtn) resetBtn.addEventListener('click', resetFilters);
    
    // Навешиваем обработчик на кнопку поиска похожих игроков
    const findSimilarBtn = document.getElementById('find-similar-btn');
    if (findSimilarBtn) {
        findSimilarBtn.addEventListener('click', () => loadSimilarPlayers(true));
    }
});

// Генерация списка сезонов с 2015/2016 по 2025/2026
function generateSeasonsList() {
    const seasons = [];
    for (let year = 2015; year <= 2025; year++) {
        seasons.push(`${year}/${year + 1}`);
    }
    return seasons;
}

// Функция заполнения сезонов для выбора в настройках
function populateSimilarSeasons() {
    const seasonSelect = document.getElementById('similar-season');
    if (!seasonSelect) return;
    
    // Очищаем существующие опции
    seasonSelect.innerHTML = '<option value="">Автоопределение</option>';
    
    // Добавляем сезоны от 2015/2016 до 2025/2026
    for (let year = 2025; year >= 2015; year--) {
        const season = `${year}/${year + 1}`;
        const option = document.createElement('option');
        option.value = season;
        option.textContent = season;
        seasonSelect.appendChild(option);
    }
}

// Функция загрузки матчей с пагинацией
async function loadFixturesWithPagination(page = 1, append = false) {
    if (!playerId || isLoadingFixtures) return;
    
    isLoadingFixtures = true;
    const offset = (page - 1) * fixturesLimit + 1;
    
    try {
        const url = `/api/player/${playerId}/fixtures?limit=${fixturesLimit}&offset=${offset}`;
        console.log('Загрузка матчей из:', url);
        
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Не удалось загрузить матчи: ${response.status}`);
        }
        
        const fixtures = await response.json();
        console.log('Матчи загружены:', fixtures);
        
        // Обновляем общее количество (если API возвращает total, иначе используем длину массива)
        if (fixtures.length < fixturesLimit) {
            fixturesTotalCount = offset + fixtures.length;
        } else {
            fixturesTotalCount = offset + fixturesLimit + 1; // Приблизительное значение
        }
        
        renderFixturesList(fixtures, append);
        
    } catch (error) {
        console.error('Ошибка загрузки матчей:', error);
        const list = document.getElementById('fixtures-list');
        if (!append && list) {
            list.innerHTML = '<p>Матчи не найдены</p>';
        }
    } finally {
        isLoadingFixtures = false;
    }
}

// Функция отрисовки списка матчей
function renderFixturesList(fixtures, append = false) {
    const list = document.getElementById('fixtures-list');
    if (!list) return;
    
    let html = '';
    
    if (append && list.innerHTML !== '<p>Матчи не найдены</p>') {
        html = list.innerHTML;
    } else {
        html = '';
    }
    
    if (!fixtures || fixtures.length === 0) {
        if (!append) {
            list.innerHTML = '<p>Матчи не найдены</p>';
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
        const goals = f.athlete_goals || 0;
        const assists = f.athlete_assists || 0;
        const rating = f.athlete_rating ? f.athlete_rating.toFixed(1) : '-';
        const minutes = f.athlete_minutes_played || 0;
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
                <div class="match-status ${statusClass}" style="display: flex; justify-content: space-between; margin-top: 10px; padding-top: 10px; border-top: 1px solid rgba(255,255,255,0.2);">
                    <span>${statusText}</span>
                    <span>⚽ ${goals} голов | 🎯 ${assists} передач | ⭐ ${rating} | ⏱️ ${minutes} мин.</span>
                </div>
            </div>
        `;
    }).join('');
    
    // Добавляем кнопку "Показать ещё" если есть еще данные
    const currentCount = (fixturesCurrentPage + 1) * fixturesLimit;
    if (currentCount < fixturesTotalCount && !append) {
        html += `
            <div style="text-align: center; margin-top: 1rem;">
                <button id="load-more-fixtures" class="btn-secondary" onclick="loadMoreFixtures()">Показать ещё ↓</button>
            </div>
        `;
    }
    
    list.innerHTML = html;
}

// Функция загрузки следующих страниц
async function loadMoreFixtures() {
    fixturesCurrentPage++;
    await loadFixturesWithPagination(fixturesCurrentPage, true);
}

// Функция загрузки похожих игроков с параметрами
async function loadSimilarPlayers(showLoading = true) {
    if (!playerId) return;
    
    // Получаем параметры из UI
    const topK = parseInt(document.getElementById('similar-top-k')?.value || 5);
    const minMinutes = parseInt(document.getElementById('similar-min-minutes')?.value || 100);
    let season = document.getElementById('similar-season')?.value || null;
    
    // Если сезон не выбран, используем текущий или первый доступный
    if (!season && availableSeasons.length > 0) {
        season = availableSeasons[0];
    }
    if (!season) {
        season = getCurrentSeason();
    }
    
    currentSimilarParams = { top_k: topK, min_minutes: minMinutes, season };
    
    console.log('Загрузка похожих игроков с параметрами:', currentSimilarParams);
    
    // Показываем индикатор загрузки
    const statsGrid = document.getElementById('player-stats-grid');
    const existingBlock = document.getElementById('player-similarity');
    
    if (showLoading) {
        if (existingBlock) {
            existingBlock.innerHTML = '<div style="text-align: center; padding: 2rem;">🔄 Загрузка похожих игроков...</div>';
        } else if (statsGrid) {
            statsGrid.insertAdjacentHTML('afterend', '<div id="temp-similar-loading" style="margin-top: 2rem; padding: 1.5rem; background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%); border-radius: 12px; color: white; text-align: center;">🔄 Загрузка похожих игроков...</div>');
        }
    }
    
    try {
        const token = TokenManager.getToken();
        const url = `/api/analytics/player_similarity?player_id=${playerId}&season=${encodeURIComponent(season)}&top_k=${topK}&min_minutes=${minMinutes}`;
        
        const headers = {};
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        
        const response = await fetch(url, { headers });
        
        if (!response.ok) {
            throw new Error(`Ошибка API: ${response.status}`);
        }
        
        const similarPlayers = await response.json();
        console.log('Похожие игроки получены:', similarPlayers);
        
        // Обновляем блок похожих игроков
        await updateSimilarPlayersBlock(similarPlayers, season);
        
    } catch (error) {
        console.error('Ошибка загрузки похожих игроков:', error);
        showFilterMessage('Не удалось загрузить похожих игроков: ' + error.message, 'error');
        
        // Показываем ошибку в блоке
        const existingBlock = document.getElementById('player-similarity');
        if (existingBlock) {
            existingBlock.innerHTML = `
                <div style="text-align: center; padding: 2rem;">
                    <p>❌ Не удалось загрузить похожих игроков</p>
                    <p style="font-size: 0.8rem; opacity: 0.7;">${error.message}</p>
                    <button onclick="loadSimilarPlayers(true)" class="btn-primary" style="margin-top: 1rem;">🔄 Попробовать снова</button>
                </div>
            `;
        }
    } finally {
        // Удаляем временный индикатор загрузки
        const tempLoading = document.getElementById('temp-similar-loading');
        if (tempLoading) tempLoading.remove();
    }
}

// Функция обновления блока похожих игроков
async function updateSimilarPlayersBlock(similarPlayers, season) {
    if (!similarPlayers || !similarPlayers.players || similarPlayers.players.length === 0) {
        const existingBlock = document.getElementById('player-similarity');
        if (existingBlock) {
            existingBlock.innerHTML = `
                <div style="text-align: center; padding: 2rem;">
                    <p>🔍 Похожие игроки не найдены</p>
                    <p style="font-size: 0.8rem; opacity: 0.7;">Попробуйте изменить параметры поиска</p>
                </div>
            `;
        }
        return;
    }
    
    // Загружаем детали похожих игроков
    const similarPlayersDetails = [];
    for (const similar of similarPlayers.players.slice(0, currentSimilarParams.top_k)) {
        const similarDetail = await playerAPI.getDetails(similar.player_id).catch(e => {
            console.error(`Ошибка загрузки похожего игрока ${similar.player_id}:`, e);
            return null;
        });
        if (similarDetail) {
            similarPlayersDetails.push({ ...similar, details: similarDetail });
        }
    }
    
    const positionNames = { 'G': 'Вратарь', 'D': 'Защитник', 'M': 'Полузащитник', 'F': 'Нападающий' };
    const positionDisplay = positionNames[similarPlayers.position] || similarPlayers.position || 'Н/Д';
    
    const similarityHtml = `
        <div id="player-similarity" class="similar-players-section" style="margin-top: 2rem; padding: 1.5rem; background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%); border-radius: 12px; color: white;">
            <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; margin-bottom: 1rem;">
                <h3 style="margin: 0;">🔄 Похожие игроки (${positionDisplay})</h3>
                <div style="font-size: 0.7rem; opacity: 0.7;">
                    ⚙️ Топ ${currentSimilarParams.top_k} | Мин ${currentSimilarParams.min_minutes} мин | Сезон ${season}
                </div>
            </div>
            <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1rem;">
                ${similarPlayersDetails.map((player, index) => {
                    const similarityPercent = (player.similarity * 100).toFixed(1);
                    const playerName = `${player.details.first_name || ''} ${player.details.last_name || ''}`.trim() || 'Неизвестно';
                    const playerPhoto = player.details.url_photo || '';
                    const playerPosition = player.details.position || 'Н/Д';
                    const positionDisplayPlayer = positionNames[playerPosition] || playerPosition;
                    
                    let similarityColor = '#e74c3c';
                    if (similarityPercent >= 70) similarityColor = '#2ecc71';
                    else if (similarityPercent >= 50) similarityColor = '#f39c12';
                    
                    return `
                        <div class="similar-player-card" style="background: rgba(255,255,255,0.15); border-radius: 8px; padding: 1rem; cursor: pointer; transition: transform 0.2s;" onclick="goToPlayer(${player.player_id})">
                            <div style="display: flex; align-items: center; gap: 15px;">
                                ${playerPhoto ? 
                                    `<img src="${playerPhoto}" alt="${playerName}" style="width: 60px; height: 60px; border-radius: 50%; object-fit: cover;">` : 
                                    `<div style="width: 60px; height: 60px; background: rgba(255,255,255,0.3); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 2rem;">👤</div>`
                                }
                                <div style="flex: 1;">
                                    <div style="font-weight: bold; font-size: 1.1rem;">${playerName}</div>
                                    <div style="font-size: 0.8rem; opacity: 0.9;">${positionDisplayPlayer}</div>
                                    <div style="margin-top: 5px;">
                                        <div style="background: rgba(255,255,255,0.2); border-radius: 10px; height: 6px; width: 100%;">
                                            <div style="background: ${similarityColor}; width: ${similarityPercent}%; height: 6px; border-radius: 10px;"></div>
                                        </div>
                                        <div style="font-size: 0.8rem; margin-top: 3px;">Совпадение: ${similarityPercent}%</div>
                                    </div>
                                </div>
                                <div style="font-size: 1.2rem;">#${index + 1}</div>
                            </div>
                        </div>
                    `;
                }).join('')}
            </div>
            <div style="margin-top: 1rem; font-size: 0.7rem; text-align: center; opacity: 0.7;">
                ${similarPlayers.details || `На основе сезона ${season} | Кандидатов: ${similarPlayers.candidates || 'Н/Д'}`}
            </div>
        </div>
    `;
    
    const statsGrid = document.getElementById('player-stats-grid');
    const existingBlock = document.getElementById('player-similarity');
    
    if (existingBlock) {
        existingBlock.outerHTML = similarityHtml;
    } else if (statsGrid) {
        statsGrid.insertAdjacentHTML('afterend', similarityHtml);
    }
    
    showFilterMessage(`Найдено ${similarPlayersDetails.length} похожих игроков`, 'success');
}

async function loadPlayerStatsWithFilters(season, dateFrom, dateTo) {
    let url = `/api/player/${playerId}/stats`;
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
    
    console.log('Загрузка статистики из:', url);
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

async function applySeasonFilter() {
    const seasonFilter = document.getElementById('season-filter').value;
    
    if (!seasonFilter) {
        alert('Пожалуйста, выберите сезон');
        return;
    }
    
    console.log('Применение фильтра по сезону:', seasonFilter);
    
    const statsGrid = document.getElementById('player-stats-grid');
    statsGrid.innerHTML = '<div class="loading-spinner">Загрузка статистики...</div>';
    
    const stats = await loadPlayerStatsWithFilters(seasonFilter, null, null);
    
    if (stats && Object.keys(stats).length > 0) {
        currentStats = stats;
        updateStatsGrid(stats);
        showFilterMessage(`Отображается статистика за сезон ${seasonFilter}`);
    } else {
        statsGrid.innerHTML = '<p>Данные за выбранный сезон отсутствуют</p>';
        showFilterMessage(`Данные за сезон ${seasonFilter} не найдены`, 'error');
    }
}

async function applyDateFilter() {
    const dateFrom = document.getElementById('date-from').value;
    const dateTo = document.getElementById('date-to').value;
    
    if (!dateFrom && !dateTo) {
        alert('Пожалуйста, выберите хотя бы одну дату');
        return;
    }
    
    console.log('Применение фильтра по датам:', { dateFrom, dateTo });
    
    const statsGrid = document.getElementById('player-stats-grid');
    statsGrid.innerHTML = '<div class="loading-spinner">Загрузка статистики...</div>';
    
    const stats = await loadPlayerStatsWithFilters(null, dateFrom, dateTo);
    
    if (stats && Object.keys(stats).length > 0) {
        currentStats = stats;
        updateStatsGrid(stats);
        let message = 'Отображаются статистические данные';
        if (dateFrom && dateTo) {
            message += ` с ${dateFrom} по ${dateTo}`;
        } else if (dateFrom) {
            message += ` с ${dateFrom}`;
        } else if (dateTo) {
            message += ` до ${dateTo}`;
        }
        showFilterMessage(message);
    } else {
        statsGrid.innerHTML = '<p>Данные за выбранный диапазон отсутствуют</p>';
        showFilterMessage('Данные за выбранный диапазон отсутствуют', 'error');
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

function resetFilters() {
    document.getElementById('season-filter').value = '';
    document.getElementById('date-from').value = '';
    document.getElementById('date-to').value = '';
    
    const statsGrid = document.getElementById('player-stats-grid');
    statsGrid.innerHTML = '<div class="loading-spinner">Загрузка статистики...</div>';
    
    loadPlayerStatsWithFilters(null, null, null).then(stats => {
        if (stats && Object.keys(stats).length > 0) {
            currentStats = stats;
            updateStatsGrid(stats);
            showFilterMessage('Отображаются статистические данные за все время');
        } else {
            statsGrid.innerHTML = '<p>Данные отсутствуют</p>';
        }
    });
}

function updateStatsGrid(stats) {
    const statsGrid = document.getElementById('player-stats-grid');
    
    if (!stats || Object.keys(stats).length === 0) {
        statsGrid.innerHTML = '<p>Статистика за выбранный период недоступна</p>';
        return;
    }
    
    const statsToShow = [
        { key: 'total_matches', label: 'Матчи', isPercent: false },
        { key: 'goals', label: 'Голы', isPercent: false },
        { key: 'assists', label: 'Передачи', isPercent: false },
        { key: 'avg_rating', label: 'Ср. рейтинг', isPercent: false },
        { key: 'minutes_played', label: 'Минуты', isPercent: false },
        { key: 'yellow_cards', label: 'Желтые карточки', isPercent: false },
        { key: 'red_cards', label: 'Красные карточки', isPercent: false },
        { key: 'total_tackles', label: 'Отборы', isPercent: false },
        { key: 'interceptions', label: 'Перехваты', isPercent: false },
        { key: 'saves', label: 'Сейвы', isPercent: false },
        { key: 'pass_accuracy', label: 'Точность пасов', isPercent: true },
        { key: 'dribble_accuracy', label: 'Точность дриблинга', isPercent: true }
    ];
    
    const passAccuracy = stats.complete_passes && stats.pass_attempts && stats.pass_attempts > 0
        ? (stats.complete_passes / stats.pass_attempts) * 100 
        : null;
    
    const displayStats = [];
    for (const stat of statsToShow) {
        let value = stats[stat.key];
        if (stat.key === 'pass_accuracy' && passAccuracy !== null) {
            value = passAccuracy;
        }
        if (value !== undefined && value !== null && value !== 0) {
            let displayValue = value;
            if (typeof value === 'number') {
                if (stat.isPercent || stat.key.includes('accuracy') || stat.key.includes('percent')) {
                    displayValue = value.toFixed(1) + '%';
                } else if (Number.isInteger(value)) {
                    displayValue = value;
                } else {
                    displayValue = value.toFixed(2);
                }
            }
            displayStats.push({ label: stat.label, value: displayValue });
        }
    }
    
    if (stats.goals && stats.minutes_played && stats.minutes_played > 0) {
        const goalsPer90 = (stats.goals / stats.minutes_played) * 90;
        displayStats.push({ label: 'Голов/90', value: goalsPer90.toFixed(2) });
    }
    
    if (stats.saves && stats.minutes_played && stats.minutes_played > 0) {
        const savesPer90 = (stats.saves / stats.minutes_played) * 90;
        displayStats.push({ label: 'Сейвов/90', value: savesPer90.toFixed(2) });
    }
    
    if (displayStats.length > 0) {
        statsGrid.innerHTML = displayStats.map(stat => `
            <div class="stat-box">
                <div class="stat-value">${stat.value}</div>
                <div class="stat-label">${stat.label}</div>
            </div>
        `).join('');
    } else {
        statsGrid.innerHTML = '<p>Статистика недоступна</p>';
    }
}

async function loadPlayerData() {
    try {
        console.log('Загрузка данных игрока для ID:', playerId);
        
        const [details, teams, fixtures] = await Promise.all([
            playerAPI.getDetails(playerId).catch(e => {
                console.error('Ошибка деталей:', e);
                return null;
            }),
            playerAPI.getTeams(playerId).catch(e => {
                console.error('Ошибка команд:', e);
                return null;
            }),
            playerAPI.getFixtures(playerId).catch(e => {
                console.error('Ошибка матчей:', e);
                return null;
            }),
        ]);

        console.log('Детали игрока:', details);
        console.log('Команды игрока:', teams);
        console.log('Матчи игрока:', fixtures);

        if (!details) {
            document.getElementById('player-detail').innerHTML = '<p class="loading">Игрок не найден</p>';
            return;
        }

        currentTeams = teams;
        currentFixtures = fixtures;

        // Генерируем список сезонов от 2015/2016 до 2025/2026
        const allSeasons = generateSeasonsList();
        
        if (teams && Array.isArray(teams)) {
            teams.forEach(team => {
                if (team.season && !allSeasons.includes(team.season)) {
                    allSeasons.push(team.season);
                }
            });
        }
        
        availableSeasons = [...new Set(allSeasons)].sort().reverse();
        
        const seasonFilter = document.getElementById('season-filter');
        if (seasonFilter) {
            while (seasonFilter.options.length > 1) {
                seasonFilter.remove(1);
            }
            
            availableSeasons.forEach(season => {
                const option = document.createElement('option');
                option.value = season;
                option.textContent = season;
                seasonFilter.appendChild(option);
            });
            console.log('Сезонов загружено:', availableSeasons.length);
        }
        
        // Заполняем сезоны для выбора в настройках похожих игроков
        populateSimilarSeasons();

        const stats = await loadPlayerStatsWithFilters(null, null, null);
        if (stats) {
            currentStats = stats;
            updateStatsGrid(stats);
        }

        // Обновляем заголовок с фото
        const playerName = `${details.first_name || ''} ${details.last_name || ''}`.trim() || 'Игрок';
        const photo = details.url_photo || '';
        const nationFlag = details.nation?.url_flag || '';
        const nationName = details.nation?.name || '';

        document.getElementById('player-name').innerHTML = `
            ${photo ? `<img src="${photo}" alt="${playerName}" style="height: 60px; border-radius: 50%; vertical-align: middle; margin-right: 15px;">` : ''}
            ${playerName}
        `;

        const positionSpan = document.getElementById('player-position');
        if (details.position) {
            const positionNames = { 'G': 'Вратарь', 'D': 'Защитник', 'M': 'Полузащитник', 'F': 'Нападающий' };
            positionSpan.textContent = positionNames[details.position] || details.position;
        }

        const nationSpan = document.getElementById('player-nation');
        if (nationName) {
            nationSpan.innerHTML = `${nationFlag ? `<img src="${nationFlag}" style="width: 20px; vertical-align: middle; margin-right: 5px;">` : '🌍'} ${nationName}`;
        }

        // Дополнительная информация
        if (details.date_of_birth) {
            const birthDate = new Date(details.date_of_birth);
            const formattedDate = birthDate.toLocaleDateString('ru-RU', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            });
            const ageDifMs = Date.now() - birthDate.getTime();
            const ageDate = new Date(ageDifMs);
            const age = Math.abs(ageDate.getUTCFullYear() - 1970);
            document.getElementById('player-birthdate').innerHTML = `${formattedDate} <span style="color: #666;">(${age} лет)</span>`;
        } else {
            document.getElementById('player-birthdate').textContent = 'Не указано';
        }

        if (details.height) {
            const heightCm = details.height;
            const heightFeet = Math.floor(heightCm / 30.48);
            const heightInches = Math.round((heightCm % 30.48) / 2.54);
            document.getElementById('player-height').innerHTML = `${heightCm} см (${heightFeet}'${heightInches}")`;
        } else {
            document.getElementById('player-height').textContent = 'Не указано';
        }

        if (details.preffered_foot) {
            const footNames = { 'Left': '👈 Левая', 'Right': '👉 Правая', 'Both': '✌️ Обе' };
            document.getElementById('player-foot').innerHTML = footNames[details.preffered_foot] || details.preffered_foot;
        } else {
            document.getElementById('player-foot').textContent = 'Не указано';
        }

        if (details.current_status) {
            const statusBadge = details.current_status === 'Active' 
                ? '<span style="color: #2ecc71;">✅ Активен</span>' 
                : '<span style="color: #e74c3c;">❌ Неактивен</span>';
            document.getElementById('player-status').innerHTML = statusBadge;
        } else {
            document.getElementById('player-status').textContent = 'Неизвестно';
        }

        // Загружаем похожих игроков с параметрами по умолчанию
        await loadSimilarPlayers(false);

        // Команды (без колонки Position)
        if (teams && Array.isArray(teams)) {
            const tbody = document.getElementById('teams-tbody');
            tbody.innerHTML = teams.map(team => `
                <tr onclick="goToTeam(${team.team_id})" style="cursor:pointer;">
                    <td style="display: flex; align-items: center; gap: 10px;">
                        ${team.url_logo ? `<img src="${team.url_logo}" alt="${team.name}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                        <strong>${team.name}</strong>
                    </td>
                    <td>${team.season || 'Н/Д'}</td>
                </tr>
            `).join('');
        }

        // Матчи с пагинацией
        fixturesCurrentPage = 1;
        await loadFixturesWithPagination(fixturesCurrentPage, false);
        
        // Проверяем избранное
        if (TokenManager.hasToken()) {
            try {
                const favorites = await playerAPI.getFavorites();
                if (favorites && Array.isArray(favorites)) {
                    isFavorite = favorites.some(f => f.athlete_id === parseInt(playerId));
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
        console.error('Ошибка загрузки данных игрока:', error);
        document.getElementById('player-detail').innerHTML = '<p class="loading">Ошибка загрузки данных игрока. Пожалуйста, попробуйте снова.</p>';
    }
}

function getCurrentSeason() {
    const currentYear = new Date().getFullYear();
    const currentMonth = new Date().getMonth();
    if (currentMonth >= 6) {
        return `${currentYear}/${currentYear + 1}`;
    } else {
        return `${currentYear - 1}/${currentYear}`;
    }
}

function switchPlayerTab(tabName) {
    const tabs = document.querySelectorAll('.tab');
    const panes = document.querySelectorAll('.tab-pane');
    
    tabs.forEach(tab => tab.classList.remove('active'));
    panes.forEach(pane => pane.classList.remove('active'));
    
    if (event && event.target) {
        event.target.classList.add('active');
    }
    
    const activePane = document.getElementById(`player-${tabName}`);
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
            await playerAPI.removeFavorite(playerId);
            isFavorite = false;
            btn.textContent = '★ Добавить в избранное';
        } else {
            await playerAPI.addFavorite(playerId);
            isFavorite = true;
            btn.textContent = '★ Удалить из избранного';
        }
    } catch (error) {
        console.error('Ошибка переключения избранного:', error);
        showFilterMessage('Не удалось обновить избранное', 'error');
    }
}

function goToMatch(id) {
    if (id && id > 0) {
        window.location.href = `/match?id=${id}`;
    }
}

function goToTeam(id) {
    if (id && id > 0) {
        window.location.href = `/team?id=${id}`;
    }
}

function goToPlayer(id) {
    if (id && id > 0) {
        window.location.href = `/player?id=${id}`;
    }
}

// Делаем функции доступными глобально для кнопки "Попробовать снова" и "Показать ещё"
window.loadSimilarPlayers = loadSimilarPlayers;
window.loadMoreFixtures = loadMoreFixtures;