package storage

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
)

func LoadAllExistingData(baseDir string) (map[string]bool, map[string]bool) {
	ids := make(map[string]bool)
	titles := make(map[string]bool)

	// Check output folder
	loadFromFolder(filepath.Join(baseDir, "output"), ids, titles)
	// Check discard folder
	loadFromFolder(filepath.Join(baseDir, "discard"), ids, titles)

	return ids, titles
}

func loadFromFolder(dir string, ids map[string]bool, titles map[string]bool) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
			loadFromFile(filepath.Join(dir, file.Name()), ids, titles)
		}
	}
}

func loadFromFile(filename string, ids map[string]bool, titles map[string]bool) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return
	}

	if len(records) < 2 {
		return
	}

	header := records[0]
	idxID := -1
	idxCompany := -1
	idxTitle := -1

	for i, h := range header {
		switch strings.ToLower(h) {
		case "jobid", "job_id":
			idxID = i
		case "company":
			idxCompany = i
		case "title":
			idxTitle = i
		}
	}

	for i, row := range records {
		if i == 0 {
			continue
		}
		if idxID != -1 && idxID < len(row) {
			ids[row[idxID]] = true
		}
		if idxCompany != -1 && idxTitle != -1 && idxCompany < len(row) && idxTitle < len(row) {
			titles[row[idxCompany]+"|"+row[idxTitle]] = true
		}
	}
}
