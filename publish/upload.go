package publish

import (
    "github.com/jlaffaye/ftp"
    "log"
    "os"
    "photo-materials/utils"
    "strings"
    "sync"
    "time"
)

type ftpConfig struct {
    url        string
    port       string
    user       string
    passwd     string
    remotePath []string
}

func loadFtpConfig() ftpConfig {
    c := utils.LoadFtpConfig("../config.ini")
    return ftpConfig{
        c.Url,
        c.Port,
        c.User,
        c.Passwd,
        c.RemotePath,
    }
}

func (fc ftpConfig) connectUploadDir(matNo *string, pathNo *int) (conn *ftp.ServerConn, err error) {
    // Connect
    conn, err = ftp.Dial(fc.url+":"+fc.port, ftp.DialWithTimeout(5*time.Second), ftp.DialWithDisabledEPSV(true))
    if err != nil {
        log.Println(err)
        return nil, err
    }

    // Login
    if err = conn.Login(fc.user, fc.passwd); err != nil {
        log.Println(err)
        return nil, err
    }

    // Change to the upload parent dir
    if err = conn.ChangeDir(fc.remotePath[*pathNo]); err != nil {
        log.Println(err)
        return nil, err
    }

    // Make upload dir if not exist & change to that
    // (I think there is a more elegant logic...)
    entries, err := conn.NameList(".")
    if err != nil {
        log.Println(err)
        return nil, err
    }

    uploadDir := "no" + *matNo
    existDir := false
    for _, entry := range entries {
        if strings.Contains(entry, uploadDir) {
            existDir = true
        }
    }
    if !existDir {
        if err = conn.MakeDir(uploadDir); err != nil {
            log.Println(err)
            return nil, err
        }
        log.Printf("[ACT] connectUploadDir [MSG] MakeDir: %s\n", uploadDir)
    }

    if err = conn.ChangeDir(uploadDir); err != nil {
        log.Println(err)
        return nil, err
    }

    return conn, nil
}

type ftpConnects struct {
    cons []*ftp.ServerConn
}

func (cs *ftpConnects) closeConnects() {
    for i, con := range cs.cons {
        if err := con.Quit(); err != nil {
            log.Fatalln(err)
        }
        log.Printf("[ACT] closeConnects [MSG] Conn(%d) closed\n", i+1)
    }
}

type UploadList struct {
    Size  string
    Files []string
}

func (ul UploadList) upload(conn *ftp.ServerConn, wg *sync.WaitGroup) {
    defer wg.Done()
    for _, file := range ul.Files {
        f, err := os.Open(file)
        if err != nil {
            log.Fatalln(err)
        }
        if err = conn.Stor(file, f); err != nil {
            f.Close()
            log.Fatalln(err)
        }
        f.Close()
        log.Printf("[ACT] upload [MSG] Completed: %s\n", file)
    }
    log.Printf("[ACT] upload [MSG] Conn(%s) done\n", ul.Size)
}

// RunPublish uploads some files(jpg/zip) to the target FTP server.
// As of 22.06.11, FTPS conn is not available due to problems with "jlaffaye/ftp".
// Maybe issue#253
// https://github.com/jlaffaye/ftp/issues/253
func RunPublish(matNo *string, remotePathNo *int) {
    // Get listings by size
    l, m, s := utils.FindFiles(".", true)

    var uls []UploadList
    uls = append(uls, UploadList{Size: "L", Files: l})
    uls = append(uls, UploadList{Size: "M", Files: m})
    uls = append(uls, UploadList{Size: "S", Files: s})

    // Create 3 connections
    fc := loadFtpConfig()

    var fs ftpConnects
    for i := 0; i < 3; i++ {
        conn, err := fc.connectUploadDir(matNo, remotePathNo)
        if err != nil {
            if len(fs.cons) > 0 {
                fs.closeConnects()
            }
            log.Fatalln("[ACT] RunPublish [MSG] connectUploadDir:", err)
        }
        fs.cons = append(fs.cons, conn)
        log.Printf("[ACT] RunPublish [MSG] Conn(%d) connected\n", i+1)
    }

    // Uploads by size
    var wg sync.WaitGroup
    for i, ul := range uls {
        wg.Add(1)
        go ul.upload(fs.cons[i], &wg)
    }
    wg.Wait()

    // Close 3 connections
    fs.closeConnects()

    // Create URL text
    RunUrlText(matNo, remotePathNo, uls)
}
