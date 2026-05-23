// Search page functionality
let searchResultsCache = [];

document.addEventListener('DOMContentLoaded', async () => {
    await checkAuth();
    updateAuthUI();
    
    // Получаем query из URL path /search/query
    const path = window.location.pathname;
    let query = '';
    
    if (path.startsWith('/search/')) {
        query = decodeURIComponent(path.substring(8));
    } else {
        const params = new URLSearchParams(window.location.search);
        query = params.get('q');
    }
    
    if (query) {
        document.getElementById('search-input').value = query;
        await performSearch();
    }

    document.getElementById('login-form').addEventListener('submit', handleLogin);
    document.getElementById('register-form').addEventListener('submit', handleRegister);
});

// Функция для загрузки деталей игроков
async function loadPlayerDetails(athleteId) {
    try {
        const response = await fetch(`/api/player/${athleteId}/details`);
        if (response.ok) {
            return await response.json();
        }
    } catch (error) {
        console.error(`Error loading player ${athleteId} details:`, error);
    }
    return null;
}

// Функция для загрузки деталей команд
async function loadTeamDetails(teamId) {
    try {
        const response = await fetch(`/api/team/${teamId}/details`);
        if (response.ok) {
            return await response.json();
        }
    } catch (error) {
        console.error(`Error loading team ${teamId} details:`, error);
    }
    return null;
}

// Функция для загрузки деталей турниров
async function loadTournamentDetails(tournamentId) {
    try {
        const response = await fetch(`/api/tournament/${tournamentId}/details`);
        if (response.ok) {
            return await response.json();
        }
    } catch (error) {
        console.error(`Error loading tournament ${tournamentId} details:`, error);
    }
    return null;
}

// Функция для загрузки деталей менеджеров
async function loadManagerDetails(managerId) {
    try {
        const response = await fetch(`/api/manager/${managerId}/details`);
        if (response.ok) {
            return await response.json();
        }
    } catch (error) {
        console.error(`Error loading manager ${managerId} details:`, error);
    }
    return null;
}

async function performSearch() {
    const query = document.getElementById('search-input').value.trim();
    if (!query) {
        document.getElementById('search-results').innerHTML = '<p class="loading">Enter a search term</p>';
        return;
    }

    // Обновляем URL без перезагрузки
    window.history.pushState({}, '', `/search/${encodeURIComponent(query)}`);

    const resultsDiv = document.getElementById('search-results');
    resultsDiv.innerHTML = '<p class="loading">Searching...</p>';

    try {
        const response = await searchAPI.search(query);
        console.log('Search response:', response);
        
        if (!response || (response.total === 0 && (!response.results || response.results.length === 0))) {
            resultsDiv.innerHTML = '<p class="loading">No results found</p>';
            return;
        }

        // Собираем все результаты в один массив с дополнительными данными
        const allResults = [];
        
        if (response.results && Array.isArray(response.results)) {
            for (const item of response.results) {
                const result = {
                    type: item.type,
                    score: item.score,
                    data: item.data
                };
                
                // Загружаем дополнительные данные в зависимости от типа
                if (item.type === 'player' && item.data.athlete_id) {
                    const details = await loadPlayerDetails(item.data.athlete_id);
                    if (details) {
                        result.details = details;
                    }
                } else if (item.type === 'team' && item.data.team_id) {
                    const details = await loadTeamDetails(item.data.team_id);
                    if (details) {
                        result.details = details;
                    }
                } else if (item.type === 'tournament' && item.data.tournament_id) {
                    const details = await loadTournamentDetails(item.data.tournament_id);
                    if (details) {
                        result.details = details;
                    }
                } else if (item.type === 'manager' && item.data.manager_id) {
                    const details = await loadManagerDetails(item.data.manager_id);
                    if (details) {
                        result.details = details;
                    }
                }
                
                allResults.push(result);
            }
        }
        
        // Сортируем по убыванию score
        allResults.sort((a, b) => b.score - a.score);
        
        // Отрисовываем все результаты в одном списке
        let html = '<div class="search-results-list">';
        
        for (const result of allResults) {
            html += createSearchResultCard(result);
        }
        
        html += '</div>';
        
        if (allResults.length === 0) {
            html = '<p class="loading">No results found</p>';
        }
        
        resultsDiv.innerHTML = html;
        
    } catch (error) {
        console.error('Search error:', error);
        resultsDiv.innerHTML = '<p class="loading">Search failed. Please try again.</p>';
    }
}

