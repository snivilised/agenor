package filtering

// 📦 pkg: filtering - this package is required because filters are required
// not just but the filter plugin, but others too like hibernation. The filter
// required by hibernation could have been implemented by the filter plugin,
// but doing so in this fashion would have mean introducing coupling of
// hibernation on filter; ie how to allow hibernation to access the filter(s)
// created by filter?
//	Instead, we factor out the filter creation code to this package, so that
// hibernation can create and apply filters as it needs, without depending on
// filter. So filter, now doesn't own the filter implementations, rather it's
// simply responsible for the plugin aspects of filtering, not implementation
// or creation.
//
