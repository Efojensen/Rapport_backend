package models

import "time"

type SameFields struct {
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	BirthDate   time.Time `json:"birthDate"`
	PhoneNumber string    `json:"phoneNumber"`
}

type Student struct {
	SameFields
	Hostel  string `json:"hostel"`
	College string `json:"college"`
	RefId   string `json:"refId"`
}

type TeachAsst struct {
	SameFields
	Department string `json:"department"`
	College    string `json:"college"`
	StaffId    string `json:"staffId"`
}

type Lecturer struct {
	SameFields
	Department string `json:"department"`
	College    string `json:"college"`
	StaffId    string `json:"staffId"`
}

type Other struct {
	SameFields
	StaffId    string `json:"staffId"`
}