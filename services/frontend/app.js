// Global State
let currentUser = null;
let currentMatchDate = new Date();

// Initialize app
document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    
    // Инициализируем календарь
    initDatePicker();
    
    await loadRecentMatches();
    await loadFavorites();
    
    // Event listeners
    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
    
    // Date navigation events
    const prevBtn = document.getElementById('prev-day');
    const nextBtn = document.getElementById('next-day');
    const todayBtn = document.getElementById('today-date-filter');
    
    if (prevBtn) {
        prevBtn.addEventListener('click', () => changeDate(-1));
    }
    if (nextBtn) {
        nextBtn.addEventListener('click', () => changeDate(1));
    }
    if (todayBtn) {
        todayBtn.addEventListener('click', () => goToToday());
    }
});

function initDatePicker() {
    const dateInput = document.getElementById('match-date-calendar');
    if (dateInput) {
        flatpickr(dateInput, {
            dateFormat: "Y-m-d",
            locale: "ru",
            onChange: function(selectedDates, dateStr) {
                if (dateStr) {
                    currentMatchDate = new Date(dateStr);
                    updateDateDisplay();
                    loadRecentMatches();
                }
            }
        });
    }
    updateDateDisplay();
}

function updateDateDisplay() {
    const dateInput = document.getElementById('match-date-calendar');
    if (dateInput) {
        const year = currentMatchDate.getFullYear();
        const month = String(currentMatchDate.getMonth() + 1).padStart(2, '0');
        const day = String(currentMatchDate.getDate()).padStart(2, '0');
        dateInput.value = `${year}-${month}-${day}`;
        
        // Обновляем flatpickr если он уже инициализирован
        if (dateInput._flatpickr) {
            dateInput._flatpickr.setDate(currentMatchDate);
        }
    }
}

function changeDate(delta) {
    const newDate = new Date(currentMatchDate);
    newDate.setDate(newDate.getDate() + delta);
    currentMatchDate = newDate;
    updateDateDisplay();
    loadRecentMatches();
}

function goToToday() {
    currentMatchDate = new Date();
    updateDateDisplay();
    loadRecentMatches();
}

// Проверка при загрузке страницы
console.log('Page loaded, checking stored token:', TokenManager.getToken());

async function checkAuth() {
    const token = TokenManager.getToken();
    console.log('=== CHECK AUTH START ===');
    
    if (!token) {
        console.log('No token found');
        currentUser = null;
        updateAuthUI();
        return;
    }

    try {
        const response = await fetch('/me', {
            method: 'GET',
            headers: {
                'Authorization': token
            }
        });
        
        console.log('/me response status:', response.status);
        
        if (response.ok) {
            currentUser = await response.json();
            console.log('User authenticated:', currentUser);
            updateAuthUI();
            if (document.getElementById('favorites-section')) {
                document.getElementById('favorites-section').style.display = 'block';
                if (typeof loadFavorites === 'function') {
                    loadFavorites();
                }
            }
        } else {
            console.log('Auth failed');
            currentUser = null;
            updateAuthUI();
        }
    } catch (error) {
        console.error('Auth check failed:', error);
        currentUser = null;
        updateAuthUI();
    }
    console.log('=== CHECK AUTH END ===');
}

function updateAuthUI() {
    const authBtn = document.getElementById('auth-btn');
    const logoutBtn = document.getElementById('logout-btn');
    const userName = document.getElementById('user-name');
    
    if (currentUser) {
        authBtn.style.display = 'none';
        logoutBtn.style.display = 'inline-block';
        userName.textContent = currentUser.username || currentUser.user_id;
        userName.style.display = 'inline-block';
    } else {
        authBtn.style.display = 'inline-block';
        logoutBtn.style.display = 'none';
        userName.style.display = 'none';
    }
}

async function handleLogin(e) {
    e.preventDefault();
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;
    const errorDiv = document.getElementById('login-error');
    
    try {
        errorDiv.textContent = '';
        
        const response = await fetch('/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
        });
        
        const responseText = await response.text();
        console.log('Login response:', response.status, responseText);
        
        if (!response.ok) {
            errorDiv.textContent = responseText || 'Login failed. Check your credentials.';
            return;
        }
        
        let token = null;
        try {
            const data = JSON.parse(responseText);
            token = data.token;
        } catch (e) {
            token = responseText;
        }
        
        if (token) {
            TokenManager.setToken(token);
            currentUser = { username };
            updateAuthUI();
            toggleAuthPanel();
            document.getElementById('login-form').reset();
            
            if (document.getElementById('favorites-section')) {
                document.getElementById('favorites-section').style.display = 'block';
                loadFavorites();
            }
        } else {
            errorDiv.textContent = 'Login failed: No token received';
        }
    } catch (error) {
        console.error('Login error:', error);
        errorDiv.textContent = 'Login failed. Please try again later.';
    }
}

