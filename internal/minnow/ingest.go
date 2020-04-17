package minnow

import (
	"fmt"
	"os"
)

type IngestInfo struct {
	ProcessedBy  []ProcessorId
	MetadataPath Path
	DataPath     Path
}

type DirectoryIngester struct {
	IngestPath Path
	MinAge     time.Duration
	IngestChan chan IngestInfo
	Logger     *log.Logger
}

func NewDirectoryIngester(path Path, minAge time.Duration, ingestChan chan<- IngestInfo) (*DirectoryIngester, error) {
	if !path.Exists() {
		return nil, fmt.Errorf("%s does not exist", path)
	}

	logger := log.New(os.Stdout, "DirectoryIngester: ", 0)
	return &DirectoryIngester{path, minAge, ingestChan, logger}, nil
}

func (ingester *DirectoryIngester) Once(processedBy []ProcessorId) error {
	pairsToIngest := make([]IngestInfo, 0)
	metadataPaths, _ := ingester.IngestPath.Glob("*" + PropertiesExtension)

	for _, metadataPath := range metadataPaths {
		dataPath := metadataPath.WithSuffix("") // lop off the extension

		if !dataPath.Exists() {
			ingester.Logger.Print("%s does not have corresponding data file", metadataPath)
			continue
		}

		now := time.Now()

		if metadataPath.Age(now) > ingester.MinAge && dataPath.Age(now) > ingester.MinAge {
			ingestInfo := IngestInfo{processedBy, metadataPath, datapath}
			ingester.IngestChan <- ingestInfo
		}
	}
}

func (ingester *DirectoryIngester) Periodic(interval time.Duration) {
	for range time.Ticker(interval) {
		err := ingester.Once(make([]ProcessorId, 0))

		if err != nil {
			ingester.Logger.Print(err.Error())
		}
	}
}
