package publish

import (
    "log"
    "os"
    "photo-materials/utils"
)

type urlTextConfig struct {
    remoteHome string
    remotePath []string
}

func loadUrlTextConfig() urlTextConfig {
    c := utils.LoadPathConfig("../config.ini")
    return urlTextConfig{
        c.RemoteHome,
        c.RemotePath,
    }
}

func (uc urlTextConfig) writeUrlText(matNo *string, pathNo *int, uls []UploadList) {
    f, err := os.Create("url.txt")
    if err != nil {
        log.Fatalln(err)
    }
    defer utils.Close(f)

    // Set url text
    // e.g. https://nnamm.com/ + path/to/upload-dir/ + no + 000 + /
    urlBase := uc.remoteHome + uc.remotePath[*pathNo] + "no" + *matNo + "/"

    // Repeat for the num of elements in the list
    for i := 0; i < len(uls[0].Files); i++ {
        for _, ul := range uls {
            filename := ul.Files[i]
            if _, err = f.WriteString(urlBase + filename + "\n"); err != nil {
                log.Fatalln(err)
            }
        }
        if _, err = f.WriteString("\n"); err != nil {
            log.Fatalln(err)
        }
    }
    log.Println("[ACT] writeUrlText [MSG] url.txt done")
}

// RunUrlText creates a text file with lists of URLs for note publication work
func RunUrlText(matNo *string, pathNo *int, uls []UploadList) {
    uc := loadUrlTextConfig()
    uc.writeUrlText(matNo, pathNo, uls)
}
