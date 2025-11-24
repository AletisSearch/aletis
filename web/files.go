package web

import (
	"embed"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io/fs"
)

//go:embed dist/**
var staticFs embed.FS

var dynamicFS fs.FS

type ManifestChunk struct {
	Src            string   `json:"src,omitempty"`
	File           string   `json:"file"`
	CSS            []string `json:"css,omitempty"`
	Assets         []string `json:"assets,omitempty"`
	IsEntry        bool     `json:"isEntry,omitempty"`
	Name           string   `json:"name,omitempty"`
	Names          []string `json:"names,omitempty"`
	IsDynamicEntry bool     `json:"isDynamicEntry,omitempty"`
	Imports        []string `json:"imports,omitempty"`
	DynamicImports []string `json:"dynamicImports,omitempty"`
}

type Manifest map[string]ManifestChunk

var AssetsFs fs.FS

var ErrFileNotFound = errors.New("file not found")

var pManifest Manifest

func GetAssetUri(name string) (string, error) {
	if IsDev {
		return "/assets/" + name, nil
	}
	return GetFilepath(name)
}

func GetFilepath(name string) (string, error) {
	var (
		chunk ManifestChunk
		ok    bool
	)
	if IsDev {
		f, err := dynamicFS.Open(".vite/manifest.json")
		if err != nil {
			return "", fmt.Errorf("failed to open manifest.json in GetAssetUri: %w", err)
		}
		defer f.Close()

		var m Manifest

		if err = json.UnmarshalRead(f, &m); err != nil {
			return "", fmt.Errorf("failed to unmarshal manifest.json in GetAssetUr: %w", err)
		}
		chunk, ok = m[name]
	} else {
		chunk, ok = pManifest[name]
	}

	if !ok {
		return "", fmt.Errorf("[%s] error: %w", name, ErrFileNotFound)
	}
	return "/" + chunk.File, nil
}

func GetLinkPreload() (string, error) {
	styleLink, err := GetAssetUri("main.css")
	if err != nil {
		return "", err
	}

	return "<" + styleLink + ">;rel=preload;as=style;fetchpriority=high", nil
}
