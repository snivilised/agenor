package collections_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ginkgo ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok

	"github.com/snivilised/agenor/collections"
	"github.com/snivilised/agenor/internal/third/lo"
)

const (
	ForwardIterator = true
	ReverseIterator = false
)

func reason(message string) string {
	return fmt.Sprintf("💥 failed because: '%v'", message)
}

type (
	record struct {
		name string
	}
	sleeve interface {
		song() string
	}

	iteratorTE struct {
		Given      string
		Should     string
		forward    bool
		sleeves    []sleeve
		recordPtrs []*record
		records    []record
		numbersN32 []int32
	}
	beginTE struct {
		iteratorTE
	}
)

func FormatIteratorTestDescription(entry *iteratorTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

func FormatBeginTestDescription(entry *beginTE) string {
	return fmt.Sprintf("Given: %v 🧪 should: %v", entry.Given, entry.Should)
}

func (e *record) song() string {
	return e.name
}

func getSleeveIt(forward bool, sequence []sleeve) collections.Iterator[sleeve] {
	var zero sleeve

	return lo.TernaryF(
		forward,
		func() collections.Iterator[sleeve] {
			return collections.ForwardIt(sequence, zero)
		},
		func() collections.Iterator[sleeve] {
			return collections.ReverseIt(sequence, zero)
		},
	)
}

func getSleeveRunIt(forward bool, sequence []sleeve) collections.RunnableIterator[sleeve, error] { // RunnableIterator
	var zero sleeve

	return lo.TernaryF(
		forward,
		func() collections.RunnableIterator[sleeve, error] {
			return collections.ForwardRunIt[sleeve, error](sequence, zero)
		},
		func() collections.RunnableIterator[sleeve, error] {
			return collections.ReverseRunIt[sleeve, error](sequence, zero)
		},
	)
}

func getRecordPtrIt(forward bool, sequence []*record) collections.Iterator[*record] {
	var zero *record

	return lo.TernaryF(
		forward,
		func() collections.Iterator[*record] {
			return collections.ForwardIt(sequence, zero)
		},
		func() collections.Iterator[*record] {
			return collections.ReverseIt(sequence, zero)
		},
	)
}

func getRecordsIt(forward bool, sequence []record) collections.Iterator[record] {
	zero := record{}

	return lo.TernaryF(
		forward,
		func() collections.Iterator[record] {
			return collections.ForwardIt(sequence, zero)
		},
		func() collections.Iterator[record] {
			return collections.ReverseIt(sequence, zero)
		},
	)
}

func getInt32It(forward bool, sequence []int32) collections.Iterator[int32] {
	var zero int32

	return lo.TernaryF(
		forward,
		func() collections.Iterator[int32] {
			return collections.ForwardIt(sequence, zero)
		},
		func() collections.Iterator[int32] {
			return collections.ReverseIt(sequence, zero)
		},
	)
}

var _ = Describe("Iterators", func() {
	Context("BeginX", func() {
		DescribeTable("interface",
			func(entry *beginTE) {
				it := getSleeveIt(entry.forward, entry.sleeves)

				Expect(it).NotTo(BeNil(), reason(entry.Should))
			},
			FormatBeginTestDescription,

			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:   "forward; empty sequence",
					Should:  "forward iterator",
					forward: ForwardIterator,
					sleeves: []sleeve{},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:   "forward; non empty sequence",
					Should:  "return forward iterator",
					forward: ForwardIterator,
					sleeves: []sleeve{&record{}},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:   "reverse; empty sequence",
					Should:  "return reverse iterator",
					sleeves: []sleeve{},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:   "reverse; non sequence",
					Should:  "return reverse iterator",
					sleeves: []sleeve{&record{}},
				},
			}),
		)

		DescribeTable("pointer to struct",
			func(entry *beginTE) {
				it := getRecordPtrIt(entry.forward, entry.recordPtrs)

				Expect(it).NotTo(BeNil(), reason(entry.Should))
			},
			FormatBeginTestDescription,

			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:      "forward; empty sequence",
					Should:     "return forward iterator",
					forward:    true,
					recordPtrs: []*record{},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:      "forward; non empty sequence",
					Should:     "return forward iterator",
					forward:    true,
					recordPtrs: []*record{{name: "norman *** rockwell"}},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:      "reverse; empty sequence",
					Should:     "return reverse iterator",
					recordPtrs: []*record{},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:      "reverse; non sequence",
					Should:     "return reverse iterator",
					recordPtrs: []*record{{name: "mariners apartment complex"}},
				},
			}),
		)

		DescribeTable("struct",
			func(entry *beginTE) {
				it := getRecordsIt(entry.forward, entry.records)

				Expect(it).NotTo(BeNil(), reason(entry.Should))
			},
			FormatBeginTestDescription,

			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:   "forward; empty sequence",
					Should:  "return forward iterator",
					forward: true,
					records: []record{},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:   "forward; non empty sequence",
					Should:  "return forward iterator",
					forward: true,
					records: []record{{name: "venice"}},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:   "reverse; empty sequence",
					Should:  "return nil reverse iterator",
					records: []record{},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:   "reverse; non sequence",
					Should:  "return reverse iterator",
					records: []record{{name: "*** it, i love you"}},
				},
			}),
		)

		DescribeTable("int32",
			func(entry *beginTE) {
				it := getInt32It(entry.forward, entry.numbersN32)

				Expect(it).NotTo(BeNil(), reason(entry.Should))
			},
			FormatBeginTestDescription,

			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:      "forward; empty sequence",
					Should:     "return forward iterator",
					forward:    true,
					numbersN32: []int32{},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:      "forward; non empty sequence",
					Should:     "return forward iterator",
					forward:    true,
					numbersN32: []int32{42},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:      "reverse; empty sequence",
					Should:     "return nil reverse iterator",
					numbersN32: []int32{},
				},
			}),
			Entry(nil, &beginTE{
				iteratorTE: iteratorTE{
					Given:      "reverse; non sequence",
					Should:     "return reverse iterator",
					numbersN32: []int32{42},
				},
			}),
		)
	})

	DescribeTable("Start",
		func(entry *beginTE) {
			it := lo.TernaryF(
				entry.forward,
				func() collections.Iterator[sleeve] {
					return collections.ForwardIt(entry.sleeves, nil)
				},
				func() collections.Iterator[sleeve] {
					return collections.ReverseIt(entry.sleeves, nil)
				},
			)

			Expect(it).NotTo(BeNil(), reason(entry.Should))
		},
		FormatBeginTestDescription,

		Entry(nil, &beginTE{
			iteratorTE: iteratorTE{
				Given:   "forward; empty sequence",
				Should:  "return forward iterator",
				forward: true,
				sleeves: []sleeve{},
			},
		}),
		Entry(nil, &beginTE{
			iteratorTE: iteratorTE{
				Given:   "reverse; empty sequence",
				Should:  "return nil reverse iterator",
				sleeves: []sleeve{},
			},
		}),
	)

	Context("Reset", func() {
		Context("forward", func() {
			When("Invoked", func() {
				It("🧪 should: re-assign content of iterator", func() {
					forwardIt := getSleeveIt(ForwardIterator, []sleeve{
						&record{name: "the next best american record"},
						&record{name: "the greatest"},
						&record{name: "bartender"},
					})
					for current := forwardIt.Start(); forwardIt.Valid(); current = forwardIt.Next() {
						_ = current
					}
					forwardIt.Reset([]sleeve{
						&record{name: "happiness is a butterfly"},
						&record{name: "hope is a dangerous thing ..."},
					})
					actual := forwardIt.Start().song()
					expected := "happiness is a butterfly"
					Expect(actual).To(Equal(expected))
				})
			})
		})

		Context("reverse", func() {
			When("Invoked", func() {
				It("🧪 should: re-assign content of iterator", func() {
					reverseIt := getSleeveIt(ReverseIterator, []sleeve{
						&record{name: "the next best american record"},
						&record{name: "the greatest"},
						&record{name: "bartender"},
					})
					for current := reverseIt.Start(); reverseIt.Valid(); current = reverseIt.Next() {
						_ = current
					}
					reverseIt.Reset([]sleeve{
						&record{name: "happiness is a butterfly"},
						&record{name: "hope is a dangerous thing ..."},
					})
					actual := reverseIt.Start().song()
					expected := "hope is a dangerous thing ..."
					Expect(actual).To(Equal(expected))
				})
			})
		})
	})

	Context("forward iterator", func() {
		When("empty sequence", func() {
			It("🧪 should: complete without error", func() {
				forwardIt := getSleeveIt(ForwardIterator, []sleeve{})

				for _ = forwardIt.Start(); forwardIt.Valid(); _ = forwardIt.Next() {
					Fail("!!! should not be invoked for empty collection")
				}
			})
		})
	})

	Context("reverse iterator", func() {
		When("empty sequence", func() {
			It("🧪 should: complete without error", func() {
				reverseIt := getSleeveIt(ReverseIterator, []sleeve{})

				for _ = reverseIt.Start(); reverseIt.Valid(); _ = reverseIt.Next() {
					Fail("!!! should not be invoked for empty collection")
				}
			})
		})
	})

	Context("single item sequence", func() {
		Context("forward iterator", func() {
			var forwardIt collections.Iterator[sleeve]
			BeforeEach(func() {
				forwardIt = getSleeveIt(ForwardIterator, []sleeve{&record{name: "love song"}})
			})

			Context("Start", func() {
				When("given: iterator in initial state", func() {
					It("🧪 should: return the single item", func() {
						item := forwardIt.Start()
						Expect(item.song()).To(Equal("love song"))
					})
				})
			})

			Context("Next", func() {
				When("given: iterator after Start", func() {
					It("🧪 should: return the zero value", func() {
						_ = forwardIt.Start()
						item := forwardIt.Next()
						Expect(item).To(BeNil())
					})
				})
			})
		})

		Context("reverse iterator", func() {
			var reverseIt collections.Iterator[sleeve]
			BeforeEach(func() {
				reverseIt = getSleeveIt(ReverseIterator, []sleeve{&record{name: "love song"}})
			})

			Context("Start", func() {
				When("given: iterator in initial state", func() {
					It("🧪 should: return the single item", func() {
						item := reverseIt.Start()
						Expect(item.song()).To(Equal("love song"))
					})
				})
			})

			Context("Valid", func() {
				When("given: iterator after Start", func() {
					It("🧪 should: return true", func() {
						_ = reverseIt.Start()
						Expect(reverseIt.Valid()).To(BeTrue())
					})
				})
			})

			Context("Next", func() {
				When("given: iterator after Start", func() {
					It("🧪 should: return the zero value", func() {
						_ = reverseIt.Start()
						item := reverseIt.Next()
						Expect(item).To(BeNil())
					})
				})
			})
		})
	})

	Context("multi item sequence", func() {
		Context("forward iterator", func() {
			var forwardIt collections.Iterator[sleeve]

			BeforeEach(func() {
				forwardIt = getSleeveIt(ForwardIterator, []sleeve{
					&record{name: "01 - cinnamon girl"},
					&record{name: "02 - how to disappear"},
					&record{name: "03 - california"},
				})
			})

			Context("Start", func() {
				When("given: iterator in initial state", func() {
					It("🧪 should: return the single item", func() {
						item := forwardIt.Start()
						Expect(item.song()).To(Equal("01 - cinnamon girl"))
					})
				})
			})

			Context("Valid", func() {
				When("given: iterator after Start", func() {
					It("🧪 should: return true", func() {
						_ = forwardIt.Start()
						Expect(forwardIt.Valid()).To(BeTrue())
					})
				})
			})

			Context("Next", func() {
				When("given: iterator after Start", func() {
					It("🧪 should: return the second item", func() {
						_ = forwardIt.Start()
						item := forwardIt.Next()
						Expect(item.song()).To(Equal("02 - how to disappear"))
					})
				})
			})

			Context("Next", func() {
				When("given: iterator in midway state", func() {
					It("🧪 should: return the second item", func() {
						_ = forwardIt.Start()
						_ = forwardIt.Next()
						item := forwardIt.Next()
						Expect(item.song()).To(Equal("03 - california"))
					})
				})
			})

			Context("full iteration", func() {
				It("🧪 should: iterate entire sequence (standard)", func() {
					actual := []string{}
					for current := forwardIt.Start(); forwardIt.Valid(); current = forwardIt.Next() {
						song := current.song()
						actual = append(actual, song)
						GinkgoWriter.Printf("===> 🔈🔈🔈 song: '%v'\n", song)
					}
					expected := []string{"01 - cinnamon girl", "02 - how to disappear", "03 - california"}
					Expect(actual).To(HaveExactElements(expected))
				})

				It("🧪 should: iterate entire sequence (do-while)", func() {
					actual := []string{}
					for current := forwardIt.Start(); ; {
						song := current.song()
						actual = append(actual, song)
						GinkgoWriter.Printf("===> 🔈🔈🔈 song: '%v'\n", song)

						current = forwardIt.Next()
						if !forwardIt.Valid() {
							break
						}
					}
					expected := []string{"01 - cinnamon girl", "02 - how to disappear", "03 - california"}
					Expect(actual).To(HaveExactElements(expected))
				})
			})
		})

		Context("reverse iterator", func() {
			var reverseIt collections.Iterator[sleeve]

			BeforeEach(func() {
				reverseIt = getSleeveIt(ReverseIterator, []sleeve{
					&record{name: "01 - cinnamon girl"},
					&record{name: "02 - how to disappear"},
					&record{name: "03 - california"},
				})
			})

			Context("Start", func() {
				When("given: iterator in initial state", func() {
					It("🧪 should: return the single item", func() {
						item := reverseIt.Start()
						Expect(item.song()).To(Equal("03 - california"))
					})
				})
			})

			Context("Valid", func() {
				When("given: iterator after Start", func() {
					It("🧪 should: return true", func() {
						_ = reverseIt.Start()
						Expect(reverseIt.Valid()).To(BeTrue())
					})
				})
			})

			Context("Next", func() {
				When("given: iterator after Start", func() {
					It("🧪 should: return the second item", func() {
						_ = reverseIt.Start()
						item := reverseIt.Next()
						Expect(item.song()).To(Equal("02 - how to disappear"))
					})
				})
			})

			Context("Next", func() {
				When("given: iterator in midway state", func() {
					It("🧪 should: return the second item", func() {
						_ = reverseIt.Start()
						_ = reverseIt.Next()
						item := reverseIt.Next()
						Expect(item.song()).To(Equal("01 - cinnamon girl"))
					})
				})
			})

			Context("full iteration", func() {
				It("🧪 should: iterate entire sequence (standard)", func() {
					actual := []string{}
					for current := reverseIt.Start(); reverseIt.Valid(); current = reverseIt.Next() {
						song := current.song()
						actual = append(actual, song)
						GinkgoWriter.Printf("===> 🔈🔈🔈 song: '%v'\n", song)
					}
					expected := []string{"03 - california", "02 - how to disappear", "01 - cinnamon girl"}
					Expect(actual).To(HaveExactElements(expected))
				})

				It("🧪 should: iterate entire sequence (do-while)", func() {
					actual := []string{}
					for current := reverseIt.Start(); ; {
						song := current.song()
						actual = append(actual, song)
						GinkgoWriter.Printf("===> 🔈🔈🔈 song: '%v'\n", song)

						current = reverseIt.Next()
						if !reverseIt.Valid() {
							break
						}
					}
					expected := []string{"03 - california", "02 - how to disappear", "01 - cinnamon girl"}
					Expect(actual).To(HaveExactElements(expected))
				})
			})

		})

		Context("runnable", Ordered, func() {
			var sleeves []sleeve

			BeforeAll(func() {
				sleeves = []sleeve{
					&record{name: "07 - cinnamon girl"},
					&record{name: "08 - how to disappear"},
					&record{name: "09 - california"},
					&record{name: "BONUS - 01"},
					&record{name: "BONUS - 02"},
				}
			})

			Context("forward", func() {
				When("while condition is never invalidated", func() {
					It("🧪 should: invoke each for all items in sequence", func() {
						const (
							expected = 5
							forward  = true
						)

						iterator := getSleeveRunIt(forward, sleeves)
						actual := 0
						each := func(_ sleeve) error {
							actual++

							return nil
						}
						while := func(_ sleeve, _ error) bool {
							return true
						}

						iterator.RunAll(each, while)
						Expect(actual).To(Equal(expected))
					})
				})

				When("while condition is invalidated before end of sequence", func() {
					It("🧪 should: invoke each for item until while fails", func() {
						const (
							expected = 4
							forward  = true
						)

						iterator := getSleeveRunIt(forward, sleeves)
						actual := 0
						each := func(_ sleeve) error {
							actual++

							return nil
						}
						while := func(s sleeve, _ error) bool {
							return strings.HasPrefix(s.song(), "0")
						}

						iterator.RunAll(each, while)
						Expect(actual).To(Equal(expected))
					})
				})
			})

			Context("reverse", Ordered, func() {
				When("while condition is never invalidated", func() {
					It("🧪 should: invoke each for all items in sequence", func() {
						const (
							expected = 5
							forward  = false
						)

						iterator := getSleeveRunIt(forward, sleeves)
						actual := 0
						each := func(_ sleeve) error {
							actual++

							return nil
						}
						while := func(_ sleeve, _ error) bool {
							return true
						}

						iterator.RunAll(each, while)
						Expect(actual).To(Equal(expected))
					})
				})

				When("while condition is invalidated before end of sequence", func() {
					It("🧪 should: invoke each for item until while fails", func() {
						const (
							expected = 3
							forward  = false
						)

						iterator := getSleeveRunIt(forward, sleeves)
						actual := 0
						each := func(_ sleeve) error {
							actual++

							return nil
						}
						while := func(s sleeve, _ error) bool {
							return strings.HasPrefix(s.song(), "BONUS")
						}

						iterator.RunAll(each, while)
						Expect(actual).To(Equal(expected))
					})
				})
			})
		})
	})
})
