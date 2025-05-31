package permission

const (
	// Профиль пользователя
	ReadOwnProfile   = "read_own_profile"
	UpdateOwnProfile = "update_own_profile"

	// Работа с внешними приложениями
	ReadOwnExternalApps   = "read_own_external_apps"
	UpdateOwnExternalApps = "update_own_external_apps"

	// Работа с устройствами
	RegisterDevice = "register_device"
	RemoveDevice   = "remove_device"
	ReadOwnDevices = "read_own_devices"

	// Работа с данными о здоровье (только свои)
	ReadOwnSteps         = "read_own_steps"
	WriteOwnSteps        = "write_own_steps"
	ReadOwnRRIntervals    = "read_own_rr_intervals"
	WriteOwnRRIntervals   = "write_own_rr_intervals"
	ReadOwnTemperatures  = "read_own_temperatures"
	WriteOwnTemperatures = "write_own_temperatures"
	ReadOwnWeights       = "read_own_weights"
	WriteOwnWeights      = "write_own_weights"
	ReadOwnSleeps        = "read_own_sleeps"
	WriteOwnSleeps       = "write_own_sleeps"

	// Административные
	ReadAllUsers      = "read_all_users"
	DeleteUser        = "delete_user"
	AssignRoles       = "assign_roles"
	ReadAllDevices    = "read_all_devices"
	ReadAllHealthData = "read_all_health_data"
)
