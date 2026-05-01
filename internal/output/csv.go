package output

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"github.com/builderorin/aptscom/internal/parser"
)

func WriteCSV(path string, leads []parser.ApartmentLead) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{
		"name",
		"url",
		"telephone",
		"price",
		"low_price",
		"high_price",
		"street_address",
		"address_locality",
		"address_region",
		"postal_code",
		"address_country",
		"latitude",
		"longitude",
		"amenity_features",
	}
	if err := w.Write(headers); err != nil {
		return err
	}

	for _, lead := range leads {
		row := []string{
			lead.Name,
			lead.URL,
			lead.Telephone,
			intPtrToString(offerPrice(lead.Offers)),
			intPtrToString(offerLowPrice(lead.Offers)),
			intPtrToString(offerHighPrice(lead.Offers)),
			addressField(lead.MainEntity, func(a *parser.PostalAddress) string { return a.StreetAddress }),
			addressField(lead.MainEntity, func(a *parser.PostalAddress) string { return a.AddressLocality }),
			addressField(lead.MainEntity, func(a *parser.PostalAddress) string { return a.AddressRegion }),
			addressField(lead.MainEntity, func(a *parser.PostalAddress) string { return a.PostalCode }),
			addressField(lead.MainEntity, func(a *parser.PostalAddress) string { return a.AddressCountry }),
			geoField(lead.MainEntity, true),
			geoField(lead.MainEntity, false),
			amenitiesString(lead.MainEntity),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return w.Error()
}

func offerPrice(o *parser.Offer) *int {
	if o == nil {
		return nil
	}
	return o.Price
}

func offerLowPrice(o *parser.Offer) *int {
	if o == nil {
		return nil
	}
	return o.LowPrice
}

func offerHighPrice(o *parser.Offer) *int {
	if o == nil {
		return nil
	}
	return o.HighPrice
}

func intPtrToString(v *int) string {
	if v == nil {
		return ""
	}
	return strconv.Itoa(*v)
}

func addressField(main *parser.ApartmentMainEntity, f func(*parser.PostalAddress) string) string {
	if main == nil || main.Address == nil {
		return ""
	}
	return f(main.Address)
}

func geoField(main *parser.ApartmentMainEntity, lat bool) string {
	if main == nil || main.Geo == nil {
		return ""
	}
	if lat {
		return strconv.FormatFloat(main.Geo.Latitude, 'f', -1, 64)
	}
	return strconv.FormatFloat(main.Geo.Longitude, 'f', -1, 64)
}

func amenitiesString(main *parser.ApartmentMainEntity) string {
	if main == nil || len(main.AmenityFeatures) == 0 {
		return ""
	}
	parts := make([]string, 0, len(main.AmenityFeatures))
	for _, a := range main.AmenityFeatures {
		if a.Name != "" {
			parts = append(parts, a.Name)
		}
	}
	return strings.Join(parts, "; ")
}
