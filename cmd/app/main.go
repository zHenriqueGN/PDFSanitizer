package main

import (
	"log"
	"os"

	"github.com/zHenriqueGN/PDFSanitizer/pkg/controller"
	"github.com/zHenriqueGN/PDFSanitizer/pkg/models"
)

func main() {
	logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	folders, err := controller.GetRootFolders()
	if err != nil {
		log.Fatal(err)
	}

	PDFsMap := make(map[string]models.PDFFile)

	for _, folder := range folders {
		err = controller.MapPDFFiles(folder, &PDFsMap)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = controller.CreateSanitizedPDFsFolder(PDFsMap)
	if err != nil {
		log.Fatal(err)
	}

}
