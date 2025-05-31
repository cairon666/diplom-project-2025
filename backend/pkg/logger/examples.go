package logger

/*
=== 2. В продакшн коде ===

func main() {
    // Создаем настоящий логгер
    logger, err := logger.NewProd()
    if err != nil {
        panic(err)
    }

    // Используем интерфейс
    userService := NewUserService(db, logger)
    authService := NewAuthService(userService, logger)
}

=== 3. В тестах ===

func TestUserService_GetByID(t *testing.T) {
    // Создаем mock логгер
    mockLogger := logger.NewMockLogger()

    // Создаем сервис с mock логгером
    service := NewUserService(mockDB, mockLogger)

    // Тестируем
    user, err := service.GetByID(123)

    // Проверяем логи
    assert.Equal(t, 1, mockLogger.GetCallsCount())
    assert.True(t, mockLogger.HasCallWithMessage("Getting user by ID"))

    calls := mockLogger.GetCallsByLevel("INFO")
    assert.Len(t, calls, 1)
}

=== 4. Использование с контекстом ===

func (s *UserService) GetByID(ctx context.Context, id int64) (*User, error) {
    // Создаем логгер с контекстом
    log := logger.WithContextFields(s.logger, ctx)

    log.Info("Getting user by ID", logger.WithUserID(id))

    user, err := s.repo.GetByID(id)
    if err != nil {
        log.Error("Failed to get user",
            logger.WithError(err),
            logger.WithUserID(id),
        )
        return nil, err
    }

    log.Info("User retrieved successfully", logger.WithUserID(id))
    return user, nil
}

=== 5. Использование LoggedService ===

type UserService struct {
    *logger.LoggedService
    repo UserRepository
}

func NewUserService(repo UserRepository, logger logger.ILogger) *UserService {
    return &UserService{
        LoggedService: logger.NewLoggedService(logger, "user_service"),
        repo:          repo,
    }
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    s.LogInfo(ctx, "Creating new user", logger.String("email", user.Email))

    if err := s.repo.Create(user); err != nil {
        s.LogError(ctx, err, "Failed to create user",
            logger.String("email", user.Email),
        )
        return err
    }

    s.LogInfo(ctx, "User created successfully",
        logger.String("email", user.Email),
        logger.WithUserID(user.ID),
    )
    return nil
}

=== 6. HTTP Middleware ===

func main() {
    logger, _ := logger.NewProd()

    r := gin.New()
    r.Use(logger.GinMiddleware(logger))

    r.GET("/users/:id", getUserHandler)
    r.Run()
}

=== 7. Интеграция с ошибками ===

func (s *UserService) GetByID(ctx context.Context, id int64) (*User, error) {
    log := s.WithContext(ctx)

    user, err := s.repo.GetByID(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            // Логируем как INFO, так как это ожидаемая ситуация
            log.Info("User not found", logger.WithUserID(id))
            return nil, apperrors.NotFoundf("User with ID %d not found", id)
        }

        // Логируем как ERROR для неожиданных ошибок
        log.Error("Database error while getting user",
            logger.WithError(err),
            logger.WithUserID(id),
        )
        return nil, apperrors.InternalErrorf("Failed to get user: %v", err)
    }

    log.Info("User retrieved successfully", logger.WithUserID(id))
    return user, nil
}

=== 8. Тестирование с проверкой логов ===

func TestUserService_GetByID_NotFound(t *testing.T) {
    mockLogger := logger.NewMockLogger()
    mockRepo := &MockUserRepository{}

    service := &UserService{
        LoggedService: logger.NewLoggedService(mockLogger, "user_service"),
        repo:          mockRepo,
    }

    // Настраиваем mock
    mockRepo.On("GetByID", int64(123)).Return(nil, sql.ErrNoRows)

    // Тестируем
    ctx := context.Background()
    user, err := service.GetByID(ctx, 123)

    // Проверяем результат
    assert.Nil(t, user)
    assert.True(t, apperrors.IsNotFound(err))

    // Проверяем логи
    calls := mockLogger.GetCalls()
    assert.Len(t, calls, 1)
    assert.Equal(t, "INFO", calls[0].Level)
    assert.Equal(t, "User not found", calls[0].Message)
}

=== 9. No-op логгер для тестов ===

func TestSomethingQuiet(t *testing.T) {
    // Используем no-op логгер когда логи не важны
    service := NewUserService(mockDB, &logger.NoOpLogger{})

    // Тестируем без шума в логах
    result := service.DoSomething()
    assert.NotNil(t, result)
}

=== 10. Structured logging с полями ===

func (s *AuthService) Login(ctx context.Context, email, password string) (*Token, error) {
    log := s.WithContext(ctx).With(
        logger.String("email", email),
        logger.WithOperation("login"),
    )

    log.Info("Login attempt started")

    user, err := s.userService.GetByEmail(ctx, email)
    if err != nil {
        if apperrors.IsNotFound(err) {
            log.Warn("Login attempt with non-existent email")
            return nil, apperrors.LoginNotRegistered()
        }
        log.Error("Failed to get user during login", logger.WithError(err))
        return nil, err
    }

    if !s.checkPassword(password, user.PasswordHash) {
        log.Warn("Login attempt with wrong password",
            logger.WithUserID(user.ID),
        )
        return nil, apperrors.WrongPassword()
    }

    token, err := s.generateToken(user)
    if err != nil {
        log.Error("Failed to generate token",
            logger.WithError(err),
            logger.WithUserID(user.ID),
        )
        return nil, apperrors.InternalError()
    }

    log.Info("Login successful",
        logger.WithUserID(user.ID),
        logger.String("token_id", token.ID),
    )

    return token, nil
}

*/