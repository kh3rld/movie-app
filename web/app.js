const searchInput = document.getElementById('searchInput');
const searchBtn = document.getElementById('searchBtn');
const resultsDiv = document.getElementById('results');

let debounceTimeout;
let currentPage = 1;
let totalPages = 1;
let currentQuery = '';
let currentGenre = '';
let currentType = 'movie';
let abortController = null;

const genreSelect = document.createElement('select');
genreSelect.id = 'genreSelect';
genreSelect.innerHTML = '<option value="">All Genres</option>';
resultsDiv.parentNode.insertBefore(genreSelect, resultsDiv);

// --- Watchlist UI ---
const watchlistBtn = document.createElement('button');
watchlistBtn.textContent = 'My Watchlist';
watchlistBtn.style.marginLeft = '1rem';
document.querySelector('h1').appendChild(watchlistBtn);

const trendingBtn = document.createElement('button');
trendingBtn.textContent = 'Trending';
trendingBtn.style.marginLeft = '1rem';
document.querySelector('h1').appendChild(trendingBtn);

const recBtn = document.createElement('button');
recBtn.textContent = 'Recommendations';
recBtn.style.marginLeft = '1rem';
document.querySelector('h1').appendChild(recBtn);

let watchlist = [];

function fetchWatchlist() {
    fetch('/api/watchlist')
        .then(res => res.json())
        .then(data => {
            watchlist = data.watchlist || [];
            renderWatchlist();
        });
}

function renderWatchlist() {
    resultsDiv.innerHTML = '<h2>Your Watchlist</h2>';
    if (!watchlist.length) {
        resultsDiv.innerHTML += '<p>Your watchlist is empty.</p>';
        return;
    }
    watchlist.forEach(item => {
        const div = document.createElement('div');
        div.className = 'result-item';
        const img = document.createElement('img');
        img.src = item.poster || 'https://via.placeholder.com/80x120?text=No+Image';
        img.alt = item.title;
        const info = document.createElement('div');
        info.className = 'info';
        info.innerHTML = `<strong>${item.title}</strong><div class='meta'>${item.type.toUpperCase()}${item.watched ? ' • Watched' : ''}</div>`;
        const removeBtn = document.createElement('button');
        removeBtn.textContent = 'Remove';
        removeBtn.onclick = () => {
            fetch(`/api/watchlist?id=${item.id}&type=${item.type}`, { method: 'DELETE' }).then(fetchWatchlist);
        };
        const markBtn = document.createElement('button');
        markBtn.textContent = 'Mark Watched';
        markBtn.disabled = item.watched;
        markBtn.onclick = () => {
            fetch(`/api/watchlist?id=${item.id}&type=${item.type}`, { method: 'PATCH' }).then(fetchWatchlist);
        };
        info.appendChild(removeBtn);
        info.appendChild(markBtn);
        div.appendChild(img);
        div.appendChild(info);
        resultsDiv.appendChild(div);
    });
}

watchlistBtn.onclick = fetchWatchlist;

// --- Trending UI ---
trendingBtn.onclick = () => {
    resultsDiv.innerHTML = '<h2>Trending</h2>';
    fetch('/api/trending?type=' + currentType)
        .then(res => res.json())
        .then(data => {
            renderResults(data.results);
        });
};

// --- Recommendations UI ---
recBtn.onclick = () => {
    resultsDiv.innerHTML = '<h2>Recommended for You</h2>';
    fetch('/api/recommendations')
        .then(res => res.json())
        .then(data => {
            renderResults(data.results);
        });
};

// --- Detail Modal ---
function showDetail(id, type) {
    fetch(`/api/detail?id=${id}&type=${type}`)
        .then(res => res.json())
        .then(data => {
            const modal = document.createElement('div');
            modal.style.position = 'fixed';
            modal.style.top = '0';
            modal.style.left = '0';
            modal.style.width = '100vw';
            modal.style.height = '100vh';
            modal.style.background = 'rgba(0,0,0,0.85)';
            modal.style.display = 'flex';
            modal.style.alignItems = 'center';
            modal.style.justifyContent = 'center';
            modal.style.zIndex = '1000';
            modal.onclick = e => { if (e.target === modal) modal.remove(); };
            const box = document.createElement('div');
            box.style.background = '#23272f';
            box.style.padding = '2rem';
            box.style.borderRadius = '12px';
            box.style.maxWidth = '420px';
            box.style.width = '90vw';
            box.style.color = '#fff';
            box.innerHTML = `
                <img src="${data.poster || 'https://via.placeholder.com/120x180?text=No+Image'}" style="width:120px;float:left;margin-right:1.5rem;border-radius:8px;">
                <h2>${data.title}</h2>
                <p><strong>Release:</strong> ${data.release_date || ''}</p>
                <p><strong>Cast:</strong> ${(data.cast || []).join(', ')}</p>
                <p><strong>Plot:</strong> ${data.plot || ''}</p>
                <p><strong>Ratings:</strong> ${Object.entries(data.ratings || {}).map(([k,v]) => `${k.toUpperCase()}: ${v}`).join(' | ')}</p>
                <button id="addToWatchlistBtn">Add to Watchlist</button>
                <button onclick="this.closest('div').parentNode.remove()">Close</button>
            `;
            modal.appendChild(box);
            document.body.appendChild(modal);
            document.getElementById('addToWatchlistBtn').onclick = () => {
                fetch('/api/watchlist', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        id: data.id,
                        type: currentType,
                        title: data.title,
                        poster: data.poster,
                        watched: false
                    })
                }).then(() => { modal.remove(); });
            };
        });
}

