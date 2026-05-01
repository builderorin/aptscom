# aptscom

A small CLI for scraping apartment listing leads from Apartments.com category pages.

It fetches the page HTML, reads the first `script[type="application/ld+json"]` block, parses the apartment listing data, and writes both JSON and CSV outputs.

## Install

### Run directly from source

```bash
go run ./cmd/aptscom "https://www.apartments.com/off-campus-housing/ca/santa-clara/santa-clara-university/"
```

### Build locally

```bash
go build -o aptscom ./cmd/aptscom
./aptscom "https://www.apartments.com/off-campus-housing/ca/san-jose/san-jose-state-university/"
```

### Install with Go

```bash
go install github.com/builderorin/aptscom/cmd/aptscom@latest
```

Then run:

```bash
aptscom "https://www.apartments.com/off-campus-housing/ca/santa-clara/santa-clara-university/"
```

## Output

The CLI writes output relative to your current working directory, not the binary location.

For example, if you run this from `~/Desktop/test`:

```bash
cd ~/Desktop/test
aptscom "https://www.apartments.com/off-campus-housing/ca/santa-clara/santa-clara-university/"
```

it will create:

```bash
scraped/json/santa-clara-university.json
scraped/csv/santa-clara-university.csv
```

## What it extracts

The parser currently pulls these fields into CSV:

- `name`
- `url`
- `telephone`
- `price`
- `low_price`
- `high_price`
- `street_address`
- `address_locality`
- `address_region`
- `postal_code`
- `address_country`
- `latitude`
- `longitude`
- `amenity_features`

## Notes

- Uses the first `application/ld+json` script tag on the page.
- `offers` supports both `Offer` and `AggregateOffer` shapes.
- `amenity_features` is currently a flattened string column in the CSV.
