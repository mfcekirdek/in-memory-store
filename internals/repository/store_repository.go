package repository

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const JSONFileSuffix = "-data.json"

type StoreRepository interface {
	Get(key string) string
	Set(key string, value string)
	Flush() map[string]string
	GetStore() map[string]string
}

type storeRepository struct {
	store map[string]string
}

func NewStoreRepository(path string) StoreRepository {
	repository := &storeRepository{}
	repository.store = repository.loadStoreDataFromFile(path)
	return repository
}

func (s *storeRepository) GetStore() map[string]string {
	return s.store
}

func (s *storeRepository) Flush() map[string]string {
	s.store = map[string]string{}
	return s.store
}

func (s *storeRepository) Get(key string) string {
	if value, ok := s.store[key]; ok {
		return value
	}
	return ""
}

func (s *storeRepository) Set(key string, value string) {
	s.store[key] = value
}

func (s *storeRepository) loadStoreDataFromFile(path string) map[string]string {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Println("Could not create storage directory", err)
		return map[string]string{}
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("Could not read files on the path", err)
		return map[string]string{}
	}

	jsonFilePath := findJSONFilePath(path, files)
	store := loadJSONFileToMap(jsonFilePath)

	return store
}

func loadJSONFileToMap(jsonFilePath string) map[string]string {
	store := map[string]string{}
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		log.Println(err)
		return store
	}
	log.Printf("Successfully Opened the json file %s", jsonFilePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &store)
	if err != nil {
		log.Println(err)
	}
	return store
}

func findJSONFilePath(dir string, files []fs.FileInfo) string {
	storeFiles := filterValidDataFiles(files)
	SortTimestampDescend(storeFiles)

	if len(storeFiles) > 0 {
		return filepath.Join(dir, storeFiles[0].Name())
	}
	return ""
}

func SortTimestampDescend(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		l1, _ := strconv.Atoi(getTimestampFromFilename(files[i].Name()))
		l2, _ := strconv.Atoi(getTimestampFromFilename(files[j].Name()))
		return l1 > l2
	})
}

func getTimestampFromFilename(filename string) string {
	timestamp := filename[:len(filename)-len(JSONFileSuffix)]
	return timestamp
}

func filterValidDataFiles(files []os.FileInfo) []os.FileInfo {
	result := make([]os.FileInfo, 0)
	for _, file := range files {
		if strings.Contains(file.Name(), JSONFileSuffix) {
			timestamp := getTimestampFromFilename(file.Name())
			if _, err := strconv.Atoi(timestamp); err == nil {
				result = append(result, file)
			}
		}
	}
	return result
}
