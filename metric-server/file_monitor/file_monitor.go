package file_monitor

import "sync"

type FileMonitor struct {
	filesMutex map[string]sync.Mutex
}

func (f *FileMonitor) ReadFromFile(name string) (string, error) {
	return "", nil
}

func (f *FileMonitor) WriteOnFile(name string) error {
	return nil
}
