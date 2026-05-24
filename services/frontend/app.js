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
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    if (loginForm) loginForm.addEventListener('submit', handleLogin);
    if (registerForm) registerForm.addEventListener('submit', handleRegister);
    
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

async function checkAuth() {
    const token = TokenManager.getToken();
    
    if (!token) {
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
        
        if (response.ok) {
            currentUser = await response.json();
            updateAuthUI();
            if (document.getElementById('favorites-section')) {
                document.getElementById('favorites-section').style.display = 'block';
                loadFavorites();
            }
        } else {
            currentUser = null;
            updateAuthUI();
        }
    } catch (error) {
        console.error('Ошибка проверки аутентификации:', error);
        currentUser = null;
        updateAuthUI();
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
    if (targetButton) targetButton.classList.add('active');
    const form = document.getElementById(`${tabName}-form`);
    if (form) form.classList.add('active');
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
        const [favPlayers, favTeams, favTournaments, favManagers] = await Promise.all([
            TokenManager.hasToken() ? playerAPI.getFavorites() : [],
            TokenManager.hasToken() ? teamAPI.getFavorites() : [],
            TokenManager.hasToken() ? tournamentAPI.getFavorites() : [],
            TokenManager.hasToken() ? managerAPI.getFavorites() : []
        ]);
        
        let html = '';
        
        if (favPlayers?.length > 0) {
            html += '<h3>⭐ Избранные игроки</h3><div class="favorites-grid">';
            for (const player of favPlayers) {
                const details = await loadPlayerDetails(player.athlete_id || player.id);
                html += createFavoritePlayerCard(player, details);
            }
            html += '</div>';
        }
        
        if (favTeams?.length > 0) {
            html += '<h3>🏆 Избранные команды</h3><div class="favorites-grid">';
            for (const team of favTeams) {
                const details = await loadTeamDetails(team.team_id || team.id);
                html += createFavoriteTeamCard(team, details);
            }
            html += '</div>';
        }
        
        if (favTournaments?.length > 0) {
            html += '<h3>🏅 Избранные турниры</h3><div class="favorites-grid">';
            for (const tournament of favTournaments) {
                const details = await loadTournamentDetails(tournament.tournament_id || tournament.id);
                html += createFavoriteTournamentCard(tournament, details);
            }
            html += '</div>';
        }
        
        if (favManagers?.length > 0) {
            html += '<h3>👨‍✈️ Избранные тренеры</h3><div class="favorites-grid">';
            for (const manager of favManagers) {
                const details = await loadManagerDetails(manager.manager_id || manager.id);
                html += createFavoriteManagerCard(manager, details);
            }
            html += '</div>';
        }

        if (!html) {
            html = '<p class="empty-favorites">✨ Пока нет избранных. Найдите и добавьте игроков, команды, турниры и тренеров.</p>';
        }

        container.innerHTML = html;
    } catch (error) {
        console.error('Ошибка загрузки избранного:', error);
        container.innerHTML = '<p class="error">Не удалось загрузить избранное</p>';
    }
}

// Функции для загрузки деталей
async function loadPlayerDetails(athleteId) {
    try {
        const response = await fetch(`/api/player/${athleteId}/details`);
        if (response.ok) return await response.json();
    } catch (error) {
        console.error(`Ошибка загрузки деталей игрока ${athleteId}:`, error);
    }
    return null;
}

async function loadTeamDetails(teamId) {
    try {
        const response = await fetch(`/api/team/${teamId}/details`);
        if (response.ok) return await response.json();
    } catch (error) {
        console.error(`Ошибка загрузки деталей команды ${teamId}:`, error);
    }
    return null;
}

async function loadTournamentDetails(tournamentId) {
    try {
        const response = await fetch(`/api/tournament/${tournamentId}/details`);
        if (response.ok) return await response.json();
    } catch (error) {
        console.error(`Ошибка загрузки деталей турнира ${tournamentId}:`, error);
    }
    return null;
}

async function loadManagerDetails(managerId) {
    try {
        const response = await fetch(`/api/manager/${managerId}/details`);
        if (response.ok) return await response.json();
    } catch (error) {
        console.error(`Ошибка загрузки деталей тренера ${managerId}:`, error);
    }
    return null;
}

function createFavoritePlayerCard(player, details) {
    const firstName = details?.first_name || player.first_name || '';
    const lastName = details?.last_name || player.last_name || '';
    const fullName = `${firstName} ${lastName}`.trim() || player.name || 'Неизвестно';
    const photo = details?.url_photo || player.url_photo || '';
    const position = details?.position || player.position || 'N/A';
    const nation = details?.nation?.name || player.nation || '';
    const playerId = player.athlete_id || player.id;
    
    const positionNames = { 'G': 'Вратарь', 'D': 'Защитник', 'M': 'Полузащитник', 'F': 'Нападающий' };
    const positionDisplay = positionNames[position] || position;
    const positionIcons = { 'G': '🧤', 'D': '🛡️', 'M': '⚡', 'F': '🎯' };
    const positionIcon = positionIcons[position] || '⚽';
    
    return `
        <div class="favorite-card" onclick="goToPlayer(${playerId})">
            <div class="favorite-card-content">
                <div class="favorite-avatar">
                    ${photo ? 
                        `<img src="${photo}" alt="${escapeHtml(fullName)}" onerror="this.style.display='none'; this.parentElement.innerHTML='<div class=\\'avatar-placeholder\\'>⚽</div>'">` : 
                        `<div class="avatar-placeholder">⚽</div>`
                    }
                </div>
                <div class="favorite-info">
                    <div class="favorite-name">
                        <strong>${escapeHtml(fullName)}</strong>
                        ${nation ? `<span class="nation-badge">🌍 ${escapeHtml(nation)}</span>` : ''}
                    </div>
                    <div class="favorite-details">
                        <span class="detail-badge">${positionIcon} ${escapeHtml(positionDisplay)}</span>
                    </div>
                </div>
            </div>
            <div class="favorite-type-badge player-badge">⚽ Игрок</div>
        </div>
    `;
}

