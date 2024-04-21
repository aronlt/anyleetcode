package common

import (
	"github.com/aronlt/toolkit/tio"
	"os"
	"path/filepath"
	"sync"
)

const DifficultyFilePath = "data/difficulty.json"
const TagsFilePath = "data/tags.json"
const DataFilePath = "data/data.json"

var fontOnce sync.Once
var fontPath string

var cacheOnce sync.Once
var cachePath string

func GetFontPath() string {
	fontOnce.Do(
		func() {
			dirname, err := os.UserHomeDir()
			if err != nil {
				panic(err)
			}
			dirname = filepath.Join(dirname, "Documents/any-leetcode/")
			err = tio.Mkdirs(dirname)
			if err != nil {
				panic(err)
			}
			fontPath = filepath.Join(dirname, "simkai.ttf")

		})
	return fontPath
}
func GetCachePath() string {
	cacheOnce.Do(
		func() {
			dirname, err := os.UserHomeDir()
			if err != nil {
				panic(err)
			}
			cachePath = filepath.Join(dirname, "Documents/any-leetcode/cache")
			err = tio.Mkdirs(cachePath)
			if err != nil {
				panic(err)
			}
		})
	return cachePath
}
