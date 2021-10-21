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

func NewStoreService(repo repository.StoreRepository, flushInterval int, path string) StoreService {
	service := &storeService{repository: repo, storageDirPath: path}
	store := service.loadStoreDataFromFile(path)
	service.repository.LoadStore(store)
	go service.BackgroundTask(flushInterval, saveToJSONFile)
	return service
}

func (s *storeService) Get(key string) (map[string]string, error) {
	value := s.repository.Get(key)
	if value == "" {
		return nil, utils.ErrNotFound
	}
	return map[string]string{key: value}, nil
}

func (s *storeService) Set(key, value string) (map[string]string, bool) {
	keyAlreadyExist := s.repository.Set(key, value)
	return map[string]string{key: value}, keyAlreadyExist
}

func (s *storeService) Flush() map[string]string {
	return s.repository.Flush()
}

func (s *storeService) BackgroundTask(interval int, task func(filepath string, store map[string]string) error) {
	dateTicker := time.NewTicker(time.Duration(interval) * time.Minute)
	for range dateTicker.C {
		timestamp := time.Now().Unix()
		fileName := fmt.Sprintf("%d%s", timestamp, JSONFileSuffix)
		filePath := filepath.Join(s.storageDirPath, fileName)
		err := task(filePath, s.repository.GetStore())
		if err != nil {
			log.Println("Could not write to json file --> ", filePath)
		}
	}
}

func (s *storeService) loadStoreDataFromFile(path string) map[string]string {
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

func saveToJSONFile(filePath string, store map[string]string) error {
	log.Println("Saving store to file -> ", filePath)
	file, err := json.MarshalIndent(store, "", " ")
	if err != nil {
		return err
	}
	var perm fs.FileMode = 0600
	err = ioutil.WriteFile(filePath, file, perm)
	if err != nil {
		return err
	}
	return nil
}

func loadJSONFileToMap(jsonFilePath string) map[string]string {
	store := map[string]string{}
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		log.Println("Could not find a json file to load data. -> Initializing an empty store..")
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
