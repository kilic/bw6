package bw6

const fpNumberOfLimbs = 12
const fpByteSize = 96
const fpBitSize = 761
const twelveWordBitSize = 768
const frBitSize = 377
const sixWordBitSize = 384

// Base field
// p = 0x122e824fb83ce0ad187c94004faff3eb926186a81d14688528275ef8087be41707ba638e584e91903cebaff25b423048689c8ed12f9fd9071dcd3dc73ebff2e98a116c25667a8f8160cf8aeeaf0a437e6913e6870000082f49d00000000008b
// r = 2 ^ 768

// -p^(-1) mod 2^64
var inp uint64 = 0xa5593568fa798dd

// supress linter warning: this variable used in assembly code
var _ = inp

// modulus = p
var modulus = fe{0xf49d00000000008b, 0xe6913e6870000082, 0x160cf8aeeaf0a437, 0x98a116c25667a8f8, 0x71dcd3dc73ebff2e, 0x8689c8ed12f9fd90, 0x03cebaff25b42304, 0x707ba638e584e919, 0x528275ef8087be41, 0xb926186a81d14688, 0xd187c94004faff3e, 0x0122e824fb83ce0a}

// r1 = r mod p
var r1 = &fe{0x0202ffffffff85d5, 0x5a5826358fff8ce7, 0x9e996e43827faade, 0xda6aff320ee47df4, 0xece9cb3e1d94b80b, 0xc0e667a25248240b, 0xa74da5bfdcad3905, 0x2352e7fe462f2103, 0x7b56588008b1c87c, 0x45848a63e711022f, 0xd7a81ebb9f65a9df, 0x0051f77ef127e87d}

// one = r1
var one = new(fe).set(r1)

// negativeOne = p - r
var negativeOne = &fe{0xf29a000000007ab6, 0x8c391832e000739b, 0x77738a6b6870f959, 0xbe36179047832b03, 0x84f3089e56574722, 0xc5a3614ac0b1d984, 0x5c81153f4906e9fe, 0x4d28be3a9f55c815, 0xd72c1d6f77d5f5c5, 0x73a18e069ac04458, 0xf9dfaa846595555f, 0x00d0f0a60a5be58c}

// r2 = r^2 mod p
var r2 = &fe{0xc686392d2d1fa659, 0x7b14c9b2f79484ab, 0x7fa1e825c1d2b459, 0xd6ec28f848329d88, 0x4afb427b73a1ed40, 0x972c69400d5930ae, 0x2c7a26bf8c995976, 0xac52e458c6e57af9, 0xac731bfa0c536dfe, 0x121e5c630b103f50, 0x8f1b0953b886cda4, 0x00ad253c2da8d807}

// nonResidue = -4
var nonResidue1 = &fe{0xe12e00000001e9c2, 0x63c1e3faa001cd69, 0xb1b4384fcbe29cf6, 0xc79630bc713d5a1d, 0x30127ac071851e2d, 0x0979f350dcd36af1, 0x6a66defed8b361f2, 0x53abac78b24d4e23, 0xb7ab89dede485a92, 0x5c3a0745675e8452, 0x446f17918c5f5700, 0x00fdf24e3267fa1e}

// nonResidue3 = u
var nonResidue3 = &fe3{fe{0}, *new(fe).set(one), fe{0}}

// pPlus1Over4 = (p + 1) / 4
var pPlus1Over4 = bigFromHex("0x48ba093ee0f382b461f250013ebfcfae49861aa07451a214a09d7be021ef905c1ee98e39613a4640f3aebfc96d08c121a2723b44be7f641c7734f71cfaffcba62845b09599ea3e05833e2bbabc290df9a44f9a1c000020bd27400000000023")

// pMinus1Over2 = (p - 1) / 2
var pMinus1Over2 = bigFromHex("0x9174127dc1e70568c3e4a0027d7f9f5c930c3540e8a34429413af7c043df20b83dd31c72c2748c81e75d7f92da11824344e476897cfec838ee69ee39f5ff974c508b612b33d47c0b067c577578521bf3489f34380000417a4e800000000045")

