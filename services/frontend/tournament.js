// Tournament page functionality
let tournamentId;
let positionHistoryData = null;
let positionChart = null;

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();

    const params = new URLSearchParams(window.location.search);
    tournamentId = params.get('id');

    if (tournamentId) {
        await loadTournamentData();
        if (TokenManager.hasToken()) {
            document.getElementById('fav-btn').style.display = 'inline-block';
        }
    }

    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
});

async function loadTournamentData() {
    try {
        console.log('Loading tournament data for ID:', tournamentId);
        
        const [details, table, teamsStats, playersStats, fixtures, positionHistory] = await Promise.all([
            tournamentAPI.getDetails(tournamentId).catch(e => {
                console.error('Details error:', e);
                return null;
            }),
            tournamentAPI.getTable(tournamentId).catch(e => {
                console.error('Table error:', e);
                return null;
            }),
            tournamentAPI.getTeamsStats(tournamentId).catch(e => {
                console.error('Teams stats error:', e);
                return null;
            }),
            tournamentAPI.getPlayersStats(tournamentId).catch(e => {
                console.error('Players stats error:', e);
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
        console.log('Teams stats:', teamsStats);
        console.log('Players stats:', playersStats);
        console.log('Fixtures:', fixtures);
        console.log('Position history:', positionHistory);

        if (!details) {
            document.getElementById('tournament-detail').innerHTML = '<p class="loading">Tournament not found</p>';
            return;
        }

        // Сохраняем данные истории позиций
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
        } else {
            document.getElementById('table-tbody').innerHTML = '<tr><td colspan="9">No standings available</td></tr>';
        }

        // Teams Stats (статистика команд)
        if (teamsStats && Array.isArray(teamsStats)) {
            const tbody = document.getElementById('teams-stats-tbody');
            const teamDetailsMap = new Map();
            for (const stat of teamsStats.slice(0, 10)) {
                const teamDetail = await teamAPI.getDetails(stat.team_id).catch(() => null);
                if (teamDetail) {
                    teamDetailsMap.set(stat.team_id, teamDetail);
                }
            }
            
            tbody.innerHTML = teamsStats.slice(0, 20).map(stat => {
                const teamDetail = teamDetailsMap.get(stat.team_id);
                const teamName = teamDetail?.name || `Team ${stat.team_id}`;
                const teamLogo = teamDetail?.url_logo || '';
                const possession = stat.average_ball_possession ? stat.average_ball_possession.toFixed(1) : 0;
                const goalsPer90 = stat.goals_per_90 ? stat.goals_per_90.toFixed(2) : 0;
                const goalsConcededPer90 = stat.goals_conceded_per_90 ? stat.goals_conceded_per_90.toFixed(2) : 0;
                
                return `
                    <tr onclick="goToTeam(${stat.team_id})" style="cursor:pointer;">
                        <td style="display: flex; align-items: center; gap: 10px;">
                            ${teamLogo ? `<img src="${teamLogo}" alt="${teamName}" style="height: 25px; width: 25px; object-fit: contain;">` : ''}
                            ${teamName}
                        </td>
                        <td>${stat.goals || 0}</td>
                        <td>${stat.total_shots || 0}</td>
                        <td>${possession}%</td>
                        <td>${goalsPer90}</td>
                        <td>${goalsConcededPer90}</td>
                    </tr>
                `;
            }).join('');
        } else {
            document.getElementById('teams-stats-tbody').innerHTML = '<tr><td colspan="6">No team stats available</td></tr>';
        }

        // Players Stats (статистика игроков)
        if (playersStats && Array.isArray(playersStats)) {
            const tbody = document.getElementById('players-stats-tbody');
            const playerDetailsMap = new Map();
            
            for (const stat of playersStats.slice(0, 20)) {
                const playerDetail = await playerAPI.getDetails(stat.athlete_id).catch(() => null);
                if (playerDetail) {
                    playerDetailsMap.set(stat.athlete_id, playerDetail);
                }
            }
            
            tbody.innerHTML = playersStats.slice(0, 30).map(stat => {
                const playerDetail = playerDetailsMap.get(stat.athlete_id);
                const playerName = playerDetail?.first_name && playerDetail?.last_name 
                    ? `${playerDetail.first_name} ${playerDetail.last_name}`
                    : `Player ${stat.athlete_id}`;
                const playerPhoto = playerDetail?.url_photo || '';
                const rating = stat.avg_rating ? stat.avg_rating.toFixed(1) : 0;
                const goalsPer90 = stat.goals_per_90 ? stat.goals_per_90.toFixed(2) : 0;
                const assistsPer90 = stat.assists_per_90 ? stat.assists_per_90.toFixed(2) : 0;
                
                return `
                    <tr onclick="goToPlayer(${stat.athlete_id})" style="cursor:pointer;">
                        <td style="display: flex; align-items: center; gap: 10px;">
                            ${playerPhoto ? `<img src="${playerPhoto}" alt="${playerName}" style="width: 30px; height: 30px; border-radius: 50%; object-fit: cover;">` : '👤'}
                            <strong>${playerName}</strong>
                        </td>
                        <td>${stat.team_name || 'N/A'}</td>
                        <td>${stat.goals || 0}</td>
                        <td>${stat.assists || 0}</td>
                        <td>${rating}</td>
                        <td>${goalsPer90}</td>
                        <td>${assistsPer90}</td>
                    </tr>
                `;
            }).join('');
        } else {
            document.getElementById('players-stats-tbody').innerHTML = '<tr><td colspan="7">No player stats available</td></tr>';
        }

        // Fixtures (матчи турнира)
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
            
            if (allMatches.length === 0) {
                list.innerHTML = '<p>No fixtures available</p>';
            }
        } else {
            document.getElementById('fixtures-list').innerHTML = '<p>No fixtures available</p>';
        }

        // Инициализируем график истории позиций
        if (positionHistoryData && Array.isArray(positionHistoryData) && positionHistoryData.length > 0) {
            console.log('Initializing position chart with data:', positionHistoryData);
            initPositionChart(positionHistoryData);
        } else {
            console.log('No position history data available');
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
    
    // Подготавливаем данные для графика
    const rounds = historyData.map((_, index) => `Round ${index + 1}`);
    console.log('Rounds:', rounds);
    
    // Собираем все команды и их позиции по турам
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
    
    console.log('Teams found:', teamsMap.size);
    const teams = Array.from(teamsMap.values()).sort((a, b) => a.name.localeCompare(b.name));
    
    // Заполняем селектор выбора команд
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
    
    // Сохраняем данные для глобального доступа
    window.positionChartData = {
        rounds: rounds,
        teams: teams
    };
    
    console.log('Chart data prepared, teams count:', teams.length);
    
    // Создаем график со всеми командами
    updatePositionChart();
}

function updatePositionChart() {
    console.log('updatePositionChart called');
    
    if (!window.positionChartData) {
        console.log('No positionChartData');
        return;
    }
    
    const teamSelector = document.getElementById('team-selector');
    if (!teamSelector) {
        console.log('team-selector not found');
        return;
    }
    
    const selectedValues = Array.from(teamSelector.selectedOptions).map(opt => opt.value);
    const showAll = selectedValues.includes('all') || selectedValues.length === 0;
    
    // Фильтруем команды
    let teamsToShow = window.positionChartData.teams;
    if (!showAll) {
        teamsToShow = window.positionChartData.teams.filter(team => 
            selectedValues.includes(team.name)
        );
    }
    
    console.log('Teams to show:', teamsToShow.length);
    
    // Подготавливаем данные для Chart.js
    const datasets = teamsToShow.map(team => {
        // Генерируем стабильный цвет для команды
        const hue = (team.name.length * 37) % 360;
        const color = `hsl(${hue}, 70%, 55%)`;
        
        // Фильтруем null значения (команда не играла в этом туре)
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
    
    if (datasets.length === 0) {
        console.log('No datasets to show');
        return;
    }
    
    // Если нет Chart.js, загружаем его
    if (typeof Chart === 'undefined') {
        console.log('Loading Chart.js...');
        const script = document.createElement('script');
        script.src = 'https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js';
        script.onload = () => {
            console.log('Chart.js loaded, creating chart...');
            createChart(datasets);
        };
        script.onerror = () => {
            console.error('Failed to load Chart.js');
        };
        document.head.appendChild(script);
    } else {
        console.log('Chart.js already loaded, creating chart...');
        createChart(datasets);
    }
}

function createChart(datasets) {
    const canvas = document.getElementById('position-chart');
    if (!canvas) {
        console.log('Canvas element not found');
        return;
    }
    
    const ctx = canvas.getContext('2d');
    
    if (positionChart) {
        positionChart.destroy();
    }
    
    // Находим максимальное количество команд для определения max оси Y
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
                    labels: {
                        font: { size: 10 },
                        boxWidth: 12
                    }
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
                    title: {
                        display: true,
                        text: 'Position',
                        font: { weight: 'bold' }
                    },
                    min: 1,
                    max: maxTeams,
                    ticks: {
                        stepSize: 1,
                        callback: function(value) {
                            return value;
                        }
                    }
                },
                x: {
                    title: {
                        display: true,
                        text: 'Round',
                        font: { weight: 'bold' }
                    }
                }
            }
        }
    });
    
    console.log('Chart created successfully');
}

function resetChartSelection() {
    const teamSelector = document.getElementById('team-selector');
    if (teamSelector) {
        // Снимаем все выделения
        Array.from(teamSelector.options).forEach(opt => {
            opt.selected = false;
        });
        // Выбираем "Show All"
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
    } else {
        tabs.forEach(tab => {
            if (tab.textContent.toLowerCase().includes(tabName)) {
                tab.classList.add('active');
            }
        });
    }
    
    const activePane = document.getElementById(`tournament-${tabName}`);
    if (activePane) {
        activePane.classList.add('active');
    }
    
    // Обновляем график при переключении на вкладку истории
    if (tabName === 'history' && positionChart) {
        setTimeout(() => {
            if (positionChart) positionChart.resize();
        }, 100);
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