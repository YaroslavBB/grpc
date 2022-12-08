package entity

import "errors"

var (
	// ErrNoData данные не найдены
	ErrNoData = errors.New("данные не найдены")
	// ErrIncorrectFileType некорректный тип файла
	ErrIncorrectFileType = errors.New("некорректный тип файла")
)
