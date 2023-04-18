// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Order struct {
	ID            string  `json:"ID"`
	UserID        string  `json:"UserID"`
	DriverID      string  `json:"DriverID"`
	DriverName    string  `json:"DriverName"`
	DriverPhone   string  `json:"DriverPhone"`
	DriverRaiting float64 `json:"DriverRaiting"`
	TaxiType      string  `json:"TaxiType"`
	From          string  `json:"From"`
	To            string  `json:"To"`
	Date          string  `json:"Date"`
	Status        Status  `json:"Status"`
}

type OrderInfo struct {
	TaxiType string `json:"TaxiType"`
	From     string `json:"From"`
	To       string `json:"To"`
}

type OrderState struct {
	ID    string `json:"ID"`
	State Status `json:"State"`
}

type Raiting struct {
	ID      string  `json:"ID"`
	Raiting float64 `json:"Raiting"`
}

type Status string

const (
	StatusWaiting    Status = "waiting"
	StatusFound      Status = "found"
	StatusInProgress Status = "inProgress"
	StatusFinished   Status = "finished"
	StatusCanceled   Status = "canceled"
)

var AllStatus = []Status{
	StatusWaiting,
	StatusFound,
	StatusInProgress,
	StatusFinished,
	StatusCanceled,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusWaiting, StatusFound, StatusInProgress, StatusFinished, StatusCanceled:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
