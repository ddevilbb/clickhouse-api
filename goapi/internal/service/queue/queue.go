package queue

import (
	"clickhouse-api/internal/helper"
	"clickhouse-api/internal/repository/test_data"
	"clickhouse-api/internal/service/config"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	dateFormat = "20060102"
	queueName  = "queue"
)

var ErrNoFiles = errors.New("no files")

type FileQueue struct {
	Path              string
	Operation         string
	Table             string
	Extension         string
	FileNum           int
	CheckFileInterval int
	LineCount         int
	Date              string
	LockedFiles       map[string]bool
	mu                sync.Mutex
	repository        *test_data.Repository
}

func New(conf *config.Config, operation string, table string, repository *test_data.Repository) *FileQueue {
	queue := new(FileQueue)
	queue.Path = fmt.Sprintf("%s/%s", conf.QueueConfig.Path, operation)
	queue.Operation = operation
	queue.Table = table
	queue.Extension = conf.QueueConfig.Extension
	queue.FileNum = 1
	queue.CheckFileInterval = conf.QueueConfig.CheckFileInterval
	queue.Date = time.Now().Format(dateFormat)
	queue.repository = repository
	return queue
}

func (queue *FileQueue) checkDir(create bool) error {
	_, err := os.Stat(queue.Path)
	if os.IsNotExist(err) && create {
		return os.MkdirAll(queue.Path, 0777)
	}
	return err
}

func (queue *FileQueue) prepareFilename() string {
	return fmt.Sprintf(
		"%s-%s-%s-%s.que",
		queue.Operation,
		queue.Table,
		queue.Date,
		strconv.Itoa(queue.FileNum),
	)
}

func (queue *FileQueue) Add(content string) error {
	queue.mu.Lock()
	defer queue.mu.Unlock()
	err := queue.checkDir(true)
	if err != nil {
		return err
	}
	queue.LineCount++
	if queue.LineCount > 1000 {
		queue.FileNum++
		queue.LineCount = 1
	}
	currentDate := time.Now().Format(dateFormat)
	if currentDate != queue.Date {
		queue.Date = currentDate
		queue.FileNum = 1
	}
	err = helper.AppendToFile(path.Join(queue.Path, queue.prepareFilename()), content)
	return err
}

func (queue *FileQueue) GetFile() (string, error) {
	const getFileError = queueName + "get file error: "
	err := queue.checkDir(false)
	if os.IsNotExist(err) {
		return "", ErrNoFiles
	}
	if err != nil {
		return "", fmt.Errorf(getFileError+"%+v", err)
	}
	files, err := os.ReadDir(queue.Path)
	if err != nil {
		return "", fmt.Errorf(getFileError+"%+v", err)
	}
	queueFiles := make([]string, 0)
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".que" {
			queueFiles = append(queueFiles, f.Name())
		}
	}
	sort.Strings(queueFiles)
	for _, f := range queueFiles {
		found, _ := queue.LockedFiles[f]
		if !found {
			queueFiles = nil
			return f, nil
		}
	}
	queueFiles = nil
	return "", ErrNoFiles
}

func (queue *FileQueue) GetFileData(filename string) (string, error) {
	s, err := os.ReadFile(path.Join(queue.Path, filename))
	if err != nil {
		return "", fmt.Errorf(queueName+" get file data error: %+v", err)
	}
	return string(s), nil
}

func (queue *FileQueue) RemoveFile(filename string) (err error) {
	err = os.Remove(path.Join(queue.Path, filename))
	if err != nil {
		return fmt.Errorf(queueName+" remove file error: %+v", err)
	}
	return nil
}

func (queue *FileQueue) ProcessNextFile() error {
	const processNextFileError = queueName + " process next file error: "
	queue.mu.Lock()
	defer queue.mu.Unlock()
	file, err := queue.GetFile()
	if err != nil && !errors.Is(err, ErrNoFiles) {
		return fmt.Errorf(processNextFileError+"%+v", err)
	}
	if file == "" {
		return nil
	}
	fileData, err := queue.GetFileData(file)
	if err != nil {
		return fmt.Errorf(processNextFileError+"%+v", err)
	}
	if fileData != "" {
		lineList := strings.Split(fileData, "\n")
		err = queue.processOperation(lineList)
		lineList = nil
		if err != nil {
			return fmt.Errorf(processNextFileError+"%+v", err)
		}
	}
	err = queue.RemoveFile(file)
	if err != nil {
		queue.LockedFiles[file] = true
		return fmt.Errorf(processNextFileError+"%+v", err)
	}
	return err
}

func (queue *FileQueue) Listen() {
	queue.LockedFiles = make(map[string]bool)
	go func(queue *FileQueue) {
		ticker := time.NewTicker(time.Second * time.Duration(queue.CheckFileInterval))
		for range ticker.C {
			err := queue.ProcessNextFile()
			if err != nil {
				if !errors.Is(err, ErrNoFiles) {
					log.Error(queueName+" operation "+queue.Operation+" listen error: %+v\n", err)
				}
			}
		}
	}(queue)
}

func (queue *FileQueue) processOperation(list []string) error {
	switch queue.Operation {
	case "update":
		return queue.repository.BulkUpdate(list)
	case "delete":
		return queue.repository.BulkDelete(list)
	default:
		return queue.repository.BulkInsert(list)
	}
}
