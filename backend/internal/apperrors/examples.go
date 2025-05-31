package apperrors

// Примеры использования архитектуры ошибок

/*
Основные способы создания ошибок:

1. Простые предопределенные ошибки:
   err := apperrors.NotFound()
   err := apperrors.InvalidToken()
   err := apperrors.ValidationFailed()
   err := apperrors.Unauthorized()

2. Ошибки с кастомными сообщениями:
   err := apperrors.NotFoundf("User with ID %d not found", userID)
   err := apperrors.InternalErrorf("Failed to connect to database: %v", dbErr)
   err := apperrors.ValidationFailedf("Field %s is invalid", fieldName)

3. Ошибки с дополнительными полями:
   err := apperrors.InvalidParams().WithField("field", "email").WithField("reason", "invalid format")
   err := apperrors.AlreadyExists().WithFields(map[string]interface{}{
       "resource": "user",
       "field": "email",
       "value": email,
   })

4. Изменение существующих ошибок:
   baseErr := apperrors.NotFound()
   customErr := baseErr.WithMessage("Specific resource not found")

5. Использование Builder для сложных ошибок:
   err := apperrors.NewBuilder(CodeValidationFailed, "Validation failed", 400).
       WithField("field", "email").
       WithField("expected", "valid email format").
       WithField("actual", email).
       Build()

6. Создание полностью кастомных ошибок:
   err := apperrors.New(CodeCustom, "Custom error message", 422)
   err := apperrors.Newf(CodeCustom, 422, "Error for user %s: %s", username, reason)

Практические примеры в коде:

// В сервисе пользователей
func (s *UserService) GetByID(id int) (*User, error) {
    user, err := s.repo.GetByID(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, apperrors.NotFoundf("User with ID %d not found", id)
        }
        return nil, apperrors.InternalErrorf("Failed to get user: %v", err)
    }
    return user, nil
}

// В валидации
func (v *Validator) ValidateEmail(email string) error {
    if !isValidEmail(email) {
        return apperrors.ValidationFailed().
            WithField("field", "email").
            WithField("value", email).
            WithField("message", "Invalid email format")
    }
    return nil
}

// В аутентификации
func (a *AuthService) Authenticate(token string) error {
    if token == "" {
        return apperrors.Unauthorized().WithMessage("Token is required")
    }

    if !a.isValidToken(token) {
        return apperrors.InvalidToken().
            WithField("token", token[:min(len(token), 10)]).
            WithField("reason", "expired or malformed")
    }

    return nil
}

// Комплексная валидация с множественными ошибками
func (v *Validator) ValidateUserData(email, password, name string) error {
    var validationErrors = make(map[string]interface{})

    if email == "" {
        validationErrors["email"] = "Email is required"
    } else if !isValidEmail(email) {
        validationErrors["email"] = "Invalid email format"
    }

    if len(password) < 8 {
        validationErrors["password"] = map[string]interface{}{
            "message": "Password is too short",
            "min_length": 8,
            "actual_length": len(password),
        }
    }

    if len(validationErrors) > 0 {
        return apperrors.ValidationFailed().
            WithMessage("User data validation failed").
            WithField("validation_errors", validationErrors)
    }

    return nil
}

// HTTP обработка ошибок
func HandleUserRequest(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        HandleError(c, apperrors.InvalidParamsf("Invalid user ID: %s", c.Param("id")))
        return
    }

    user, err := userService.GetByID(userID)
    if err != nil {
        HandleError(c, err)
        return
    }

    c.JSON(200, user)
}

*/