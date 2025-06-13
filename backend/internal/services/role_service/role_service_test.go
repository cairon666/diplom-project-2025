package role_service_test

import (
	"testing"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/models"
	"github.com/cairon666/vkr-backend/internal/services/role_service"
	"github.com/cairon666/vkr-backend/internal/services/role_service/mocks"
	"github.com/cairon666/vkr-backend/pkg/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createServiceWithMock(t *testing.T) (*role_service.RoleService, *mocks.MockRoleRepos) {
	t.Helper()
	mockRepo := mocks.NewMockRoleRepos(t)
	service := role_service.NewRoleService(mockRepo)

	return service, mockRepo
}

func setupGetRoleByNameMock(mockRepo *mocks.MockRoleRepos, roleName string, role models.Role, err error) {
	mockRepo.
		EXPECT().
		GetRoleByName(mock.Anything, roleName).
		Return(role, err)
}

func TestNewRoleService(t *testing.T) {
	t.Parallel()
	service, mockRepo := createServiceWithMock(t)

	assert.NotNil(t, service)
	assert.NotNil(t, mockRepo)
}

type hasPermissionTestCase struct {
	name           string
	userID         uuid.UUID
	permission     string
	repoPerms      []models.Permission
	repoError      error
	expectedResult bool
	expectedError  error
}

func createHasPermissionTestCases() []hasPermissionTestCase {
	return []hasPermissionTestCase{
		{
			name:       "user has permission",
			userID:     uuid.New(),
			permission: "read_users",
			repoPerms: []models.Permission{
				{ID: 1, Name: "read_users"},
				{ID: 2, Name: "write_users"},
			},
			repoError:      nil,
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:       "user does not have permission",
			userID:     uuid.New(),
			permission: "delete_users",
			repoPerms: []models.Permission{
				{ID: 1, Name: "read_users"},
				{ID: 2, Name: "write_users"},
			},
			repoError:      nil,
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "empty permissions",
			userID:         uuid.New(),
			permission:     "read_users",
			repoPerms:      []models.Permission{},
			repoError:      nil,
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "repository error",
			userID:         uuid.New(),
			permission:     "read_users",
			repoPerms:      nil,
			repoError:      apperrors.InternalErrorf("database error"),
			expectedResult: false,
			expectedError:  apperrors.InternalErrorf("database error"),
		},
	}
}

func TestRoleService_HasPermission(t *testing.T) {
	t.Parallel()
	tests := createHasPermissionTestCases()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			service, mockRepo := createServiceWithMock(t)

			mockRepo.
				EXPECT().
				GetPermissionsByUserID(mock.Anything, testCase.userID).
				Return(testCase.repoPerms, testCase.repoError)

			result, err := service.HasPermission(t.Context(), testCase.userID, testCase.permission)

			assert.Equal(t, testCase.expectedResult, result)
			testutils.AssertError(t, testCase.expectedError, err)
		})
	}
}

type hasOneOfPermissionsTestCase struct {
	name           string
	userID         uuid.UUID
	permissions    []string
	repoPerms      []models.Permission
	repoError      error
	expectedResult bool
	expectedError  error
}

func createHasOneOfPermissionsTestCases() []hasOneOfPermissionsTestCase {
	return []hasOneOfPermissionsTestCase{
		{
			name:        "user has one of permissions",
			userID:      uuid.New(),
			permissions: []string{"read_users", "write_users", "delete_users"},
			repoPerms: []models.Permission{
				{ID: 1, Name: "read_users"},
				{ID: 2, Name: "create_posts"},
			},
			repoError:      nil,
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:        "user has multiple of permissions",
			userID:      uuid.New(),
			permissions: []string{"read_users", "write_users"},
			repoPerms: []models.Permission{
				{ID: 1, Name: "read_users"},
				{ID: 2, Name: "write_users"},
			},
			repoError:      nil,
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:        "user has none of permissions",
			userID:      uuid.New(),
			permissions: []string{"delete_users", "admin_access"},
			repoPerms: []models.Permission{
				{ID: 1, Name: "read_users"},
				{ID: 2, Name: "write_users"},
			},
			repoError:      nil,
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "empty permissions list",
			userID:         uuid.New(),
			permissions:    []string{},
			repoPerms:      []models.Permission{{ID: 1, Name: "read_users"}},
			repoError:      nil,
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "repository error",
			userID:         uuid.New(),
			permissions:    []string{"read_users"},
			repoPerms:      nil,
			repoError:      apperrors.InternalErrorf("database error"),
			expectedResult: false,
			expectedError:  apperrors.InternalErrorf("database error"),
		},
	}
}

