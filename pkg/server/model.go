package server

import (
	"github.com/chutommy/eetgateway/pkg/eet"
)

// HTTPPingResponse represents a response structure of HTTP responses for pings.
type HTTPPingResponse struct {
	EETGatewayStatus string `json:"eet_gateway"`
	TaxAdminStatus   string `json:"tax_admin"`
}

func encodePingResponse(taxAdmin string) *HTTPPingResponse {
	return &HTTPPingResponse{
		EETGatewayStatus: "online", // is able to response
		TaxAdminStatus:   taxAdmin,
	}
}

// HTTPEETRequest represents a binding structure to HTTP requests for sending sales.
type HTTPEETRequest struct {
	CertID       string `json:"cert_id,omitempty" binding:"required"`
	CertPassword string `json:"cert_password,omitempty" binding:""`

	UUIDZpravy      eet.UUIDType   `json:"uuid_zpravy" binding:"omitempty,uuid_zpravy"`
	DatOdesl        eet.DateTime   `json:"dat_odesl" binding:""`
	PrvniZaslani    bool           `json:"prvni_zaslani" binding:""`
	Overeni         bool           `json:"overeni" binding:""`
	DICPopl         eet.CZDICType  `json:"dic_popl" binding:"required,dic"`
	DICPoverujiciho eet.CZDICType  `json:"dic_poverujiciho" binding:"omitempty,dic,necsfield=Dicpopl"`
	IDProvoz        int            `json:"id_provoz" binding:"required,id_provoz"`
	IDPokl          eet.String20   `json:"id_pokl" binding:"required,id_pokl"`
	PoradCis        eet.String25   `json:"porad_cis" binding:"required,porad_cis"`
	DatTrzby        eet.DateTime   `json:"dat_trzby" binding:""`
	CelkTrzba       eet.CastkaType `json:"celk_trzba" binding:"required,fin_poloz"`
	ZaklNepodlDPH   eet.CastkaType `json:"zakl_nepodl_dph" binding:"omitempty,fin_poloz"`
	ZaklDan1        eet.CastkaType `json:"zakl_dan1" binding:"omitempty,fin_poloz"`
	Dan1            eet.CastkaType `json:"dan1" binding:"omitempty,fin_poloz"`
	ZaklDan2        eet.CastkaType `json:"zakl_dan2" binding:"omitempty,fin_poloz"`
	Dan2            eet.CastkaType `json:"dan2" binding:"omitempty,fin_poloz"`
	ZaklDan3        eet.CastkaType `json:"zakl_dan3" binding:"omitempty,fin_poloz"`
	Dan3            eet.CastkaType `json:"dan3" binding:"omitempty,fin_poloz"`
	CestSluz        eet.CastkaType `json:"cest_sluz" binding:"omitempty,fin_poloz"`
	PouzitZboz1     eet.CastkaType `json:"pouzit_zboz1" binding:"omitempty,fin_poloz"`
	PouzitZboz2     eet.CastkaType `json:"pouzit_zboz2" binding:"omitempty,fin_poloz"`
	PouzitZboz3     eet.CastkaType `json:"pouzit_zboz3" binding:"omitempty,fin_poloz"`
	UrcenoCerpzZuct eet.CastkaType `json:"urceno_cerp_zuct" binding:"omitempty,fin_poloz"`
	CerpZuct        eet.CastkaType `json:"cerp_zuct" binding:"omitempty,fin_poloz"`
	Rezim           eet.RezimType  `json:"rezim" binding:"omitempty,rezim"`
}

func decodeEETRequest(req *HTTPEETRequest) *eet.TrzbaType {
	return &eet.TrzbaType{
		Hlavicka: eet.TrzbaHlavickaType{
			Uuidzpravy:   req.UUIDZpravy,
			Datodesl:     req.DatOdesl,
			Prvnizaslani: req.PrvniZaslani,
			Overeni:      req.Overeni,
		},
		Data: eet.TrzbaDataType{
			Dicpopl:         req.DICPopl,
			Dicpoverujiciho: req.DICPoverujiciho,
			Idprovoz:        req.IDProvoz,
			Idpokl:          req.IDPokl,
			Poradcis:        req.PoradCis,
			Dattrzby:        req.DatTrzby,
			Celktrzba:       req.CelkTrzba,
			Zaklnepodldph:   req.ZaklNepodlDPH,
			Zakldan1:        req.ZaklDan1,
			Dan1:            req.Dan1,
			Zakldan2:        req.ZaklDan2,
			Dan2:            req.Dan2,
			Zakldan3:        req.ZaklDan3,
			Dan3:            req.Dan3,
			Cestsluz:        req.CestSluz,
			Pouzitzboz1:     req.PouzitZboz1,
			Pouzitzboz2:     req.PouzitZboz2,
			Pouzitzboz3:     req.PouzitZboz3,
			Urcenocerpzuct:  req.UrcenoCerpzZuct,
			Cerpzuct:        req.CerpZuct,
			Rezim:           req.Rezim,
		},
	}
}

// HTTPEETResponse represents a reponse structure to HTTP sale requests.
type HTTPEETResponse struct {
	GatewayError string                    `json:"gateway_error,omitempty"`
	Dat          *eet.DateTime             `json:"dat,omitempty"`
	Fik          eet.FikType               `json:"fik,omitempty"`
	Zprava       string                    `json:"zprava,omitempty"`
	Kod          int                       `json:"kod,omitempty"`
	Test         bool                      `json:"test,omitempty"`
	Varovani     []eet.OdpovedVarovaniType `json:"varovani,omitempty"`
}

func encodeEETResponse(err error, odpoved *eet.OdpovedType) *HTTPEETResponse {
	if err != nil {
		return &HTTPEETResponse{
			GatewayError: err.Error(),
		}
	} else if odpoved != nil {
		// select the non-empty date/time
		cas := odpoved.Hlavicka.Datprij
		if (cas == eet.DateTime{}) {
			cas = odpoved.Hlavicka.Datodmit
		}

		return &HTTPEETResponse{
			Dat:      &cas,
			Fik:      odpoved.Potvrzeni.Fik,
			Zprava:   odpoved.Chyba.Zprava,
			Kod:      odpoved.Chyba.Kod,
			Test:     odpoved.Potvrzeni.Test || odpoved.Chyba.Test,
			Varovani: odpoved.Varovani,
		}
	}

	return &HTTPEETResponse{}
}

// HTTPEETRequest represents a binding structure to HTTP requests for storing certificate.
type HTTPCreateCertRequest struct {
	ID             string `json:"id" binding:""`
	Password       string `json:"password" binding:""`
	PKCS12Data     []byte `json:"pkcs12_data" binding:"required"`
	PKCS12Password string `json:"pkcs12_password" binding:"required"`
}

// HTTPEETResponse represents a reponse structure to HTTP request for storing certificate .
type HTTPCreateCertResponse struct {
	GatewayError string `json:"gateway_error,omitempty"`
	ID           string `json:"id,omitempty"`
}

func encodeCreateCertResponse(err error, id *string) *HTTPCreateCertResponse {
	if err != nil {
		return &HTTPCreateCertResponse{
			GatewayError: err.Error(),
		}
	} else if id != nil {
		return &HTTPCreateCertResponse{
			ID: *id,
		}
	}

	return &HTTPCreateCertResponse{}
}
