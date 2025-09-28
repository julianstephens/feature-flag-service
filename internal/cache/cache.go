package cache

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/julianstephens/go-utils/helpers"
)

const cacheDirName = "featurectl"

// Dir returns the full path to the cache directory.
func Dir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, cacheDirName), nil
}

// ensureDir ensures the cache directory exists.
func ensureDir() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	if err := helpers.Ensure(dir, true); err != nil {
		return "", err
	}
	return dir, nil
}

// WriteBytes writes raw bytes to a named cache file.
func WriteBytes(name string, data []byte) error {
	dir, err := ensureDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, name)
	return os.WriteFile(path, data, 0o600)
}

// ReadBytes reads raw bytes from a named cache file.
func ReadBytes(name string) ([]byte, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, name)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return data, err
}

// WriteJSON writes a value (as JSON) to a named cache file.
func WriteJSON[T any](name string, value T) error {
	dir, err := ensureDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, name)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	return enc.Encode(value)
}

// ReadJSON reads a value (from JSON) from a named cache file.
func ReadJSON[T any](name string, out *T) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, name)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	return dec.Decode(out)
}

// Remove removes a named cache file.
func Remove(name string) error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, name)
	return os.Remove(path)
}

// Clear removes all cache files in the featurectl cache directory.
func Clear() error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	files, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, fname := range files {
		if err := os.Remove(filepath.Join(dir, fname)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}
