// Team page functionality
let teamId;
let currentTeamData = {};

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();

    const params = new URLSearchParams(window.location.search);
    teamId = params.get('id');

    if (teamId) {
        await loadTeamData();
        if (TokenManager.hasToken()) {
            document.getElementById('fav-btn').style.display = 'inline-block';
        }
    }

    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
});

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
        
        const players = await teamAPI.getPlayers(teamId).catch(e => {
            console.error('Players error:', e);
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
        
        // Загружаем аналитику формы команды
        let teamForm = null;
        
        // Пробуем получить сезон из разных источников
        let season = null;
        
        // 1. Из турнира команды
        if (details && details.tournament && details.tournament.season) {
            season = details.tournament.season;
        }
        
        // 2. Из первого матча в расписании (если есть)
        if (!season && fixtures && fixtures.length > 0) {
            const firstFixture = fixtures[0];
            const tournament = firstFixture.tournament || firstFixture.data?.tournament;
            if (tournament && tournament.season) {
                season = tournament.season;
                console.log('Season from fixture:', season);
            }
        }
        
        // 3. Из следующего матча
        if (!season && nextGame && nextGame.tournament && nextGame.tournament.season) {
            season = nextGame.tournament.season;
            console.log('Season from nextGame:', season);
        }
        
        // 4. Из турнирной таблицы (если есть)
        if (!season && standings && standings.length > 0 && standings[0].season) {
            season = standings[0].season;
            console.log('Season from standings:', season);
        }
        
        // 5. Текущий сезон по умолчанию (2025/2026 или определяем по дате)
        if (!season) {
            const currentYear = new Date().getFullYear();
            const nextYear = currentYear + 1;
            // Определяем сезон: если сейчас середина/конец года, используем currentYear/currentYear+1
            const currentMonth = new Date().getMonth();
            if (currentMonth >= 6) { // Июль и позже
                season = `${currentYear}/${currentYear + 1}`;
            } else {
                season = `${currentYear - 1}/${currentYear}`;
            }
            console.log('Using default season:', season);
        }
        
        if (teamId && season) {
            try {
                const formUrl = `/api/analytics/team_form?team_id=${teamId}&season=${encodeURIComponent(season)}&matches_back=10&half_life_matches=5`;
                console.log('Fetching team form from:', formUrl);
                const formResponse = await fetch(formUrl);
                if (formResponse.ok) {
                    teamForm = await formResponse.json();
                    console.log('Team form analytics:', teamForm);
                } else {
                    console.error('Team form response not OK:', formResponse.status);
                }
            } catch (e) {
                console.error('Error loading team form:', e);
            }
        } else {
            console.log('Cannot fetch team form: missing teamId or season', { teamId, season });
        }

        console.log('Team details:', details);
        console.log('Team stats:', stats);
        console.log('Next game:', nextGame);
        console.log('Standings:', standings);
        console.log('Players (raw):', players);
        console.log('Manager:', manager);
        console.log('Fixtures:', fixtures);
        console.log('Team form:', teamForm);

        if (!details) {
            document.getElementById('team-detail').innerHTML = '<p class="loading">Team not found</p>';
            return;
        }

        currentTeamData = details;

        // Update header with logo if available
        const teamName = details.name || details.team_name || 'Team';
        document.getElementById('team-name').innerHTML = `
            ${details.url_logo ? `<img src="${details.url_logo}" alt="${teamName}" style="height: 40px; vertical-align: middle; margin-right: 10px;">` : ''}
            ${teamName}
        `;

        // Stats grid with rounded values
        if (stats) {
            const statsGrid = document.getElementById('team-stats-grid');
            const statsData = stats.data || stats;
            const statsToShow = Object.entries(statsData).slice(0, 8);
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

        // Добавляем блок аналитики формы (в overview после статистики)
        if (teamForm && teamForm.form_index !== undefined) {
            console.log('Adding team form analytics block...');
            
            const formIndex = (teamForm.form_index * 100).toFixed(1);
            const attackIndex = (teamForm.attack_index * 100).toFixed(1);
            const defenseIndex = (teamForm.defense_index * 100).toFixed(1);
            const confidence = (teamForm.confidence * 100).toFixed(1);
            const trend = teamForm.trend || 0;
            
            // Определяем направление тренда
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
            
            // Определяем оценку формы
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
            
            // Находим элемент team-stats-grid и добавляем после него
            const statsGrid = document.getElementById('team-stats-grid');
            if (statsGrid) {
                // Проверяем, нет ли уже такого блока
                if (!document.getElementById('team-form-analytics')) {
                    statsGrid.insertAdjacentHTML('afterend', formHtml);
                    console.log('Team form block added successfully');
                } else {
                    console.log('Team form block already exists');
                }
            } else {
                console.log('statsGrid element not found');
                // Альтернативный вариант - добавить в начало overview
                const overviewDiv = document.getElementById('team-overview');
                if (overviewDiv) {
                    overviewDiv.insertAdjacentHTML('afterbegin', formHtml);
                    console.log('Team form block added to overview');
                }
            }
        } else {
            console.log('No team form data available:', teamForm);
            // Показываем сообщение, что данные формы недоступны
            const statsGrid = document.getElementById('team-stats-grid');
            if (statsGrid && !document.getElementById('team-form-analytics')) {
                const formUnavailableHtml = `
                    <div id="team-form-analytics" class="form-analytics-section" style="margin-top: 2rem; padding: 1.5rem; background: linear-gradient(135deg, #7f8c8d 0%, #95a5a6 100%); border-radius: 12px; color: white;">
                        <h3 style="text-align: center; margin-bottom: 1rem;">📊 Team Form Analytics</h3>
                        <div style="text-align: center;">
                            <p>Form data not available for this team</p>
                            <p style="font-size: 0.8rem; opacity: 0.8;">Not enough matches played or data processing in progress</p>
                        </div>
                    </div>
                `;
                statsGrid.insertAdjacentHTML('afterend', formUnavailableHtml);
                console.log('Added form unavailable message');
            }
        }

        // Next game with team logos
        if (nextGame && Object.keys(nextGame).length > 0) {
            const game = Array.isArray(nextGame) ? nextGame[0] : nextGame;
            const gameData = game.data || game;
            const homeTeam = gameData.home_team?.name || gameData.home_team_name || 'Home';
            const awayTeam = gameData.away_team?.name || gameData.away_team_name || 'Away';
            const homeLogo = gameData.home_team?.url_logo || '';
            const awayLogo = gameData.away_team?.url_logo || '';
            const matchId = gameData.match_id || gameData.id;
            const date = gameData.date ? new Date(gameData.date) : new Date();
            
            document.getElementById('next-game').innerHTML = `
                <div class="match-card" style="cursor:pointer;" onclick="goToMatch(${matchId})">
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

        // Standings with team logos
        if (standings && Array.isArray(standings)) {
            const tbody = document.getElementById('standing-tbody');
            tbody.innerHTML = standings.map(standing => {
                const teamData = standing.team || standing;
                const teamName = teamData.name || teamData.team_name || 'Team';
                const teamIdStanding = teamData.team_id || teamData.id;
                const teamLogo = teamData.url_logo || '';
                
                return `
                    <tr onclick="goToTeam(${teamIdStanding})" style="cursor:pointer;">
                        <td><strong>${standing.position || standing.pos || 'N/A'}</strong></td>
                        <td style="display: flex; align-items: center; gap: 10px;">
                            ${teamLogo ? `<img src="${teamLogo}" alt="${teamName}" style="height: 30px; width: 30px; object-fit: contain;">` : ''}
                            ${teamName}
                        </td>
                        <td>${standing.points || 0}</td>
                        <td>${standing.matches_played || standing.played || 0}</td>
                        <td>${standing.wins || 0}/${standing.draws || 0}/${standing.losses || 0}</td>
                      </tr>
                `;
            }).join('');
        } else {
            document.getElementById('standing-tbody').innerHTML = '<tr><td colspan="5">No standings available</tr>';
        }

        // Players
        if (players && typeof players === 'object') {
            const tbody = document.getElementById('players-tbody');
            
            const allPlayers = [];
            const positions = ['G', 'D', 'M', 'F'];
            const positionNames = { 'G': 'Goalkeeper', 'D': 'Defender', 'M': 'Midfielder', 'F': 'Forward' };
            
            positions.forEach(pos => {
                if (players[pos] && Array.isArray(players[pos])) {
                    players[pos].forEach(player => {
                        allPlayers.push({
                            ...player,
                            position_display: positionNames[pos],
                            position_code: pos
                        });
                    });
                }
            });
            
            console.log('Processed players:', allPlayers);
            
            if (allPlayers.length > 0) {
                tbody.innerHTML = allPlayers.map(player => {
                    const name = `${player.first_name || ''} ${player.last_name || ''}`.trim() || 'Unknown';
                    const position = player.position_display || player.position || 'N/A';
                    const number = player.number || player.jersey_number || 'N/A';
                    const playerId = player.athlete_id || player.id;
                    const photo = player.url_photo || '';
                    const nation = player.nation?.name || '';
                    const flag = player.nation?.url_flag || '';
                    
                    return `
                        <tr onclick="goToPlayer(${playerId})" style="cursor:pointer;">
                            <td style="display: flex; align-items: center; gap: 10px;">
                                ${photo ? `<img src="${photo}" alt="${name}" style="width: 40px; height: 40px; border-radius: 50%; object-fit: cover;">` : 
                                          `<div style="width: 40px; height: 40px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center;">📷</div>`}
                                <div>
                                    <strong>${name}</strong>
                                    ${nation ? `<div style="font-size: 0.8rem; color: #666;">${flag ? `<img src="${flag}" style="width: 20px; vertical-align: middle;">` : ''} ${nation}</div>` : ''}
                                </div>
                              </td>
                              <td>${position}</td>
                              <td>${number}</td>
                          </tr>
                    `;
                }).join('');
            } else {
                tbody.innerHTML = '<tr><td colspan="3">No players available</td></tr>';
            }
        } else {
            document.getElementById('players-tbody').innerHTML = '<tr><td colspan="3">No players available</td></tr>';
        }

        // Manager
        if (manager) {
            const mgr = Array.isArray(manager) ? manager[0] : manager;
            const managerData = mgr.data || mgr;
            const firstName = managerData.first_name || '';
            const lastName = managerData.last_name || '';
            const name = `${firstName} ${lastName}`.trim() || managerData.name || 'Unknown';
            const managerId = managerData.manager_id || managerData.id;
            const nation = managerData.nation?.name || managerData.nation || 'Unknown';
            const flag = managerData.nation?.url_flag || '';
            
            document.getElementById('manager-card').innerHTML = `
                <div class="card" style="cursor:pointer;" onclick="goToManager(${managerId})">
                    <h3>${name}</h3>
                    <p>Manager</p>
                    <p class="result-info">${flag ? `<img src="${flag}" style="width: 20px; vertical-align: middle;">` : ''} ${nation}</p>
                </div>
            `;
        } else {
            document.getElementById('manager-card').innerHTML = '<p>No manager information available</p>';
        }

        // Fixtures with team logos
        if (fixtures && Array.isArray(fixtures)) {
            const list = document.getElementById('fixtures-list');
            list.innerHTML = fixtures.slice(0, 10).map(f => {
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
                
                return `
                    <div class="card match-card" style="cursor:pointer;" onclick="goToMatch(${matchId})">
                        <div class="match-header">${date.toLocaleDateString()}</div>
                        <div class="match-score">
                            <div class="team-container">
                                ${homeLogo ? `<img src="${homeLogo}" alt="${homeTeam}" style="height: 40px; width: 40px; object-fit: contain;">` : '<div style="width: 40px;"></div>'}
                                <span class="team-name">${homeTeam}</span>
                            </div>
                            <span class="score">${homeScore} - ${awayScore}</span>
                            <div class="team-container">
                                ${awayLogo ? `<img src="${awayLogo}" alt="${awayTeam}" style="height: 40px; width: 40px; object-fit: contain;">` : '<div style="width: 40px;"></div>'}
                                <span class="team-name">${awayTeam}</span>
                            </div>
                        </div>
                        <div class="match-status">${status}</div>
                    </div>
                `;
            }).join('');
        } else {
            document.getElementById('fixtures-list').innerHTML = '<p>No fixtures available</p>';
        }
        
        // Tournament info
        if (details.tournament) {
            const tournamentInfo = document.getElementById('team-info');
            if (tournamentInfo) {
                tournamentInfo.innerHTML = `${details.tournament.name || 'Tournament'} ${details.tournament.season || ''}`;
            }
        }
        
    } catch (error) {
        console.error('Error loading team data:', error);
        document.getElementById('team-detail').innerHTML = '<p class="loading">Error loading team data. Please try again.</p>';
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
    window.location.href = `/match.html?id=${id}`;
}

function goToPlayer(id) {
    window.location.href = `/player.html?id=${id}`;
}

function goToTeam(id) {
    window.location.href = `/team.html?id=${id}`;
}

function goToManager(id) {
    window.location.href = `/manager.html?id=${id}`;
}