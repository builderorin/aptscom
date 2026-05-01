package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/builderorin/aptscom/internal/client"
	"github.com/builderorin/aptscom/internal/output"
	"github.com/builderorin/aptscom/internal/parser"
	"github.com/builderorin/aptscom/internal/slug"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: aptscom <apartments.com-url>")
	}

	url := os.Args[1]
	base := slug.FromURL(url)

	html, err := client.FetchHTML(url)
	if err != nil {
		log.Fatalf("fetch html: %v", err)
	}

	leads, err := parser.ParseApartmentLeads(html)
	if err != nil {
		log.Fatalf("parse apartment leads: %v", err)
	}

	jsonPath := base + ".json"
	csvPath := base + ".csv"

	jsonBytes, err := json.MarshalIndent(leads, "", "  ")
	if err != nil {
		log.Fatalf("marshal json: %v", err)
	}
	if err := os.WriteFile(jsonPath, jsonBytes, 0o644); err != nil {
		log.Fatalf("write json: %v", err)
	}

	if err := output.WriteCSV(csvPath, leads); err != nil {
		log.Fatalf("write csv: %v", err)
	}

	fmt.Printf("wrote %s\n", jsonPath)
	fmt.Printf("wrote %s\n", csvPath)
	fmt.Printf("apartment_leads=%d\n", len(leads))
}
