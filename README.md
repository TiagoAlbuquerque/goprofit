# GoProfit

**Market arbitrage tool for EVE Online** - Find profitable trade routes by analyzing buy and sell orders across different regions.

## Features

- 🔍 **Real-time market analysis** - Fetches market data from EVE Online ESI API
- 📊 **Smart deal detection** - Identifies profitable arbitrage opportunities
- 🚀 **Optimized shopping lists** - Calculates best items to trade based on cargo space and ISK
- 🌐 **Web interface** - Modern dark-themed dashboard with live updates
- ⚡ **High performance** - Concurrent processing with connection pooling

## Quick Start

```bash
# Build
cd src
go build -o ../goprofit.exe

# Run
cd ..
./goprofit.exe
```

Access the web interface at `http://localhost:8080`

## Configuration

Edit `data_conf.json` to customize:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `cargo` | Ship cargo capacity (m³) | 122.4 |
| `max_invest` | Maximum ISK per shopping list | 1.5B |
| `min_pm3` | Minimum profit per m³ | 100K |
| `tax` | Sales tax rate | 3.6% |

## Architecture

```
goprofit/
├── src/                    # Go source code
│   ├── controller/         # Main orchestration
│   ├── deals/              # Arbitrage logic
│   ├── shoppingLists/      # Shopping list optimization
│   ├── items/              # Item cache
│   ├── orders/             # Market orders
│   ├── regions/            # Region management
│   ├── locations/          # System/station info
│   ├── server/             # HTTP server
│   └── utils/              # Utilities & worker pool
├── public/                 # Web frontend
└── data_*.json             # Cached data files
```

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | Web interface |
| `GET /api/shopping-lists` | Top 9 profitable routes (JSON) |
| `GET /api/status` | Current fetch progress |
| `GET /api/items?q=` | Search items |

## Performance

- **Worker pool** with 50 concurrent HTTP requests
- **Connection pooling** for HTTP client reuse
- **Lock-free reads** on stable data structures
- **O(1) deal deduplication** using hash maps

## Dependencies

- [cheggaaa/pb](https://github.com/cheggaaa/pb) - Progress bar

## License

MIT
