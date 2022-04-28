package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"time"
)

func main() {
	layout := "2006-01-02 15:04:05"
	date, err := time.Parse(layout, "2022-04-28 20:15:00")
	if err != nil {
		fmt.Println("Error parsing date")
		return
	}
	writeMetric("metric_1", date.Unix(), 3)
}

func writeMetric(metricId string, timestamp int64, value float32) {

	bytes := make([]byte, 8+4)
	binary.BigEndian.PutUint64(bytes[:8], uint64(timestamp))
	binary.BigEndian.PutUint32(bytes[8:], math.Float32bits(value))

	timeParsed := time.Unix(timestamp, 0)
	year, month, day := timeParsed.Date()
	hours, minutes, _ := timeParsed.Clock()
	fileToWrite := fmt.Sprintf("%s_%d%02d%02d_%02d%02d", metricId, year, month, day, hours, minutes)
	file, err := os.OpenFile(fileToWrite, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Could not write file")
		return
	}
	for i := 0; i < len(bytes); {
		written, err := file.Write(bytes[i:])
		if err != nil {
			fmt.Println("Could not write file")
			return
		}
		i = i + written
	}
	file.Close()
}
