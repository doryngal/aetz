package services

import (
	"binai.net/v2/config"
	"binai.net/v2/internal/models"
	"binai.net/v2/internal/repository"
	"binai.net/v2/internal/shared/utils"
	"context"
	"errors"
	"fmt"

	"math/rand"
	"time"
)

type AuthService struct {
	Repo repository.AuthRepository
	SMTP struct {
		Host     string
		Port     int
		User     string
		Password string
		From     string
	}
}

func NewAuthService(repo repository.AuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		Repo: repo,
		SMTP: cfg.SMTP,
	}
}

func (a *AuthService) Register(user *models.User) error {
	isRegistered, err := a.Repo.IsEmailRegistered(user.Email)
	if err != nil {
		return err
	}
	if isRegistered {
		return errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword // Сохраняем хеш
	return a.Repo.CreateUser(user)
}

func (a *AuthService) Login(email, password string) (string, string, error) {
	user, err := a.Repo.GetUserByEmail(email)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	err = utils.CheckPasswordHash(user.PasswordHash, password)
	if err != nil {
		return "", "", errors.New("password is not correct")
	}

	token, err := utils.GenerateNewAccessToken(user)
	if err != nil {
		return "", "", err
	}

	return token, user.Role, nil
}

func (s *AuthService) ConfirmationCode(ctx context.Context, email string) error {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	// Проверка времени последней отправки
	if user.LastResetSentAt != nil && time.Since(*user.LastResetSentAt) < time.Minute {
		return errors.New("вы можете запросить код восстановления только раз в минуту")
	}

	// Генерация 4-значного кода
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	expiresAt := time.Now().Add(10 * time.Minute)
	now := time.Now()

	// Обновление пользователя в БД
	err = s.Repo.UpdateResetCode(ctx, user.ID, code, expiresAt, now)
	if err != nil {
		return err
	}

	// Генерация HTML-письма
	//body := GenerateConfirmationEmail(code)
	//subject := "Код восстановления пароля / Құпия сөзді қалпына келтіру коды"

	// Отправка кода на почту
	//err = notification.SendEmail(
	//	s.SMTP.Host, s.SMTP.Port,
	//	s.SMTP.User, s.SMTP.Password,
	//	s.SMTP.From, user.Email,
	//	subject, body,
	//)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, email, code string) (string, string, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	// Проверка кода
	if user.ResetCode == nil || *user.ResetCode != code || user.ResetCodeExpiresAt == nil || time.Now().After(*user.ResetCodeExpiresAt) {
		return "", "", errors.New("invalid or expired reset code")
	}

	// Очистка кода восстановления
	err = s.Repo.ClearResetCode(ctx, user.ID)

	token, err := utils.GenerateNewAccessToken(user)
	if err != nil {
		return "", "", err
	}

	return token, user.Role, nil
}

func GenerateConfirmationEmail(code string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Binai.kz</title>
	</head>
	<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
		<h2 style="color: #0056b3;">Код для восстановления пароля</h2>
		<p>Здравствуйте,</p>
		<p>Ваш код для восстановления пароля:</p>
		<h3 style="color: #ff6600;">%s</h3>
		<p>Этот код действителен в течение <strong>10 минут</strong>.</p>
		<hr>
		<h2 style="color: #0056b3;">Құпия сөзді қалпына келтіру коды</h2>
		<p>Сәлеметсіз бе,</p>
		<p>Құпия сөзді қалпына келтіру кодыңыз:</p>
		<h3 style="color: #ff6600;">%s</h3>
		<p>Бұл код <strong>10 минут</strong> ішінде жарамды.</p>
		<p>SkillServe платформасын таңдағаныңызға рахмет!</p>
		<p>С уважением,<br>Команда Binai.kz</p>
		<p>Құрметпен,<br>Binai.kz командасы</p>
	</body>
	</html>
	`, code, code)
}
