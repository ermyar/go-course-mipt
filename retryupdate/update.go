//go:build !solution

package retryupdate

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

var (
	conflictError *kvapi.ConflictError
	authError     *kvapi.AuthError
	apiError      *kvapi.APIError
)

func getRealError(err error) error {
	if err == nil {
		return nil
	}
	if errors.As(err, &apiError) {
		return err.(*kvapi.APIError).Unwrap()
	}
	return nil
}

var m map[uuid.UUID]bool

func set(c kvapi.Client, key string, newValue *string, oldVersion *uuid.UUID,
	updateFn func(oldValue *string) (newValue string, err error)) error {
	setReq := kvapi.SetRequest{Key: key, Value: *newValue, NewVersion: uuid.Must(uuid.NewV4())}
	if oldVersion != nil {
		setReq.OldVersion = *oldVersion
	}
	fmt.Println("hello1 ", setReq)
	_, err := c.Set(&setReq)
	err1 := getRealError(err)
	if err1 == nil {
		return nil
	} else if errors.As(err1, &authError) {
		fmt.Println(err)
		return err
	} else if errors.Is(err1, kvapi.ErrKeyNotFound) {
		newValue, err := updateFn(nil)
		if err != nil {
			return err
		}
		return set(c, key, &newValue, nil, updateFn)
	} else if errors.As(err1, &conflictError) {
		tmp := err1.(*kvapi.ConflictError)
		fmt.Println(tmp)
		if m[tmp.ExpectedVersion] {
			return nil
		}

		return getAndSet(c, key, updateFn)
	} else {
		m[setReq.NewVersion] = true
		return set(c, key, newValue, oldVersion, updateFn)
	}
}

func getAndSet(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	getResp, err := c.Get(&kvapi.GetRequest{Key: key})
	err1 := getRealError(err)
	if err == nil {
		newValue, err := updateFn(&getResp.Value)
		if err != nil {
			return err
		}
		set(c, key, &newValue, &getResp.Version, updateFn)
	} else if errors.Is(err1, kvapi.ErrKeyNotFound) {
		newValue, err := updateFn(nil)
		if err != nil {
			return err
		}
		return set(c, key, &newValue, nil, updateFn)
	} else if errors.As(err1, &authError) {
		return err
	} else {
		return getAndSet(c, key, updateFn)
	}
	return nil
}

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	m = make(map[uuid.UUID]bool)
	defer func() {
		for u := range m {
			delete(m, u)
		}
	}()

	return getAndSet(c, key, updateFn)
}
