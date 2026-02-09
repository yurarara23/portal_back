package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Member struct {
    gorm.Model
    UserName      string `json:"username"`
    Password      string `json:"password"`
    Department    string `json:"department"`
    Grade         string `json:"grade"`
    StudentNumber string `json:"student_num"`
}

type Rental struct {
	gorm.Model
	MemberID uint   // MembersテーブルのIDと紐付く（外部キー）
    ItemName string // 借りたもの
	ImageURL string `json:"image_url"`
	IsRented bool `json:"isRented"`
}

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/nu_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("データベース接続に失敗しました")
	}

	db.AutoMigrate(&Member{}, &Rental{})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())


	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Backend!")
	})

	// ログインして「チケット」を渡す処理
	// 6. ログイン (POST /login)
    e.POST("/login", func(c echo.Context) error {
        // フロントから送られてくるログイン情報用の構造体
        type LoginRequest struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        req := new(LoginRequest)
        if err := c.Bind(req); err != nil {
            return err
        }

        var member Member
        // DBから「ユーザー名」と「パスワード」が一致する人を探す
        result := db.Where("user_name = ? AND password = ?", req.Username, req.Password).First(&member)

        if result.Error != nil {
            // 見つからなければ401（認証エラー）を返す
            return c.JSON(http.StatusUnauthorized, map[string]string{"message": "ユーザー名またはパスワードが違います"})
        }

        // ログイン成功！ ユーザー情報を返す（本来はここでトークンを発行する）
        return c.JSON(http.StatusOK, map[string]interface{}{
            "message": "ログイン成功",
            "user":    member,
        })
    })

	e.POST("/members", func(c echo.Context) error {
        m := new(Member)
        if err := c.Bind(m); err != nil {
            return err
        }
        db.Create(&m)
        return c.JSON(http.StatusCreated, m)//そのまま返して登録されたことを見せてくれる
    })

	e.GET("/members", func(c echo.Context) error {
        var members []Member
        // DBから全員取得
        db.Find(&members)
        return c.JSON(http.StatusOK, members)
    })

	// 3. 備品を借りる (POST /rentals)
    e.POST("/rentals", func(c echo.Context) error {
        r := new(Rental)
        if err := c.Bind(r); err != nil {
            return err
        }
        
        // 借りたフラグを立てる
        r.IsRented = true
        
        // DBに保存（ここでMemberIDが保存されるので、誰が借りたか紐付きます）
        if err := db.Create(&r).Error; err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "保存に失敗しました"})
        }
        
        return c.JSON(http.StatusCreated, r)
    })

    // 4. 誰が何を借りているか一覧を表示 (GET /rentals)
    e.GET("/rentals", func(c echo.Context) error {
        var rentals []Rental
        // Preload("Member") を使うと、Memberの情報も一緒に取ってこれますが、
        // まずは単純に一覧を取得します
        db.Find(&rentals)
        return c.JSON(http.StatusOK, rentals)
    })

	// 5. 借りる・返すを切り替える (PATCH /rentals/toggle)
    e.PATCH("/rentals/toggle", func(c echo.Context) error {
        // リクエストの中身（誰が何を）を受け取るための構造体
        type ToggleRequest struct {
            MemberID uint   `json:"member_id"`
            ItemName string `json:"item_name"`
        }
        req := new(ToggleRequest)
        if err := c.Bind(req); err != nil {
            return err
        }

        var rental Rental
        // DBから「その人が借りているそのアイテム」を探す
        // 未返却(is_rented = true)のものを優先的に探す
        result := db.Where("member_id = ? AND item_name = ?", req.MemberID, req.ItemName).Last(&rental)

        if result.Error != nil {
            return c.JSON(http.StatusNotFound, map[string]string{"message": "記録が見つかりません"})
        }

        // 状態を反転させる (true ↔ false)
        rental.IsRented = !rental.IsRented

        // DBを更新
        db.Save(&rental)

        status := "返却しました"
        if rental.IsRented {
            status = "借りました"
        }

        return c.JSON(http.StatusOK, map[string]interface{}{
            "message":   status,
            "isRented":  rental.IsRented,
            "item_name": rental.ItemName,
        })
    })

	// 7. 特定の備品の現在のレンタル状態を取得 (GET /rentals/status)
	e.GET("/rentals/status", func(c echo.Context) error {
    	itemName := c.QueryParam("item_name")
		
    	var rental Rental
    	// そのアイテムの最新のレコードを1件取得
    	result := db.Where("item_name = ?", itemName).Last(&rental)
    
    	if result.Error != nil {
     	   // まだ一度も借りられていない場合
     	   return c.JSON(http.StatusOK, map[string]interface{}{
     	       "isRented": false,
     	       "username": "",
     	   })
   		}

   	 	// 借りている人の名前も一緒に返したい場合
    	var member Member
    	db.First(&member, rental.MemberID)

    	return c.JSON(http.StatusOK, map[string]interface{}{
        	"isRented": rental.IsRented,
        	"username": member.UserName,
    	})
	})

	fmt.Println("Server started at :8080")
	e.Logger.Fatal(e.Start(":8080"))
}