func TestRoleService_HasOneOfPermissions(t *testing.T) {
	t.Parallel()
	tests := createHasOneOfPermissionsTestCases()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			service, mockRepo := createServiceWithMock(t)

			mockRepo.
				EXPECT().
				GetPermissionsByUserID(mock.Anything, testCase.userID).
				Return(testCase.repoPerms, testCase.repoError)

			result, err := service.HasOneOfPermissions(t.Context(), testCase.userID, testCase.permissions)

			assert.Equal(t, testCase.expectedResult, result)
			testutils.AssertError(t, testCase.expectedError, err)
		})
	}
}

type assignRoleTestCase struct {
	name          string
	userID        uuid.UUID
	roleName      string
	role          models.Role
	getRoleError  error
	assignError   error
	expectedError error
}

func createAssignRoleTestCases() []assignRoleTestCase {
	return []assignRoleTestCase{
		{
			name:          "successful assignment",
			userID:        uuid.New(),
			roleName:      "admin",
			role:          models.Role{ID: 1, Name: "admin"},
			getRoleError:  nil,
			assignError:   nil,
			expectedError: nil,
		},
		{
			name:          "role not found",
			userID:        uuid.New(),
			roleName:      "nonexistent",
			role:          models.Role{},
			getRoleError:  apperrors.NotFound(),
			assignError:   nil,
			expectedError: apperrors.NotFound(),
		},
		{
			name:          "assign role error",
			userID:        uuid.New(),
			roleName:      "admin",
			role:          models.Role{ID: 1, Name: "admin"},
			getRoleError:  nil,
			assignError:   apperrors.InternalErrorf("database error"),
			expectedError: apperrors.InternalErrorf("database error"),
		},
	}
}

func setupAssignRoleToUserMocks(mockRepo *mocks.MockRoleRepos, testCase assignRoleTestCase) {
	setupGetRoleByNameMock(mockRepo, testCase.roleName, testCase.role, testCase.getRoleError)
	if testCase.getRoleError == nil {
		mockRepo.
			EXPECT().
			AssignRoleToUser(mock.Anything, testCase.userID, testCase.role.ID).
			Return(testCase.assignError)
	}
}

func TestRoleService_AssignRoleToUser(t *testing.T) {
	t.Parallel()
	tests := createAssignRoleTestCases()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			service, mockRepo := createServiceWithMock(t)

			setupAssignRoleToUserMocks(mockRepo, testCase)

			err := service.AssignRoleToUser(t.Context(), testCase.userID, testCase.roleName)

			testutils.AssertError(t, testCase.expectedError, err)
		})
	}
}

func TestRoleService_GetPermissionsByExternalAppUserID(t *testing.T) {
	t.Parallel()
	userID := uuid.New()
	expectedPerms := []models.Permission{
		{ID: 1, Name: "read_data"},
		{ID: 2, Name: "write_data"},
	}

	service, mockRepo := createServiceWithMock(t)

	mockRepo.EXPECT().GetPermissionsByExternalAppID(mock.Anything, userID).Return(expectedPerms, nil)

	result, err := service.GetPermissionsByExternalAppUserID(t.Context(), userID)

	require.NoError(t, err)
	assert.Equal(t, expectedPerms, result)
}

func TestRoleService_GetRolesByExternalAppID(t *testing.T) {
	t.Parallel()
	userID := uuid.New()
	expectedRoles := []models.Role{
		{ID: 1, Name: "external_reader"},
		{ID: 2, Name: "external_writer"},
	}

	service, mockRepo := createServiceWithMock(t)

	mockRepo.EXPECT().GetRolesByExternalAppID(mock.Anything, userID).Return(expectedRoles, nil)

	result, err := service.GetRolesByExternalAppID(t.Context(), userID)

	require.NoError(t, err)
	assert.Equal(t, expectedRoles, result)
}

type assignRoleToExternalAppTestCase struct {
	name          string
	externalID    uuid.UUID
	roleName      string
	role          models.Role
	getRoleError  error
	assignError   error
	expectedError error
}

func createAssignRoleToExternalAppTestCases() []assignRoleToExternalAppTestCase {
	return []assignRoleToExternalAppTestCase{
		{
			name:          "successful assignment",
			externalID:    uuid.New(),
			roleName:      "external_reader",
			role:          models.Role{ID: 1, Name: "external_reader"},
			getRoleError:  nil,
			assignError:   nil,
			expectedError: nil,
		},
		{
			name:          "role not found",
			externalID:    uuid.New(),
			roleName:      "nonexistent",
			role:          models.Role{},
			getRoleError:  apperrors.NotFound(),
			assignError:   nil,
			expectedError: apperrors.NotFound(),
		},
		{
			name:          "assign role error",
			externalID:    uuid.New(),
			roleName:      "external_reader",
			role:          models.Role{ID: 1, Name: "external_reader"},
			getRoleError:  nil,
			assignError:   apperrors.InternalErrorf("database error"),
			expectedError: apperrors.InternalErrorf("database error"),
		},
	}
}

