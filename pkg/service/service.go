package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"log/slog"
	"test-rest-api/pkg/model"
	"test-rest-api/pkg/repository"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	salt            = "dfkdsjfklsdfj_sldkoi3242343242"
	signingKey      = "sdaskjdhkjahriw3or3asdsadad"
	tokenTTL        = 6 * time.Minute
	refreshTokenTTL = 24 * time.Hour * 30
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int64 `json:"userId"`
}

type Service struct {
	logger *slog.Logger
	repo   *repository.Repository
}

func (s *Service) AddRefferer(id int64) {
	s.repo.AddRefferer(id)
}

func NewService(logger *slog.Logger, repo *repository.Repository) *Service {
	return &Service{
		logger: logger,
		repo:   repo,
	}
}
func (s *Service) AddTask(task *model.TaskInput) error {
	return s.repo.AddTask(task)
}

func (s *Service) CompleteTask(id int64, task *model.CompleteTask) error {
	return s.repo.ComplitTask(id, task)
}

func (s *Service) RemoveTask(id int64) error {
	return s.repo.RemoveTask(id)
}

func (s *Service) CreateUser(user model.User) (int64, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}
func (s *Service) GetUserStatus(id int64) (*model.UserStatus, error) {

	return s.repo.GetUserStatus(id)
}
func (s *Service) GetLeaders() ([]*model.UsersWithPoints, error) {

	return s.repo.GetLeaders()
}

func (s *Service) UpdateUser(user model.User) (int64, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.UpdateUser(user)
}

func (s *Service) GenerateToken(email, password string) (string, string, *model.User, error) {
	user, err := s.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		s.logger.Error("Get User error:%v", err)
		return "", "", user, errors.New("User with the credentials not found!")
	}
	t, e := s.generateToken(user.Id, tokenTTL)
	at, e := s.generateToken(user.Id, refreshTokenTTL)
	return t, at, user, e
}
func (s *Service) GetUser(email string) (*model.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return user, err
	}
	return user, err
}

func (s *Service) ParseToken(accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	s.logger.Info("claims", claims)
	return claims.UserId, nil
}

func (auth *Service) generateToken(id int64, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix()},
		id,
	})
	return token.SignedString([]byte(signingKey))
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
