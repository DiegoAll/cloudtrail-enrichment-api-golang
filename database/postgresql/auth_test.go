package postgresql

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

// TestMain se ejecuta una vez antes de todas las pruebas del paquete.
func TestMain(m *testing.M) {
	// Inicializar el logger para evitar panics.
	log.Println("Initializing logger for tests...")
	logger.Init()

	// Ejecutar todas las pruebas.
	os.Exit(m.Run())
}

// TestAuthPostgresRepository es la suite de pruebas para los métodos del repositorio.
func TestAuthPostgresRepository(t *testing.T) {
	// Configuración de la hora para mocks consistentes
	fixedTime := time.Date(2025, time.August, 6, 21, 0, 0, 0, time.UTC)

	// Crear una conexión de base de datos mockeada
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("se esperaba que la conexión al mock fuera exitosa, se obtuvo: %v", err)
	}
	defer db.Close()

	// Crear una instancia del repositorio con el mock de la base de datos
	repo := NewAuthPostgresRepository(db)
	ctx := context.Background()

	// --- Prueba de InsertUser ---
	t.Run("InsertUser - success", func(t *testing.T) {
		testUUID := uuid.New().String()
		user := &models.User{
			UUID:         testUUID,
			Email:        "test@example.com",
			PasswordHash: "hashedpassword",
			Role:         "user",
		}
		// Se espera una consulta que devuelva el ID
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(testUUID, user.Email, user.PasswordHash, user.Role, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.InsertUser(ctx, user)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo: %v", err)
		}
		if user.ID != 1 {
			t.Errorf("ID de usuario incorrecto. Se esperaba 1, se obtuvo: %d", user.ID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	t.Run("InsertUser - db error", func(t *testing.T) {
		testUUID := uuid.New().String()
		user := &models.User{
			UUID:         testUUID,
			Email:        "test@example.com",
			PasswordHash: "hashedpassword",
			Role:         "user",
		}
		expectedErr := errors.New("simulated database error")
		mock.ExpectQuery("INSERT INTO users").
			WithArgs(testUUID, user.Email, user.PasswordHash, user.Role, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(expectedErr)

		err := repo.InsertUser(ctx, user)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("error incorrecto. Se esperaba '%v', se obtuvo '%v'", expectedErr, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	// --- Prueba de GetUserByEmail ---
	t.Run("GetUserByEmail - success", func(t *testing.T) {
		email := "test@example.com"
		expectedUser := &models.User{
			ID:           1,
			UUID:         uuid.New().String(),
			Email:        email,
			PasswordHash: "hashedpassword",
			Role:         "user",
			CreatedAt:    fixedTime,
			UpdatedAt:    fixedTime,
		}
		rows := sqlmock.NewRows([]string{"id", "uuid", "email", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.UUID, expectedUser.Email, expectedUser.PasswordHash, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt)
		mock.ExpectQuery("SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users").
			WithArgs(email).
			WillReturnRows(rows)

		user, err := repo.GetUserByEmail(ctx, email)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo: %v", err)
		}
		if user.Email != expectedUser.Email {
			t.Errorf("email de usuario incorrecto. Se esperaba '%s', se obtuvo: '%s'", expectedUser.Email, user.Email)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	t.Run("GetUserByEmail - not found", func(t *testing.T) {
		email := "notfound@example.com"
		mock.ExpectQuery("SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users").
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByEmail(ctx, email)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("error incorrecto. Se esperaba '%v', se obtuvo '%v'", sql.ErrNoRows, err)
		}
		if user != nil {
			t.Error("se esperaba un usuario nulo, se obtuvo un usuario")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	// --- Prueba de GetUserByUUID ---
	t.Run("GetUserByUUID - success", func(t *testing.T) {
		userUUID := uuid.New().String()
		expectedUser := &models.User{
			ID:           1,
			UUID:         userUUID,
			Email:        "test@example.com",
			PasswordHash: "hashedpassword",
			Role:         "user",
			CreatedAt:    fixedTime,
			UpdatedAt:    fixedTime,
		}
		rows := sqlmock.NewRows([]string{"id", "uuid", "email", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.UUID, expectedUser.Email, expectedUser.PasswordHash, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt)
		mock.ExpectQuery("SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users").
			WithArgs(userUUID).
			WillReturnRows(rows)

		user, err := repo.GetUserByUUID(ctx, userUUID)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo: %v", err)
		}
		if user.UUID != expectedUser.UUID {
			t.Errorf("UUID de usuario incorrecto. Se esperaba '%s', se obtuvo: '%s'", expectedUser.UUID, user.UUID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	t.Run("GetUserByUUID - not found", func(t *testing.T) {
		userUUID := uuid.New().String()
		mock.ExpectQuery("SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users").
			WithArgs(userUUID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByUUID(ctx, userUUID)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("error incorrecto. Se esperaba '%v', se obtuvo '%v'", sql.ErrNoRows, err)
		}
		if user != nil {
			t.Error("se esperaba un usuario nulo, se obtuvo un usuario")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	// --- Prueba de InsertToken ---
	t.Run("InsertToken - success", func(t *testing.T) {
		token := &models.Token{
			UserID:    1,
			Email:     "test@example.com",
			Token:     "test_token",
			TokenHash: "test_token_hash",
			Expiry:    time.Now().Add(time.Hour),
			Role:      "user",
		}
		mock.ExpectExec("INSERT INTO tokens").
			WithArgs(token.UserID, token.Email, token.Token, token.TokenHash, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), token.Role).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.InsertToken(ctx, token)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	t.Run("InsertToken - db error", func(t *testing.T) {
		token := &models.Token{
			UserID: 1,
		}
		expectedErr := errors.New("simulated database error")
		mock.ExpectExec("INSERT INTO tokens").
			WithArgs(token.UserID, token.Email, token.Token, token.TokenHash, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), token.Role).
			WillReturnError(expectedErr)

		err := repo.InsertToken(ctx, token)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("error incorrecto. Se esperaba '%v', se obtuvo '%v'", expectedErr, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	// --- Prueba de GetTokenByTokenHash ---
	t.Run("GetTokenByTokenHash - success", func(t *testing.T) {
		tokenHash := "test_token_hash"
		expectedToken := &models.Token{
			ID:        1,
			UserID:    1,
			Email:     "test@example.com",
			Token:     "test_token",
			TokenHash: tokenHash,
			Expiry:    fixedTime.Add(time.Hour),
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
			Role:      "user",
		}
		rows := sqlmock.NewRows([]string{"id", "user_id", "email", "token", "token_hash", "expiry", "created_at", "updated_at", "role"}).
			AddRow(expectedToken.ID, expectedToken.UserID, expectedToken.Email, expectedToken.Token, expectedToken.TokenHash, expectedToken.Expiry, expectedToken.CreatedAt, expectedToken.UpdatedAt, expectedToken.Role)
		mock.ExpectQuery("SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens").
			WithArgs(tokenHash).
			WillReturnRows(rows)

		token, err := repo.GetTokenByTokenHash(ctx, tokenHash)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo: %v", err)
		}
		if token.TokenHash != expectedToken.TokenHash {
			t.Errorf("token hash incorrecto. Se esperaba '%s', se obtuvo: '%s'", expectedToken.TokenHash, token.TokenHash)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	t.Run("GetTokenByTokenHash - not found", func(t *testing.T) {
		tokenHash := "not_found_hash"
		mock.ExpectQuery("SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens").
			WithArgs(tokenHash).
			WillReturnError(sql.ErrNoRows)

		token, err := repo.GetTokenByTokenHash(ctx, tokenHash)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("error incorrecto. Se esperaba '%v', se obtuvo '%v'", sql.ErrNoRows, err)
		}
		if token != nil {
			t.Error("se esperaba un token nulo, se obtuvo un token")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	// --- Prueba de DeleteTokensByUserID ---
	t.Run("DeleteTokensByUserID - success", func(t *testing.T) {
		userID := 1
		mock.ExpectExec("DELETE FROM tokens").
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteTokensByUserID(ctx, userID)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	t.Run("DeleteTokensByUserID - db error", func(t *testing.T) {
		userID := 1
		expectedErr := errors.New("simulated database error")
		mock.ExpectExec("DELETE FROM tokens").
			WithArgs(userID).
			WillReturnError(expectedErr)

		err := repo.DeleteTokensByUserID(ctx, userID)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("error incorrecto. Se esperaba '%v', se obtuvo '%v'", expectedErr, err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	// --- Prueba de GetTokenByToken ---
	t.Run("GetTokenByToken - success", func(t *testing.T) {
		tokenString := "test_token"
		expectedToken := &models.Token{
			ID:        1,
			UserID:    1,
			Email:     "test@example.com",
			Token:     tokenString,
			TokenHash: "test_token_hash",
			Expiry:    fixedTime.Add(time.Hour),
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
			Role:      "user",
		}
		rows := sqlmock.NewRows([]string{"id", "user_id", "email", "token", "token_hash", "expiry", "created_at", "updated_at", "role"}).
			AddRow(expectedToken.ID, expectedToken.UserID, expectedToken.Email, expectedToken.Token, expectedToken.TokenHash, expectedToken.Expiry, expectedToken.CreatedAt, expectedToken.UpdatedAt, expectedToken.Role)
		mock.ExpectQuery("SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens").
			WithArgs(tokenString).
			WillReturnRows(rows)

		token, err := repo.GetTokenByToken(ctx, tokenString)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo: %v", err)
		}
		if token.Token != expectedToken.Token {
			t.Errorf("token incorrecto. Se esperaba '%s', se obtuvo: '%s'", expectedToken.Token, token.Token)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	t.Run("GetTokenByToken - not found", func(t *testing.T) {
		tokenString := "not_found_token"
		mock.ExpectQuery("SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens").
			WithArgs(tokenString).
			WillReturnError(sql.ErrNoRows)

		token, err := repo.GetTokenByToken(ctx, tokenString)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("error incorrecto. Se esperaba '%v', se obtuvo '%v'", sql.ErrNoRows, err)
		}
		if token != nil {
			t.Error("se esperaba un token nulo, se obtuvo un token")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	// --- Prueba de GetUserForToken ---
	t.Run("GetUserForToken - success", func(t *testing.T) {
		userID := 1
		expectedUser := &models.User{
			ID:           userID,
			UUID:         uuid.New().String(),
			Email:        "test@example.com",
			PasswordHash: "hashedpassword",
			Role:         "user",
			CreatedAt:    fixedTime,
			UpdatedAt:    fixedTime,
		}
		rows := sqlmock.NewRows([]string{"id", "uuid", "email", "password_hash", "role", "created_at", "updated_at"}).
			AddRow(expectedUser.ID, expectedUser.UUID, expectedUser.Email, expectedUser.PasswordHash, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt)
		mock.ExpectQuery("SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users").
			WithArgs(userID).
			WillReturnRows(rows)

		user, err := repo.GetUserForToken(ctx, userID)
		if err != nil {
			t.Errorf("se esperaba un error nulo, se obtuvo: %v", err)
		}
		if user.ID != expectedUser.ID {
			t.Errorf("ID de usuario incorrecto. Se esperaba '%d', se obtuvo: '%d'", expectedUser.ID, user.ID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

	t.Run("GetUserForToken - not found", func(t *testing.T) {
		userID := 2
		mock.ExpectQuery("SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users").
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserForToken(ctx, userID)
		if err == nil {
			t.Error("se esperaba un error, se obtuvo nil")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("error incorrecto. Se esperaba '%v', se obtuvo '%v'", sql.ErrNoRows, err)
		}
		if user != nil {
			t.Error("se esperaba un usuario nulo, se obtuvo un usuario")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("no se cumplieron las expectativas del mock: %s", err)
		}
	})

}
