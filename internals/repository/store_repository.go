package repository

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type StoreRepository interface {
	Get(key string) string
	Set(key string, value string)
	Flush()
}

type storeRepository struct {
	storageDirPath string
	flushInterval  int
	store          map[string]string
}

func NewStoreRepository(path string, flushInterval int) StoreRepository {
	repository := &storeRepository{storageDirPath: path, flushInterval: flushInterval}
	repository.store = repository.loadStoreDataFromFile(path)
	return repository
}

func (s *storeRepository) Flush() {
	s.store = map[string]string{}
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

func findJSONFilePath(dirPath string, files []fs.FileInfo) string {
	suffix := "-data.json"
	SortFileNameDescend(files)
	for _, file := range files {
		if strings.Contains(file.Name(), suffix) {
			filePath := filepath.Join(dirPath, file.Name())
			log.Println(file.Name(), file.Size())
			return filePath
		}
	}
	log.Println("Json file not found on the given path..")
	return ""
}

func SortFileNameDescend(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})
}
