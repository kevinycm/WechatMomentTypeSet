package backend

import "fmt"

// calculateFivePicturesLayout defines templates and calculates dimensions for 5 pictures using a geometric approach.
func (e *ContinuousLayoutEngine) calculateFivePicturesLayout(pictures []Picture) (TemplateLayout, error) {
	if len(pictures) != 5 {
		return TemplateLayout{}, fmt.Errorf("incorrect number of pictures for 5-pic layout")
	}

	// Define minimums (Rule 3.8)
	minAllowedHeight := 2500.00
	minAllowedPictureWidth := 1666.67
	spacing := e.imageSpacing
	AW := e.availableWidth

	// Get Aspect Ratios (W/H)
	ARs := make([]float64, 5)
	avgARtop2 := 1.0
	avgARtop3 := 1.0
	validARcount := 0
	for i, pic := range pictures {
		if pic.Height > 0 && pic.Width > 0 {
			ARs[i] = float64(pic.Width) / float64(pic.Height)
			if i < 2 {
				avgARtop2 += ARs[i]
			}
			if i < 3 {
				avgARtop3 += ARs[i]
			}
			validARcount++
		} else {
			ARs[i] = 1.0 // Default AR
			if i < 2 {
				avgARtop2 += 1.0
			}
			if i < 3 {
				avgARtop3 += 1.0
			}
			validARcount++
		}
	}
	if validARcount < 5 {
		return TemplateLayout{}, fmt.Errorf("invalid picture dimensions found")
	}
	avgARtop2 /= 2.0
	avgARtop3 /= 3.0

	// --- Template Selection Heuristic ---
	// If the first 2 are wider on average, try 2T3B first.
	// Otherwise, try 3T2B first.
	preferredOrder := []string{"3T2B", "2T3B"} // Default order
	if avgARtop2 > avgARtop3 {                 // Heuristic check
		preferredOrder = []string{"2T3B", "3T2B"}
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
		case "2T3B":
			layout, meetsMin, err = e.calculateLayout_2T3B(pictures, ARs, AW, spacing, minAllowedHeight, minAllowedPictureWidth)
		case "3T2B":
			layout, meetsMin, err = e.calculateLayout_3T2B(pictures, ARs, AW, spacing, minAllowedHeight, minAllowedPictureWidth)
		default:
			continue // Should not happen
		}

		if err == nil {
			if meetsMin {
				fmt.Printf("Info: Using 5-pic template: %s\n", templateName)
				return layout, nil // Found a valid layout that meets minimums
			} else if !fallbackFound {
				// Store the first calculated layout that doesn't meet minimums as a fallback
				bestFallbackLayout = layout
				fallbackFound = true
				fmt.Printf("Info: Storing 5-pic template %s as fallback (violates minimums).\n", templateName)
			}
		} else if firstError == nil {
			firstError = fmt.Errorf("template %s failed: %w", templateName, err)
		}
	}

	// --- Handle Results ---
	if fallbackFound {
		fmt.Println("Warning: No 5-pic template met minimum dimensions, using best fallback.")
		return bestFallbackLayout, nil // Return the best fallback we found
	} else if firstError != nil {
		// All attempts failed with errors
		return TemplateLayout{}, fmt.Errorf("all 5-pic templates failed: %w", firstError)
	} else {
		// Should not be reached if at least one template calculation was attempted
		return TemplateLayout{}, fmt.Errorf("unexpected error in 5-pic layout selection")
	}
}
