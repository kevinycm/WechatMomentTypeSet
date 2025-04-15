package backend

import "fmt"

// calculateSevenPicturesLayout defines templates and calculates dimensions for 7 pictures.
func (e *ContinuousLayoutEngine) calculateSevenPicturesLayout(pictures []Picture) (TemplateLayout, error) {
	if len(pictures) != 7 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 7-pic layout")
	}

	// Define minimums (Rule 3.8)
	minAllowedHeight := 2500.00
	minAllowedPictureWidth := 1666.67
	spacing := e.imageSpacing
	AW := e.availableWidth

	// Get Aspect Ratios (W/H)
	ARs := make([]float64, 7)
	avgARtop3 := 0.0
	avgARtop4 := 0.0
	validARcount := 0
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
			if i < 3 {
				avgARtop3 += ARs[i]
			}
			if i < 4 {
				avgARtop4 += ARs[i]
			}
			validARcount++
		} else {
			ARs[i] = 1.0 // Default AR
			if i < 3 {
				avgARtop3 += 1.0
			}
			if i < 4 {
				avgARtop4 += 1.0
			}
			validARcount++
		}
	}
	if validARcount < 7 {
		return TemplateLayout{}, fmt.Errorf("invalid picture dimensions found for 7-pic layout")
	}
	avgARtop3 /= 3.0
	avgARtop4 /= 4.0

	// --- Template Selection Heuristic ---
	// If the first 3 are wider on average than the first 4, try 3T4B first.
	preferredOrder := []string{"4T3B", "3T4B"} // Default order (4 Top, 3 Bottom)
	if avgARtop3 > avgARtop4 {                 // Heuristic check
		preferredOrder = []string{"3T4B", "4T3B"} // Prefer 3 Top, 4 Bottom
	}

	var bestFallbackLayout TemplateLayout
	var fallbackFound bool = false
	var firstError error = nil

	// --- Try Templates in Preferred Order ---
	for _, templateName := range preferredOrder {
		var layout TemplateLayout
		var meetsMin bool
		var err error

		switch templateName {
		case "3T4B":
			layout, meetsMin, err = e.calculateLayout_3T4B(pictures, ARs, AW, spacing, minAllowedHeight, minAllowedPictureWidth)
		case "4T3B":
			layout, meetsMin, err = e.calculateLayout_4T3B(pictures, ARs, AW, spacing, minAllowedHeight, minAllowedPictureWidth)
		default:
			continue // Should not happen
		}

		if err == nil {
			if meetsMin {
				fmt.Printf("Info: Using 7-pic template: %s\n", templateName)
				return layout, nil // Found a valid layout that meets minimums
			} else if !fallbackFound {
				// Store the first calculated layout that doesn't meet minimums as a fallback
				bestFallbackLayout = layout
				fallbackFound = true
				fmt.Printf("Info: Storing 7-pic template %s as fallback (violates minimums).\n", templateName)
			}
		} else if firstError == nil {
			firstError = fmt.Errorf("template %s failed: %w", templateName, err)
		}
	}

	// --- Handle Results ---
	if fallbackFound {
		fmt.Println("Warning: No 7-pic template met minimum dimensions, using best fallback.")
		return bestFallbackLayout, nil // Return the best fallback we found
	} else if firstError != nil {
		// All attempts failed with errors
		return TemplateLayout{}, fmt.Errorf("all 7-pic templates failed: %w", firstError)
	} else {
		// Should not be reached if at least one template calculation was attempted
		return TemplateLayout{}, fmt.Errorf("unexpected error in 7-pic layout selection")
	}
}
