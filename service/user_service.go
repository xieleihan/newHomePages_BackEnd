package service

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gin/db"
	"gin/model"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SRP常数定义
var (
	// RFC 5054 3072-bit素数（768个十六进制字符）
	N = new(big.Int)
	g = big.NewInt(2)
	k = big.NewInt(3)
)

func init() {
	// 初始化SRP参数 N (RFC 5054 3072-bit)
	// 必须和前端的 RFC5054_N_HEX 完全一致！
	N.SetString("FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE65381FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE65381FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD129024E088A67CC74020BBEA63B139B22514A08798E3404DDEF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7EDEE386BFB5A899FA5AE9F24117C4B1FE649286651ECE65381FFFFFFFFFFFFFFFF", 16)
}

// RegisterService 用户注册服务
// 验证用户信息并创建新用户账户
// 使用SRP协议，Salt和Verifier由客户端生成
func RegisterService(req model.Register) error {
	// 检查用户名与邮箱是否已存在
	var existingUser model.User
	if err := db.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		// 用户名或邮箱已存在
		fmt.Println("用户名或邮箱已存在 - 用户名:", req.Username, "邮箱:", req.Email)
		return errors.New("用户名或邮箱已存在")
	}

	// 验证邮箱验证码
	if !VerifyCode(req.Email, req.EmailVerificationCode) {
		fmt.Println("邮箱验证码验证失败 - 邮箱:", req.Email)
		return errors.New("邮箱验证码无效或已过期")
	}

	// 验证图片验证码
	if !VerifyCaptcha(req.HumanCheckKey, req.HumanCheckCode) {
		fmt.Println("图片验证码验证失败 - Key:", req.HumanCheckKey)
		return errors.New("图形验证码无效或已过期")
	}

	// 检查 Verifier 长度
	if len(req.Verifier) < 768 {
		fmt.Println("警告: 注册请求中 Verifier 长度不足 768，可能存在问题。当前长度:", len(req.Verifier))
	}

	// 创建新用户 - 只保存必要的字段
	newUser := model.User{
		Username:  req.Username,
		Email:     req.Email,
		Salt:      req.Salt,
		Verifier:  req.Verifier,
		UserId:    uuid.New().String(),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	if err := db.DB.Create(&newUser).Error; err != nil {
		fmt.Println("创建用户失败:", err)
		return errors.New("创建用户失败: " + err.Error())
	}

	fmt.Println("用户注册成功 - 用户名:", newUser.Username, "邮箱:", newUser.Email)
	return nil
}

// LoginService 用户登录第一步服务
// 使用SRP协议，返回服务器公钥B和盐值Salt
func LoginService(req model.Login) (model.LoginResponse, error) {
	var user model.User
	if err := db.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		// 用户不存在
		fmt.Println("用户不存在 - 用户名:", req.Username)
		return model.LoginResponse{}, errors.New("用户不存在")
	}

	fmt.Println("用户查询成功 - 用户名:", user.Username)

	// 检查 Verifier 长度，防止数据库截断导致计算错误
	if len(user.Verifier) < 768 {
		fmt.Printf("严重错误: 数据库中 Verifier 长度不足 (期望 768, 实际 %d)。请检查数据库字段类型是否为 TEXT。\n", len(user.Verifier))
		return model.LoginResponse{}, errors.New("服务器内部错误: 用户数据异常")
	}

	// 使用SRP算法生成服务器公钥B
	// B = (k*v + g^b) mod N
	// 其中：k=3, v=Verifier, g=2, b=随机私钥, N=大素数

	// 1. 将Verifier从hex字符串转换为big.Int
	verifier := new(big.Int)
	verifier.SetString(user.Verifier, 16)

	// 2. 生成随机服务器私钥b (0 < b < N)
	b, err := generateRandomBigInt(N)
	if err != nil {
		fmt.Println("生成随机数失败:", err)
		return model.LoginResponse{}, errors.New("生成随机数失败")
	}

	// 3. 计算g^b mod N
	gb := new(big.Int)
	gb.Exp(g, b, N)

	// 4. 计算k*v mod N
	kv := new(big.Int)
	kv.Mul(k, verifier)
	kv.Mod(kv, N)

	// 5. 计算B = (k*v + g^b) mod N
	B := new(big.Int)
	B.Add(kv, gb)
	B.Mod(B, N)

	// 转换为hex字符串
	BHex := fmt.Sprintf("%x", B)

	// 将会话数据存储到Redis中(供第二步使用)
	sessionKey := "srp:session:" + req.Username + ":" + time.Now().Format("20060102150405")
	sessionData := map[string]interface{}{
		"b":         b.String(),
		"B":         BHex,
		"A":         req.A,
		"v":         user.Verifier,
		"N":         N.String(),
		"salt":      user.Salt, // 存储 Salt
		"username":  req.Username,
		"timestamp": time.Now().Unix(),
	}

	// 序列化为JSON并存储到Redis
	sessionJSON, err := json.Marshal(sessionData)
	if err != nil {
		fmt.Println("会话数据序列化失败:", err)
		return model.LoginResponse{}, errors.New("会话数据序列化失败")
	}

	// 存储到Redis，设置1小时过期时间
	if err := db.RDB.Set(db.Ctx, sessionKey, string(sessionJSON), time.Hour).Err(); err != nil {
		fmt.Println("存储会话数据到Redis失败:", err)
		return model.LoginResponse{}, errors.New("存储会话数据失败")
	}

	fmt.Println(" 会话数据已存储到Redis - Key:", sessionKey)

	response := model.LoginResponse{
		Salt: user.Salt,
		B:    BHex,
	}

	fmt.Println("登录第一步成功 - 返回Salt和公钥B")
	return response, nil
}

