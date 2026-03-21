// Package agenor is the front line user facing interface to this module.
// It sits on the top of the code stack and is allowed to use anything, but
// nothing else can depend on definitions here, except unit tests.
package agenor

// sub package description:
//

// This high level list assumes everything can use core and enums; dependencies
// can only point downwards. NB: These restrictions do not apply to the unit tests;
// eg, "life_test" defines tests that are dependent on "pref", but "life" is prohibited
// from using "pref".
// ============================================================================
// 🔆 user interface layer
// agenor: [everything]
// ---
//
// 🔆 feature layer
// resume: ["pref", "opts", "kernel"]
// sampling: ["filter"]
// hiber: ["filter", "services"]
// filter: []
//
// 🔆 central layer
// kernel: []
// enclave: [pref, override]
// opts: [pref]
// override: [tapable], !("enclave")
// ---
//
// 🔆 support layer
// pref: ["life", "services", "persist(to-be-confirmed)"] actually, persist should be part of pref
// persist: []
// services: []
// ---
//
// 🔆 intermediary layer
// life: [], !("pref")
// ---
//
// 🔆 platform layer
// tapable: [core]
// core: []
// enums: [none]
// tfs:
// ---
// ============================================================================
//
