package types

import "github.com/google/uuid"

func PopUUID(data *[]byte) (u uuid.UUID, err error) {
	u, err = uuid.FromBytes(*data)
	if err != nil {
		return u, err
	}
	*data = (*data)[16:]
	return u, err
}
