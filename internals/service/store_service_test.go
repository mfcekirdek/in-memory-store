//+build unit

package service

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"gitlab.com/mfcekirdek/in-memory-store/internals/repository"
	"gitlab.com/mfcekirdek/in-memory-store/mocks"
)

func TestNewStoreService(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	mockCategoryRepository := mocks.NewMockStoreRepository(mockController)
	mockCategoryRepository.
		EXPECT().
		LoadStore(gomock.Any()).
		Return().
		Times(1)

	type args struct {
		repo     repository.StoreRepository
		interval int
		path     string
	}
	tests := []struct {
		name string
		args args
		want StoreService
	}{
		{"New storeService instance", args{
			repo:     mockCategoryRepository,
			interval: 10,
			path:     "../storage",
		}, &storeService{
			repository:     mockCategoryRepository,
			storageDirPath: "../storage",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStoreService(tt.args.repo, tt.args.interval, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStoreService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterValidDataFiles(t *testing.T) {
	timestamp := time.Now().Unix()
	suffix := "-data.json"
	validFileName := fmt.Sprintf("%d%s", timestamp, suffix)

	path := "/tmp/testDir"
	_ = os.MkdirAll(path, os.ModePerm)
	defer os.RemoveAll(path)

	createFile(path, validFileName)
	wanted, _ := ioutil.ReadDir(path)

	createFile(path, "a.txt")
	createFile(path, "b-data.json")
	files, _ := ioutil.ReadDir(path)
	type args struct {
		files []os.FileInfo
	}
	tests := []struct {
		name string
		args args
		want []os.FileInfo
	}{
		{"Filter valid 'TIMESTAMP-data.json' files.", args{files: files}, wanted},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterValidDataFiles(tt.args.files); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterValidDataFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createFile(path, validFileName string) {
	filePath := filepath.Join(path, validFileName)
	_, _ = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func Test_findJSONFilePath(t *testing.T) {
	now := time.Now()
	yesterdaysTimestamp := now.AddDate(0, 0, -1).Unix()
	todaysTimestamp := now.Unix()
	tomorrowsTimestamp := now.AddDate(0, 0, 1).Unix()
	suffix := "-data.json"

	yesterdaysFile := fmt.Sprintf("%d%s", yesterdaysTimestamp, suffix)
	todaysFile := fmt.Sprintf("%d%s", todaysTimestamp, suffix)
	tomorrowsFile := fmt.Sprintf("%d%s", tomorrowsTimestamp, suffix)

	path1 := "/tmp/testDir1"
	_ = os.MkdirAll(path1, os.ModePerm)
	defer os.RemoveAll(path1)
	files1, _ := ioutil.ReadDir(path1)
	want1 := ""

	path2 := "/tmp/testDir2"
	_ = os.MkdirAll(path2, os.ModePerm)
	defer os.RemoveAll(path2)
	createFile(path2, yesterdaysFile)
	createFile(path2, todaysFile)
	createFile(path2, tomorrowsFile)
	files2, _ := ioutil.ReadDir(path2)
	want2 := filepath.Join(path2, tomorrowsFile)

	type args struct {
		dir   string
		files []fs.FileInfo
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"File does not exist in given file list", args{
			dir:   path1,
			files: files1,
		}, want1},
		{"File exists in given file list", args{
			dir:   path2,
			files: files2,
		}, want2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findJSONFilePath(tt.args.dir, tt.args.files); got != tt.want {
				t.Errorf("findJSONFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTimestampFromFilename(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"File format is valid", args{filename: "1634775852-data.json"}, "1634775852"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTimestampFromFilename(tt.args.filename); got != tt.want {
				t.Errorf("getTimestampFromFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_saveToMap(t *testing.T) {
	path := "/tmp/testDir"
	_ = os.MkdirAll(path, os.ModePerm)
	defer os.RemoveAll(path)

	filePath := filepath.Join(path, "tmp.json")
	store := map[string]string{"foo": "bar"}
	file, _ := json.MarshalIndent(store, "", " ")
	var perm fs.FileMode = 0600
	_ = ioutil.WriteFile(filePath, file, perm)

	filePath2 := filepath.Join(path, "tmp2.json")
	store2 := map[string]int{"foo": 123}
	file2, _ := json.MarshalIndent(store2, "", " ")
	_ = ioutil.WriteFile(filePath2, file2, perm)

	type args struct {
		jsonFilePath string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"File does not exist", args{jsonFilePath: "/wrong/path"}, nil},
		{"File exists with valid content", args{jsonFilePath: filePath}, store},
		{"File exists not valid content", args{jsonFilePath: filePath2}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := saveToMap(tt.args.jsonFilePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("saveToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_saveToJSONFile(t *testing.T) {
	path := "/tmp/tmp.json"
	defer os.RemoveAll(path)

	type args struct {
		filePath string
		store    map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"invalid path",
			args{
				filePath: "/invalid/path",
				store:    map[string]string{"foo": "bar"},
			},
			true},
		{"valid path",
			args{
				filePath: path,
				store:    map[string]string{"foo": "bar"},
			},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := saveToJSONFile(tt.args.filePath, tt.args.store); (err != nil) != tt.wantErr {
				t.Errorf("saveToJSONFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_backgroundTask(t *testing.T) {
	path := "/tmp/testDir"
	_ = os.MkdirAll(path, os.ModePerm)
	defer os.RemoveAll(path)

	type args struct {
		interval       time.Duration
		storageDirPath string
		store          map[string]string
		task           func(filepath string, store map[string]string) error
	}
	tests := []struct {
		name          string
		args          args
		wantFileCount int
	}{
		{"Start a goroutine to save the store to file at intervals", args{
			interval:       time.Second * 1,
			storageDirPath: "/invalid/path",
			store:          map[string]string{"foo": "bar"},
			task:           saveToJSONFile,
		}, 0},
		{"Start a goroutine to save the store to file at intervals", args{
			interval:       time.Second * 1,
			storageDirPath: path,
			store:          map[string]string{"foo": "bar"},
			task:           saveToJSONFile,
		}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go backgroundTask(tt.args.interval, tt.args.storageDirPath, tt.args.store, tt.args.task)
			time.Sleep(time.Second*time.Duration(tt.wantFileCount) + time.Millisecond*300)
			files, _ := ioutil.ReadDir(path)
			if got := len(files); got != tt.wantFileCount {
				t.Errorf("backgroundTask() = %v, want %v", got, tt.wantFileCount)
			}
		})
	}
}

func Test_storeService_Flush(t *testing.T) {
	mockController := gomock.NewController(t)
	mockRepository := mocks.NewMockStoreRepository(mockController)
	mockRepository.
		EXPECT().
		Flush().
		Return(map[string]string{}).
		Times(1)

	type fields struct {
		repository     repository.StoreRepository
		storageDirPath string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		{"returns what repository returns", fields{
			repository:     mockRepository,
			storageDirPath: "/tmp",
		}, map[string]string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storeService{
				repository:     tt.fields.repository,
				storageDirPath: tt.fields.storageDirPath,
			}
			if got := s.Flush(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flush() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storeService_Set(t *testing.T) {
	mockController := gomock.NewController(t)
	mockRepository := mocks.NewMockStoreRepository(mockController)
	mockRepository.
		EXPECT().
		Set(gomock.Any(), gomock.Any()).
		Return(true).
		Times(1)

	type fields struct {
		repository     repository.StoreRepository
		storageDirPath string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
		want1  bool
	}{
		{"returns what repository returns", fields{
			repository:     mockRepository,
			storageDirPath: "/tmp/testDir",
		}, args{
			key:   "foo",
			value: "bar",
		}, map[string]string{"foo": "bar"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storeService{
				repository:     tt.fields.repository,
				storageDirPath: tt.fields.storageDirPath,
			}
			got, got1 := s.Set(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Set() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_storeService_Get(t *testing.T) {
	mockController := gomock.NewController(t)
	mockRepositoryExists := mocks.NewMockStoreRepository(mockController)
	mockRepositoryExists.
		EXPECT().
		Get(gomock.Any()).
		Return("bar").
		Times(1)

	mockRepositoryNotExists := mocks.NewMockStoreRepository(mockController)
	mockRepositoryNotExists.
		EXPECT().
		Get(gomock.Any()).
		Return("").
		Times(1)

	type fields struct {
		repository     repository.StoreRepository
		storageDirPath string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		{"returns what repository returns - exists", fields{
			repository:     mockRepositoryExists,
			storageDirPath: "/tmp/testDir",
		}, args{key: "foo"}, map[string]string{"foo": "bar"}, false},
		{"returns what repository returns - not exists", fields{
			repository:     mockRepositoryNotExists,
			storageDirPath: "/tmp/testDir",
		}, args{key: "foo"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storeService{
				repository:     tt.fields.repository,
				storageDirPath: tt.fields.storageDirPath,
			}
			got, err := s.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readStoreDataFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{"Could not create directory", args{path: "/path"}, map[string]string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readStoreDataFromFile(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readStoreDataFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
