let currentView = 'market';
let fetchInterval = null;

async function fetchLists() {
    if (currentView !== 'market') return;
    try {
        const response = await fetch('/api/shopping-lists');
        if (!response.ok) throw new Error('Network response was not ok');
        const data = await response.json();
        renderLists(data);
    } catch (error) {
        console.error('Error fetching lists:', error);
        const status = document.getElementById('market-status');
        status.innerText = 'Connection Error';
        status.style.color = '#ef4444';
    }
}

async function fetchStatus() {
    try {
        const response = await fetch('/api/status');
        if (!response.ok) return;
        const data = await response.json();
        updateProgress(data);
    } catch (error) {
        console.error('Error fetching status:', error);
    }
}

function updateProgress(data) {
    const container = document.getElementById('progress-container');
    const bar = document.getElementById('progress-bar');
    const text = document.getElementById('progress-text');
    const statusLabel = document.getElementById('market-status');

    if (data.total > 0 && data.current < data.total) { // Active progress
        container.classList.remove('hidden');
        const percentage = (data.current / data.total) * 100;
        bar.style.setProperty('--progress', `${percentage}%`);
        text.innerText = `Updating... ${data.current} / ${data.total}`;

        statusLabel.innerText = 'Fetching Data...';
        statusLabel.style.color = '#facc15'; // Yellow
    } else {
        container.classList.add('hidden');
        // Only set to Live Updates if we successfully rendered lists recently
        if (statusLabel.innerText !== 'Connection Error') {
            statusLabel.innerText = 'Live Updates';
            statusLabel.style.color = '#4ade80'; // Green
        }
    }
}


function renderLists(lists) {
    const container = document.getElementById('shopping-lists');

    // Don't clear if null, might be just a failed fetch during update
    if (!lists) return;

    container.innerHTML = '';

    if (lists.length === 0) {
        container.innerHTML = '<div style="grid-column: 1/-1; text-align: center; color: var(--text-secondary);">Waiting for market data...</div>';
        return;
    }

    lists.forEach(list => {
        const card = document.createElement('div');
        card.className = 'card';
        card.innerHTML = `
            <div class="card-header">
                <div class="route">${list.From} &rarr; ${list.To}</div>
                <div class="jumps">${list.Jumps} jumps</div>
            </div>
            <div class="stat-row">
                <span class="label">Profit</span>
                <span class="value profit">${list.Profit}</span>
            </div>
             <div class="stat-row">
                <span class="label">Investment</span>
                <span class="value">${list.Investment}</span>
            </div>
             <div class="stat-row">
                <span class="label">Volume</span>
                <span class="value">${list.Volume} m³</span>
            </div>
            <div class="stat-row">
                <span class="label">ROI</span>
                <span class="value">${list.ROI}%</span>
            </div>
             <div class="items-container">
                <div class="items-header">Items</div>
                <table class="items-table">
                    <thead>
                        <tr>
                            <th>Qty</th>
                            <th>Name</th>
                            <th>Buy</th>
                            <th>Sell</th>
                            <th>Profit</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${list.Items ? list.Items.map(item => `
                            <tr>
                                <td>${item.Quantity}</td>
                                <td>${item.Name}</td>
                                <td>${item.BuyPrice}</td>
                                <td>${item.SellPrice}</td>
                                <td class="profit">${item.Profit}</td>
                            </tr>
                        `).join('') : ''}
                    </tbody>
                </table>
            </div>
        `;
        container.appendChild(card);
    });
}

function switchView(view) {
    currentView = view;

    // Update Nav
    document.querySelectorAll('.nav-link').forEach(l => l.classList.remove('active'));
    document.querySelector(`.nav-link[onclick="switchView('${view}')"]`).classList.add('active');

    // Update Views
    const marketView = document.getElementById('market-view');
    const itemsView = document.getElementById('items-view');

    if (view === 'market') {
        marketView.classList.remove('hidden');
        itemsView.classList.add('hidden');
        fetchLists(); // Fetch immediately
    } else {
        marketView.classList.add('hidden');
        itemsView.classList.remove('hidden');
    }
}

var searchTimeout;
function handleSearch(event) {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
        const query = event.target.value;
        if (query.length > 2) {
            searchItems(query);
        }
    }, 300);
}

async function searchItems(query) {
    try {
        const response = await fetch(`/api/items?q=${encodeURIComponent(query)}`);
        const items = await response.json();
        renderItems(items);
    } catch (e) {
        console.error(e);
    }
}

function renderItems(items) {
    const container = document.getElementById('items-list');
    container.innerHTML = '';

    if (!items || items.length === 0) {
        container.innerHTML = '<div style="grid-column: 1/-1; text-align: center; color: var(--text-secondary);">No items found.</div>';
        return;
    }

    items.forEach(item => {
        const card = document.createElement('div');
        card.className = 'card';
        card.innerHTML = `
            <div class="card-header">
                <div class="route">${item.name}</div>
            </div>
            <div class="stat-row">
                <span class="label">Volume</span>
                <span class="value">${item.volume} m³</span>
            </div>
        `;
        container.appendChild(card);
    });
}

// Initial fetch and periodic update
fetchLists();
setInterval(fetchLists, 2000);
setInterval(fetchStatus, 500);
