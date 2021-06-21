package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	DEFAULT_LOCALE = "en"

	USER_EMAIL_MAX_LENGTH = 32

	USER_FIRST_NAME_MAX_RUNES = 64
	USER_LAST_NAME_MAX_RUNES  = 64

	USER_PASSWORD_MAX_LENGTH       = 32
	USER_CONTACT_NUMBER_MAX_LENGTH = 20
)

// User contains the details about the user.
type User struct {
	ID                      string            `json:"id"`
	CreatedAt               int64             `json:"created_at,omitempty"`
	UpdatedAt               int64             `json:"updated_at,omitempty"`
	DeletedAt               int64             `json:"deleted_at"`
	Username                string            `json:"username"`
	Password                string            `json:"password,omitempty"`
	AuthData                *string           `json:"auth_data,omitempty"`
	AuthService             string            `json:"auth_service"`
	Email                   string            `json:"email"`
	EmailVerified           bool              `json:"email_verified,omitempty"`
	FirstName               string            `json:"first_name"`
	LastName                string            `json:"last_name"`
	ContactNumber           string            `json:"contact_number"`
	FacebookUserID          string            `json:"facebook_user_id"`
	GoogleUserID            string            `json:"google_user_id"`
	Locale                  string            `json:"locale"`
	Timezone                map[string]string `json:"timezone"`
	TermsOfServiceId        string            `db:"-" json:"terms_of_service_id,omitempty"`
	TermsOfServiceCreatedAt int64             `db:"-" json:"terms_of_service_created_at,omitempty"`
}

type UserPatch struct {
	Username  *string           `json:"username"`
	Password  *string           `json:"password,omitempty"`
	FirstName *string           `json:"first_name"`
	LastName  *string           `json:"last_name"`
	Email     *string           `json:"email"`
	Locale    *string           `json:"locale"`
	Timezone  map[string]string `json:"timezone"`
	RemoteId  *string           `json:"remote_id"`
}

type UserAuth struct {
	AuthData    *string `json:"auth_data,omitempty"`
	AuthService string  `json:"auth_service,omitempty"`
}

var validUsernameChars = regexp.MustCompile(`^[a-z0-9\.\-_]+$`)

