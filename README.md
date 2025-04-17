# WeChat Moment Layout Engine

This project provides a layout engine for WeChat Moments content, implemented in Go. It fetches moment data from a database and handles the layout of text and images onto dynamically generated A4-sized pages.

## Project Structure

```
.
├── golang/                  # Go implementation
│   ├── backend/
│   │   ├── server.go       # HTTP server implementation
│   │   ├── real_data.go    # Handles loading data from DB
│   │   ├── typeset.go      # Older types/layout logic (partially used?)
│   │   └── calculate/      # Core continuous layout engine logic
│   ├── frontend/
│   │   └── continuous.html # Continuous layout interface
│   └── main.go             # Application entry point
├── requirements.md         # Project requirements (may need update)
└── README.md              # This file
```

## Features

- Layout engine for WeChat Moments content fetched from a database.
- Support for text and image layout.
- Continuous layout with dynamic page creation.
- Automatic page breaking for long content.
- Responsive image scaling and positioning.
- Smart spacing algorithm between entries and elements.
- Advanced picture layout handling for 1-9+ pictures, including complex row/split layouts.
- Dynamic row-based layout for handling ultra-wide/tall images.
- Interactive zoom controls for continuous layout view.
- Groups moments by month, adding interstitial pages.

## Layout Algorithm

### Continuous Layout Engine
The Go implementation uses the `calculate.ContinuousLayoutEngine` for laying out WeChat Moments content across multiple pages:

1.  **Entry Spacing**: Configurable spacing between entries, elements within an entry, and images.
2.  **Page Management**: A4 page size (2480x3508 pixels @ 300 DPI, coordinates often handled at 72 DPI), dynamic page creation, configurable margins.
3.  **Layout Rules**: Handles time, text, and pictures. Includes advanced layout logic for various picture counts, aspect ratios, minimum heights, and complex splits.
4.  **Optimization**: Prevents awkward content splits, maintains visual hierarchy, and preserves image aspect ratios.

For more details, see `golang/README_golang.md`.

## Requirements

- Go 1.21 or later
- Access to the MySQL database specified in `golang/main.go`

## Installation & Running

1. Clone the repository.
2. Navigate to the Go implementation directory:
   ```bash
   cd golang
   ```
3. Ensure Go dependencies are downloaded:
   ```bash
   go mod tidy
   ```
4. Configure the database connection string (DSN) in `main.go` if necessary.
5. Run the server:
   ```bash
   go run main.go
   ```
6. Open your browser and navigate to `http://localhost:8888/` (default port).

## API Endpoints

- `GET /continuous-layout-real`: Fetches real moment data from the database, performs layout calculations, groups by month with interstitial pages, and returns the full layout as JSON (coordinates converted to 72 DPI).
  - Example: `http://localhost:8888/continuous-layout-real`

## License

This project is licensed under the MIT License. 