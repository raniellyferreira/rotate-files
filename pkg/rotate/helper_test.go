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

package rotate_test

import (
	"testing"

	"github.com/raniellyferreira/rotate-files/pkg/rotate"
)

func TestGetBucketAndPrefix(t *testing.T) {
	t.Run("Teste com caminho completo válido", func(t *testing.T) {
		fullPath := "s3://my-bucket/my-prefix/file.txt"
		expectedBucket := "my-bucket"
		expectedPrefix := "my-prefix/file.txt"

		bucket, prefix := rotate.GetBucketAndPrefix(fullPath)

		if bucket != expectedBucket || prefix != expectedPrefix {
			t.Errorf("Resultado incorreto. Esperado: (%s, %s), Obtido: (%s, %s)", expectedBucket, expectedPrefix, bucket, prefix)
		}
	})

	t.Run("Teste com caminho sem prefixo", func(t *testing.T) {
		fullPath := "s3://my-bucket/"
		expectedBucket := "my-bucket"
		expectedPrefix := ""

		bucket, prefix := rotate.GetBucketAndPrefix(fullPath)

		if bucket != expectedBucket || prefix != expectedPrefix {
			t.Errorf("Resultado incorreto. Esperado: (%s, %s), Obtido: (%s, %s)", expectedBucket, expectedPrefix, bucket, prefix)
		}
	})

	t.Run("Teste com caminho inválido", func(t *testing.T) {
		fullPath := "s3://"
		expectedBucket := ""
		expectedPrefix := ""

		bucket, prefix := rotate.GetBucketAndPrefix(fullPath)

		if bucket != expectedBucket || prefix != expectedPrefix {
			t.Errorf("Resultado incorreto. Esperado: (%s, %s), Obtido: (%s, %s)", expectedBucket, expectedPrefix, bucket, prefix)
		}
	})
}
