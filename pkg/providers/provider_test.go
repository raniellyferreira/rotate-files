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

package providers_test

import (
	"fmt"
	"testing"

	"github.com/golang-module/carbon"
	"github.com/raniellyferreira/rotate-files/pkg/providers"
)

func TestFileInfo(t *testing.T) {
	t.Run("Teste de criação de FileInfo com valores válidos", func(t *testing.T) {
		expectedPath := "/path/to/file.txt"
		expectedSize := int64(1024)
		expectedTimestamp := carbon.Now()

		fileInfo := &providers.FileInfo{
			Path:      expectedPath,
			Size:      expectedSize,
			Timestamp: expectedTimestamp,
		}

		if fileInfo.Path != expectedPath {
			t.Errorf("Esperado Path: %s, Obtido: %s", expectedPath, fileInfo.Path)
		}

		if fileInfo.Size != expectedSize {
			t.Errorf("Esperado Size: %d, Obtido: %d", expectedSize, fileInfo.Size)
		}

		if !fileInfo.Timestamp.Eq(expectedTimestamp) {
			t.Errorf("Esperado Timestamp: %v, Obtido: %v", expectedTimestamp, fileInfo.Timestamp)
		}
	})

	t.Run("Teste de criação de FileInfo com valores zero", func(t *testing.T) {
		fileInfo := &providers.FileInfo{}

		if fileInfo.Path != "" {
			t.Errorf("Esperado Path vazio, Obtido: %s", fileInfo.Path)
		}

		if fileInfo.Size != 0 {
			t.Errorf("Esperado Size: 0, Obtido: %d", fileInfo.Size)
		}

		if !fileInfo.Timestamp.IsZero() {
			t.Errorf("Esperado Timestamp zero, Obtido: %v", fileInfo.Timestamp)
		}
	})

	t.Run("Teste com tamanho negativo", func(t *testing.T) {
		fileInfo := &providers.FileInfo{
			Path: "/test/file.txt",
			Size: -100,
		}

		if fileInfo.Size != -100 {
			t.Errorf("Esperado Size: -100, Obtido: %d", fileInfo.Size)
		}
	})
}

// MockProvider implementa a interface Provider para testes
type MockProvider struct {
	DeleteFunc    func(fullPath string) error
	ListFilesFunc func(fullPath string) ([]*providers.FileInfo, error)
}

func (m *MockProvider) Delete(fullPath string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(fullPath)
	}
	return nil
}

func (m *MockProvider) ListFiles(fullPath string) ([]*providers.FileInfo, error) {
	if m.ListFilesFunc != nil {
		return m.ListFilesFunc(fullPath)
	}
	return []*providers.FileInfo{}, nil
}

func TestProviderInterface(t *testing.T) {
	t.Run("Teste de implementação da interface Provider", func(t *testing.T) {
		// Testa se MockProvider implementa corretamente a interface Provider
		var provider providers.Provider = &MockProvider{}

		// Testa método Delete
		err := provider.Delete("/test/path")
		if err != nil {
			t.Errorf("Erro inesperado no Delete: %v", err)
		}

		// Testa método ListFiles
		files, err := provider.ListFiles("/test/path")
		if err != nil {
			t.Errorf("Erro inesperado no ListFiles: %v", err)
		}

		if files == nil {
			t.Error("ListFiles retornou nil em vez de slice vazio")
		}
	})

	t.Run("Teste de MockProvider com funções personalizadas", func(t *testing.T) {
		expectedError := "erro personalizado"
		expectedFiles := []*providers.FileInfo{
			{Path: "/file1.txt", Size: 100, Timestamp: carbon.Now()},
			{Path: "/file2.txt", Size: 200, Timestamp: carbon.Now().AddDays(1)},
		}

		mock := &MockProvider{
			DeleteFunc: func(fullPath string) error {
				if fullPath == "" {
					return nil
				}
				return nil
			},
			ListFilesFunc: func(fullPath string) ([]*providers.FileInfo, error) {
				if fullPath == "/error/path" {
					return nil, fmt.Errorf(expectedError)
				}
				return expectedFiles, nil
			},
		}

		// Testa Delete com sucesso
		err := mock.Delete("/valid/path")
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		// Testa ListFiles com sucesso
		files, err := mock.ListFiles("/valid/path")
		if err != nil {
			t.Errorf("Erro inesperado: %v", err)
		}

		if len(files) != len(expectedFiles) {
			t.Errorf("Esperado %d arquivos, Obtido: %d", len(expectedFiles), len(files))
		}

		// Testa ListFiles com erro
		files, err = mock.ListFiles("/error/path")
		if err == nil {
			t.Error("Esperava um erro, mas nenhum ocorreu")
		}

		if err.Error() != expectedError {
			t.Errorf("Esperado erro: %s, Obtido: %v", expectedError, err)
		}

		if files != nil {
			t.Error("Esperava files nil quando há erro")
		}
	})
}