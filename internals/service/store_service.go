package service

import (
	"fmt"
	"gitlab.com/mfcekirdek/in-memory-store/internals/repository"
	"gitlab.com/mfcekirdek/in-memory-store/internals/utils"
	"log"
	"os/exec"
	"time"
)

type StoreService interface {
	Get(key string) (map[string]string, error)
	Set(key string, value string) map[string]string
	Flush() map[string]string
}

type storeService struct {
	repository repository.StoreRepository
}

func NewStoreService(repo repository.StoreRepository, flushInterval int) StoreService {
	service := &storeService{repository: repo}
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

func (s *storeService) Set(key string, value string) map[string]string {
	s.repository.Set(key, value)
	return map[string]string{key: value}
}

func (s *storeService) Flush() map[string]string {
	return s.repository.Flush()
}

func (s *storeService) BackgroundTask(interval int, task func() []byte) {
	dateTicker := time.NewTicker(time.Duration(interval) * time.Minute)
	for x := range dateTicker.C {
		fmt.Println(x)
		response := task()
		fmt.Println(string(response))
	}
}

func saveToJSONFile() []byte {
	timestamp := time.Now().Unix()
	log.Println("Saving to file..", timestamp)
	out, err := exec.Command("date").Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

// todo Repository'i biraz incelt, service'e taşı logicleri