func setupAssignRoleToExternalAppMocks(mockRepo *mocks.MockRoleRepos, testCase assignRoleToExternalAppTestCase) {
	setupGetRoleByNameMock(mockRepo, testCase.roleName, testCase.role, testCase.getRoleError)
	if testCase.getRoleError == nil {
		mockRepo.
			EXPECT().
			AssignRoleToExternalApp(mock.Anything, testCase.externalID, testCase.role.ID).
			Return(testCase.assignError)
	}
}

func TestRoleService_AssignRoleToExternalApp(t *testing.T) {
	t.Parallel()
	tests := createAssignRoleToExternalAppTestCases()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			service, mockRepo := createServiceWithMock(t)

			setupAssignRoleToExternalAppMocks(mockRepo, testCase)

			err := service.AssignRoleToExternalApp(t.Context(), testCase.externalID, testCase.roleName)

			testutils.AssertError(t, testCase.expectedError, err)
		})
	}
}

type assignRolesToExternalAppTestCase struct {
	name           string
	externalID     uuid.UUID
	roleNames      []string
	roles          []models.Role
	getRoleErrors  []error
	assignError    error
	expectedError  error
	expectedErrMsg string
}

func createAssignRolesToExternalAppTestCases() []assignRolesToExternalAppTestCase {
	return []assignRolesToExternalAppTestCase{
		{
			name:       "successful assignment",
			externalID: uuid.New(),
			roleNames:  []string{"external_reader", "external_writer"},
			roles: []models.Role{
				{ID: 1, Name: "external_reader"},
				{ID: 2, Name: "external_writer"},
			},
			getRoleErrors: []error{nil, nil},
			assignError:   nil,
			expectedError: nil,
		},
		{
			name:       "first role not found",
			externalID: uuid.New(),
			roleNames:  []string{"nonexistent", "external_writer"},
			roles: []models.Role{
				{},
				{ID: 2, Name: "external_writer"},
			},
			getRoleErrors:  []error{apperrors.NotFound(), nil},
			assignError:    nil,
			expectedError:  apperrors.DataProcessingErrorf("failed to get role nonexistent:"),
			expectedErrMsg: "failed to get role nonexistent:",
		},
		{
			name:       "second role not found",
			externalID: uuid.New(),
			roleNames:  []string{"external_reader", "nonexistent"},
			roles: []models.Role{
				{ID: 1, Name: "external_reader"},
				{},
			},
			getRoleErrors:  []error{nil, apperrors.NotFound()},
			assignError:    nil,
			expectedError:  apperrors.DataProcessingErrorf("failed to get role nonexistent:"),
			expectedErrMsg: "failed to get role nonexistent:",
		},
		{
			name:       "assign roles error",
			externalID: uuid.New(),
			roleNames:  []string{"external_reader", "external_writer"},
			roles: []models.Role{
				{ID: 1, Name: "external_reader"},
				{ID: 2, Name: "external_writer"},
			},
			getRoleErrors:  []error{nil, nil},
			assignError:    apperrors.InternalErrorf("database error"),
			expectedError:  apperrors.DataProcessingErrorf("failed to assign roles:"),
			expectedErrMsg: "failed to assign roles:",
		},
	}
}

func setupMocksForAssignRolesToExternalApp(t *testing.T, testCase assignRolesToExternalAppTestCase) *mocks.MockRoleRepos {
	t.Helper()
	_, mockRepo := createServiceWithMock(t)

	hasError := false
	for i, roleName := range testCase.roleNames {
		mockRepo.EXPECT().GetRoleByName(mock.Anything, roleName).Return(testCase.roles[i], testCase.getRoleErrors[i])
		if testCase.getRoleErrors[i] != nil {
			hasError = true

			break
		}
	}

	if !hasError {
		expectedRoleIDs := make([]int32, len(testCase.roles))
		for i, role := range testCase.roles {
			expectedRoleIDs[i] = role.ID
		}
		mockRepo.EXPECT().AssignRolesToExternalApp(mock.Anything, testCase.externalID, expectedRoleIDs).Return(testCase.assignError)
	}

	return mockRepo
}

func TestRoleService_AssignRolesToExternalApp(t *testing.T) {
	t.Parallel()
	tests := createAssignRolesToExternalAppTestCases()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			mockRepo := setupMocksForAssignRolesToExternalApp(t, testCase)
			service := role_service.NewRoleService(mockRepo)

			err := service.AssignRolesToExternalApp(t.Context(), testCase.externalID, testCase.roleNames)

			if testCase.expectedError != nil {
				require.Error(t, err)
				if testCase.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), testCase.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
