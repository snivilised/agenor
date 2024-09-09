package collections_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/traverse/collections"
	lab "github.com/snivilised/traverse/internal/laboratory"
)

var (
	rainbow   = []string{"richard", "of", "york", "gave", "battle", "in", "vain"}
	scrambled = []string{"york", "vain", "battle", "of", "richard", "gave", "in"}
)

func assertColoursAreInOrder(set *collections.PositionalSet[string]) {
	anchor, _ := set.Position("ANCHOR")

	for _, colour := range rainbow {
		pos, _ := set.Position(colour)
		Expect(pos < anchor).To(BeTrue(), lab.Reason(
			fmt.Sprintf("position(%v) of colour: %v should be less than anchor's(%v)",
				pos, colour, anchor),
		))
	}
}

var _ = Describe("PositionalSet", func() {
	type (
		orderingStringTE struct {
			given  string
			should string
			roles  []string
		}
	)

	var (
		set *collections.PositionalSet[string]
	)

	BeforeEach(func() {
		set = collections.NewPositionalSet(rainbow, "ANCHOR")
	})

	Context("Count", func() {
		When("no items added", func() {
			It("🧪 should: contain just the anchor", func() {
				Expect(set.Count()).To(Equal(1), lab.Reason("only anchor should be present"))
			})
		})
	})

	Context("Insert", func() {
		When("requested item is the anchor", func() {
			It("🧪 should: not insert", func() {
				Expect(set.Insert("ANCHOR")).To(BeFalse(), lab.Reason("inserting anchor is invalid"))
				Expect(set.Count()).To(Equal(1), lab.Reason("only anchor should be present"))
			})
		})

		When("valid item requested", func() {
			It("🧪 should: insert", func() {
				Expect(set.Insert("richard")).To(BeTrue(), lab.Reason("richard is in order list"))
				Expect(set.Count()).To(Equal(2), lab.Reason("richard, anchor"))
			})
		})

		When("valid item already present", func() {
			It("🧪 should: not insert", func() {
				set.Insert("richard")
				Expect(set.Insert("richard")).To(BeFalse(), lab.Reason("richard already in order list"))
				Expect(set.Count()).To(Equal(2), lab.Reason("richard, anchor"))
			})
		})

		When("invalid item requested", func() {
			It("🧪 should: not insert", func() {
				Expect(set.Insert("gold")).To(BeFalse(), lab.Reason("gold not in order list"))
				Expect(set.Count()).To(Equal(1), lab.Reason("only anchor should be present"))
			})
		})
	})

	Context("All", func() {
		When("All valid items requested", func() {
			It("🧪 should: insert all", func() {
				Expect(set.All(
					"richard", "of", "york", "gave", "battle", "in", "vain",
				)).To(BeTrue(), lab.Reason("all items are valid"))
				Expect(set.Count()).To(Equal(8), lab.Reason("should contain all items"))
			})
		})

		When("Not all are valid", func() {
			It("🧪 should: insert only valid", func() {
				Expect(set.All(
					"richard", "gold", "of", "silver", "york", "bronze",
				)).To(BeFalse(), lab.Reason("all items are valid"))
				Expect(set.Count()).To(Equal(4), lab.Reason("should contain valid items"))
			})
		})
	})

	Context("Delete", func() {
		When("requested item is the anchor", func() {
			It("🧪 should: not delete", func() {
				set.Delete("ANCHOR")
				Expect(set.Count()).To(Equal(1), lab.Reason("anchor should still be present"))
			})
		})

		When("requested valid item is present", func() {
			It("🧪 should: delete", func() {
				set.Insert("york")
				set.Delete("york")
				Expect(set.Count()).To(Equal(1), lab.Reason("york should deleted"))
			})
		})

		When("requested valid item is not present", func() {
			It("🧪 should: not delete", func() {
				set.Delete("york")
				Expect(set.Count()).To(Equal(1), lab.Reason("only anchor should be present"))
			})
		})

		When("requested valid item is not valid", func() {
			It("🧪 should: not delete", func() {
				set.Delete("silver")
				Expect(set.Count()).To(Equal(1), lab.Reason("only anchor should be present"))
			})
		})
	})

	Context("Position", func() {
		When("multiple items inserted in order", func() {
			It("🧪 should: return position less than anchor", func() {
				set.All(rainbow...)
				assertColoursAreInOrder(set)
			})
		})

		When("multiple items inserted out of order", func() {
			It("🧪 should: return position less than anchor", func() {
				set.All(scrambled...)
				assertColoursAreInOrder(set)
			})

			It("🧪 should: contain correct positions when colours compared", func() {
				set.All(scrambled...)
				richard, _ := set.Position("richard")
				of, _ := set.Position("of")
				york, _ := set.Position("york")
				Expect(richard < of).To(BeTrue())
				Expect(richard < york).To(BeTrue())
			})
		})
	})

	Context("Items", func() {
		When("multiple items inserted in order", func() {
			It("🧪 should: return items defined by order", func() {
				set.All(rainbow...)
				expected := append(append([]string{}, rainbow...), "ANCHOR")
				Expect(set.Items()).To(HaveExactElements(expected))
			})
		})

		When("multiple items inserted out of order", func() {
			It("🧪 should: return items defined by order", func() {
				set.All(scrambled...)
				expected := append(append([]string{}, rainbow...), "ANCHOR")
				Expect(set.Items()).To(HaveExactElements(expected))
			})
		})

		When("partial items inserted in order", func() {
			It("🧪 should: return items defined by order", func() {
				set.All("vain", "battle", "york")
				expected := []string{"york", "battle", "vain", "ANCHOR"}
				Expect(set.Items()).To(HaveExactElements(expected))

				set.Delete("battle")
				set.Insert("of")

				expected = []string{"of", "york", "vain", "ANCHOR"}
				Expect(set.Items()).To(HaveExactElements(expected))
			})
		})
	})

	Context("Contains", func() {
		When("item is present", func() {
			It("🧪 should: return true", func() {
				set.All(rainbow...)
				Expect(set.Contains("richard")).To(BeTrue(),
					lab.Reason("richard IS contained in set"),
				)
			})
		})

		When("item is NOT present", func() {
			It("🧪 should: return false", func() {
				set.All(rainbow...)
				Expect(set.Contains("james")).To(BeFalse(),
					lab.Reason("james is NOT contained in set"),
				)
			})
		})
	})
})
