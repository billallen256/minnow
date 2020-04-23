package minnow

import (
	"log"
	"os"
	"time"
)

type IngestDirInfo struct {
	IngestPath  Path
	MinAge      time.Duration
	ProcessedBy []ProcessorId
}

type DirectoryIngester struct {
	ingestDirChan chan IngestDirInfo
	dispatchChan  chan DispatchInfo
	Logger        *log.Logger
}

func NewDirectoryIngester(ingestDirChan chan IngestDirInfo, dispatchChan chan DispatchInfo) *DirectoryIngester {
	logger := log.New(os.Stdout, "DirectoryIngester: ", 0)
	return &DirectoryIngester{ingestDirChan, dispatchChan, logger}
}

func (ingester *DirectoryIngester) Run() {
	for ingestDirInfo := range ingester.ingestDirChan {
		metadataPaths, _ := ingestDirInfo.IngestPath.Glob("*" + PropertiesExtension)

		for _, metadataPath := range metadataPaths {
			dataPath := metadataPath.WithSuffix("") // lop off the extension

			if !dataPath.Exists() {
				ingester.Logger.Printf("%s does not have corresponding data file", metadataPath)
				continue
			}

			now := time.Now()
			metadataAge, err := metadataPath.Age(now)

			if err != nil {
				ingester.Logger.Print(err.Error())
				continue
			}

			dataAge, err := dataPath.Age(now)

			if err != nil {
				ingester.Logger.Print(err.Error())
				continue
			}

			if metadataAge > ingestDirInfo.MinAge && dataAge > ingestDirInfo.MinAge {
				dispatchInfo := DispatchInfo{metadataPath, dataPath, ingestDirInfo.ProcessedBy}
				ingester.dispatchChan <- dispatchInfo
			}
		}
	}
}
