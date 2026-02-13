package models

import "gorm.io/gorm"

type Member struct {
    gorm.Model
    Username      string `json:"username"`
    Password      string `json:"password"`
    Department    string `json:"department"`
    Grade         string `json:"grade"`
    Student_num string `json:"student_num"`
}