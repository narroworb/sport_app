// Player page functionality
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

// Пагинация для fixtures
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
    seasonSelect.innerHTML = '<option value="">Auto-detect</option>';
    
    // Добавляем сезоны от 2015/2016 до 2025/2026
    for (let year = 2025; year >= 2015; year--) {
        const season = `${year}/${year + 1}`;
        const option = document.createElement('option');
        option.value = season;
        option.textContent = season;
        seasonSelect.appendChild(option);
    }
}

// Функция загрузки fixtures с пагинацией
async function loadFixturesWithPagination(page = 1, append = false) {
    if (!playerId || isLoadingFixtures) return;
    
    isLoadingFixtures = true;
    const offset = (page - 1) * fixturesLimit + 1;
    
    try {
        const url = `/api/player/${playerId}/fixtures?limit=${fixturesLimit}&offset=${offset}`;
        console.log('Loading fixtures from:', url);
        
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error(`Failed to load fixtures: ${response.status}`);
        }
        
        const fixtures = await response.json();
        console.log('Fixtures loaded:', fixtures);
        
        // Обновляем общее количество (если API возвращает total, иначе используем длину массива)
        if (fixtures.length < fixturesLimit) {
            fixturesTotalCount = offset + fixtures.length;
        } else {
            fixturesTotalCount = offset + fixturesLimit + 1; // Приблизительное значение
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

// Функция отрисовки списка fixtures
function renderFixturesList(fixtures, append = false) {
    const list = document.getElementById('fixtures-list');
    if (!list) return;
    
    let html = '';
    
    if (append && list.innerHTML !== '<p>No fixtures available</p>') {
        html = list.innerHTML;
    } else {
        html = '';
    }
    
    if (!fixtures || fixtures.length === 0) {
        if (!append) {
            list.innerHTML = '<p>No fixtures available</p>';
        }
        return;
    }
    
    html += fixtures.map(f => {
        const homeTeam = f.home_team?.name || 'Home';
        const awayTeam = f.away_team?.name || 'Away';
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
                <div class="match-status ${statusClass}" style="display: flex; justify-content: space-between; margin-top: 10px; padding-top: 10px; border-top: 1px solid rgba(255,255,255,0.2);">
                    <span>${statusText}</span>
                    <span>⚽ ${goals} goals | 🎯 ${assists} assists | ⭐ ${rating} | ⏱️ ${minutes}'</span>
                </div>
            </div>
        `;
    }).join('');
    
    // Добавляем кнопку "Load More" если есть еще данные
    const currentCount = (fixturesCurrentPage + 1) * fixturesLimit;
    if (currentCount < fixturesTotalCount && !append) {
        html += `
            <div style="text-align: center; margin-top: 1rem;">
                <button id="load-more-fixtures" class="btn-secondary" onclick="loadMoreFixtures()">Load More ↓</button>
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
    
    console.log('Loading similar players with params:', currentSimilarParams);
    
    // Показываем индикатор загрузки
    const statsGrid = document.getElementById('player-stats-grid');
    const existingBlock = document.getElementById('player-similarity');
    
    if (showLoading) {
        if (existingBlock) {
            existingBlock.innerHTML = '<div style="text-align: center; padding: 2rem;">🔄 Loading similar players...</div>';
        } else if (statsGrid) {
            statsGrid.insertAdjacentHTML('afterend', '<div id="temp-similar-loading" style="margin-top: 2rem; padding: 1.5rem; background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%); border-radius: 12px; color: white; text-align: center;">🔄 Loading similar players...</div>');
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
            throw new Error(`API error: ${response.status}`);
        }
        
        const similarPlayers = await response.json();
        console.log('Similar players received:', similarPlayers);
        
        // Обновляем блок похожих игроков
        await updateSimilarPlayersBlock(similarPlayers, season);
        
    } catch (error) {
        console.error('Error loading similar players:', error);
        showFilterMessage('Failed to load similar players: ' + error.message, 'error');
        
        // Показываем ошибку в блоке
        const existingBlock = document.getElementById('player-similarity');
        if (existingBlock) {
            existingBlock.innerHTML = `
                <div style="text-align: center; padding: 2rem;">
                    <p>❌ Failed to load similar players</p>
                    <p style="font-size: 0.8rem; opacity: 0.7;">${error.message}</p>
                    <button onclick="loadSimilarPlayers(true)" class="btn-primary" style="margin-top: 1rem;">🔄 Try Again</button>
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
                    <p>🔍 No similar players found</p>
                    <p style="font-size: 0.8rem; opacity: 0.7;">Try adjusting search parameters</p>
                </div>
            `;
        }
        return;
    }
    
    // Загружаем детали похожих игроков
    const similarPlayersDetails = [];
    for (const similar of similarPlayers.players.slice(0, currentSimilarParams.top_k)) {
        const similarDetail = await playerAPI.getDetails(similar.player_id).catch(e => {
            console.error(`Error loading similar player ${similar.player_id}:`, e);
            return null;
        });
        if (similarDetail) {
            similarPlayersDetails.push({ ...similar, details: similarDetail });
        }
    }
    
    const similarityHtml = `
        <div id="player-similarity" class="similar-players-section" style="margin-top: 2rem; padding: 1.5rem; background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%); border-radius: 12px; color: white;">
            <div style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; margin-bottom: 1rem;">
                <h3 style="margin: 0;">🔄 Similar Players (${similarPlayers.position || 'N/A'})</h3>
                <div style="font-size: 0.7rem; opacity: 0.7;">
                    ⚙️ Top ${currentSimilarParams.top_k} | Min ${currentSimilarParams.min_minutes} min | Season ${season}
                </div>
            </div>
            <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1rem;">
                ${similarPlayersDetails.map((player, index) => {
                    const similarityPercent = (player.similarity * 100).toFixed(1);
                    const playerName = `${player.details.first_name || ''} ${player.details.last_name || ''}`.trim() || 'Unknown';
                    const playerPhoto = player.details.url_photo || '';
                    const playerPosition = player.details.position || 'N/A';
                    const positionNames = { 'G': 'Goalkeeper', 'D': 'Defender', 'M': 'Midfielder', 'F': 'Forward' };
                    const positionDisplay = positionNames[playerPosition] || playerPosition;
                    
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
                                    <div style="font-size: 0.8rem; opacity: 0.9;">${positionDisplay}</div>
                                    <div style="margin-top: 5px;">
                                        <div style="background: rgba(255,255,255,0.2); border-radius: 10px; height: 6px; width: 100%;">
                                            <div style="background: ${similarityColor}; width: ${similarityPercent}%; height: 6px; border-radius: 10px;"></div>
                                        </div>
                                        <div style="font-size: 0.8rem; margin-top: 3px;">Match: ${similarityPercent}%</div>
                                    </div>
                                </div>
                                <div style="font-size: 1.2rem;">#${index + 1}</div>
                            </div>
                        </div>
                    `;
                }).join('')}
            </div>
            <div style="margin-top: 1rem; font-size: 0.7rem; text-align: center; opacity: 0.7;">
                ${similarPlayers.details || `Based on season ${season} | Candidates: ${similarPlayers.candidates || 'N/A'}`}
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
    
    showFilterMessage(`Found ${similarPlayersDetails.length} similar players`, 'success');
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
    
    console.log('Fetching stats from:', url);
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

async function applySeasonFilter() {
    const seasonFilter = document.getElementById('season-filter').value;
    
    if (!seasonFilter) {
        alert('Please select a season');
        return;
    }
    
    console.log('Applying season filter:', seasonFilter);
    
    const statsGrid = document.getElementById('player-stats-grid');
    statsGrid.innerHTML = '<div class="loading-spinner">Loading statistics...</div>';
    
    const stats = await loadPlayerStatsWithFilters(seasonFilter, null, null);
    
    if (stats && Object.keys(stats).length > 0) {
        currentStats = stats;
        updateStatsGrid(stats);
        showFilterMessage(`Showing stats for season ${seasonFilter}`);
    } else {
        statsGrid.innerHTML = '<p>No stats available for selected season</p>';
        showFilterMessage(`No data found for season ${seasonFilter}`, 'error');
    }
}

async function applyDateFilter() {
    const dateFrom = document.getElementById('date-from').value;
    const dateTo = document.getElementById('date-to').value;
    
    if (!dateFrom && !dateTo) {
        alert('Please select at least one date');
        return;
    }
    
    console.log('Applying date filter:', { dateFrom, dateTo });
    
    const statsGrid = document.getElementById('player-stats-grid');
    statsGrid.innerHTML = '<div class="loading-spinner">Loading statistics...</div>';
    
    const stats = await loadPlayerStatsWithFilters(null, dateFrom, dateTo);
    
    if (stats && Object.keys(stats).length > 0) {
        currentStats = stats;
        updateStatsGrid(stats);
        let message = 'Showing stats';
        if (dateFrom && dateTo) {
            message += ` from ${dateFrom} to ${dateTo}`;
        } else if (dateFrom) {
            message += ` from ${dateFrom}`;
        } else if (dateTo) {
            message += ` until ${dateTo}`;
        }
        showFilterMessage(message);
    } else {
        statsGrid.innerHTML = '<p>No stats available for selected date range</p>';
        showFilterMessage('No data found for selected date range', 'error');
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
    statsGrid.innerHTML = '<div class="loading-spinner">Loading statistics...</div>';
    
    loadPlayerStatsWithFilters(null, null, null).then(stats => {
        if (stats && Object.keys(stats).length > 0) {
            currentStats = stats;
            updateStatsGrid(stats);
            showFilterMessage('Showing all-time statistics');
        } else {
            statsGrid.innerHTML = '<p>No stats available</p>';
        }
    });
}

function updateStatsGrid(stats) {
    const statsGrid = document.getElementById('player-stats-grid');
    
    if (!stats || Object.keys(stats).length === 0) {
        statsGrid.innerHTML = '<p>No statistics available for this period</p>';
        return;
    }
    
    const statsToShow = [
        { key: 'total_matches', label: 'Matches', isPercent: false },
        { key: 'goals', label: 'Goals', isPercent: false },
        { key: 'assists', label: 'Assists', isPercent: false },
        { key: 'avg_rating', label: 'Avg Rating', isPercent: false },
        { key: 'minutes_played', label: 'Minutes', isPercent: false },
        { key: 'yellow_cards', label: 'Yellow Cards', isPercent: false },
        { key: 'red_cards', label: 'Red Cards', isPercent: false },
        { key: 'total_tackles', label: 'Tackles', isPercent: false },
        { key: 'interceptions', label: 'Interceptions', isPercent: false },
        { key: 'saves', label: 'Saves', isPercent: false },
        { key: 'pass_accuracy', label: 'Pass Accuracy', isPercent: true },
        { key: 'dribble_accuracy', label: 'Dribble Accuracy', isPercent: true }
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
        displayStats.push({ label: 'Goals/90', value: goalsPer90.toFixed(2) });
    }
    
    if (stats.saves && stats.minutes_played && stats.minutes_played > 0) {
        const savesPer90 = (stats.saves / stats.minutes_played) * 90;
        displayStats.push({ label: 'Saves/90', value: savesPer90.toFixed(2) });
    }
    
    if (displayStats.length > 0) {
        statsGrid.innerHTML = displayStats.map(stat => `
            <div class="stat-box">
                <div class="stat-value">${stat.value}</div>
                <div class="stat-label">${stat.label}</div>
            </div>
        `).join('');
    } else {
        statsGrid.innerHTML = '<p>No stats available</p>';
    }
}

async function loadPlayerData() {
    try {
        console.log('Loading player data for ID:', playerId);
        
        const [details, teams, fixtures] = await Promise.all([
            playerAPI.getDetails(playerId).catch(e => {
                console.error('Details error:', e);
                return null;
            }),
            playerAPI.getTeams(playerId).catch(e => {
                console.error('Teams error:', e);
                return null;
            }),
            playerAPI.getFixtures(playerId).catch(e => {
                console.error('Fixtures error:', e);
                return null;
            }),
        ]);

        console.log('Player details:', details);
        console.log('Player teams:', teams);
        console.log('Player fixtures:', fixtures);

        if (!details) {
            document.getElementById('player-detail').innerHTML = '<p class="loading">Player not found</p>';
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
            console.log('Seasons loaded:', availableSeasons.length);
        }
        
        // Заполняем сезоны для выбора в настройках похожих игроков
        populateSimilarSeasons();

        const stats = await loadPlayerStatsWithFilters(null, null, null);
        if (stats) {
            currentStats = stats;
            updateStatsGrid(stats);
        }

        // Update header with photo
        const playerName = `${details.first_name || ''} ${details.last_name || ''}`.trim() || 'Player';
        const photo = details.url_photo || '';
        const nationFlag = details.nation?.url_flag || '';
        const nationName = details.nation?.name || '';

        document.getElementById('player-name').innerHTML = `
            ${photo ? `<img src="${photo}" alt="${playerName}" style="height: 60px; border-radius: 50%; vertical-align: middle; margin-right: 15px;">` : ''}
            ${playerName}
        `;

        const positionSpan = document.getElementById('player-position');
        if (details.position) {
            const positionNames = { 'G': 'Goalkeeper', 'D': 'Defender', 'M': 'Midfielder', 'F': 'Forward' };
            positionSpan.textContent = positionNames[details.position] || details.position;
        }

        const nationSpan = document.getElementById('player-nation');
        if (nationName) {
            nationSpan.innerHTML = `${nationFlag ? `<img src="${nationFlag}" style="width: 20px; vertical-align: middle; margin-right: 5px;">` : '🌍'} ${nationName}`;
        }

        // Дополнительная информация
        if (details.date_of_birth) {
            const birthDate = new Date(details.date_of_birth);
            const formattedDate = birthDate.toLocaleDateString('en-GB', {
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            });
            const ageDifMs = Date.now() - birthDate.getTime();
            const ageDate = new Date(ageDifMs);
            const age = Math.abs(ageDate.getUTCFullYear() - 1970);
            document.getElementById('player-birthdate').innerHTML = `${formattedDate} <span style="color: #666;">(${age} years)</span>`;
        } else {
            document.getElementById('player-birthdate').textContent = 'Not specified';
        }

        if (details.height) {
            const heightCm = details.height;
            const heightFeet = Math.floor(heightCm / 30.48);
            const heightInches = Math.round((heightCm % 30.48) / 2.54);
            document.getElementById('player-height').innerHTML = `${heightCm} cm (${heightFeet}'${heightInches}")`;
        } else {
            document.getElementById('player-height').textContent = 'Not specified';
        }

        if (details.preffered_foot) {
            const footNames = { 'Left': '👈 Left', 'Right': '👉 Right', 'Both': '✌️ Both' };
            document.getElementById('player-foot').innerHTML = footNames[details.preffered_foot] || details.preffered_foot;
        } else {
            document.getElementById('player-foot').textContent = 'Not specified';
        }

        if (details.current_status) {
            const statusBadge = details.current_status === 'Active' 
                ? '<span style="color: #2ecc71;">✅ Active</span>' 
                : '<span style="color: #e74c3c;">❌ Inactive</span>';
            document.getElementById('player-status').innerHTML = statusBadge;
        } else {
            document.getElementById('player-status').textContent = 'Unknown';
        }

        // Загружаем похожих игроков с параметрами по умолчанию
        await loadSimilarPlayers(false);

        // Teams (без колонки Position)
        if (teams && Array.isArray(teams)) {
            const tbody = document.getElementById('teams-tbody');
            tbody.innerHTML = teams.map(team => `
                <tr onclick="goToTeam(${team.team_id})" style="cursor:pointer;">
                    <td style="display: flex; align-items: center; gap: 10px;">
                        ${team.url_logo ? `<img src="${team.url_logo}" alt="${team.name}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                        <strong>${team.name}</strong>
                    </td>
                    <td>${team.season || 'N/A'}</td>
                </tr>
            `).join('');
        }

        // Fixtures с пагинацией
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
                        btn.textContent = isFavorite ? '★ Remove from Favorites' : '★ Add to Favorites';
                    }
                }
            } catch (e) {
                console.error('Error checking favorites:', e);
            }
        }
        
    } catch (error) {
        console.error('Error loading player data:', error);
        document.getElementById('player-detail').innerHTML = '<p class="loading">Error loading player data. Please try again.</p>';
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
            btn.textContent = '★ Add to Favorites';
        } else {
            await playerAPI.addFavorite(playerId);
            isFavorite = true;
            btn.textContent = '★ Remove from Favorites';
        }
    } catch (error) {
        console.error('Error toggling favorite:', error);
        showFilterMessage('Failed to update favorites', 'error');
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

// Делаем функции доступными глобально для кнопки Try Again и Load More
window.loadSimilarPlayers = loadSimilarPlayers;
window.loadMoreFixtures = loadMoreFixtures;