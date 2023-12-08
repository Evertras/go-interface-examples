package main

import (
	"fmt"
	"os"
)

// GSLDataStore knows how to get GSL data
type GSLDataStore struct {
	championFile string
}

// NewGSLDataStore returns a GSLDataStore ready to tell us about the GSL
func NewGSLDataStore(championFile string) *GSLDataStore {
	return &GSLDataStore{
		championFile,
	}
}

// GetCurrentChampion returns the name of the current GSL champion
func (s *GSLDataStore) GetCurrentChampion() (string, error) {
	contents, err := os.ReadFile(s.championFile)

	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(contents), nil
}
