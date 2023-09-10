package handler

import (
	"strconv"
)

type InputDTO struct {
	Id string `json:"id"`
}

func (i *InputDTO) toModel() (uint64, error) {
	id, err := strconv.Atoi(i.Id)
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}
