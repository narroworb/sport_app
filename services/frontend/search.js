// Функционал страницы поиска

document.addEventListener('DOMContentLoaded', async () => {
    try {
        await checkAuth();
        updateAuthUI();
        
        // Получаем query из URL
        const path = window.location.pathname;
        let query = '';
        
        if (path.startsWith('/search/')) {
            query = decodeURIComponent(path.substring(8));
        } else {
            const params = new URLSearchParams(window.location.search);
            query = params.get('q');
        }
        
        if (query) {
            const searchInput = document.getElementById('search-input');
            if (searchInput) {
                searchInput.value = query;
            }
            await performSearch();
        }

        const loginForm = document.getElementById('login-form');
        const registerForm = document.getElementById('register-form');
        if (loginForm) loginForm.addEventListener('submit', handleLogin);
        if (registerForm) registerForm.addEventListener('submit', handleRegister);
        
        const searchInput = document.getElementById('search-input');
        if (searchInput) {
            searchInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    performSearch();
                }
            });
        }
    } catch (error) {
        console.error('Ошибка инициализации:', error);
    }
});

async function performSearch() {
    // Проверяем существование элементов с защитой
    const searchInput = document.getElementById('search-input');
    if (!searchInput) {
        console.error('Элемент search-input не найден');
        return;
    }
    
    const query = searchInput.value.trim();
    
    const resultsWrapper = document.getElementById('search-results-wrapper');
    if (!resultsWrapper) {
        console.error('Элемент search-results-wrapper не найден');
        // Пытаемся найти альтернативный контейнер или создаём его
        const container = document.querySelector('.container .section');
        if (container) {
            const newWrapper = document.createElement('div');
            newWrapper.id = 'search-results-wrapper';
            container.appendChild(newWrapper);
            resultsWrapper = newWrapper;
        } else {
            return;
        }
    }
    
    if (!query) {
        resultsWrapper.innerHTML = '<p class="loading">Введите поисковый запрос</p>';
        return;
    }

    // Обновляем URL без перезагрузки
    window.history.pushState({}, '', `/search/${encodeURIComponent(query)}`);

    resultsWrapper.innerHTML = '<div class="search-loading">🔍 Поиск...</div>';

    try {
        if (typeof searchAPI === 'undefined' || !searchAPI.search) {
            throw new Error('API поиска не загружен');
        }
        
        const response = await searchAPI.search(query);
        console.log('Ответ поиска:', response);
        
        if (!response || !response.results || response.results.length === 0) {
            resultsWrapper.innerHTML = `
                <div class="search-no-results">
                    <div class="no-results-icon">😔</div>
                    <p>Ничего не найдено по запросу "${escapeHtml(query)}"</p>
                    <p style="font-size: 0.9rem; color: #666;">Попробуйте изменить поисковый запрос</p>
                </div>
            `;
            return;
        }

        // Сортируем результаты по убыванию score
        const sortedResults = [...response.results].sort((a, b) => b.score - a.score);
        
        let html = `
            <div class="search-results-container">
                <div class="search-stats">
                    <span class="stats-icon">🔍</span>
                    Найдено: ${response.total || sortedResults.length} результатов по запросу "${escapeHtml(query)}"
                </div>
                <div class="search-results-list">
        `;
        
        for (const result of sortedResults) {
            html += createSearchResultCard(result);
        }
        
        html += `
                </div>
            </div>
        `;
        
        resultsWrapper.innerHTML = html;
        
    } catch (error) {
        console.error('Ошибка поиска:', error);
        resultsWrapper.innerHTML = `
            <div class="search-error">
                <div class="error-icon">❌</div>
                <p>Ошибка поиска</p>
                <p style="font-size: 0.9rem;">Пожалуйста, попробуйте снова</p>
                <button onclick="performSearch()" class="btn-primary" style="margin-top: 1rem;">Повторить</button>
            </div>
        `;
    }
}

