package domain

type P2pPeerRegistryCmd struct {
	MultiAddr   string `form:"multiAddr" binding:"required" json:"multiAddr"`
	OneTimePass string `form:"oneTimePass" binding:"required" json:"oneTimePass"`
}
