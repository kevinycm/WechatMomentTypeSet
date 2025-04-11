# WechatMomentTypeSet Go Implementation

This is the Go implementation of the WechatMomentTypeSet project, which provides a layout engine for WeChat Moments content.

## Project Structure

```
golang/
├── backend/
│   ├── server.go    # HTTP server implementation
│   └── typeset.go   # Layout engine implementation
├── frontend/
│   └── index.html   # Frontend interface
└── main.go          # Application entry point
```

## Features

- Layout engine for WeChat Moments content
- Support for text and image layout
- Automatic page breaking for long content
- Responsive image scaling and positioning

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
2. Open your browser and navigate to `http://localhost:8888`
3. Enter a test case ID (121-131) to view the layout

## Test Cases

The project includes several test cases with different combinations of text and images:

- ID 121-129: Various combinations of text and images
- ID 130: Long text only
- ID 131: Images only

## API Endpoints

- `GET /layout/{id}`: Get layout for a specific test case
  - Example: `http://localhost:8888/layout/121`

## License

This project is licensed under the MIT License. 