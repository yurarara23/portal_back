package models

import "gorm.io/gorm"

type Member struct {
    gorm.Model
    UserName      string `json:"username"`
    Password      string `json:"password"`
    Department    string `json:"department"`
    Grade         string `json:"grade"`
    StudentNumber string `json:"student_num"`
}

type Rental struct {
    gorm.Model
    MemberID uint   
    ItemName string 
    ImageURL string `json:"image_url"`
    IsRented bool   `json:"isRented"`
}