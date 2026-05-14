package tracks

import (
	"regexp"
	"strings"
)

var punctuationRE = regexp.MustCompile(`[^\w\s]`)

func normalize(s string) string {
	s = strings.ToLower(s)
	s = punctuationRE.ReplaceAllString(s, " ")
	return strings.Join(strings.Fields(s), " ")
}

func MatchTrack(title, description string) (string, string, bool) {
	normTitle := normalize(title)
	normDesc := normalize(description)
	titleWords := strings.Fields(normTitle)
	if len(titleWords) == 0 {
		return "", "", false
	}

	bestAreaID := ""
	bestTrackName := ""
	maxScore := 0

	for _, area := range AllAreas {
		for _, track := range area.Tracks {
			score := calculateDetailedScore(normTitle, titleWords, normDesc, area, track)
			if score > 0 {
				if score > maxScore {
					maxScore = score
					bestAreaID = area.ID
					bestTrackName = track.Name
				} else if score == maxScore {
					if len(track.Name) > len(bestTrackName) {
						bestAreaID = area.ID
						bestTrackName = track.Name
					}
				}
			}
		}
	}

	// Threshold raised: Matches must be meaningful.
	if maxScore >= 25 {
		return bestAreaID, bestTrackName, true
	}

	return "", "", false
}

var wordWeights = map[string]int{
	"qa":          40, // Bumped QA weight to ensure it outscores generic "data"
	"tqa":         40,
	"cloud":       30,
	"aws":         30,
	"azure":       30,
	"devops":      30,
	"android":     30,
	"php":         30,
	"javascript":  30,
	"react":       30,
	"java":        30,
	"python":      30,
	"ux":          30,
	"ui":          30,
	"seo":         30,
	"scrum":       30,
	"agile":       30,
	"security":    30,
	"sql":         30,
	"oracle":      30,
	"plsql":       30,
	"robotics":    30,
	"multimedia":  30,
	"ecommerce":   30,
	"data":        20, // Still strong but less than QA/Cloud/etc.
	"science":     20,
	"product":     25,
	"project":     25,
	"marketing":   20,
	"growth":      20,
	"automation":  20,
	"fullstack":   15,
	"frontend":    15,
	"backend":     15,
	"sre":         30,
	"talent":      20,
	"recruiter":   20,
	"hr":          20,
	"designer":    15,
	"developer":   5,
	"engineer":    5,
	"analyst":     5,
	"manager":     5,
	"support":     5,
	"technician":  5,
	"specialist":  5,
	"admin":       10,
	"administrator": 10,
}

// areaWords maps Area IDs to specific keywords that help identify them.
var areaKeywords = map[string][]string{
	"TQA":     {"qa", "testing", "test", "tester"},
	"CLOUD":   {"aws", "azure", "google cloud", "gcp", "cloud"},
	"IA":      {"ai", "ia", "artificial intelligence", "ml", "learning"},
	"SINF":    {"security", "cybersecurity", "ciberseguridad", "sinf", "soc", "penetration"}, // SINF = Seguridad Info
	"DATA":    {"data", "scientist", "scientist", "datum", "analitics"},
	"REDES":   {"network", "redes", "cisico", "switch", "router"},
}

func calculateDetailedScore(normTitle string, titleWords []string, normDesc string, area Area, track Track) int {
	normTrack := normalize(track.Name)
	trackWords := strings.Fields(normTrack)
	normAreaName := normalize(area.Name)
	normAreaID := normalize(area.ID)

	score := 0

	// 1. Exact track name in title (Excellent match)
	if strings.Contains(normTitle, normTrack) {
		score += 60
	}

	// 2. Exact match on high-signal words in Title
	foundSignalInTitle := false
	for _, tw := range trackWords {
		weight := wordWeights[tw]
		if weight == 0 {
			weight = 5
		}
		
		matchFound := false
		for _, jw := range titleWords {
			if tw == jw {
				score += weight
				matchFound = true
				if weight >= 25 {
					foundSignalInTitle = true
				}
				break
			}
		}

		if !matchFound && strings.Contains(normDesc, tw) {
			score += weight / 2
			if weight >= 25 {
				foundSignalInTitle = true
			}
		}
	}

	// 3. Area Correlation (Keyword-based)
	// Check if the title matches keywords specific to the Area.
	if keywords, ok := areaKeywords[area.ID]; ok {
		for _, kw := range keywords {
			if strings.Contains(normTitle, kw) {
				score += 25
				foundSignalInTitle = true
			}
		}
	}
	
	// Also check Area ID/Name directly (but only as whole word for short IDs)
	if len(normAreaID) > 2 {
		if strings.Contains(" "+normTitle+" ", " "+normAreaID+" ") {
			score += 15
		}
	}
	if strings.Contains(normTitle, strings.TrimPrefix(normAreaName, "area ")) {
		score += 15
	}

	// 4. Verification Check: High Signal words MUST BE PRESENT
	// If a track name has a VERY strong word (e.g., "QA", "Cloud", "PHP"),
	// it MUST appear in the title or description.
	for _, tw := range trackWords {
		if wordWeights[tw] >= 30 {
			if !strings.Contains(normTitle, tw) && !strings.Contains(normDesc, tw) {
				return 0 // Wrong context
			}
		}
	}

	// 5. Penalize mismatched strong keywords
	// If the title has "QA" but the track is "Data Scientist" (which has NO QA keywords),
	// we should probably avoid the match if another track HAS "QA".
	// (This is implicitly handled by higher scores for the "right" track).

	// Threshold: If no significant signal word from the track was found in the title/desc,
	// and the track is not an exact match, discard.
	if score > 0 && !foundSignalInTitle && score < 25 {
		return 0
	}

	return score
}
