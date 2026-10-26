package report

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Record struct {
	SKU        string
	Qty        int
	PriceCents int // unit price
}

type Exporter interface {
	Name() string
	Export(records []Record) ([]byte, error)
}

type JSONExporter struct{}

type jsonRecord struct {
	SKU        string `json:"sku"`
	Qty        int    `json:"qty"`
	PriceCents int    `json:"price_cents"`
	LineCents  int    `json:"line_cents"`
}

type CSVExporter struct {
	Comma rune // if 0, use ','
}

func validate(records []Record) error {
	for _, v := range records {
		if strings.TrimSpace(v.SKU) == "" {
			return errors.New("empty SKU")
		}
		if v.Qty < 0 || v.PriceCents < 0 {
			return errors.New("Qty or PriceCents less than zero")
		}
	}
	return nil
}

func (JSONExporter) Name() string {
	return "json"
}

func (JSONExporter) Export(records []Record) ([]byte, error) {
	if err := validate(records); err != nil {
		return nil, err
	}

	result := make([]jsonRecord, 0, len(records))
	for _, v := range records {
		result = append(result, jsonRecord{
			SKU:        v.SKU,
			Qty:        v.Qty,
			PriceCents: v.PriceCents,
			LineCents:  v.Qty * v.PriceCents,
		})
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (CSVExporter) Name() string {
	return "csv"
}

func (e CSVExporter) Export(records []Record) ([]byte, error) {
	if err := validate(records); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	
	if e.Comma != 0 {
		w.Comma = e.Comma
	}

	err := w.Write([]string{"sku", "qty", "price_cents", "line_cents"})
	if err != nil {
		return nil, err
	}

	for _, v := range records {
		row := []string{
			v.SKU,
			strconv.Itoa(v.Qty),
			strconv.Itoa(v.PriceCents),
			strconv.Itoa(v.Qty * v.PriceCents),
		}

		if err := w.Write(row); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// RenderAll runs every exporter. Keys are exporter names.
func RenderAll(records []Record, exporters []Exporter) (map[string][]byte, error) {
	result := make(map[string][]byte, len(exporters))
	
	for _, e := range exporters {
		name := e.Name()
		
		if _, ok := result[name]; ok {
			return nil, fmt.Errorf("report: duplicate exporter name %q", name)
		}
		
		data, err := e.Export(records)
		if err != nil {
			return nil, err
		}
		
		result[name] = data
	}
	
	return result, nil
}

func CommaOf(e Exporter) (rune, bool) {
	csvExp, ok := e.(CSVExporter)
	if ok {
		return csvExp.Comma, true
	}
	return 0, false
}