package repository_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"cloudtrail-enrichment-api-golang/database/postgresql" // CAMBIO: Importación corregida a ruta absoluta
	"cloudtrail-enrichment-api-golang/models"              // CAMBIO: Mover esta línea a la agrupación de paquetes locales

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

// setUpDBMock establece un mock de la base de datos para pruebas.
func setUpDBMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("un error ocurrió al crear el mock de la base de datos: %v", err)
	}
	return db, mock
}

// CustomArgsMatcher es un implementador de driver.Matcher que compara argumentos de forma flexible.
type CustomArgsMatcher struct {
	ExpectedArgs []driver.Value
}

func (m CustomArgsMatcher) Match(args []driver.Value) error {
	if len(m.ExpectedArgs) != len(args) {
		return fmt.Errorf("se esperaban %d argumentos, se recibieron %d", len(m.ExpectedArgs), len(args))
	}
	for i, expected := range m.ExpectedArgs {
		actual := args[i]
		// Se ignoran las comparaciones de tiempo y UUID debido a la naturaleza dinámica.
		if _, ok := expected.(time.Time); ok {
			continue
		}
		if _, ok := expected.(uuid.UUID); ok {
			continue
		}
		if expected != actual {
			return fmt.Errorf("el argumento %d no coincide. Se esperaba %v, se obtuvo %v", i, expected, actual)
		}
	}
	return nil
}

// TestAuthRepository_InsertUser prueba el método InsertUser.
func TestAuthRepository_InsertUser(t *testing.T) {
	db, mock := setUpDBMock(t)
	defer db.Close()

	// Llama a la función del paquete 'postgresql'
	repo := postgresql.NewAuthPostgresRepository(db)
	ctx := context.Background()
	now := time.Now()
	userID := 1
	userUUID := uuid.New().String()

	userToInsert := &models.User{
		UUID:         userUUID,
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "user",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Caso de éxito
	t.Run("success", func(t *testing.T) {
		// Se espera una consulta INSERT y se simula el retorno del ID del usuario.
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (uuid, email, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`)).
			WithArgs(CustomArgsMatcher{ExpectedArgs: []driver.Value{userToInsert.UUID, userToInsert.Email, userToInsert.PasswordHash, userToInsert.Role, userToInsert.CreatedAt, userToInsert.UpdatedAt}}).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

		err := repo.InsertUser(ctx, userToInsert)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo %v", err)
		}
		if userToInsert.ID != userID {
			t.Errorf("se esperaba que el ID del usuario fuera %d, se obtuvo %d", userID, userToInsert.ID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de error
	t.Run("failure", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (uuid, email, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`)).
			WillReturnError(errors.New("error de base de datos simulado"))

		err := repo.InsertUser(ctx, userToInsert)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})
}

// TestAuthRepository_GetUserByEmail prueba el método GetUserByEmail.
func TestAuthRepository_GetUserByEmail(t *testing.T) {
	db, mock := setUpDBMock(t)
	defer db.Close()
	// Llama a la función del paquete 'postgresql'
	repo := postgresql.NewAuthPostgresRepository(db)
	ctx := context.Background()

	// Datos de usuario de prueba
	testUser := &models.User{
		ID:           1,
		UUID:         uuid.New().String(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Caso de éxito
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "uuid", "email", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(testUser.ID, testUser.UUID, testUser.Email, testUser.PasswordHash, testUser.Role, testUser.CreatedAt, testUser.UpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`)).
			WithArgs(testUser.Email).
			WillReturnRows(rows)

		user, err := repo.GetUserByEmail(ctx, testUser.Email)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo %v", err)
		}
		if user == nil || user.Email != testUser.Email {
			t.Errorf("el usuario retornado no coincide con el esperado")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de usuario no encontrado (sql.ErrNoRows)
	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`)).
			WithArgs("nonexistent@example.com").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByEmail(ctx, "nonexistent@example.com")
		if err == nil || !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("se esperaba sql.ErrNoRows, se obtuvo %v", err)
		}
		if user != nil {
			t.Errorf("se esperaba un usuario nulo, se obtuvo %v", user)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de error en la consulta
	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`)).
			WithArgs(testUser.Email).
			WillReturnError(errors.New("error de base de datos simulado"))

		_, err := repo.GetUserByEmail(ctx, testUser.Email)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})
}

// TestAuthRepository_GetUserByUUID prueba el método GetUserByUUID.
func TestAuthRepository_GetUserByUUID(t *testing.T) {
	db, mock := setUpDBMock(t)
	defer db.Close()
	// Llama a la función del paquete 'postgresql'
	repo := postgresql.NewAuthPostgresRepository(db)
	ctx := context.Background()

	// Datos de usuario de prueba
	testUser := &models.User{
		ID:           1,
		UUID:         uuid.New().String(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Caso de éxito
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "uuid", "email", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(testUser.ID, testUser.UUID, testUser.Email, testUser.PasswordHash, testUser.Role, testUser.CreatedAt, testUser.UpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE uuid = $1`)).
			WithArgs(testUser.UUID).
			WillReturnRows(rows)

		user, err := repo.GetUserByUUID(ctx, testUser.UUID)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo %v", err)
		}
		if user == nil || user.UUID != testUser.UUID {
			t.Errorf("el usuario retornado no coincide con el esperado")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de usuario no encontrado (sql.ErrNoRows)
	t.Run("user not found", func(t *testing.T) {
		nonExistentUUID := uuid.New().String()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE uuid = $1`)).
			WithArgs(nonExistentUUID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByUUID(ctx, nonExistentUUID)
		if err == nil || !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("se esperaba sql.ErrNoRows, se obtuvo %v", err)
		}
		if user != nil {
			t.Errorf("se esperaba un usuario nulo, se obtuvo %v", user)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de error en la consulta
	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`)).
			WithArgs(testUser.UUID).
			WillReturnError(errors.New("error de base de datos simulado"))

		_, err := repo.GetUserByUUID(ctx, testUser.UUID)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})
}

