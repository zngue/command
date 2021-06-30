package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zngue/go_helper/pkg"
	"github.com/zngue/go_helper/pkg/sign_chan"
	"io/ioutil"
	"log"
	"os/exec"
)

func main() {

	if conErr := pkg.NewConfig(); conErr != nil {
		log.Fatal(conErr)
		return
	}
	port := viper.GetString("AppPort")
	run, errs := pkg.GinRun(port, func(engine *gin.Engine) {

		commands := engine.Group("command")
		commands.POST("shell", func(c *gin.Context) {
			all, _ := ioutil.ReadAll(c.Request.Body)
			query := c.DefaultQuery("typeName", "")
			if query == "" {
				c.JSON(200, gin.H{
					"code": 100,
				})
				return
			}
			command := fmt.Sprintf("./shell/%s.sh ", query)
			cmd := exec.Command("/bin/bash", "-c", command)
			output, err := cmd.Output()
			if err != nil {
				c.JSON(200, gin.H{
					"code":    100,
					"message": err.Error(),
				})
				return
			} else {
				c.JSON(200, gin.H{
					"code":    200,
					"message": string(output),
					"data":    string(all),
				})
				return
			}
		})
	})
	if errs != nil {
		sign_chan.SignLog(errs)
	}
	go func() {
		err := run.ListenAndServe()
		if err != nil {
			sign_chan.SignLog(err)
		}
	}()
	sign_chan.SignChalNotify()
	sign_chan.ListClose(func(ctx context.Context) error {
		return run.Shutdown(ctx)
	})

}
