package providers

import "anuncios/internal/models"

type Provider interface {
	GetName() string
	FetchJobs(isDuplicate func(models.Item) bool, process func(models.Item)) error
	GetMaxJobs() int
}
