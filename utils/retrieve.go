package utils

import (
    "io/fs"
    "log"
    "path/filepath"
    "strings"
)

func FindFiles(root string, all bool) (sizeL []string, sizeM []string, sizeS []string) {
    var l, m, s []string
    err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
        if err != nil {
            log.Fatal(err)
        }

        if info.IsDir() {
            return nil
        }

        // Image files
        switch {
        case strings.Contains(info.Name(), "-L.jpg"):
            l = append(l, path)
        case strings.Contains(info.Name(), "-M.jpg"):
            m = append(m, path)
        case strings.Contains(info.Name(), "-S.jpg"):
            s = append(s, path)
        }

        // and archive files
        if all {
            switch {
            case strings.Contains(info.Name(), "-L.zip"):
                l = append(l, path)
            case strings.Contains(info.Name(), "-M.zip"):
                m = append(m, path)
            case strings.Contains(info.Name(), "-S.zip"):
                s = append(s, path)
            }
        }

        return nil
    })

    if err != nil {
        log.Fatalln(err)
    }

    return l, m, s
}
