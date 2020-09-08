package bw6

const N_LIMBS = 12
const FE_BYTE_SIZE = 96
const FE_BIT_SIZE = 761
const TWELWE_WORD_BYTE_SIZE = 768

/*
	Base field
	p = 0x122e824fb83ce0ad187c94004faff3eb926186a81d14688528275ef8087be41707ba638e584e91903cebaff25b423048689c8ed12f9fd9071dcd3dc73ebff2e98a116c25667a8f8160cf8aeeaf0a437e6913e6870000082f49d00000000008b
	r = 2^768
*/

// -p^(-1) mod 2^64
var inp uint64 = 744663313386281181

// modulus = p
var modulus = fe{
	0xf49d00000000008b,
	0xe6913e6870000082,
	0x160cf8aeeaf0a437,
	0x98a116c25667a8f8,
	0x71dcd3dc73ebff2e,
	0x8689c8ed12f9fd90,
	0x03cebaff25b42304,
	0x707ba638e584e919,
	0x528275ef8087be41,
	0xb926186a81d14688,
	0xd187c94004faff3e,
	0x0122e824fb83ce0a,
}

// zero = 0
var zero = &fe{
	0x0000000000000000,
	0x0000000000000000,
	0x0000000000000000,
	0x0000000000000000,
	0x0000000000000000,
	0x0000000000000000,
	0x0000000000000000,
	0x0000000000000000,
	0x0000000000000000,
	0x0000000000000000,
	0x0000000000000000,
	0x0000000000000000,
}

// r1 = r mod p
var r1 = &fe{
	0x0202ffffffff85d5,
	0x5a5826358fff8ce7,
	0x9e996e43827faade,
	0xda6aff320ee47df4,
	0xece9cb3e1d94b80b,
	0xc0e667a25248240b,
	0xa74da5bfdcad3905,
	0x2352e7fe462f2103,
	0x7b56588008b1c87c,
	0x45848a63e711022f,
	0xd7a81ebb9f65a9df,
	0x0051f77ef127e87d,
}

// one = r1
var one = new(fe).set(r1)

// negativeOne = p - r
var negativeOne = &fe{
	0xf29a000000007ab6,
	0x8c391832e000739b,
	0x77738a6b6870f959,
	0xbe36179047832b03,
	0x84f3089e56574722,
	0xc5a3614ac0b1d984,
	0x5c81153f4906e9fe,
	0x4d28be3a9f55c815,
	0xd72c1d6f77d5f5c5,
	0x73a18e069ac04458,
	0xf9dfaa846595555f,
	0x00d0f0a60a5be58c,
}

// r2 = r^2 mod p
var r2 = &fe{
	0xc686392d2d1fa659,
	0x7b14c9b2f79484ab,
	0x7fa1e825c1d2b459,
	0xd6ec28f848329d88,
	0x4afb427b73a1ed40,
	0x972c69400d5930ae,
	0x2c7a26bf8c995976,
	0xac52e458c6e57af9,
	0xac731bfa0c536dfe,
	0x121e5c630b103f50,
	0x8f1b0953b886cda4,
	0x00ad253c2da8d807,
}

// nonResidue = -4
var nonResidue = &fe{
	0xe12e00000001e9c2,
	0x63c1e3faa001cd69,
	0xb1b4384fcbe29cf6,
	0xc79630bc713d5a1d,
	0x30127ac071851e2d,
	0x0979f350dcd36af1,
	0x6a66defed8b361f2,
	0x53abac78b24d4e23,
	0xb7ab89dede485a92,
	0x5c3a0745675e8452,
	0x446f17918c5f5700,
	0x00fdf24e3267fa1e,
}

// pPlus1Over4 = (p + 1) / 4
var pPlus1Over4 = bigFromHex("0x48ba093ee0f382b461f250013ebfcfae49861aa07451a214a09d7be021ef905c1ee98e39613a4640f3aebfc96d08c121a2723b44be7f641c7734f71cfaffcba62845b09599ea3e05833e2bbabc290df9a44f9a1c000020bd27400000000023")

// pMinus1Over2 = (p - 1) / 2
var pMinus1Over2 = bigFromHex("0x9174127dc1e70568c3e4a0027d7f9f5c930c3540e8a34429413af7c043df20b83dd31c72c2748c81e75d7f92da11824344e476897cfec838ee69ee39f5ff974c508b612b33d47c0b067c577578521bf3489f34380000417a4e800000000045")

// parameter of p
var x = bigFromHex("0x8508c00000000001")

/*
	Curve
	y^2 = x+3 + b
*/

// Group order
// q = x^6 - 2x^5 + 2x^3 + x + 1
var q = bigFromHex("0x1ae3a4617c510eac63b05c06ca1493b1a22d9f300f5138f1ef3622fba094800170b5d44300000008508c00000000001")

// b coefficient for G1
// b = -1
var b = &fe{0xaa270000000cfff3, 0x53cc0032fc34000a, 0x478fe97a6b0a807f, 0xb1d37ebee6ba24d7, 0x8ec9733bbf78ab2f, 0x09d645513d83de7e}

// G1 cofactor
var cofactorG1 = bigFromHex("0xad1972339049ce762c77d5ac34cb12efc856a0853c9db94cc61c554757551c0c832ba4061000003b3de580000000007c")

// G1 generator
// x = 0x01075b020ea190c8b277ce98a477beaee6a0cfb7551b27f0ee05c54b85f56fc779017ffac15520ac11dbfcd294c2e746a17a54ce47729b905bd71fa0c9ea097103758f9a280ca27f6750dd0356133e82055928aca6af603f4088f3af66e5b43d
// y = 0x0058b84e0a6fc574e6fd637b45cc2a420f952589884c9ec61a7348d2a2e573a3265909f1af7e0dbac5b8fa1771b5b806cc685d31717a4c55be3fb90b6fc2cdd49f9df141b3053253b2b08119cad0fb93ad1cb2be0b20d2a1bafc8f2db4e95363
var g1One = PointG1{
	fe{0xd6e42d7614c2d770,
		0x4bb886eddbc3fc21,
		0x64648b044098b4d2,
		0x1a585c895a422985,
		0xf1a9ac17cf8685c9,
		0x352785830727aea5,
		0xddf8cb12306266fe,
		0x6913b4bfbc9e949a,
		0x3a4b78d67ba5f6ab,
		0x0f481c06a8d02a04,
		0x91d4e7365c43edac,
		0x00f4d17cd48beca5},
	fe{0x97e805c4bd16411f,
		0x870d844e1ee6dd08,
		0x1eba7a37cb9eab4d,
		0xd544c4df10b9889a,
		0x8fe37f21a33897be,
		0xe9bf99a43a0885d2,
		0xd7ee0c9e273de139,
		0xaa6a9ec7a38dd791,
		0x8f95d3fcf765da8e,
		0x42326e7db7357c99,
		0xe217e407e218695f,
		0x009d1eb23b7cf684},
	*new(fe).set(r1),
}
