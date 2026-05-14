package providers

import (
	"anuncios/internal/models"
	"anuncios/internal/parser"
	"anuncios/internal/utils"
	"fmt"
	"regexp"
	"strings"
)

var findJobITRe = regexp.MustCompile(`(?i)(.*)\s+(?:busca\s+un|busca\s+una|se\s+busca)\s+(.*)`)

type FindJobIT struct {
	URLs    []string
	MaxJobs int
}

func (p *FindJobIT) GetName() string { return "FindJobIT" }
func (p *FindJobIT) GetMaxJobs() int { return p.MaxJobs }

func (p *FindJobIT) FetchJobs(isDuplicate func(models.Item) bool, process func(models.Item)) error {
	count := 0
	for _, url := range p.URLs {
		if p.MaxJobs > 0 && count >= p.MaxJobs {
			break
		}
		fmt.Printf("Fetching: %s\n", url)
		dataStr, err := utils.FetchURL(url)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", url, err)
			continue
		}

		itemBlocks := strings.Split(dataStr, "<item>")
		if len(itemBlocks) < 2 {
			continue
		}
		itemBlocks = itemBlocks[1:]

		for _, block := range itemBlocks {
			if endIdx := strings.Index(block, "</item>"); endIdx != -1 {
				block = block[:endIdx]
			}

			item := models.Item{
				Title:       parser.ExtractTag(block, "title"),
				Region:      parser.ExtractTag(block, "region"),
				Country:     parser.ExtractTag(block, "country"),
				Estado:      parser.ExtractTag(block, "state"),
				Description: parser.ExtractTag(block, "description"),
				PubDate:     parser.ExtractTag(block, "pubDate"),
				ExpiresAt:   parser.ExtractTag(block, "expires_at"),
				Guid:        parser.ExtractTag(block, "guid"),
				Link:        parser.ExtractTag(block, "link"),
				Location:    parser.ExtractTag(block, "location"),
				Company:     parser.ExtractTag(block, "company"),
				JobIdTag:    parser.ExtractTag(block, "jobid"),
			}

			if item.Title == "" && item.Description == "" {
				continue
			}

			// FindJobIT specific title splitting
			if matches := findJobITRe.FindStringSubmatch(item.Title); len(matches) > 2 {
				item.Company = strings.TrimSpace(matches[1])
				item.Title = strings.TrimSpace(matches[2])
			}

			if !isDuplicate(item) {
				process(item)
				count++
				if p.MaxJobs > 0 && count >= p.MaxJobs {
					break
				}
			}
		}
	}
	return nil
}