// TestAuthRepository_InsertToken prueba el método InsertToken.
func TestAuthRepository_InsertToken(t *testing.T) {
	db, mock := setUpDBMock(t)
	defer db.Close()
	// Llama a la función del paquete 'postgresql'
	repo := postgresql.NewAuthPostgresRepository(db)
	ctx := context.Background()

	testToken := &models.Token{
		UserID:    1,
		Email:     "test@example.com",
		Token:     "testtoken",
		TokenHash: "hashed_testtoken",
		Expiry:    time.Now().Add(time.Hour),
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Caso de éxito
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO tokens (user_id, email, token, token_hash, expiry, created_at, updated_at, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)).
			WithArgs(testToken.UserID, testToken.Email, testToken.Token, testToken.TokenHash, CustomArgsMatcher{ExpectedArgs: []driver.Value{testToken.Expiry}}, CustomArgsMatcher{ExpectedArgs: []driver.Value{testToken.CreatedAt}}, CustomArgsMatcher{ExpectedArgs: []driver.Value{testToken.UpdatedAt}}, testToken.Role).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.InsertToken(ctx, testToken)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de error
	t.Run("failure", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO tokens (user_id, email, token, token_hash, expiry, created_at, updated_at, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)).
			WillReturnError(errors.New("error de base de datos simulado"))

		err := repo.InsertToken(ctx, testToken)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})
}

