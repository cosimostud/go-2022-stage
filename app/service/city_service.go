package service

import (
	"context"
	"mysql/app/entity"
)

type CityService interface {
	CreateCity(ctx context.Context, city *entity.City) error
	DeleteCity(ctx context.Context, id int64) error
	UpdateCity(ctx context.Context, id int64, cup CityUpdate) error
	//Finds cities by id, population, populationGte, population
	FindCities(ctx context.Context, filter CityFilter) (cities entity.Cities, err error)
	FindIdByName(ctx context.Context, name string) (id *int64, err error)
}

type CityUpdate struct {
	Population *int
}

type CityFilter struct {
	Id            *int64
	Name          *string
	Names         []string
	Population    *int
	PopulationGte *int
	PopulationLte *int

	Offset int
	Limit  int
}