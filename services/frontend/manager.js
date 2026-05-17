// Manager page functionality
let managerId;
let isFavorite = false;

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();

    const params = new URLSearchParams(window.location.search);
    managerId = params.get('id');

    if (managerId) {
        await loadManagerData();
        if (TokenManager.hasToken()) {
            document.getElementById('fav-btn').style.display = 'inline-block';
        }
    }

    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
});

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

        // Update header with photo and flag
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

        // Stats grid - используем данные из stats.full
        if (stats && stats.full) {
            const statsGrid = document.getElementById('manager-stats-grid');
            const statsData = stats.full;
            
            // Важные статистики для отображения
            const statsToShow = [
                { key: 'total_matches', label: 'Matches', isPercent: false },
                { key: 'win_percentage', label: 'Win %', isPercent: true, multiply: 100 },
                { key: 'avg_points', label: 'Avg Points', isPercent: false },
                { key: 'average_ball_possession', label: 'Possession %', isPercent: false, suffix: '%' },
                { key: 'goals', label: 'Goals', isPercent: false },
                { key: 'goals_per_90', label: 'Goals/90', isPercent: false },
                { key: 'goals_conceded', label: 'Goals Conceded', isPercent: false },
                { key: 'goals_conceded_per_90', label: 'Conceded/90', isPercent: false },
                { key: 'yellow_cards', label: 'Yellow Cards', isPercent: false },
                { key: 'red_cards', label: 'Red Cards', isPercent: false },
                { key: 'average_pass_accuracy', label: 'Pass Accuracy', isPercent: true, multiply: 100 }
            ];
            
            const displayStats = [];
            for (const stat of statsToShow) {
                let value = statsData[stat.key];
                if (value !== undefined && value !== null) {
                    let displayValue = value;
                    if (stat.multiply) {
                        displayValue = value * stat.multiply;
                    }
                    if (typeof displayValue === 'number') {
                        if (stat.isPercent || stat.key.includes('percentage') || stat.key.includes('accuracy')) {
                            displayValue = displayValue.toFixed(1) + '%';
                        } else if (stat.suffix === '%') {
                            displayValue = displayValue.toFixed(1) + '%';
                        } else if (stat.key.includes('_per_90') && displayValue < 10) {
                            displayValue = displayValue.toFixed(2);
                        } else if (Number.isInteger(displayValue)) {
                            displayValue = displayValue;
                        } else {
                            displayValue = displayValue.toFixed(2);
                        }
                    }
                    displayStats.push({ label: stat.label, value: displayValue });
                }
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
        } else {
            document.getElementById('manager-stats-grid').innerHTML = '<p>No stats available</p>';
        }

        // Teams - обрабатываем объект вида {"Burnley":["2023/2024"],"FC Bayern München":["2024/2025","2025/2026"]}
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
                <div class="card team-card" style="cursor:pointer; padding: 1rem; margin-bottom: 1rem;" onclick="searchTeam('${encodeURIComponent(team.name)}')">
                    <h3>${team.name}</h3>
                    <p>${team.season}</p>
                </div>
            `).join('');
            
            if (teamsArray.length === 0) {
                list.innerHTML = '<p>No team history available</p>';
            }
        } else {
            document.getElementById('teams-list').innerHTML = '<p>No team history available</p>';
        }

        // Fixtures (матчи менеджера)
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
                        <div class="match-status" style="margin-top: 10px; padding-top: 10px; border-top: 1px solid #eee; text-align: center;">
                            ${status}
                        </div>
                    </div>
                `;
            }).join('');
            
            if (fixtures.length === 0) {
                list.innerHTML = '<p>No fixtures available</p>';
            }
        } else {
            document.getElementById('fixtures-list').innerHTML = '<p>No fixtures available</p>';
        }
        
        // Проверяем избранное
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

function formatLabel(key) {
    return key.replace(/_/g, ' ')
              .split(' ')
              .map(w => w.charAt(0).toUpperCase() + w.slice(1))
              .join(' ');
}

function goToMatch(id) {
    if (id && id > 0) {
        window.location.href = `/match.html?id=${id}`;
    }
}

function searchTeam(teamName) {
    if (teamName) {
        window.location.href = `/search.html?q=${encodeURIComponent(teamName)}`;
    }
}