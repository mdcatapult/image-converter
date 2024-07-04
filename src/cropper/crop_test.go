/*
 * Copyright 2024 Medicines Discovery Catapult
 * Licensed under the Apache License, Version 2.0 (the "Licence");
 * you may not use this file except in compliance with the Licence.
 * You may obtain a copy of the Licence at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the Licence is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the Licence for the specific language governing permissions and
 * limitations under the Licence.
 */

package cropper

import (
	"fmt"
	"testing"

	"github.com/go-playground/assert/v2"

	"errors"
)

func TestValidateCoords(t *testing.T) {

	croppedImageSize := int64(100)
	rawImageHeight := 1000
	rawImageWidth := 1000

	imageMetadata := fmt.Sprintf("Height = %d\n Width = %d", rawImageHeight, rawImageWidth)

	for _, test := range []struct {
		name        string
		x, y        int64
		expectedErr error
	}{
		{
			name:        "x too small",
			x:           -50, // less than 0, so x is too small
			y:           500,
			expectedErr: errors.New(ErrOutOfBounds),
		}, {
			name:        "y too small",
			x:           500,
			y:           -50,
			expectedErr: errors.New(ErrOutOfBounds),
		}, {
			name:        "x too large",
			x:           950, // 950 + croppedImageSize is over 1000 (size as per imageMetadata), so x is too large
			y:           500,
			expectedErr: errors.New(ErrOutOfBounds),
		}, {
			name:        "y too large",
			x:           500,
			y:           950,
			expectedErr: errors.New(ErrOutOfBounds),
		}, {
			name:        "x and y in bounds",
			x:           500,
			y:           500,
			expectedErr: nil,
		}} {

		err := validateCoords(test.x, test.y, croppedImageSize, "", func(string) (string, error) {
			return imageMetadata, nil
		})

		if test.expectedErr == nil {
			assert.Equal(t, test.expectedErr, nil)
		} else {
			assert.Equal(t, test.expectedErr.Error(), err.Error())
		}
	}
}
