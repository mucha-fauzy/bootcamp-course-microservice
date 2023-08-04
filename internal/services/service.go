package services

import (
	"bootcamp_course_microservice/internal/models"
	"bootcamp_course_microservice/internal/repository"
)

type Service interface {
	ReadCourseByUserID(userID string) ([]models.Course, error)
	CreateCourse(course *models.Course) error
}

type ServiceImpl struct {
	Repo repository.Repository
}

func ProvideService(r repository.Repository) *ServiceImpl {
	return &ServiceImpl{
		Repo: r,
	}
}

func (s *ServiceImpl) CreateCourse(course *models.Course) error {
	return s.Repo.CreateCourse(course)
}

func (s *ServiceImpl) ReadCourseByUserID(userID string) ([]models.Course, error) {
	return s.Repo.ReadCourseByUserID(userID)
}
