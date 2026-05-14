package models

type Item struct {
	Title       string
	Region      string
	Country     string
	Estado      string
	Description string
	PubDate     string
	ExpiresAt   string
	Guid        string
	Link        string
	Location    string
	Company     string
	JobIdTag    string
}

type HimalayasResponse struct {
	TotalCount int            `json:"totalCount"`
	Limit      int            `json:"limit"`
	Offset     int            `json:"offset"`
	Jobs       []HimalayasJob `json:"jobs"`
}

type HimalayasJob struct {
	Title                string   `json:"title"`
	CompanyName          string   `json:"companyName"`
	Description          string   `json:"description"`
	PubDate              int64    `json:"pubDate"`
	ExpiryDate           int64    `json:"expiryDate"`
	ApplicationLink      string   `json:"applicationLink"`
	Guid                 string   `json:"guid"`
	LocationRestrictions []string `json:"locationRestrictions"`
}

type Provider struct {
	Name    string
	URLs    []string
	Type    string
	MaxJobs int
}
