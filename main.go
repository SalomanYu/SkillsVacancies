package main

import (
	"fmt"
	"github.com/SalomanYu/SkillsVacancies/joiner"
	"time"
)


func main(){
	start := time.Now().Unix()
	// joiner.JoinEdwicaProfessions()
	joiner.CombineTheFoundPairsWithVacancies()
	fmt.Println(time.Now().Unix()-start, " sec.")

}