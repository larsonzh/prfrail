package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"sort"
)

type fixture struct {
	SchemaVersion  string       `json:"schemaVersion"`
	Source         *source      `json:"source,omitempty"`
	BaselineSHA256 string       `json:"baselineSha256,omitempty"`
	Mode           string       `json:"mode,omitempty"`
	Cases          []shadowCase `json:"cases"`
}

type source struct {
	Repository     string `json:"repository"`
	Path           string `json:"path"`
	SHA256         string `json:"sha256"`
	Bytes          int    `json:"bytes"`
	ExportedAt     string `json:"exportedAt"`
	Classification string `json:"classification"`
}

type shadowCase struct {
	ID                    string `json:"id"`
	Input                 input  `json:"input"`
	Result                result `json:"result"`
	FailureClassification string `json:"failureClassification"`
}

type input struct {
	Query    string `json:"query"`
	StartRIR string `json:"startRir"`
}

type result struct {
	Authority string `json:"authority"`
	Host      string `json:"host"`
}

func main() {
	baselinePath := flag.String("baseline", "baseline.json", "path to the frozen whois baseline")
	shadowPath := flag.String("shadow", "shadow.json", "path to the read-only shadow result")
	flag.Parse()

	if err := compareFiles(*baselinePath, *shadowPath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("PASS: 9/9 cases preserve input, result, and failure classification")
}

func compareFiles(baselinePath, shadowPath string) error {
	baseline, err := loadFixture(baselinePath)
	if err != nil {
		return fmt.Errorf("baseline: %w", err)
	}
	shadow, err := loadFixture(shadowPath)
	if err != nil {
		return fmt.Errorf("shadow: %w", err)
	}
	if err := validateEnvelope(baseline, shadow); err != nil {
		return err
	}
	return compareFixtures(baseline, shadow)
}

func validateEnvelope(baseline, shadow fixture) error {
	if baseline.Source == nil || baseline.Source.SHA256 == "" {
		return errors.New("baseline-source-digest-missing")
	}
	if shadow.BaselineSHA256 != baseline.Source.SHA256 {
		return errors.New("baseline-digest-mismatch")
	}
	if shadow.Mode != "offline-read-only" {
		return errors.New("invalid-shadow-mode")
	}
	return nil
}

func loadFixture(path string) (fixture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return fixture{}, err
	}
	var decoded fixture
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&decoded); err != nil {
		return fixture{}, err
	}
	if decoded.SchemaVersion != "1.0.0" {
		return fixture{}, fmt.Errorf("unsupported schemaVersion %q", decoded.SchemaVersion)
	}
	if len(decoded.Cases) == 0 {
		return fixture{}, errors.New("fixture has no cases")
	}
	return decoded, nil
}

func compareFixtures(baseline, shadow fixture) error {
	baselineCases, err := indexCases(baseline.Cases)
	if err != nil {
		return fmt.Errorf("baseline: %w", err)
	}
	shadowCases, err := indexCases(shadow.Cases)
	if err != nil {
		return fmt.Errorf("shadow: %w", err)
	}
	if len(baselineCases) != len(shadowCases) {
		return fmt.Errorf("case-count-mismatch: baseline=%d shadow=%d", len(baselineCases), len(shadowCases))
	}

	ids := make([]string, 0, len(baselineCases))
	for id := range baselineCases {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		expected := baselineCases[id]
		actual, ok := shadowCases[id]
		if !ok {
			return fmt.Errorf("missing-case: %s", id)
		}
		if !reflect.DeepEqual(expected.Input, actual.Input) {
			return fmt.Errorf("input-mismatch: %s", id)
		}
		if !reflect.DeepEqual(expected.Result, actual.Result) {
			return fmt.Errorf("result-mismatch: %s", id)
		}
		if expected.FailureClassification != actual.FailureClassification {
			return fmt.Errorf("failure-classification-mismatch: %s", id)
		}
	}
	return nil
}

func indexCases(cases []shadowCase) (map[string]shadowCase, error) {
	indexed := make(map[string]shadowCase, len(cases))
	for _, item := range cases {
		if item.ID == "" || item.Input.Query == "" || item.Input.StartRIR == "" || item.Result.Authority == "" || item.Result.Host == "" || item.FailureClassification == "" {
			return nil, errors.New("case has an empty required field")
		}
		if _, exists := indexed[item.ID]; exists {
			return nil, fmt.Errorf("duplicate-case: %s", item.ID)
		}
		indexed[item.ID] = item
	}
	return indexed, nil
}
