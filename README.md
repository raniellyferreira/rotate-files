# Rotate Files

[![Go Report Card](https://goreportcard.com/badge/github.com/raniellyferreira/rotate-files)](https://goreportcard.com/report/github.com/raniellyferreira/rotate-files) ![License](https://img.shields.io/github/license/raniellyferreira/rotate-files)

Rotate files locally or in S3 bucket based on custom backup rotation scheme

## Installation

```bash
curl -SsL https://raw.githubusercontent.com/raniellyferreira/rotate-files/master/environment/scripts/get | bash
```

## Usage

```bash
rotate help
```

```console
Rotate files locally or in S3 bucket based on backup rotation scheme

Usage:
   rotate <path> {flags}
   rotate <command> {flags}

Commands: 
   help                          displays usage informationn
   version                       displays version number

Arguments: 
   path                          local directory path or s3:// path

Flags: 
   -d, --daily                   number of daily backups to preserve (default: 7)
   -D, --dry-run                 simulate deletion process (default: false)
   -h, --help                    displays usage information of the application or a command
   -h, --hourly                  number of hourly backups to preserve (default: 24)
   -m, --monthly                 number of monthly backups to preserve (default: 12)
   -v, --version                 displays version number
   -w, --weekly                  number of weekly backups to preserve (default: 14)
   -y, --yearly                  number of yearly backups to preserve, set 0 to preserver always (default: 0)
```

### Use example

```bash
rotate s3://example-bucket/backups/ --hourly 24 --daily 7 --weekly 10 --monthly 12
```

rotate output:

```console
2023/02/15 13:28:07 Starting rotation on s3://example-bucket/backups
2023/02/15 13:28:07 Yearly matched:
2023/02/15 13:28:07   backups/mysql-2021-03-03-0607.tar.gz 2021-03-03 09:09:34
2023/02/15 13:28:07 Monthly matched:
2023/02/15 13:28:07   backups/mysql-2023-01-16-0201.tar.gz 2023-01-16 05:40:08
2023/02/15 13:28:07   backups/mysql-2022-12-19-0201.tar.gz 2022-12-19 05:40:56
2023/02/15 13:28:07   backups/mysql-2022-11-01-0201.tar.gz 2022-11-01 05:36:10
2023/02/15 13:28:07   backups/mysql-2022-10-01-0201.tar.gz 2022-10-01 05:32:32
2023/02/15 13:28:07   backups/mysql-2022-09-25-0245.tar.gz 2022-09-25 06:11:53
2023/02/15 13:28:07   backups/mysql-2022-08-01-0201.tar.gz 2022-08-01 05:15:59
2023/02/15 13:28:07   backups/mysql-2022-07-01-0201.tar.gz 2022-07-01 05:16:58
2023/02/15 13:28:07   backups/mysql-2022-06-01-0201.tar.gz 2022-06-01 05:16:04
2023/02/15 13:28:07   backups/mysql-2022-05-01-0201.tar.gz 2022-05-01 05:16:06
2023/02/15 13:28:07   backups/mysql-2022-04-01-0201.tar.gz 2022-04-01 05:11:35
2023/02/15 13:28:07   backups/mysql-2022-03-01-0201.tar.gz 2022-03-01 05:15:01
2023/02/15 13:28:07   backups/mysql-2022-01-01-0201.tar.gz 2022-01-01 05:37:37
2023/02/15 13:28:07 Weekly matched:
2023/02/15 13:28:07   backups/mysql-2023-01-01-0201.tar.gz 2023-01-01 05:37:16
2023/02/15 13:28:07 Daily matched:
2023/02/15 13:28:07   backups/mysql-2023-02-14-0201.tar.gz 2023-02-14 05:42:08
2023/02/15 13:28:07   backups/mysql-2023-02-13-0201.tar.gz 2023-02-13 05:41:46
2023/02/15 13:28:07   backups/mysql-2023-02-12-0201.tar.gz 2023-02-12 05:40:40
2023/02/15 13:28:07   backups/mysql-2023-02-11-0201.tar.gz 2023-02-11 05:40:28
2023/02/15 13:28:07   backups/mysql-2023-02-10-0201.tar.gz 2023-02-10 05:46:37
2023/02/15 13:28:07   backups/mysql-2023-02-09-0201.tar.gz 2023-02-09 05:46:39
2023/02/15 13:28:07   backups/mysql-2023-02-08-0201.tar.gz 2023-02-08 05:48:30
2023/02/15 13:28:07 Hourly matched:
2023/02/15 13:28:07   backups/mysql-2023-02-15-1245.tar.gz 2023-02-15 16:27:01
2023/02/15 13:28:07   backups/mysql-2023-02-15-1126.tar.gz 2023-02-15 15:09:08
2023/02/15 13:28:07   backups/mysql-2023-02-15-1043.tar.gz 2023-02-15 14:25:47
2023/02/15 13:28:07 Deleted:
2023/02/15 13:28:08   backups/mysql-2023-02-06-0201.tar.gz 2023-02-06 05:45:07
2023/02/15 13:28:08   backups/mysql-2023-02-01-0201.tar.gz 2023-02-01 05:40:06
2023/02/15 13:28:09   backups/mysql-2023-01-30-0201.tar.gz 2023-01-30 05:39:51
2023/02/15 13:28:09   backups/mysql-2023-01-23-0201.tar.gz 2023-01-23 05:40:46
2023/02/15 13:28:09   backups/mysql-2023-01-09-0201.tar.gz 2023-01-09 05:38:39
2023/02/15 13:28:09   backups/mysql-2023-01-02-0201.tar.gz 2023-01-02 05:37:57
2023/02/15 13:28:09   backups/mysql-2022-12-12-0201.tar.gz 2022-12-12 05:40:55
2023/02/15 13:28:09   backups/mysql-2022-12-01-0201.tar.gz 2022-12-01 05:39:11
```
