package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
)

func genField(y, x int) [][]string {
	// Generate a slice of slices
	field := make([][]string, y)

	for i := 0; i < y; i++ {
		// Generate the slices of slice i
		field[i] = make([]string, x)
		for j := 0; j < x; j++ {
			// Fill the field
			field[i][j] = " "
		}
	}
	return field
}

func genEnemies(field *[][]string, enemyLook string, p Player) {
	// Spawn new Enemies
	for j := 0; j < len((*field)[0]); j++ {
		if rand.Intn(10-1) >= 9-p.lvl {
			(*field)[0][j] = enemyLook
		}
	}
}

func calcEnemies(field *[][]string, enemyLook string) {
	// Go through all rows, except the last
	for i := len(*field) - 3; i >= 0; i-- {
		for j := 0; j < len((*field)[i]); j++ {
			// Scroll down every row (enemies etc)
			if (*field)[i][j] == enemyLook {
				(*field)[i+1][j] = (*field)[i][j]
				(*field)[i][j] = " "
			}

		}
	}
}

func chkEnemies(field [][]string, enemyLook string) bool {
	// Check for remaining enemies on the field
	for i := len(field) - 3; i >= 0; i-- {
		for j := 0; j < len((field)[i]); j++ {
			if field[i][j] == enemyLook {
				// If enemy was found, return false
				return false
			}
		}
	}
	// if no enemy was found, return true
	return true
}

func genGunFire(field *[][]string, p Player) {
	yMax := len(*field) - 3
	xMax := len((*field)[yMax])

	switch p.gunLvl {
	case 1, 2, 3:
		// Look if we can set the fires left and right
		// If not, set it on the other side
		for i := 0; i < p.gunLvl; i++ {
			if p.xPos-i < 0 {
				(*field)[yMax][xMax-i+p.xPos] = p.gunLook
			} else {
				(*field)[yMax][p.xPos-i] = p.gunLook
			}
			if p.xPos+i >= xMax {
				(*field)[yMax][p.xPos+i-xMax] = p.gunLook
			} else {
				(*field)[yMax][p.xPos+i] = p.gunLook
			}
		}
	}

}

func calcGunFire(field *[][]string, enemyLook string, p *Player) {
	scrollGunfire := true

	// Go through all rows, except the player one, backwards
	for i := len(*field) - 3; i >= 0; i-- {
		// Scroll up Gunfires and set points
		if scrollGunfire {
			switch p.gunLvl {
			case 1, 2, 3:
				// Search for Gunfire points
				for j := 0; j < len((*field)[i]); j++ {
					if (*field)[i][j] == p.gunLook {
						// Scroll them up
						if i == 0 {
							(*field)[i][j] = " "
						} else if (*field)[i-1][j] == enemyLook {
							(*field)[i-1][j] = " "
							(*field)[i][j] = " "
							p.points++
						} else {
							(*field)[i-1][j] = p.gunLook
							(*field)[i][j] = " "
						}
						scrollGunfire = false
					}
				}
			case 4:
				// If SHooting in multiple directions
			}
		}
	}
}

func checkDamage(field [][]string, p *Player, enemyLook string) {
	yPoint := len(field) - 2
	for i := 0; i < len(field[yPoint]); i++ {
		if field[yPoint][i] == enemyLook {
			field[yPoint][i] = " "
			p.hp--
		}
	}
}

func setPlayerPosition(lr string, field *[][]string, p *Player) {
	// Go only through the player row
	last := len(*field) - 2

	// Change the player position
	switch lr {
	case "left":
		// if Player already on left bounds, send them to the right
		if (*field)[last][1] != " " {
			// IF IN BONUS LEVEL
			/*(*field)[last][len((*field)[last])-1] = (*field)[last][1]
			(*field)[last][1] = " "
			p.xPos = len((*field)[last]) - 1
			*/
		} else {
			for i := 1; i <= len((*field)[last])-1; i++ {
				if i == p.xPos {
					// Set every object to the next left column
					(*field)[last][i-1] = p.look
					(*field)[last][i] = " "
				}
			}
			p.xPos--
		}
	case "right":
		// if Player already on right bounds, send them to the left
		if (*field)[last][len((*field)[last])-1] != " " {
			// if in bonus level
			/*		(*field)[last][1] = (*field)[last][len((*field)[last])-1]
					(*field)[last][len((*field)[last])-1] = " "
					p.xPos = 1
			*/
		} else {
			for i := len((*field)[last]) - 2; i >= 0; i-- {
				if i == p.xPos {
					// Set the player to the next right column
					(*field)[last][i+1] = p.look
					(*field)[last][i] = " "
				}
			}
			p.xPos++
		}
	}
}

