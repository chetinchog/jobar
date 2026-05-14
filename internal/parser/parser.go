package parser

import (
	"anuncios/internal/models"
	"anuncios/internal/utils"
	"regexp"
	"strings"
)

var (
	reCountry  = regexp.MustCompile(`(?i)<(?:country|país)[^>]*>([\s\S]*?)</(?:country|país)>`)
	reTitle    = regexp.MustCompile(`(?i)<(?:title|título)[^>]*>([\s\S]*?)</(?:title|título)>`)
	reDesc     = regexp.MustCompile(`(?i)<(?:description|descripción)[^>]*>([\s\S]*?)</(?:description|descripción)>`)
	reState    = regexp.MustCompile(`(?i)<(?:state|estado)[^>]*>([\s\S]*?)</(?:state|estado)>`)
	reRegion   = regexp.MustCompile(`(?i)<(?:region|región)[^>]*>([\s\S]*?)</(?:region|región)>`)
	reLocation = regexp.MustCompile(`(?i)<(?:location|ubicación|city)[^>]*>([\s\S]*?)</(?:location|ubicación|city)>`)
	reCompany  = regexp.MustCompile(`(?i)<(?:company|empresa|organization|org)[^>]*>([\s\S]*?)</(?:company|empresa|organization|org)>`)
	reJobIDTag = regexp.MustCompile(`(?i)<(?:jobid|id)[^>]*>([\s\S]*?)</(?:jobid|id)>`)
)

func ExtractTagCustom(content, openTag, closeTag string) string {
	start := strings.Index(content, openTag)
	if start == -1 {
		openTagPrefix := strings.TrimRight(openTag, ">") + " "
		start = strings.Index(content, openTagPrefix)
		if start == -1 {
			return ""
		}
		endOpen := strings.Index(content[start:], ">")
		if endOpen == -1 {
			return ""
		}
		start = start + endOpen + 1
	} else {
		start += len(openTag)
	}

	end := strings.Index(content[start:], closeTag)
	if end == -1 {
		return ""
	}
	return strings.TrimSpace(content[start : start+end])
}

func ExtractTag(content, tag string) string {
	var re *regexp.Regexp
	switch strings.ToLower(tag) {
	case "country", "pais", "país":
		re = reCountry
	case "title", "título", "titulo":
		re = reTitle
	case "description", "descripción", "descripcion":
		re = reDesc
	case "state", "estado":
		re = reState
	case "region", "región":
		re = reRegion
	case "location", "ubicación", "ubicacion":
		re = reLocation
	case "company", "empresa":
		re = reCompany
	case "jobid":
		re = reJobIDTag
	default:
		return ExtractTagCustom(content, "<"+tag+">", "</"+tag+">")
	}

	if re != nil {
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			data := matches[1]
			data = strings.TrimPrefix(data, "<![CDATA[")
			data = strings.TrimSuffix(data, "]]>")
			return strings.TrimSpace(data)
		}
	}

	return ExtractTagCustom(content, "<"+tag+">", "</"+tag+">")
}

func NormalizeJobID(item models.Item) string {
	jobID := item.JobIdTag
	if jobID == "" {
		jobID = strings.TrimRight(item.Guid, "/")
		if lastSlash := strings.LastIndex(jobID, "/"); lastSlash != -1 {
			jobID = jobID[lastSlash+1:]
		}
	}
	return jobID
}

func MapToPaisID(item models.Item, companyName string, providerName string) string {
	paisID := utils.ExtractISO(item.Country)
	if paisID == "" {
		paisID = utils.StateToISO(utils.CleanForSearch(item.Country))
	}
	if paisID == "" {
		paisID = utils.StateToISO(utils.CleanForSearch(item.Estado))
	}
	if paisID == "" {
		paisID = utils.StateToISO(utils.CleanForSearch(item.Region))
	}
	if paisID == "" {
		paisID = utils.StateToISO(utils.CleanForSearch(item.Location))
	}
	if paisID == "" {
		lowSearchClean := utils.CleanForSearch(item.Description + " " + companyName)
		patterns := []string{"sede central:", "sede en:", "sede en ", "sede: ", "con sede en ", "based in ", "ubicación:", "ubicacion:", "headquarters:"}
		for _, op := range patterns {
			if idx := strings.Index(lowSearchClean, op); idx != -1 {
				endIdx := idx + 80
				if endIdx > len(lowSearchClean) {
					endIdx = len(lowSearchClean)
				}
				if iso := utils.StateToISO(lowSearchClean[idx:endIdx]); iso != "" {
					paisID = iso
					break
				}
			}
		}
		if paisID == "" {
			paisID = utils.StateToISO(lowSearchClean)
		}
	}

	// NN Fallback logic
	if paisID == "" {
		lowLoc := strings.ToLower(item.Region + " " + item.Country + " " + item.Location + " " + item.Description)
		if strings.Contains(lowLoc, "anywhere in the world") || strings.Contains(lowLoc, "worldwide") || strings.Contains(lowLoc, "global") || strings.Contains(lowLoc, "remoto") || providerName == "RemoteOK" {
			paisID = "NN"
		}
	}
	return paisID
}

func FormatLocation(item models.Item, providerName string) string {
	finalLocation := item.Location
	if finalLocation == "" {
		if providerName == "RemoteOK" {
			finalLocation = "Worldwide"
		} else {
			finalLocation = item.Region
		}
	}
	if item.Country != "" {
		if finalLocation != "" && !strings.Contains(strings.ToLower(finalLocation), strings.ToLower(item.Country)) {
			finalLocation += " (" + item.Country + ")"
		} else if finalLocation == "" {
			finalLocation = item.Country
		}
	}
	if item.Estado != "" {
		if finalLocation != "" && !strings.Contains(strings.ToLower(finalLocation), strings.ToLower(item.Estado)) {
			finalLocation += " - " + item.Estado
		} else if finalLocation == "" {
			finalLocation = item.Estado
		}
	}
	return finalLocation
}

