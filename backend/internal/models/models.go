package models

import (
	"time"
)

type WebProjects struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=15,max=100"`
	Description string    `json:"description" db:"description" validate:"required,min=20,max=1500"`
	Img         string    `json:"img" db:"img" validate:"reuired,url"`
	Price       float64   `json:"price" db:"price" validate:"required,min=0"`
	TimeDevelop int       `json:"time_develop" db:"time_develop" validate:"required,min=1,max=1825"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdateAt    time.Time `json:"update_at" db:"update_at"`
}

type MobileProjects struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=15,max=100"`
	Description string    `json:"description" db:"description" validate:"required,min=20,max=1500"`
	Img         string    `json:"img" db:"img" validate:"reuired,url"`
	Price       float64   `json:"price" db:"price" validate:"required,min=0"`
	TimeDevelop int       `json:"time_develop" db:"time_develop" validate:"required,min=1,max=1825"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdateAt    time.Time `json:"update_at" db:"update_at"`
}

type BotsProjects struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=15,max=100"`
	Description string    `json:"description" db:"description" validate:"required,min=20,max=1500"`
	Img         string    `json:"img" db:"img" validate:"reuired,url"`
	Price       float64   `json:"price" db:"price" validate:"required,min=0"`
	TimeDevelop int       `json:"time_develop" db:"time_develop" validate:"required,min=1,max=1825"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdateAt    time.Time `json:"update_at" db:"update_at"`
}

type HealthResponse struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Timestamp map[string]interface{} `json:"timestamp"`
	Database  string                 `json:"database,omitempty"`
}

type CreateWebProjectRequest struct {
	Name string `json:"name" validate:"required,min=15,max=100"`
	Description string `json:"description" validate:"required,min=20,max=1500"`
	Img string `json:"img" validate:"required,url"`
	Price float64 `json:"price" validate:"required,min=0"`
	TimeDevelop int `json:"time_develop" validate:"required,min=1,max=1825"`
}

type CreateMobileProjectRequest struct {
	Name string `json:"name" validate:"required,min=15,max=100"`
	Description string `json:"description" validate:"required,min=20,max=1500"`
	Img string `json:"img" validate:"required,url"`
	Price float64 `json:"price" validate:"required,min=0"`
	TimeDevelop int `json:"time_develop" validate:"required,min=1,max=1825"`
}

type CreateBotsProjectRequest struct {
	Name string `json:"name" validate:"required,min=15,max=100"`
	Description string `json:"description" validate:"required,min=20,max=1500"`
	Img string `json:"img" validate:"required,url"`
	Price float64 `json:"price" validate:"required,min=0"`
	TimeDevelop int `json:"time_develop" validate:"required,min=1,max=1825"`
}

type GetProjectRequest struct {
    ID int `json:"id" uri:"id" validate:"required,min=1"`
}