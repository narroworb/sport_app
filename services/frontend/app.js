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
console.log('Страница загружена, проверяем сохраненный токен:', TokenManager.getToken());

async function checkAuth() {
    const token = TokenManager.getToken();
    console.log('=== НАЧАЛО ПРОВЕРКИ АУТЕНТИФИКАЦИИ ===');
    
    if (!token) {
        console.log('Токен не найден');
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
        
        console.log('Статус ответа /me:', response.status);
        
        if (response.ok) {
            currentUser = await response.json();
            console.log('Пользователь аутентифицирован:', currentUser);
            updateAuthUI();
            if (document.getElementById('favorites-section')) {
                document.getElementById('favorites-section').style.display = 'block';
                if (typeof loadFavorites === 'function') {
                    loadFavorites();
                }
            }
        } else {
            console.log('Ошибка аутентификации');
            currentUser = null;
            updateAuthUI();
        }
    } catch (error) {
        console.error('Ошибка проверки аутентификации:', error);
        currentUser = null;
        updateAuthUI();
    }
    console.log('=== КОНЕЦ ПРОВЕРКИ АУТЕНТИФИКАЦИИ ===');
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
        console.log('Ответ входа:', response.status, responseText);
        
        if (!response.ok) {
            errorDiv.textContent = responseText || 'Ошибка входа. Проверьте данные и попробуйте снова.';
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
            errorDiv.textContent = 'Ошибка входа: токен не получен';
        }
    } catch (error) {
        console.error('Ошибка входа:', error);
        errorDiv.textContent = 'Ошибка входа. Попробуйте позже.';
    }
}

async function handleRegister(e) {
    e.preventDefault();
    const username = document.getElementById('register-username').value;
    const password = document.getElementById('register-password').value;
    const passwordConfirm = document.getElementById('register-password-confirm').value;
    const errorDiv = document.getElementById('register-error');
    
    if (password !== passwordConfirm) {
        errorDiv.textContent = 'Пароли не совпадают';
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
        console.log('Ответ регистрации:', response.status, responseText);
        
        if (!response.ok) {
            if (response.status === 409) {
                errorDiv.textContent = `Пользователь "${username}" уже существует. Выберите другой логин или войдите.`;
            } else if (response.status === 400) {
                errorDiv.textContent = responseText || 'Неверный формат логина или пароля.';
            } else {
                errorDiv.textContent = responseText || `Ошибка регистрации: ${response.status}`;
            }
            return;
        }
        
        errorDiv.style.color = '#2ecc71';
        errorDiv.textContent = responseText || 'Регистрация прошла успешно! Войдите в систему.';
        
        setTimeout(() => {
            switchTab('login');
            document.getElementById('register-form').reset();
            errorDiv.textContent = '';
            errorDiv.style.color = '#e74c3c';
        }, 2000);
    } catch (error) {
        console.error('Ошибка регистрации:', error);
        errorDiv.textContent = 'Ошибка регистрации. Попробуйте позже.';
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

function switchTab(tabName, event) {
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));
    
    const targetButton = event?.target || document.querySelector(`.tab-btn[onclick*="${tabName}"]`);
    targetButton?.classList.add('active');
    document.getElementById(`${tabName}-form`).classList.add('active');
}

