package permissions

const (
	PermissionMerchantManage = "merchant:manage"
	PermissionMerchantRead   = "merchant:read"
	PermissionStoreManage    = "store:manage"
	PermissionStoreRead      = "store:read"
	PermissionStaffManage    = "staff:manage"
	PermissionBillingRead    = "billing:read"
	PermissionBillingManage  = "billing:manage"
	PermissionCreatorLink    = "creator:link"
)

var MerchantAdminPermissions = []string{
	PermissionMerchantManage,
	PermissionMerchantRead,
	PermissionStoreManage,
	PermissionStoreRead,
	PermissionStaffManage,
	PermissionBillingRead,
	PermissionBillingManage,
}

var StoreManagerPermissions = []string{
	PermissionStoreManage,
	PermissionStoreRead,
}
