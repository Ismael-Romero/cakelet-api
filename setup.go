package main

import (
	"errors"
	"fmt"
	"os"
)

func Setup(cfg *Config) {

	var err error

	if cfg.MediaSource.DirMedia == "" {
		panic("Skipping directory creation: folder path DirMedia is empty.")
	}

	if cfg.MediaSource.DirMedia == "" {
		panic("Skipping directory creation: folder path DirLogs is empty.")
	}

	err = makeDirectory(cfg.MediaSource.DirMedia)
	if err != nil {
		panic(err)
	}

	err = makeDirectory(cfg.MediaSource.DirLogs)
	if err != nil {
		panic(err)
	}

}

func makeDirectory(folder string) error {

	info, err := os.Stat(folder)

	if err == nil && info.IsDir() {
		fmt.Printf("Directory '%s' already exists. Skipping creation to protect data.\n", folder)
		return nil
	}

	if err == nil && !info.IsDir() {
		return fmt.Errorf("path '%s' already exists but it is a file, not a directory", folder)
	}

	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error checking directory '%s': %w", folder, err)
	}

	err = os.MkdirAll(folder, 0700)
	if err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", folder, err)
	}

	fmt.Printf("Directory '%s' created successfully.\n", folder)
	return nil
}
