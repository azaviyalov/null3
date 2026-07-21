package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) WithTx(ctx context.Context, fn func(repo *Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&Repository{db: tx})
	})
}

func (r *Repository) SessionRepository() *session.Repository {
	return session.NewRepository(r.db)
}

func (r *Repository) GetUserByID(ctx context.Context, id uint) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("get user %d: %w", id, err)
	}
	return &user, nil
}

func (r *Repository) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("login = ?", login).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("get user by login: %w", err)
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (r *Repository) UpdateUserPassword(ctx context.Context, userID uint, passwordHash string) error {
	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("password_hash", passwordHash).Error; err != nil {
		return fmt.Errorf("update password for user %d: %w", userID, err)
	}
	return nil
}

func (r *Repository) CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) (*PasswordResetToken, error) {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return nil, fmt.Errorf("create password reset token: %w", err)
	}
	return token, nil
}

func (r *Repository) GetPasswordResetTokenByHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error) {
	var token PasswordResetToken
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("get password reset token: %w", err)
	}
	return &token, nil
}

func (r *Repository) DeletePasswordResetTokensByUser(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&PasswordResetToken{}).Error; err != nil {
		return fmt.Errorf("delete password reset tokens for user %d: %w", userID, err)
	}
	return nil
}

func (r *Repository) CreateInvite(ctx context.Context, invite *Invite) (*Invite, error) {
	if err := r.db.WithContext(ctx).Create(invite).Error; err != nil {
		return nil, fmt.Errorf("create invite: %w", err)
	}
	return invite, nil
}

func (r *Repository) GetInviteByHash(ctx context.Context, tokenHash string) (*Invite, error) {
	var invite Invite
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&invite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("get invite: %w", err)
	}
	return &invite, nil
}

func (r *Repository) SaveInvite(ctx context.Context, invite *Invite) error {
	if err := r.db.WithContext(ctx).Save(invite).Error; err != nil {
		return fmt.Errorf("save invite %d: %w", invite.ID, err)
	}
	return nil
}
