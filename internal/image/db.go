package image

import (
	"database/sql"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"

	"github.com/macarrie/relique/internal/api_helpers"
	"github.com/macarrie/relique/internal/client"
	"github.com/macarrie/relique/internal/db"
	"github.com/macarrie/relique/internal/module"
	"github.com/macarrie/relique/internal/repo"
)

func (img *Image) Save() (int64, error) {
	tx, err := db.Handler().Begin()
	// Defers are stacked, defer are executed in reverse order of stacking
	defer func() {
		if err != nil {
			img.GetLog().With(
				slog.Any("error", err),
			).Debug("Rollback image save")
			tx.Rollback()
		}
	}()

	if err != nil {
		return 0, fmt.Errorf("cannot start transaction to save image: %w", err)
	}

	if img.ID != 0 {
		id, err := img.Update(tx)
		if err != nil || id == 0 {
			return 0, fmt.Errorf("cannot update image: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return 0, fmt.Errorf("cannot commit image save transaction: %w", err)
		}

		return id, err
	}

	img.GetLog().Debug("Saving image into database")

	request := sq.Insert("images").SetMap(sq.Eq{
		"uuid":               img.Uuid,
		"created_at":         img.CreatedAt,
		"module_type":        img.Module.ModuleType,
		"client_name":        img.Client.Name,
		"repo_name":          img.Repository.GetName(),
		"number_of_elements": img.NumberOfElements,
		"number_of_files":    img.NumberOfFiles,
		"number_of_folders":  img.NumberOfFolders,
		"size_on_disk":       img.SizeOnDisk,
	})
	query, args, err := request.ToSql()
	if err != nil {
		return 0, fmt.Errorf("cannot build sql query: %w", err)
	}

	result, err := tx.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("cannot save image into db: %w", err)
	}

	img.ID, err = result.LastInsertId()
	if img.ID == 0 || err != nil {
		return 0, fmt.Errorf("cannot get last insert ID: %w", err)
	}

	img.GetLog().Debug("Commit image save transaction")
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("cannot commit image save transaction: %w", err)
	}

	return img.ID, nil
}

