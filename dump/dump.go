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

func DumpRun(searchIn string, saveDumpTo string, removePrefix string) {
	if searchIn == "" || saveDumpTo == "" {
		log.Println("searchIn or saveDumpTo is empty")
		return
	}
	log.Printf("dumping %s", searchIn)

	filesMeta := getFilesMeta(searchIn)
	filesCounter := 0
	var filesSizeCounter uint64 = 0

	d := dto.Dump{Path: searchIn}

	err := filepath.Walk(searchIn,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			pathWithoutPrefix := strings.TrimPrefix(path, removePrefix)

			d.Files = append(d.Files, dto.File{
				Name: pathWithoutPrefix,
				Hash: getMd5(path),
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
		log.Fatalln(err)
	}

	if dumpJson, dumpJsonErr := json.Marshal(d); dumpJsonErr != nil {
		log.Fatalf("couldn't marshal dump: %v", dumpJsonErr)
	} else if writeErr := os.WriteFile(saveDumpTo, dumpJson, 0600); writeErr != nil {
		log.Fatalf("couldn't save dump to %s: %v", saveDumpTo, writeErr)
	}

	log.Printf("saved dump to %s", saveDumpTo)
}

func getFilesMeta(searchIn string) FilesMeta {
	log.Println("calculating files meta...")

	meta := FilesMeta{}

	err := filepath.Walk(searchIn,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
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

func getMd5(file string) string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
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
