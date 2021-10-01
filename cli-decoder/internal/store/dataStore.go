package store

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type DataStore struct {
	config   *Config
	hashList map[string]string
}

func NewStore(config *Config) (*DataStore, func(), error) {
	hashFilePath := filepath.Join(config.FilePath, config.HashFileName)

	file, err := os.OpenFile(hashFilePath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s, %v", hashFilePath, err)
	}
	defer file.Close()

	store, err := getStore(file)
	store.config = config

	closeFunc := store.updateHashListFile

	if err != nil {
		return nil, nil, fmt.Errorf("problem creating hash list file, %v ", err)
	}

	return store, closeFunc, nil
}

func getStore(file *os.File) (*DataStore, error) {
	err := initHashListFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initialising hash store File, %v", err)
	}

	store := &DataStore{
		hashList: make(map[string]string),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hash := scanner.Text()
		store.hashList[hash] = hash
	}

	return store, nil
}

func initHashListFile(file *os.File) error {
	file.Seek(0, 0)
	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting File info from File %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte(""))
		file.Seek(0, 0)
	}

	return nil
}

func (s DataStore) setHash(hash string) {
	s.hashList[hash] = hash
	go s.updateHashListFile()
}

func (s DataStore) updateHashListFile() {
	hashFilePath := filepath.Join(s.config.FilePath, s.config.HashFileName)

	file, err := os.OpenFile(hashFilePath, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		fmt.Errorf("problem opening %s, %v\n", hashFilePath, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, data := range s.hashList {
		_, err := writer.WriteString(data + "\n")
		if err != nil {
			fmt.Errorf("problem save file hash\n")
		}
	}
	writer.Flush()
}

func (s *DataStore) SavaFile(name string, extension string, content []byte) error {
	if _, ok := s.hashList[name]; ok {
		return nil
	}
	path := filepath.Join(s.config.FilePath, name+extension)

	err := os.WriteFile(path, content, 0644)
	if err == nil {
		s.setHash(name)
	}

	return err
}
