package auth

import (
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/models"
	"FullStackOfYear/backend/types/auth"
	"context"
	"crypto/rand"
	"encoding/base64"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type service struct {
	users         *mongo.Collection
	tokenGen      TokenGenerator
	emailSender   EmailSender
	oauthProvider OAuthProvider
}

func NewService(
	users *mongo.Collection,
	tokenGen TokenGenerator,
	emailSender EmailSender,
	oauthProvider OAuthProvider,
) Service {
	return &service{
		users:         users,
		tokenGen:      tokenGen,
		emailSender:   emailSender,
		oauthProvider: oauthProvider,
	}
}

func (s *service) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	var user *models.User
	var err error

	switch req.LoginType {
	case auth.EmailLogin:
		user, err = s.handleEmailLogin(ctx, req)
	case auth.GoogleLogin:
		user, err = s.handleGoogleLogin(ctx, req)
	default:
		return nil, errors.NewAppError(errors.BadRequest, "Unsupported login type")
	}

	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.GenerateTokenPair(ctx, user.ID.Hex(), user.Role.Type, req.DeviceID)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		ExpiresIn:      3600,
		TokenType:      "Bearer",
		UserID:         user.ID.Hex(),
		Email:          user.Email,
		Username:       user.Username,
		CreatedAt:      user.CreatedAt,
		OAuthConnected: user.OAuth != nil && user.OAuth.Google != nil && user.OAuth.Google.Connected,
	}, nil
}

func (s *service) handleEmailLogin(ctx context.Context, req *auth.LoginRequest) (*models.User, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.NewAppError(errors.BadRequest, "Email and password required")
	}

	user, err := s.ValidateEmailPassword(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) handleGoogleLogin(ctx context.Context, req *auth.LoginRequest) (*models.User, error) {
	if req.Code == "" || req.State == "" {
		return nil, errors.NewAppError(errors.BadRequest, "OAuth code and state required")
	}

	token, err := s.oauthProvider.ExchangeCode(ctx, req.Code)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to exchange OAuth code")
	}

	userInfo, err := s.oauthProvider.GetUserInfo(ctx, token)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to get user info")
	}

	return s.findOrCreateGoogleUser(ctx, userInfo)
}

func (s *service) ValidateEmailPassword(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.GetErrorCode(err) == errors.NotFound {
			// 用户不存在时延迟返回，防止时间泄露
			time.Sleep(100 * time.Millisecond)
			return nil, errors.NewAppError(errors.Unauthorized, "邮箱或密码错误")
		}
		return nil, err
	}

	// 检查邮箱验证状态
	if !user.Status.EmailVerified {
		return nil, errors.NewAppError(errors.Unauthorized, "请先验证邮箱")
	}

	// 检查账户锁定状态
	if user.Status.IsLocked {
		if user.Status.LockExpires.After(time.Now()) {
			return nil, errors.NewAppError(errors.Forbidden,
				"账户已锁定，请在 "+user.Status.LockExpires.Sub(time.Now()).String()+" 后重试")
		}
		// 自动解锁
		user.Status.IsLocked = false
		user.Status.LockReason = ""
		user.Status.LockExpires = time.Time{}

		_, err = s.users.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
			"$set": bson.M{
				"status.isLocked":    false,
				"status.lockReason":  "",
				"status.lockExpires": time.Time{},
			},
		})
		if err != nil {
			return nil, errors.NewAppError(errors.InternalError, "更新用户状态失败")
		}
	}

	// 验证密码
	if !CheckPasswordHash(password, user.Password) {
		// 更新登录失败次数
		failedAttempts := user.Stats.FailedLoginAttempts + 1
		update := bson.M{
			"$inc": bson.M{"stats.failedLoginAttempts": 1},
		}

		// 如果失败次数过多，锁定账户
		if failedAttempts >= 5 {
			lockExpires := time.Now().Add(30 * time.Minute)
			update["$set"] = bson.M{
				"status.isLocked":    true,
				"status.lockReason":  "登录失败次数过多",
				"status.lockExpires": lockExpires,
			}
		}

		_, err = s.users.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
		if err != nil {
			return nil, errors.NewAppError(errors.InternalError, "更新用户状态失败")
		}

		return nil, errors.NewAppError(errors.Unauthorized, "邮箱或密码错误")
	}

	// 重置登录失败次数并更新最后登录时间
	_, err = s.users.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"stats.failedLoginAttempts": 0,
			"stats.lastLoginAt":         time.Now(),
		},
	})
	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "更新用户状态失败")
	}

	return user, nil
}

