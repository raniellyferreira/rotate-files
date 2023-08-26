/*
Copyright The Rotate Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package files_test

import (
	"os"
	"testing"

	"github.com/raniellyferreira/rotate-files/pkg/files"
)

func TestListDir(t *testing.T) {
	t.Run("Teste com diretório válido", func(t *testing.T) {
		dirPath := "testdir"
		filePaths := []string{
			"testdir/file1.txt",
			"testdir/file2.txt",
		}

		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dirPath)

		for _, path := range filePaths {
			file, err := os.Create(path)
			if err != nil {
				t.Fatal(err)
			}
			file.Close()
		}

		files, err := files.ListDir(dirPath)
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		expectedLen := len(filePaths)
		if len(files) != expectedLen {
			t.Errorf("Resultado incorreto. Esperado: %d arquivos, Obtido: %d arquivos", expectedLen, len(files))
		}

		for _, path := range filePaths {
			err = os.Remove(path)
			if err != nil {
				t.Fatal(err)
			}
		}
	})

	t.Run("Teste com diretório inexistente", func(t *testing.T) {
		dirPath := "nonexistentdir"

		_, err := files.ListDir(dirPath)
		if err == nil {
			t.Errorf("Esperava um erro, mas nenhum ocorreu")
		}
	})
}

func TestDeleteLocalFile(t *testing.T) {
	t.Run("Teste com arquivo existente", func(t *testing.T) {
		path := "testfile.txt"

		file, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()

		err = files.DeleteLocalFile(path)
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		_, err = os.Stat(path)
		if !os.IsNotExist(err) {
			t.Errorf("O arquivo não foi excluído corretamente")
		}
	})

	t.Run("Teste com arquivo inexistente", func(t *testing.T) {
		path := "nonexistentfile.txt"

		err := files.DeleteLocalFile(path)
		if err == nil {
			t.Errorf("Esperava um erro, mas nenhum ocorreu")
		}
	})
}