async function handleRegister(e) {
    e.preventDefault();
    const username = document.getElementById('register-username').value;
    const password = document.getElementById('register-password').value;
    const passwordConfirm = document.getElementById('register-password-confirm').value;
    const errorDiv = document.getElementById('register-error');
    
    if (password !== passwordConfirm) {
        errorDiv.textContent = 'Passwords do not match';
        return;
    }

    try {
        errorDiv.textContent = '';
        errorDiv.style.color = '#e74c3c';
        
        const response = await fetch('/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
        });
        
        const responseText = await response.text();
        console.log('Registration response:', response.status, responseText);
        
        if (!response.ok) {
            if (response.status === 409) {
                errorDiv.textContent = `Username "${username}" already exists. Please choose another username or login.`;
            } else if (response.status === 400) {
                errorDiv.textContent = responseText || 'Invalid username or password format.';
            } else {
                errorDiv.textContent = responseText || `Registration failed: ${response.status}`;
            }
            return;
        }
        
        errorDiv.style.color = '#2ecc71';
        errorDiv.textContent = responseText || 'Registration successful! Please login.';
        
        setTimeout(() => {
            switchTab('login');
            document.getElementById('register-form').reset();
            errorDiv.textContent = '';
            errorDiv.style.color = '#e74c3c';
        }, 2000);
    } catch (error) {
        console.error('Registration error:', error);
        errorDiv.textContent = 'Registration failed. Please try again later.';
    }
}

function logout() {
    TokenManager.removeToken();
    currentUser = null;
    updateAuthUI();
    if (document.getElementById('favorites-section')) {
        document.getElementById('favorites-section').style.display = 'none';
    }
    window.location.href = '/';
}

function toggleAuthPanel() {
    const panel = document.getElementById('auth-panel');
    panel.style.display = panel.style.display === 'none' ? 'flex' : 'none';
}

function switchTab(tabName) {
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));
    
    event.target.classList.add('active');
    document.getElementById(`${tabName}-form`).classList.add('active');
}

async function loadRecentMatches() {
    const container = document.getElementById('recent-matches');
    if (!container) return;

    try {
        container.innerHTML = '<p class="loading">Loading matches...</p>';
        
        const year = currentMatchDate.getFullYear();
        const month = String(currentMatchDate.getMonth() + 1).padStart(2, '0');
        const day = String(currentMatchDate.getDate()).padStart(2, '0');
        const dateStr = `${year}-${month}-${day}`;
        
        const matches = await fixtureAPI.getByDate(dateStr);

        if (Array.isArray(matches) && matches.length > 0) {
            container.innerHTML = matches.slice(0, 12).map(match => createMatchCard(match)).join('');
        } else {
            container.innerHTML = '<p class="loading">No matches found for this date</p>';
        }
    } catch (error) {
        console.error('Error loading matches:', error);
        container.innerHTML = '<p class="loading">Failed to load matches</p>';
    }
}

async function loadFavorites() {
    const container = document.getElementById('favorites-content');
    if (!container) return;

    try {
        const favPlayers = TokenManager.hasToken() ? await playerAPI.getFavorites() : [];
        const favTeams = TokenManager.hasToken() ? await teamAPI.getFavorites() : [];
        const favTournaments = TokenManager.hasToken() ? await tournamentAPI.getFavorites() : [];
        
        let html = '';
        
        if (favPlayers?.length > 0) {
            html += '<h3>Favorite Players</h3><div class="players-grid">';
            html += favPlayers.map(p => createPlayerCard(p)).join('');
            html += '</div>';
        }
        
        if (favTeams?.length > 0) {
            html += '<h3>Favorite Teams</h3><div class="teams-grid">';
            html += favTeams.map(t => createTeamCard(t)).join('');
            html += '</div>';
        }
        
        if (favTournaments?.length > 0) {
            html += '<h3>Favorite Tournaments</h3><div class="tournaments-grid">';
            html += favTournaments.map(t => createTournamentCard(t)).join('');
            html += '</div>';
        }

        if (!html) {
            html = '<p>No favorites yet. Search and add your favorite players, teams, and tournaments!</p>';
        }

        container.innerHTML = html;
    } catch (error) {
        console.error('Error loading favorites:', error);
    }
}

