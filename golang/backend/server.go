package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Server represents the HTTP server
type Server struct {
	port     int
	basePath string
	db       string // Add database connection string
}

// NewServer creates a new server instance
func NewServer(port int, basePath string, dbDSN string) *Server {
	return &Server{
		port:     port,
		basePath: basePath,
		db:       dbDSN,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Initialize real data
	if err := InitRealData(s.db); err != nil {
		log.Printf("Warning: Failed to initialize real data: %v", err)
	}

	// Serve static files
	fs := http.FileServer(http.Dir(filepath.Join(s.basePath, "frontend")))
	http.Handle("/", fs)

	// API endpoints
	http.HandleFunc("/layout/", s.handleLayout)
	http.HandleFunc("/continuous-layout-sample", s.handleContinuousLayoutSample)
	http.HandleFunc("/continuous-layout-real", s.handleContinuousLayoutReal)

	// Start server
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Server starting on port %d...", s.port)
	return http.ListenAndServe(addr, nil)
}

// handleLayout handles the layout API endpoint
func (s *Server) handleLayout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 3 {
		http.Error(w, "Invalid path format", http.StatusBadRequest)
		return
	}

	idStr := pathParts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	// Get test case
	testCase, ok := SampleData[id]
	if !ok {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Process layout
	engine := NewLayoutEngine(testCase)
	result, err := engine.ProcessEntry()
	if err != nil {
		http.Error(w, "Error processing layout", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Encode and send response
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Entry represents a single moment entry with time, text and pictures
type Entry struct {
	Time     string    `json:"time"`
	Text     string    `json:"text"`
	Pictures []Picture `json:"pictures"`
}

// formatTime converts time string from "2025-03-20 12:30:15" to "2025年3月20日 12:30"
func formatTime(timeStr string) string {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return timeStr // Return original string if parsing fails
	}
	return fmt.Sprintf("%d年%d月%d日 %02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func (s *Server) handleContinuousLayoutSample(w http.ResponseWriter, r *http.Request) {
	// Convert all test cases to entries and sort by ID
	entries := make([]Entry, 0, len(SampleData))
	ids := make([]int, 0, len(SampleData))

	// First collect all IDs
	for id := range SampleData {
		ids = append(ids, id)
	}

	// Sort IDs
	sort.Ints(ids)

	// Create entries in sorted order with original time strings
	originalTimes := make(map[int]string)
	for _, id := range ids {
		testCase := SampleData[id]
		originalTimes[id] = testCase.Time
		entry := Entry{
			Time:     formatTime(testCase.Time),
			Text:     testCase.Text,
			Pictures: testCase.Pictures,
		}
		entries = append(entries, entry)
	}

	// Create layout engine and process entries
	engine := NewContinuousLayoutEngine(entries)
	pages, err := engine.ProcessEntries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format the response to match frontend expectations
	formattedPages := make([]map[string]interface{}, len(pages))
	for i, page := range pages {
		// Create text areas and texts slices
		textAreas := make([][][]float64, len(page.TextAreas))
		texts := make([]string, len(page.TextAreas))
		for j, area := range page.TextAreas {
			textAreas[j] = [][]float64{
				{area[0][0], area[0][1]},
				{area[1][0], area[1][1]},
			}
			texts[j] = page.Texts[j]
		}

		// Create pictures slice
		pictures := make([]map[string]interface{}, len(page.Pictures))
		for j, pic := range page.Pictures {
			pictures[j] = map[string]interface{}{
				"url": pic.URL,
				"area": [][]float64{
					{pic.Area[0][0], pic.Area[0][1]},
					{pic.Area[1][0], pic.Area[1][1]},
				},
			}
		}

		// Create time area with safety checks
		var timeArea [][]float64
		if len(page.TimeArea) >= 2 && len(page.TimeArea[0]) >= 2 && len(page.TimeArea[1]) >= 2 {
			timeArea = [][]float64{
				{page.TimeArea[0][0], page.TimeArea[0][1]},
				{page.TimeArea[1][0], page.TimeArea[1][1]},
			}
		} else {
			// Default time area if not properly set
			timeArea = [][]float64{
				{100, 100},
				{2380, 204},
			}
		}

		// Get the original time string from the page time by parsing it back
		yearMonth := ""
		if t, err := time.Parse("2006年1月2日 15:04", page.Time); err == nil {
			yearMonth = fmt.Sprintf("%d.%02d", t.Year(), t.Month())
		}

		formattedPage := map[string]interface{}{
			"time":       page.Time,
			"time_area":  timeArea,
			"text_areas": textAreas,
			"texts":      texts,
			"pictures":   pictures,
			"year_month": yearMonth,
		}

		formattedPages[i] = formattedPage
	}

	// Return the result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pages": formattedPages,
	})
}

func (s *Server) handleContinuousLayoutReal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get all IDs and sort by release_time in descending order
	var ids []int
	for id := range RealData {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		timeI, _ := time.Parse("2006-01-02 15:04:05", RealData[ids[i]].Time)
		timeJ, _ := time.Parse("2006-01-02 15:04:05", RealData[ids[j]].Time)
		return timeI.After(timeJ)
	})

	// Convert all entries to the required format
	entries := make([]Entry, 0, len(ids))
	for _, id := range ids {
		testCase := RealData[id]
		entry := Entry{
			Time:     formatTime(testCase.Time),
			Text:     testCase.Text,
			Pictures: testCase.Pictures,
		}
		entries = append(entries, entry)
	}

	// Create layout engine and process entries
	engine := NewContinuousLayoutEngine(entries)
	pages, err := engine.ProcessEntries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Format the response to match frontend expectations
	formattedPages := make([]map[string]interface{}, len(pages))
	for i, page := range pages {
		// Create text areas and texts slices
		textAreas := make([][][]float64, len(page.TextAreas))
		texts := make([]string, len(page.TextAreas))
		for j, area := range page.TextAreas {
			textAreas[j] = [][]float64{
				{area[0][0], area[0][1]},
				{area[1][0], area[1][1]},
			}
			texts[j] = page.Texts[j]
		}

		// Create pictures slice
		pictures := make([]map[string]interface{}, len(page.Pictures))
		for j, pic := range page.Pictures {
			pictures[j] = map[string]interface{}{
				"url": pic.URL,
				"area": [][]float64{
					{pic.Area[0][0], pic.Area[0][1]},
					{pic.Area[1][0], pic.Area[1][1]},
				},
			}
		}

		// Get the year and month from the page time
		yearMonth := ""
		if t, err := time.Parse("2006年1月2日 15:04", page.Time); err == nil {
			yearMonth = fmt.Sprintf("%d.%02d", t.Year(), t.Month())
		}

		formattedPage := map[string]interface{}{
			"time":        page.Time,
			"time_area":   page.TimeArea,
			"text_areas":  textAreas,
			"texts":       texts,
			"pictures":    pictures,
			"year_month":  yearMonth,
			"page_number": i + 1,
		}

		formattedPages[i] = formattedPage
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pages": formattedPages,
	})
}
