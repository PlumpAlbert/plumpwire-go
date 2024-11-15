package models

type WG_Client struct {
	ID                   string `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Enabled              bool   `json:"enabled,omitempty"`
	Address              string `json:"address,omitempty"`
	PublicKey            string `json:"publicKey,omitempty"`
	CreatedAt            string `json:"createdAt,omitempty"`
	UpdatedAt            string `json:"updatedAt,omitempty"`
	ExpiredAt            string `json:"expiredAt,omitempty"`
	OneTimeLink          string `json:"oneTimeLink,omitempty"`
	OneTimeLinkExpiresAt string `json:"oneTimeLinkExpiresAt,omitempty"`
	DownloadableConfig   bool   `json:"downloadableConfig,omitempty"`
	PersistentKeepalive  string `json:"persistentKeepalive,omitempty"`
	LatestHandshakeAt    string `json:"latestHandshakeAt,omitempty"`
	TransferRx           int    `json:"transferRx,omitempty"`
	TransferTx           int    `json:"transferTx,omitempty"`
	Endpoint             string `json:"endpoint,omitempty"`

	DeviceName string
}
