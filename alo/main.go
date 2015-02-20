package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type matrix [][]int

const (
	matrixN = 10000000
	seed    = 100000
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	m := make([][]int, matrixN, matrixN)
	for i := 0; i < matrixN; i++ {
		m[i] = []int{r.Intn(seed), r.Intn(seed)}
	}
	p := [][]int{{1, 2}, {3, 4}, {5, 4}, {7, 8}, {1, 2}}
	mat := [][]int{{1, 2}, {3, 4}, {4, 6}, {7, 8}, {9, 2}}
	fmt.Println("Running initial CS402Test1 code (Brute Force) with small optimizations.")
	fmt.Println("test matrix (has dups):", p)
	fmt.Println("matrix is unique:", bruteForce(p))
	fmt.Println("test matrix (no dups):", mat)
	fmt.Println("matrix is unique:", bruteForce(mat))

	fmt.Println("\nRunning Divide and Conquer code.")
	sort.Sort(matrix(p))   // some pre-sorting.
	sort.Sort(matrix(mat)) // some pre-sorting.
	fmt.Println("test matrix (has dups):", p)
	fmt.Println("matrix is unique:", divideAndConquer(p))
	fmt.Println("test matrix (no dups):", mat)
	fmt.Println("matrix is unique", divideAndConquer(mat))

	fmt.Printf("\n %d array randomly seeded  of points. between 0 & %d \n", matrixN, seed)
	fmt.Println("matrix is unique", bruteForce(m))
	fmt.Println("\n")
	sort.Sort(matrix(m))
	fmt.Println("matrix is unique", divideAndConquer(m))
}

// Code from CS402Test1 code.
func bruteForce(p [][]int) bool {
	defer timeTrack(time.Now(), "Brute Force")
	n := len(p)
	for a := 0; a < n-1; a++ {
		for b := a + 1; b < n; b++ {
			k := len(p[b])
			flag := false
			for c := 0; c < k; c++ {
				if p[a][c] != p[b][c] {
					flag = true
					break
				}
			}
			if flag == false {
				return false
			}
		}
	}
	return true
}

// The divide and Conquer code. It uses closest pair algorithm.
// If the closest pair is 0, we know we have a duplicate point.
func divideAndConquer(p [][]int) bool {
	defer timeTrack(time.Now(), "Divide & Conquer")
	i := closest(p)
	if i == 0 {
		return false
	}
	return true
}

// closest returns the smallest distance between two points.
func closest(p [][]int) int {
	return int(closestUtil(p))
}

func closestUtil(p [][]int) float64 {
	// If there are 2 or 3 points, then use brute force (bruteForceF in our case.)
	if len(p) <= 3 {
		return bruteForceF(p)
	}

	// Find the middle point
	n := len(p) / 2
	midPoint := p[n]

	// Consider the vertical line passing through the middle point
	// calculate the smallest distance dl on left of middle point and
	// dr on right side.

	dl := closestUtil(p[:n])
	dr := closestUtil(p[n:])

	// Find the smaller of two distances.
	d := min(dl, dr)

	// Build an array that contains points close (closer than d)
	// to the line passing through the middle point
	strip := make([][]int, n)
	j := 0
	for i := 0; i < n; i++ {
		if math.Abs(float64(p[i][0]-midPoint[0])) < d {
			strip[j] = p[i]
			j++
		}
	}

	// Find the closest points in strip.  Return the minimum of d and closest
	// distance is strip[]
	return min(d, stripClosest(strip, j, d))
}

// // min finds the minimum of two float values.
func min(x, y float64) float64 {
	if x < y {
		return x
	} else {
		return y
	}
}

// stripClosest is a utility function to find the distance between the closest points of
// strip of given size. All points in strip[] are sorted according to
// y coordinate. They all have an upper bound on minimum distance as d.
// Note that this method seems to be a O(n^2) method, but it's a O(n)
// method as the inner loop runs 6 times at the very most.
func stripClosest(strip [][]int, size int, d float64) float64 {
	min := d // Initialize the minimum distance as d

	// Pick all points one by one and try the next points till the difference
	// between y coordinates is smaller than d.
	// This is a proven fact that this loop runs at most 6 times
	for i := 0; i < size; i++ {
		for j := i + 1; j < size && (float64(strip[j][1]-strip[i][1])) < min; j++ {
			if dist(strip[i], strip[j]) < min {
				min = dist(strip[i], strip[j])
			}
		}
	}

	return min
}

// dist calculates the square root between the x and y coordinates of each point
// to return the distance between the two points.
func dist(p1, p2 []int) float64 {
	return math.Sqrt(float64((p1[0]-p2[0])*(p1[0]-p2[0]) +
		(p1[1]-p2[1])*(p1[1]-p2[1])))
}

// bruteForceF is actually used by the divideAndConquer version of closest pairs
// if n is less than 3. See chapter 5 for more details.
func bruteForceF(p [][]int) float64 {
	n := len(p)
	min := math.MaxFloat64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if dist(p[i], p[j]) < min {
				min = dist(p[i], p[j])
			}
		}
	}

	return min
}

// Implemented the methods below so we can satisfy the sorting interface.
// This will all us to have a nice presorted list of points in nondecreasing order (just like the book!).
func (m matrix) Len() int {
	return len(m)
}
func (m matrix) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m matrix) Less(i, j int) bool {
	k := len(m[i])
	for c := 0; c < k; c++ {
		if m[i][c] == m[j][c] {
			if c == k-1 {
				return true
			}
		} else {
			return m[i][c] < m[j][c]
		}
	}
	return m[i][0] < m[j][0]
}

// timeTrack starts a timer when the function starts and then prints the value
// when the function ends.
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}
