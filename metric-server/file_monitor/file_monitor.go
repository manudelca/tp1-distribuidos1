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
		mapStringToMutex: util.NewMapStringToMutex(),
	}
}

func (f *FileMonitor) ReadLine(name string, i int) (string, error) {
	return "", errors.New("ReadLine not implemented yet")
}

func (f *FileMonitor) WriteLineOnFile(line string, fileName string) error {
	fileMutex, ok := f.mapStringToMutex.Get(fileName)
	if !ok {
		mut := sync.Mutex{}
		fileMutex = f.mapStringToMutex.Add(fileName, &mut)
	}
	fileMutex.Lock()
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	for i := 0; i < len(line); {
		written, err := file.WriteString(line[i:])
		if err != nil {
			fileMutex.Unlock()
			return err
		}
		i = i + written
	}
	file.Close()
	fileMutex.Unlock()
	return nil
}
