package hiber

// hiber represents the facility to be able to start navigation in hibernated state,
// ie we navigate but dont invoke a client action, until a certain condition occurs,
// specified by a node matching a filter. This is what used to be called listening
// in extendio. We could call these conditions, waypoints. We could wake or sleep
// type waypoints
//
// Hibernation depends on filtering.
//
