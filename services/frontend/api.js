// API Configuration
// Using local endpoints through Nginx proxy for development
// In production, update these URLs to your domain
const API_BASE = '/api';
const AUTH_BASE = window.location.origin;

// Token Management
class TokenManager {
    static setToken(token) {
        try {
            localStorage.setItem('token', token);
            console.log('Token saved to localStorage');
        } catch (e) {
            console.error('Failed to save token:', e);
        }
    }

    static getToken() {
        try {
            const token = localStorage.getItem('token');
            console.log('Token retrieved from localStorage, exists:', !!token);
            return token;
        } catch (e) {
            console.error('Failed to get token:', e);
            return null;
        }
    }

    static removeToken() {
        try {
            localStorage.removeItem('token');
            console.log('Token removed from localStorage');
        } catch (e) {
            console.error('Failed to remove token:', e);
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
        throw new Error('Unauthorized');
    }

    if (!response.ok) {
        const error = await response.text();
        throw new Error(error || `API Error: ${response.status}`);
    }

    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
        return await response.json();
    }
    return response;
}

// Auth API
const authAPI = {
    register: async (username, password) => {
        const response = await fetch(`${AUTH_BASE}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password }),
        });
        if (!response.ok) {
            const error = await response.text();
            throw new Error(error || 'Registration failed');
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
            throw new Error(error || 'Login failed');
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

// Search API
const searchAPI = {
    search: async (query) => {
        return await apiCall(`/search?q=${encodeURIComponent(query)}`);
    }
};

// Player API
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

// Team API
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

// Tournament API
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

// Manager API
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

// Fixture API
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

// Analytics API
const analyticsAPI = {
    getStats: async (endpoint) => {
        return await apiCall(`/analytics${endpoint}`);
    }
};