// parameter of p where p is actuall parameterized polynomial p(x)
var x = bigFromHex("0x8508c00000000001")
var xIsNeg = false
var ateLoop1 = bigFromHex("0x8508c00000000002")
var ateLoop1Neg = false
var ateLoop2 = computeNaf(bigFromHex("0x23ed1347970dec008a442f991fffffffffffffffffffffff"))
var ateLoop2Neg = false

/*
	Curve
	y^2 = x+3 + b
*/

// Group order
// q = x^6 - 2x^5 + 2x^3 + x + 1
var q = bigFromHex("0x1ae3a4617c510eac63b05c06ca1493b1a22d9f300f5138f1ef3622fba094800170b5d44300000008508c00000000001")

// b coefficient for G1
// b = -1
var b = new(fe).set(negativeOne)

// b2 coefficient for G2
// b2 = 4
var b2 = &fe{0x136efffffffe16c9, 0x82cf5a6dcffe3319, 0x6458c05f1f0e0741, 0xd10ae605e52a4eda, 0x41ca591c0266e100, 0x7d0fd59c3626929f, 0x9967dc004d00c112, 0x1ccff9c033379af5, 0x9ad6ec10a23f63af, 0x5cec11251a72c235, 0x8d18b1ae789ba83e, 0x0024f5d6c91bd3ec}

// G1 cofactor
var cofactorG1 = bigFromHex("0xad1972339049ce762c77d5ac34cb12efc856a0853c9db94cc61c554757551c0c832ba4061000003b3de580000000007c")

// G2 cofactor
var cofactorG2 = bigFromHex("0xad1972339049ce762c77d5ac34cb12efc856a0853c9db94cc61c554757551c0c832ba4061000003b3de5800000000075")

// G1 generator
var g1One = Point{
	fe{0xd6e42d7614c2d770, 0x4bb886eddbc3fc21, 0x64648b044098b4d2, 0x1a585c895a422985, 0xf1a9ac17cf8685c9, 0x352785830727aea5, 0xddf8cb12306266fe, 0x6913b4bfbc9e949a, 0x3a4b78d67ba5f6ab, 0x0f481c06a8d02a04, 0x91d4e7365c43edac, 0x00f4d17cd48beca5},
	fe{0x97e805c4bd16411f, 0x870d844e1ee6dd08, 0x1eba7a37cb9eab4d, 0xd544c4df10b9889a, 0x8fe37f21a33897be, 0xe9bf99a43a0885d2, 0xd7ee0c9e273de139, 0xaa6a9ec7a38dd791, 0x8f95d3fcf765da8e, 0x42326e7db7357c99, 0xe217e407e218695f, 0x009d1eb23b7cf684},
	*new(fe).set(one),
}

// G2 Generator
var g2One = Point{
	fe{0x3d902a84cd9f4f78, 0x864e451b8a9c05dd, 0xc2b3c0d6646c5673, 0x17a7682def1ecb9d, 0xbe31a1e0fb768fe3, 0x4df125e09b92d1a6, 0x0943fce635b02ee9, 0xffc8e7ad0605e780, 0x8165c00a39341e95, 0x8ccc2ae90a0f094f, 0x73a8b8cc0ad09e0c, 0x011027e203edd9f4},
	fe{0x9a159be4e773f67c, 0x6b957244aa8f4e6b, 0xa27b70c9c945a38c, 0xacb6a09fda11d0ab, 0x3abbdaa9bb6b1291, 0xdbdf642af5694c36, 0xb6360bb9560b369f, 0xac0bd1e822b8d6da, 0xfa355d17afe6945f, 0x8d6a0fc1fbcad35e, 0x72a63c7874409840, 0x0114976e5b0db280},
	*new(fe).set(one),
}

// G2 Twist type
var twistType = TWIST_TYPE_M

/*
	Frobenius Coefficients
*/