func (img *Image) Update(tx *sql.Tx) (int64, error) {
	img.GetLog().Debug("Updating image details into database")

	request := sq.Update("images").SetMap(sq.Eq{
		"module_type":        img.Module.ModuleType,
		"client_name":        img.Client.Name,
		"repo_name":          img.Repository.GetName(),
		"number_of_elements": img.NumberOfElements,
		"number_of_files":    img.NumberOfFiles,
		"number_of_folders":  img.NumberOfFolders,
		"size_on_disk":       img.SizeOnDisk,
	}).Where(
		"uuid = ?",
		img.Uuid,
	)
	query, args, err := request.ToSql()
	if err != nil {
		return 0, fmt.Errorf("cannot build sql query: %w", err)
	}

	var result sql.Result
	if tx == nil {
		result, err = db.Handler().Exec(query, args...)
	} else {
		result, err = tx.Exec(query, args...)
	}
	if err != nil {
		return 0, fmt.Errorf("cannot update image into db: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected != 1 || err != nil {
		return 0, fmt.Errorf("no rows affected: %w", err)
	}

	return img.ID, nil
}

func GetByUuid(uuid string) (Image, error) {
	slog.With(
		slog.String("uuid", uuid),
	).Debug("Looking for image in database")

	request := sq.Select(
		"id",
		"uuid",
		"created_at",
		"client_name",
		"module_type",
		"repo_name",
		"number_of_elements",
		"number_of_files",
		"number_of_folders",
		"size_on_disk",
	).From("images").Where("uuid = ?", uuid)
	query, args, err := request.ToSql()
	if err != nil {
		return Image{}, fmt.Errorf("cannot build sql query: %w", err)
	}

	row := db.Handler().QueryRow(query, args...)

	var img Image
	if err := row.Scan(&img.ID,
		&img.Uuid,
		&img.CreatedAt,
		&img.ClientName,
		&img.ModuleType,
		&img.RepoName,
		&img.NumberOfElements,
		&img.NumberOfFiles,
		&img.NumberOfFolders,
		&img.SizeOnDisk,
	); err == sql.ErrNoRows {
		return Image{}, fmt.Errorf("no image with UUID '%s' found in db", uuid)
	} else if err != nil {
		return Image{}, fmt.Errorf("cannot retrieve image from db: %w", err)
	}

	imgStorageFolderPath, err := img.GetStorageFolderPath()
	if err != nil {
		return Image{}, fmt.Errorf("cannot get image storage folder path: %w", err)
	}

	modFilePath := fmt.Sprintf("%s/module.toml", imgStorageFolderPath)
	mod, err := module.LoadFromFile(modFilePath)
	if err != nil {
		return Image{}, fmt.Errorf("linked module cannot be loaded from file: %w", err)
	}
	img.Module = mod

	clFilePath := fmt.Sprintf("%s/client.toml", imgStorageFolderPath)
	cl, err := client.LoadFromFile(clFilePath)
	if err != nil {
		return Image{}, fmt.Errorf("linked client cannot be loaded from file: %w", err)
	}
	img.Client = cl

	repoFilePath := fmt.Sprintf("%s/repo.toml", imgStorageFolderPath)
	r, err := repo.LoadFromFile(repoFilePath)
	if err != nil {
		return Image{}, fmt.Errorf("linked repo cannot be loaded from file: %w", err)
	}
	img.Repository = r

	return img, nil
}

func Search(p api_helpers.PaginationParams, modulesInstallPath string) ([]Image, error) {
	slog.Debug("Searching for jobs in db")
	var imgs []Image

	// TODO: Prepare request and clean data to avoid SQL injections
	// TODO: handle status and backup type
	request := sq.Select(
		"uuid",
	).From(
		"images",
	)
	if p.Limit > 0 {
		request = request.Limit(p.Limit)
	}
	if p.Offset > 0 {
		request = request.Offset(p.Offset)
	}

	// TODO/ Handle pagination
	// TODO/ Handle search parameters
	request = request.OrderBy("images.id DESC")

	query, args, err := request.ToSql()
	if err != nil {
		return []Image{}, fmt.Errorf("cannot build sql query: %w", err)
	}

	rows, err := db.Handler().Query(query, args...)
	if err == sql.ErrNoRows {
		return imgs, nil
	} else if err != nil {
		return imgs, fmt.Errorf("cannot search jobs IDs from db: %w", err)
	}

	uuids := make([]string, 0)
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			slog.With(
				slog.Any("error", err),
			).Error("Cannot parse image uuid from db")
		}
		if u != "" {
			uuids = append(uuids, u)
		}
	}

	if len(uuids) == 0 {
		// No previous job found
		return imgs, nil
	}

	for _, imgUuid := range uuids {
		imgFromDB, err := GetByUuid(imgUuid)
		if err != nil {
			slog.With(
				slog.Any("error", err),
				slog.String("uuid", imgUuid),
			).Error("Cannot get image with uuid from db")
			continue
		}
		if imgFromDB.ID == 0 {
			slog.With(
				slog.Any("error", err),
				slog.String("uuid", imgUuid),
			).Error("No image with this uuid found in db")
			continue
		}

		imgs = append(imgs, imgFromDB)
	}

	return imgs, nil
}

// TODO: Handle search parameters to have a selective count
func Count() (uint64, error) {
	var count uint64

	request := sq.Select(
		"COUNT(*)",
	).From(
		"images",
	)

	request = request.OrderBy("images.id DESC")

	query, args, err := request.ToSql()
	if err != nil {
		return 0, fmt.Errorf("cannot build sql query: %w", err)
	}

	queryErr := db.Handler().QueryRow(query, args...).Scan(&count)
	if queryErr == sql.ErrNoRows {
		return 0, nil
	} else if queryErr != nil {
		return 0, fmt.Errorf("cannot count images from db: %w", err)
	}

	return count, nil
}
