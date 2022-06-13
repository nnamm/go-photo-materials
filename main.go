package main

import (
    "flag"
    "log"
    "os"
    "photo-materials/material"
    "photo-materials/publish"
)

func main() {
    matNo := flag.String("n", "000", "Material number")
    targetPath := flag.String("d", ".", "Working directory")
    remotePathNo := flag.Int("p", 0, "Select remote directory'(0=Test, 1-5=Product)")
    flag.Parse()

    if err := os.Chdir(*targetPath); err != nil {
        log.Fatalln(err)
    }

    material.CreateZipFiles(matNo)

    publish.RunPublish(matNo, remotePathNo)
}
