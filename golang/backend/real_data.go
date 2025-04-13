package backend

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// NewMoment represents the database table structure
type NewMoment struct {
	ID             int64     `json:"id"`
	ReleaseTime    time.Time `json:"release_time"`
	Text           string    `json:"text"`
	MediaInfos     string    `json:"media_infos"`
	QiniuMediaURLs string    `json:"qiniu_media_urls"`
}

// RealData stores the data fetched from MySQL
var RealData map[int]TestCase

// InitRealData initializes the RealData map by fetching data from MySQL
func InitRealData(dsn string) error {
	// Connect to MySQL
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		return err
	}

	// Query the database
	rows, err := db.Query(`
		SELECT 
			id,
			release_time,
			text,
			media_infos,
			qiniu_media_urls
		FROM new_moment 
		WHERE deleted_at IS NULL 
		AND type = 1
		AND media_infos IS NOT NULL 
		AND qiniu_media_urls IS NOT NULL
		AND (LENGTH(qiniu_media_urls) - LENGTH(REPLACE(qiniu_media_urls, ',', '')) + 1) / 2 <= 2
		ORDER BY release_time DESC
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	RealData = make(map[int]TestCase)

	// Process each row
	for rows.Next() {
		var moment NewMoment
		err := rows.Scan(
			&moment.ID,
			&moment.ReleaseTime,
			&moment.Text,
			&moment.MediaInfos,
			&moment.QiniuMediaURLs,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Process media information
		pictures, err := processPictureInfo(moment.MediaInfos, moment.QiniuMediaURLs)
		if err != nil {
			log.Printf("Error processing picture info for ID %d: %v", moment.ID, err)
			continue
		}

		// Create TestCase entry
		testCase := TestCase{
			ID:       int(moment.ID),
			Time:     moment.ReleaseTime.Format("2006-01-02 15:04:05"),
			Text:     moment.Text,
			Pictures: pictures,
		}

		RealData[int(moment.ID)] = testCase
	}

	return rows.Err()
}

// processPictureInfo converts media_infos and qiniu_media_urls into Picture slice
func processPictureInfo(mediaInfos, qiniuMediaURLs string) ([]Picture, error) {
	if mediaInfos == "" || qiniuMediaURLs == "" {
		return nil, nil
	}

	// Split media infos (width,height pairs)
	dimensions := strings.Split(mediaInfos, ",")
	if len(dimensions)%2 != 0 {
		return nil, nil
	}

	// Split URLs
	urls := strings.Split(qiniuMediaURLs, ",")

	var pictures []Picture
	urlIndex := 1 // Start from index 1 (second URL)

	// Process each width,height pair
	for i := 0; i < len(dimensions); i += 2 {
		// Skip if we've reached the end of URLs
		if urlIndex >= len(urls) {
			break
		}

		// Parse width and height
		var width, height int
		_, err := fmt.Sscanf(dimensions[i]+","+dimensions[i+1], "%d,%d", &width, &height)
		if err != nil {
			log.Printf("Error parsing dimensions: %v", err)
			continue
		}

		// Create Picture entry
		picture := Picture{
			Width:  width,
			Height: height,
			URL:    urls[urlIndex],
		}
		pictures = append(pictures, picture)
		urlIndex += 2 // Move to next even-indexed URL
	}

	return pictures, nil
}

// GetRealDataByID returns a specific test case from RealData
func GetRealDataByID(id int) (TestCase, bool) {
	testCase, exists := RealData[id]
	return testCase, exists
}
