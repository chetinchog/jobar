package providers

import (
	"anuncios/internal/models"
	"anuncios/internal/parser"
	"anuncios/internal/utils"
	"fmt"
	"strings"
)

type Jobicy struct {
	URLs    []string
	MaxJobs int
}

func (p *Jobicy) GetName() string { return "Jobicy" }
func (p *Jobicy) GetMaxJobs() int { return p.MaxJobs }

func (p *Jobicy) FetchJobs(isDuplicate func(models.Item) bool, process func(models.Item)) error {
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

			// Use content:encoded for full description; fall back to description
			desc := parser.ExtractTagCustom(block, "<content:encoded>", "</content:encoded>")
			if desc == "" {
				desc = parser.ExtractTag(block, "description")
			}
			// Strip CDATA wrappers
			desc = strings.TrimPrefix(desc, "<![CDATA[")
			desc = strings.TrimSuffix(desc, "]]>")

			item := models.Item{
				Title:       parser.ExtractTag(block, "title"),
				Description: strings.TrimSpace(desc),
				PubDate:     parser.ExtractTag(block, "pubDate"),
				Guid:        parser.ExtractTag(block, "guid"),
				Link:        parser.ExtractTagCustom(block, "<link>", "</link>"),
				Company:     parser.ExtractTagCustom(block, "<job_listing:company>", "</job_listing:company>"),
				Location:    parser.ExtractTagCustom(block, "<job_listing:location>", "</job_listing:location>"),
				JobIdTag:    parser.ExtractTag(block, "jobid"),
			}

			// Strip CDATA from company and location
			item.Company = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(item.Company, "<![CDATA["), "]]>"))
			item.Location = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(item.Location, "<![CDATA["), "]]>"))

			if item.Title == "" && item.Description == "" {
				continue
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
