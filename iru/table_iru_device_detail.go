package iru

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TYPES

// deviceDetailAssignedUser is the shape of the assigned_user nested in the general section.
type deviceDetailAssignedUser struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Active     bool   `json:"active"`
	IsArchived bool   `json:"is_archived"`
}

// deviceDetailGeneral is the general section of the device detail response.
type deviceDetailGeneral struct {
	DeviceID                    string                    `json:"device_id"`
	DeviceName                  string                    `json:"device_name"`
	Platform                    string                    `json:"platform"`
	FirstEnrollment             string                    `json:"first_enrollment"`
	LastEnrollment              string                    `json:"last_enrollment"`
	AssignedUser                *deviceDetailAssignedUser `json:"assigned_user"`
	BlueprintName               string                    `json:"blueprint_name"`
	BlueprintUUID               string                    `json:"blueprint_uuid"`
	AssetTag                    string                    `json:"asset_tag"`
	Model                       string                    `json:"model"`
	OSVersion                   string                    `json:"os_version"`
	SupplementalBuildVersion    string                    `json:"supplemental_build_version"`
	SupplementalOSVersionExtra  string                    `json:"supplemental_os_version_extra"`
	SystemVersion               string                    `json:"system_version"`
	BootVolume                  string                    `json:"boot_volume"`
	TimeSinceBoot               string                    `json:"time_since_boot"`
	LastUser                    string                    `json:"last_user"`
}

// deviceDetailMDM is the mdm section of the device detail response.
// Note: mdm_enabled and supervised are returned as strings "True"/"False".
type deviceDetailMDM struct {
	MDMEnabled     string   `json:"mdm_enabled"`
	Supervised     string   `json:"supervised"`
	InstallDate    string   `json:"install_date"`
	LastCheckIn    string   `json:"last_check_in"`
	MDMEnabledUser []string `json:"mdm_enabled_user"`
}

// deviceDetailActivationLock is the activation_lock section.
type deviceDetailActivationLock struct {
	BypassCodeFailed                       bool `json:"bypass_code_failed"`
	UserActivationLockEnabled              bool `json:"user_activation_lock_enabled"`
	DeviceActivationLockEnabled            bool `json:"device_activation_lock_enabled"`
	ActivationLockAllowedWhileSupervised   bool `json:"activation_lock_allowed_while_supervised"`
	ActivationLockSupported                bool `json:"activation_lock_supported"`
}

// deviceDetailFileVault is the filevault section.
type deviceDetailFileVault struct {
	FileVaultEnabled        bool   `json:"filevault_enabled"`
	FileVaultRecoveryKeyType string `json:"filevault_recoverykey_type"`
	FileVaultPRKEscrowed    bool   `json:"filevault_prk_escrowed"`
	FileVaultNextRotation   string `json:"filevault_next_rotation"`
	FileVaultRegenRequired  bool   `json:"filevault_regen_required"`
}

// deviceDetailKandjiAgent is the kandji_agent section.
type deviceDetailKandjiAgent struct {
	AgentInstalled string `json:"agent_installed"`
	InstallDate    string `json:"install_date"`
	LastCheckIn    string `json:"last_check_in"`
	AgentVersion   string `json:"agent_version"`
}

// deviceDetailHardware is the hardware_overview section.
type deviceDetailHardware struct {
	ModelName          string `json:"model_name"`
	ModelIdentifier    string `json:"model_identifier"`
	ProcessorName      string `json:"processor_name"`
	ProcessorSpeed     string `json:"processor_speed"`
	NumberOfProcessors string `json:"number_of_processors"`
	TotalCores         string `json:"total_number_of_cores"`
	Memory             string `json:"memory"`
	BatteryHealth      string `json:"battery_health"`
	UDID               string `json:"udid"`
	SerialNumber       string `json:"serial_number"`
}

// deviceDetailADE is the automated_device_enrollment section.
type deviceDetailADE struct {
	AutoEnrollEligible bool `json:"auto_enroll_eligible"`
	AutoEnrolled       bool `json:"auto_enrolled"`
}

