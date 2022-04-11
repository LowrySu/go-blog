package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-blog/services/conf"
	"go-blog/services/store"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/cristalhq/jwt/v3"
)

var (
	jwtSigner   jwt.Signer
	jwtVerifier jwt.Verifier
)

func jwtSetup(conf conf.Config) {
	// 建立签名者和验证者

	var err error
	key := []byte(conf.JwtSecret)

	jwtSigner, err = jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		log.Panic().Err(err).Msg("Error creating JWT signer")
	}

	jwtVerifier, err = jwt.NewVerifierHS(jwt.HS256, key)
	if err != nil {
		log.Panic().Err(err).Msg("Error creating JWT verifier")
	}
}

func generateJWT(user *store.User) string {
	// 输出用户的id和有效时间
	claims := &jwt.RegisteredClaims{
		ID:        fmt.Sprint(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
	}
	// 建立生成器
	builder := jwt.NewBuilder(jwtSigner)
	// 生成token
	token, err := builder.Build(claims)
	if err != nil {
		log.Panic().Err(err).Msg("Error building JWT")
	}
	return token.String()

}

func verifyJWT(tokenStr string) (int, error) {
	// 从token内提取id

	// 解析token
	token, err := jwt.Parse([]byte(tokenStr))
	if err != nil {
		log.Panic().Err(err).Str("tokenStr", tokenStr).Msg("Error parsing JWT")
		return 0, err
	}

	// 验证token
	if err := jwtVerifier.Verify(token.Payload(), token.Signature()); err != nil {
		log.Error().Err(err).Msg("Error verifying token")
		return 0, err
	}

	// 解析
	var claims jwt.StandardClaims
	if err := json.Unmarshal(token.RawClaims(), &claims); err != nil {
		log.Error().Err(err).Msg("Error unmarshalling JWT claims ")
		return 0, err
	}

	// 查看是否过期
	if notExpired := claims.IsValidAt(time.Now()); !notExpired {
		return 0, errors.New("Token expired.")
	}

	// 查看id是否有效
	id, err := strconv.Atoi(claims.ID)
	if err != nil {
		log.Error().Err(err).Str("claims.ID", claims.ID).Msg("Error converting claims ID to number")
		return 0, errors.New("ID in token is not valid")
	}

	return id, err

}
