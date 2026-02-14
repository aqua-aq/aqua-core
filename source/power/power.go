package power

type BindingPower byte

const (
	PowerLowest           BindingPower = iota
	PowerPatternAssigment              // `a, b = v`
	PowerDirectAssignment              // `a = b`
	PowerQuestion                      // `??`
	PowerOr                            // `or`
	PowerAnd                           // `and`
	PowerXor                           // `xor`
	PowerEquality                      // `==`, `~=`
	PowerComparison                    // `<`, `>`, `<=`, `>=`, `in`
	PowerShift                         // `<<`, `>>`
	PowerAdditive                      // `+`, `-`
	PowerMultiplicative                // `*`, `/`, `%` `//`
	PowerPrefix                        // `not`, `typeof`, `-`, `&`
	PowerBind                          // `->`
	PowerPostfix                       // `[]`, `++`, `--`, `.`, `.>`, `?.`, `?.>`
)