function createSearchResultCard(result) {
    const type = result.type;
    const data = result.data;
    
    // Типы сущностей
    const typeInfo = {
        player: { label: 'Игрок', icon: '⚽', gradient: 'linear-gradient(135deg, #3498db, #2980b9)' },
        team: { label: 'Команда', icon: '🏆', gradient: 'linear-gradient(135deg, #2ecc71, #27ae60)' },
        tournament: { label: 'Турнир', icon: '🏆', gradient: 'linear-gradient(135deg, #e74c3c, #c0392b)' },
        manager: { label: 'Тренер', icon: '👨‍✈️', gradient: 'linear-gradient(135deg, #f39c12, #e67e22)' }
    };
    
    const info = typeInfo[type] || { label: type.toUpperCase(), icon: '📌', gradient: 'linear-gradient(135deg, #95a5a6, #7f8c8d)' };
    
    let content = '';
    let linkUrl = '';
    
    // Получаем ID в зависимости от типа
    let entityId = null;
    if (type === 'player') {
        entityId = data.id || data.athlete_id;
        linkUrl = `/player?id=${entityId}`;
    } else if (type === 'team') {
        entityId = data.team_id || data.id;
        linkUrl = `/team?id=${entityId}`;
    } else if (type === 'tournament') {
        entityId = data.tournament_id || data.id;
        linkUrl = `/tournament?id=${entityId}`;
    } else if (type === 'manager') {
        entityId = data.manager_id || data.id;
        linkUrl = `/manager?id=${entityId}`;
    }
    
    if (type === 'player') {
        const firstName = data.first_name || '';
        const lastName = data.last_name || '';
        const fullName = `${firstName} ${lastName}`.trim() || 'Неизвестно';
        const photo = data.url_photo || '';
        const position = data.position || 'N/A';
        const nation = data.nation || '';
        const status = data.current_status || '';
        const height = data.height || '';
        
        const positionNames = { 'G': 'Вратарь', 'D': 'Защитник', 'M': 'Полузащитник', 'F': 'Нападающий' };
        const positionDisplay = positionNames[position] || position;
        const positionIcons = { 'G': '🧤', 'D': '🛡️', 'M': '⚡', 'F': '🎯' };
        const positionIcon = positionIcons[position] || '⚽';
        
        const statusClass = status === 'Active' ? 'status-active' : 'status-retired';
        const statusText = status === 'Active' ? 'Активен' : 'Завершил карьеру';
        
        content = `
            <div class="result-card-content">
                <div class="result-avatar">
                    ${photo ? 
                        `<img src="${photo}" alt="${escapeHtml(fullName)}" onerror="this.style.display='none'; this.parentElement.innerHTML='<div class=\\'avatar-placeholder\\'>${info.icon}</div>'">` : 
                        `<div class="avatar-placeholder">${info.icon}</div>`
                    }
                </div>
                <div class="result-info">
                    <div class="result-name">
                        <strong>${escapeHtml(fullName)}</strong>
                        ${nation ? `<span class="nation-badge">🌍 ${escapeHtml(nation)}</span>` : ''}
                        <span class="status-badge ${statusClass}">${statusText}</span>
                    </div>
                    <div class="result-details">
                        <span class="detail-badge">${positionIcon} ${escapeHtml(positionDisplay)}</span>
                        ${height ? `<span class="detail-badge">📏 ${height} см</span>` : ''}
                    </div>
                </div>
            </div>
        `;
    }
    else if (type === 'manager') {
        const firstName = data.first_name || '';
        const lastName = data.last_name || '';
        const fullName = `${firstName} ${lastName}`.trim() || 'Неизвестно';
        const photo = data.url_photo || '';
        const nation = data.nation || '';
        
        content = `
            <div class="result-card-content">
                <div class="result-avatar">
                    ${photo ? 
                        `<img src="${photo}" alt="${escapeHtml(fullName)}" onerror="this.style.display='none'; this.parentElement.innerHTML='<div class=\\'avatar-placeholder\\'>${info.icon}</div>'">` : 
                        `<div class="avatar-placeholder">${info.icon}</div>`
                    }
                </div>
                <div class="result-info">
                    <div class="result-name">
                        <strong>${escapeHtml(fullName)}</strong>
                        ${nation ? `<span class="nation-badge">🌍 ${escapeHtml(nation)}</span>` : ''}
                    </div>
                    <div class="result-details">
                        <span class="detail-badge">👨‍✈️ Главный тренер</span>
                    </div>
                </div>
            </div>
        `;
    }
    else if (type === 'team') {
        const teamName = data.name || 'Неизвестно';
        const logo = data.url_logo || '';
        
        content = `
            <div class="result-card-content">
                <div class="result-avatar team-avatar">
                    ${logo ? 
                        `<img src="${logo}" alt="${escapeHtml(teamName)}" onerror="this.style.display='none'; this.parentElement.innerHTML='<div class=\\'avatar-placeholder\\'>${info.icon}</div>'">` : 
                        `<div class="avatar-placeholder">${info.icon}</div>`
                    }
                </div>
                <div class="result-info">
                    <div class="result-name">
                        <strong>${escapeHtml(teamName)}</strong>
                    </div>
                </div>
            </div>
        `;
    }
    else if (type === 'tournament') {
        const tournamentName = data.name || 'Неизвестно';
        const logo = data.url_logo || '';
        const season = data.season || '';
        
        content = `
            <div class="result-card-content">
                <div class="result-avatar">
                    ${logo ? 
                        `<img src="${logo}" alt="${escapeHtml(tournamentName)}" onerror="this.style.display='none'; this.parentElement.innerHTML='<div class=\\'avatar-placeholder\\'>${info.icon}</div>'">` : 
                        `<div class="avatar-placeholder">${info.icon}</div>`
                    }
                </div>
                <div class="result-info">
                    <div class="result-name">
                        <strong>${escapeHtml(tournamentName)}</strong>
                    </div>
                    ${season ? `<div class="result-details"><span class="detail-badge">📅 ${escapeHtml(season)}</span></div>` : ''}
                </div>
            </div>
        `;
    }
    
    return `
        <div class="result-card" onclick="goToPage('${linkUrl}')">
            ${content}
            <div class="result-type-badge" style="background: ${info.gradient}">
                ${info.icon} ${info.label}
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

function goToPage(url) {
    window.location.href = url;
}