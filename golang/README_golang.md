# WechatMomentTypeSet Go Implementation

This is the Go implementation of the WechatMomentTypeSet project, which provides a layout engine for WeChat Moments content.

## Project Structure

```
golang/
├── backend/
│   ├── server.go             # HTTP server implementation
│   ├── typeset.go            # Single page layout engine
│   └── continuous_layout.go  # Continuous layout engine
├── frontend/
│   ├── index.html           # Single page layout interface
│   └── continuous.html      # Continuous layout interface
└── main.go                  # Application entry point
```

## Features

- Layout engine for WeChat Moments content
- Support for text and image layout
- Two layout modes:
  - Single page layout with fixed dimensions
  - Continuous layout with dynamic page creation
- Automatic page breaking for long content
- Responsive image scaling and positioning
- Smart spacing algorithm between entries and elements
- Interactive zoom controls for continuous layout view

## Layout Algorithm

### Continuous Layout Engine
The continuous layout engine (`ContinuousLayoutEngine`) implements a sophisticated algorithm for laying out WeChat Moments content across multiple pages:

1. **Entry Spacing**:
   - 150px spacing between entries
   - 30px spacing between elements within an entry
   - 20px spacing between images

2. **Page Management**:
   - A4 page size (2480x3508 pixels)
   - Dynamic page creation based on content
   - Margins: Top/Bottom 100/160px, Left/Right 100px

3. **Layout Rules**:
   - Time area: Fixed height of 104px
   - Text: Font size 50px, line height 75px
   - Single images: Default height of 1695px
   - Multiple images: 9-grid layout with aspect ratio preservation
   - Smart pagination to avoid orphaned elements

4. **Optimization Features**:
   - Prevents splitting of time and first content element
   - Maintains visual hierarchy with consistent spacing
   - Preserves image aspect ratios while maximizing space usage
   - Handles various content combinations (time+text, time+images, time+text+images)

## Requirements

- Go 1.16 or later

## Installation

1. Clone the repository
2. Navigate to the golang directory:
   ```bash
   cd golang
   ```
3. Run the server:
   ```bash
   go run main.go
   ```

## Usage

1. Start the server (default port: 8888)
2. Open your browser and navigate to:
   - Single page layout: `http://localhost:8888`
   - Continuous layout: `http://localhost:8888/continuous`
3. Use the interface controls to:
   - Load sample data
   - Adjust zoom level (continuous layout)
   - Clear current layout

## API Endpoints

- `GET /layout/{id}`: Get single page layout for a specific test case
- `GET /continuous-layout-sample`: Get sample data for continuous layout
  - Example: `http://localhost:8888/continuous-layout-sample`

## Test Cases

The project includes several test cases with different combinations of text and images:

- ID 121-129: Various combinations of text and images
- ID 130: Long text only
- ID 131: Images only

## License

This project is licensed under the MIT License. 