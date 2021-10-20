package service

import (
	"encoding/json"
	"fmt"
	"gitlab.com/mfcekirdek/in-memory-store/internals/repository"
	"gitlab.com/mfcekirdek/in-memory-store/internals/utils"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

const JSONFileSuffix = "-data.json"

type StoreService interface {
	Get(key string) (map[string]string, error)
	Set(key string, value string) map[string]string
	Flush() map[string]string
}

type storeService struct {
	repository     repository.StoreRepository
	storageDirPath string
}

func NewStoreService(repo repository.StoreRepository, flushInterval int, path string) StoreService {
	service := &storeService{repository: repo, storageDirPath: path}
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

func (s *storeService) Set(key, value string) map[string]string {
	s.repository.Set(key, value)
	return map[string]string{key: value}
}

func (s *storeService) Flush() map[string]string {
	return s.repository.Flush()
}

func (s *storeService) BackgroundTask(interval int, task func(filepath string, store map[string]string) error) {
	dateTicker := time.NewTicker(time.Duration(interval) * time.Second)
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

func saveToJSONFile(filePath string, store map[string]string) error {
	log.Println("Saving store to file -> ", filePath)
	file, err := json.MarshalIndent(store, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, file, 0600)
	if err != nil {
		return err
	}
	return nil
}

// todo Repository'i biraz incelt, service'e taşı logicleri
