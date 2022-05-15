package entities

type Model_domain struct {
	Domain string `json:"domain"`
}

type Controller_clientSaveTogel struct {
	Client_Username string `json:"client_username"`
	Client_Company  string `json:"client_company"`
	Idtrxkeluaran   string `json:"idtrxkeluaran"`
	Idcomppasaran   string `json:"idcomppasaran"`
	Pasarancode     string `json:"pasarancode"`
	Pasaranperiode  string `json:"pasaranperiode"`
	Devicemember    string `json:"devicemember"`
	Formipaddress   string `json:"formipaddress"`
	Timezone        string `json:"timezone"`
	Totalbayarbet   int    `json:"totalbayarbet"`
	List4d          string `json:"list4d"`
	Hostname        string `json:"hostname"`
	Token           string `json:"token"`
}
