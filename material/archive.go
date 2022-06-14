package material

import (
    "archive/zip"
    "io"
    "log"
    "os"
    "photo-materials/utils"
)

type archiveList struct {
    size  string
    files []string
}

func (al archiveList) newZipFile(matNo string) {
    f, err := os.Create("No" + matNo + "-" + al.size + ".zip")
    if err != nil {
        log.Fatalln(err)
    }
    defer utils.Close(f)

    zipWriter := zip.NewWriter(f)
    defer utils.Close(zipWriter)

    for _, file := range al.files {
        addFileToZip(zipWriter, file)
    }
}

func addFileToZip(zipWriter *zip.Writer, file string) {
    fileToZip, err := os.Open(file)
    if err != nil {
        log.Fatalln(err)
    }
    defer utils.Close(fileToZip)

    info, err := fileToZip.Stat()
    if err != nil {
        log.Fatalln(err)
    }

    header, err := zip.FileInfoHeader(info)
    if err != nil {
        log.Fatalln(err)
    }

    header.Name = file
    header.Method = zip.Deflate
    writer, err := zipWriter.CreateHeader(header)
    if err != nil {
        log.Fatalln(err)
    }

    if _, err = io.Copy(writer, fileToZip); err != nil {
        log.Fatalln(err)
    }
}

// CreateZipFiles selects target files and create 3-sizes of material zips.
func CreateZipFiles(matNo *string) {
    // Get listings by 3 size
    l, m, s := utils.FindFiles(".", false)

    var als []archiveList
    als = append(als, archiveList{size: "L", files: l})
    als = append(als, archiveList{size: "M", files: m})
    als = append(als, archiveList{size: "S", files: s})

    // Create archive files
    for _, al := range als {
        log.Printf("[ACT] CreateZipFiles [MSG] Archive list(%s): %s\n", al.size, al.files)
        al.newZipFile(*matNo)
        log.Printf("[ACT] CreateZipFiles [MSG] Archive completed(%s)\n", al.size)
    }
}