// generateRandomBigInt 生成范围内的随机大整数 (0, max)
func generateRandomBigInt(max *big.Int) (*big.Int, error) {
	// 使用crypto/rand生成随机数
	// rand.Int返回范围[0, max)的随机大整数
	result, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}
	// 确保不是0
	if result.Sign() == 0 {
		result.Add(result, big.NewInt(1))
	}
	return result, nil
}

// LoginStep2Service 用户登录第二步服务
// 使用SRP协议，验证客户端证据消息M1，返回服务器证据消息M2
func LoginStep2Service(req model.LoginStep2) (model.LoginStep2Response, error) {
	var user model.User
	if err := db.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		fmt.Println("用户不存在 - 用户名:", req.Username)
		return model.LoginStep2Response{}, errors.New("用户不存在")
	}

	// 从Redis中获取会话数据
	sessionPattern := "srp:session:" + req.Username + ":*"

	// 获取所有匹配的session key
	keys, err := db.RDB.Keys(db.Ctx, sessionPattern).Result()
	if err != nil || len(keys) == 0 {
		fmt.Println(" 未找到会话数据 - 用户名:", req.Username)
		return model.LoginStep2Response{}, errors.New("会话已过期或不存在")
	}

	sessionKey := keys[0]
	sessionJSON, err := db.RDB.Get(db.Ctx, sessionKey).Result()
	if err != nil {
		fmt.Println("获取会话数据失败:", err)
		return model.LoginStep2Response{}, errors.New("获取会话数据失败")
	}

	// 解析会话数据
	var sessionData map[string]interface{}
	if err := json.Unmarshal([]byte(sessionJSON), &sessionData); err != nil {
		fmt.Println("会话数据反序列化失败:", err)
		return model.LoginStep2Response{}, errors.New("会话数据格式错误")
	}

	fmt.Println("从Redis获取会话数据成功")

	// SRP第二步验证算法
	// 1. 从会话数据中获取 b、A、B、v、N
	bStr := sessionData["b"].(string)
	b := new(big.Int)
	b.SetString(bStr, 10)

	AHex := sessionData["A"].(string)
	A := new(big.Int)
	A.SetString(AHex, 16)

	BHex := sessionData["B"].(string)
	B := new(big.Int)
	B.SetString(BHex, 16)

	vHex := sessionData["v"].(string)
	v := new(big.Int)
	v.SetString(vHex, 16)

	NStr := sessionData["N"].(string)
	N := new(big.Int)
	N.SetString(NStr, 10)

	username := sessionData["username"].(string)

	// 2. 计算 u = H(PAD(A) | PAD(B)) - 使用SHA256
	u := calculateU(A, B)
	fmt.Println("计算u成功:", u)

	// 3. 计算 S = (A * v^u)^b mod N (服务器端)
	// 先计算 v^u mod N
	vu := new(big.Int)
	vu.Exp(v, u, N)

	// 计算 A * v^u mod N
	Avu := new(big.Int)
	Avu.Mul(A, vu)
	Avu.Mod(Avu, N)

	// 计算 (A * v^u)^b mod N
	S := new(big.Int)
	S.Exp(Avu, b, N)

	// fmt.Println("计算共享密钥S成功:", S) // S是敏感信息，生产环境不应打印
	// 打印S的前几位用于调试
	sBytes := S.Bytes()
	if len(sBytes) > 10 {
		fmt.Printf(" S (前10字节): %x\n", sBytes[:10])
	}

	// 4. 计算 K = H(S)
	K := hashBigInt(S)
	fmt.Printf(" K (Hex): %x\n", K)

	// 5. 验证 M1 = H(H(N) XOR H(g) | H(I) | s | A | B | K)
	// 计算预期的M1
	HN := hashBigInt(N)
	Hg := hashBigInt(g)

	// H(N) XOR H(g)
	HNxorHg := make([]byte, len(HN))
	for i := range HN {
		HNxorHg[i] = HN[i] ^ Hg[i]
	}

	// H(I) = H(username)
	HI := hashString(username)

	// Salt作为s - 需要从hex字符串解码为字节
	salt := sessionData["salt"].(string)
	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		fmt.Println("Salt hex解码失败:", err)
		saltBytes = []byte(salt) // 降级到字符串字节
	}

	// 构建M1的输入
	// 注意：这里需要确保顺序和客户端完全一致
	// M1 = H(H(N) XOR H(g) | H(I) | s | A | B | K)
	var M1Input []byte
	M1Input = append(M1Input, HNxorHg...)
	M1Input = append(M1Input, HI...)
	M1Input = append(M1Input, saltBytes...)
	M1Input = append(M1Input, A.Bytes()...)
	M1Input = append(M1Input, B.Bytes()...)
	M1Input = append(M1Input, K...)

	fmt.Printf(" M1输入组成 - HNxorHg长度: %d, HI长度: %d, salt长度: %d, A长度: %d, B长度: %d, K长度: %d\n",
		len(HNxorHg), len(HI), len(saltBytes), len(A.Bytes()), len(B.Bytes()), len(K))
	fmt.Printf(" HNxorHg: %x\n", HNxorHg)
	fmt.Printf(" HI: %x\n", HI)
	fmt.Printf(" salt bytes: %x\n", saltBytes)
	fmt.Printf(" A.Bytes(): %x\n", A.Bytes())
	fmt.Printf(" B.Bytes(): %x\n", B.Bytes())
	fmt.Printf(" K: %x\n", K)

	expectedM1 := hashBytes(M1Input)
	expectedM1Hex := fmt.Sprintf("%x", expectedM1)

	// 验证客户端发送的M1是否匹配
	if !strings.EqualFold(req.M1, expectedM1Hex) {
		fmt.Println(" M1验证失败 - 期望:", expectedM1Hex, "实际:", req.M1)
		db.RDB.Del(db.Ctx, sessionKey)
		return model.LoginStep2Response{}, errors.New("登录验证失败")
	}

	fmt.Println(" M1验证成功")

	// 6. 计算 M2 = H(A | M1 | K)
	M2Input := append(A.Bytes(), expectedM1...)
	M2Input = append(M2Input, K...)
	M2 := hashBytes(M2Input)
	M2Hex := fmt.Sprintf("%x", M2)

	fmt.Println(" 计算M2成功")

	// 7. 返回M2
	response := model.LoginStep2Response{
		M2: M2Hex,
	}

	// 删除已使用的会话数据
	db.RDB.Del(db.Ctx, sessionKey)

	fmt.Println(" 登录第二步成功 - 返回服务器证据M2")
	return response, nil
}

