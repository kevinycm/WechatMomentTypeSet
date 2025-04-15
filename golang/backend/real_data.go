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
		AND (LENGTH(qiniu_media_urls) - LENGTH(REPLACE(qiniu_media_urls, ',', '')) + 1) / 2 = 6
		ORDER BY release_time DESC
	`)

	// rows, err := db.Query(`
	// 	SELECT
	// 		id,
	// 		release_time,
	// 		text,
	// 		media_infos,
	// 		qiniu_media_urls
	// 	FROM new_moment
	// 	WHERE deleted_at IS NULL
	// 	AND type = 1
	// 	AND media_infos IS NOT NULL
	// 	AND qiniu_media_urls IS NOT NULL
	// 	AND (LENGTH(qiniu_media_urls) - LENGTH(REPLACE(qiniu_media_urls, ',', '')) + 1) / 2 <= 2
	// 	ORDER BY release_time DESC
	// `)
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
	log.Printf("Processing picture info: mediaInfos='%s', qiniuMediaURLs='%s'", mediaInfos, qiniuMediaURLs)
	if mediaInfos == "" || qiniuMediaURLs == "" {
		log.Println(" MediaInfos or QiniuMediaURLs is empty, returning nil.")
		return nil, nil
	}

	// Split media infos (width,height pairs)
	dimensions := strings.Split(mediaInfos, ",")
	log.Printf(" Split dimensions: %v (count: %d)", dimensions, len(dimensions))
	if len(dimensions)%2 != 0 {
		log.Printf(" Error: Odd number of dimensions (%d).", len(dimensions))
		// Return nil or an empty slice depending on desired behavior for malformed data
		return nil, fmt.Errorf("odd number of dimensions found in mediaInfos: %s", mediaInfos)
	}

	// Split URLs (expecting id,url,id,url...)
	urls := strings.Split(qiniuMediaURLs, ",")
	log.Printf(" Split URLs: %v (count: %d)", urls, len(urls))
	if len(urls)%2 != 0 {
		log.Printf(" Warning: Odd number of URL components (%d). Check format.", len(urls))
		// Continue processing, but might result in fewer pictures than dimensions
	}

	var pictures []Picture
	numPics := len(dimensions) / 2
	log.Printf(" Expected number of pictures based on dimensions: %d", numPics)

	for i := 0; i < numPics; i++ {
		dimIndex := i * 2
		urlValueIndex := i*2 + 1 // URL is the second element in each pair (index 1, 3, 5...)

		// Check if indices are within bounds
		if dimIndex+1 >= len(dimensions) {
			log.Printf(" Error: Dimension index %d out of bounds (len=%d)", dimIndex+1, len(dimensions))
			break // Stop processing if dimensions are insufficient
		}
		if urlValueIndex >= len(urls) {
			log.Printf(" Error: URL index %d out of bounds (len=%d)", urlValueIndex, len(urls))
			break // Stop processing if URLs are insufficient
		}

		widthStr := dimensions[dimIndex]
		heightStr := dimensions[dimIndex+1]
		url := urls[urlValueIndex]
		log.Printf("  Processing pic %d: widthStr='%s', heightStr='%s', url='%s'", i, widthStr, heightStr, url)

		// Parse width and height
		var width, height int
		_, err := fmt.Sscanf(widthStr+","+heightStr, "%d,%d", &width, &height)
		if err != nil {
			log.Printf("Error parsing dimensions for pic %d ('%s', '%s'): %v. Skipping picture.", i, widthStr, heightStr, err)
			continue
		}

		// Basic validation
		if width <= 0 || height <= 0 {
			log.Printf("Warning: Invalid parsed dimensions for pic %d (width=%d, height=%d). Skipping picture.", i, width, height)
			continue
		}
		if url == "" {
			log.Printf("Warning: Empty URL for pic %d. Skipping picture.", i)
			continue
		}

		// Create Picture entry
		picture := Picture{
			Index:  i,
			Width:  width,
			Height: height,
			URL:    url,
		}
		pictures = append(pictures, picture)
		log.Printf("  Successfully processed pic %d.", i)
	}

	log.Printf(" Finished processing. Created %d pictures.", len(pictures))
	return pictures, nil
}

// GetRealDataByID returns a specific test case from RealData
func GetRealDataByID(id int) (TestCase, bool) {
	testCase, exists := RealData[id]
	return testCase, exists
}
