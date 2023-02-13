# Rotate Files

Rotate files locally or in S3 bucket based on backup rotation scheme

<https://en.wikipedia.org/wiki/Backup_rotation_scheme#Grandfather-father-son>

## Usage

```bash
rotate --help
```

```console
Rotate files locally or in S3 bucket based on backup rotation scheme

Usage:
   rotate [path] {flags}
   rotate <command> {flags}

Commands: 
   help                          displays usage informationn
   version                       displays version number

Arguments: 
   path                          local directory path or s3:// path (default: ./)

Flags: 
   -d, --daily                   number of daily backups to preserve (default: 7)
   -d, --dry-run                 simulate deletion process (default: false)
   -h, --help                    displays usage information of the application or a command (default: false)
   -h, --hourly                  number of hourly backups to preserve (default: 24)
   -m, --monthly                 number of monthly backups to preserve (default: 12)
   -v, --version                 displays version number (default: false)
   -w, --weekly                  number of weekly backups to preserve (default: 14)
   -y, --yearly                  number of yearly backups to preserve, set -1 to preserver always (default: -1)
```