// IsValid validates the user and returns an error if it isn't configured correctly.
func (u *User) IsValid() *AppError {

	if !IsValidId(u.ID) {
		return InvalidUserError("id", "")
	}

	if u.CreatedAt == 0 {
		return InvalidUserError("create_at", u.ID)
	}

	if u.UpdatedAt == 0 {
		return InvalidUserError("update_at", u.ID)
	}

	if len(u.Email) > USER_EMAIL_MAX_LENGTH || u.Email == "" || !IsValidEmail(u.Email) {
		return InvalidUserError("email", u.ID)
	}

	if utf8.RuneCountInString(u.FirstName) > USER_FIRST_NAME_MAX_RUNES {
		return InvalidUserError("first_name", u.ID)
	}

	if utf8.RuneCountInString(u.LastName) > USER_LAST_NAME_MAX_RUNES {
		return InvalidUserError("last_name", u.ID)
	}

	if u.AuthData != nil && *u.AuthData != "" && u.AuthService == "" {
		return InvalidUserError("auth_data_type", u.ID)
	}

	if u.Password != "" && u.AuthData != nil && *u.AuthData != "" {
		return InvalidUserError("auth_data_pwd", u.ID)
	}

	if len(u.Password) > USER_PASSWORD_MAX_LENGTH {
		return InvalidUserError("password_limit", u.ID)
	}

	if !IsValidLocale(u.Locale) {
		return InvalidUserError("locale", u.ID)
	}

	if len(u.Timezone) > 0 {
		if tzJSON, err := json.Marshal(u.Timezone); err != nil {
			return NewAppError("User.IsValid", "model.user.is_valid.marshal.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else if utf8.RuneCount(tzJSON) > USER_TIMEZONE_MAX_RUNES {
			return InvalidUserError("timezone_limit", u.ID)
		}
	}

	if len(u.ContactNumber) > USER_CONTACT_NUMBER_MAX_LENGTH {
		return InvalidUserError("contact_number_limit", u.ID)
	}

	return nil
}

func InvalidUserError(fieldName string, userId string) *AppError {
	id := fmt.Sprintf("model.user.is_valid.%s.app_error", fieldName)
	details := ""
	if userId != "" {
		details = "user_id=" + userId
	}
	return NewAppError("User.IsValid", id, nil, details, http.StatusBadRequest)
}

func NormalizeUsername(username string) string {
	return strings.ToLower(username)
}

func NormalizeEmail(email string) string {
	return strings.ToLower(email)
}

// PreSave will set the Id and Username if missing.  It will also fill
// in the CreatedAt, UpdatedAt times.  It will also hash the password.  It should
// be run before saving the user to the db.
func (u *User) PreSave() {
	if u.ID == "" {
		u.ID = NewId()
	}

	if u.Username == "" {
		u.Username = NewId()
	}

	if u.AuthData != nil && *u.AuthData == "" {
		u.AuthData = nil
	}

	u.Username = SanitizeUnicode(u.Username)
	u.FirstName = SanitizeUnicode(u.FirstName)
	u.LastName = SanitizeUnicode(u.LastName)

	u.Username = NormalizeUsername(u.Username)
	u.Email = NormalizeEmail(u.Email)

	u.CreatedAt = GetMillis()
	u.UpdatedAt = u.CreatedAt

	if u.Locale == "" {
		u.Locale = DEFAULT_LOCALE
	}

	if u.Timezone == nil {
		u.Timezone = timezones.DefaultUserTimezone()
	}

	if u.Password != "" {
		u.Password = HashPassword(u.Password)
	}
}

// PreUpdate should be run before updating the user in the db.
func (u *User) PreUpdate() {
	u.Username = SanitizeUnicode(u.Username)
	u.FirstName = SanitizeUnicode(u.FirstName)
	u.LastName = SanitizeUnicode(u.LastName)

	u.Username = NormalizeUsername(u.Username)
	u.Email = NormalizeEmail(u.Email)
	u.UpdatedAt = GetMillis()

	u.FirstName = SanitizeUnicode(u.FirstName)
	u.LastName = SanitizeUnicode(u.LastName)

	if u.AuthData != nil && *u.AuthData == "" {
		u.AuthData = nil
	}

}

func (u *User) Patch(patch *UserPatch) {
	if patch.Username != nil {
		u.Username = *patch.Username
	}

	if patch.FirstName != nil {
		u.FirstName = *patch.FirstName
	}

	if patch.LastName != nil {
		u.LastName = *patch.LastName
	}

	if patch.Email != nil {
		u.Email = *patch.Email
	}

	if patch.Locale != nil {
		u.Locale = *patch.Locale
	}

	if patch.Timezone != nil {
		u.Timezone = patch.Timezone
	}

}

// Generate a valid strong etag so the browser can cache the results
func (u *User) Etag(showFullName, showEmail bool) string {
	return Etag(u.ID, u.UpdatedAt, u.TermsOfServiceId, u.TermsOfServiceCreatedAt, showFullName, showEmail)
}

// Remove any private data from the user object
func (u *User) Sanitize(options map[string]bool) {
	u.Password = ""
	u.AuthData = NewString("")

	if len(options) != 0 && !options["email"] {
		u.Email = ""
	}
	if len(options) != 0 && !options["fullname"] {
		u.FirstName = ""
		u.LastName = ""
	}

	if len(options) != 0 && !options["authservice"] {
		u.AuthService = ""
	}
}

// Remove any input data from the user object that is not user controlled
func (u *User) SanitizeInput(isAdmin bool) {
	if !isAdmin {
		u.AuthData = NewString("")
		u.AuthService = ""
		u.EmailVerified = false
	}
}

func (u *User) SanitizeProfile(options map[string]bool) {

	u.Sanitize(options)
}

func (u *User) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	} else if u.FirstName != "" {
		return u.FirstName
	} else if u.LastName != "" {
		return u.LastName
	} else {
		return ""
	}
}

func (u *User) getDisplayName(baseName, nameFormat string) string {
	displayName := baseName

	if fullName := u.GetFullName(); fullName != "" {
		displayName = fullName
	}

	return displayName
}

func (u *User) GetDisplayName(nameFormat string) string {
	displayName := u.Username

	return u.getDisplayName(displayName, nameFormat)
}

func (u *User) ToPatch() *UserPatch {
	return &UserPatch{
		Username:  &u.Username,
		Password:  &u.Password,
		FirstName: &u.FirstName,
		LastName:  &u.LastName,
		Email:     &u.Email,
		Locale:    &u.Locale,
		Timezone:  u.Timezone,
	}
}

func (u *UserPatch) SetField(fieldName string, fieldValue string) {
	switch fieldName {
	case "FirstName":
		u.FirstName = &fieldValue
	case "LastName":
		u.LastName = &fieldValue
	case "Email":
		u.Email = &fieldValue
	}
}

// UserFromJson will decode the input and return a User
func UserFromJson(data io.Reader) *User {
	var user *User
	json.NewDecoder(data).Decode(&user)
	return user
}

func UserAuthFromJson(data io.Reader) *UserAuth {
	var user *UserAuth
	json.NewDecoder(data).Decode(&user)
	return user
}

// HashPassword generates a hash using the bcrypt.GenerateFromPassword
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}

	return string(hash)
}
