package power

type BindingPower byte

const (
	PowerLowest         BindingPower = iota
	PowerAssignment                  // `=`
	PowerQuestion                    // `??`
	PowerOr                          // `or`
	PowerAnd                         // `and`
	PowerXor                         // `xor`
	PowerEquality                    // `==`, `~=`
	PowerComparison                  // `<`, `>`, `<=`, `>=`, `in`
	PowerShift                       // `<<`, `>>`
	PowerAdditive                    // `+`, `-`
	PowerMultiplicative              // `*`, `/`, `%` `//`
	PowerPrefix                      // prefix
	PowerBind                        // `->`
	PowerPostfix                     // `[]`, `++`, `--`, `.`, `.>`, `->`
)