async function loadRecentMatches() {
    const container = document.getElementById('recent-matches');
    if (!container) return;

    try {
        container.innerHTML = '<p class="loading">Загрузка матчей...</p>';
        
        const year = currentMatchDate.getFullYear();
        const month = String(currentMatchDate.getMonth() + 1).padStart(2, '0');
        const day = String(currentMatchDate.getDate()).padStart(2, '0');
        const dateStr = `${year}-${month}-${day}`;
        
        const matches = await fixtureAPI.getByDate(dateStr);

        if (Array.isArray(matches) && matches.length > 0) {
            container.innerHTML = matches.slice(0, 12).map(match => createMatchCard(match)).join('');
        } else {
            container.innerHTML = '<p class="loading">Матчи на эту дату не найдены</p>';
        }
    } catch (error) {
        console.error('Ошибка загрузки матчей:', error);
        container.innerHTML = '<p class="loading">Не удалось загрузить матчи</p>';
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
            html += '<h3>⭐ Избранные игроки</h3><div class="players-grid">';
            html += favPlayers.map(p => createPlayerCard(p)).join('');
            html += '</div>';
        }
        
        if (favTeams?.length > 0) {
            html += '<h3>🏆 Избранные команды</h3><div class="teams-grid">';
            html += favTeams.map(t => createTeamCard(t)).join('');
            html += '</div>';
        }
        
        if (favTournaments?.length > 0) {
            html += '<h3>🏅 Избранные турниры</h3><div class="tournaments-grid">';
            html += favTournaments.map(t => createTournamentCard(t)).join('');
            html += '</div>';
        }

        if (!html) {
            html = '<p class="empty-favorites">✨ Пока нет избранных. Найдите и добавьте игроков, команды и турниры.</p>';
        }

        container.innerHTML = html;
    } catch (error) {
        console.error('Ошибка загрузки избранного:', error);
        container.innerHTML = '<p class="error">Не удалось загрузить избранное</p>';
    }
}

function createMatchCard(match) {
    const date = new Date(match.date);
    const timeStr = date.toLocaleTimeString('ru-RU', {hour: '2-digit', minute:'2-digit'});
    const dateStr = date.toLocaleDateString('ru-RU');
    
    const homeTeam = match.home_team?.name || 'Команда А';
    const awayTeam = match.away_team?.name || 'Команда Б';
    const homeTeamLogo = match.home_team?.url_logo || '';
    const awayTeamLogo = match.away_team?.url_logo || '';
    const homeScore = match.home_team_score ?? 0;
    const awayScore = match.away_team_score ?? 0;
    const status = match.status || 'Not started';
    const matchId = match.match_id ?? 0;
    const tournamentName = match.tournament?.name || '';
    const tournamentLogo = match.tournament?.url_logo || '';
    
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
        <div class="card match-card" onclick="goToMatch(${matchId})">
            <div class="match-header" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.75rem;">
                <div style="display: flex; align-items: center; gap: 0.5rem;">
                    ${tournamentLogo ? `<img src="${tournamentLogo}" alt="${tournamentName}" style="height: 20px; width: 20px; object-fit: contain;">` : '🏆'}
                    <span style="font-size: 0.8rem; color: #666;">${tournamentName || 'Турнир'}</span>
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
    const name = player.first_name || player.FirstName || player.name || 'Игрок';
    const position = player.position || player.Position || 'Н/Д';
    const id = player.athlete_id || player.AthleteID || player.id;
    
    const positionMap = {
        'Goalkeeper': 'Вратарь',
        'Defender': 'Защитник',
        'Midfielder': 'Полузащитник',
        'Forward': 'Нападающий'
    };
    
    const russianPosition = positionMap[position] || position;
    
    return `
        <div class="card player-card" onclick="goToPlayer(${id})">
            <h3>${name}</h3>
            <p>${russianPosition}</p>
            <div class="card-meta">
                <span class="badge">${russianPosition}</span>
            </div>
        </div>
    `;
}

function createTeamCard(team) {
    const name = team.name || team.Name || 'Команда';
    const id = team.team_id || team.TeamID || team.id;
    
    return `
        <div class="card team-card" onclick="goToTeam(${id})">
            <h3>${name}</h3>
            <p>⚽ Футбольный клуб</p>
        </div>
    `;
}

function createTournamentCard(tournament) {
    const name = tournament.name || tournament.Name || 'Турнир';
    const season = tournament.season || tournament.Season || '';
    const id = tournament.tournament_id || tournament.TournamentID || tournament.id;
    
    return `
        <div class="card tournament-card" onclick="goToTournament(${id})">
            <h3>${name}</h3>
            ${season ? `<p>📅 Сезон ${season}</p>` : ''}
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