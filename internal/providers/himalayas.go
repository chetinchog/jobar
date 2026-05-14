package providers

import (
	"anuncios/internal/models"
	"anuncios/internal/utils"
	"encoding/json"
	"fmt"
	"time"
)

type Himalayas struct {
	URLs    []string
	MaxJobs int
}

func (p *Himalayas) GetName() string { return "Himalayas" }
func (p *Himalayas) GetMaxJobs() int { return p.MaxJobs }

func (p *Himalayas) FetchJobs(isDuplicate func(models.Item) bool, process func(models.Item)) error {
	count := 0
	offset, limit := 0, 20
	baseUrl := p.URLs[0]

	for {
		if p.MaxJobs > 0 && count >= p.MaxJobs {
			break
		}
		url := fmt.Sprintf("%s?limit=%d&offset=%d", baseUrl, limit, offset)
		fmt.Printf("Fetching: %s\n", url)
		dataStr, err := utils.FetchURL(url)
		if err != nil {
			if err.Error() == "RATE_LIMIT" {
				fmt.Println("Rate limit exceeded. Waiting 60 seconds...")
				time.Sleep(60 * time.Second)
				continue
			}
			return err
		}

		var resp models.HimalayasResponse
		if err := json.Unmarshal([]byte(dataStr), &resp); err != nil {
			return err
		}
		if len(resp.Jobs) == 0 {
			break
		}

		for _, hJob := range resp.Jobs {
			item := models.Item{
				Title:       hJob.Title,
				Company:     hJob.CompanyName,
				Description: hJob.Description,
				Link:        hJob.ApplicationLink,
				Guid:        hJob.Guid,
				Location:    "",
			}
			if hJob.ExpiryDate > 0 {
				item.ExpiresAt = time.Unix(hJob.ExpiryDate, 0).Format("2006-01-02 15:04:05")
			}
			if len(hJob.LocationRestrictions) > 0 {
				item.Location = hJob.LocationRestrictions[0]
				for i := 1; i < len(hJob.LocationRestrictions); i++ {
					item.Location += ", " + hJob.LocationRestrictions[i]
				}
			} else {
				item.Location = "Worldwide"
			}

			if !isDuplicate(item) {
				process(item)
				count++
				if p.MaxJobs > 0 && count >= p.MaxJobs {
					break
				}
			}
		}

		offset += limit
		if offset >= resp.TotalCount {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}
