package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type GraphDocument struct {
	Context any         `json:"@context,omitempty"`
	Graph   []GraphNode `json:"@graph,omitempty"`
}

type GraphNode struct {
	Type       string     `json:"@type"`
	ID         string     `json:"@id,omitempty"`
	URL        string     `json:"url,omitempty"`
	Name       string     `json:"name,omitempty"`
	MainEntity *ItemList  `json:"mainEntity,omitempty"`
	About      any        `json:"about,omitempty"`
	Raw        struct{}   `json:"-"`
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
	Context            any             `json:"@context,omitempty"`
	Type               any             `json:"@type,omitempty"`
	ID                 string          `json:"@id,omitempty"`
	URL                string          `json:"url,omitempty"`
	Name               string          `json:"name,omitempty"`
	Description        string          `json:"description,omitempty"`
	Image              any             `json:"image,omitempty"`
	Telephone          string          `json:"telephone,omitempty"`
	Latitude           any             `json:"latitude,omitempty"`
	Longitude          any             `json:"longitude,omitempty"`
	NumberOfRooms      any             `json:"numberOfRooms,omitempty"`
	NumberOfBathrooms  any             `json:"numberOfBathroomsTotal,omitempty"`
	FloorSize          json.RawMessage `json:"floorSize,omitempty"`
	Address            json.RawMessage `json:"address,omitempty"`
	Geo                json.RawMessage `json:"geo,omitempty"`
	PetsAllowed        any             `json:"petsAllowed,omitempty"`
	AmenityFeature     json.RawMessage `json:"amenityFeature,omitempty"`
	AggregateRating    json.RawMessage `json:"aggregateRating,omitempty"`
	ContainedInPlace   json.RawMessage `json:"containedInPlace,omitempty"`
	Offers             json.RawMessage `json:"offers,omitempty"`
	PotentialAction    json.RawMessage `json:"potentialAction,omitempty"`
	AdditionalProperty json.RawMessage `json:"additionalProperty,omitempty"`
}

func main() {
	path := "ldjson/0.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("read %s: %v", path, err)
	}

	var graphDoc GraphDocument
	if err := json.Unmarshal(raw, &graphDoc); err != nil {
		log.Fatalf("unmarshal graph doc: %v", err)
	}

	var itemList *ItemList
	for _, node := range graphDoc.Graph {
		if node.MainEntity != nil && node.MainEntity.Type == "ItemList" {
			itemList = node.MainEntity
			break
		}
	}
	if itemList == nil {
		log.Fatalf("no ItemList found in %s", path)
	}

	fmt.Printf("source=%s\n", path)
	fmt.Printf("graph_nodes=%d\n", len(graphDoc.Graph))
	fmt.Printf("number_of_items=%d\n", itemList.NumberOfItems)
	fmt.Printf("parsed_items=%d\n\n", len(itemList.ItemListElement))

	leads := make([]ApartmentLead, 0, len(itemList.ItemListElement))
	for _, li := range itemList.ItemListElement {
		var lead ApartmentLead
		if len(li.Item) > 0 {
			if err := json.Unmarshal(li.Item, &lead); err != nil {
				log.Fatalf("unmarshal apartment item at position %d: %v", li.Position, err)
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

	out, err := json.MarshalIndent(leads[0], "", "  ")
	if err != nil {
		log.Fatalf("marshal sample lead: %v", err)
	}

	fmt.Println("first_apartment_lead=")
	fmt.Println(string(out))
}
