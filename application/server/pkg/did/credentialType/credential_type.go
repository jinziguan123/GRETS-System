package credentialType

type CredentialType string

const (
	CredentialTypeIdentity     CredentialType = "IdentityCredential"
	CredentialTypeOrganization CredentialType = "OrganizationCredential"
	CredentialTypeRole         CredentialType = "RoleCredential"
	CredentialTypeAsset        CredentialType = "AssetCredential"
)