function createSearchResultCard(result) {
    const type = result.type;
    const score = result.score ? (result.score * 100).toFixed(1) : null;
    const data = result.data;
    const details = result.details;
    
    // Определяем тип на русском и иконку
    const typeInfo = {
        player: { label: 'Игрок', icon: '👤', color: '#3498db' },
        team: { label: 'Команда', icon: '🏆', color: '#2ecc71' },
        tournament: { label: 'Турнир', icon: '🏆', color: '#e74c3c' },
        manager: { label: 'Тренер', icon: '👨‍✈️', color: '#f39c12' }
    };
    
    const info = typeInfo[type] || { label: type.toUpperCase(), icon: '📌', color: '#95a5a6' };
    
    let content = '';
    
    if (type === 'player') {
        const firstName = details?.first_name || data.first_name || '';
        const lastName = details?.last_name || data.last_name || '';
        const fullName = `${firstName} ${lastName}`.trim() || data.name || 'Unknown';
        const photo = details?.url_photo || '';
        const position = details?.position || data.position || 'N/A';
        const nationName = details?.nation?.name || '';
        const nationFlag = details?.nation?.url_flag || '';
        
        const positionNames = { 'G': 'Вратарь', 'D': 'Защитник', 'M': 'Полузащитник', 'F': 'Нападающий' };
        const positionDisplay = positionNames[position] || position;
        
        content = `
            <div style="display: flex; align-items: center; gap: 15px;">
                ${photo ? `<img src="${photo}" alt="${fullName}" style="width: 60px; height: 60px; border-radius: 50%; object-fit: cover;">` : 
                          `<div style="width: 60px; height: 60px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 2rem;">👤</div>`}
                <div style="flex: 1;">
                    <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
                        <strong style="font-size: 1.1rem;">${fullName}</strong>
                        ${nationFlag ? `<img src="${nationFlag}" alt="${nationName}" style="height: 20px;">` : ''}
                        ${nationName ? `<span style="font-size: 0.8rem; color: #666;">${nationName}</span>` : ''}
                    </div>
                    <div style="font-size: 0.85rem; color: #666;">${positionDisplay}</div>
                    ${score ? `<div style="font-size: 0.75rem; color: #3498db; margin-top: 4px;">Match: ${score}%</div>` : ''}
                </div>
            </div>
        `;
    } 
    else if (type === 'team') {
        const teamName = details?.name || data.name || 'Unknown';
        const logo = details?.url_logo || data.url_logo || '';
        const tournament = details?.tournament?.name || '';
        
        content = `
            <div style="display: flex; align-items: center; gap: 15px;">
                ${logo ? `<img src="${logo}" alt="${teamName}" style="width: 50px; height: 50px; object-fit: contain;">` : 
                         `<div style="width: 50px; height: 50px; background: #ddd; border-radius: 8px; display: flex; align-items: center; justify-content: center; font-size: 1.5rem;">🏆</div>`}
                <div style="flex: 1;">
                    <strong style="font-size: 1.1rem;">${teamName}</strong>
                    ${tournament ? `<div style="font-size: 0.8rem; color: #666;">${tournament}</div>` : ''}
                    ${score ? `<div style="font-size: 0.75rem; color: #3498db; margin-top: 4px;">Match: ${score}%</div>` : ''}
                </div>
            </div>
        `;
    }
    else if (type === 'tournament') {
        const tournamentName = details?.name || data.name || 'Unknown';
        const logo = details?.url_logo || data.url_logo || '';
        const season = details?.season || data.season || '';
        const country = details?.country?.name || '';
        const countryFlag = details?.country?.url_flag || '';
        
        content = `
            <div style="display: flex; align-items: center; gap: 15px;">
                ${logo ? `<img src="${logo}" alt="${tournamentName}" style="width: 50px; height: 50px; object-fit: contain;">` : 
                         `<div style="width: 50px; height: 50px; background: #ddd; border-radius: 8px; display: flex; align-items: center; justify-content: center; font-size: 1.5rem;">🏆</div>`}
                <div style="flex: 1;">
                    <strong style="font-size: 1.1rem;">${tournamentName}</strong>
                    <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
                        ${season ? `<span style="font-size: 0.8rem; color: #666;">${season}</span>` : ''}
                        ${countryFlag ? `<img src="${countryFlag}" alt="${country}" style="height: 16px;">` : ''}
                        ${country ? `<span style="font-size: 0.8rem; color: #666;">${country}</span>` : ''}
                    </div>
                    ${score ? `<div style="font-size: 0.75rem; color: #3498db; margin-top: 4px;">Match: ${score}%</div>` : ''}
                </div>
            </div>
        `;
    }
    else if (type === 'manager') {
        const firstName = details?.first_name || '';
        const lastName = details?.last_name || '';
        const fullName = `${firstName} ${lastName}`.trim() || data.name || 'Unknown';
        const photo = details?.url_photo || '';
        const nationName = details?.nation?.name || '';
        const nationFlag = details?.nation?.url_flag || '';
        
        content = `
            <div style="display: flex; align-items: center; gap: 15px;">
                ${photo ? `<img src="${photo}" alt="${fullName}" style="width: 60px; height: 60px; border-radius: 50%; object-fit: cover;">` : 
                          `<div style="width: 60px; height: 60px; background: #ddd; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 2rem;">👨‍✈️</div>`}
                <div style="flex: 1;">
                    <div style="display: flex; align-items: center; gap: 8px; flex-wrap: wrap;">
                        <strong style="font-size: 1.1rem;">${fullName}</strong>
                        ${nationFlag ? `<img src="${nationFlag}" alt="${nationName}" style="height: 20px;">` : ''}
                        ${nationName ? `<span style="font-size: 0.8rem; color: #666;">${nationName}</span>` : ''}
                    </div>
                    <div style="font-size: 0.85rem; color: #666;">Manager</div>
                    ${score ? `<div style="font-size: 0.75rem; color: #3498db; margin-top: 4px;">Match: ${score}%</div>` : ''}
                </div>
            </div>
        `;
    }
    
    // Генерация ID для ссылки
    let linkId = '';
    let linkUrl = '';
    if (type === 'player') {
        const playerId = details?.athlete_id || data.athlete_id || data.id;
        linkId = playerId;
        linkUrl = `/player?id=${playerId}`;
    } else if (type === 'team') {
        const teamId = details?.team_id || data.team_id || data.id;
        linkId = teamId;
        linkUrl = `/team?id=${teamId}`;
    } else if (type === 'tournament') {
        const tournamentId = details?.tournament_id || data.tournament_id || data.id;
        linkId = tournamentId;
        linkUrl = `/tournament?id=${tournamentId}`;
    } else if (type === 'manager') {
        const managerId = details?.manager_id || data.manager_id || data.id;
        linkId = managerId;
        linkUrl = `/manager?id=${managerId}`;
    }
    
    return `
        <div class="result-card" onclick="goToPage('${linkUrl}')" style="display: flex; justify-content: space-between; align-items: center; padding: 1.25rem;">
            <div style="flex: 1;">
                ${content}
            </div>
            <div style="display: flex; flex-direction: column; align-items: flex-end; gap: 5px;">
                <span class="result-type" style="background: ${info.color};">${info.icon} ${info.label}</span>
                <span style="font-size: 0.7rem; color: #666;">Score: ${(score || 0)}%</span>
            </div>
        </div>
    `;
}

function goToPage(url) {
    window.location.href = url;
}

// Функции для перехода (если их нет в app.js)
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