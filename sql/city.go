package sql

import (
	"context"
	"database/sql"
	"mysql/app/apperr"
	"mysql/app/entity"
	"mysql/app/service"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var _ service.CityService = (*CityService)(nil)

type CityService struct {
	db *sql.DB
}

func NewCityService(db *sql.DB) *CityService {
	return &CityService{db}
}

func (s *CityService) CreateCity(ctx context.Context, city *entity.City) error {

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createCity(ctx, tx, city); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *CityService) DeleteCity(ctx context.Context, id int64) error {

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteCity(ctx, tx, id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *CityService) UpdateCity(ctx context.Context, id int64, cup service.CityUpdate) error {

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := updateCity(ctx, tx, id, cup); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *CityService) FindCities(ctx context.Context, filter service.CityFilter) (cities entity.Cities, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	cities, err = findCities(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, apperr.Errorf(apperr.EINTERNAL, "errore: %v", err)
	}

	return cities, nil
}

func (s *CityService) FindIdByName(ctx context.Context, name string) (id *int64, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	id, err = findIdByName(ctx, tx, name)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, apperr.Errorf(apperr.EINTERNAL, "errore: %v", err)
	}

	return id, nil
}

func createCity(ctx context.Context, tx *sql.Tx, city *entity.City) error {

	if err := city.Validate(); err != nil {
		return err
	}

	if res, err := tx.ExecContext(ctx, "INSERT INTO cities(name, population) VALUES (?,?)", city.Name, city.Population); err != nil {
		return apperr.Errorf(apperr.EINTERNAL, "errore nell'inserimento: %v", err)
	} else if city.Id, err = res.LastInsertId(); err != nil {
		return apperr.Errorf(apperr.EINTERNAL, "errore recupero id: %v", err)
	}

	return nil
}

func deleteCity(ctx context.Context, tx *sql.Tx, id int64) error {

	if city, err := findCityById(ctx, tx, id); err != nil {
		return err
	} else if err := city.Validate(); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM cities WHERE id = ?", id); err != nil {
		return apperr.Errorf(apperr.EINTERNAL, "errore nella cancellazione della città: %v", err)
	}

	return nil
}

func updateCity(ctx context.Context, tx *sql.Tx, id int64, cup service.CityUpdate) error {

	if _, err := findCityById(ctx, tx, id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "UPDATE cities SET Population = ? WHERE id = ?;", cup.Population, id); err != nil {
		return apperr.Errorf(apperr.EINTERNAL, "errore nell'aggiornamento della città: %v", err)
	}

	return nil
}

func findCityById(ctx context.Context, tx *sql.Tx, id int64) (*entity.City, error) {

	c, err := findCities(ctx, tx, service.CityFilter{Id: &id})
	if err != nil {
		return nil, err
	} else if len(c) == 0 {
		return nil, apperr.Errorf(apperr.ENOTFOUND, "city not found")
	}

	return c[0], nil
}

func findIdByName(ctx context.Context, tx *sql.Tx, name string) (id *int64, err error) {

	rows, err := tx.QueryContext(ctx, "SELECT id FROM cities WHERE name = ?", name)
	if err != nil {
		return nil, apperr.Errorf(apperr.EINTERNAL, "errore nella ricerca della città: %v", err)
	}

	defer rows.Close()

	var Id int64

	for rows.Next() {

		var city entity.City

		if err := rows.Scan(
			&city.Id,
		); err != nil {
			return nil, apperr.Errorf(apperr.EINTERNAL, "failed to scan city: %v", err)
		}

		Id = city.Id
	}
	if err := rows.Err(); err != nil {
		return nil, apperr.Errorf(apperr.EINTERNAL, "failed to iterate over cities: %v", err)
	}

	return &Id, nil
}

func findCityByPopulation(ctx context.Context, tx *sql.Tx, population int) (entity.Cities, error) {

	u, err := findCities(ctx, tx, service.CityFilter{Population: &population})
	if err != nil {
		return nil, err
	} else if len(u) == 0 {
		return nil, apperr.Errorf(apperr.ENOTFOUND, "city not found")
	}

	return u, nil
}

func findCityByPopulationGte(ctx context.Context, tx *sql.Tx, population int) (entity.Cities, error) {

	u, err := findCities(ctx, tx, service.CityFilter{PopulationGte: &population})
	if err != nil {
		return nil, err
	} else if len(u) == 0 {
		return nil, apperr.Errorf(apperr.ENOTFOUND, "city not found")
	}

	return u, nil
}

func findCityByPopulationLte(ctx context.Context, tx *sql.Tx, population int) (entity.Cities, error) {

	u, err := findCities(ctx, tx, service.CityFilter{PopulationLte: &population})
	if err != nil {
		return nil, err
	} else if len(u) == 0 {
		return nil, apperr.Errorf(apperr.ENOTFOUND, "city not found")
	}

	return u, nil
}

func findCities(ctx context.Context, tx *sql.Tx, filter service.CityFilter) (_ entity.Cities, err error) {

	where, args := []string{"1 = 1"}, []interface{}{}

	if v := filter.Id; v != nil {
		where = append(where, "id = ?")
		args = append(args, *v)
	}
	if v := filter.Name; v != nil {
		where = append(where, "name = ?")
		args = append(args, *v)
	}
	if v := filter.Population; v != nil {
		where = append(where, "population = ?")
		args = append(args, *v)
	}
	if v := filter.PopulationGte; v != nil {
		where = append(where, "population >= ?")
		args = append(args, *v)
	}
	if v := filter.PopulationLte; v != nil {
		where = append(where, "population <= ?")
		args = append(args, *v)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT
		    id,
		    name,
		    population
		FROM cities
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`, args...,
	)
	if err != nil {
		return nil, apperr.Errorf(apperr.EINTERNAL, "failed to query cities: %v", err)
	}
	defer rows.Close()

	cities := make(entity.Cities, 0)

	for rows.Next() {

		var city entity.City

		if err := rows.Scan(
			&city.Id,
			&city.Name,
			&city.Population,
		); err != nil {
			return nil, apperr.Errorf(apperr.EINTERNAL, "failed to scan city: %v", err)
		}

		cities = append(cities, &city)
	}
	if err := rows.Err(); err != nil {
		return nil, apperr.Errorf(apperr.EINTERNAL, "failed to iterate over cities: %v", err)
	}

	return cities, nil
}
