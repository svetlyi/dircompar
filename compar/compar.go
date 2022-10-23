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

dump1Loop:
	for _, dump1File := range dump1.Files {
		for _, dump2File := range dump2.Files {
			if dump1File.Name == dump2File.Name {
				if dump1File.Hash != dump2File.Hash {
					differentFiles = append(differentFiles, dump1File.Name)
				}
				continue dump1Loop // found dump1 in dump2
			}
		}
		// didn't find dump1 in dump2
		onlyInDump1 = append(onlyInDump1, dump1File)
	}

dump2Loop:
	for _, dump2File := range dump2.Files {
		for _, dump1File := range dump1.Files {
			if dump2File.Name == dump1File.Name {
				continue dump2Loop // found dump2 in dump1
			}
		}
		// didn't find dump1 in dump2
		onlyInDump2 = append(onlyInDump1, dump2File)
	}

	log.Printf("=== only in %s: ===", dump1Path)
	for _, d1File := range onlyInDump1 {
		log.Println(d1File.Name)
	}

	log.Printf("=== only in %s: ===", dump2Path)
	for _, d2File := range onlyInDump2 {
		log.Println(d2File.Name)
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
