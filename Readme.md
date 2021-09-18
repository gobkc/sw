# SW: A simple web Framework

-----

sw try to use the official library as much as possible

#### step 1. Hello World

````
    app := sw.NewApp()
    app.GET("/hello1", func(this *sw.This) {
        data := sw.M{
            "hello": "world",
        }
        this.Json(http.StatusOK, data)
    })
    if err := app.Run(":8999"); err != nil {
        log.Println(err.Error())
    }
````

#### step 2. Use routing group

````
    app := sw.NewApp()
    group1 := app.Group("/test1/")
    {
        group1.GET("/hello1", func(this *sw.This) {
            data := sw.M{
                "hello": "world",
            }
            this.Json(http.StatusOK, data)
        })
    }
    if err := app.Run(":8999"); err != nil {
        log.Println(err.Error())
    }
````

#### step 3. Cross domain -> using built-in middleware

````
    app := sw.NewApp()
    app.Use(sw.Cors()).GET("/hello1", func(this *sw.This) {
        data := sw.M{
            "hello": "world",
        }
        this.Json(http.StatusOK, data)
    })
    if err := app.Run(":8999"); err != nil {
        log.Println(err.Error())
    }
````

#### step 4. Basic Auth -> using built-in middleware

````
    app := sw.NewApp()
    app.Use(sw.BasicAuth("abc","123")).GET("/hello1", func(this *sw.This) {
        data := sw.M{
            "hello": "world",
        }
        this.Json(http.StatusOK, data)
    })
    if err := app.Run(":8999"); err != nil {
        log.Println(err.Error())
    }
````

#### step 5. Support Vue -> using built-in routing

- First, define public variables

````
//go:embed static/*
var Files embed.FS
````

- using built-in routing

````
    app := sw.NewApp()
    app.VUE("./static", &Files)
    if err := app.Run(":8999"); err != nil {
        log.Println(err.Error())
    }
````

#### step 6. Get routing parameters

````
    app := sw.NewApp()
    app.GET("/users/:name", func(this *sw.This) {
        user := this.Param("name","zhang")
        this.Html(http.StatusOK, user)
    })
    if err := app.Run(":8999"); err != nil {
        log.Println(err.Error())
    }

````

#### step 7. "GET" parameters

````
    app := sw.NewApp()
    app.GET("/list", func(this *sw.This) {
        pg := this.Path("page","1")
        this.Html(http.StatusOK, pg)
    })
    if err := app.Run(":8999"); err != nil {
        log.Println(err.Error())
    }

````

- testing

````
curl http://localhost:8999/list?page=2
````

#### step 8. "POST" parameters(form)

````
    app := sw.NewApp()
    app.POST("/new-user", func(this *sw.This) {
        param := Param1{}
        if err := this.ParseForm(&param); err != nil {
            log.Println(err)
        }
        this.Json(http.StatusOK, param)
    })
    if err := app.Run(":8999"); err != nil {
        log.Println(err.Error())
    }
````

#### step 9. "POST" parameters(json)

````
    app := sw.NewApp()
    app.POST("/new-user", func(this *sw.This) {
        param := Param1{}
        if err := this.ParseJson(&param); err != nil {
            log.Println(err)
        }
        this.Json(http.StatusOK, param)
    })
    if err := app.Run(":8999"); err != nil {
        log.Println(err.Error())
    }
````

