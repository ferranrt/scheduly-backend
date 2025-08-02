package config

import "sync"

/*
Required constants that are not injectable from the environment variables
*/

type InternalConstants struct {
	DefaultLanguage string
}

var (
	internalSingleton         sync.Once
	InternalConstantsInstance *InternalConstants
)

func NewInternalConstants() *InternalConstants {
	internalSingleton.Do(func() {
		InternalConstantsInstance = &InternalConstants{
			DefaultLanguage: "en",
		}
	})

	return InternalConstantsInstance
}