func bonus(p *Player) {
	if p.gunLvl < 3 {
		p.gunLvl++

		// Change the look of the Gun
		switch p.gunLvl {
		case 2, 3:
			p.gunLook = "|"
		}
	}
}

func setNewLevel(field *[][]string, p *Player) {

	// Generate the "new Level String" on screen
	lvlString := fmt.Sprint("LEVEL ", p.lvl, " COMPLETE!")
	for i := 0; i < len(lvlString); i++ {
		(*field)[3][i+10] = string(lvlString[i])
	}
	// Generate the "press enter string" on screen
	contString := "PRESS ENTER TO CONTINUE"
	for i := 0; i < len(contString); i++ {
		(*field)[4][i+7] = string(contString[i])
	}
	// Reset Player HP
	p.hp = 10
	// Set new level
	p.lvl++
}

func setGameOver(field *[][]string, p *Player) {

	// Generate the "Game Over" String on screen
	lvlString := fmt.Sprint("GAME OVER!")
	for i := 0; i < len(lvlString); i++ {
		(*field)[3][i+14] = string(lvlString[i])
	}
	// Generate the "press enter string" on screen
	contString := "Press ENTER to restart!"
	for i := 0; i < len(contString); i++ {
		(*field)[4][i+7] = string(contString[i])
	}
	// Reset Player HP
	p.hp = 10
	// reset level
	p.lvl = 1
	// reset points
	p.points = 0
}

func valueBar(field *[][]string, points int, lvl int, hp int) {
	// Set a counter and a temporary step counter to watch the horizontal state
	ctr := 0
	steps := 0
	// How much whitespaces between the strings?
	ws := 3

	// Set the Points string
	pString := fmt.Sprint("Points: ", points)
	for i := 0; i < len(pString); i++ {
		(*field)[len(*field)-1][i] = string(pString[i])
		ctr++
	}

	// Get the point digits
	pd := 0
	switch {
	case points >= 10 && points <= 99:
		pd = 1
	case points >= 100 && points <= 999:
		pd = 2
	case points >= 1000 && points <= 9999:
		pd = 3
	}

	// Set whitespaces
	for i := 0; i < ws-pd; i++ {
		(*field)[len(*field)-1][i+ctr] = " "
		steps++
	}
	ctr += steps
	steps = 0
	// Create HP line
	hpString := "HP: "
	for i := 0; i <= 10; i++ {
		if hp >= i {
			hpString = fmt.Sprint(hpString, "O")
		}
	}
	// Add hp String
	for i := 0; i < len(hpString); i++ {
		(*field)[len(*field)-1][i+ctr] = string(hpString[i])
		steps++
	}
	ctr += steps
	steps = 0
	// Add whitespaces
	for i := ctr; i <= 30; i++ {
		(*field)[len(*field)-1][i] = " "
		steps++
	}
	ctr += steps
	steps = 0
	// Add level String
	lString := fmt.Sprint("Level: ", lvl)
	for i := 0; i < len(lString); i++ {
		(*field)[len(*field)-1][i+30] = string(lString[i])
		steps++
	}
	ctr += steps
	// Fill the rest of the field with whitespaces
	for i := ctr; i < len((*field)[len(*field)-1]); i++ {
		(*field)[len(*field)-1][i] = " "
	}
}

func keyPress(keyP *string) {

	keyboard.Open()
	defer keyboard.Close()
	// Wait for the next keypress
	_, key, _ := keyboard.GetKey()

	switch key {
	case keyboard.KeyArrowLeft:
		*keyP = "left"
	case keyboard.KeyArrowRight:
		*keyP = "right"
	case keyboard.KeySpace:
		*keyP = "space"
	case keyboard.KeyEnter:
		*keyP = "enter"
	case keyboard.KeyEsc:
		*keyP = "ESC"
		fmt.Println("Pressed ESC key, exiting.")
	default:
		*keyP = "none"
	}
}

