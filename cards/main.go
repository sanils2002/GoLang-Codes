package main

func main() {
	cards := newDeck()
	cards.shuffle()
	cards.print()
	// cards := newDeckFromFile("my_cards")
	// cards.print()
	// cards.saveToFile("my_cards")

	// hand, remainingCards := deal(cards, 5)

	// hand.print()
	// remainingCards.print()

	// fmt.Println(cards.toString())
}
