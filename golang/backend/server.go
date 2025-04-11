package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

// Server represents the HTTP server
type Server struct {
	port     int
	basePath string
}

// NewServer creates a new server instance
func NewServer(port int, basePath string) *Server {
	return &Server{
		port:     port,
		basePath: basePath,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Serve static files
	fs := http.FileServer(http.Dir(filepath.Join(s.basePath, "frontend")))
	http.Handle("/", fs)

	// API endpoints
	http.HandleFunc("/layout/", s.handleLayout)

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
	testCase, ok := sampleData[id]
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

// Sample data
var sampleData = map[int]TestCase{
	121: {
		ID:   121,
		Time: "2025-03-20 12:30:15",
		Text: "这是一个需要跨多页的长文本aaaaaaaaaaaf发生的方式发发发发放水阀代发沙发沙发撒发撒发达说法都发发撒打发萨法沙发沙发发多少范德萨范德萨发沙发沙发的阿范德萨发生发生发生...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n",
		Pictures: []Picture{
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"},
			{Width: 1620, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkUp1vXcxc4J0kVhFwxPRnkLNqvFqER99R2FC0BDmCYx8"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnorDmccicxxRQhREKbJjQGuQocfNrSvyvrytnoGHcwSWmDHFQTMbdDj4"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9MyDLXiaZ7YtL7DgIDWPqNMS8odq91EdX586jQx2UDvlo"},
			{Width: 1080, Height: 1440, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9rhKE57uzPHBIt4ldv1btOMa0ibW5zxlKRYXQaQMico61Q"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9vlFfMRibmJ3m0nJR94FaibxBFEiaw3Dq3UM3fs7cD1ReqA"},
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFU0ic58W0HarI0kZFBdia9cXEwibbzG2RqEYr0bmiaYMwJ7E6OwMp6haQkk"},
			{Width: 1080, Height: 1440, URL: "https://img.diandibianji.com/8u9KefYVGSVTCWiabXdfPI8xRknynoVGZ1rC9M13tRXSU3A1libL6xT8eTkbrtRcRtXOR2C33FSU8"},
			{Width: 1440, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSVTCWiabXdfPIib3W7TLFloYBj0w7WOWtxxawey8bHgg4Tyqzrkwre1V8dNA7AlQj4fc"},
		},
	},
	122: {
		ID:   122,
		Time: "2025-03-20 12:30:15",
		Text: `早上刚出发的时候，
她接到了一个电话，她接到了一个电话她接到了一个电话她接到了一个电话她接到了一个电话她接到了一个电话她接到了一个电话她接到了一个电话她接到了一个电话她接到了一个电话她接到了一个电话她接到了一个电话她接到了一个电话她接到了一个电话
“你出发没？”，
“刚出发，怎么了？”，
“哦，没事，我也马上出发。”，
“好，我8点46到。”，
“哦，我8点52到（小朋友的妈妈告诉他的）。”，
两人像老朋友马上要见面一样的聊着……，

小朋友好像是叫涵涵，
我以为是个女孩子，
然后……，
是个男孩子，

到了之后才发现，
今天穿的都是purple，
可明明这件紫色外套是我去顺丰快递点刚拿的，
意外，
又那么巧合，

到了之后才发现，
今天穿的都是purple，
可明明这件紫色外套是我去顺丰快递点刚拿的，
意外，
又那么巧合，

到了之后才发现，
今天穿的都是purple，
可明明这件紫色外套是我去顺丰快递点刚拿的，
意外，
又那么巧合，

到了之后才发现，
今天穿的都是purple，
可明明这件紫色外套是我去顺丰快递点刚拿的，
意外，
又那么巧合，

到了之后才发现，
今天穿的都是purple，
可明明这件紫色外套是我去顺丰快递点刚拿的，
意外，
又那么巧合，

到了之后才发现，
今天穿的都是purple，
可明明这件紫色外套是我去顺丰快递点刚拿的，
意外，
又那么巧合，

到了之后才发现，
今天穿的都是purple，
可明明这件紫色外套是我去顺丰快递点刚拿的，
意外，
又那么巧合，

第三十七次记录，
关于你的#CAMPNOW`,
		Pictures: []Picture{
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"},
			{Width: 1620, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkUp1vXcxc4J0kVhFwxPRnkLNqvFqER99R2FC0BDmCYx8"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnorDmccicxxRQhREKbJjQGuQocfNrSvyvrytnoGHcwSWmDHFQTMbdDj4"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9MyDLXiaZ7YtL7DgIDWPqNMS8odq91EdX586jQx2UDvlo"},
			{Width: 1080, Height: 1440, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9rhKE57uzPHBIt4ldv1btOMa0ibW5zxlKRYXQaQMico61Q"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9vlFfMRibmJ3m0nJR94FaibxBFEiaw3Dq3UM3fs7cD1ReqA"},
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFU0ic58W0HarI0kZFBdia9cXEwibbzG2RqEYr0bmiaYMwJ7E6OwMp6haQkk"},
			{Width: 1080, Height: 1440, URL: "https://img.diandibianji.com/8u9KefYVGSVTCWiabXdfPI8xRknynoVGZ1rC9M13tRXSU3A1libL6xT8eTkbrtRcRtXOR2C33FSU8"},
		},
	},
	123: {
		ID:   123,
		Time: "2025-03-20 12:30:15",
		Text: "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n",
		Pictures: []Picture{
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"},
			{Width: 1620, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkUp1vXcxc4J0kVhFwxPRnkLNqvFqER99R2FC0BDmCYx8"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnorDmccicxxRQhREKbJjQGuQocfNrSvyvrytnoGHcwSWmDHFQTMbdDj4"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9MyDLXiaZ7YtL7DgIDWPqNMS8odq91EdX586jQx2UDvlo"},
			{Width: 1080, Height: 1440, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9rhKE57uzPHBIt4ldv1btOMa0ibW5zxlKRYXQaQMico61Q"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9vlFfMRibmJ3m0nJR94FaibxBFEiaw3Dq3UM3fs7cD1ReqA"},
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFU0ic58W0HarI0kZFBdia9cXEwibbzG2RqEYr0bmiaYMwJ7E6OwMp6haQkk"},
		},
	},
	124: {
		ID:   124,
		Time: "2025-03-20 12:30:15",
		Text: "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n",
		Pictures: []Picture{
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"},
			{Width: 1620, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkUp1vXcxc4J0kVhFwxPRnkLNqvFqER99R2FC0BDmCYx8"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnorDmccicxxRQhREKbJjQGuQocfNrSvyvrytnoGHcwSWmDHFQTMbdDj4"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9MyDLXiaZ7YtL7DgIDWPqNMS8odq91EdX586jQx2UDvlo"},
			{Width: 1080, Height: 1440, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9rhKE57uzPHBIt4ldv1btOMa0ibW5zxlKRYXQaQMico61Q"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9vlFfMRibmJ3m0nJR94FaibxBFEiaw3Dq3UM3fs7cD1ReqA"},
		},
	},
	125: {
		ID:   125,
		Time: "2025-03-20 12:30:15",
		Text: "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n",
		Pictures: []Picture{
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"},
			{Width: 1620, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkUp1vXcxc4J0kVhFwxPRnkLNqvFqER99R2FC0BDmCYx8"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnorDmccicxxRQhREKbJjQGuQocfNrSvyvrytnoGHcwSWmDHFQTMbdDj4"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9MyDLXiaZ7YtL7DgIDWPqNMS8odq91EdX586jQx2UDvlo"},
			{Width: 1080, Height: 1440, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9rhKE57uzPHBIt4ldv1btOMa0ibW5zxlKRYXQaQMico61Q"},
		},
	},
	126: {
		ID:   126,
		Time: "2025-03-20 12:30:15",
		Text: "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n",
		Pictures: []Picture{
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"},
			{Width: 1620, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkUp1vXcxc4J0kVhFwxPRnkLNqvFqER99R2FC0BDmCYx8"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnorDmccicxxRQhREKbJjQGuQocfNrSvyvrytnoGHcwSWmDHFQTMbdDj4"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9MyDLXiaZ7YtL7DgIDWPqNMS8odq91EdX586jQx2UDvlo"},
		},
	},
	127: {
		ID:   127,
		Time: "2025-03-20 12:30:15",
		Text: "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n",
		Pictures: []Picture{
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"},
			{Width: 1620, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkUp1vXcxc4J0kVhFwxPRnkLNqvFqER99R2FC0BDmCYx8"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnorDmccicxxRQhREKbJjQGuQocfNrSvyvrytnoGHcwSWmDHFQTMbdDj4"},
		},
	},
	128: {
		ID:   128,
		Time: "2025-03-20 12:30:15",
		Text: "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n",
		Pictures: []Picture{
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"},
			{Width: 1620, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkUp1vXcxc4J0kVhFwxPRnkLNqvFqER99R2FC0BDmCYx8"},
		},
	},
	129: {
		ID:   129,
		Time: "2025-03-20 12:30:15",
		Text: "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n",
		Pictures: []Picture{
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"},
		},
	},
	130: {
		ID:       130,
		Time:     "2025-03-20 12:30:15",
		Text:     "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n" + "这是一个需要跨多页的长文本...\n",
		Pictures: []Picture{},
	},
	131: {
		ID:   131,
		Time: "2025-03-20 12:30:15",
		Text: "",
		Pictures: []Picture{
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkFyklENiaLE4qZQMa4f8rViacHtmAO6RicznsqlO4iaa0LJI"},
			{Width: 1620, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFatvViceRfiaBkUp1vXcxc4J0kVhFwxPRnkLNqvFqER99R2FC0BDmCYx8"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnorDmccicxxRQhREKbJjQGuQocfNrSvyvrytnoGHcwSWmDHFQTMbdDj4"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9MyDLXiaZ7YtL7DgIDWPqNMS8odq91EdX586jQx2UDvlo"},
			{Width: 1080, Height: 1440, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9rhKE57uzPHBIt4ldv1btOMa0ibW5zxlKRYXQaQMico61Q"},
			{Width: 810, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSUU5adtZJnPnricVO6bUEiav9vlFfMRibmJ3m0nJR94FaibxBFEiaw3Dq3UM3fs7cD1ReqA"},
			{Width: 1080, Height: 1620, URL: "https://img.diandibianji.com/8u9KefYVGSUWsdSd01LaFU0ic58W0HarI0kZFBdia9cXEwibbzG2RqEYr0bmiaYMwJ7E6OwMp6haQkk"},
			{Width: 1080, Height: 1440, URL: "https://img.diandibianji.com/8u9KefYVGSVTCWiabXdfPI8xRknynoVGZ1rC9M13tRXSU3A1libL6xT8eTkbrtRcRtXOR2C33FSU8"},
			{Width: 1440, Height: 1080, URL: "https://img.diandibianji.com/8u9KefYVGSVTCWiabXdfPIib3W7TLFloYBj0w7WOWtxxawey8bHgg4Tyqzrkwre1V8dNA7AlQj4fc"},
		},
	},
}
