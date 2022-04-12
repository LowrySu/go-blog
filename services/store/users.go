package store

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// 实体类: User
type User struct {
	ID             int
	Username       string `binding:"required,min=5,max=30"`
	Password       string `pg:"-" binding:"required,min=7,max=32"`
	HashedPassword []byte `json:"-"`
	Salt           []byte `json:"-"`
	CreatedAt      time.Time
	ModifiedAt     time.Time
	Posts          []*Post `json:"-" pg:"fk:user_id,rel:has-many,on_delete:CASCADE"`
}

var _ pg.AfterSelectHook = (*User)(nil)

func (user *User) AfterSelect(ctx context.Context) error {
	if user.Posts == nil {
		user.Posts = []*Post{}
	}
	return nil
}

func AddUser(user *User) error {
	// 获取随机数种子
	salt, err := GenerateSalt()
	if err != nil {
		return err
	}

	// 对密码加密
	toHash := append([]byte(user.Password), salt...)                               // 拼接密码和随机数种子
	hashedPassword, err := bcrypt.GenerateFromPassword(toHash, bcrypt.DefaultCost) // 计算哈希值
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
		return err
	}

	// 设置加密的密码
	user.Salt = salt
	user.HashedPassword = hashedPassword

	// 插入数据
	_, err = db.Model(user).Returning("*").Insert()
	if err != nil {
		log.Error().Err(err).Msg("Error inserting new user")
		return dbError(err)
	}
	return nil
}

func Authenticate(username, password string) (*User, error) {
	user := new(User)

	// 根据用户名搜索数据,并赋值到user的属性内
	if err := db.Model(user).Where(
		"username = ?", username).Select(); err != nil {
		return nil, err
	}

	// 获取机密种子,并拼接密码
	salt := user.Salt
	salted := append([]byte(password), salt...)
	// 验证密码是否正确
	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, salted); err != nil {
		return nil, err
	}
	return user, nil
}

func GenerateSalt() ([]byte, error) {
	// 生成一个16个随机数的列表
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func FetchUser(id int) (*User, error) {
	// 用id获取用户

	user := new(User)
	user.ID = id
	err := db.Model(user).Returning("*").WherePK().Select()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching user")
		return nil, err
	}
	return user, nil
}
