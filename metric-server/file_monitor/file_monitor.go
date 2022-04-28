package file_monitor

import (
	"errors"
	"io"
	"os"
	"sync"

	"github.com/manudelca/tp1-distribuidos1/metric-server/util"
)

const MetricSaveLen = 8 + 4

type FileMonitor struct {
	mapStringToMutex *util.MapStringToMutex
}

func NewFileMonitor() *FileMonitor {
	return &FileMonitor{
		mapStringToMutex: util.NewMapStringToMutex(),
	}
}

func (f *FileMonitor) FileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		return !errors.Is(err, os.ErrNotExist)
	}
	return true
}

func (f *FileMonitor) ReadMetric(fileName string, metricIndex int) ([]byte, error) {
	fileMutex, ok := f.mapStringToMutex.Get(fileName)
	if !ok {
		return nil, errors.New("Non exisiting file")
	}
	fileMutex.Lock()
	file, err := os.Open(fileName)
	if err != nil {
		fileMutex.Unlock()
		return nil, err
	}
	result := make([]byte, MetricSaveLen)
	for i := 0; i < MetricSaveLen; {
		read, err := file.ReadAt(result, int64(metricIndex)*int64(MetricSaveLen))
		if errors.Is(err, io.EOF) {
			file.Close()
			fileMutex.Unlock()
			return make([]byte, 0), nil
		} else if err != nil {
			file.Close()
			fileMutex.Unlock()
			return nil, err
		}
		i = i + read
	}
	file.Close()
	fileMutex.Unlock()
	return result, nil
}

func (f *FileMonitor) WriteMetricOnFile(bytes []byte, fileName string) error {
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
	for i := 0; i < len(bytes); {
		written, err := file.Write(bytes[i:])
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
