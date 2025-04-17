package calculate

import (
	"fmt"
)

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
	// Clamp numPics to the valid range [1, 9] for slice access
	idx := numPics
	if idx < 1 {
		idx = 1
	}
	if idx > 9 {
		idx = 9
	}

	switch picType {
	case "wide":
		return e.minWideHeight // Wide min height is constant
	case "tall":
		return e.minTallHeight // Tall min height is constant
	case "landscape":
		// Make sure the slice has been initialized and the index is valid
		if len(e.minLandscapeHeights) > idx-1 {
			return e.minLandscapeHeights[idx-1]
		}
		fmt.Printf("Warning: minLandscapeHeights not properly initialized or index out of bounds (%d)\n", idx)
		return 800.0 // Return default landscape height as fallback
	case "portrait":
		// Make sure the slice has been initialized and the index is valid
		if len(e.minPortraitHeights) > idx-1 {
			return e.minPortraitHeights[idx-1]
		}
		fmt.Printf("Warning: minPortraitHeights not properly initialized or index out of bounds (%d)\n", idx)
		return 1000.0 // Return default portrait height as fallback
	default: // square, unknown
		// Use landscape height as fallback
		if len(e.minLandscapeHeights) > idx-1 {
			return e.minLandscapeHeights[idx-1]
		}
		fmt.Printf("Warning: Fallback minLandscapeHeights not properly initialized or index out of bounds (%d)\n", idx)
		return 800.0 // Return default landscape height as fallback
	}
}
