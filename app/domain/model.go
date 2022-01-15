package domain

type P2pPeerRegistryCmd struct {
	MultiAddr string `json:"multiAddr"`
	Otp       string `json:"otp"`
	Mode      string `json:"mode"`
}

type P2pBootstrapNodeRegistryCmd struct {
	NodeId string `json:"nodeId"`
	Port   int    `json:"port"`
}
