package models

import "time"

type (
	Session struct {
		Login     string
		SID       string
		ExpiresAt time.Time
	}

	UserItem struct {
		Login string `json:"login"`
	}

	PostItem struct {
		Id       string    `json:"id"`
		Content  string    `json:"content"`
		CreateAt time.Time `json:"created_at"`
	}

	Post struct {
		ID              int       `json:"id"`
		UserID          int       `json:"user_id"`
		Content         string    `json:"content"`
		CreatedAt       time.Time `json:"created_at"`
		CommentsAllowed bool      `json:"comments_allowed"`
	}

	Comment struct {
		ID        int       `json:"id"`
		UserID    int       `json:"user_id"`
		PostID    int       `json:"post_id"`
		ParentID  int       `json:"parent_id"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}

	House struct {
		HouseID            int       `json:"house_id"`
		Address            string    `json:"address"`
		YearBuilt          int       `json:"year_built"`
		Developer          string    `json:"developer"`
		CreatedAt          time.Time `json:"created_at"`
		LastApartmentAdded time.Time `json:"last_apartment_added"`
	}

	HouseFlat struct {
		ApartmentID     int    `json:"apartment_id"`
		ApartmentNumber int    `json:"apartment_number"`
		Price           int    `json:"price"`
		Rooms           int    `json:"rooms"`
		HouseID         int    `json:"house_id"`
		Address         string `json:"address"`
		Status          string `json:"status"`
	}
)
