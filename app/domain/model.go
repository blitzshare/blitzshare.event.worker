package domain

type P2pPeerRegistryCmd struct {
	MultiAddr   string `json:"multiAddr"`
	OneTimePass string `json:"oneTimePass"`
}

type P2pBootstrapNodeRegistryCmd struct {
	NodeId string `json:"nodeId"`
	Port   int    `json:"port"`
}
