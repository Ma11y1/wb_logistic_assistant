package models

import "errors"

type UserInfo struct {
	ID           int              `json:"id"`
	Verified     bool             `json:"verified"`
	RoleIDs      []string         `json:"role_ids"`
	Roles        []*UserInfoRole  `json:"roles"`
	Permissions  []string         `json:"permissions"`
	UserDetails  *UserInfoDetails `json:"user_details"`
	DriverRoleID int              `json:"driver_role_id"`
}

func (i *UserInfo) Validate() error {
	if i.ID == 0 {
		return errors.New("id cannot be 0")
	}
	err := i.UserDetails.Validate()
	if err != nil {
		return err
	}
	return nil
}

type UserInfoRole struct {
	UserRoleUID  string `json:"user_role_uid"`
	UserRoleName string `json:"user_role_name"`
}

type UserInfoDetails struct {
	Name                 string `json:"name"`
	PhoneNumber          string `json:"phone_number"`
	SupplierID           int    `json:"supplier_id"`
	FreelancerEmployeeID int    `json:"freelancer_employee_id"`
	VatID                int    `json:"vat_id"`
	VatName              string `json:"vat_name"`
	Telegram             string `json:"telegram"`
}

func (i *UserInfoDetails) Validate() error {
	if i.Name == "" {
		return errors.New("name is empty")
	}
	if i.PhoneNumber == "" {
		return errors.New("phone_number is empty")
	}
	if i.SupplierID == 0 {
		return errors.New("supplier_id is empty")
	}
	return nil
}
