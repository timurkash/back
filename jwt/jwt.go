package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/timurkash/back/header"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
)

const Authorization = "Authorization"
const SecondsInDay = int64(86400)

type (
	Policy struct {
		Issuer   string
		Audience string
		JwksUrl  string
	}
	Key struct {
		Alg string `json:"alg"`
		E   string `json:"e"`
		Kid string `json:"kid"`
		N   string `json:"n"`
		Kty string `json:"kty"`
		Use string `json:"use"`
	}
	Keys struct {
		Keys []*Key `json:"keys"`
	}
)

var (
	pubKeys            []*rsa.PublicKey
	policy             *Policy
	notAuthorizedError = errors.New("NotAuthorized")
)

//func (p *Policy)

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

func (p *Policy) GetPolicy() *Policy {
	return p
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
	for _, key := range keys.Keys {
		rsaPublicKey, err := getRsaPublicKey(key)
		if err != nil {
			return err
		}
		pubKeys = append(pubKeys, rsaPublicKey)
	}
	return nil
}

func getRsaPublicKey(jwk *Key) (*rsa.PublicKey, error) {
	if jwk.Kty != "RSA" {
		return nil, fmt.Errorf("invalid key type: %s", jwk.Kty)
	}
	nb, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}
	e := 0
	// The default exponent is usually 65537, so just compare the
	// base64 for [1,0,1] or [0,1,0,1]
	if jwk.E == "AQAB" || jwk.E == "AAEAAQ" {
		e = 65537
	} else {
		// need to decode "e" as a big-endian int
		return nil, fmt.Errorf("need to decode e: %s", jwk.E)
	}
	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nb),
		E: e,
	}, nil
}

func GetJwtFromHeader(r *http.Request) (*string, *string, *string, error) {
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

func VerifyToken(token *string, part2 *string) error {
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
