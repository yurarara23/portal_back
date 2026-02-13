package models

import "gorm.io/gorm"

type Member struct {
    gorm.Model
    username      string `json:"username"`
    password      string `json:"password"`
    department    string `json:"department"`
    grade         string `json:"grade"`
    student_num string `json:"student_num"`
}