package main

import (
	"anuncios/internal/models"
	"anuncios/internal/parser"
	"anuncios/internal/providers"
	"anuncios/internal/storage"
	"anuncios/internal/tracks"
	"anuncios/internal/utils"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	unlimited := flag.Bool("i", false, "Run without maxjobs limit")
	useTracks := flag.Bool("t", false, "Enable area/track matching and columns")
	flag.Parse()

	allProviders := []providers.Provider{
		&providers.WeWorkRemotely{
			MaxJobs: 10,
			URLs: []string{
				"https://weworkremotely.com/categories/remote-customer-support-jobs.rss",
				"https://weworkremotely.com/categories/remote-product-jobs.rss",
				"https://weworkremotely.com/categories/remote-full-stack-programming-jobs.rss",
				"https://weworkremotely.com/categories/remote-back-end-programming-jobs.rss",
				"https://weworkremotely.com/categories/remote-front-end-programming-jobs.rss",
				"https://weworkremotely.com/categories/remote-programming-jobs.rss",
				"https://weworkremotely.com/categories/remote-sales-and-marketing-jobs.rss",
				"https://weworkremotely.com/categories/remote-management-and-finance-jobs.rss",
				"https://weworkremotely.com/categories/remote-design-jobs.rss",
				"https://weworkremotely.com/categories/remote-devops-sysadmin-jobs.rss",
				"https://weworkremotely.com/categories/all-other-remote-jobs.rss",
				"https://weworkremotely.com/remote-jobs.rss",
			},
		},
		&providers.Remotive{
			MaxJobs: 10,
			URLs: []string{
				"https://remotive.com/remote-jobs/feed/software-development",
				"https://remotive.com/remote-jobs/feed/customer-service",
				"https://remotive.com/remote-jobs/feed/design",
				"https://remotive.com/remote-jobs/feed/marketing",
				"https://remotive.com/remote-jobs/feed/sales-business",
				"https://remotive.com/remote-jobs/feed/product",
				"https://remotive.com/remote-jobs/feed/project-management",
				"https://remotive.com/remote-jobs/feed/ai-ml",
				"https://remotive.com/remote-jobs/feed/data",
				"https://remotive.com/remote-jobs/feed/devops",
				"https://remotive.com/remote-jobs/feed/finance",
				"https://remotive.com/remote-jobs/feed/human-resources",
				"https://remotive.com/remote-jobs/feed/qa",
				"https://remotive.com/remote-jobs/feed/writing",
				"https://remotive.com/remote-jobs/feed/legal",
				"https://remotive.com/remote-jobs/feed/medical",
				"https://remotive.com/remote-jobs/feed/education",
				"https://remotive.com/remote-jobs/feed/all-others",
				"https://remotive.com/remote-jobs/feed",
			},
		},
		&providers.FindJobIT{
			MaxJobs: 10,
			URLs: []string{
				"https://findjobit.com/jobs/role/ingeniero-calidad-qa/feed",
				"https://findjobit.com/jobs/role/programador-full-stack/feed",
				"https://findjobit.com/jobs/role/ingeniero-de-infrestructura/feed",
				"https://findjobit.com/jobs/role/disenador-ui-ux/feed",
				"https://findjobit.com/jobs/role/ingeniero-backend/feed",
				"https://findjobit.com/jobs/role/ingeniero-devops/feed",
				"https://findjobit.com/jobs/role/scrum-master/feed",
				"https://findjobit.com/jobs/role/coach-agil/feed",
				"https://findjobit.com/jobs/role/programador-frontend/feed",
				"https://findjobit.com/jobs/role/analista-de-negocio/feed",
				"https://findjobit.com/jobs/role/arquitecto-software/feed",
				"https://findjobit.com/jobs/role/lider-tecnico/feed",
				"https://findjobit.com/jobs/feed",
			},
		},
		&providers.RemoteOK{
			MaxJobs: 10,
			URLs:    []string{"https://remoteok.com/rss"},
		},
		&providers.Jobicy{
			MaxJobs: 10,
			URLs:    []string{"https://jobicy.com/feed/job_feed"},
		},
		&providers.Himalayas{
			MaxJobs: 10,
			URLs:    []string{"https://himalayas.app/jobs/api"},
		},
	}

	if *unlimited {
		for _, p := range allProviders {
			switch p := p.(type) {
			case *providers.WeWorkRemotely:
				p.MaxJobs = 0
			case *providers.Remotive:
				p.MaxJobs = 0
			case *providers.FindJobIT:
				p.MaxJobs = 0
			case *providers.RemoteOK:
				p.MaxJobs = 0
			case *providers.Himalayas:
				p.MaxJobs = 0
			case *providers.Jobicy:
				p.MaxJobs = 0
			}
		}
		fmt.Println("Running in UNLIMITED mode (-i set).")
	}

	if *useTracks {
		fmt.Println("Area/Track matching ENABLED (-t set).")
	}

	outputDir := "output"
	discardDir := "discard"
	os.MkdirAll(outputDir, 0755)
	os.MkdirAll(discardDir, 0755)

	// Load all existing CSVs for deduplication (from both outcome and discard)
	existingIDs, seenData := storage.LoadAllExistingData(".")

	timestamp := time.Now().Format("20060102150405")
	csvFilename := fmt.Sprintf("%s/JOBS_%s.csv", outputDir, timestamp)
	discardFilename := fmt.Sprintf("%s/DISCARD_%s.csv", discardDir, timestamp)

	var csvFile *os.File
	var writer *csv.Writer
	var discardFile *os.File
	var discardWriter *csv.Writer

	defer func() {
		if csvFile != nil {
			csvFile.Close()
		}
		if discardFile != nil {
			discardFile.Close()
		}
	}()

	totalFound, totalProcessed, okCount, issuesCount, alreadyCount, alreadyTitleCount := 0, 0, 0, 0, 0, 0
	missingDataReport := []string{}
	oneYearFromNow := time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05")

	for _, p := range allProviders {
		fmt.Printf("\n--- Processing Provider: %s ---\n", p.GetName())

		err := p.FetchJobs(func(item models.Item) bool {
			// ID Normalization
			jobID := parser.NormalizeJobID(item)
			if jobID != "" && existingIDs[jobID] {
				alreadyCount++
				return true
			}

			companyName := strings.ReplaceAll(item.Company, "\"", "")
			jobTitle := strings.ReplaceAll(item.Title, "\"", "")

			if companyName != "" && jobTitle != "" {
				key := companyName + "|" + jobTitle
				if seenData[key] {
					alreadyTitleCount++
					return true
				}
			}
			return false
		}, func(item models.Item) {
			totalFound++

			// ID Normalization
			jobID := parser.NormalizeJobID(item)

			companyName := strings.ReplaceAll(item.Company, "\"", "")
			jobTitle := strings.ReplaceAll(item.Title, "\"", "")

			if companyName != "" && jobTitle != "" {
				key := companyName + "|" + jobTitle
				seenData[key] = true
			}

			// Country Mapping & NN Fallback logic
			paisID := parser.MapToPaisID(item, companyName, p.GetName())

			// Location Formatting
			finalLocation := parser.FormatLocation(item, p.GetName())

			expiration := utils.ParseDate(item.ExpiresAt)
			if expiration == "" {
				expiration = oneYearFromNow
			}

			desc := utils.CleanHTML(item.Description)
			totalProcessed++

			missingFields := []string{}
			if jobID == "" {
				missingFields = append(missingFields, "job_id")
			}
			if paisID == "" {
				missingFields = append(missingFields, "pais_id")
			}
			if companyName == "" {
				missingFields = append(missingFields, "company_name")
			}
			if jobTitle == "" {
				missingFields = append(missingFields, "title")
			}
			if desc == "" {
				missingFields = append(missingFields, "description")
			}
			if item.Link == "" {
				missingFields = append(missingFields, "apply_url")
			}
			if expiration == "" {
				missingFields = append(missingFields, "expiration_date")
			}
			if finalLocation == "" {
				missingFields = append(missingFields, "location")
			}

			// Track Matching (title + description) - Only if flag -t is set
			areaID, trackName, trackFound := "", "", true
			if *useTracks {
				areaID, trackName, trackFound = tracks.MatchTrack(jobTitle, desc)
			}

			if len(missingFields) > 0 {
				issuesCount++
				missingDataReport = append(missingDataReport, fmt.Sprintf("[%s] Job ID: %s is missing fields: %s", p.GetName(), jobID, strings.Join(missingFields, ", ")))
			} else if *useTracks && !trackFound {
				// DISCARD logic (only if -t is used)
				if discardWriter == nil {
					f, err := os.Create(discardFilename)
					if err != nil {
						fmt.Println("Error creating Discard CSV:", err)
						return
					}
					discardFile = f
					discardWriter = csv.NewWriter(discardFile)
					headers := []string{"area_id", "track", "provider_id", "expiration_date", "pais_id", "job_id", "company_name", "title", "description", "apply_url", "location"}
					discardWriter.Write(headers)
				}
				discardWriter.Write([]string{"", "", p.GetName(), expiration, paisID, jobID, companyName, jobTitle, desc, item.Link, finalLocation})
				discardWriter.Flush()
				existingIDs[jobID] = true
				okCount++
			} else {
				// OK (Matched or skipping tracks)
				if writer == nil {
					f, err := os.Create(csvFilename)
					if err != nil {
						fmt.Println("Error creating JOBS CSV:", err)
						return
					}
					csvFile = f
					writer = csv.NewWriter(csvFile)
					headers := []string{"provider_id", "expiration_date", "pais_id", "job_id", "company_name", "title", "description", "apply_url", "location"}
					if *useTracks {
						headers = append([]string{"area_id", "track"}, headers...)
					}
					writer.Write(headers)
				}

				row := []string{p.GetName(), expiration, paisID, jobID, companyName, jobTitle, desc, item.Link, finalLocation}
				if *useTracks {
					row = append([]string{areaID, trackName}, row...)
				}
				writer.Write(row)
				writer.Flush()
				existingIDs[jobID] = true
				okCount++
			}
		})

		if err != nil {
			fmt.Printf("Error fetching jobs from %s: %v\n", p.GetName(), err)
			continue
		}
	}

	if okCount > 0 {
		fmt.Printf("\nSaved %d new records to %s\n", okCount, csvFilename)
	} else {
		fmt.Println("\nNo new jobs found to save.")
	}

	fmt.Printf("\n--- Multi-Provider RSS/API Conversion Summary ---\n")
	fmt.Printf("Jobs Found: %d\nTotal NEW processed: %d\nRecords OK: %d\nRecords Issues: %d\nAlready in CSV (ID): %d\nAlready in CSV (Title): %d\n", totalFound, totalProcessed, okCount, issuesCount, alreadyCount, alreadyTitleCount)

	if len(missingDataReport) > 0 {
		fmt.Println("\n--- Missing Data Report ---")
		for _, report := range missingDataReport {
			fmt.Println(report)
		}
	} else {
		fmt.Println("\nNo missing data found in new entries!")
	}
}
