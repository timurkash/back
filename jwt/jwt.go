package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"gitlab.com/mcsolutions/lib/back/common/header"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"time"
)

const Authorization = "Authorization"
const SecondsInDay = int64(86400)

type (
	Policy struct {
		Issuer   string
		Audience string
		JwksUrl  string
	}

	Keys struct {
		Keys []map[string]string `json:"keys"`
	}

	Payload struct {
		Iss           string `json:"iss"`
		Aud           string `json:"aud"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Roles         string `json:"roles"`
		Iat           int64  `json:"iat"`
		Exp           int64  `json:"exp"`
		UserId        string `json:"user_id"`
		IsAdmin       bool   `json:"isAdmin"`
		IsPsy         bool   `json:"isPsy"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Lang          string `json:"lang"`
	}

	EmailRoles struct {
		Iat     int64
		Email   string
		Uid     string
		Roles   []string
		IsAdmin bool
		IsPsy   bool
		Name    string
		Picture string
		Lang    string
	}
)

var (
	pubKeys            []*rsa.PublicKey
	policy             *Policy
	notAuthorizedError = errors.New("NotAuthorized")
)

func SetPolicy(p *Policy) error {
	if p == nil {
		return errors.New("policy is nil")
	}
	if p.Audience == "" {
		return errors.New("policy.Audience is empty")
	}
	policy = p
	return nil
}

func getPubKeys() error {
	if policy == nil {
		return errors.New("policy not defined")
	}
	if policy.JwksUrl == "" {
		return errors.New("certifications url not defined")
	}
	resp, err := http.Get(policy.JwksUrl)
	if err != nil {
		return err
	}
	bytes_, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	keys := &Keys{}
	if err := json.Unmarshal(bytes_, keys); err != nil {
		return err
	}
	pubKeys = nil
	for _, jwk := range keys.Keys {
		rsaPublicKey, err := getRsaPublicKey(jwk)
		if err != nil {
			return err
		}
		pubKeys = append(pubKeys, rsaPublicKey)
	}
	return nil
}

func getRsaPublicKey(jwk map[string]string) (*rsa.PublicKey, error) {
	if jwk["kty"] != "RSA" {
		return nil, fmt.Errorf("invalid key type:", jwk["kty"])
	}
	nb, err := base64.RawURLEncoding.DecodeString(jwk["n"])
	if err != nil {
		return nil, err
	}
	e := 0
	// The default exponent is usually 65537, so just compare the
	// base64 for [1,0,1] or [0,1,0,1]
	if jwk["e"] == "AQAB" || jwk["e"] == "AAEAAQ" {
		e = 65537
	} else {
		// need to decode "e" as a big-endian int
		return nil, fmt.Errorf("need to deocde e:", jwk["e"])
	}
	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nb),
		E: e,
	}, nil
}

func UnauthorizedCode(err error) int {
	if err == notAuthorizedError {
		return http.StatusUnauthorized
	}
	return http.StatusForbidden
}

func HasAuthorization(r *http.Request) bool {
	return r.Header.Get(Authorization) != ""
}

func GetJwtFromHeader(r *http.Request) (*string, *string, *string, error) {
	if !HasAuthorization(r) {
		return nil, nil, nil, notAuthorizedError
	}
	bearer, err := header.GetHeaderRequired(r, Authorization)
	if err != nil {
		return nil, nil, nil, err
	}
	if !strings.HasPrefix(bearer, "Bearer ey") {
		return nil, nil, nil, errors.New("bad authorization header")
	}
	bearer = bearer[7:]
	if strings.Index(bearer, " ") >= 0 {
		return nil, nil, nil, errors.New("jwt has space")
	}
	parts := strings.Split(bearer, ".")
	if len(parts) != 3 {
		return nil, nil, nil, errors.New("jwt does not consist of 3 part separated by a dot")
	}

	token := strings.Join(parts[0:2], ".")
	part1 := parts[1]
	part2 := parts[2]
	return &token, &part1, &part2, nil
}

func IsJwtValid(r *http.Request, seconds ...int64) (*EmailRoles, int, error) {
	token, part1, part2, err := GetJwtFromHeader(r)
	if err != nil {
		return nil, http.StatusUnauthorized, err
	}
	if err := verifyToken(token, part2); err != nil {
		return nil, http.StatusConflict, err
	}
	bytes, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(*part1)
	if err != nil {
		return nil, http.StatusConflict, err
	}
	if policy == nil {
		return nil, http.StatusConflict, errors.New("policy not defined")
	}
	now := time.Now().Unix()
	payload := new(Payload)
	if err := json.Unmarshal(bytes, payload); err != nil {
		return nil, http.StatusConflict, err
	}
	if payload.Iss != policy.Issuer {
		return nil, http.StatusConflict, errors.New("issuers not matched")
	}
	if payload.Aud != payload.Aud {
		return nil, http.StatusConflict, errors.New("audience not matched")
	}
	validPeriod := SecondsInDay
	if len(seconds) > 0 {
		validPeriod = seconds[0]
	}
	if now > payload.Iat+validPeriod {
		return nil, 419, errors.New("token has expired")
	}
	if !payload.EmailVerified {
		return nil, http.StatusForbidden, fmt.Errorf("email of %s not verified", payload.Email)
	}
	return &EmailRoles{
		Email:   payload.Email,
		Uid:     payload.UserId,
		Roles:   strings.Split(payload.Roles, ","),
		IsAdmin: payload.IsAdmin,
		IsPsy:   payload.IsPsy,
		Iat:     payload.Iat,
		Name:    payload.Name,
		//Picture: payload.Picture,
		Lang: payload.Lang,
	}, http.StatusOK, nil
}

func verifyToken(token *string, part2 *string) error {
	if pubKeys == nil {
		if err := getPubKeys(); err != nil {
			return err
		}
	}
	for _, pubKey := range pubKeys {
		if err := jwtgo.SigningMethodRS256.Verify(*token, *part2, pubKey); err == nil {
			return nil
		}
	}
	//expired pubKeys reloading
	if err := getPubKeys(); err != nil {
		return err
	}
	for _, pubKey := range pubKeys {
		if err := jwtgo.SigningMethodRS256.Verify(*token, *part2, pubKey); err == nil {
			return nil
		}
	}
	return errors.New("bad token")
}
