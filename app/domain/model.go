package domain

type P2pPeerRegistryCmd struct {
	MultiAddr string `binding:"required" json:"multiAddr"`
	Otp       string `binding:"required" json:"otp"`
	Mode      string `binding:"required" json:"mode"`
	Token     string `binding:"required" json:"token"`
}

type P2pPeerDeregisterCmd struct {
	Otp   string `form:"otp" binding:"required" json:"otp"`
	Token string `form:"otp" binding:"required" json:"token"`
}

type P2pBootstrapNodeRegistryCmd struct {
	NodeId string `json:"nodeId"`
	Port   int    `json:"port"`
}
