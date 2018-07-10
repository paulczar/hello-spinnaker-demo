// The MIT License (MIT)

// Copyright (c) 2016 Claudemiro

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

/*
Package health is a easy to use, extensible health check library.

Example

    package main

    import (
        "net/http"
        "database/sql"

        "github.com/dimiro1/health"
        "github.com/dimiro1/health/url"
        "github.com/dimiro1/health/db"
        "github.com/dimiro1/health/redis"
        _ "github.com/go-sql-driver/mysql"
    )

    func main() {
        database, _ := sql.Open("mysql", "/")
        mysql := db.NewMySQLChecker(database)

        companies := health.NewCompositeChecker()
        companies.AddChecker("Microsoft", url.NewChecker("https://www.microsoft.com/"))
        companies.AddChecker("Oracle", url.NewChecker("https://www.oracle.com/"))
        companies.AddChecker("Google", url.NewChecker("https://www.google.com/"))

        handler := health.NewHandler()
        handler.AddChecker("Go", url.NewChecker("https://golang.org/"))
        handler.AddChecker("Big Companies", companies)
        handler.AddChecker("MySQL", mysql)
        handler.AddChecker("Redis", redis.NewChecker("tcp", ":6379"))

        http.Handle("/health/", handler)
        http.ListenAndServe(":8080", nil)
    }

Executing a curl

    $ curl localhost:8080/health/

If everything is ok the server must respond with HTTP Status 200 OK and have following json in the body.

    {
        "Big Companies": {
            "Google": {
                "code": 200,
                "status": "UP"
            },
            "Microsoft": {
                "code": 200,
                "status": "UP"
            },
            "Oracle": {
                "code": 200,
                "status": "UP"
            },
            "status": "UP"
        },
        "Go": {
            "code": 200,
            "status": "UP"
        },
        "MySQL": {
            "status": "UP",
            "version": "10.1.9-MariaDB"
        },
        "Redis": {
            "status": "UP",
            "version": "3.0.5"
        },
        "status": "UP"
    }


The server responds with HTTP Status 503 Service Unavailable if the ckeck is Down and the json response could be something like this.

    {
        "Big Companies": {
            "Google": {
                "code": 200,
                "status": "UP"
            },
            "Microsoft": {
                "code": 200,
                "status": "UP"
            },
            "Oracle": {
                "code": 200,
                "status": "UP"
            },
            "status": "UP"
        },
        "Go": {
            "code": 200,
            "status": "UP"
        },
        "MySQL": {
            "status": "DOWN",
            "error": "Error 1044: Access denied for user ''@'localhost' to database 'invalid-database'",
        },
        "status": "DOWN"
    }

Motivation

It is very important to verify the status of your system, not only the system itself, but all its dependencies,
If your system is not Up you can easily know what is the cause of the problem only looking the health check.

Also it serves as a kind of basic itegration test between the systems.

Inspiration

I took a lot of ideas from the spring framework (http://spring.io/).

Installation

This package is a go getable packake.

    $ go get github.com/dimiro1/health

API

The API is stable and I do not have any plans to break compatibility, but I recommend you to vendor this dependency in your project, as it is a good practice.

Testing

You have to install the test dependencies.


    $ go get gopkg.in/DATA-DOG/go-sqlmock.v1

or you can go get this package with the -t flag

    $ go get -t github.com/dimiro1/health

Implementing custom checkers

The key interface is health.Checker, you only have to implement a type that satisfies that interface.

    type Checker interface {
        Check() Health
    }

Here an example of Disk Space usage (unix only).

    package main

    import (
        "syscall"
        "os"
    )

    type DiskSpaceChecker struct {
        Dir       string
        Threshold uint64
    }

    func NewDiskSpaceChecker(dir string, threshold uint64) DiskSpaceChecker {
        return DiskSpaceChecker{Dir: dir, Threshold: threshold}
    }

    func (d DiskSpaceChecker) Check() health.Health {
        health := health.NewHealth()

        var stat syscall.Statfs_t

        wd, err := os.Getwd()

        if err != nil {
            health.Down().AddInfo("error", err.Error()) // Why the check is Down
            return health
        }

        syscall.Statfs(wd, &stat)

        diskFreeInBytes := stat.Bavail * uint64(stat.Bsize)

        if diskFreeInBytes < d.Threshold {
            health.Down()
        } else {
            health.Up()
        }

        health.
            AddInfo("free", diskFreeInBytes).
            AddInfo("threshold", d.Threshold)

        return health
    }

Important

The **status** key in the json have priority over a "status" key added by a Checker, so if some checker add a "status" key to the json, it will not be rendered

*/
package health
