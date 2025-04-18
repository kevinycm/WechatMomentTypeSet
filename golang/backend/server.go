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
	"wechatmomenttypeset/backend/waterfall"
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
	http.HandleFunc("/continuous-layout-real", s.handleContinuousLayoutReal)

	// Start server
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Server starting on port %d...", s.port)
	return http.ListenAndServe(addr, nil)
}

// getYearMonthKey returns a sortable key for year-month grouping
func getYearMonthKey(timeStr string) string {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d-%02d", t.Year(), t.Month())
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
	yearMonthGroups := make(map[string][]waterfall.Entry)
	yearMonthKeys := make([]string, 0)
	for _, id := range ids {
		element := RealData[id]
		// Convert models.NewMoment (represented by Element here) to calculate.Entry
		entry := waterfall.Entry{
			ID:       int64(element.ID),
			Time:     element.Time,
			Text:     element.Text,
			Pictures: element.Pictures, // Directly use the existing Pictures field
		}
		yearMonthKey := getYearMonthKey(element.Time)
		if yearMonthKey == "" {
			continue
		}

		if _, exists := yearMonthGroups[yearMonthKey]; !exists {
			yearMonthKeys = append(yearMonthKeys, yearMonthKey)
		}

		yearMonthGroups[yearMonthKey] = append(yearMonthGroups[yearMonthKey], entry)
	}

	// 按年月降序排序
	sort.Sort(sort.Reverse(sort.StringSlice(yearMonthKeys)))

	// 处理每个年月组的数据
	var allPages []waterfall.ContinuousLayoutPage
	pageNumber := 1

	for _, yearMonthKey := range yearMonthKeys {
		entries := yearMonthGroups[yearMonthKey]
		if len(entries) == 0 {
			continue
		}

		// 添加插页
		parts := strings.Split(yearMonthKey, "-")
		if len(parts) != 2 {
			continue
		}
		year, _ := strconv.Atoi(parts[0])
		month, _ := strconv.Atoi(parts[1])
		yearMonth := fmt.Sprintf("%d年%d月", year, month)

		insertPage := waterfall.ContinuousLayoutPage{
			Page:      pageNumber,
			IsInsert:  true,
			YearMonth: yearMonth,
			Entries:   []waterfall.PageEntry{},
		}
		allPages = append(allPages, insertPage)
		pageNumber++

		// 处理该年月的条目
		engine := waterfall.NewContinuousLayoutEngine(entries)
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

	// 将所有页面的坐标转换为72DPI
	for i := range allPages {
		allPages[i] = convertPageTo72DPI(allPages[i])
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pages": allPages,
	})
}

// convertTo72DPI converts coordinates from 300DPI to 72DPI
func convertTo72DPI(value float64) float64 {
	return value * 72 / 300
}

// convertAreaTo72DPI converts an area's coordinates from 300DPI to 72DPI
func convertAreaTo72DPI(area [][]float64) [][]float64 {
	if len(area) != 2 || len(area[0]) != 2 || len(area[1]) != 2 {
		return area
	}
	return [][]float64{
		{convertTo72DPI(area[0][0]), convertTo72DPI(area[0][1])},
		{convertTo72DPI(area[1][0]), convertTo72DPI(area[1][1])},
	}
}

// convertPageTo72DPI converts all coordinates in a page from 300DPI to 72DPI
func convertPageTo72DPI(page waterfall.ContinuousLayoutPage) waterfall.ContinuousLayoutPage {
	// 不要转换页码，保持原样
	// page.Page = int(convertTo72DPI(float64(page.Page)))

	// Convert each entry
	for i := range page.Entries {
		entry := &page.Entries[i]

		// Convert time area
		if entry.TimeArea != nil {
			entry.TimeArea = convertAreaTo72DPI(entry.TimeArea)
		}

		// Convert text areas
		for j := range entry.TextAreas {
			entry.TextAreas[j] = convertAreaTo72DPI(entry.TextAreas[j])
		}

		// Convert pictures
		for j := range entry.Pictures {
			entry.Pictures[j].Area = convertAreaTo72DPI(entry.Pictures[j].Area)
		}
	}

	return page
}
