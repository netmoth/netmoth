package sensor

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/netmoth/netmoth/internal/connection"
)

type logSave struct {
	fileName string
	mutex    sync.Mutex
}

type logData struct {
	SensorMetadata *Metadata
	Connections    []connection.Connection
}

func newSaver(logName string, meta *Metadata) (*logSave, error) {
	f, err := os.Create(logName)
	if err != nil {
		return nil, err
	}

	logFile := &logData{
		SensorMetadata: meta,
	}

	initJSON, err := json.MarshalIndent(logFile, "", "  ")
	if err != nil {
		return nil, err
	}

	if _, err := f.Write(initJSON); err != nil {
		return nil, err
	}

	return &logSave{
		fileName: logName,
	}, nil
}

func (l *logSave) save(c connection.Connection) {
	l.mutex.Lock()
	contents, err := os.ReadFile(l.fileName)
	if err != nil {
		log.Println(err)
	}

	data := new(logData)
	if err := json.Unmarshal(contents, data); err != nil {
		log.Println(err)
	}

	data.Connections = append(data.Connections, c)

	newContents, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile(l.fileName, newContents, 0644)
	if err != nil {
		log.Println(err)
	}
	l.mutex.Unlock()
}