// DeviceDetail is the full response from GET /api/v1/devices/{id}/details.
type DeviceDetail struct {
	General          *deviceDetailGeneral        `json:"general"`
	MDM              *deviceDetailMDM             `json:"mdm"`
	ActivationLock   *deviceDetailActivationLock  `json:"activation_lock"`
	FileVault        *deviceDetailFileVault        `json:"filevault"`
	KandjiAgent      *deviceDetailKandjiAgent      `json:"kandji_agent"`
	HardwareOverview *deviceDetailHardware         `json:"hardware_overview"`
	ADE              *deviceDetailADE              `json:"automated_device_enrollment"`
	Volumes          interface{}                    `json:"volumes"`
}

//// TABLE DEFINITION

func tableIruDeviceDetail() *plugin.Table {
	return &plugin.Table{
		Name:        "iru_device_detail",
		Description: "Rich device details including MDM, FileVault, hardware, and agent status.",
		List: &plugin.ListConfig{
			Hydrate: listDeviceDetails,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "device_id", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// General
			{Name: "device_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.DeviceID"), Description: "Unique identifier for the device."},
			{Name: "device_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.DeviceName"), Description: "Name of the device."},
			{Name: "platform", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.Platform"), Description: "Device platform (Mac, iPhone, iPad, AppleTV)."},
			{Name: "model", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.Model"), Description: "Hardware model name."},
			{Name: "os_version", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.OSVersion"), Description: "Operating system version."},
			{Name: "system_version", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.SystemVersion"), Description: "Full system version string."},
			{Name: "supplemental_build_version", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.SupplementalBuildVersion"), Description: "Supplemental OS build version."},
			{Name: "asset_tag", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.AssetTag"), Description: "Asset tag assigned to the device."},
			{Name: "blueprint_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.BlueprintName"), Description: "Name of the assigned blueprint."},
			{Name: "blueprint_uuid", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.BlueprintUUID"), Description: "UUID of the assigned blueprint."},
			{Name: "boot_volume", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.BootVolume"), Description: "Name of the boot volume."},
			{Name: "time_since_boot", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.TimeSinceBoot"), Description: "Time elapsed since last boot."},
			{Name: "last_user", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.LastUser"), Description: "Username of the last logged-in user."},
			{Name: "first_enrollment", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.FirstEnrollment"), Description: "Timestamp of first enrollment."},
			{Name: "last_enrollment", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.LastEnrollment"), Description: "Timestamp of last enrollment."},

			// Assigned user (from General)
			{Name: "assigned_user_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.AssignedUser.ID"), Description: "ID of the assigned user."},
			{Name: "assigned_user_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.AssignedUser.Name"), Description: "Name of the assigned user."},
			{Name: "assigned_user_email", Type: proto.ColumnType_STRING, Transform: transform.FromField("General.AssignedUser.Email"), Description: "Email of the assigned user."},

			// MDM
			{Name: "mdm_enabled", Type: proto.ColumnType_STRING, Transform: transform.FromField("MDM.MDMEnabled"), Description: "Whether MDM is enabled (True/False)."},
			{Name: "supervised", Type: proto.ColumnType_STRING, Transform: transform.FromField("MDM.Supervised"), Description: "Whether the device is supervised (True/False)."},
			{Name: "mdm_install_date", Type: proto.ColumnType_STRING, Transform: transform.FromField("MDM.InstallDate"), Description: "Date MDM was installed."},
			{Name: "mdm_last_check_in", Type: proto.ColumnType_STRING, Transform: transform.FromField("MDM.LastCheckIn"), Description: "Last MDM check-in timestamp."},
			{Name: "mdm_enabled_user", Type: proto.ColumnType_JSON, Transform: transform.FromField("MDM.MDMEnabledUser"), Description: "Users with MDM enabled."},

			// Activation Lock
			{Name: "activation_lock_supported", Type: proto.ColumnType_BOOL, Transform: transform.FromField("ActivationLock.ActivationLockSupported"), Description: "Whether Activation Lock is supported."},
			{Name: "user_activation_lock_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("ActivationLock.UserActivationLockEnabled"), Description: "Whether user Activation Lock is enabled."},
			{Name: "device_activation_lock_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("ActivationLock.DeviceActivationLockEnabled"), Description: "Whether device Activation Lock is enabled."},

			// FileVault
			{Name: "filevault_enabled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("FileVault.FileVaultEnabled"), Description: "Whether FileVault encryption is enabled."},
			{Name: "filevault_recovery_key_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("FileVault.FileVaultRecoveryKeyType"), Description: "FileVault recovery key type (Personal, Institutional, etc.)."},
			{Name: "filevault_prk_escrowed", Type: proto.ColumnType_BOOL, Transform: transform.FromField("FileVault.FileVaultPRKEscrowed"), Description: "Whether the personal recovery key is escrowed."},
			{Name: "filevault_next_rotation", Type: proto.ColumnType_STRING, Transform: transform.FromField("FileVault.FileVaultNextRotation"), Description: "Next scheduled FileVault key rotation."},

			// Kandji Agent
			{Name: "agent_installed", Type: proto.ColumnType_STRING, Transform: transform.FromField("KandjiAgent.AgentInstalled"), Description: "Whether the Kandji agent is installed (True/False)."},
			{Name: "agent_version", Type: proto.ColumnType_STRING, Transform: transform.FromField("KandjiAgent.AgentVersion"), Description: "Installed Kandji agent version."},
			{Name: "agent_install_date", Type: proto.ColumnType_STRING, Transform: transform.FromField("KandjiAgent.InstallDate"), Description: "Date the Kandji agent was installed."},
			{Name: "agent_last_check_in", Type: proto.ColumnType_STRING, Transform: transform.FromField("KandjiAgent.LastCheckIn"), Description: "Last Kandji agent check-in timestamp."},

			// Hardware
			{Name: "processor_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("HardwareOverview.ProcessorName"), Description: "Processor name (e.g. Apple M4 Pro)."},
			{Name: "total_cores", Type: proto.ColumnType_STRING, Transform: transform.FromField("HardwareOverview.TotalCores"), Description: "Total number of CPU cores."},
			{Name: "memory", Type: proto.ColumnType_STRING, Transform: transform.FromField("HardwareOverview.Memory"), Description: "Total RAM (e.g. 24 GB LPDDR5)."},
			{Name: "battery_health", Type: proto.ColumnType_STRING, Transform: transform.FromField("HardwareOverview.BatteryHealth"), Description: "Battery health status."},
			{Name: "serial_number", Type: proto.ColumnType_STRING, Transform: transform.FromField("HardwareOverview.SerialNumber"), Description: "Device serial number."},
			{Name: "udid", Type: proto.ColumnType_STRING, Transform: transform.FromField("HardwareOverview.UDID"), Description: "Device UDID."},
			{Name: "model_identifier", Type: proto.ColumnType_STRING, Transform: transform.FromField("HardwareOverview.ModelIdentifier"), Description: "Hardware model identifier (e.g. Mac16,8)."},

			// ADE
			{Name: "auto_enrolled", Type: proto.ColumnType_BOOL, Transform: transform.FromField("ADE.AutoEnrolled"), Description: "Whether the device was auto-enrolled via ADE."},
			{Name: "auto_enroll_eligible", Type: proto.ColumnType_BOOL, Transform: transform.FromField("ADE.AutoEnrollEligible"), Description: "Whether the device is eligible for ADE auto-enrollment."},

			// Volumes
			{Name: "volumes", Type: proto.ColumnType_JSON, Description: "List of storage volumes on the device."},
		},
	}
}

//// HYDRATE FUNCTIONS

func listDeviceDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		return nil, err
	}

	// If a specific device_id is requested, only fetch details for that device.
	if deviceID := d.EqualsQualString("device_id"); deviceID != "" {
		var detail DeviceDetail
		if err := client.get(ctx, fmt.Sprintf("/api/v1/devices/%s/details", deviceID), nil, &detail); err != nil {
			return nil, fmt.Errorf("getting device details for %s: %w", deviceID, err)
		}
		d.StreamListItem(ctx, detail)
		return nil, nil
	}

	// Otherwise, fetch all device IDs first, then details for each.
	deviceIDs, err := getAllDeviceIDs(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("listing device IDs: %w", err)
	}

	for _, deviceID := range deviceIDs {
		var detail DeviceDetail
		if err := client.get(ctx, fmt.Sprintf("/api/v1/devices/%s/details", deviceID), nil, &detail); err != nil {
			plugin.Logger(ctx).Warn("iru_device_detail: failed to fetch details", "device_id", deviceID, "error", err)
			continue
		}
		d.StreamListItem(ctx, detail)
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
