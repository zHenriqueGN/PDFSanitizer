package main

import (
	"log"

	"github.com/zHenriqueGN/PDFSanitizer/pkg/controller"
	"github.com/zHenriqueGN/PDFSanitizer/pkg/models"
)

func main() {
	logFile, logger, err := controller.CreateLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	folders, err := controller.GetRootFolders()
	if err != nil {
		logger.Fatal(err)
	}

	PDFsMap := make(map[string]models.PDFFile)

	for _, folder := range folders {
		err = controller.MapPDFFiles(folder, &PDFsMap)
		if err != nil {
			logger.Fatal(err)
		}
	}

	err = controller.CreateSanitizedPDFsFolder(PDFsMap, logger)
	if err != nil {
		logger.Fatal(err)
	}

}