function createFavoriteTeamCard(team, details) {
    const teamName = details?.name || team.name || 'Неизвестно';
    const logo = details?.url_logo || team.url_logo || '';
    const tournament = details?.tournament?.name || '';
    const teamId = team.team_id || team.id;
    
    return `
        <div class="favorite-card" onclick="goToTeam(${teamId})">
            <div class="favorite-card-content">
                <div class="favorite-avatar team-avatar">
                    ${logo ? 
                        `<img src="${logo}" alt="${escapeHtml(teamName)}" onerror="this.style.display='none'; this.parentElement.innerHTML='<div class=\\'avatar-placeholder\\'>🏆</div>'">` : 
                        `<div class="avatar-placeholder">🏆</div>`
                    }
                </div>
                <div class="favorite-info">
                    <div class="favorite-name">
                        <strong>${escapeHtml(teamName)}</strong>
                    </div>
                    ${tournament ? `<div class="favorite-details"><span class="detail-badge">🏆 ${escapeHtml(tournament)}</span></div>` : ''}
                </div>
            </div>
            <div class="favorite-type-badge team-badge">🏆 Команда</div>
        </div>
    `;
}

function createFavoriteTournamentCard(tournament, details) {
    const tournamentName = details?.name || tournament.name || 'Неизвестно';
    const logo = details?.url_logo || tournament.url_logo || '';
    const season = details?.season || tournament.season || '';
    const tournamentId = tournament.tournament_id || tournament.id;
    
    return `
        <div class="favorite-card" onclick="goToTournament(${tournamentId})">
            <div class="favorite-card-content">
                <div class="favorite-avatar">
                    ${logo ? 
                        `<img src="${logo}" alt="${escapeHtml(tournamentName)}" onerror="this.style.display='none'; this.parentElement.innerHTML='<div class=\\'avatar-placeholder\\'>🏆</div>'">` : 
                        `<div class="avatar-placeholder">🏆</div>`
                    }
                </div>
                <div class="favorite-info">
                    <div class="favorite-name">
                        <strong>${escapeHtml(tournamentName)}</strong>
                    </div>
                    ${season ? `<div class="favorite-details"><span class="detail-badge">📅 ${escapeHtml(season)}</span></div>` : ''}
                </div>
            </div>
            <div class="favorite-type-badge tournament-badge">🏅 Турнир</div>
        </div>
    `;
}

function createFavoriteManagerCard(manager, details) {
    const firstName = details?.first_name || manager.first_name || '';
    const lastName = details?.last_name || manager.last_name || '';
    const fullName = `${firstName} ${lastName}`.trim() || manager.name || 'Неизвестно';
    const photo = details?.url_photo || manager.url_photo || '';
    const nation = details?.nation?.name || manager.nation || '';
    const managerId = manager.manager_id || manager.id;
    
    return `
        <div class="favorite-card" onclick="goToManager(${managerId})">
            <div class="favorite-card-content">
                <div class="favorite-avatar">
                    ${photo ? 
                        `<img src="${photo}" alt="${escapeHtml(fullName)}" onerror="this.style.display='none'; this.parentElement.innerHTML='<div class=\\'avatar-placeholder\\'>👨‍✈️</div>'">` : 
                        `<div class="avatar-placeholder">👨‍✈️</div>`
                    }
                </div>
                <div class="favorite-info">
                    <div class="favorite-name">
                        <strong>${escapeHtml(fullName)}</strong>
                        ${nation ? `<span class="nation-badge">🌍 ${escapeHtml(nation)}</span>` : ''}
                    </div>
                    <div class="favorite-details">
                        <span class="detail-badge">👨‍✈️ Главный тренер</span>
                    </div>
                </div>
            </div>
            <div class="favorite-type-badge manager-badge">👨‍✈️ Тренер</div>
        </div>
    `;
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
            <div class="match-header">
                <div class="tournament-info">
                    ${tournamentLogo ? `<img src="${tournamentLogo}" alt="${tournamentName}" style="height: 18px; width: 18px; object-fit: contain;">` : '🏆'}
                    <span>${escapeHtml(tournamentName) || 'Турнир'}</span>
                </div>
                <div class="match-time">${dateStr} • ${timeStr}</div>
            </div>
            <div class="match-score">
                <div class="team-container home">
                    ${homeTeamLogo ? `<img src="${homeTeamLogo}" alt="${escapeHtml(homeTeam)}" class="team-logo">` : ''}
                    <span class="team-name">${escapeHtml(homeTeam)}</span>
                </div>
                <span class="score">${homeScore} - ${awayScore}</span>
                <div class="team-container away">
                    <span class="team-name">${escapeHtml(awayTeam)}</span>
                    ${awayTeamLogo ? `<img src="${awayTeamLogo}" alt="${escapeHtml(awayTeam)}" class="team-logo">` : ''}
                </div>
            </div>
            <div class="match-status ${statusClass}">
                ${statusText}
            </div>
        </div>
    `;
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
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