package traverse

// traverse is the front line user facing interface to this module. It sits
// on the top of the code stack and is allowed to use anything, but nothing
// else can depend on definitions here, except unit tests.

// sub package description:
//
// core:
// - requires: ["tbd"]
// - prohibited: [everything, except enums]
//
// cycle:
// - requires: ["tbd"]
// - prohibited: ["prefs"]
//
// enums:
// - requires: [nothing]
// - prohibited: [everything]
//
// hiber:
// - requires: ["refine"]
// - prohibited: ["tbd"]
//
// i18n:
//
// kernel:
//
// persist:
//
// pref:
//
// refine:
//
// resume:
// - requires: [""]
// - prohibited: ["tbd"]
//
// sampling:
// - requires: ["refine"]
// - prohibited: ["tbd"]
//
