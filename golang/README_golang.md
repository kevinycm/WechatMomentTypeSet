# WechatMomentTypeSet Go Implementation

This is the Go implementation of the WechatMomentTypeSet project, which provides a layout engine for WeChat Moments content.

## Project Structure

```
golang/
├── backend/
│   ├── server.go             # HTTP server implementation
│   ├── real_data.go          # Handles loading data from DB
│   ├── typeset.go            # Older types/layout logic (partially used?)
│   └── calculate/            # Core continuous layout engine logic
├── frontend/
│   └── continuous.html      # Continuous layout interface
└── main.go                  # Application entry point
```

## Features

- Layout engine for WeChat Moments content from a database
- Support for text and image layout
- Continuous layout with dynamic page creation
- Automatic page breaking for long content
- Responsive image scaling and positioning
- Smart spacing algorithm between entries and elements
- Handles various picture counts and layouts (including complex splits)
- Interactive zoom controls for continuous layout view

## Layout Algorithm

### Continuous Layout Engine
The continuous layout engine (`calculate.ContinuousLayoutEngine`) implements a sophisticated algorithm for laying out WeChat Moments content across multiple pages:

1. **Entry Spacing**:
   - Configurable spacing between entries, elements within an entry, and images.

2. **Page Management**:
   - A4 page size (2480x3508 pixels @ 300 DPI, coordinates often handled at 72 DPI).
   - Dynamic page creation based on content flow.
   - Configurable margins.

3. **Layout Rules**:
   - Handles time, text, and pictures.
   - Advanced picture layout handling for 1-9+ pictures, including aspect ratio considerations, minimum height checks, and complex row/split layouts.
   - Dynamic row-based layout for handling ultra-wide/tall images.
   - Smart pagination to avoid orphaned elements.

4. **Optimization Features**:
   - Prevents splitting of time and first content element.
   - Maintains visual hierarchy with consistent spacing.
   - Preserves image aspect ratios while maximizing space usage.
   - Handles various content combinations (time+text, time+images, time+text+images).

## Requirements

- Go 1.21 or later
- Access to the MySQL database specified in `main.go`

## Installation

1. Clone the repository
2. Navigate to the golang directory:
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

## Usage

1. Start the server (default port: 8888). Ensure the database is accessible.
2. Open your browser and navigate to `http://localhost:8888/`.
3. The continuous layout view will load, displaying moments fetched from the database.
4. Use the interface controls to:
   - Adjust zoom level.
   - Navigate pages (if applicable, depending on frontend implementation).

## API Endpoints

- `GET /continuous-layout-real`: Fetches real moment data from the database, performs layout calculations, groups by month with interstitial pages, and returns the full layout as JSON (coordinates converted to 72 DPI).
  - Example: `http://localhost:8888/continuous-layout-real`

## License

This project is licensed under the MIT License. 