// Global State
let currentUser = null;

// Initialize app
document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    await loadRecentMatches();
    await loadTournaments();
    
    // Event listeners
    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
});

// Authentication
async function checkAuth() {
    if (!TokenManager.hasToken()) {
        return;
    }

    try {
        currentUser = await authAPI.getMe();
        if (currentUser) {
            updateAuthUI();
            if (document.getElementById('favorites-section')) {
                document.getElementById('favorites-section').style.display = 'block';
            }
        }
    } catch (error) {
        console.error('Auth check failed:', error);
        TokenManager.removeToken();
    }
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
        
        // Пытаемся распарсить JSON токен
        let token = null;
        try {
            const data = JSON.parse(responseText);
            token = data.token;
        } catch (e) {
            // Если ответ не JSON, возможно токен пришел как plain text
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
        
        // Получаем текст ответа (не JSON)
        const responseText = await response.text();
        console.log('Registration response:', response.status, responseText);
        
        if (!response.ok) {
            // Обрабатываем разные статусы
            if (response.status === 409) {
                errorDiv.textContent = `Username "${username}" already exists. Please choose another username or login.`;
            } else if (response.status === 400) {
                errorDiv.textContent = responseText || 'Invalid username or password format.';
            } else {
                errorDiv.textContent = responseText || `Registration failed: ${response.status}`;
            }
            return;
        }
        
        // Регистрация успешна
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

// Content Loading
async function loadRecentMatches() {
    const container = document.getElementById('recent-matches');
    if (!container) return;

    try {
        const today = new Date().toISOString().split('T')[0];
        const matches = await fixtureAPI.getByDate(today);

        if (Array.isArray(matches) && matches.length > 0) {
            container.innerHTML = matches.slice(0, 6).map(match => createMatchCard(match)).join('');
        } else {
            container.innerHTML = '<p class="loading">No matches today</p>';
        }
    } catch (error) {
        console.error('Error loading matches:', error);
        const message = error.message || '';
        if (message.includes('404') || message.toLowerCase().includes('fixtures not found')) {
            container.innerHTML = '<p class="loading">No matches today</p>';
        } else {
            container.innerHTML = '<p class="loading">Failed to load matches</p>';
        }
    }
}

async function loadTournaments() {
    const container = document.getElementById('tournaments');
    if (!container) return;

    try {
        // In a real app, you'd have an endpoint to list tournaments
        // For now, we'll show a message
        container.innerHTML = `
            <div class="card tournament-card">
                <h3>Browse Tournaments</h3>
                <p>Use the search feature to find tournaments and explore their standings and statistics.</p>
                <a href="/search.html" class="btn-primary" style="text-decoration: none; text-align: center;">Go to Search</a>
            </div>
        `;
    } catch (error) {
        console.error('Error loading tournaments:', error);
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
    const dateStr = date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
    
    // Правильные пути для вашего API
    const homeTeam = match.home_team?.name || 'Team A';
    const awayTeam = match.away_team?.name || 'Team B';
    const homeScore = match.home_team_score ?? 0;
    const awayScore = match.away_team_score ?? 0;
    const status = match.status || 'Not started';
    const matchId = match.match_id ?? 0;

    return `
        <div class="card match-card" onclick="goToMatch(${matchId})">
            <div class="match-header">${dateStr}</div>
            <div class="match-score">
                <span class="team-name">${homeTeam}</span>
                <span class="score">${homeScore} - ${awayScore}</span>
                <span class="team-name">${awayTeam}</span>
            </div>
            <div class="match-status">${status}</div>
            ${match.tournament?.name ? `<div class="tournament-name">${match.tournament.name}</div>` : ''}
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

// Navigation
function performSearch() {
    const query = document.getElementById('main-search').value;
    if (query.trim()) {
        window.location.href = `/search.html?q=${encodeURIComponent(query)}`;
    }
}

function goToPlayer(id) {
    window.location.href = `/player.html?id=${id}`;
}

function goToTeam(id) {
    window.location.href = `/team.html?id=${id}`;
}

function goToTournament(id) {
    window.location.href = `/tournament.html?id=${id}`;
}

function goToManager(id) {
    window.location.href = `/manager.html?id=${id}`;
}

function goToMatch(id) {
    window.location.href = `/match.html?id=${id}`;
}
