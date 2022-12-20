package dump

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/svetlyi/dircompar/dto"
)

type FilesMeta struct {
	size  uint64
	count uint64
}

func DumpRun(searchIn string, saveDumpTo string, removePrefix string, skipWithErrors bool) {
	if searchIn == "" || saveDumpTo == "" {
		log.Println("searchIn or saveDumpTo is empty")
		return
	}
	log.Printf("dumping %s", searchIn)

	filesMeta := getFilesMeta(searchIn, skipWithErrors)
	filesCounter := 0
	var filesSizeCounter uint64 = 0

	d := dto.Dump{Path: searchIn}

	err := filepath.Walk(searchIn,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if skipWithErrors {
					log.Printf("skipping %s: %v", path, err)
					return nil
				}
				return fmt.Errorf("couldn't open folder: %v", err)
			}
			if info.IsDir() {
				return nil
			}
			pathWithoutPrefix := strings.TrimPrefix(path, removePrefix)

			md5, err := getMd5(path)
			if err != nil {
				if skipWithErrors {
					log.Printf("skipping %s: %v", path, err)
					return nil
				} else {
					return fmt.Errorf("couldn't calculate md5 for %s: %v", path, err)
				}
			}
			d.Files = append(d.Files, dto.File{
				Name: pathWithoutPrefix,
				Hash: md5,
				Size: info.Size(),
			})
			filesCounter++
			filesSizeCounter += uint64(info.Size())

			log.Printf(
				"[%d/%d][%s/%s (%f %%)]: %s",
				filesCounter,
				filesMeta.count,
				byteCountSI(filesSizeCounter),
				byteCountSI(filesMeta.size),
				(float64(filesSizeCounter)/float64(filesMeta.size))*100,
				pathWithoutPrefix,
			)
			return nil
		})

	if err != nil {
		log.Fatalf("error traversing %s: %v", searchIn, err)
	}

	if dumpJson, dumpJsonErr := json.Marshal(d); dumpJsonErr != nil {
		log.Fatalf("couldn't marshal dump: %v", dumpJsonErr)
	} else if writeErr := os.WriteFile(saveDumpTo, dumpJson, 0600); writeErr != nil {
		log.Fatalf("couldn't save dump to %s: %v", saveDumpTo, writeErr)
	}

	log.Printf("saved dump to %s", saveDumpTo)
}

func getFilesMeta(searchIn string, skipWithErrors bool) FilesMeta {
	log.Println("calculating files meta...")

	meta := FilesMeta{}

	err := filepath.Walk(searchIn,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if skipWithErrors {
					log.Printf("calculating meta: skipping %s: %v", path, err)
					return nil
				}
				return err
			}
			if info.IsDir() {
				return nil
			}
			meta.count++
			meta.size += uint64(info.Size())
			return nil
		})

	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("found %d files with overall size %s", meta.count, byteCountSI(meta.size))

	return meta
}

func getMd5(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("couldn't open file %s: %v", file, err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("couldn't read file %s: %v", file, err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func byteCountSI(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
