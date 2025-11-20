package repositories

import (
	"context"
	"fmt"
	"strings"

	pgx "github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type CarRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewCarRepositoryImpl(db *pgxpool.Pool) ports.CarRepository {
	return &CarRepositoryImpl{db: db}
}

func (r *CarRepositoryImpl) CreateCar(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	const insertCarQuery = `
		INSERT INTO cars (
			name,
			image_url,
			only_with_driver,
			car_mark_id,
			car_category_id,
			price_per_day
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var markID any
	if car.Mark != nil {
		markID = car.Mark.ID
	}

	var categoryID any
	if car.Category != nil {
		categoryID = car.Category.ID
	}

	err = tx.QueryRow(ctx, insertCarQuery,
		car.Name,
		car.ImageUrl,
		car.OnlyWithDriver,
		markID,
		categoryID,
		car.PricePerDay,
	).Scan(&car.ID)
	if err != nil {
		return nil, err
	}

	if len(car.Tags) > 0 {
		const insertCarTagQuery = `
			INSERT INTO car_car_tags (car_id, car_tag_id)
			VALUES ($1, $2)
		`

		for _, tag := range car.Tags {
			if tag == nil {
				continue
			}
			if _, err = tx.Exec(ctx, insertCarTagQuery, car.ID, tag.ID); err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetCarByID(ctx, car.ID)
}

func (r *CarRepositoryImpl) GetCarByID(ctx context.Context, id int64) (*entities.Car, error) {
	const querySelect = `
		SELECT 
			c.id,
			c.name,
			COALESCE(c.image_url, ''),
			c.only_with_driver,
			c.price_per_day,
			c.created_at,
			c.updated_at,
			cm.id,
			cm.name,
			cc.id,
			cc.name
		FROM cars c
		LEFT JOIN car_marks cm ON c.car_mark_id = cm.id
		LEFT JOIN car_categories cc ON c.car_category_id = cc.id
		WHERE c.id = $1
	`

	var car entities.Car
	var pricePerDay int64
	var markID *int64
	var markName *string
	var categoryID *int64
	var categoryName *string

	err := r.db.QueryRow(ctx, querySelect, id).Scan(
		&car.ID,
		&car.Name,
		&car.ImageUrl,
		&car.OnlyWithDriver,
		&pricePerDay,
		&car.CreatedAt,
		&car.UpdatedAt,
		&markID,
		&markName,
		&categoryID,
		&categoryName,
	)
	if err != nil {
		return nil, err
	}

	if markID != nil {
		car.Mark = &entities.CarMark{
			ID:   *markID,
			Name: derefString(markName),
		}
	}
	if categoryID != nil {
		car.Category = &entities.CarCategory{
			ID:   *categoryID,
			Name: derefString(categoryName),
		}
	}

	const queryFetchTags = `
		SELECT
			ct.id,
			ct.name,
			ct.created_at,
			ct.updated_at
		FROM car_tags ct
		JOIN car_car_tags cct ON ct.id = cct.car_tag_id
		WHERE cct.car_id = $1
	`

	rows, err := r.db.Query(ctx, queryFetchTags, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	car.Tags = make([]*entities.CarTag, 0)

	for rows.Next() {
		tag := &entities.CarTag{}
		if err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		); err != nil {
			return nil, err
		}
		car.Tags = append(car.Tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &car, nil
}

func (r CarRepositoryImpl) UpdateCar(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const updateCarQuery = `
		UPDATE cars
		SET
			name = $1,
			image_url = $2,
			only_with_driver = $3,
			car_mark_id = $4,
			car_category_id = $5,
			price_per_day = $6,
			updated_at = NOW()
		WHERE id = $7
	`

	var markID any
	if car.Mark != nil {
		markID = car.Mark.ID
	}

	var categoryID any
	if car.Category != nil {
		categoryID = car.Category.ID
	}

	price := int64(car.PricePerDay)

	cmdTag, err := tx.Exec(ctx, updateCarQuery,
		car.Name,
		car.ImageUrl,
		car.OnlyWithDriver,
		markID,
		categoryID,
		price,
		car.ID,
	)
	if err != nil {
		return nil, err
	}
	if cmdTag.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	const deleteTagsQuery = `DELETE FROM car_car_tags WHERE car_id = $1`
	if _, err := tx.Exec(ctx, deleteTagsQuery, car.ID); err != nil {
		return nil, err
	}

	if len(car.Tags) > 0 {
		const insertTagQuery = `
			INSERT INTO car_car_tags (car_id, car_tag_id)
			VALUES ($1, $2)
		`
		for _, tag := range car.Tags {
			if tag == nil {
				continue
			}
			if _, err := tx.Exec(ctx, insertTagQuery, car.ID, tag.ID); err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	updated, err := r.GetCarByID(ctx, car.ID)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (r CarRepositoryImpl) DeleteCar(ctx context.Context, id int64) error {
	const deleteCarQuery = `DELETE FROM cars WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, deleteCarQuery, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r CarRepositoryImpl) ListCars(ctx context.Context, offset, limit int64, name string, markID int64, categoryID int64) (int64, []*entities.Car, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	conditions := []string{"1=1"}
	args := []any{}
	argPos := 1

	if name != "" {
		conditions = append(conditions, fmt.Sprintf("c.name ILIKE $%d", argPos))
		args = append(args, "%"+name+"%")
		argPos++
	}

	if markID != 0 {
		conditions = append(conditions, fmt.Sprintf("cm.id = $%d", argPos))
		args = append(args, markID)
		argPos++
	}

	if categoryID != 0 {
		conditions = append(conditions, fmt.Sprintf("cc.id = $%d", argPos))
		args = append(args, categoryID)
		argPos++
	}

	whereSQL := strings.Join(conditions, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM cars c
		LEFT JOIN car_marks cm ON c.car_mark_id = cm.id
		WHERE %s
	`, whereSQL)

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return 0, nil, err
	}
	if total == 0 {
		return 0, []*entities.Car{}, nil
	}

	selectQuery := fmt.Sprintf(`
		SELECT 
			c.id,
			c.name,
			COALESCE(c.image_url, ''),
			c.only_with_driver,
			c.price_per_day,
			c.created_at,
			c.updated_at,
			cm.id,
			cm.name,
			cc.id,
			cc.name
		FROM cars c
		LEFT JOIN car_marks cm ON c.car_mark_id = cm.id
		LEFT JOIN car_categories cc ON c.car_category_id = cc.id
		WHERE %s
		ORDER BY c.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereSQL, argPos, argPos+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, selectQuery, args...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	cars := make([]*entities.Car, 0)

	for rows.Next() {
		var car entities.Car
		var pricePerDay int64

		var markIDPtr *int64
		var markNamePtr *string
		var categoryIDPtr *int64
		var categoryNamePtr *string

		if err := rows.Scan(
			&car.ID,
			&car.Name,
			&car.ImageUrl,
			&car.OnlyWithDriver,
			&pricePerDay,
			&car.CreatedAt,
			&car.UpdatedAt,
			&markIDPtr,
			&markNamePtr,
			&categoryIDPtr,
			&categoryNamePtr,
		); err != nil {
			return 0, nil, err
		}
		if markIDPtr != nil {
			car.Mark = &entities.CarMark{
				ID:   *markIDPtr,
				Name: derefString(markNamePtr),
			}
		}

		if categoryIDPtr != nil {
			car.Category = &entities.CarCategory{
				ID:   *categoryIDPtr,
				Name: derefString(categoryNamePtr),
			}
		}

		tags, err := r.getTagsByCarID(ctx, car.ID)
		if err != nil {
			return 0, nil, err
		}
		car.Tags = tags

		c := car
		cars = append(cars, &c)
	}

	if err := rows.Err(); err != nil {
		return 0, nil, err
	}

	return total, cars, nil
}

func (r CarRepositoryImpl) getTagsByCarID(ctx context.Context, carID int64) ([]*entities.CarTag, error) {
	const query = `
		SELECT
			ct.id,
			ct.name,
			ct.created_at,
			ct.updated_at
		FROM car_tags ct
		JOIN car_car_tags cct ON ct.id = cct.car_tag_id
		WHERE cct.car_id = $1
	`

	rows, err := r.db.Query(ctx, query, carID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]*entities.CarTag, 0)

	for rows.Next() {
		tag := &entities.CarTag{}
		if err := rows.Scan(
			&tag.ID,
			&tag.Name,
			&tag.CreatedAt,
			&tag.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
