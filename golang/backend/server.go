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

// formatYearMonth extracts and formats year-month from time string
func formatYearMonth(timeStr string) string {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return "" // Return empty string if parsing fails
	}
	return fmt.Sprintf("%d年%d月", t.Year(), t.Month())
}

// getYearMonthKey returns a sortable key for year-month grouping
func getYearMonthKey(timeStr string) string {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d-%02d", t.Year(), t.Month())
}

func (s *Server) handleContinuousLayoutSample(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Convert TestCase to Entry
	var entries []Entry
	for _, testCase := range SampleData {
		entries = append(entries, Entry{
			Time:     testCase.Time,
			Text:     testCase.Text,
			Pictures: testCase.Pictures,
		})
	}

	engine := NewContinuousLayoutEngine(entries)
	pages, err := engine.ProcessEntries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pages": pages,
	})
}

func (s *Server) handleContinuousLayoutReal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取所有ID并按降序排序
	var ids []int
	for id := range RealData {
		ids = append(ids, id)
	}

	// 按时间降序排序
	sort.Slice(ids, func(i, j int) bool {
		timeI, errI := time.Parse("2006-01-02 15:04:05", RealData[ids[i]].Time)
		timeJ, errJ := time.Parse("2006-01-02 15:04:05", RealData[ids[j]].Time)

		// 如果解析出错，将其放到最后
		if errI != nil {
			return false
		}
		if errJ != nil {
			return true
		}

		return timeI.After(timeJ)
	})

	// 按年月分组
	yearMonthGroups := make(map[string][]Entry)
	yearMonthKeys := make([]string, 0)
	for _, id := range ids {
		testCase := RealData[id]
		yearMonthKey := getYearMonthKey(testCase.Time)
		if yearMonthKey == "" {
			continue
		}

		if _, exists := yearMonthGroups[yearMonthKey]; !exists {
			yearMonthKeys = append(yearMonthKeys, yearMonthKey)
		}

		entry := Entry{
			Time:     formatTime(testCase.Time),
			Text:     testCase.Text,
			Pictures: testCase.Pictures,
		}
		yearMonthGroups[yearMonthKey] = append(yearMonthGroups[yearMonthKey], entry)
	}

	// 按年月降序排序
	sort.Sort(sort.Reverse(sort.StringSlice(yearMonthKeys)))

	// 处理每个年月组的数据
	var allPages []ContinuousLayoutPage
	pageNumber := 1

	for _, yearMonthKey := range yearMonthKeys {
		entries := yearMonthGroups[yearMonthKey]
		if len(entries) == 0 {
			continue
		}

		// 添加插页
		// 使用当前年月组的时间来格式化年月
		// 从yearMonthKey中提取年月信息
		parts := strings.Split(yearMonthKey, "-")
		if len(parts) != 2 {
			continue
		}
		year, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		yearMonth := fmt.Sprintf("%d年%d月", year, month)

		insertPage := ContinuousLayoutPage{
			Page:      pageNumber,
			IsInsert:  true,
			YearMonth: yearMonth,
			Entries:   []PageEntry{},
		}
		allPages = append(allPages, insertPage)
		pageNumber++

		// 处理该年月的条目
		engine := NewContinuousLayoutEngine(entries)
		pages, err := engine.ProcessEntries()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 更新页码并添加年月信息
		for i := range pages {
			pages[i].Page = pageNumber
			pages[i].YearMonth = yearMonth
			pageNumber++
		}

		allPages = append(allPages, pages...)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pages": allPages,
	})
}
