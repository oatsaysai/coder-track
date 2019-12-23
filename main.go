package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	input := make([]string, 0)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	if scanner.Err() != nil {
		panic("")
	}

	// addOrSum(input)
	// transpose(input)
	// expenseTracking(input)
	// expenseTrackingByDay(input)
	findFastestRoute(input)
}

func findFastestRoute(input []string) {

	nodes := make(map[string]*Node, 0)
	for _, line := range input {
		strArr := strings.Split(line, " ")
		if nodes[strArr[0]] == nil {
			nodes[strArr[0]] = &Node{Name: strArr[0]}
		}
	}

	graph := Graph{}
	for _, line := range input {
		strArr := strings.Split(line, " ")
		dis, _ := strconv.Atoi(strArr[2])
		graph.AddEdge(nodes[strArr[0]], nodes[strArr[1]], dis)
	}

	fmt.Println(graph.Dijkstra(nodes["HOME"], nodes["DEST"]))

}

type Graph struct {
	Edges []*Edge
	Nodes []*Node
}

type Edge struct {
	Parent *Node
	Child  *Node
	Cost   int
}

type Node struct {
	Name string
}

// AddEdge adds an Edge to the Graph
func (g *Graph) AddEdge(parent, child *Node, cost int) {
	edge := &Edge{
		Parent: parent,
		Child:  child,
		Cost:   cost,
	}

	g.Edges = append(g.Edges, edge)
	g.AddNode(parent)
	g.AddNode(child)
}

// AddNode adds a Node to the Graph list of Nodes, if the the node wasn't already added
func (g *Graph) AddNode(node *Node) {
	var isPresent bool
	for _, n := range g.Nodes {
		if n == node {
			isPresent = true
		}
	}

	if !isPresent {
		g.Nodes = append(g.Nodes, node)
	}
}

const Infinity = int(^uint(0) >> 1)

// Dijkstra implements THE Dijkstra algorithm
// Returns the shortest path from startNode to all the other Nodes
func (g *Graph) Dijkstra(startNode, stopNode *Node) (shortestPathValue int) {

	// First, we instantiate a "Cost Table", it will hold the information:
	// "From startNode, what's is the cost to all the other Nodes?"
	// When initialized, It looks like this:
	// NODE  COST
	//  A     0    // The startNode has always the lowest cost to itself, in this case, 0
	//  B    Inf   // the distance to all the other Nodes are unknown, so we mark as Infinity
	//  C    Inf
	// ...
	costTable := g.NewCostTable(startNode)

	// An empty list of "visited" Nodes. Everytime the algorithm runs on a Node, we add it here
	var visited []*Node

	// A loop to visit all Nodes
	for len(visited) != len(g.Nodes) {

		// Get closest non visited Node (lower cost) from the costTable
		node := getClosestNonVisitedNode(costTable, visited)

		// Mark Node as visited
		visited = append(visited, node)

		// Get Node's Edges (its neighbors)
		nodeEdges := g.GetNodeEdges(node)

		for _, edge := range nodeEdges {

			// The distance to that neighbor, let's say B is the cost from the costTable + the cost to get there (Edge cost)
			// In the first run, the costTable says it's "Infinity"
			// Plus the actual cost, let's say "5"
			// The distance becomes "5"
			// var distanceToNeighbor int

			distanceToNeighbor := costTable[node] + edge.Cost

			// If the distance above is lesser than the distance currently in the costTable for that neighbor
			if distanceToNeighbor < costTable[edge.Child] && costTable[node] != Infinity {
				// Update the costTable for that neighbor
				costTable[edge.Child] = distanceToNeighbor
			}
		}
	}

	// Make the costTable nice to read :)
	shortestPathValue = -1
	for node, cost := range costTable {
		// fmt.Println(cost)
		if node == stopNode {
			// fmt.Println(cost)
			if shortestPathValue == -1 || cost < shortestPathValue {
				shortestPathValue = cost
			}

		}
	}

	return shortestPathValue
}

// NewCostTable returns an initialized cost table for the Dijkstra algorithm work with
// by default, the lowest cost is assigned to the startNode – so the algorithm starts from there
// all the other Nodes in the Graph receives the Infinity value
func (g *Graph) NewCostTable(startNode *Node) map[*Node]int {
	costTable := make(map[*Node]int)
	costTable[startNode] = 0

	for _, node := range g.Nodes {
		if node != startNode {
			costTable[node] = Infinity
		}
	}

	return costTable
}

// GetNodeEdges returns all the Edges that start with the specified Node
// In other terms, returns all the Edges connecting to the Node's neighbors
func (g *Graph) GetNodeEdges(node *Node) (edges []*Edge) {
	for _, edge := range g.Edges {
		if edge.Parent == node {
			edges = append(edges, edge)
		}
	}

	return edges
}

