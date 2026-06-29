package profiles

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	repo *Repository
	r2   *R2Uploader
}

func NewService(repo *Repository, r2 *R2Uploader) *Service {
	return &Service{
		repo: repo,
		r2:   r2,
	}
}

func (s *Service) EnsureProfile(ctx context.Context, ownerID string) (*Profile, error) {
	existing, err := s.repo.FindProfileByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return existing, nil
	}

	user, err := s.repo.FindUserByID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user account not found")
	}

	email := GetUserEmail(user)
	baseUsername := user.Username

	if baseUsername == "" && email != "" {
		baseUsername = strings.Split(email, "@")[0]
	}

	if baseUsername == "" {
		baseUsername = "user-" + ownerID[:min(6, len(ownerID))]
	}

	username := NormalizeUsername(baseUsername)
	slug := Slugify(username)
	now := time.Now()

	displayName := username
	if value, ok := user.Profile["displayName"].(string); ok && strings.TrimSpace(value) != "" {
		displayName = strings.TrimSpace(value)
	} else if value, ok := user.Profile["name"].(string); ok && strings.TrimSpace(value) != "" {
		displayName = strings.TrimSpace(value)
	} else if value, ok := user.Profile["fullName"].(string); ok && strings.TrimSpace(value) != "" {
		displayName = strings.TrimSpace(value)
	}

	profile := &Profile{
		OwnerID:      ownerID,
		Username:     username,
		Slug:         slug,
		Email:        email,
		Phone:        "",
		Language:     "en",
		DisplayName:  displayName,
		Bio:          "",
		SocialLinks:  []SocialLink{},
		Status:       "active",
		IsBlocked:    false,
		Custom:       map[string]interface{}{},
		CreatedAt:    now,
		UpdatedAt:    now,
		Privacy:      DefaultPrivacy(),
		Consent:      DefaultConsent(now),
		Verification: Verification{IsVerified: len(user.Emails) > 0 && user.Emails[0].Verified},
	}

	err = s.repo.InsertProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func DefaultPrivacy() Privacy {
	return Privacy{
		ShowEmail:       false,
		ShowPhone:       false,
		ShowLocation:    true,
		ShowSocialLinks: true,
		ShowBio:         true,
		ShowCustom:      false,
	}
}

func DefaultConsent(now time.Time) Consent {
	return Consent{
		Marketing: false,
		Email:     false,
		UpdatedAt: &now,
	}
}

func (s *Service) GetMe(ctx context.Context, ownerID string) (*Profile, error) {
	err := s.ensureOwner(ownerID)
	if err != nil {
		return nil, err
	}

	return s.EnsureProfile(ctx, ownerID)
}

