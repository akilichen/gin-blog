package main

import (
	"context"
	"fmt"
	"gin-blog/pkg/setting"
	"gin-blog/routers"
	"log"
	"net/http"
	os2 "os"
	"os/exec"
	"os/signal"
	"time"
)

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os2.Signal)
	signal.Notify(quit, os2.Interrupt)
	<-quit

	log.Println("shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	log.Println("server will begin to restart after 1s...")
	time.Sleep(1 * time.Second)
	cmd := exec.Command("go", "run", "D:/Go-project/gin-blog/main.go")
	cmd.Stdout = os2.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal("Server restart:", err)
	}
	log.Println("Server restart success...")
}
