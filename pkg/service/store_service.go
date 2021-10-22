// Package service operates the business logic.
// This layer uses repository layer and is used by handler layer.
package service

import (
	"encoding/json"
	"fmt"
	"gitlab.com/mfcekirdek/in-memory-store/pkg/repository"
	"gitlab.com/mfcekirdek/in-memory-store/pkg/utils"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Suffix of valid store files
const JSONFileSuffix = "-data.json"

//StoreService interface has Get, Set and Flush functions.
type StoreService interface {
	Get(key string) (map[string]string, error)
	Set(key string, value string) (map[string]string, bool)
	Flush() map[string]string
}

// storeService implements the StoreService interface.
// Contains StoreRepository and data storage path.
type storeService struct {
	repository     repository.StoreRepository
	storageDirPath string
}

// NewStoreService creates a new storeService instance.
func NewStoreService(repo repository.StoreRepository, interval int, path string) StoreService {
	service := &storeService{repository: repo, storageDirPath: path}
	store := readStoreDataFromFile(path)
	service.repository.LoadStore(store)
	intervalDuration := 60 * time.Second * time.Duration(interval)
	go backgroundTask(intervalDuration, path, store, saveToJSONFile)
	return service
}

// Get function fetches store data from repository layer and returns it.
// Takes <string> key parameter and returns <map{key: value}, nil> if the key is found in the store.
// Returns <nil, error> if not found
func (s *storeService) Get(key string) (map[string]string, error) {
	value := s.repository.Get(key)
	if value == "" {
		return nil, utils.ErrNotFound
	}
	return map[string]string{key: value}, nil
}

// Set function calls repository layer Set function and sets value of an item in the store by using given key,value parameters.
// Takes <string> key and <string> value parameters.
// If key alreadys exists in the store, returns <map{key: value}, true>
// If not returns <map{key: value}, false>
func (s *storeService) Set(key, value string) (map[string]string, bool) {
	keyAlreadyExist := s.repository.Set(key, value)
	return map[string]string{key: value}, keyAlreadyExist
}

// Flush function calls repository layer Flush function, deletes all items in the store and returns the output.
func (s *storeService) Flush() map[string]string {
	return s.repository.Flush()
}

// Starts a goroutine to save the store to file at intervals.
func backgroundTask(interval time.Duration, storageDir string, store map[string]string, task func(filepath string, store map[string]string) error) {
	dateTicker := time.NewTicker(interval)
	for now := range dateTicker.C {
		timestamp := now.Unix()
		fileName := fmt.Sprintf("%d%s", timestamp, JSONFileSuffix)
		filePath := filepath.Join(storageDir, fileName)
		err := task(filePath, store)
		if err != nil {
			log.Println("Could not write to json file --> ", filePath)
		}
	}
}

// Takes <string> directory path parameter.
// Looks for datastore json files whose filenames are in TIMESTAMP-data.json format.
// Finds the most recent one, reads and returns the data from it.
func readStoreDataFromFile(path string) map[string]string {
	_ = os.MkdirAll(path, os.ModePerm)
	files, _ := ioutil.ReadDir(path)
	jsonFilePath := findJSONFilePath(path, files)
	store := saveToMap(jsonFilePath)
	return store
}

// Takes <string, map> filepath and store parameters.
// Writes store contents to the given file.
func saveToJSONFile(filePath string, store map[string]string) error {
	log.Println("Saving store to file -> ", filePath)
	var perm fs.FileMode = 0600
	file, _ := json.MarshalIndent(store, "", " ")
	err := ioutil.WriteFile(filePath, file, perm)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Takes <string> filepath and returns <map> store.
// Reads store data from the given file and returns it
func saveToMap(jsonFilePath string) map[string]string {
	store := map[string]string{}
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		log.Println("Could not find a json file to load data. -> Initializing an empty store..")
		return map[string]string{}
	}
	log.Printf("Successfully Opened the json file %s", jsonFilePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &store)
	if err != nil {
		log.Println(err)
		return map[string]string{}
	}
	return store
}

// Takes <string> directory path and <list> fileInfos.
// Filters valid formatted files in the given fileInfos.
// Sorts them comparing their TIMESTAMPS.
// Returns the concatenation of given directory path and the most recent store file.
func findJSONFilePath(dir string, files []fs.FileInfo) string {
	storeFiles := filterValidDataFiles(files)
	sortTimestampDescend(storeFiles)

	if len(storeFiles) > 0 {
		return filepath.Join(dir, storeFiles[0].Name())
	}
	return ""
}

// Takes <list> fileInfos and sorts them.
func sortTimestampDescend(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		l1, _ := strconv.Atoi(getTimestampFromFilename(files[i].Name()))
		l2, _ := strconv.Atoi(getTimestampFromFilename(files[j].Name()))
		return l1 > l2
	})
}

// Takes <string> filename and returns the timestamp of the file by splitting filename.
func getTimestampFromFilename(filename string) string {
	timestamp := filename[:len(filename)-len(JSONFileSuffix)]
	return timestamp
}

// Takes <list> fileInfos and returns the filtered valid formatted files.
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
