package models

type SmallBidError struct {
	input string
}

type ForeignKeyConstraint struct {
	input string
}

type UserBidOnOwnBid struct {
	input string
}

type UserDoesNotExist struct {
	input string
}

type NoRowsError struct {
	input string
}

type AlreadyExistsInDatabaseError struct {
	input string
}

type ItemDoesNotExist struct {
	input string
}

func (s SmallBidError) Error() string {
	return s.input
}

func (f ForeignKeyConstraint) Error() string {
	return f.input
}

func (f UserDoesNotExist) Error() string {
	return f.input
}

func (f ItemDoesNotExist) Error() string {
	return f.input
}

func (f UserBidOnOwnBid) Error() string {
	return f.input
}

func (f NoRowsError) Error() string {
	return f.input
}

func (f AlreadyExistsInDatabaseError) Error() string {
	return f.input
}
