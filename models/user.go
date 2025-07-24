package models

import "time"

type SameFields struct {
	FirstName   string    `json:"firstName" bson:"firstName"`
	LastName    string    `json:"lastName" bson:"lastName"`
	Email       string    `json:"email" bson:"email"`
	BirthDate   time.Time `json:"birthDate" bson:"birthDate"`
	PhoneNumber string    `json:"phoneNumber" bson:"phoneNumber"`
	Role        string    `json:"role" bson:"role"`
}

type Student struct {
	SameFields
	Hostel  string `json:"hostel" bson:"hostel"`
	College string `json:"college" bson:"college"`
	RefId   string `json:"refId" bson:"refId"`
}

type TeachAsst struct {
	SameFields
	Department string `json:"department" bson:"department"`
	College    string `json:"college" bson:"college"`
	StaffId    string `json:"staffId" bson:"staffId"`
}

type Lecturer struct {
	SameFields
	Department string `json:"department" bson:"department"`
	College    string `json:"college" bson:"college"`
	StaffId    string `json:"staffId" bson:"staffId"`
}

type Other struct {
	SameFields
	StaffId string `json:"staffId" bson:"staffId"`
}