// calculateU 计算 u = H(PAD(A) | PAD(B))
// PAD是指将字节数组填充到相同长度
func calculateU(A, B *big.Int) *big.Int {
	// 获取N的字节长度作为PAD长度
	nlen := N.BitLen() / 8
	if N.BitLen()%8 > 0 {
		nlen++
	}

	fmt.Printf(" 计算 u - N字节长度 (PAD长度): %d\n", nlen)

	// PAD A 和 B
	aBytes := make([]byte, nlen)
	bBytes := make([]byte, nlen)

	aData := A.Bytes()
	bData := B.Bytes()

	copy(aBytes[nlen-len(aData):], aData)
	copy(bBytes[nlen-len(bData):], bData)

	// 连接 A | B
	input := append(aBytes, bBytes...)

	fmt.Printf(" PAD(A) 长度: %d, PAD(B) 长度: %d, 总长度: %d\n", nlen, nlen, len(input))

	// 计算 H(PAD(A) | PAD(B))
	hash := hashBytes(input)
	u := new(big.Int)
	u.SetBytes(hash)

	fmt.Printf(" u值: %d\n", u)

	return u
}

// hashBigInt 计算 H(bigint) - 将big.Int转换为字节后进行哈希
func hashBigInt(num *big.Int) []byte {
	return hashBytes(num.Bytes())
}

// hashString 计算 H(string)
func hashString(str string) []byte {
	return hashBytes([]byte(str))
}

// hashBytes 计算 H(bytes) - 使用SHA512
func hashBytes(data []byte) []byte {
	hash := sha512.Sum512(data)
	return hash[:]
}
