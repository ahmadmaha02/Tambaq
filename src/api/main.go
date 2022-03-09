package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID             uint   `gorm:"primarykey" json:"id"`
	Name           string `json:"nama"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Alamat         string `json:"alamat"`
	Jenis_Budidaya string `json:"jenis_budidaya"`
	Lokasi_Tambak  string `json:"lokasi_tambak"`
	Luas_Kolam     string `json:"luas_kolam"`
}

type Ikan struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Kategori    string `json:"kategori"`
	Jenis_Ikan  string `json:"jenis_ikan"`
	Harga       string `json:"harga"`
	TokoID      uint   `json:"toko_id"`
	Provinsi    string `json:"provinsi"`
	Kota        string `json:"kota"`
	Bulan_Panen string `json:"bulan_panen"`
}

type Toko struct {
	ID uint `gorm:"primarykey" json:"id"`
}

type Tweet struct {
	ID        uint `gorm:"primarykey" json:"id"`
	UserID    uint `json:"user_id"`
	User      User
	Content   string        `json:"name"`
	RepliedTo sql.NullInt64 `json:"replied_to"`
	CreatedAt time.Time     `json:"created_at"`
}

var db *gorm.DB
var r *gin.Engine

func InitDB() error {
	_db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/intern_workshop?parseTime=true"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = _db
	err = db.AutoMigrate(&User{}, &Tweet{}, &Ikan{})

	tweet := Tweet{
		ID: 1,
	}
	db.Preload("User").Take(&tweet)
	fmt.Println(tweet)
	if err != nil {
		return err
	}
	return nil
}

func InitGin() {
	r = gin.Default()
	r.Use(cors.Default())
}

type postRegisterBody struct {
	Name           string `json:"nama"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Alamat         string `json:"alamat"`
	Jenis_Budidaya string `json:"jenis_budidaya"`
	Lokasi_Tambak  string `json:"lokasi_tambak"`
	Luas_Kolam     string `json:"luas_kolam"`
}

type postLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type postTambahIkanBody struct {
	Kategori    string `json:"kategori"`
	Jenis_Ikan  string `json:"jenis_ikan"`
	Harga       string `json:"harga"`
	TokoID      uint   `json:"toko_id"`
	Provinsi    string `json:"provinsi"`
	Kota        string `json:"kota"`
	Bulan_Panen string `json:"bulan_panen"`
}

type patchUserBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type postCreateTweetBody struct {
	Content   string `json:"content"`
	RepliedTo uint   `json:"replied_to,omitempty"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		header = header[len("Bearer "):]
		token, err := jwt.Parse(header, func(t *jwt.Token) (interface{}, error) {
			return []byte("passwordBuatSigning"), nil
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "JWT validation error.",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("id", claims["id"])
			c.Next()
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "JWT invalid.",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
	}
}

func InitRouter() {
	r.POST("/api/auth/register", func(c *gin.Context) {

		var body postRegisterBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Body is invalid.",
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		user := User{
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
		}
		if result := db.Create(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Berhasil Membuat Akun",
			"status":  "Sukses",
			"data": gin.H{
				"id": user.ID,
			},
		})
	})

	r.POST("/api/auth/register-member", func(c *gin.Context) {
		_, isEmailExists := c.GetQuery("email")
		if !isEmailExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Email Sudah Digunakan.",
				"status":  "Register Gagal.",
			})
			return
		}
		var body postRegisterBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		user := User{
			Name:           body.Name,
			Email:          body.Email,
			Alamat:         body.Alamat,
			Jenis_Budidaya: body.Jenis_Budidaya,
			Lokasi_Tambak:  body.Lokasi_Tambak,
			Luas_Kolam:     body.Luas_Kolam,
		}
		if result := db.Create(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User successfully registered.",
			"data": gin.H{
				"id": user.ID,
			},
		})
	})

	r.POST("/api/auth/login", func(c *gin.Context) {
		var body postLoginBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		user := User{}
		if result := db.Where("email = ?", body.Email).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if user.Password == body.Password {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id":  user.ID,
				"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
			})
			tokenString, err := token.SignedString([]byte("passwordBuatSigning"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when generating the token.",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Login Berhasil",
				"status":  "Sukses",
				"data": gin.H{
					"id":    user.ID,
					"name":  user.Name,
					"token": tokenString,
				},
			})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Password is incorrect.",
			})
			return
		}
	})

	r.GET("/api/auth/ikan-segar", func(c *gin.Context) {
		ikan := Ikan{}
		if result := db.Where("kategori = ?", "ikan segar").Find(&ikan); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful.",
			"data":    ikan,
		})
	})

	r.GET("/api/auth/ikan-frozen", func(c *gin.Context) {
		ikan := Ikan{}
		if result := db.Where("kategori = ?", "ikan frozen").Find(&ikan); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful.",
			"data":    ikan,
		})
	})

	r.GET("/api/auth/bibit-ikan", func(c *gin.Context) {
		ikan := Ikan{}
		if result := db.Where("kategori = ?", "bibit ikan").Find(&ikan); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful.",
			"data":    ikan,
		})
	})

	r.POST("/api/auth/tambah-ikan", func(c *gin.Context) {
		var body postTambahIkanBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		ikan := Ikan{
			Kategori:    body.Kategori,
			Jenis_Ikan:  body.Jenis_Ikan,
			Harga:       body.Harga,
			TokoID:      body.TokoID,
			Provinsi:    body.Provinsi,
			Kota:        body.Kota,
			Bulan_Panen: body.Bulan_Panen,
		}
		if result := db.Create(&ikan); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Ikan Berhasil Ditambahkan",
			"data":    ikan,
		})
	})

	r.GET("/user", AuthMiddleware(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
		if result := db.Where("id = ?", id).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful",
			"data":    user,
		})
	})

	r.GET("/user/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		user := User{}
		if result := db.Where("id = ?", id).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful.",
			"data":    user,
		})
	})

	r.PATCH("/user/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		var body patchUserBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
				"error":   err.Error(),
			})
			return
		}
		user := User{
			ID:       uint(parsedId),
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
		}
		result := db.Model(&user).Updates(user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result = db.Where("id = ?", parsedId).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update successful.",
			"data":    user,
		})
	})

	r.GET("/user/search", func(c *gin.Context) {
		name, isNameExists := c.GetQuery("name")
		email, isEmailExists := c.GetQuery("email")
		username, isUsernameExists := c.GetQuery("username")
		if !isNameExists && !isEmailExists && !isUsernameExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
			})
			return
		}

		var queryResults []User
		trx := db
		if isNameExists {
			trx = trx.Where("name LIKE ?", "%"+name+"%")
		}
		if isEmailExists {
			trx = trx.Where("email LIKE ?", "%"+email+"%")
		}
		if isUsernameExists {
			trx = trx.Where("username LIKE ?", "%"+username+"%")
		}

		if result := trx.Find(&queryResults); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Search successful",
			"data": gin.H{
				"query": gin.H{
					"name":     name,
					"email":    email,
					"username": username,
				},
				"result": queryResults,
			},
		})
	})

	r.DELETE("/user/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
			})
			return
		}
		user := User{
			ID: uint(parsedId),
		}
		if result := db.Delete(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when deleting from the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete successful.",
		})
	})

	r.POST("/tweet", AuthMiddleware(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var body postCreateTweetBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		repliedTo := sql.NullInt64{
			Int64: 0,
			Valid: false,
		}
		if body.RepliedTo != 0 {
			repliedTo.Int64 = int64(body.RepliedTo)
			repliedTo.Valid = true
		}
		tweet := Tweet{
			Content:   body.Content,
			UserID:    uint(id.(float64)),
			CreatedAt: time.Now(),
			RepliedTo: repliedTo,
		}
		if result := db.Create(&tweet); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Tweet successfully created.",
			"data": gin.H{
				"id": tweet.ID,
			},
		})
	})

	r.GET("/tweet/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		withReplies := false
		withRepliesStr, isWithRepliesSet := c.GetQuery("with_replies")
		if isWithRepliesSet {
			_withReplies, err := strconv.ParseBool(withRepliesStr)
			if err != nil {
				withReplies = false
			}
			withReplies = _withReplies
		}

		tweet := Tweet{}
		if result := db.Where("id = ?", id).Take(&tweet); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		var repliedTo *int64
		var tweetReplies []Tweet
		var tweetRepliesCleaned []gin.H
		if tweet.RepliedTo.Valid {
			repliedTo = &tweet.RepliedTo.Int64
		}
		if withReplies {
			if result := db.Where("replied_to = ?", id).Find(&tweetReplies); result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when querying the database.",
					"error":   result.Error.Error(),
				})
				return
			}
			for _, reply := range tweetReplies {
				tweetRepliesCleaned = append(tweetRepliesCleaned, gin.H{
					"id":         reply.ID,
					"user_id":    reply.UserID,
					"content":    reply.Content,
					"created_at": reply.CreatedAt,
				})
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful.",
			"data": gin.H{
				"id":         tweet.ID,
				"user_id":    tweet.UserID,
				"content":    tweet.Content,
				"replied_to": repliedTo,
				"created_at": tweet.CreatedAt,
				"replies":    tweetRepliesCleaned,
			},
		})
	})

	r.DELETE("/tweet/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
			})
			return
		}
		tweet := Tweet{
			ID: uint(parsedId),
		}
		if result := db.Where("id = ?", id).Take(&tweet); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		var repliedTo *int64
		if tweet.RepliedTo.Valid {
			repliedTo = &tweet.RepliedTo.Int64
		}
		result := db.Delete(&tweet)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when deleting from the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "ID not found.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete successful.",
			"data": gin.H{
				"id":         tweet.ID,
				"user_id":    tweet.UserID,
				"content":    tweet.Content,
				"replied_to": repliedTo,
			},
		})
	})
}

func StartServer() error {
	return r.Run()
}

func main() {
	if err := InitDB(); err != nil {
		fmt.Println("Database error on init!")
		fmt.Println(err.Error())
		return
	}
	InitGin()
	InitRouter()
	if err := StartServer(); err != nil {
		fmt.Println("Server error!")
		fmt.Println(err.Error())
		return
	}
}
