package file_monitor

import (
	"errors"
	"os"
	"sync"

	"github.com/manudelca/tp1-distribuidos1/metric-server/util"
)

type FileMonitor struct {
	mapStringToMutex *util.MapStringToMutex
}

func NewFileMonitor() *FileMonitor {
	return &FileMonitor{
		mapStringToMutex: &util.MapStringToMutex{},
	}
}

func (f *FileMonitor) ReadFromFile(name string) (string, error) {
	return "", nil
}

func (f *FileMonitor) WriteLineOnFile(line string, fileName string) error {
	fileMutex, ok := f.mapStringToMutex.Get(fileName)
	if !ok {
		mut := sync.Mutex{}
		fileMutex = f.mapStringToMutex.Add(fileName, &mut)
	}
	fileMutex.Lock()
	file, err := os.Open(fileName)
	if errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(fileName)
		if err != nil {
			fileMutex.Unlock()
			return err
		}
	}
	for i := 0; i < len(line); {
		written, err := file.WriteString(line[i:])
		if err != nil {
			fileMutex.Unlock()
			return err
		}
		i = i + written
	}

	fileMutex.Unlock()
	return nil
}
