/*
 *		D.T.Os (Data Transfer Objects)
 *
 *		All structs for packaging data to pass
 *		through to the pages
 *
 */

package inf

type FormUser struct {
	Name            string `form:"name"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	NewPassword     string `form:"new_password"`
	PasswordConfirm string `form:"pass_conf"`
	StayLoggedIn    string `form:"stay_logged_in"`
}

type FormItem struct {
	Name        string `form:"name"`
	Description string `form:"desc"`
	Price       string `form:"price"`
	LinkText    string `form:"link-text"`
	LinkURL     string `form:"link-url"`
}

type UserFormError struct {
	Name            string
	Email           string
	NewPassword     string
	Password        string
	PasswordConfirm string
}

type ItemFormError struct {
	Name  string
	Price string
	Link  string
}

type StatusPageData struct {
	Colour          string
	MainMessage     string
	NextPageURL     string
	NextPageMessage string
}

type ConfirmPageData struct {
	MainTitle       string
	MainDescription string
	YesURL          string
	YesColour       string
	YesText         string
	NoURL           string
	NoColour        string
	NoText          string
}
