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
        const info = document.createElement('div');
        info.className = 'info';
        info.innerHTML = `<strong>${item.title}</strong><div class='meta'>${item.year || ''}</div>`;
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