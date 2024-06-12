package assert

// Assert will panic if a condition is false and report the given state.
//
// An Assert creates a contract for invariants which can be used to catch bugs
// at runtime. Asserts can be recursively deleted prior to compiling for release.
func Assert(condition bool, state any) {
    if !condition {
        panic(state)
    }
}
