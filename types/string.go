package types

func PopString(data *[]byte) (s string, err error) {
	length, err := PopVarInt(data)
	if err != nil {
		return
	}
	s = string((*data)[:length])
	*data = (*data)[length:]
	return
}
