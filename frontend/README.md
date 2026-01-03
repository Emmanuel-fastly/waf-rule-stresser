# WAF Rate Limit Tester - Frontend

Modern React frontend for testing NGWAF rules & rate-limit rules.

## Tech Stack

- **Vite 5.x** - Fast build tool
- **React 18.3** - UI library
- **Tailwind CSS 3.x** - Utility-first CSS framework

## Getting Started

### Install Dependencies

```bash
npm install
```

### Run Development Server

```bash
npm run dev
```

The frontend will start on [http://localhost:3000](http://localhost:3000)

### Build for Production

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

## Features

### Test Configuration

- Target URL input
- Request count and duration
- HTTP method selection (GET, POST, PUT, DELETE, PATCH)
- Test mode (Baseline or Burst)
- Traffic type (Normal or Attack)
- User-agent configuration
- Custom headers support
- Request body for POST/PUT/PATCH

### Results Display

- Summary cards (Total, Success, Errors, Blocked)
- Performance metrics (Avg, Min, Max, P50, P95, P99)
- Rate limit detection alert
- Request log with pagination
- Filtering by status
- Detailed request modal

### Export

- Export to JSON format
- Export to CSV format
- Files saved to `../exports/` directory

## API Integration

The frontend communicates with the Go backend via these endpoints:

- `POST /api/test/start` - Start a new test
- `POST /api/test/export` - Export results
- `GET /api/exports/list` - List all exports

Vite's proxy automatically forwards `/api/*` requests to `http://localhost:8080`.

## Troubleshooting

### Port 3000 already in use

Change the port in `vite.config.js`:

```js
server: {
  port: 3001, // or any available port
}
```

### Cannot connect to backend

1. Ensure backend is running on port 8080
2. Check browser console for CORS errors
3. Verify Vite proxy configuration

### Tailwind styles not loading

```bash
npm install
npm run dev
```
