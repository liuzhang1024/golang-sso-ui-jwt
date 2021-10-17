package ssojwt

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func ValidatTicket(config SSOConfig, ticket string) (bodyBytes []byte, err error) {

	url := fmt.Sprintf("%sserviceValidate?ticket=%s&service=%s", config.CasURL, ticket, config.ServiceUrl)
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	return
}

func Unmarshal(bodyBytes []byte) (model ServiceResponse, err error) {
	err = xml.Unmarshal(bodyBytes, &model)
	if err != nil {
		err = fmt.Errorf("error in unmarshaling: %w", err)
		return
	}

	data, err := ReadOrgcode()
	if err != nil {
		err = fmt.Errorf("error in reading json: %w", err)
	}

	model.AuthenticationSuccess.Attributes.Jusuran = data[model.AuthenticationSuccess.Attributes.Kd_org]
	return
}

func ReadOrgcode() (data map[string]Jurusan, err error) {
	abs, err := filepath.Abs("../static/orgcode.json")
	if err != nil {
		return
	}

	file, err := ioutil.ReadFile(abs)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(file), &data)
	return
}

type ServiceResponse struct {
	XMLName               xml.Name              `xml:"serviceResponse" json:"-"`
	AuthenticationSuccess AuthenticationSuccess `xml:"authenticationSuccess"`
}

type AuthenticationSuccess struct {
	XMLName    xml.Name   `xml:"authenticationSuccess" json:"-"`
	User       string     `xml:"user" json:"user"`
	Attributes Attributes `xml:"attributes" json:"attributes"`
}

type Attributes struct {
	XMLName    xml.Name `xml:"attributes" json:"-"`
	Ldap_cn    string   `xml:"ldap_cn" xml:"ldap_cn"`
	Kd_org     string   `xml:"kd_org" json:"kd_org"`
	Peran_user string   `xml:"peran_user" json:"peran_user"`
	Nama       string   `xml:"nama" json:"nama"`
	Npm        string   `xml:"npm" json:"npm"`
	Jusuran    Jurusan  `json:"jurusan"`
}

type Jurusan struct {
	Faculty      string `json:"faculty"`
	ShortFaculty string `json:"shortFaculty"`
	Major        string `json:"major"`
	Program      string `json:"program"`
}
