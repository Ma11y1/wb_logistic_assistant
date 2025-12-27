package wb_logistic

import (
	wb_models "wb_logistic_assistant/external/wb_logistic_api/models"
	"wb_logistic_assistant/internal/errors"
	"wb_logistic_assistant/internal/models"
)

func (i *Initializer) GetLoginStorage() string {
	return i.storage.ConfigStore().GetWBLogisticLogin()
}

func (i *Initializer) SetLoginStorage(login string) {
	i.storage.ConfigStore().SetWBLogisticLogin(login)
}

func (i *Initializer) GetAccessTokenStorage() (*wb_models.AuthAccessToken, error) {
	tokenModel := i.storage.ConfigStore().GetWBLogisticAccessToken()
	if tokenModel == nil {
		return nil, errors.New("Initializer.GetAccessTokenStorage()", "access token not found in storage")
	}

	token := &wb_models.AuthAccessToken{
		TokenType:    tokenModel.GetTokenType(),
		AccessToken:  tokenModel.GetAccessToken(),
		RefreshToken: tokenModel.GetRefreshToken(),
		ExpiresIn:    tokenModel.GetExpiresIn(),
	}

	if err := token.Validate(); err != nil {
		return nil, errors.Wrap(err, "Initializer.GetAccessTokenStorage()", "invalid access token")
	}

	return token, nil
}

func (i *Initializer) SetAccessTokenStorage(token *wb_models.AuthAccessToken) error {
	if token == nil {
		return errors.New("Initializer.SetAccessTokenStorage()", "access token is nil")
	}
	err := token.Validate()
	if err != nil {
		return errors.Wrap(err, "Initializer.SetUserInfoStorage()", "invalid access token")
	}
	tokenModel := &models.WBLogisticAccessTokenModel{}
	tokenModel.SetAll(
		token.TokenType,
		token.AccessToken,
		token.RefreshToken,
		token.ExpiresIn,
	)
	i.storage.ConfigStore().SetWBLogisticAccessToken(tokenModel)
	return nil
}

func (i *Initializer) GetSessionTokenStorage() (*wb_models.AuthSessionToken, error) {
	tokenModel := i.storage.ConfigStore().GetWBLogisticSessionToken()
	if tokenModel == nil {
		return nil, errors.New("Initializer.GetSessionTokenStorage()", "session token not found in storage")
	}

	token := &wb_models.AuthSessionToken{
		Source: tokenModel.GetSource(),
		Token: wb_models.AuthSessionTokenData{
			TokenType:   tokenModel.GetTokenType(),
			AccessToken: tokenModel.GetAccessToken(),
			ExpiresIn:   tokenModel.GetExpiresIn(),
		},
	}

	if err := token.Validate(); err != nil {
		return nil, errors.Wrap(err, "Initializer.GetSessionTokenStorage()", "invalid session token")
	}

	return token, nil
}

func (i *Initializer) SetSessionTokenStorage(token *wb_models.AuthSessionToken) error {
	if token == nil {
		return errors.New("Initializer.SetSessionTokenStorage()", "session token is nil")
	}
	err := token.Validate()
	if err != nil {
		return errors.Wrap(err, "Initializer.SetUserInfoStorage()", "invalid session token")
	}
	tokenModel := &models.WBLogisticSessionTokenModel{}
	tokenModel.SetAll(
		token.Token.TokenType,
		token.Source,
		token.Token.AccessToken,
		token.Token.ExpiresIn,
	)
	i.storage.ConfigStore().SetWBLogisticSessionToken(tokenModel)
	return nil
}

func (i *Initializer) GetUserInfoStorage() (*wb_models.UserInfo, error) {
	infoModel := i.storage.ConfigStore().GetWBLogisticUserInfo()
	if infoModel == nil {
		return nil, errors.New("Initializer.GetUserInfoStorage()", "user info not found in storage")
	}

	roleIDs := make([]string, len(infoModel.GetRoleIDs()))
	for j, roleID := range infoModel.GetRoleIDs() {
		roleIDs[j] = roleID
	}

	roles := make([]*wb_models.UserInfoRole, len(infoModel.GetRoles()))
	for j, roleModel := range infoModel.GetRoles() {
		role := &wb_models.UserInfoRole{
			UserRoleUID:  roleModel.GetUserRoleUID(),
			UserRoleName: roleModel.GetUserRoleName(),
		}
		roles[j] = role
	}

	userDetailsModel := infoModel.GetUserDetails()
	userDetails := &wb_models.UserInfoDetails{
		Name:                 userDetailsModel.GetName(),
		PhoneNumber:          userDetailsModel.GetPhoneNumber(),
		SupplierID:           userDetailsModel.GetSupplierID(),
		FreelancerEmployeeID: userDetailsModel.GetFreelancerEmployeeID(),
		VatID:                userDetailsModel.GetVatID(),
		VatName:              userDetailsModel.GetVatName(),
		Telegram:             userDetailsModel.GetTelegram(),
	}

	permissions := make([]string, len(infoModel.GetPermissions()))
	for j, permission := range infoModel.GetPermissions() {
		permissions[j] = permission
	}

	info := &wb_models.UserInfo{
		ID:           infoModel.GetID(),
		Verified:     infoModel.IsVerified(),
		RoleIDs:      roleIDs,
		Roles:        roles,
		Permissions:  permissions,
		UserDetails:  userDetails,
		DriverRoleID: infoModel.GetDriverRoleID(),
	}

	if err := info.Validate(); err != nil {
		return nil, errors.Wrap(err, "Initializer.GetUserInfoStorage()", "invalid user info")
	}

	return info, nil
}

func (i *Initializer) SetUserInfoStorage(info *wb_models.UserInfo) error {
	if info == nil {
		return errors.New("Initializer.SetUserInfoStorage()", "user info is nil")
	}
	err := info.Validate()
	if err != nil {
		return errors.Wrap(err, "Initializer.SetUserInfoStorage()", "invalid user info")
	}

	infoModel := &models.WBLogisticUserInfoModel{}

	roleIDs := make([]string, len(info.RoleIDs))
	for j, roleID := range info.RoleIDs {
		roleIDs[j] = roleID
	}

	roles := make([]*models.WBLogisticUserInfoRoleModel, len(info.Roles))
	for j, role := range info.Roles {
		roleModel := &models.WBLogisticUserInfoRoleModel{}
		roleModel.SetAll(role.UserRoleUID, role.UserRoleName)
		roles[j] = roleModel
	}

	permissions := make([]string, len(info.Permissions))
	for j, permission := range info.Permissions {
		permissions[j] = permission
	}

	userDetails := &models.WBLogisticUserInfoDetailsModel{}
	userDetails.SetAll(
		info.UserDetails.Name,
		info.UserDetails.PhoneNumber,
		info.UserDetails.SupplierID,
		info.UserDetails.FreelancerEmployeeID,
		info.UserDetails.VatID,
		info.UserDetails.VatName,
		info.UserDetails.Telegram,
	)

	infoModel.SetAll(
		info.ID,
		info.Verified,
		roleIDs,
		roles,
		permissions,
		userDetails,
		info.DriverRoleID,
	)

	i.storage.ConfigStore().SetWBLogisticUserInfo(infoModel)

	return nil
}

func (i *Initializer) UpdateStorage() error {
	err := i.storage.Save(i.config.Storage().Path())
	if err != nil {
		return errors.Wrap(err, "Initializer.UpdateStorage()", "failed to update storage")
	}
	return nil
}
