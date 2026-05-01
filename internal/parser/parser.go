package parser

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GraphDocument struct {
	Context any         `json:"@context,omitempty"`
	Graph   []GraphNode `json:"@graph,omitempty"`
}

type GraphNode struct {
	Type       string    `json:"@type"`
	ID         string    `json:"@id,omitempty"`
	URL        string    `json:"url,omitempty"`
	Name       string    `json:"name,omitempty"`
	MainEntity *ItemList `json:"mainEntity,omitempty"`
}

type ItemList struct {
	Type            string     `json:"@type"`
	ItemListOrder   string     `json:"itemListOrder,omitempty"`
	NumberOfItems   int        `json:"numberOfItems,omitempty"`
	ItemListElement []ListItem `json:"itemListElement,omitempty"`
}

type ListItem struct {
	Type     string          `json:"@type"`
	Position int             `json:"position"`
	URL      string          `json:"url"`
	Name     string          `json:"name"`
	Item     json.RawMessage `json:"item"`
}

type ApartmentLead struct {
	Context         any                  `json:"@context,omitempty"`
	Type            any                  `json:"@type,omitempty"`
	ID              string               `json:"@id,omitempty"`
	URL             string               `json:"url,omitempty"`
	Name            string               `json:"name,omitempty"`
	Description     string               `json:"description,omitempty"`
	Image           any                  `json:"image,omitempty"`
	Telephone       string               `json:"telephone,omitempty"`
	Offers          *Offer               `json:"offers,omitempty"`
	PotentialAction *Action              `json:"potentialAction,omitempty"`
	MainEntity      *ApartmentMainEntity `json:"mainEntity,omitempty"`
}

type ApartmentMainEntity struct {
	Type            string           `json:"@type,omitempty"`
	ID              string           `json:"@id,omitempty"`
	Name            string           `json:"name,omitempty"`
	HasMap          string           `json:"hasMap,omitempty"`
	Address         *PostalAddress   `json:"address,omitempty"`
	Geo             *GeoCoordinates  `json:"geo,omitempty"`
	AmenityFeatures []AmenityFeature `json:"amenityFeature,omitempty"`
}

type PostalAddress struct {
	Type            string `json:"@type,omitempty"`
	StreetAddress   string `json:"streetAddress,omitempty"`
	AddressLocality string `json:"addressLocality,omitempty"`
	AddressRegion   string `json:"addressRegion,omitempty"`
	PostalCode      string `json:"postalCode,omitempty"`
	AddressCountry  string `json:"addressCountry,omitempty"`
}

type GeoCoordinates struct {
	Type      string  `json:"@type,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type AmenityFeature struct {
	Type  string `json:"@type,omitempty"`
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type Offer struct {
	Type          string `json:"@type,omitempty"`
	PriceCurrency string `json:"priceCurrency,omitempty"`
	Price         *int   `json:"price,omitempty"`
	LowPrice      *int   `json:"lowPrice,omitempty"`
	HighPrice     *int   `json:"highPrice,omitempty"`
	Availability  string `json:"availability,omitempty"`
	OfferCount    int    `json:"offerCount,omitempty"`
}

type Action struct {
	Type string `json:"@type,omitempty"`
}

func ParseApartmentLeads(html string) ([]ApartmentLead, error) {
	raw, err := firstLDJSON(html)
	if err != nil {
		return nil, err
	}

	var graphDoc GraphDocument
	if err := json.Unmarshal(raw, &graphDoc); err != nil {
		return nil, err
	}

	var itemList *ItemList
	for _, node := range graphDoc.Graph {
		if node.MainEntity != nil && node.MainEntity.Type == "ItemList" {
			itemList = node.MainEntity
			break
		}
	}
	if itemList == nil {
		return nil, errors.New("no ItemList found in ld+json graph")
	}

	leads := make([]ApartmentLead, 0, len(itemList.ItemListElement))
	for _, li := range itemList.ItemListElement {
		var lead ApartmentLead
		if len(li.Item) > 0 {
			if err := json.Unmarshal(li.Item, &lead); err != nil {
				return nil, err
			}
		}

		if lead.URL == "" {
			lead.URL = li.URL
		}
		if lead.Name == "" {
			lead.Name = li.Name
		}

		leads = append(leads, lead)
	}

	return leads, nil
}

func firstLDJSON(html string) ([]byte, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	selection := doc.Find(`script[type="application/ld+json"]`).First()
	if selection.Length() == 0 {
		return nil, errors.New("no application/ld+json script tag found")
	}

	raw := strings.TrimSpace(selection.Text())
	if raw == "" {
		return nil, errors.New("first application/ld+json script tag was empty")
	}

	return []byte(raw), nil
}
