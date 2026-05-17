// Match page functionality
let matchId;

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();

    const params = new URLSearchParams(window.location.search);
    matchId = params.get('id');

    if (matchId) {
        await loadMatchData();
    }

    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
});

// Функция для определения статуса матча
function getMatchStatusInfo(details) {
    const matchDate = details.date ? new Date(details.date) : new Date();
    const now = new Date();
    const status = details.status || 'Scheduled';
    
    // Сыгранный матч
    if (status === 'Ended') {
        return { isFinished: true, isLive: false, isPending: false, message: null };
    }
    
    // Матч в процессе (если есть такая информация)
    if (status === 'Live' || status === 'In Progress' || status === 'First Half' || status === 'Second Half') {
        return { isFinished: false, isLive: true, isPending: false, message: '⚽ Match is currently in progress! Live updates coming soon.' };
    }
    
    // Матч еще не начался
    if (matchDate > now) {
        const timeUntil = Math.ceil((matchDate - now) / (1000 * 60));
        return { 
            isFinished: false, 
            isLive: false, 
            isPending: false, 
            message: `⏰ Match starts in approximately ${timeUntil} minutes` 
        };
    }
    
    // Матч уже должен был закончиться, но статус не обновлен
    if (matchDate < now && status !== 'Ended') {
        return { 
            isFinished: false, 
            isLive: false, 
            isPending: true, 
            message: '⏳ The match has ended. Statistics will be updated shortly. Please check back in a few minutes.' 
        };
    }
    
    return { isFinished: false, isLive: false, isPending: false, message: null };
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
            document.getElementById('match-detail').innerHTML = '<p class="loading">Match not found</p>';
            return;
        }

        // Получаем информацию о статусе матча
        const matchStatusInfo = getMatchStatusInfo(details);
        const isFinished = matchStatusInfo.isFinished;
        const isLive = matchStatusInfo.isLive;
        const isPending = matchStatusInfo.isPending;
        
        console.log('Match status info:', matchStatusInfo);

        // Загружаем детали команд по их ID
        const homeTeamId = details.home_team?.team_id || details.home_team_id;
        const awayTeamId = details.away_team?.team_id || details.away_team_id;
        
        let homeTeamDetails = null;
        let awayTeamDetails = null;
        
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
        
        // Загружаем аналитику только для предстоящих матчей (не начавшихся)
        let matchProbabilities = null;
        let tournamentStandings = null;
        let tournamentId = null;
        
        if (!isFinished && !isLive && !isPending && details.tournament?.tournament_id) {
            tournamentId = details.tournament.tournament_id;
            const season = details.tournament?.season || '2024/2025';
            
            try {
                // Загружаем вероятности исхода матча
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
                // Загружаем таблицу турнира
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
        
        // Уникальные ID игроков
        const playerIds = [...new Set(allPlayers.map(p => p.athlete?.athlete_id).filter(id => id && id > 0))];
        console.log('Player IDs to fetch:', playerIds);
        
        // Загружаем детали игроков (ограничим первыми 20 для производительности)
        const playerDetailsMap = new Map();
        for (const playerId of playerIds.slice(0, 20)) {
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

        // Обновляем заголовок с фото
        document.getElementById('home-team').innerHTML = `
            ${homeTeamLogo ? `<img src="${homeTeamLogo}" alt="${homeTeamName}" style="height: 40px; vertical-align: middle; margin-right: 10px;">` : ''}
            ${homeTeamName}
        `;
        document.getElementById('home-team').style.cursor = 'pointer';
        document.getElementById('home-team').dataset.id = homeTeamId;
        document.getElementById('home-team').onclick = () => goToTeam(homeTeamId);
        
        document.getElementById('away-team').innerHTML = `
            ${awayTeamLogo ? `<img src="${awayTeamLogo}" alt="${awayTeamName}" style="height: 40px; vertical-align: middle; margin-right: 10px;">` : ''}
            ${awayTeamName}
        `;
        document.getElementById('away-team').style.cursor = 'pointer';
        document.getElementById('away-team').dataset.id = awayTeamId;
        document.getElementById('away-team').onclick = () => goToTeam(awayTeamId);

        document.getElementById('match-score').textContent = `${homeScore} : ${awayScore}`;
        
        const matchDate = details.date ? new Date(details.date) : new Date();
        const dateStr = matchDate.toLocaleDateString() + ' ' + matchDate.toLocaleTimeString([], {hour:'2-digit', minute:'2-digit'});
        document.getElementById('match-date').textContent = dateStr;
        
        // Обновляем статус с дополнительной информацией
        let statusText = details.status || 'Scheduled';
        if (isLive) {
            statusText = '🔴 LIVE - Match in progress';
        } else if (isPending) {
            statusText = '⏳ Match ended - Awaiting statistics update';
        } else if (!isFinished && matchDate > new Date()) {
            statusText = '📅 Upcoming match';
        }
        document.getElementById('match-status').textContent = statusText;

        document.getElementById('home-team-name-header').innerHTML = homeTeamLogo ? `<img src="${homeTeamLogo}" alt="${homeTeamName}" style="height: 30px; vertical-align: middle; margin-right: 10px;">${homeTeamName}` : homeTeamName;
        document.getElementById('away-team-name-header').innerHTML = awayTeamLogo ? `<img src="${awayTeamLogo}" alt="${awayTeamName}" style="height: 30px; vertical-align: middle; margin-right: 10px;">${awayTeamName}` : awayTeamName;

        // Добавляем информационное сообщение о статусе в overview
        const overviewDiv = document.getElementById('match-overview');
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
        
        // Добавляем информационное сообщение в начало overview
        if (infoMessageHtml) {
            overviewDiv.insertAdjacentHTML('afterbegin', infoMessageHtml);
        }

        // Если матч предстоящий - показываем аналитику
        if (!isFinished && !isLive && !isPending && matchProbabilities) {
            // Преобразуем вероятности в проценты и округляем
            const homeProb = (matchProbabilities.p_home_win || 0) * 100;
            const drawProb = (matchProbabilities.p_draw || 0) * 100;
            const awayProb = (matchProbabilities.p_away_win || 0) * 100;
            
            // Создаем HTML для вероятностей
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
            
            // Добавляем вероятности в overview
            overviewDiv.insertAdjacentHTML('beforeend', probabilitiesHtml);
        }
        
        // Показываем таблицу турнира для предстоящих матчей
        if (!isFinished && !isLive && !isPending && tournamentStandings && tournamentStandings.length > 0) {
            const standingsHtml = `
                <div class="standings-preview" style="margin-top: 2rem;">
                    <h3>🏆 Tournament Standings</h3>
                    <table class="table">
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

        // Team Stats - показываем только для сыгранных матчей или live
        if (teamStats && (isFinished || isLive)) {
            const homePossession = teamStats.ball_possession_home_team ?? 0;
            const awayPossession = teamStats.ball_possession_away_team ?? 0;
            const homeShots = teamStats.total_shots_home_team ?? 0;
            const awayShots = teamStats.total_shots_away_team ?? 0;

            document.getElementById('home-shots').textContent = homeShots;
            document.getElementById('away-shots').textContent = awayShots;
            document.getElementById('home-possession').textContent = homePossession + '%';
            document.getElementById('away-possession').textContent = awayPossession + '%';

            const tbody = document.getElementById('team-stats-tbody');
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
        } else if (isPending) {
            // Для матчей, ожидающих обновления статистики
            const teamStatsDiv = document.getElementById('match-team-stats');
            const tbody = document.getElementById('team-stats-tbody');
            if (tbody) {
                tbody.innerHTML = '<tr><td colspan="3" style="text-align: center;">⏳ Match statistics are being processed and will appear shortly...</td></tr>';
            }
            
            document.getElementById('home-shots').textContent = '-';
            document.getElementById('away-shots').textContent = '-';
            document.getElementById('home-possession').textContent = '-';
            document.getElementById('away-possession').textContent = '-';
        } else if (!isFinished && !isLive) {
            // Для предстоящих матчей
            const teamStatsDiv = document.getElementById('match-team-stats');
            const tbody = document.getElementById('team-stats-tbody');
            if (tbody) {
                tbody.innerHTML = '<tr><td colspan="3" style="text-align: center;">📊 Statistics will be available after the match</td></tr>';
            }
            
            document.getElementById('home-shots').textContent = '-';
            document.getElementById('away-shots').textContent = '-';
            document.getElementById('home-possession').textContent = '-';
            document.getElementById('away-possession').textContent = '-';
        }

        // Players Stats - показываем только для сыгранных матчей или live
        if (playersStats && (isFinished || isLive)) {
            // Home players
            if (playersStats.home_team && Array.isArray(playersStats.home_team)) {
                const homePlayers = playersStats.home_team.filter(p => p.start_player === true || p.minutes_played > 0);
                document.getElementById('home-players-tbody').innerHTML = homePlayers.map(p => {
                    const playerId = p.athlete?.athlete_id;
                    const playerDetails = playerDetailsMap.get(playerId);
                    const name = playerDetails?.first_name && playerDetails?.last_name 
                        ? `${playerDetails.first_name} ${playerDetails.last_name}`
                        : `${p.athlete?.first_name || ''} ${p.athlete?.last_name || ''}`.trim() || 'Unknown';
                    const photo = playerDetails?.url_photo || '';
                    const rating = p.rating ? p.rating.toFixed(1) : '-';
                    
                    return `
                        <tr onclick="goToPlayer(${playerId})" style="cursor:pointer;">
                            <td style="display: flex; align-items: center; gap: 10px;">
                                ${photo ? `<img src="${photo}" alt="${name}" style="width: 35px; height: 35px; border-radius: 50%; object-fit: cover;">` : 
                                          `<div style="width: 35px; height: 35px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center;">👤</div>`}
                                <strong>${name}</strong>
                              </td>
                            <td>${p.goals || 0}</td>
                            <td>${rating}</td>
                        </tr>
                    `;
                }).join('');
            }

            // Away players
            if (playersStats.away_team && Array.isArray(playersStats.away_team)) {
                const awayPlayers = playersStats.away_team.filter(p => p.start_player === true || p.minutes_played > 0);
                document.getElementById('away-players-tbody').innerHTML = awayPlayers.map(p => {
                    const playerId = p.athlete?.athlete_id;
                    const playerDetails = playerDetailsMap.get(playerId);
                    const name = playerDetails?.first_name && playerDetails?.last_name 
                        ? `${playerDetails.first_name} ${playerDetails.last_name}`
                        : `${p.athlete?.first_name || ''} ${p.athlete?.last_name || ''}`.trim() || 'Unknown';
                    const photo = playerDetails?.url_photo || '';
                    const rating = p.rating ? p.rating.toFixed(1) : '-';
                    
                    return `
                        <tr onclick="goToPlayer(${playerId})" style="cursor:pointer;">
                            <td style="display: flex; align-items: center; gap: 10px;">
                                ${photo ? `<img src="${photo}" alt="${name}" style="width: 35px; height: 35px; border-radius: 50%; object-fit: cover;">` : 
                                          `<div style="width: 35px; height: 35px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center;">👤</div>`}
                                <strong>${name}</strong>
                              </td>
                            <td>${p.goals || 0}</td>
                            <td>${rating}</td>
                        </tr>
                    `;
                }).join('');
            }
        } else if (isPending) {
            document.getElementById('home-players-tbody').innerHTML = '<tr><td colspan="3" style="text-align: center;">⏳ Player statistics are being processed...</td></tr>';
            document.getElementById('away-players-tbody').innerHTML = '<tr><td colspan="3" style="text-align: center;">⏳ Player statistics are being processed...</td></tr>';
        } else if (!isFinished && !isLive) {
            document.getElementById('home-players-tbody').innerHTML = '<tr><td colspan="3" style="text-align: center;">📊 Player stats will be available after the match</td></tr>';
            document.getElementById('away-players-tbody').innerHTML = '<tr><td colspan="3" style="text-align: center;">📊 Player stats will be available after the match</td></tr>';
        }

        // Goalies Stats - показываем только для сыгранных матчей или live
        if (goaliesStats && (isFinished || isLive)) {
            // Home goalies
            if (goaliesStats.home_team && Array.isArray(goaliesStats.home_team)) {
                const homeGoalies = goaliesStats.home_team.filter(g => g.start_player === true || g.minutes_played > 0);
                document.getElementById('home-goalies-tbody').innerHTML = homeGoalies.map(g => {
                    const playerId = g.athlete?.athlete_id;
                    const playerDetails = playerDetailsMap.get(playerId);
                    const name = playerDetails?.first_name && playerDetails?.last_name 
                        ? `${playerDetails.first_name} ${playerDetails.last_name}`
                        : `${g.athlete?.first_name || ''} ${g.athlete?.last_name || ''}`.trim() || 'Unknown';
                    const photo = playerDetails?.url_photo || '';
                    const rating = g.rating ? g.rating.toFixed(1) : '-';
                    
                    return `
                        <tr onclick="goToPlayer(${playerId})" style="cursor:pointer;">
                            <td style="display: flex; align-items: center; gap: 10px;">
                                ${photo ? `<img src="${photo}" alt="${name}" style="width: 35px; height: 35px; border-radius: 50%; object-fit: cover;">` : 
                                          `<div style="width: 35px; height: 35px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center;">🧤</div>`}
                                <strong>${name}</strong>
                              </td>
                            <td>${g.saves || 0}</td>
                            <td>${rating}</td>
                         </tr>
                    `;
                }).join('');
            }

            // Away goalies
            if (goaliesStats.away_team && Array.isArray(goaliesStats.away_team)) {
                const awayGoalies = goaliesStats.away_team.filter(g => g.start_player === true || g.minutes_played > 0);
                document.getElementById('away-goalies-tbody').innerHTML = awayGoalies.map(g => {
                    const playerId = g.athlete?.athlete_id;
                    const playerDetails = playerDetailsMap.get(playerId);
                    const name = playerDetails?.first_name && playerDetails?.last_name 
                        ? `${playerDetails.first_name} ${playerDetails.last_name}`
                        : `${g.athlete?.first_name || ''} ${g.athlete?.last_name || ''}`.trim() || 'Unknown';
                    const photo = playerDetails?.url_photo || '';
                    const rating = g.rating ? g.rating.toFixed(1) : '-';
                    
                    return `
                        <tr onclick="goToPlayer(${playerId})" style="cursor:pointer;">
                            <td style="display: flex; align-items: center; gap: 10px;">
                                ${photo ? `<img src="${photo}" alt="${name}" style="width: 35px; height: 35px; border-radius: 50%; object-fit: cover;">` : 
                                          `<div style="width: 35px; height: 35px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center;">🧤</div>`}
                                <strong>${name}</strong>
                              </td>
                            <td>${g.saves || 0}</td>
                            <td>${rating}</td>
                         </tr>
                    `;
                }).join('');
            }
        } else if (isPending) {
            document.getElementById('home-goalies-tbody').innerHTML = '<tr><td colspan="3" style="text-align: center;">⏳ Goalkeeper statistics are being processed...</td></tr>';
            document.getElementById('away-goalies-tbody').innerHTML = '<tr><td colspan="3" style="text-align: center;">⏳ Goalkeeper statistics are being processed...</td></tr>';
        } else if (!isFinished && !isLive) {
            document.getElementById('home-goalies-tbody').innerHTML = '<tr><td colspan="3" style="text-align: center;">📊 Goalkeeper stats will be available after the match</td></tr>';
            document.getElementById('away-goalies-tbody').innerHTML = '<tr><td colspan="3" style="text-align: center;">📊 Goalkeeper stats will be available after the match</td></tr>';
        }
        
    } catch (error) {
        console.error('Error loading match data:', error);
        document.getElementById('match-detail').innerHTML = '<p class="loading">Error loading match data. Please try again.</p>';
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
        window.location.href = `/player.html?id=${id}`;
    }
}

function goToTeam(id) {
    if (id && id > 0) {
        window.location.href = `/team.html?id=${id}`;
    }
}

function goToMatch(id) {
    if (id && id > 0) {
        window.location.href = `/match.html?id=${id}`;
    }
}