func (s *service) HandleOAuthLogin(ctx context.Context, userInfo *auth.OAuthUserInfo) (*auth.LoginResponse, error) {
	user, err := s.findOrCreateGoogleUser(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	//accessToken, refreshToken, err := s.tokenGen.GenerateTokenPair(
	//	user.ID.Hex(),
	//	user.Role.Type,
	//	"oauth_"+user.ID.Hex(),
	//)
	deviceID := "oauth_device_" + user.ID.Hex()

	accessToken, refreshToken, err := s.tokenGen.GenerateTokenPair(
		user.ID.Hex(),
		user.Role.Type,
		deviceID,
	)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to generate tokens")
	}

	return &auth.LoginResponse{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		ExpiresIn:      3600,
		TokenType:      "Bearer",
		UserID:         user.ID.Hex(),
		Email:          user.Email,
		Username:       user.Username,
		CreatedAt:      user.CreatedAt,
		OAuthConnected: true,
		OAuthProvider:  "google",
	}, nil
}

func (s *service) findOrCreateGoogleUser(ctx context.Context, userInfo *auth.OAuthUserInfo) (*models.User, error) {
	user, err := s.findUserByGoogleID(ctx, userInfo.ID)
	if err == nil {
		return s.updateGoogleUser(ctx, user, userInfo)
	}

	user, err = s.GetUserByEmail(ctx, userInfo.Email)
	if err == nil {
		return s.linkGoogleAccount(ctx, user, userInfo)
	}

	return s.createGoogleUser(ctx, userInfo)
}

func (s *service) findUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	var user models.User
	err := s.users.FindOne(ctx, bson.M{"oauth.google.id": googleID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewAppError(errors.NotFound, "User not found")
		}
		return nil, errors.NewAppError(errors.InternalError, "Database error")
	}
	return &user, nil
}

func (s *service) updateGoogleUser(ctx context.Context, user *models.User, info *auth.OAuthUserInfo) (*models.User, error) {
	now := time.Now()

	googleOAuth := &models.GoogleOAuth{
		ID:          info.ID,
		Email:       info.Email,
		Connected:   true,
		ConnectedAt: now,
	}

	user.UpdateOAuthInfo("google", googleOAuth)

	_, err := s.users.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to update user")
	}

	return user, nil
}

func (s *service) linkGoogleAccount(ctx context.Context, user *models.User, info *auth.OAuthUserInfo) (*models.User, error) {
	if !user.Status.EmailVerified {
		return nil, errors.NewAppError(errors.BadRequest, "Email must be verified before linking Google account")
	}

	now := time.Now()
	googleOAuth := &models.GoogleOAuth{
		ID:          info.ID,
		Email:       info.Email,
		Connected:   true,
		ConnectedAt: now,
	}

	user.UpdateOAuthInfo("google", googleOAuth)

	user.AddSecurityLog(models.SecurityLog{
		Action:      "google_account_linked",
		Timestamp:   now,
		Description: "Google account linked to existing email account",
	})

	_, err := s.users.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to link Google account")
	}

	return user, nil
}

