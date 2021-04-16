package helper

import (
	"sort"
	"time"
)

const OBJECT_PRESERVE = "PRESERVE"
const OBJECT_WAITING = "WAITING"
const OBJECT_DELETE = "DELETE"

type RotationObject struct {
	Bucket    *string
	Path      *string
	CreatedAt *time.Time
	Status    string
}

type RotationScheme struct {
	Hourly  int
	Daily   int
	Weekly  int
	Monthly int
	Yearly  int
	DryRun  bool
}

func RotateObjects(filesObj []RotationObject, rotationScheme *RotationScheme) {

	// Ordena os objetos pela data de criação
	sort.Slice(filesObj, func(i, j int) bool {
		return filesObj[i].CreatedAt.Before(*filesObj[j].CreatedAt)
	})

	// TODO parei aqui fazer a rotaçao dos arquivos
}
