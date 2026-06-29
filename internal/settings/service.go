package settings

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	col   *mongo.Collection
	users *mongo.Collection
	mail  EmailSender
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		col:   db.Collection("user_settings"),
		users: db.Collection("users"),
		mail:  NewPostmarkEmailSender(),
	}
}

func defaultSettings(userID string) UserSettings {
	now := time.Now().UTC()

	return UserSettings{
		UserID:   userID,
		Theme:    "system",
		Language: "en",
		Timezone: "UTC",
		Notifications: NotificationSettings{
			Email:               true,
			Push:                true,
			SMS:                 false,
			SchoolUpdates:       true,
			BillingAlerts:       true,
			SecurityAlerts:      true,
			MarketingEmails:     false,
			ParentUpdates:       true,
			AttendanceAlerts:    true,
			AssignmentAlerts:    true,
			FeePaymentReminders: true,
		},
		Privacy: PrivacySettings{
			ProfileVisibility: "private",
			ShowEmail:         false,
			ShowPhone:         false,
			ShowOnlineStatus:  false,
			AllowMessages:     true,
		},
		Security: SecuritySettings{
			TwoFactorEnabled: false,
			LoginAlerts:      true,
		},
		Product: ProductSettings{
			DefaultSchoolID: "",
			DefaultRole:     "",
			DashboardMode:   "school",
			CompactUI:       false,
			ReduceMotion:    false,
		},
		Account: AccountDeletion{
			IsDeleted:         false,
			DeletionGraceDays: 30,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s *Service) GetOrCreate(ctx context.Context, userID string) (*UserSettings, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	var result UserSettings
	err := s.col.FindOne(ctx, bson.M{"userId": userID}).Decode(&result)
	if err == nil {
		return &result, nil
	}

	if err != mongo.ErrNoDocuments {
		return nil, err
	}

	settings := defaultSettings(userID)

	_, err = s.col.InsertOne(ctx, settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func (s *Service) Update(ctx context.Context, userID string, req UpdateSettingsRequest) (*UserSettings, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	set := bson.M{"updatedAt": time.Now().UTC()}

	if req.Theme != nil {
		theme := cleanLower(*req.Theme)
		if theme != "light" && theme != "dark" && theme != "system" {
			return nil, errors.New("invalid theme")
		}
		set["theme"] = theme
	}

	if req.Language != nil {
		language := cleanLower(*req.Language)
		if language != "en" && language != "fr" && language != "es" {
			return nil, errors.New("invalid language")
		}
		set["language"] = language
	}

	if req.Timezone != nil {
		timezone := strings.TrimSpace(*req.Timezone)
		if timezone == "" {
			return nil, errors.New("timezone cannot be empty")
		}
		set["timezone"] = timezone
	}

	if req.Notifications != nil {
		set["notifications"] = req.Notifications
	}

	if req.Privacy != nil {
		visibility := cleanLower(req.Privacy.ProfileVisibility)
		if visibility == "" {
			visibility = "private"
		}
		if visibility != "public" && visibility != "school" && visibility != "private" {
			return nil, errors.New("invalid profile visibility")
		}
		req.Privacy.ProfileVisibility = visibility
		set["privacy"] = req.Privacy
	}

	if req.Security != nil {
		set["security"] = req.Security
	}

	if req.Product != nil {
		mode := cleanLower(req.Product.DashboardMode)
		if mode == "" {
			mode = "school"
		}
		if mode != "school" && mode != "studio" && mode != "market" {
			return nil, errors.New("invalid dashboard mode")
		}
		req.Product.DashboardMode = mode
		set["product"] = req.Product
	}

	defaultDoc := defaultSettings(userID)

	update := bson.M{
		"$setOnInsert": bson.M{
			"userId":        defaultDoc.UserID,
			"createdAt":     defaultDoc.CreatedAt,
			"theme":         defaultDoc.Theme,
			"language":      defaultDoc.Language,
			"timezone":      defaultDoc.Timezone,
			"notifications": defaultDoc.Notifications,
			"privacy":       defaultDoc.Privacy,
			"security":      defaultDoc.Security,
			"product":       defaultDoc.Product,
			"account":       defaultDoc.Account,
		},
		"$set": set,
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result UserSettings
	err := s.col.FindOneAndUpdate(ctx, bson.M{"userId": userID}, update, opts).Decode(&result)
	if err != nil {
		return nil, err
	}

	s.sendUserEmail(ctx, userID, "Your Blomiva settings were updated", "Your account settings were updated successfully.")

	return &result, nil
}

func (s *Service) RequestAccountDeletion(ctx context.Context, userID string, reason string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("user id is required")
	}

	now := time.Now().UTC()
	scheduled := now.AddDate(0, 0, 30)

	_, err := s.col.UpdateOne(
		ctx,
		bson.M{"userId": userID},
		bson.M{
			"$set": bson.M{
				"account.isDeleted":         true,
				"account.deleteRequestedAt": now,
				"account.deleteScheduledAt": scheduled,
				"account.deletionReason":    strings.TrimSpace(reason),
				"account.deletionGraceDays": 30,
				"account.restoredAt":        nil,
				"updatedAt":                 now,
			},
			"$setOnInsert": bson.M{
				"userId":    userID,
				"createdAt": now,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	_, _ = s.users.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"isDeleted": true,
				"deletedAt": now,
				"updatedAt": now,
			},
		},
	)

	s.sendUserEmail(
		ctx,
		userID,
		"Your Blomiva account deletion was requested",
		"Your account has been scheduled for deletion. You have 30 days to restore it before permanent removal.",
	)

	return nil
}

func (s *Service) RestoreAccount(ctx context.Context, userID string, note string) error {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return errors.New("user id is required")
	}

	now := time.Now().UTC()

	_, err := s.col.UpdateOne(
		ctx,
		bson.M{"userId": userID},
		bson.M{
			"$set": bson.M{
				"account.isDeleted":          false,
				"account.restoredAt":         now,
				"account.deletionCancelNote": strings.TrimSpace(note),
				"updatedAt":                  now,
			},
			"$unset": bson.M{
				"account.deleteRequestedAt": "",
				"account.deleteScheduledAt": "",
				"account.deletionReason":    "",
			},
		},
	)
	if err != nil {
		return err
	}

	_, _ = s.users.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{
			"$set": bson.M{
				"isDeleted": false,
				"updatedAt": now,
			},
			"$unset": bson.M{
				"deletedAt": "",
			},
		},
	)

	s.sendUserEmail(
		ctx,
		userID,
		"Your Blomiva account was restored",
		"Your account deletion request has been cancelled and your account has been restored.",
	)

	return nil
}

func (s *Service) sendUserEmail(ctx context.Context, userID string, subject string, message string) {
	email := s.findUserEmail(ctx, userID)
	if email == "" || s.mail == nil {
		return
	}

	_ = s.mail.Send(ctx, email, subject, message)
}

func (s *Service) findUserEmail(ctx context.Context, userID string) string {
	var user struct {
		Email string `bson:"email"`
	}

	err := s.users.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err == nil && strings.TrimSpace(user.Email) != "" {
		return strings.TrimSpace(user.Email)
	}

	err = s.users.FindOne(ctx, bson.M{"id": userID}).Decode(&user)
	if err == nil && strings.TrimSpace(user.Email) != "" {
		return strings.TrimSpace(user.Email)
	}

	err = s.users.FindOne(ctx, bson.M{"userId": userID}).Decode(&user)
	if err == nil && strings.TrimSpace(user.Email) != "" {
		return strings.TrimSpace(user.Email)
	}

	return ""
}

func cleanLower(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
