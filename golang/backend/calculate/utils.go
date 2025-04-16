package calculate

import (
	"fmt"
)

// Layout represents a picture layout configuration
type Layout struct {
}

// Helper function to get picture type based on Aspect Ratio (AR)
func GetPictureType(aspectRatio float64) string {
	if aspectRatio >= 3.0 {
		return "wide"
	} else if aspectRatio <= 1.0/3.0 {
		return "tall"
	} else if aspectRatio > 1.0 && aspectRatio < 3.0 {
		return "landscape"
	} else if aspectRatio > 1.0/3.0 && aspectRatio < 1.0 {
		return "portrait"
	} else if aspectRatio == 1.0 {
		return "square" // Added for completeness, rules don't specify minimums for square
	} else {
		return "unknown" // Should not happen with valid positive dimensions
	}
}

// Helper function to get the required minimum height for a given picture type and group size.
// UPDATED SIGNATURE AND LOGIC
func GetRequiredMinHeight(e *ContinuousLayoutEngine, picType string, numPics int) float64 {
	switch picType {
	case "wide":
		return e.minWideHeight // Wide min height is constant
	case "tall":
		return e.minTallHeight // Tall min height is constant
	case "landscape":
		if numPics >= 8 {
			return e.minLandscapeHeightVeryLargeGroup // Use specific value for >= 8 pics
		} else if numPics >= 5 {
			return e.minLandscapeHeightLargeGroup // Use 600 for 5-7 pics
		} else {
			return e.minLandscapeHeight // Use base 400 for < 5 pics
		}
	case "portrait":
		if numPics >= 8 {
			return e.minPortraitHeightVeryLargeGroup // Use specific value for >= 8 pics
		} else if numPics >= 5 {
			return e.minPortraitHeightLargeGroup // Use 900 for 5-7 pics
		} else {
			return e.minPortraitHeight // Use base 600 for < 5 pics
		}
	default: // square, unknown
		// Use landscape height as fallback, respecting numPics tiers
		if numPics >= 8 {
			return e.minLandscapeHeightVeryLargeGroup
		} else if numPics >= 5 {
			return e.minLandscapeHeightLargeGroup
		} else {
			return e.minLandscapeHeight
		}
	}
}

// checkMinHeights verifies if all pictures in a calculated layout meet their type-specific minimum height requirements.
// UPDATED TO CALL getRequiredMinHeight WITH numPics
func CheckMinHeights(e *ContinuousLayoutEngine, layout TemplateLayout, types []string, numPics int) bool {
	if len(layout.Dimensions) != numPics || len(types) != numPics {
		fmt.Printf("Warning: checkMinHeights received mismatched lengths (Dimensions: %d, Types: %d, expected: %d)\n", len(layout.Dimensions), len(types), numPics)
		return false // Data inconsistency
	}

	for i := 0; i < numPics; i++ {
		picType := types[i]
		// Pass numPics to get the correct minimum height
		requiredMinHeight := GetRequiredMinHeight(e, picType, numPics)

		// Safely access dimensions
		if i >= len(layout.Dimensions) || len(layout.Dimensions[i]) != 2 {
			fmt.Printf("Warning: Invalid dimensions data for picture %d in checkMinHeights.\n", i)
			return false // Invalid layout data
		}
		actualHeight := layout.Dimensions[i][1]

		if actualHeight < requiredMinHeight-1e-6 { // Use tolerance for float comparison
			fmt.Printf("Debug: Min height check failed for pic %d (Type: %s, NumPics: %d). Required: %.2f, Actual: %.2f\n", i, picType, numPics, requiredMinHeight, actualHeight)
			return false // Minimum height not met
		}
	}
	return true // All checks passed
}
