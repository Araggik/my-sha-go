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

	isRetry := true

	getReq := &kvapi.GetRequest{Key: key}

	for isRetry {
		getRes, err := c.Get(getReq)

		resultErr = err

		if resultErr == nil {
			newVal, err := updateFn(&getRes.Value)

			resultErr = err

			if resultErr == nil {
				setReg := &kvapi.SetRequest{
					Key:        key,
					Value:      newVal,
					OldVersion: getRes.Version,
					NewVersion: uuid.Must(uuid.NewV4()),
				}

				_, err := c.Set(setReg)

				resultErr = err

				if errors.As(resultErr, &authErr) {
					isRetry = false
				}
			} else {
				isRetry = false
			}
		} else if errors.As(err, &authErr) {
			isRetry = false
		}

		if resultErr == nil {
			isRetry = false
		}
	}

	return resultErr
}
