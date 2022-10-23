package compar

import (
	"encoding/json"
	"log"
	"os"

	"github.com/svetlyi/dircompar/dto"
)

var onlyInDump1 []dto.File
var onlyInDump2 []dto.File
var differentFiles []string

func Compare(dump1Path string, dump2Path string) {
	if dump1Path == "" || dump2Path == "" {
		log.Println("dump1 or dump2 path is empty")
		return
	}

	dump1 := getDump(dump1Path)
	dump2 := getDump(dump2Path)

	var dump1FilesMap map[string]dto.File = make(map[string]dto.File)
	var dump2FilesMap map[string]dto.File = make(map[string]dto.File)

	for _, dump1File := range dump1.Files {
		dump1FilesMap[dump1File.GetCleanUnixName()] = dump1File
	}
	for _, dump2File := range dump2.Files {
		dump2FilesMap[dump2File.GetCleanUnixName()] = dump2File
	}

	for dump1FileName, dump1File := range dump1FilesMap {
		if dump2File, existsInDump2 := dump2FilesMap[dump1FileName]; existsInDump2 {
			if dump1File.Hash != dump2File.Hash {
				differentFiles = append(differentFiles, dump1File.GetCleanUnixName())
			}
			delete(dump2FilesMap, dump1FileName)
		} else {
			onlyInDump1 = append(onlyInDump1, dump1File)
		}
	}

	for dump2FileName, dump2File := range dump2FilesMap {
		if _, existsInDump1 := dump1FilesMap[dump2FileName]; !existsInDump1 {
			onlyInDump2 = append(onlyInDump2, dump2File)
		}
	}

	log.Printf("=== only in %s: ===", dump1Path)
	for _, d1File := range onlyInDump1 {
		log.Println(d1File.GetCleanUnixName())
	}

	log.Printf("=== only in %s: ===", dump2Path)
	for _, d2File := range onlyInDump2 {
		log.Println(d2File.GetCleanUnixName())
	}

	log.Printf("=== different files: ===")
	for _, dFile := range differentFiles {
		log.Println(dFile)
	}
}

func getDump(f string) dto.Dump {
	fBytes, fBytesErr := os.ReadFile(f)
	if fBytesErr != nil {
		log.Fatalf("couldn't read file %s: %v", f, fBytesErr)
	}

	dump := dto.Dump{}

	if err := json.Unmarshal(fBytes, &dump); err != nil {
		log.Fatalf("couldn't unmarshal dump %s: %v", f, err)
	}

	return dump
}
