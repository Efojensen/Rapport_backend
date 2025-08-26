package models

import (
	"time"
)

type SameFields struct {
	FirstName   string             `json:"firstName" bson:"firstName"`
	LastName    string             `json:"lastName" bson:"lastName"`
	Email       string             `json:"email" bson:"email"`
	BirthDate   time.Time          `json:"birthDate" bson:"birthDate"`
	PhoneNumber string             `json:"phoneNumber" bson:"phoneNumber"`
	Role        string             `json:"role" bson:"role"`
	Username    string             `json:"username" bson:"username"`
	Password    string             `json:"password" bson:"password"`
	Avatar      string             `json:"avatar" bson:"avatar"`
}

type SOS_Alert struct {
	Name        string
	PhoneNumber string
	Role        string
	Hostel      *string
	College     *string
}

type Student struct {
	SameFields `bson:",inline"`
	Hostel     string `json:"hostel" bson:"hostel"`
	College    string `json:"college" bson:"college"`
	RefId      string `json:"refId" bson:"refId"`
}

func (std *Student) UserDetails() SOS_Alert {
	return SOS_Alert{
		Name:        std.FirstName + " " + std.LastName,
		PhoneNumber: std.PhoneNumber,
		Role:        std.Role,
		Hostel:      &std.Hostel,
		College:     &std.College,
	}
}

type TeachAsst struct {
	SameFields `bson:",inline"`
	Department string `json:"department" bson:"department"`
	College    string `json:"college" bson:"college"`
	StaffId    string `json:"staffId" bson:"staffId"`
}

func (ta *TeachAsst) UserDetails() SOS_Alert {
	return SOS_Alert{
		Name:        ta.FirstName + " " + ta.LastName,
		PhoneNumber: ta.PhoneNumber,
		Role:        ta.Role,
		College:     &ta.College,
	}
}

type Lecturer struct {
	SameFields `bson:",inline"`
	Department string `json:"department" bson:"department"`
	College    string `json:"college" bson:"college"`
	StaffId    string `json:"staffId" bson:"staffId"`
}

func (lec *Lecturer) UserDetails() SOS_Alert {
	return SOS_Alert{
		Name:        lec.FirstName + " " + lec.LastName,
		PhoneNumber: lec.PhoneNumber,
		Role:        lec.Role,
		College:     &lec.College,
	}
}

type Other struct {
	SameFields `bson:",inline"`
	StaffId    string `json:"staffId" bson:"staffId"`
}

func (otr *Other) UserDetails() SOS_Alert {
	return SOS_Alert{
		Name:        otr.FirstName + " " + otr.LastName,
		PhoneNumber: otr.PhoneNumber,
		Role:        otr.Role,
	}
}

type User interface {
	UserDetails() SOS_Alert
}
