package minnow

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func nanoTimestamp(now time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d:%09d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond())
}

func makeRandomPath(baseDir Path, purpose string) (Path, error) {
	purpose = strings.ReplaceAll(purpose, " ", "")
	name := fmt.Sprintf("%s-%s", purpose, nanoTimestamp(time.Now()))
	path := baseDir.JoinPath(Path(name))

	// If the path somehow already exists, try again
	for path.Exists() {
		name = fmt.Sprintf("%s-%s", purpose, nanoTimestamp(time.Now()))
		path = baseDir.JoinPath(Path(name))
	}

	err := path.Mkdir()

	if err != nil {
		return baseDir, err
	}

	return path, nil
}

func CopyFile(source, destination Path) error {
	if source.IsDir() {
		return fmt.Errorf("Cannot use Copy for a directory: %s", source)
	}

	if destination.IsDir() {
		destination = destination.JoinPath(Path(source.Name()))
	}

	if source == destination {
		return fmt.Errorf("Copy source and destination are the same.  Nothing to do.")
	}

	from, err := os.Open(string(source))

	if err != nil {
		return err
	}

	defer from.Close()

	perms, err := source.Permissions()

	if err != nil {
		return err
	}

	// Copy source permissions, making sure that the destination is
	// at least readable and writable.
	perms = perms | 0600
	to, err := os.OpenFile(string(destination), os.O_RDWR|os.O_CREATE, perms)

	if err != nil {
		return err
	}

	defer to.Close()

	_, err = io.Copy(to, from)

	if err != nil {
		return err
	}

	return nil
}
