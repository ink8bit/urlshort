package memory

import (
	"encoding/json"
	"fmt"
	"os"
)

type Record struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      int    `json:"user_id"`
}

func (r *Record) Save(fname string) error {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("error while marshaling data: %w", err)
	}
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0644) //nolint:gosec,gomnd // false positive
	if err != nil {
		return fmt.Errorf("error while opening file: %w", err)
	}
	defer func() {
		err = f.Close()
	}()
	if err != nil {
		return fmt.Errorf("error while closing file: %w", err)
	}
	if _, err := f.WriteString(string(data) + "\n"); err != nil {
		return fmt.Errorf("cannot write data to file: %w", err)
	}
	return nil
}