var frobeniuCoeffs31 = [3]fe{
	*new(fe).set(one),
	fe{0x7f96b51bd840c549, 0xd59782096496171f, 0x49b046fd9ce14bbc, 0x4b6163bba7527a56, 0xef6c92fb771d59f1, 0x0425bedbac1dfdc7, 0xd3ac39de759c0ffd, 0x9f43ed0e063a81d0, 0x5bd7d20b4f9a3ce2, 0x0411f03c36cf5c3c, 0x2d658fd49661c472, 0x01100249ae760b93},
	fe{0x67a04ae427bfb5f8, 0x9d32d491eb6a5cff, 0x43d03c1cb68051d4, 0x0b75ca96f69859a5, 0x0763497f5325ec60, 0x48076b5c278dd94d, 0x8ca3965ff91efd06, 0x1e6077657ea02f5d, 0xcdd6c153a8c37724, 0x28b5b634e5c22ea4, 0x9e01e3efd42e902c, 0x00e3d6815769a804},
}

var frobeniuCoeffs32 = [3]fe{
	*new(fe).set(one),
	fe{0x67a04ae427bfb5f8, 0x9d32d491eb6a5cff, 0x43d03c1cb68051d4, 0x0b75ca96f69859a5, 0x0763497f5325ec60, 0x48076b5c278dd94d, 0x8ca3965ff91efd06, 0x1e6077657ea02f5d, 0xcdd6c153a8c37724, 0x28b5b634e5c22ea4, 0x9e01e3efd42e902c, 0x00e3d6815769a804},
	fe{0x7f96b51bd840c549, 0xd59782096496171f, 0x49b046fd9ce14bbc, 0x4b6163bba7527a56, 0xef6c92fb771d59f1, 0x0425bedbac1dfdc7, 0xd3ac39de759c0ffd, 0x9f43ed0e063a81d0, 0x5bd7d20b4f9a3ce2, 0x0411f03c36cf5c3c, 0x2d658fd49661c472, 0x01100249ae760b93},
}

var frobeniusCoeffs6 = [6]fe{
	*new(fe).set(one),
	fe{0x8cfcb51bd8404a93, 0x495e69d68495a383, 0xd23cbc9234705263, 0x8d2b4c2b5fcf4f52, 0x6a798a5d20c612ce, 0x3e825d90eb6c2443, 0x772b249f2c9525fe, 0x521b2ed366e4b9bb, 0x84abb49bd7c4471d, 0x907062359c0f17e3, 0x3385e55030cc6f12, 0x3f11a3a41a2606},
	fe{0x7f96b51bd840c549, 0xd59782096496171f, 0x49b046fd9ce14bbc, 0x4b6163bba7527a56, 0xef6c92fb771d59f1, 0x0425bedbac1dfdc7, 0xd3ac39de759c0ffd, 0x9f43ed0e063a81d0, 0x5bd7d20b4f9a3ce2, 0x0411f03c36cf5c3c, 0x2d658fd49661c472, 0x01100249ae760b93},
	*new(fe).set(negativeOne),
	fe{0x67a04ae427bfb5f8, 0x9d32d491eb6a5cff, 0x43d03c1cb68051d4, 0x0b75ca96f69859a5, 0x0763497f5325ec60, 0x48076b5c278dd94d, 0x8ca3965ff91efd06, 0x1e6077657ea02f5d, 0xcdd6c153a8c37724, 0x28b5b634e5c22ea4, 0x9e01e3efd42e902c, 0x00e3d6815769a804},
	fe{0x75064ae427bf3b42, 0x10f9bc5f0b69e963, 0xcc5cb1b14e0f587b, 0x4d3fb306af152ea1, 0x827040e0fccea53d, 0x82640a1166dbffc8, 0x30228120b0181307, 0xd137b92adf4a6748, 0xf6aaa3e430ed815e, 0xb514282e4b01ea4b, 0xa422396b6e993acc, 0x0012e5db4d0dc277},
}

// x
// var x = bigFromHex("0x8508c00000000001")