func (s *Service) UpdateProfile(ctx context.Context, ownerID string, req UpdateProfileRequest) (*Profile, error) {
	err := s.ensureOwner(ownerID)
	if err != nil {
		return nil, err
	}

	_, err = s.EnsureProfile(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	set := bson.M{}

	if req.Username != nil {
		username := NormalizeUsername(*req.Username)

		if err := ValidateUsername(username); err != nil {
			return nil, err
		}

		set["username"] = username
		set["slug"] = Slugify(username)
	}

	if req.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*req.Email))

		if err := ValidateEmail(email); err != nil {
			return nil, err
		}

		existingUser, err := s.repo.FindUserByEmail(ctx, email)
		if err != nil {
			return nil, err
		}

		if existingUser != nil && existingUser.ID != ownerID {
			return nil, errors.New("this email is already in use")
		}

		set["email"] = email
		set["verification.isVerified"] = false

		err = s.repo.UpdateUser(ctx, ownerID, bson.M{
			"emails": []UserEmail{
				{
					Address:  email,
					Verified: false,
				},
			},
		})
		if err != nil {
			return nil, err
		}

		_ = s.repo.CreateAuditLog(ctx, "profile.updateEmail", bson.M{
			"ownerId": ownerID,
			"newEmail": email,
		})
	}

	if req.Phone != nil {
		set["phone"] = strings.TrimSpace(*req.Phone)
	}

	if req.Language != nil {
		lang := strings.TrimSpace(*req.Language)
		if lang != "en" && lang != "fr" && lang != "es" {
			return nil, errors.New("unsupported language")
		}
		set["language"] = lang
	}

	if req.DisplayName != nil {
		displayName := strings.TrimSpace(*req.DisplayName)
		set["displayName"] = displayName

		err = s.repo.UpdateUser(ctx, ownerID, bson.M{
			"profile.displayName": displayName,
			"profile.name":        displayName,
			"profile.fullName":    displayName,
			"profile.updatedAt":   time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}

	if req.Bio != nil {
		if len(*req.Bio) > 1000 {
			return nil, errors.New("bio is too long")
		}
		set["bio"] = strings.TrimSpace(*req.Bio)
	}

	if req.Location != nil {
		set["location"] = req.Location
	}

	if req.SocialLinks != nil {
		if len(*req.SocialLinks) > 10 {
			return nil, errors.New("maximum 10 social links allowed")
		}
		set["socialLinks"] = *req.SocialLinks
	}

	if req.Custom != nil {
		set["custom"] = req.Custom
	}

	if len(set) == 0 {
		return s.repo.FindProfileByOwnerID(ctx, ownerID)
	}

	err = s.repo.UpdateProfileByOwnerID(ctx, ownerID, set)
	if err != nil {
		return nil, err
	}

	return s.repo.FindProfileByOwnerID(ctx, ownerID)
}

func (s *Service) UnsetCustomKey(ctx context.Context, ownerID string, key string) error {
	safeKey, err := SanitizeCustomKey(key)
	if err != nil {
		return err
	}

	_, err = s.EnsureProfile(ctx, ownerID)
	if err != nil {
		return err
	}

	return s.repo.UnsetProfileCustomKey(ctx, ownerID, safeKey)
}

func (s *Service) Touch(ctx context.Context, ownerID string) error {
	if ownerID == "" {
		return nil
	}

	return s.repo.UpdateProfileByOwnerID(ctx, ownerID, bson.M{
		"lastActiveAt": time.Now(),
	})
}

func (s *Service) UpdateAvatar(ctx context.Context, ownerID string, base64Image string) (*Profile, error) {
	if s.r2 == nil {
		return nil, errors.New("R2 uploader is not configured")
	}

	profile, err := s.EnsureProfile(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	media, err := s.r2.UploadImage(
		ctx,
		base64Image,
		fmt.Sprintf("blomiva/profiles/%s/profile-image", ownerID),
		"profile-image",
	)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpdateProfileByOwnerID(ctx, ownerID, bson.M{
		"profileImage": media,
		"avatar":       media,
		"avatarUrl":    media.PublicURL,
	})
	if err != nil {
		_ = s.r2.DeleteImage(ctx, media)
		return nil, err
	}

	err = s.repo.UpdateUser(ctx, ownerID, bson.M{
		"profile.profileImage": media,
		"profile.avatar":       media,
		"profile.avatarUrl":    media.PublicURL,
		"profile.updatedAt":    time.Now(),
	})
	if err != nil {
		return nil, err
	}

	if profile.ProfileImage != nil && profile.ProfileImage.StorageKey != media.StorageKey {
		_ = s.r2.DeleteImage(context.Background(), profile.ProfileImage)
	}

	return s.repo.FindProfileByOwnerID(ctx, ownerID)
}

func (s *Service) UpdateBanner(ctx context.Context, ownerID string, base64Image string) (*Profile, error) {
	if s.r2 == nil {
		return nil, errors.New("R2 uploader is not configured")
	}

	profile, err := s.EnsureProfile(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	media, err := s.r2.UploadImage(
		ctx,
		base64Image,
		fmt.Sprintf("blomiva/profiles/%s/banner", ownerID),
		"profile-banner",
	)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpdateProfileByOwnerID(ctx, ownerID, bson.M{
		"banner":    media,
		"bannerUrl": media.PublicURL,
	})
	if err != nil {
		_ = s.r2.DeleteImage(ctx, media)
		return nil, err
	}

	err = s.repo.UpdateUser(ctx, ownerID, bson.M{
		"profile.banner":    media,
		"profile.bannerUrl": media.PublicURL,
		"profile.updatedAt": time.Now(),
	})
	if err != nil {
		return nil, err
	}

	if profile.Banner != nil && profile.Banner.StorageKey != media.StorageKey {
		_ = s.r2.DeleteImage(context.Background(), profile.Banner)
	}

	return s.repo.FindProfileByOwnerID(ctx, ownerID)
}

func (s *Service) ScheduleDelete(ctx context.Context, ownerID string) error {
	_, err := s.EnsureProfile(ctx, ownerID)
	if err != nil {
		return err
	}

	deletionDate := time.Now().Add(30 * 24 * time.Hour)

	err = s.repo.UpdateProfileByOwnerID(ctx, ownerID, bson.M{
		"scheduledForDeletion": deletionDate,
		"status":               "inactive",
	})
	if err != nil {
		return err
	}

	return s.repo.CreateAuditLog(ctx, "profile.scheduleDelete", bson.M{
		"ownerId":      ownerID,
		"deletionDate": deletionDate,
	})
}

func (s *Service) SetUserStatus(ctx context.Context, adminID string, ownerID string, active bool) error {
	if ownerID == "" {
		return errors.New("user ID is required")
	}

	if adminID == ownerID {
		return errors.New("you cannot change your own status")
	}

	status := "inactive"
	if active {
		status = "active"
	}

	err := s.repo.UpdateProfileByOwnerID(ctx, ownerID, bson.M{
		"status":    status,
		"isBlocked": !active,
	})
	if err != nil {
		return err
	}

	return s.repo.CreateAuditLog(ctx, "admin.profiles.setUserStatus", bson.M{
		"adminId":       adminID,
		"targetOwnerId": ownerID,
		"active":        active,
	})
}

func (s *Service) DeleteUser(ctx context.Context, adminID string, ownerID string) error {
	if ownerID == "" {
		return errors.New("user ID is required")
	}

	if adminID == ownerID {
		return errors.New("you cannot delete your own account")
	}

	err := s.repo.DeleteUserAndProfile(ctx, ownerID)
	if err != nil {
		return err
	}

	return s.repo.CreateAuditLog(ctx, "admin.profiles.deleteUser", bson.M{
		"adminId":       adminID,
		"targetOwnerId": ownerID,
	})
}

func (s *Service) ensureOwner(ownerID string) error {
	if ownerID == "" {
		return errors.New("not authenticated")
	}
	return nil
}