package traverse

// traverse is the front line user facing interface to this module. It sits
// on the top of the code stack and is allowed to use anything, but nothing
// else can depend on definitions here, except unit tests.

// sub package description:
//

// This high level list assumes everything can use core and enums; dependencies
// can only point downwards.
// ============================================================================
// 🔆 user interface layer
// traverse: [everything]
// ---
//
// 🔆 feature layer
// resume: ["pref"]
// sampling: ["refine"]
// hiber: ["refine", "services"]
// refine: []
//
// 🔆 central layer
// kernel: []
// ---
//
// 🔆 support layer
// pref: ["cycle", "services", "persist(to-be-confirmed)"] actually, persist should be part of pref
// persist: []
// services: []
// ---
//
// 🔆 intermediary layer
// cycle: [], !("pref")
// ---
//
// 🔆 platform layer
// core: []
// enums: [none]
// ---
// ============================================================================
//
