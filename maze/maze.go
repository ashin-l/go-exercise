package main

import (
	"fmt"
	"os"
)

func readMaze(filename string) ([][]int, []point) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	var row, col int
	fmt.Fscanf(file, "%d %d", &row, &col)

	route := make([]point, 0, row*col)
	maze := make([][]int, row)
	for i := range maze {
		maze[i] = make([]int, col)
		for j := range maze[i] {
			fmt.Fscanf(file, "%d", &maze[i][j])
		}
	}
	return maze, route
}

type point struct {
	i, j int
}

func (p point) add(r point) point {
	return point{p.i + r.i, p.j + r.j}
}

var dirs = [4]point{
	{-1, 0}, {0, -1}, {1, 0}, {0, 1},
}

func (p point) at(grid [][]int) (int, bool) {
	if p.i < 0 || p.i >= len(grid) {
		return 0, false
	}
	if p.j < 0 || p.j >= len(grid[p.i]) {
		return 0, false
	}
	return grid[p.i][p.j], true
}

func walk(maze [][]int, start, end point) ([][]int, bool) {
	steps := make([][]int, len(maze))
	bret := false
	for i := range steps {
		steps[i] = make([]int, len(maze[i]))
	}

	Q := []point{start}

	for len(Q) > 0 {
		cur := Q[0]
		Q = Q[1:]

		if cur == end {
			bret = true
			break
		}

		for _, dir := range dirs {
			next := cur.add(dir)

			val, ok := next.at(maze)
			if !ok || val == 1 {
				continue
			}

			val, ok = next.at(steps)
			if !ok || val != 0 {
				continue
			}

			if next == start {
				continue
			}

			curSteps, _ := cur.at(steps)
			steps[next.i][next.j] = curSteps + 1

			Q = append(Q, next)
		}
	}
	return steps, bret
}

func getRoute(route []point, steps [][]int) []point {
	nextval, _ := route[len(route)-1].at(steps)
	nextval--
	for nextval > 0 {
		cur := route[len(route)-1]
		for _, dir := range dirs {
			next := cur.add(dir)
			val, ok := next.at(steps)
			if ok && (val == nextval) {
				nextval--
				route = append(route, next)
				break
			}
		}
	}
	return route
}

func main() {
	maze, route := readMaze("maze.in")
	start, end := point{0, 0}, point{len(maze) - 1, len(maze[0]) - 1}

	steps, ok := walk(maze, start, end)
	if !ok {
		fmt.Println("Sorry! Do not walk out the maze!")
		return
	}
	route = append(route, end)
	route = getRoute(route, steps)
	route = append(route, start)
	stepnum, _ := end.at(steps)
	fmt.Printf("Walk out maze need %d steps!\n", stepnum)
	fmt.Println(route)
	for _, row := range steps {
		for _, val := range row {
			fmt.Printf("%3d", val)
		}
		fmt.Println()
	}
}
