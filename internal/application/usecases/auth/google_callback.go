package auth

import (
	"fmt"
)

type GoogleCallbackUseCase struct {
}


func NewGoogleCallbackUseCase() *GoogleCallbackUseCase {
	return &GoogleCallbackUseCase{
	}
}

func (uc *GoogleCallbackUseCase) Execute()  error {
	return fmt.Errorf("not implemented")
}