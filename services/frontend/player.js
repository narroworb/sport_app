// Player page functionality
let playerId;
let isFavorite = false;
let currentStats = null;
let currentTeams = null;
let currentFixtures = null;
let availableSeasons = [];

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();

    const params = new URLSearchParams(window.location.search);
    playerId = params.get('id');

    if (playerId) {
        await loadPlayerData();
        if (TokenManager.hasToken()) {
            document.getElementById('fav-btn').style.display = 'inline-block';
        }
    }

    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
    
    // Навешиваем обработчики на фильтры
    document.getElementById('apply-filters').addEventListener('click', applyFilters);
    document.getElementById('reset-filters').addEventListener('click', resetFilters);
});

async function loadPlayerStatsWithFilters(season, dateFrom, dateTo) {
    let url = `/api/player/${playerId}/stats`;
    const params = [];
    
    if (season) {
        params.push(`season=${encodeURIComponent(season)}`);
    } else {
        if (dateFrom) params.push(`dateFrom=${dateFrom}`);
        if (dateTo) params.push(`dateTo=${dateTo}`);
    }
    
    if (params.length > 0) {
        url += `?${params.join('&')}`;
    }
    
    console.log('Fetching stats from:', url);
    const response = await fetch(url);
    if (response.ok) {
        return await response.json();
    }
    return null;
}

async function applyFilters() {
    const seasonFilter = document.getElementById('season-filter').value;
    const dateFrom = document.getElementById('date-from').value;
    const dateTo = document.getElementById('date-to').value;
    
    // Если выбран сезон, игнорируем даты
    const stats = await loadPlayerStatsWithFilters(seasonFilter || null, dateFrom, dateTo);
    
    if (stats) {
        currentStats = stats;
        updateStatsGrid(stats);
    } else {
        console.error('Failed to load filtered stats');
        document.getElementById('player-stats-grid').innerHTML = '<p>No stats available for selected period</p>';
    }
}

function resetFilters() {
    document.getElementById('season-filter').value = '';
    document.getElementById('date-from').value = '';
    document.getElementById('date-to').value = '';
    
    // Перезагружаем полную статистику
    loadPlayerStatsWithFilters(null, null, null).then(stats => {
        if (stats) {
            currentStats = stats;
            updateStatsGrid(stats);
        }
    });
}

