package analytics

import "errors"

var (
	// ErrPassUserID throws when you don't provided userId
	ErrPassUserID = errors.New("You must pass a 'userId'")

	// ErrPassPreviousID throws when you don't provided previousId
	ErrPassPreviousID = errors.New("You must pass a 'previousId'")

	// ErrPassEvent throws when you don't provided event
	ErrPassEvent = errors.New("You must pass 'event'")

	// ErrPassGroupID throws when you don't provided groupId
	ErrPassGroupID = errors.New("You must pass a 'groupId'")

	// ErrPassAnonymousOrUser throws when you don't provided anonymousId & userId
	ErrPassAnonymousOrUser = errors.New("You must pass either an 'anonymousId' or 'userId'")
)