type Player struct {
	look    string
	xPos    int
	hp      int
	gunLvl  int
	gunLook string
	points  int
	lvl     int
}

func main() {
	xField := 40
	yField := 15

	// Generate playing field
	field := genField(yField, xField)

	// Setup Look
	bonusLook := "°"
	enemyLook := "X"

	// Setup player
	p := Player{"A", xField / 2, 10, 1, "^", 0, 1}
	field[yField-2][p.xPos] = p.look

	keyPressed := "none"
	spawnCounter := 0
	gunCounter := 0
	waveCounter := 0
	lvlWait := false
	gameOverWait := false

	// Turn off cursor and on if program ends
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h")

	for {
		// Scroll GunFire
		if gunCounter == 7 || (p.gunLvl > 1 && gunCounter == 1) {
			go func() {
				calcGunFire(&field, enemyLook, &p)
			}()
			gunCounter = 0
		} else {
			gunCounter++
		}

		// Look for space triggers first for parallel processing of fire and movements
		switch keyPressed {
		case "space":
			genGunFire(&field, p)
			// Open a new thread, so that other keys can be pressed in that run
			go func() {
				keyPress(&keyPressed)
			}()
			// Safety function - if nothing pressed, the function will be called in the loop
			keyPressed = "waiting"
		}

		// Look for keypresses
		switch keyPressed {
		case "none":
			go func() {
				keyPress(&keyPressed)
			}()
			// Safety function - if nothing pressed, the function will be called in the loop
			keyPressed = "waiting"
		case "left":
			// If Player hits bonus, change gunlevel
			if p.xPos-1 > 0 {
				if field[len(field)-2][p.xPos-1] == bonusLook {
					bonus(&p)
				}
			}
			setPlayerPosition("left", &field, &p)
			keyPressed = "none"
		case "right":
			// If Player hits bonus, change gunlevel (if below 2)
			if p.xPos+1 < len(field[len(field)-1]) {
				if field[len(field)-2][p.xPos+1] == bonusLook {
					bonus(&p)
				}
			}
			setPlayerPosition("right", &field, &p)
			keyPressed = "none"
		case "enter":
			if lvlWait == false && gameOverWait == false {
				keyPressed = "none"
			}
		case "ESC":
			return
		}

		// Rotate and Spawn Enemies every 5 Seconds
		if spawnCounter == 200-yField-(p.lvl*10) {
			calcEnemies(&field, enemyLook)
			waveCounter++
			// Check if some enemies hit the ground
			checkDamage(field, &p, enemyLook)
			if waveCounter <= yField-2 {
				genEnemies(&field, enemyLook, p)
			} else if p.hp <= 0 {
				// Set Game Over if HP down
				gameOverWait = true
				field = nil
				field = genField(yField, xField)
				setGameOver(&field, &p)
			} else if gameOverWait {
				// If Player wants to restart
				if keyPressed == "enter" {
					keyPressed = "none"

					waveCounter = 0
					gameOverWait = false
				}
			} else if chkEnemies(field, enemyLook) == true && lvlWait == false {
				// If all enemies are blown
				// Set the next level and a new wave
				setNewLevel(&field, &p)
				lvlWait = true
			} else if lvlWait {
				// If level set wait for user to press space
				if keyPressed == "enter" {
					keyPressed = "none"
					field = genField(yField, xField)
					waveCounter = 0
					lvlWait = false
				}
			}
			spawnCounter = 0
		} else {
			spawnCounter++
		}

		// Generate Bonus
		if rand.Intn(10000-0) == 999 {
			field[yField-2][rand.Intn(39-0)] = bonusLook
		}

		// Refresh value bar (Health Bar etc)
		valueBar(&field, p.points, p.lvl, p.hp)

		// Draw Game field
		for row := 0; row < yField; row++ {
			for col := 0; col < xField; col++ {
				fmt.Printf("%v", field[row][col])
			}
			fmt.Printf("\n")
		}

		// Let the Game field stay at the same position on console
		for row := 0; row < yField; row++ {
			fmt.Printf("\033[1A")
		}

		time.Sleep(4 * time.Millisecond)
	}
}