function updateStatsGrid(stats) {
    const statsGrid = document.getElementById('player-stats-grid');
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
    
    const passAccuracy = stats.complete_passes && stats.pass_attempts 
        ? (stats.complete_passes / stats.pass_attempts) * 100 
        : null;
    
    const displayStats = [];
    for (const stat of statsToShow) {
        let value = stats[stat.key];
        if (stat.key === 'pass_accuracy' && passAccuracy !== null) {
            value = passAccuracy;
        }
        if (value !== undefined && value !== null) {
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
        
        // Сначала загружаем базовую информацию и статистику по сезонам
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

        // Собираем доступные сезоны из команд
        if (teams && Array.isArray(teams)) {
            const seasonsSet = new Set();
            teams.forEach(team => {
                if (team.season) {
                    seasonsSet.add(team.season);
                }
            });
            availableSeasons = Array.from(seasonsSet).sort().reverse();
            
            const seasonFilter = document.getElementById('season-filter');
            if (seasonFilter) {
                availableSeasons.forEach(season => {
                    const option = document.createElement('option');
                    option.value = season;
                    option.textContent = season;
                    seasonFilter.appendChild(option);
                });
            }
        }

        // Загружаем статистику (без фильтров сначала)
        const stats = await loadPlayerStatsWithFilters(null, null, null);
        if (stats) {
            currentStats = stats;
            updateStatsGrid(stats);
        }

        // Загружаем похожих игроков
        let similarPlayers = null;
        let season = availableSeasons[0] || getCurrentSeason();
        
        try {
            const similarUrl = `/api/analytics/player_similarity?player_id=${playerId}&season=${encodeURIComponent(season)}&top_k=5&min_minutes=100`;
            console.log('Fetching similar players from:', similarUrl);
            const similarResponse = await fetch(similarUrl);
            if (similarResponse.ok) {
                similarPlayers = await similarResponse.json();
                console.log('Similar players:', similarPlayers);
            }
        } catch (e) {
            console.error('Error loading similar players:', e);
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

        // Добавляем блок похожих игроков
        if (similarPlayers && similarPlayers.players && similarPlayers.players.length > 0) {
            console.log('Adding similar players block...');
            
            const similarPlayersDetails = [];
            for (const similar of similarPlayers.players.slice(0, 5)) {
                const similarDetail = await playerAPI.getDetails(similar.player_id).catch(e => null);
                if (similarDetail) {
                    similarPlayersDetails.push({ ...similar, details: similarDetail });
                }
            }
            
            const similarityHtml = `
                <div id="player-similarity" class="similar-players-section" style="margin-top: 2rem; padding: 1.5rem; background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%); border-radius: 12px; color: white;">
                    <h3 style="text-align: center; margin-bottom: 1rem;">🔄 Similar Players (${similarPlayers.position || 'N/A'})</h3>
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
                        ${similarPlayers.details || `Based on season ${season}`}
                    </div>
                </div>
            `;
            
            const statsGrid = document.getElementById('player-stats-grid');
            if (statsGrid && !document.getElementById('player-similarity')) {
                statsGrid.insertAdjacentHTML('afterend', similarityHtml);
            }
        }

        // Teams
        if (teams && Array.isArray(teams)) {
            const tbody = document.getElementById('teams-tbody');
            tbody.innerHTML = teams.map(team => `
                <tr onclick="goToTeam(${team.team_id})" style="cursor:pointer;">
                    <td style="display: flex; align-items: center; gap: 10px;">
                        ${team.url_logo ? `<img src="${team.url_logo}" alt="${team.name}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                        <strong>${team.name}</strong>
                    </td>
                    <td>${team.season || 'N/A'}</td>
                    <td>${team.position || team.role || 'N/A'}</td>
                </tr>
            `).join('');
        }

        // Fixtures
        if (fixtures && Array.isArray(fixtures)) {
            const list = document.getElementById('fixtures-list');
            list.innerHTML = fixtures.slice(0, 15).map(f => {
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
                
                return `
                    <div class="card match-card" style="cursor:pointer; margin-bottom: 1rem;" onclick="goToMatch(${matchId})">
                        <div class="match-header">
                            ${date.toLocaleDateString()} | ${tournament}
                        </div>
                        <div class="match-score" style="display: flex; justify-content: space-between; align-items: center;">
                            <div class="team-container" style="display: flex; align-items: center; gap: 10px; min-width: 120px;">
                                ${homeLogo ? `<img src="${homeLogo}" alt="${homeTeam}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                                <span class="team-name">${homeTeam}</span>
                            </div>
                            <span class="score" style="font-size: 1.2rem;">${homeScore} - ${awayScore}</span>
                            <div class="team-container" style="display: flex; align-items: center; gap: 10px; min-width: 120px;">
                                <span class="team-name">${awayTeam}</span>
                                ${awayLogo ? `<img src="${awayLogo}" alt="${awayTeam}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                            </div>
                        </div>
                        <div class="match-status" style="display: flex; justify-content: space-between; margin-top: 10px; padding-top: 10px; border-top: 1px solid #eee;">
                            <span>${status}</span>
                            <span>⚽ ${goals} goals | 🎯 ${assists} assists | ⭐ ${rating} | ⏱️ ${minutes}'</span>
                        </div>
                    </div>
                `;
            }).join('');
        }
        
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
    }
}

function goToMatch(id) {
    if (id && id > 0) {
        window.location.href = `/match.html?id=${id}`;
    }
}

function goToTeam(id) {
    if (id && id > 0) {
        window.location.href = `/team.html?id=${id}`;
    }
}

function goToPlayer(id) {
    if (id && id > 0) {
        window.location.href = `/player.html?id=${id}`;
    }
}