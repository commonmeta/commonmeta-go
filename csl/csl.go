// Package csl provides a function to read CSL JSON and convert it to commonmeta.
package csl

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/schemautils"
	"github.com/xeipuuv/gojsonschema"
)

type content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type CSL struct {
	ID             string             `json:"id"`
	Type           string             `json:"type"`
	Abstract       string             `json:"abstract,omitempty"`
	Accessed       map[string][][]int `json:"accessed,omitempty"`
	Author         []Author           `json:"author,omitempty"`
	ContainerTitle string             `json:"container-title,omitempty"`
	DOI            string             `json:"DOI,omitempty"`
	ISSN           string             `json:"ISSN,omitempty"`
	Issue          string             `json:"issue,omitempty"`
	Issued         map[string][][]int `json:"issued,omitempty"`
	Keyword        string             `json:"keyword,omitempty"`
	Language       string             `json:"language,omitempty"`
	License        string             `json:"license,omitempty"`
	Page           string             `json:"page,omitempty"`
	Publisher      string             `json:"publisher,omitempty"`
	Submitted      map[string][][]int `json:"submitted,omitempty"`
	Title          string             `json:"title,omitempty"`
	URL            string             `json:"URL,omitempty"`
	Version        string             `json:"version,omitempty"`
	Volume         string             `json:"volume,omitempty"`
}

type Author struct {
	Given   string `json:"given,omitempty"`
	Family  string `json:"family,omitempty"`
	Literal string `json:"literal,omitempty"`
}

var CMToCSLMappings = map[string]string{
	"Article":               "article",
	"JournalArticle":        "article-journal",
	"Book":                  "book",
	"BookChapter":           "chapter",
	"Collection":            "collection",
	"Dataset":               "dataset",
	"Document":              "document",
	"Entry":                 "entry",
	"Event":                 "event",
	"Figure":                "figure",
	"Image":                 "graphic",
	"LegalDocument":         "legal_case",
	"Manuscript":            "manuscript",
	"Map":                   "map",
	"Audiovisual":           "motion_picture",
	"Patent":                "patent",
	"Performance":           "performance",
	"Journal":               "periodical",
	"PersonalCommunication": "personal_communication",
	"Report":                "report",
	"Review":                "review",
	"Software":              "software",
	"Presentation":          "speech",
	"Standard":              "standard",
	"Dissertation":          "thesis",
	"WebPage":               "webpage",
}

// Read reads CSL JSON and converts it to commonmeta.
func Read(content content) (commonmeta.Data, error) {
	var data commonmeta.Data

	data.ID = content.ID
	return data, nil
}

// Convert converts commonmeta metadata to CSL JSON.
func Convert(data commonmeta.Data) (CSL, error) {
	var csl CSL

	csl.ID = data.ID
	csl.Type = CMToCSLMappings[data.Type]
	if data.Type == "Software" && data.Version != "" {
		csl.Type = "book"
	} else if csl.Type == "" {
		csl.Type = "Document"
	}
	csl.ContainerTitle = data.Container.Title
	doi, _ := doiutils.ValidateDOI(data.ID)
	csl.DOI = doi
	csl.Issue = data.Container.Issue
	if len(data.Subjects) > 0 {
		for _, subject := range data.Subjects {
			if subject.Subject != "" {
				csl.Keyword += subject.Subject
			}
		}
	}
	csl.Language = data.Language
	csl.Page = data.Container.Pages()
	if len(data.Titles) > 0 {
		csl.Title = data.Titles[0].Title
	}
	csl.URL = data.URL
	csl.Volume = data.Container.Volume
	if len(data.Contributors) > 0 {
		var author Author
		for _, contributor := range data.Contributors {
			if slices.Contains(contributor.ContributorRoles, "Author") {
				if contributor.FamilyName != "" {
					author = Author{
						Given:  contributor.GivenName,
						Family: contributor.FamilyName,
					}
				} else {
					author = Author{
						Literal: contributor.Name,
					}

				}
				csl.Author = append(csl.Author, author)
			}
		}
	}

	if data.Date.Published != "" {
		csl.Issued = dateutils.GetDateParts(data.Date.Published)
	}
	if data.Date.Submitted != "" {
		csl.Submitted = dateutils.GetDateParts(data.Date.Submitted)
	}
	if data.Date.Accessed != "" {
		csl.Accessed = dateutils.GetDateParts(data.Date.Accessed)
	}

	if len(data.Descriptions) > 0 {
		csl.Abstract = data.Descriptions[0].Description
	}
	csl.Publisher = data.Publisher.Name
	csl.Version = data.Version

	return csl, nil
}

// Write writes commonmeta metadata.
func Write(data commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	csl, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err := json.Marshal(csl)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output, "csl-data")
	if !validation.Valid() {
		return nil, validation.Errors()
	}

	return output, nil
}

// WriteList writes a list of commonmeta metadata.
func WriteList(list []commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	var cslList []CSL
	for _, data := range list {
		csl, err := Convert(data)
		if err != nil {
			fmt.Println(err)
		}
		cslList = append(cslList, csl)
	}
	output, err := json.Marshal(cslList)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output, "csl-data")
	if !validation.Valid() {
		return nil, validation.Errors()
	}

	return output, nil
}
