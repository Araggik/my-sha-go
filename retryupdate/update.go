//go:build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"

	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

var authErr *kvapi.AuthError
var conflictErr *kvapi.ConflictError

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var resultErr error

	isGetRetry := true

	getReq := &kvapi.GetRequest{Key: key}

	for isGetRetry {
		getRes, err := c.Get(getReq)

		resultErr = err

		isKeyNoFound := errors.Is(resultErr, kvapi.ErrKeyNotFound)

		if resultErr == nil || isKeyNoFound {
			isSetRetry := true

			var prevVersion *uuid.UUID

			for isSetRetry {
				var value *string
				var oldVersion uuid.UUID
				newVersion := uuid.Must(uuid.NewV4())

				if !isKeyNoFound {
					value = &getRes.Value
					oldVersion = getRes.Version
				}

				newVal, err := updateFn(value)

				resultErr = err

				if resultErr == nil {
					setReg := &kvapi.SetRequest{
						Key:        key,
						Value:      newVal,
						OldVersion: oldVersion,
						NewVersion: newVersion,
					}

					_, err := c.Set(setReg)

					resultErr = err

					if resultErr == nil || errors.As(resultErr, &authErr) {
						isSetRetry = false
						isGetRetry = false
					} else if errors.As(resultErr, &conflictErr) {
						var a any = errors.Unwrap(err)

						originErr := a.(*kvapi.ConflictError)

						if originErr.ExpectedVersion == *prevVersion {
							isGetRetry = false
							resultErr = nil
						}

						isSetRetry = false
					} else if errors.Is(resultErr, kvapi.ErrKeyNotFound) {
						isKeyNoFound = true
					}

					prevVersion = &newVersion
				} else {
					isGetRetry = false
					isSetRetry = false
				}
			}
		} else if errors.As(err, &authErr) {
			isGetRetry = false
		}

		if resultErr == nil {
			isGetRetry = false
		}
	}

	return resultErr
}
