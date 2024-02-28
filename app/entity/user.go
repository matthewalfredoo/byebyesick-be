package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type User struct {
	Id            int64        `json:"id"`
	Email         string       `json:"email"`
	Password      string       `json:"password"`
	UserRoleId    int64        `json:"user_role_id"`
	IsVerified    bool         `json:"is_verified"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
	DeletedAt     sql.NullTime `json:"deleted_at"`
	UserProfile   *UserProfile
	DoctorProfile *DoctorProfile
}

func (u *User) GetEntityName() string {
	return "users"
}

func (u *User) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(u).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (u *User) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", u.GetEntityName(), u.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (u *User) ToUserResponse() *responsedto.UserResponse {
	return &responsedto.UserResponse{
		Id:         u.Id,
		Email:      u.Email,
		UserRoleId: u.UserRoleId,
		IsVerified: u.IsVerified,
	}
}

func (u *User) ToUserProfileResponse() *responsedto.UserProfileResponse {
	return &responsedto.UserProfileResponse{
		Id:           u.Id,
		Email:        u.Email,
		UserRoleId:   u.UserRoleId,
		IsVerified:   u.IsVerified,
		Name:         u.UserProfile.Name,
		ProfilePhoto: u.UserProfile.ProfilePhoto,
		DateOfBirth:  u.UserProfile.DateOfBirth.Format(appconstant.TimeFormatQueryParam),
	}
}

func (u *User) ToDoctorProfileResponse() *responsedto.DoctorProfileResponse {
	return &responsedto.DoctorProfileResponse{
		Id:                u.Id,
		Email:             u.Email,
		UserRoleID:        u.UserRoleId,
		IsVerified:        u.IsVerified,
		Name:              u.DoctorProfile.Name,
		ProfilePhoto:      u.DoctorProfile.ProfilePhoto,
		StartingYear:      u.DoctorProfile.StartingYear,
		DoctorCertificate: u.DoctorProfile.DoctorCertificate,
		DoctorSpecialization: &responsedto.DoctorSpecializationResponse{
			Id:   u.DoctorProfile.DoctorSpecialization.Id,
			Name: u.DoctorProfile.DoctorSpecialization.Name,
		},
		ConsultationFee: u.DoctorProfile.ConsultationFee.String(),
		IsOnline:        u.DoctorProfile.IsOnline,
	}
}

func (u *User) ToPrescriptionSickLeaveUserProfileResponse() *responsedto.PrescriptionSickLeaveUserProfileResponse {
	if u == nil {
		return nil
	}
	if u.UserProfile != nil {
		var dateOfBirth string

		if !u.UserProfile.DateOfBirth.IsZero() {
			dateOfBirth = u.UserProfile.DateOfBirth.Format(time.RFC3339)
		}

		return &responsedto.PrescriptionSickLeaveUserProfileResponse{
			Email:       u.Email,
			Name:        u.UserProfile.Name,
			DateOfBirth: dateOfBirth,
		}
	}
	if u.DoctorProfile != nil {
		return &responsedto.PrescriptionSickLeaveUserProfileResponse{
			Email:                u.Email,
			Name:                 u.DoctorProfile.Name,
			DoctorSpecialization: u.DoctorProfile.DoctorSpecialization.Name,
		}
	}
	return nil
}

func (u *User) GetProfile() *Profile {
	if u.UserProfile != nil {
		return &Profile{
			UserId:       u.Id,
			RoleId:       u.UserRoleId,
			Name:         u.UserProfile.Name,
			ProfilePhoto: u.UserProfile.ProfilePhoto,
		}
	}
	if u.DoctorProfile != nil {
		return &Profile{
			UserId:       u.Id,
			RoleId:       u.UserRoleId,
			Name:         u.DoctorProfile.Name,
			ProfilePhoto: u.DoctorProfile.ProfilePhoto,
		}
	}
	return nil
}
