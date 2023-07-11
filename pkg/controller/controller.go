package controller

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/zHenriqueGN/pdfSanitizer/pkg/models"
)

// GetRootFolders gets all the folders in the root
// where the code is executing
func GetRootFolders() (folders []string, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}

	dirEntrys, err := os.ReadDir(currentDir)
	if err != nil {
		return
	}

	for _, folder := range dirEntrys {
		if folder.IsDir() {
			if err != nil {
				return
			}
			folders = append(folders, folder.Name())
		}
	}
	return
}

// MapPDFFiles search for files in a folder and map them. Files that has a name
// that is already in the map will be ignored
func MapPDFFiles(folderPath string, PDFs *map[string]models.PDFFile) (err error) {
	dirEntrys, err := os.ReadDir(folderPath)
	if err != nil {
		return
	}

	for _, entry := range dirEntrys {
		var PDF models.PDFFile
		PDF.Name = entry.Name()
		PDF.Path = path.Join(folderPath, PDF.Name)
		if _, keyExists := (*PDFs)[PDF.Name]; !keyExists {
			(*PDFs)[PDF.Name] = PDF
		} else {
			fmt.Println("Key exists", PDF.Name)
		}
	}

	return
}

// CreateSanitizedPDFsFolder move all the PDFs of a given map and move them
// to a folder where all the pdfs are unique
func CreateSanitizedPDFsFolder(PDFsMap map[string]models.PDFFile) (err error) {
	dstFolder := "sanitized_pdfs"
	err = os.MkdirAll(dstFolder, 0777)
	if err != nil {
		return
	}

	for PDFName, PDF := range PDFsMap {
		sourcePDF, err := os.Open(PDF.Path)
		if err != nil {
			return err
		}
		defer sourcePDF.Close()

		dstPDF, err := os.Create(path.Join(dstFolder, PDFName))
		if err != nil {
			return err
		}
		defer dstPDF.Close()

		_, err = io.Copy(dstPDF, sourcePDF)
		if err != nil {
			return err
		}
	}

	return
}