function createMatchCard(match) {
    const date = new Date(match.date);
    const timeStr = date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
    const dateStr = date.toLocaleDateString();
    
    const homeTeam = match.home_team?.name || 'Team A';
    const awayTeam = match.away_team?.name || 'Team B';
    const homeTeamLogo = match.home_team?.url_logo || '';
    const awayTeamLogo = match.away_team?.url_logo || '';
    const homeScore = match.home_team_score ?? 0;
    const awayScore = match.away_team_score ?? 0;
    const status = match.status || 'Not started';
    const matchId = match.match_id ?? 0;
    const tournamentName = match.tournament?.name || '';
    const tournamentLogo = match.tournament?.url_logo || '';
    
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
            <div class="match-header" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem;">
                <div style="display: flex; align-items: center; gap: 0.5rem;">
                    ${tournamentLogo ? `<img src="${tournamentLogo}" alt="${tournamentName}" style="height: 20px; width: 20px; object-fit: contain;">` : '🏆'}
                    <span style="font-size: 0.8rem; color: #666;">${tournamentName}</span>
                </div>
                <div style="font-size: 0.8rem; color: #666;">${dateStr} • ${timeStr}</div>
            </div>
            <div class="match-score" style="display: flex; justify-content: space-between; align-items: center; gap: 0.5rem; margin-bottom: 0.5rem;">
                <div style="display: flex; align-items: center; gap: 0.5rem; flex: 1; justify-content: flex-end;">
                    ${homeTeamLogo ? `<img src="${homeTeamLogo}" alt="${homeTeam}" style="height: 35px; width: 35px; object-fit: contain;">` : ''}
                    <span class="team-name" style="font-weight: 600; font-size: 1rem;">${homeTeam}</span>
                </div>
                <span class="score" style="font-size: 1.3rem; font-weight: bold; color: #2c3e50; min-width: 60px; text-align: center;">${homeScore} - ${awayScore}</span>
                <div style="display: flex; align-items: center; gap: 0.5rem; flex: 1;">
                    <span class="team-name" style="font-weight: 600; font-size: 1rem;">${awayTeam}</span>
                    ${awayTeamLogo ? `<img src="${awayTeamLogo}" alt="${awayTeam}" style="height: 35px; width: 35px; object-fit: contain;">` : ''}
                </div>
            </div>
            <div class="match-status ${statusClass}" style="margin-top: 0.5rem; padding-top: 0.5rem; border-top: 1px solid #eee; text-align: center; font-size: 0.8rem;">
                ${statusText}
            </div>
        </div>
    `;
}

function createPlayerCard(player) {
    const name = player.first_name || player.FirstName || player.name || 'Player';
    const position = player.position || player.Position || 'N/A';
    const id = player.athlete_id || player.AthleteID || player.id;
    
    return `
        <div class="card player-card" onclick="goToPlayer(${id})">
            <h3>${name}</h3>
            <p>${position}</p>
            <div class="card-meta">
                <span class="badge">${position}</span>
            </div>
        </div>
    `;
}

function createTeamCard(team) {
    const name = team.name || team.Name || 'Team';
    const id = team.team_id || team.TeamID || team.id;
    
    return `
        <div class="card team-card" onclick="goToTeam(${id})">
            <h3>${name}</h3>
            <p>Football Club</p>
        </div>
    `;
}

function createTournamentCard(tournament) {
    const name = tournament.name || tournament.Name || 'Tournament';
    const season = tournament.season || tournament.Season || '';
    const id = tournament.tournament_id || tournament.TournamentID || tournament.id;
    
    return `
        <div class="card tournament-card" onclick="goToTournament(${id})">
            <h3>${name}</h3>
            ${season ? `<p>Season ${season}</p>` : ''}
        </div>
    `;
}

function performSearch() {
    const query = document.getElementById('main-search').value;
    if (query.trim()) {
        window.location.href = `/search?q=${encodeURIComponent(query)}`;
    }
}

function goToPlayer(id) {
    window.location.href = `/player?id=${id}`;
}

function goToTeam(id) {
    window.location.href = `/team?id=${id}`;
}

function goToTournament(id) {
    window.location.href = `/tournament?id=${id}`;
}

function goToManager(id) {
    window.location.href = `/manager?id=${id}`;
}

function goToMatch(id) {
    window.location.href = `/match?id=${id}`;
}