// TestAuthRepository_GetTokenByTokenHash prueba el método GetTokenByTokenHash.
func TestAuthRepository_GetTokenByTokenHash(t *testing.T) {
	db, mock := setUpDBMock(t)
	defer db.Close()
	// Llama a la función del paquete 'postgresql'
	repo := postgresql.NewAuthPostgresRepository(db)
	ctx := context.Background()

	// Datos de token de prueba
	testToken := &models.Token{
		ID:        1,
		UserID:    1,
		Email:     "test@example.com",
		Token:     "testtoken",
		TokenHash: "hashed_testtoken",
		Expiry:    time.Now().Add(time.Hour),
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Caso de éxito
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "email", "token", "token_hash", "expiry", "created_at", "updated_at", "role"}).
			AddRow(testToken.ID, testToken.UserID, testToken.Email, testToken.Token, testToken.TokenHash, testToken.Expiry, testToken.CreatedAt, testToken.UpdatedAt, testToken.Role)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens WHERE token_hash = $1`)).
			WithArgs(testToken.TokenHash).
			WillReturnRows(rows)

		token, err := repo.GetTokenByTokenHash(ctx, testToken.TokenHash)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo %v", err)
		}
		if token == nil || token.TokenHash != testToken.TokenHash {
			t.Errorf("el token retornado no coincide con el esperado")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de token no encontrado (sql.ErrNoRows)
	t.Run("token not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens WHERE token_hash = $1`)).
			WithArgs("nonexistent_hash").
			WillReturnError(sql.ErrNoRows)

		token, err := repo.GetTokenByTokenHash(ctx, "nonexistent_hash")
		if err == nil || !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("se esperaba sql.ErrNoRows, se obtuvo %v", err)
		}
		if token != nil {
			t.Errorf("se esperaba un token nulo, se obtuvo %v", token)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de error en la consulta
	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens WHERE token_hash = $1`)).
			WithArgs(testToken.TokenHash).
			WillReturnError(errors.New("error de base de datos simulado"))

		_, err := repo.GetTokenByTokenHash(ctx, testToken.TokenHash)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})
}

// TestAuthRepository_DeleteTokensByUserID prueba el método DeleteTokensByUserID.
func TestAuthRepository_DeleteTokensByUserID(t *testing.T) {
	db, mock := setUpDBMock(t)
	defer db.Close()
	// Llama a la función del paquete 'postgresql'
	repo := postgresql.NewAuthPostgresRepository(db)
	ctx := context.Background()
	userID := 1

	// Caso de éxito
	t.Run("success", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM tokens WHERE user_id = $1`)).
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteTokensByUserID(ctx, userID)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de error en la ejecución
	t.Run("failure", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM tokens WHERE user_id = $1`)).
			WithArgs(userID).
			WillReturnError(errors.New("error de base de datos simulado"))

		err := repo.DeleteTokensByUserID(ctx, userID)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})
}

// TestAuthRepository_GetTokenByToken prueba el método GetTokenByToken.
func TestAuthRepository_GetTokenByToken(t *testing.T) {
	db, mock := setUpDBMock(t)
	defer db.Close()
	// Llama a la función del paquete 'postgresql'
	repo := postgresql.NewAuthPostgresRepository(db)
	ctx := context.Background()

	// Datos de token de prueba
	testToken := &models.Token{
		ID:        1,
		UserID:    1,
		Email:     "test@example.com",
		Token:     "testtoken_plain",
		TokenHash: "hashed_testtoken_plain",
		Expiry:    time.Now().Add(time.Hour),
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Caso de éxito
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "email", "token", "token_hash", "expiry", "created_at", "updated_at", "role"}).
			AddRow(testToken.ID, testToken.UserID, testToken.Email, testToken.Token, testToken.TokenHash, testToken.Expiry, testToken.CreatedAt, testToken.UpdatedAt, testToken.Role)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens WHERE token = $1`)).
			WithArgs(testToken.Token).
			WillReturnRows(rows)

		token, err := repo.GetTokenByToken(ctx, testToken.Token)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo %v", err)
		}
		if token == nil || token.Token != testToken.Token {
			t.Errorf("el token retornado no coincide con el esperado")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de token no encontrado (sql.ErrNoRows)
	t.Run("token not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens WHERE token = $1`)).
			WithArgs("nonexistent_token").
			WillReturnError(sql.ErrNoRows)

		token, err := repo.GetTokenByToken(ctx, "nonexistent_token")
		if err == nil || !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("se esperaba sql.ErrNoRows, se obtuvo %v", err)
		}
		if token != nil {
			t.Errorf("se esperaba un token nulo, se obtuvo %v", token)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de error en la consulta
	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens WHERE token = $1`)).
			WithArgs(testToken.Token).
			WillReturnError(errors.New("error de base de datos simulado"))

		_, err := repo.GetTokenByToken(ctx, testToken.Token)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})
}

// TestAuthRepository_GetUserForToken prueba el método GetUserForToken.
func TestAuthRepository_GetUserForToken(t *testing.T) {
	db, mock := setUpDBMock(t)
	defer db.Close()
	// Llama a la función del paquete 'postgresql'
	repo := postgresql.NewAuthPostgresRepository(db)
	ctx := context.Background()
	userID := 1

	// Datos de usuario de prueba
	testUser := &models.User{
		ID:           userID,
		UUID:         uuid.New().String(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Caso de éxito
	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "uuid", "email", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(testUser.ID, testUser.UUID, testUser.Email, testUser.PasswordHash, testUser.Role, testUser.CreatedAt, testUser.UpdatedAt)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`)).
			WithArgs(userID).
			WillReturnRows(rows)

		user, err := repo.GetUserForToken(ctx, userID)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo %v", err)
		}
		if user == nil || user.ID != userID {
			t.Errorf("el usuario retornado no coincide con el esperado")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de usuario no encontrado (sql.ErrNoRows)
	t.Run("user not found", func(t *testing.T) {
		nonExistentUserID := 999
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`)).
			WithArgs(nonExistentUserID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserForToken(ctx, nonExistentUserID)
		if err == nil || !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("se esperaba sql.ErrNoRows, se obtuvo %v", err)
		}
		if user != nil {
			t.Errorf("se esperaba un usuario nulo, se obtuvo %v", user)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})

	// Caso de error en la consulta
	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`)).
			WithArgs(userID).
			WillReturnError(errors.New("error de base de datos simulado"))

		_, err := repo.GetUserForToken(ctx, userID)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron todas las expectativas: %v", err)
		}
	})
}