func (s *service) createGoogleUser(ctx context.Context, info *auth.OAuthUserInfo) (*models.User, error) {
	now := time.Now()

	user := &models.User{
		ID:       primitive.NewObjectID(),
		Username: info.Email,
		Email:    info.Email,
		Status: models.Status{
			EmailVerified: true,
		},
		Role: models.Role{
			Type: models.RoleUser,
		},
		Stats: models.UserStats{
			CreatedAt:   now,
			LastLoginAt: now,
		},
		OAuth: &models.OAuthInfo{
			Google: &models.GoogleOAuth{
				ID:          info.ID,
				Email:       info.Email,
				Connected:   true,
				ConnectedAt: now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := s.users.InsertOne(ctx, user)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to create user")
	}

	return user, nil
}

func (s *service) SendVerificationEmail(ctx context.Context, userID string) error {
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Status.EmailVerified {
		return errors.NewAppError(errors.BadRequest, "邮箱已验证")
	}

	verifyToken, err := s.GenerateEmailVerificationToken(ctx, userID)
	if err != nil {
		return err
	}

	// URL安全的Base64编码
	//encodedToken := base64.URLEncoding.EncodeToString([]byte(verifyToken))
	//
	//// 生成完整的验证链接
	//verifyLink := fmt.Sprintf("%s/verify-email?token=%s",
	//	"http://localhost:3000", // TODO: 从配置中获取
	//	encodedToken,
	//)

	// 直接发送验证token
	return s.emailSender.SendVerificationEmail(user.Email, user.Username, verifyToken)
}

func (s *service) VerifyEmail(ctx context.Context, encodedToken string) error {
	// 解码token
	tokenBytes, err := base64.URLEncoding.DecodeString(encodedToken)
	if err != nil {
		return errors.NewAppError(errors.BadRequest, "无效的验证链接")
	}
	verifyToken := string(tokenBytes)

	// 使用事务确保原子性
	session, err := s.users.Database().Client().StartSession()
	if err != nil {
		return errors.NewAppError(errors.InternalError, "数据库会话创建失败")
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		update := bson.M{
			"$set": bson.M{
				"status.emailVerified": true,
				"status.verifyToken":   "",
				"status.tokenExpires":  time.Time{},
			},
		}

		result := s.users.FindOneAndUpdate(
			sessCtx,
			bson.M{
				"status.verifyToken": verifyToken,
				"status.tokenExpires": bson.M{
					"$gt": time.Now(),
				},
			},
			update,
		)

		if result.Err() != nil {
			if result.Err() == mongo.ErrNoDocuments {
				return nil, errors.NewAppError(errors.BadRequest, "验证链接无效或已过期")
			}
			return nil, errors.NewAppError(errors.InternalError, "验证邮箱失败")
		}

		return nil, nil
	})

	return err
}

func (s *service) GenerateTokenPair(ctx context.Context, userID, role, deviceID string) (string, string, error) {
	return s.tokenGen.GenerateTokenPair(userID, role, deviceID)
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.tokenGen.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	return s.GenerateTokenPair(ctx, claims.UserID, claims.Role, claims.DeviceID)
}

func (s *service) RevokeTokens(ctx context.Context, userID, deviceID string) error {
	return s.tokenGen.RevokeTokens(userID, deviceID)
}

func (s *service) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewAppError(errors.BadRequest, "Invalid user ID")
	}

	var user models.User
	err = s.users.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewAppError(errors.NotFound, "User not found")
		}
		return nil, errors.NewAppError(errors.InternalError, "Database error")
	}
	return &user, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NewAppError(errors.NotFound, "User not found")
		}
		return nil, errors.NewAppError(errors.InternalError, "Database error")
	}
	return &user, nil
}

func (s *service) GenerateEmailVerificationToken(ctx context.Context, userID string) (string, error) {
	token := generateSecureToken()
	expires := time.Now().Add(24 * time.Hour)

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", errors.NewAppError(errors.BadRequest, "Invalid user ID")
	}

	update := bson.M{
		"$set": bson.M{
			"status.verifyToken":  token,
			"status.tokenExpires": expires,
		},
	}

	_, err = s.users.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return "", errors.NewAppError(errors.InternalError, "Failed to update verification token")
	}

	return token, nil
}

func (s *service) GenerateEmailChangeToken(user *models.User, newEmail string) (string, error) {
	token := generateSecureToken()
	expires := time.Now().Add(24 * time.Hour)

	update := bson.M{
		"$set": bson.M{
			"status.emailChange":  newEmail,
			"status.verifyToken":  token,
			"status.tokenExpires": expires,
		},
	}

	_, err := s.users.UpdateOne(context.Background(), bson.M{"_id": user.ID}, update)
	if err != nil {
		return "", errors.NewAppError(errors.InternalError, "Failed to generate email change token")
	}

	return token, nil
}

func generateSecureToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
