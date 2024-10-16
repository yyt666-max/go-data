package cache

import "encoding/json"

func decodeList[T any](bytes []byte) ([]T, error) {

	t := make([]T, 0)
	err := json.Unmarshal(bytes, &t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
func encode[T any](t T) ([]byte, error) {

	bytes, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return bytes, nil

}
func decode[T any](bytes []byte) (*T, error) {

	t := new(T)
	err := json.Unmarshal(bytes, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
