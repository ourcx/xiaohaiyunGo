package utils

//这里是使用bcrypt对密码进行加密和解密的工具函数
import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// 实现 bcrypt 或其他哈希算法
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
	//密码加密
}

func CheckPassword(password, hash string) bool {
	//fmt.Println(password + "hash1")
	//fmt.Println(hash + "hash2")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	//fmt.Println(err)
	if err != nil {
		return false
	}
	return true
	//密码对比
	//这里传进去的是一个hash和一个非hash的密码

}
