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
			if j == x-1 || j == 0 {
				field[i][j] = "|"
			} else if i == y-4 {
				field[i][j] = "_"
			} else {
				// Fill the field
				field[i][j] = " "
			}
		}
	}
	return field
}

func genEnemies(field *[][]string, eLook string, p Player) {
	// Spawn new Enemies
	for j := 1; j < len((*field)[0])-1; j++ {
		if rand.Intn(10-1) >= 9-p.lvl {
			(*field)[0][j] = eLook
		}
	}
}

func calcEnemies(field *[][]string, eLook string) {
	// Go through all rows, except the last
	for i := len(*field) - 3; i >= 0; i-- {
		for j := 1; j < len((*field)[i])-1; j++ {
			// Scroll down every row (enemies etc)
			if (*field)[i][j] == eLook {
				(*field)[i+1][j] = (*field)[i][j]
				(*field)[i][j] = " "
			}

		}
	}
}

func chkEnemies(field [][]string, eLook string) bool {
	// Check for remaining enemies on the field
	for i := len(field) - 3; i >= 0; i-- {
		for j := 0; j < len((field)[i]); j++ {
			if field[i][j] == eLook {
				// If enemy was found, return false
				return false
			}
		}
	}
	// if no enemy was found, return true
	return true
}

func genGunFire(field *[][]string, p Player) {
	yMax := len(*field) - 4
	xMax := len((*field)[yMax]) - 1

	switch p.gunLvl {
	case 1, 2:
		(*field)[yMax][p.xPos] = p.gunLook
	case 3, 4:
		// Look if we can set the fires left and right
		// If not, set it on the other side
		for i := 0; i < p.gunLvl-1; i++ {
			if p.xPos-i < 1 {
				(*field)[yMax][xMax-i+p.xPos-1] = p.gunLook
			} else {
				(*field)[yMax][p.xPos-i] = p.gunLook
			}
			if p.xPos+i >= xMax {
				(*field)[yMax][p.xPos+i-xMax+1] = p.gunLook
			} else {
				(*field)[yMax][p.xPos+i] = p.gunLook
			}
		}
	}

}

func calcGunFire(field *[][]string, eLook string, p *Player) {
	lastY := len(*field) - 4

	// Scroll up Gunfires and set points
	switch p.gunLvl {
	case 1, 3, 4:
		// Go through all rows, except the player one, backwards
		for i := 0; i <= lastY; i++ {
			for j := 1; j < len((*field)[i])-1; j++ {
				// Search for Gunfire points
				if (*field)[i][j] == p.gunLook {
					if i != 0 {
						// Scroll up and blow enemy
						if (*field)[i-1][j] == eLook {
							(*field)[i-1][j] = " "
							p.points++
						} else {
							// Only scroll up
							(*field)[i-1][j] = p.gunLook
						}
					}
					(*field)[i][j] = " "
				}
			}
		}
	case 2:
		// Go through all rows, except the player one, backwards
		for i := 0; i <= lastY; i++ {
			for j := 1; j < len((*field)[i])-1; j++ {
				// Search for gunfire points
				if (*field)[i][j] == p.gunLook {
					if i == 1 {
						(*field)[i-1][j] = p.gunLook
					}
					if i >= 2 {
						// Search for enemies to lock on
						switch {
						case (*field)[i-2][j-1] == eLook:
							(*field)[i-1][j-1] = p.gunLook
						case (*field)[i-2][j] == eLook:
							(*field)[i-1][j] = p.gunLook
						case (*field)[i-2][j+1] == eLook:
							(*field)[i-1][j+1] = p.gunLook
						default:
							// Scroll normal and blow enemy
							if (*field)[i-1][j] == eLook {
								(*field)[i-1][j] = " "
							} else {
								// or just scroll up
								(*field)[i-1][j] = p.gunLook
							}
						}
					}
					(*field)[i][j] = " "
				}
			}
		}
		// If SHooting in multiple directions
	}
	// Refresh the baricade
	for j := 1; j < len((*field)[lastY])-1; j++ {
		(*field)[lastY][j] = "_"
	}

}

func checkDamage(field [][]string, p *Player, eLook string) {
	lastY := len(field) - 4
	for i := 0; i < len(field[lastY]); i++ {
		if field[lastY][i] == eLook {
			field[lastY][i] = " "
			p.hp--
		}
	}
}

