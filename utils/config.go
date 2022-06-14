package utils

import (
    "gopkg.in/ini.v1"
    "log"
)

type FtpConfig struct {
    ServerName string
    Url        string
    Port       string
    User       string
    Passwd     string
    RemotePath []string
}

func LoadFtpConfig(path string) FtpConfig {
    cfg := load(path)
    return FtpConfig{
        ServerName: cfg.Section("ftp").Key("server_name").String(),
        Url:        cfg.Section("ftp").Key("server").String(),
        Port:       cfg.Section("ftp").Key("port").String(),
        User:       cfg.Section("ftp").Key("user").String(),
        Passwd:     cfg.Section("ftp").Key("passwd").String(),
        RemotePath: []string{
            cfg.Section("path").Key("test_path").String(),
            cfg.Section("path").Key("product_path1").String(),
            cfg.Section("path").Key("product_path2").String(),
            cfg.Section("path").Key("product_path3").String(),
            cfg.Section("path").Key("product_path4").String(),
            cfg.Section("path").Key("product_path5").String(),
        },
    }
}

type UrlTextConfig struct {
    RemoteHome string
    RemotePath []string
}

func LoadPathConfig(path string) UrlTextConfig {
    cfg := load(path)
    return UrlTextConfig{
        RemoteHome: cfg.Section("path").Key("remote_home").String(),
        RemotePath: []string{
            cfg.Section("path").Key("test_path").String(),
            cfg.Section("path").Key("product_path1").String(),
            cfg.Section("path").Key("product_path2").String(),
            cfg.Section("path").Key("product_path3").String(),
            cfg.Section("path").Key("product_path4").String(),
            cfg.Section("path").Key("product_path5").String(),
        },
    }
}

func load(path string) *ini.File {
    f, err := ini.Load(path)
    if err != nil {
        log.Fatalln(err)
    }
    return f
}
