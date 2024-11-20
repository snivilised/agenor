package age

import (
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/pants"
)

// Composites enable the client to build CLIs without having to duplicate
// code. A composite is useful for any CLI that supports more than 1 scenario
// that would otherwise require code duplication.
//
// 🐍 `Hydra` (multi faceted/many headed) : supports all four scenarios
// 🐰 `Hare` (speedy): only supports run
// 🐢 `Tortoise` (slow): only supports walk
// 🐡 `Goldfish` (no memory/no resume) : only supports prime

type (
	// Scenario is a function that encodes within it the semantics of walk vs run
	// and prime vs resume. When invoked, it returns the underlying navigator
	// upon which a traversal session can be executed. The defined scenarios are
	// as follows:
	// * walk/prime
	// * walk/resume
	// * run/prime
	// * run/resume
	Scenario func(facade pref.Facade, settings ...pref.Option) Navigator
)

// Hydra is a composite that a client can use to build a cli that implements
// all four scenarios.
//
// The walk scenarios would ordinarily look something like this...
//
// # Walk/Prime
//
// Example:
//
// age.Walk().Configure().Extent(age.Prime(facade, options...)).Navigate(ctx)
//
// # Walk/Resume
//
// Example:
//
// age.Walk().Configure().Extent(age.Resume(facade, options...)).Navigate(ctx)
//
// The run scenarios would ordinarily look something like this...
//
// # Run/Prime
//
// Example:
//
// age.Run().Configure().Extent(age.Prime(facade, options...)).Navigate(ctx)
//
// # Run/Resume
//
// Example:
//
// age.Run().Configure().Extent(age.Resume(facade, options...)).Navigate(ctx)
//
// and these would need to be invoked conditionally depending on flags on the CLI.
// Notice the duplication; this can be resolved using the Hydra composite
// which can run all scenarios, so 2 values needs to be specified:
// isWalk and isPrime.
//
// # Hydra
//
// Example:
//
//	var wg sync.WaitGroup
//	isWalk := \<set depending on wether user requested walk or run\>
//	isPrime := \<set depending on wether user requested prime or resume\>
//	age.Hydra(isWalk, isPrime, &wg)(facade, options...).Navigate(ctx)
func Hydra(isWalk, isPrime bool, wg pants.WaitGroup) Scenario {
	if isWalk && isPrime {
		return slowPrime
	}

	if isWalk && !isPrime {
		return slowResume
	}

	if !isWalk && isPrime {
		return func(facade pref.Facade, settings ...pref.Option) Navigator {
			return fastPrime(facade, wg, settings...)
		}
	}

	return func(facade pref.Facade, settings ...pref.Option) Navigator {
		return fastResume(facade, wg, settings...)
	}
}

// Hare is a composite that a client can use to build a cli that only implements
// the run scenarios.
//
// The run scenarios would ordinarily look something like this...
//
// # Run/Prime
//
// Example:
//
// age.Run().Configure().Extent(age.Prime(facade, options...)).Navigate(ctx)
//
// # Run/Resume
//
// Example:
//
// age.Run().Configure().Extent(age.Resume(facade, options...)).Navigate(ctx)
//
// and these would need to be invoked conditionally depending on flags on the CLI.
// Notice the duplication; this can be resolved using the Hare composite
// which can only invoke Run sessions, so the query function is being asked to
// determine if it should prime or resume:
//
// # Hare
//
// Example:
//
//	var wg sync.WaitGroup
//	isPrime := \<set depending on wether user requested prime or resume\>
//	age.Hare(isPrime, &wg)(facade, options...).Navigate(ctx)
func Hare(isPrime bool, wg pants.WaitGroup) Scenario {
	if isPrime {
		return func(facade pref.Facade, settings ...pref.Option) Navigator {
			return fastPrime(facade, wg, settings...)
		}
	}

	return func(facade pref.Facade, settings ...pref.Option) Navigator {
		return fastResume(facade, wg, settings...)
	}
}

// Tortoise is a composite that a client can use to build a cli that only implements
// the walk scenarios.
//
// The walk scenarios would ordinarily look something like this...
//
// # Walk/Prime
//
// Example:
//
// age.Walk().Configure().Extent(age.Prime(facade, options...)).Navigate(ctx)
//
// # Walk/Resume
//
// Example:
//
// age.Walk().Configure().Extent(age.Resume(facade, options...)).Navigate(ctx)
//
// and these would need to be invoked conditionally depending on flags on the CLI.
// Notice the duplication; this can be resolved using the Tortoise composite
// which can only run Walk sessions, so the query function is being asked to
// determine if it should prime or resume:
//
// # Tortoise
//
// Example:
//
//	var wg sync.WaitGroup
//	isPrime := \<set depending on wether user requested prime or resume\>
//	age.Tortoise(isPrime)(facade, options...).Navigate(ctx)
func Tortoise(isPrime bool) Scenario {
	if isPrime {
		return slowPrime
	}

	return slowResume
}

// Goldfish is a composite that a client can use to build a cli that only implements
// the prime scenarios.
//
// The prime scenarios would ordinarily look something like this...
//
// # Walk/Prime
//
// Example:
//
// age.Walk().Configure().Extent(age.Prime(facade, options...)).Navigate(ctx)
//
// # Run/Prime
//
// Example:
//
// age.Run().Configure().Extent(age.Prime(facade, options...)).Navigate(ctx)
//
// and these would need to be invoked conditionally depending on flags on the CLI.
// Notice the duplication; this can be resolved using the Goldfish composite
// which can only run Prime sessions, so the query function is being asked to
// determine if it should walk or run:
//
// # Goldfish
//
// Example:
//
//	var wg sync.WaitGroup
//	isWalk := \<set depending on wether user requested walk or run\>
//	age.Goldfish(isWalk, &wg)(facade, options...).Navigate(ctx)
func Goldfish(isWalk bool, wg pants.WaitGroup) Scenario {
	if isWalk {
		return slowPrime
	}

	return func(facade pref.Facade, settings ...pref.Option) Navigator {
		return fastPrime(facade, wg, settings...)
	}
}

func slowPrime(facade pref.Facade, settings ...pref.Option) Navigator {
	return Walk().Configure().Extent(Prime(
		facade,
		settings...,
	))
}

func slowResume(facade pref.Facade, settings ...pref.Option) Navigator {
	return Walk().Configure().Extent(Resume(
		facade,
		settings...,
	))
}

func fastPrime(facade pref.Facade, wg pants.WaitGroup, settings ...pref.Option) Navigator {
	return Run(wg).Configure().Extent(Prime(
		facade,
		settings...,
	))
}

func fastResume(facade pref.Facade, wg pants.WaitGroup, settings ...pref.Option) Navigator {
	return Run(wg).Configure().Extent(Resume(
		facade,
		settings...,
	))
}
