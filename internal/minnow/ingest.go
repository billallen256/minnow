package minnow

import (
	"log"
	"os"
	"time"
)

type IngestDirInfo struct {
	IngestPath         Path
	MinAge             time.Duration
	ProcessedBy        []ProcessorId
	RemoveOnceIngested bool
}

type DirectoryIngester struct {
	workPath      Path
	ingestDirChan chan IngestDirInfo
	dispatchChan  chan DispatchInfo
	logger        *log.Logger
}

func NewDirectoryIngester(workPath Path, ingestDirChan chan IngestDirInfo, dispatchChan chan DispatchInfo) *DirectoryIngester {
	logger := log.New(os.Stdout, "DirectoryIngester: ", 0)
	return &DirectoryIngester{workPath, ingestDirChan, dispatchChan, logger}
}

func moveToRandomPath(workPath, metadataPath, dataPath Path) (Path, Path, error) {
	randomPath, err := makeRandomPath(workPath, "dispatch")

	if err != nil {
		return metadataPath, dataPath, err
	}

	newMetadataPath := randomPath.JoinPath(Path(metadataPath.Name()))
	err = metadataPath.Rename(newMetadataPath)

	if err != nil {
		return metadataPath, dataPath, err
	}

	newDataPath := randomPath.JoinPath(Path(dataPath.Name()))
	err = dataPath.Rename(newDataPath)

	if err != nil {
		return metadataPath, dataPath, err
	}

	return newMetadataPath, newDataPath, nil
}

func (ingester *DirectoryIngester) Run() {
	for ingestDirInfo := range ingester.ingestDirChan {
		metadataPaths, _ := ingestDirInfo.IngestPath.Glob("*" + PropertiesExtension)

		for _, metadataPath := range metadataPaths {
			dataPath := metadataPath.WithSuffix("") // lop off the extension

			if !dataPath.Exists() {
				ingester.logger.Printf("%s does not have corresponding data file", metadataPath)
				continue
			}

			now := time.Now()
			metadataAge, err := metadataPath.Age(now)

			if err != nil {
				ingester.logger.Print(err.Error())
				continue
			}

			dataAge, err := dataPath.Age(now)

			if err != nil {
				ingester.logger.Print(err.Error())
				continue
			}

			if metadataAge > ingestDirInfo.MinAge && dataAge > ingestDirInfo.MinAge {
				// Move things to a random path in case we're ingesting from the
				// main ingest directory.  Files that have already been processed
				// are already in a random directory, but we're moving them anyway
				// just to be consistent.
				metadataPath, dataPath, err := moveToRandomPath(ingester.workPath, metadataPath, dataPath)

				if err != nil {
					ingester.logger.Print(err.Error())
					continue
				}

				dispatchInfo := DispatchInfo{metadataPath, dataPath, ingestDirInfo.ProcessedBy}
				ingester.dispatchChan <- dispatchInfo

				if ingestDirInfo.RemoveOnceIngested {
					err := ingestDirInfo.IngestPath.RmdirRecursive()

					if err != nil {
						ingester.logger.Printf("Could not remove ingest dir: %s", err.Error())
					}
				}
			}
		}
	}
}