func setPlayerPosition(lr string, field *[][]string, p *Player) {
	// Go only through the player row
	last := len(*field) - 3

	// Change the player position
	switch lr {
	case "left":
		// if Player already on left bounds, send them to the right
		if (*field)[last][1] == p.look {
			// IF IN BONUS LEVEL
			/*(*field)[last][len((*field)[last])-1] = (*field)[last][1]
			(*field)[last][1] = " "
			p.xPos = len((*field)[last]) - 1
			*/
		} else {
			for i := 1; i <= len((*field)[last])-2; i++ {
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
		if (*field)[last][len((*field)[last])-2] == p.look {
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

func bonus(p *Player, e *Enemy, bonusItem string) {
	switch bonusItem {
	case "G":
		if p.gunLvl < 4 {
			p.gunLvl++

			// Change the look of the Gun
			switch p.gunLvl {
			case 2:
				p.gunLook = "^"
			case 3:
				p.gunLook = "|"
			}
		}
	case "F":
		e.status = "frozen"
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
	// Generate the "Highscore" string on screen
	hsString := fmt.Sprint("Highscore: ", p.points)
	for i := 0; i < len(hsString); i++ {
		(*field)[5][i+7] = string(hsString[i])
	}
	// Generate the "press enter string" on screen
	contString := "Press ENTER to restart!"
	for i := 0; i < len(contString); i++ {
		(*field)[6][i+7] = string(contString[i])
	}
	// Reset Player vars
	p.hp = 10
	p.lvl = 1
	p.points = 0
	p.gunLvl = 1
}

func valueBar(field *[][]string, points int, lvl int, hp int, gLvl int, e Enemy) {
	lastY := len(*field) - 2
	// Set a counter and a temporary step counter to watch the horizontal state
	ctr := 0
	steps := 0
	// How much whitespaces between the strings?
	ws := 3

	// Set the Points string
	pString := fmt.Sprint("Points: ", points)
	for i := 0; i < len(pString); i++ {
		(*field)[lastY][i] = string(pString[i])
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
		(*field)[lastY][i+ctr] = " "
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
		(*field)[lastY][i+ctr] = string(hpString[i])
		steps++
	}
	ctr += steps
	steps = 0
	// Add whitespaces
	for i := ctr; i <= 30; i++ {
		(*field)[lastY][i] = " "
		steps++
	}
	ctr += steps
	steps = 0
	// Add level String
	lString := fmt.Sprint("Level: ", lvl)
	for i := 0; i < len(lString); i++ {
		(*field)[lastY][i+30] = string(lString[i])
		steps++
	}
	ctr += steps
	// Fill the rest of the field with whitespaces
	for i := ctr; i < len((*field)[len(*field)-1]); i++ {
		(*field)[lastY][i] = " "
	}

	// Create the second Bar
	ctr = 0
	steps = 0
	gString := "Gun Mode: "
	switch gLvl {
	case 1:
		gString = fmt.Sprint(gString, "Laser")
	case 2:
		gString = fmt.Sprint(gString, "Heatseeker")
	case 3:
		gString = fmt.Sprint(gString, "Rockets")
	case 4:
		gString = fmt.Sprint(gString, "The BIG one")
	}

	for i := 0; i < len(gString); i++ {
		(*field)[lastY+1][i] = string(gString[i])
		steps++
	}
	ctr += steps
	fString := "    !! FREEZE !!"
	if e.status == "frozen" {
		for i := 0; i < len(fString); i++ {
			(*field)[lastY+1][i+ctr] = string(fString[i])
			steps++
		}
		ctr += steps
	}

	// Fill the rest of the field with whitespaces
	for i := ctr; i < len((*field)[lastY+1]); i++ {
		(*field)[lastY+1][i] = " "
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

type Enemy struct {
	look   string
	status string
}

func main() {
	xField := 40
	yField := 18

	// Generate playing field
	field := genField(yField, xField)

	// Setup player
	p := Player{"A", xField / 2, 10, 1, "|", 0, 1}
	field[yField-3][p.xPos] = p.look

	// Setup Enemy
	e := Enemy{"x", "normal"}

	// Initial "how to play" message

	keyPressed := "none"
	bonusItem := ""
	frozenTimer := 0
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
		if gunCounter == 7 || (p.gunLvl == 3 && gunCounter == 3) || (p.gunLvl > 3 && gunCounter == 2) {
			go func() {
				calcGunFire(&field, e.look, &p)
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
			// If Player hits bonus, set bonus
			if p.xPos-1 > 0 {
				if field[len(field)-3][p.xPos-1] == bonusItem {
					bonus(&p, &e, bonusItem)
				}
			}
			setPlayerPosition("left", &field, &p)
			keyPressed = "none"
		case "right":
			// If Player hits bonus, set bonus
			if p.xPos+1 < len(field[len(field)-1]) {
				if field[len(field)-3][p.xPos+1] == bonusItem {
					bonus(&p, &e, bonusItem)
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
		if spawnCounter == 300-yField-(p.lvl*10) {
			// Check for bonus activities
			if e.status == "normal" {
				calcEnemies(&field, e.look)
			} else if frozenTimer > 0 {
				frozenTimer--
				if frozenTimer == 0 {
					e.status = "normal"
				}
			} else if e.status == "frozen" {
				frozenTimer = 5
			}
			waveCounter++
			// Check if some enemies hit the ground
			checkDamage(field, &p, e.look)
			if waveCounter <= yField-3+p.lvl {
				genEnemies(&field, e.look, p)
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
					field = genField(yField, xField)
					waveCounter = 0
					gameOverWait = false
				}
			} else if chkEnemies(field, e.look) == true && lvlWait == false {
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
		if rand.Intn(1000-0) == 999 {
			bonusRoll := rand.Intn(3 - 0)
			if bonusRoll == 1 {
				field[yField-3][rand.Intn(39-0)] = "G"
				bonusItem = "G"
			}
			if bonusRoll == 2 {
				field[yField-3][rand.Intn(39-0)] = "F"
				bonusItem = "F"
			}
		}

		// Refresh value bar (Health Bar etc)
		valueBar(&field, p.points, p.lvl, p.hp, p.gunLvl, e)

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