// getClosestNonVisitedNode returns the closest Node (with the lower cost) from the costTable
// **if the node hasn't been visited yet**
func getClosestNonVisitedNode(costTable map[*Node]int, visited []*Node) *Node {
	type CostTableToSort struct {
		Node *Node
		Cost int
	}
	var sorted []CostTableToSort

	// Verify if the Node has been visited already
	for node, cost := range costTable {
		var isVisited bool
		for _, visitedNode := range visited {
			if node == visitedNode {
				isVisited = true
			}
		}
		// If not, add them to the sorted slice
		if !isVisited {
			sorted = append(sorted, CostTableToSort{node, cost})
		}
	}

	// We need the Node with the lower cost from the costTable
	// So it's important to sort it
	// Here I'm using an anonymous struct to make it easier to sort a map
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Cost < sorted[j].Cost
	})

	return sorted[0].Node
}

func expenseTrackingByDay(input []string) {
	exTypes := []byte{'f', 'g', 'm', 's', 't'}
	total := make(map[int]int64, 0)
	curDay := 3

	for _, line := range input {
		var totalInDay int64
		totalInDay = 0
		for _, exType := range exTypes {
			totalInDay += getSumTypeFromRow(line, exType)
		}
		total[curDay] += totalInDay
		curDay++
		if curDay >= 7 {
			curDay = 0
		}
	}

	fmt.Printf("Mon %d\n", total[0])
	fmt.Printf("Tue %d\n", total[1])
	fmt.Printf("Wed %d\n", total[2])
	fmt.Printf("Thu %d\n", total[3])
	fmt.Printf("Fri %d\n", total[4])
	fmt.Printf("Sat %d\n", total[5])
	fmt.Printf("Sun %d\n", total[6])
}

func expenseTracking(input []string) {
	exTypes := []byte{'f', 'g', 'm', 's', 't'}
	var totolExp int64
	total := make(map[byte]int64, 0)
	for _, line := range input {
		for _, exType := range exTypes {
			// totolExp += getSumTypeFromRow(line, exType)
			total[exType] += getSumTypeFromRow(line, exType)
		}
	}
	for _, exType := range exTypes {
		totolExp += total[exType]
	}
	fmt.Printf("%.2f\n", float64(totolExp)/float64(len(input)))
	fmt.Printf("food %d\n", total['f'])
	fmt.Printf("game %d\n", total['g'])
	fmt.Printf("movie %d\n", total['m'])
	fmt.Printf("stationery %d\n", total['s'])
	fmt.Printf("transportation %d\n", total['t'])
	fmt.Printf("TOTAL %d\n", totolExp)
}

func getSumTypeFromRow(line string, exType byte) int64 {
	strArr := strings.Split(line, " ")
	var sum int64
	for _, str := range strArr {
		if str[0] == exType {
			nStr := strings.ReplaceAll(str, string(exType), "")
			n, err := strconv.ParseInt(nStr, 10, 64)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}
			sum += n
		}
	}
	return sum
}

func transpose(input []string) {
	for _, str := range input {
		fmt.Println(transposeLine(str))
	}
}

func transposeLine(str string) string {
	strArr := strings.Split(str, " ")
	res := ""
	for _, chord := range strArr {
		switch chord {
		case "C":
			res += "D"
		case "Dm":
			res += "Em"
		case "Em":
			res += "F#m"
		case "F":
			res += "G"
		case "G":
			res += "A"
		case "Am":
			res += "Bm"
		case "[C]":
			res += "[D]"
		case "[Dm]":
			res += "[Em]"
		case "[Em]":
			res += "[F#m]"
		case "[F]":
			res += "[G]"
		case "[G]":
			res += "[A]"
		case "[Am]":
			res += "[Bm]"
		default:
			res += chord
		}
		res += " "
	}
	return res
}

func addOrSum(input []string) {
	totalStr := "0"
	for _, str := range input {
		totalStr = addTwoString(totalStr, str)
	}
	fmt.Println(totalStr)
}

func addTwoString(text1, text2 string) string {
	l1 := ""
	l2 := ""

	if len(text1) >= len(text2) {
		l1 = text1
		l2 = text2
	} else {
		l1 = text2
		l2 = text1
	}

	if len(l1) > len(l2) {
		padZero := ""
		for i := 0; i < len(l1)-len(l2); i++ {
			padZero = padZero + "0"
		}
		l2 = padZero + l2
	}

	line1 := make([]int, len(l1))
	line2 := make([]int, len(l2))
	res := make([]int, len(l2))

	for i := range l1 {
		line1[i] = int(l1[i] - '0')
		line2[i] = int(l2[i] - '0')
	}

	carry := 0
	buffer := 0

	for i := len(line1) - 1; i >= 0; i-- {
		buffer = line1[i] + line2[i] + carry
		carry = 0
		if buffer >= 10 {
			carry = 1
			buffer -= 10
		}
		res[i] = buffer
	}

	strRes := ""

	if carry > 0 {
		strRes += fmt.Sprint(carry)
	}

	for _, i := range res {
		strRes += fmt.Sprint(i)
	}

	return strRes
}
