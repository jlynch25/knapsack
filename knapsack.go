package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type object struct {
	value   int
	weight  int
	include int
}

var objects []object = []object{}

type KnapSack struct {
	value   int
	weight  int
	objects []object
}

var KnapSacks []KnapSack = []KnapSack{}

type ByValue []KnapSack

func (a KnapSack) getValue() int {
	value := 0
	for _, o := range a.objects {
		value += o.value * o.include
	}
	return value
}

func (a KnapSack) getWeight() int {
	weight := 0
	for _, o := range a.objects {
		weight += o.weight * o.include
	}
	return weight
}

func main() {

	objects = []object{{78, 18, 0}, {35, 9, 0}, {89, 23, 0}, {36, 20, 0}, {94, 59, 0}, {75, 61, 0}, {74, 70, 0}, {79, 75, 0}, {80, 76, 0}, {16, 30, 0}}
	fmt.Println("knapSack size: 103")
	fmt.Printf("\nknapSack value: %v \n\n", knapSack(103))
	fmt.Println("knapSack size: 156")
	fmt.Printf("\nknapSack value: %v", knapSack(156))
}

func knapSack(cap int) KnapSack {

	totalFitness := 0.0
	var averageFitness []float64
	var averageFitnessIndex []float64
	var bestFitness []float64
	usedPopulationSize := 0
	maxValue := 425 //659 // ONLY USED FOR GRAPHING

	population := len(objects) * 10

	rand.Seed(time.Now().UnixNano())

	for i := 1; i < population; i++ {
		sack := RandSeq(objects)
		var knapsack KnapSack
		knapsack.objects = sack
		knapsack.value = knapsack.getValue()
		knapsack.weight = knapsack.getWeight()

		if knapsack.weight > cap {
			knapsack.value = 0
		}

		KnapSacks = append(KnapSacks, knapsack)
	}

	sort.Sort(sort.Reverse(ByValue(KnapSacks)))

	for j := 1; j <= 7; j++ {

		fmt.Print("░░")

		totalFitness = 0
		for _, sack := range KnapSacks {
			totalFitness = totalFitness + float64(sack.value)
		}
		average := totalFitness / float64(len(KnapSacks)*maxValue)
		averageFitness = append(averageFitness, average)
		averageFitnessIndex = append(averageFitnessIndex, float64(j))
		bestFitness = append(bestFitness, float64(KnapSacks[0].value))

		usedPopulationSize = int(float64(len(KnapSacks)) * 0.4)

		for range KnapSacks {
			rand := rand.Intn(len(KnapSacks))
			newSack := KnapSacks[rand].mutation(cap)
			KnapSacks = append(KnapSacks, newSack)
		}

		for range KnapSacks {
			rand := rand.Intn(len(KnapSacks))
			newSack := KnapSacks[rand].swap(cap)
			KnapSacks = append(KnapSacks, newSack)
		}

		for range KnapSacks {
			j := rand.Intn(len(KnapSacks))
			k := rand.Intn(len(KnapSacks))
			newSack1, newSack2 := crossover(KnapSacks[j], KnapSacks[k], cap)
			KnapSacks = append(KnapSacks, newSack1, newSack2)
		}

		sort.Sort(sort.Reverse(ByValue(KnapSacks)))

		for k := 1; k < (len(KnapSacks) - usedPopulationSize); k++ {
			KnapSacks = KnapSacks[:len(KnapSacks)-1] // Truncate slice.
		}

	}

	plotGraph(averageFitness, averageFitnessIndex, bestFitness)

	return KnapSacks[0] // the best value knapsack within weight
}

func (a KnapSack) mutation(cap int) KnapSack {
	i := rand.Intn(len(a.objects))

	flipCoin := rand.Intn(1)

	if a.objects[i].include > 0 {
		if flipCoin == 0 {
			if (a.objects[i].weight * a.objects[i].include) < cap {
				a.objects[i].include = a.objects[i].include + 1
			} else {
				a.objects[i].include = 0
			}
		} else {
			a.objects[i].include = a.objects[i].include - 1
		}
	} else {
		a.objects[i].include = a.objects[i].include + 1
	}

	a.value = a.getValue()
	a.weight = a.getWeight()

	if a.weight > cap {
		a.value = 0
		// KnapSacks = append(KnapSacks, a.swap(cap))
	}
	return a
}

func (a KnapSack) swap(cap int) KnapSack {
	i := rand.Intn(len(a.objects))
	j := rand.Intn(len(a.objects))

	temp1 := a.objects[i].include
	temp2 := a.objects[j].include

	a.objects[i].include = temp2
	a.objects[j].include = temp1

	a.value = a.getValue()
	a.weight = a.getWeight()

	if a.weight > cap {
		a.value = 0
		// a.mutation(cap)
	}
	return a
}

func crossover(a, b KnapSack, cap int) (KnapSack, KnapSack) {

	rand := rand.Intn(len(a.objects))

	temp1 := a.objects[0:rand]
	temp2 := b.objects[0:rand]

	a.objects = append(temp2, a.objects[rand:]...)
	b.objects = append(temp1, b.objects[rand:]...)

	a.value = a.getValue()
	a.weight = a.getWeight()
	if a.weight > cap {
		a.value = 0
		// KnapSacks = append(KnapSacks, a.swap(cap))
	}
	b.value = b.getValue()
	b.weight = b.getWeight()

	if b.weight > cap {
		b.value = 0
		// KnapSacks = append(KnapSacks, b.swap(cap))
	}

	return a, b
}

// Use to generate random binary strings
func RandSeq(objects []object) []object {
	var includes = []rune("0123")
	var KnapSack []object = objects
	b := make([]rune, len(objects))
	for i := range b {
		KnapSack[i].include = int(includes[rand.Intn(len(includes))])
	}
	return KnapSack
}

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Less(i, j int) bool { return a[i].value < a[j].value }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func plotGraph(results, resultIndex, bestResults []float64) {

	p := plot.New()

	p.Title.Text = "Average Fitness per Generation"
	p.Y.Label.Text = "Average Fitness"
	p.X.Label.Text = "Generations"

	pts := make(plotter.XYs, len(results))

	for i := range results {

		pts[i].Y = results[i]

		pts[i].X = resultIndex[i]
	}

	pts2 := make(plotter.XYs, len(bestResults))

	for i := range results {

		pts2[i].Y = bestResults[i]

		pts2[i].X = resultIndex[i]
	}

	plotutil.AddLinePoints(p, "Average", pts, "Best", pts2)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "Knapsack_156.png"); err != nil {
		panic(err)
	}
}
