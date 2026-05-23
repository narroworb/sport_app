// API Configuration
// Использование локальных эндпоинтов через прокси Nginx для разработки
// В продакшене обновите эти URL на ваш домен
const API_BASE = '/api';
const AUTH_BASE = window.location.origin;

// Управление токенами
class TokenManager {
    static setToken(token) {
        try {
            localStorage.setItem('token', token);
            console.log('Токен сохранен в localStorage');
        } catch (e) {
            console.error('Не удалось сохранить токен:', e);
        }
    }

    static getToken() {
        try {
            const token = localStorage.getItem('token');
            console.log('Токен получен из localStorage, существует:', !!token);
            return token;
        } catch (e) {
            console.error('Не удалось получить токен:', e);
            return null;
        }
    }

    static removeToken() {
        try {
            localStorage.removeItem('token');
            console.log('Токен удален из localStorage');
        } catch (e) {
            console.error('Не удалось удалить токен:', e);
        }
    }

    static hasToken() {
        const token = this.getToken();
        return !!token && token.length > 0;
    }
}

// API Helper
async function apiCall(endpoint, options = {}) {
    const url = endpoint.startsWith('http') ? endpoint : `${API_BASE}${endpoint}`;
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers,
    };

    const token = TokenManager.getToken();
    if (token) {
        // Для API запросов используем Bearer
        headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(url, {
        ...options,
        headers,
    });

    if (response.status === 401) {
        TokenManager.removeToken();
        if (window.location.pathname !== '/') {
            window.location.href = '/';
        }
        throw new Error('Неавторизовано');
    }

    if (!response.ok) {
        const error = await response.text();
        throw new Error(error || `Ошибка API: ${response.status}`);
    }

    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
        return await response.json();
    }
    return response;
}

// API аутентификации
const authAPI = {
    register: async (username, password) => {
        const response = await fetch(`${AUTH_BASE}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
        });
        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || 'Ошибка регистрации');
        }
        return await response.json();
    },

    login: async (username, password) => {
        const response = await fetch(`${AUTH_BASE}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
        });
        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || 'Ошибка входа');
        }
        return await response.json();
    },

    getMe: async () => {
        const token = TokenManager.getToken();
        if (!token) return null;
        
        // Для /me - без Bearer!
        const response = await fetch(`${AUTH_BASE}/me`, {
            method: 'GET',
            headers: { 
                'Authorization': token  // Только токен, без Bearer
            },
        });
        if (!response.ok) return null;
        return await response.json();
    }
};

// API поиска
const searchAPI = {
    search: async (query) => {
        return await apiCall(`/search?q=${encodeURIComponent(query)}`);
    }
};

// API игроков
const playerAPI = {
    getDetails: async (id) => {
        return await apiCall(`/player/${id}/details`);
    },

    getStats: async (id) => {
        return await apiCall(`/player/${id}/stats`);
    },

    getFixtures: async (id, limit = 10, offset = 0) => {
        return await apiCall(`/player/${id}/fixtures?limit=${limit}&offset=${offset}`);
    },

    getTeams: async (id) => {
        return await apiCall(`/player/${id}/teams`);
    },

    getFavorites: async () => {
        return await apiCall(`/player/favorite`);
    },

    addFavorite: async (id) => {
        return await apiCall(`/player/${id}/favorite`, { method: 'POST' });
    },

    removeFavorite: async (id) => {
        return await apiCall(`/player/${id}/favorite`, { method: 'DELETE' });
    }
};

// API команд
const teamAPI = {
    getDetails: async (id) => {
        return await apiCall(`/team/${id}/details`);
    },

    getStats: async (id) => {
        return await apiCall(`/team/${id}/stats`);
    },

    getNextGame: async (id) => {
        return await apiCall(`/team/${id}/next_game`);
    },

    getStandings: async (id) => {
        return await apiCall(`/team/${id}/standings`);
    },

    getPlayers: async (id) => {
        return await apiCall(`/team/${id}/players`);
    },

    getFixtures: async (id, limit = 10, offset = 0) => {
        return await apiCall(`/team/${id}/fixtures?limit=${limit}&offset=${offset}`);
    },

    getManager: async (id) => {
        return await apiCall(`/team/${id}/manager`);
    },

    getPlayersStats: async (id) => {
        return await apiCall(`/team/${id}/players_stats`);
    },

    getFavorites: async () => {
        return await apiCall(`/team/favorite`);
    },

    addFavorite: async (id) => {
        return await apiCall(`/team/${id}/favorite`, { method: 'POST' });
    },

    removeFavorite: async (id) => {
        return await apiCall(`/team/${id}/favorite`, { method: 'DELETE' });
    }
};

// API турниров
const tournamentAPI = {
    getDetails: async (id) => {
        return await apiCall(`/tournament/${id}/details`);
    },

    getTable: async (id) => {
        return await apiCall(`/tournament/${id}/table`);
    },

    getTeamsStats: async (id) => {
        return await apiCall(`/tournament/${id}/stats/teams`);
    },

    getPlayersStats: async (id) => {
        return await apiCall(`/tournament/${id}/stats/players`);
    },

    getTableGraph: async (id) => {
        return await apiCall(`/tournament/${id}/table/graph`);
    },

    getFixtures: async (id) => {
        return await apiCall(`/tournament/${id}/fixtures`);
    },

    getFavorites: async () => {
        return await apiCall(`/tournament/favorite`);
    },

    addFavorite: async (id) => {
        return await apiCall(`/tournament/${id}/favorite`, { method: 'POST' });
    },

    removeFavorite: async (id) => {
        return await apiCall(`/tournament/${id}/favorite`, { method: 'DELETE' });
    }
};

// API тренеров
const managerAPI = {
    getDetails: async (id) => {
        return await apiCall(`/manager/${id}/details`);
    },

    getStats: async (id) => {
        return await apiCall(`/manager/${id}/stats`);
    },

    getTeams: async (id) => {
        return await apiCall(`/manager/${id}/teams`);
    },

    getFixtures: async (id) => {
        return await apiCall(`/manager/${id}/fixtures`);
    },

    getFavorites: async () => {
        return await apiCall(`/manager/favorite`);
    },

    addFavorite: async (id) => {
        return await apiCall(`/manager/${id}/favorite`, { method: 'POST' });
    },

    removeFavorite: async (id) => {
        return await apiCall(`/manager/${id}/favorite`, { method: 'DELETE' });
    }
};

// API матчей
const fixtureAPI = {
    getDetails: async (id) => {
        return await apiCall(`/fixture/${id}/details`);
    },

    getPlayersStats: async (id) => {
        return await apiCall(`/fixture/${id}/stats/players`);
    },

    getGoaliesStats: async (id) => {
        return await apiCall(`/fixture/${id}/stats/goalies`);
    },

    getTeamsStats: async (id) => {
        return await apiCall(`/fixture/${id}/stats/teams`);
    },

    getByDate: async (date) => {
        return await apiCall(`/fixture?date=${date}`);
    }
};

// API аналитики
const analyticsAPI = {
    getStats: async (endpoint) => {
        return await apiCall(`/analytics${endpoint}`);
    }
};