// --- Enhance renderResults to support detail modal and add-to-watchlist ---
function renderResults(items) {
    resultsDiv.innerHTML = '';
    if (!items || items.length === 0) {
        resultsDiv.innerHTML = '<p>No results found.</p>';
        return;
    }
    items.forEach(item => {
        const div = document.createElement('div');
        div.className = 'result-item';
        const img = document.createElement('img');
        img.src = item.poster || 'https://via.placeholder.com/80x120?text=No+Image';
        img.alt = item.title;
        img.style.cursor = 'pointer';
        img.onclick = () => showDetail(item.id, currentType);
        const info = document.createElement('div');
        info.className = 'info';
        info.innerHTML = `<strong>${item.title}</strong><div class='meta'>${item.year || ''}</div>`;
        const addBtn = document.createElement('button');
        addBtn.textContent = 'Add to Watchlist';
        addBtn.onclick = () => {
            fetch('/api/watchlist', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    id: item.id,
                    type: currentType,
                    title: item.title,
                    poster: item.poster,
                    watched: false
                })
            });
        };
        info.appendChild(addBtn);
        div.appendChild(img);
        div.appendChild(info);
        resultsDiv.appendChild(div);
    });
}

function renderPagination() {
    const pagDiv = document.createElement('div');
    pagDiv.className = 'pagination';
    if (totalPages <= 1) return;
    pagDiv.innerHTML = `
        <button id="prevPage" ${currentPage === 1 ? 'disabled' : ''}>&laquo; Prev</button>
        <span>Page ${currentPage} of ${totalPages}</span>
        <button id="nextPage" ${currentPage === totalPages ? 'disabled' : ''}>Next &raquo;</button>
    `;
    resultsDiv.appendChild(pagDiv);
    document.getElementById('prevPage').onclick = () => changePage(currentPage - 1);
    document.getElementById('nextPage').onclick = () => changePage(currentPage + 1);
}

function renderLoading() {
    resultsDiv.innerHTML = '<p>Loading...</p>';
}

function renderError(msg) {
    resultsDiv.innerHTML = `<p style="color:#e50914;">${msg}</p>`;
}

function fetchGenres() {
    fetch('/api/genres?type=' + currentType)
        .then(res => res.json())
        .then(data => {
            if (data.genres) {
                genreSelect.innerHTML = '<option value="">All Genres</option>';
                data.genres.forEach(g => {
                    const opt = document.createElement('option');
                    opt.value = g.id;
                    opt.textContent = g.name;
                    genreSelect.appendChild(opt);
                });
            }
        });
}

function searchMovies(page = 1) {
    const query = searchInput.value.trim();
    currentQuery = query;
    currentPage = page;
    currentGenre = genreSelect.value;
    if (abortController) abortController.abort();
    abortController = new AbortController();
    if (!query) {
        resultsDiv.innerHTML = '';
        return;
    }
    renderLoading();
    let url = `/api/search?q=${encodeURIComponent(query)}&type=${currentType}&page=${page}`;
    if (currentGenre) url += `&genre=${currentGenre}`;
    fetch(url, { signal: abortController.signal })
        .then(res => res.json())
        .then(data => {
            if (data.error) {
                renderError(data.error);
                return;
            }
            renderResults(data.results);
            totalPages = data.total_pages || 1;
            renderPagination();
        })
        .catch(e => {
            if (e.name !== 'AbortError') renderError('Error fetching results.');
        });
}

function changePage(page) {
    if (page < 1 || page > totalPages) return;
    searchMovies(page);
}

searchBtn.addEventListener('click', () => searchMovies(1));
searchInput.addEventListener('input', () => {
    clearTimeout(debounceTimeout);
    debounceTimeout = setTimeout(() => searchMovies(1), 350);
});
genreSelect.addEventListener('change', () => searchMovies(1));

// Initial genre fetch
fetchGenres();