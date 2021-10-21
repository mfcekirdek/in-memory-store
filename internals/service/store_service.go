package service

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.com/mfcekirdek/in-memory-store/internals/repository"
	"gitlab.com/mfcekirdek/in-memory-store/internals/utils"
)

const JSONFileSuffix = "-data.json"

type StoreService interface {
	Get(key string) (map[string]string, error)
	Set(key string, value string) (map[string]string, bool)
	Flush() map[string]string
}

type storeService struct {
	repository     repository.StoreRepository
	storageDirPath string
}

func NewStoreService(repo repository.StoreRepository, interval int, path string) StoreService {
	service := &storeService{repository: repo, storageDirPath: path}
	store := readStoreDataFromFile(path)
	service.repository.LoadStore(store)
	intervalDuration := 60 * time.Second * time.Duration(interval)
	go backgroundTask(intervalDuration, path, store, saveToJSONFile)
	return service
}

func (s storeService) Get(key string) (map[string]string, error) {
	value := s.repository.Get(key)
	if value == "" {
		return nil, utils.ErrNotFound
	}
	return map[string]string{key: value}, nil
}

func (s storeService) Set(key, value string) (map[string]string, bool) {
	keyAlreadyExist := s.repository.Set(key, value)
	return map[string]string{key: value}, keyAlreadyExist
}

func (s storeService) Flush() map[string]string {
	return s.repository.Flush()
}

//Start a goroutine to save the store to file at intervals
func backgroundTask(interval time.Duration, storageDirPath string, store map[string]string, task func(filepath string, store map[string]string) error) {
	dateTicker := time.NewTicker(interval)
	for now := range dateTicker.C {
		timestamp := now.Unix()
		fileName := fmt.Sprintf("%d%s", timestamp, JSONFileSuffix)
		filePath := filepath.Join(storageDirPath, fileName)
		err := task(filePath, store)
		if err != nil {
			log.Println("Could not write to json file --> ", filePath)
		}
	}
}

func readStoreDataFromFile(path string) map[string]string {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Println("Could not create storage directory", err)
		return map[string]string{}
	}

	files, _ := ioutil.ReadDir(path)
	jsonFilePath := findJSONFilePath(path, files)
	store := saveToMap(jsonFilePath)
	return store
}

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

func saveToMap(jsonFilePath string) map[string]string {
	store := map[string]string{}
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		log.Println("Could not find a json file to load data. -> Initializing an empty store..")
		return nil
	}
	log.Printf("Successfully Opened the json file %s", jsonFilePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &store)
	if err != nil {
		log.Println(err)
		return nil
	}
	return store
}

func findJSONFilePath(dir string, files []fs.FileInfo) string {
	storeFiles := filterValidDataFiles(files)
	sortTimestampDescend(storeFiles)

	if len(storeFiles) > 0 {
		return filepath.Join(dir, storeFiles[0].Name())
	}
	return ""
}

func sortTimestampDescend(files []os.FileInfo) {
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
