# WeChat Moment Layout Engine

This project provides a layout engine for WeChat Moments content, implemented in both Python and Go. It handles the layout of text and images on A4-sized pages according to specific requirements.

## Project Structure

```
.
├── python/                  # Python implementation
│   ├── backend/
│   │   ├── server.py       # Tornado web server
│   │   └── typeset.py      # Layout engine
│   └── frontend/
│       └── index.html      # Frontend interface
├── golang/                  # Go implementation
│   ├── backend/
│   │   ├── server.go       # HTTP server
│   │   └── typeset.go      # Layout engine
│   ├── frontend/
│   │   └── index.html      # Frontend interface
│   └── main.go             # Application entry point
├── requirements.md         # Project requirements
└── README.md              # This file
```

## Features

- Layout engine for WeChat Moments content
- Support for text and image layout
- Automatic page breaking for long content
- Responsive image scaling and positioning
- 9-grid image layout algorithm
- Multi-page support
- Web interface for visualization

## Python Implementation

The Python implementation uses Tornado for the web server and provides a clean, object-oriented layout engine.

### Requirements
- Python 3.8+
- Tornado

### Running
```bash
cd python
python backend/server.py
```

## Go Implementation

The Go implementation provides a high-performance alternative with the same functionality.

### Requirements
- Go 1.16+

### Running
```bash
cd golang
go run main.go
```

## Common Features

Both implementations support:
1. Three types of content:
   - Time and text only
   - Time and images (1-9)
   - Time, text, and images
2. A4 page layout (3508 x 2480 pixels)
3. Automatic page breaking
4. Image scaling with aspect ratio preservation
5. 9-grid image layout algorithm
6. Web interface for visualization

## API Endpoints

Both implementations provide:
- `GET /layout/{id}`: Get layout for a specific test case
  - Example: `http://localhost:8888/layout/121`

## Test Cases

The project includes test cases for various scenarios:
- ID 121-129: Various combinations of text and images
- ID 130: Long text only
- ID 131: Images only

## License

This project is licensed under the MIT License. 