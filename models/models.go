package models

import (
    "gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
)

type Member struct {
    gorm.Model
    Username    string `json:"username"`
    Password    string `json:"password"`
    Department  string `json:"department"`
    Grade       string `json:"grade"`
    Student_num string `json:"student_num"`
}

func (m *Member) HashPassword() error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    m.Password = string(hashedPassword)
    return nil
}