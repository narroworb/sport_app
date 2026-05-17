// Search page functionality
document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();
    
    const params = new URLSearchParams(window.location.search);
    const query = params.get('q');
    
    if (query) {
        document.getElementById('search-input').value = query;
        await performSearch();
    }

    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
});

async function performSearch() {
    const query = document.getElementById('search-input').value.trim();
    if (!query) {
        document.getElementById('search-results').innerHTML = '<p class="loading">Enter a search term</p>';
        return;
    }

    // Обновляем URL без перезагрузки страницы
    const newUrl = `${window.location.pathname}?q=${encodeURIComponent(query)}`;
    window.history.pushState({}, '', newUrl);

    const resultsDiv = document.getElementById('search-results');
    resultsDiv.innerHTML = '<p class="loading">Searching...</p>';

    try {
        const response = await searchAPI.search(query);
        console.log('Search response:', response); // Для отладки
        
        // Проверяем формат ответа
        if (!response || (response.total === 0 && (!response.results || response.results.length === 0))) {
            resultsDiv.innerHTML = '<p class="loading">No results found</p>';
            return;
        }

        // Если API возвращает объект с полем "results" (как в вашем случае)
        if (response.results && Array.isArray(response.results)) {
            // Группируем результаты по типу
            const grouped = {
                player: [],
                team: [],
                tournament: [],
                manager: []
            };
            
            response.results.forEach(item => {
                if (grouped[item.type]) {
                    grouped[item.type].push(item.data);
                }
            });
            
            let html = '';
            
            if (grouped.player.length > 0) {
                html += '<h3>Players</h3><div class="search-results-grid">';
                html += grouped.player.map(player => createSearchPlayerResult(player)).join('');
                html += '</div>';
            }
            
            if (grouped.team.length > 0) {
                html += '<h3>Teams</h3><div class="search-results-grid">';
                html += grouped.team.map(team => createSearchTeamResult(team)).join('');
                html += '</div>';
            }
            
            if (grouped.tournament.length > 0) {
                html += '<h3>Tournaments</h3><div class="search-results-grid">';
                html += grouped.tournament.map(tournament => createSearchTournamentResult(tournament)).join('');
                html += '</div>';
            }
            
            if (grouped.manager.length > 0) {
                html += '<h3>Managers</h3><div class="search-results-grid">';
                html += grouped.manager.map(manager => createSearchManagerResult(manager)).join('');
                html += '</div>';
            }
            
            if (html === '') {
                html = '<p class="loading">No results found</p>';
            }
            
            resultsDiv.innerHTML = html;
        }
        // Если API возвращает старый формат (с players, teams, etc)
        else if (response.players || response.teams || response.tournaments || response.managers) {
            let html = '';

            if (response.players && response.players.length > 0) {
                html += '<h3>Players</h3><div class="search-results-grid">';
                html += response.players.map(p => createSearchPlayerResult(p)).join('');
                html += '</div>';
            }

            if (response.teams && response.teams.length > 0) {
                html += '<h3>Teams</h3><div class="search-results-grid">';
                html += response.teams.map(t => createSearchTeamResult(t)).join('');
                html += '</div>';
            }

            if (response.tournaments && response.tournaments.length > 0) {
                html += '<h3>Tournaments</h3><div class="search-results-grid">';
                html += response.tournaments.map(t => createSearchTournamentResult(t)).join('');
                html += '</div>';
            }

            if (response.managers && response.managers.length > 0) {
                html += '<h3>Managers</h3><div class="search-results-grid">';
                html += response.managers.map(m => createSearchManagerResult(m)).join('');
                html += '</div>';
            }

            resultsDiv.innerHTML = html || '<p class="loading">No results found</p>';
        }
        else {
            resultsDiv.innerHTML = '<p class="loading">No results found</p>';
        }
    } catch (error) {
        console.error('Search error:', error);
        resultsDiv.innerHTML = '<p class="loading">Search failed. Please try again.</p>';
    }
}

function createSearchPlayerResult(player) {
    const name = player.first_name || player.FirstName || player.name || 'Player';
    const position = player.position || player.Position || 'N/A';
    const id = player.athlete_id || player.AthleteID || player.id || player.player_id;
    
    return `
        <div class="result-card" onclick="goToPlayer(${id})">
            <span class="result-type">PLAYER</span>
            <h3>${name}</h3>
            <p class="result-info">Position: ${position}</p>
            <a href="/player.html?id=${id}">View Profile →</a>
        </div>
    `;
}

function createSearchTeamResult(team) {
    const name = team.name || team.Name || 'Team';
    const id = team.team_id || team.TeamID || team.id;
    const logo = team.url_logo || '';
    
    return `
        <div class="result-card" onclick="goToTeam(${id})">
            <span class="result-type">TEAM</span>
            ${logo ? `<img src="${logo}" alt="${name}" style="height: 40px; margin-bottom: 10px;">` : ''}
            <h3>${name}</h3>
            <p class="result-info">Football Club</p>
            <a href="/team.html?id=${id}">View Team →</a>
        </div>
    `;
}

function createSearchTournamentResult(tournament) {
    const name = tournament.name || tournament.Name || 'Tournament';
    const season = tournament.season || tournament.Season || '';
    const id = tournament.tournament_id || tournament.TournamentID || tournament.id;
    const logo = tournament.url_logo || '';
    
    return `
        <div class="result-card" onclick="goToTournament(${id})">
            <span class="result-type">TOURNAMENT</span>
            ${logo ? `<img src="${logo}" alt="${name}" style="height: 40px; margin-bottom: 10px;">` : ''}
            <h3>${name}</h3>
            <p class="result-info">${season ? 'Season ' + season : 'Championship'}</p>
            <a href="/tournament.html?id=${id}">View Tournament →</a>
        </div>
    `;
}

function createSearchManagerResult(manager) {
    const name = manager.first_name || manager.FirstName || manager.name || 'Manager';
    const id = manager.manager_id || manager.ManagerID || manager.id;
    
    return `
        <div class="result-card" onclick="goToManager(${id})">
            <span class="result-type">MANAGER</span>
            <h3>${name}</h3>
            <p class="result-info">Football Manager</p>
            <a href="/manager.html?id=${id}">View Profile →</a>
        </div>
    `;
}

// Функции для перехода (если их нет в app.js)
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