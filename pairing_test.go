package bw6

import (
	"encoding/hex"
	"math/big"
	"testing"
)

func fromHex(size int, hexStrs ...string) []byte {
	var out []byte
	if size > 0 {
		out = make([]byte, size*len(hexStrs))
	}
	for i := 0; i < len(hexStrs); i++ {
		hexStr := hexStrs[i]
		if hexStr[:2] == "0x" {
			hexStr = hexStr[2:]
		}
		if len(hexStr)%2 == 1 {
			hexStr = "0" + hexStr
		}
		bytes, err := hex.DecodeString(hexStr)
		if err != nil {
			return nil
		}
		if size <= 0 {
			out = append(out, bytes...)
		} else {
			if len(bytes) > size {
				return nil
			}
			offset := i*size + (size - len(bytes))
			copy(out[offset:], bytes)
		}
	}
	return out
}

func TestPairingMillerLoop(t *testing.T) {
	bw6 := NewEngine()
	// GT := bw6.GT()
	g1, g2 := bw6.G1, bw6.G2
	// e(a*G1, b*G2) = e(G1, G2)^c
	expected := new(fe6)
	tmp := new(fe)
	tmp, _ = fromString("0007F1343BD9C1E8952750396AD193649ECA67CB0AE37721DE07E52F02C8605929D10FC6004C720356FA9C5C9FFA1CB35166D93AD9AE073CAA5657904C9FAAD9493A9FF77330A5DCB2FE7B2A9FB4FCB426985C4C283C5E5854620F23C2D11617")
	expected[0][0].set(tmp)
	tmp, _ = fromString("00A4678EA437044B4F8BE20EEC2F108234087B76D7D3F3A8FDB9F75158EA39CABCF7E1BFA20AA188117838C4E3D1AA978315EEB7D4D07BCEB8BAADEA596CBA9E7410CD3D194E8B6AE7982D5C9B2D7101E1191E1B489C1A418361D9F0C17F3A38")
	expected[0][1].set(tmp)
	tmp, _ = fromString("000BE5A7D14480EFD9CF038C51FE798BF752B0FD5EC00AEE3678D923837D753DA90ADD45154A7BAA28BC237496A9870DF7239CEF71526C73EA066E9AF34E74EBE03B8398E8F46C0A6DAB7842598EF0DCD5EF4BC019F8D20639838F1A8ED65012")
	expected[0][2].set(tmp)
	tmp, _ = fromString("0036AB139835FC9D146E519F99A11AAF557E0DF8A971DCB5379A66466C066DB4783716D9C17657F00AD4E68DA91FCBD254F21092F35A4ECE4016B29AE863C19008DE8105503B0BF57AD89B4D707EC51E0EB28D8C49B0AEC5D11C6E7D45713592")
	expected[1][0].set(tmp)
	tmp, _ = fromString("00E98C7B72100B181876B1D159643543C803E7B84C8F15BF7DE53062DCF4B922BF7F464BB2A1F6B5891F01143E3930DC8561AEF6D09EDA8DDB03673603C20F5E3FAB9F2E7911792FF789ED9CFEF6BFC45A7AC86B7E8010E7626ED6D2D642AED8")
	expected[1][1].set(tmp)
	tmp, _ = fromString("00013D59C89577690D7083AC6BE49C052DDE8CEC2C7A776076644227AC84485372138E60ECB2835069C8059F633C0BA9967F76397861D3A3227D78EF65104E675FC65EC268DF1681F3D0601EB2EB689B644BCAAF5912E1936B70CC017E9F25ED")
	expected[1][2].set(tmp)

	if !bw6.GT().IsValid(expected) {
		t.Fatalf("expected is not in correct subgroup")
	}

	// expected, err := GT.FromBytes(
	// 	fromHex(
	// 		FE_BYTE_SIZE,
	// 		"0007F1343BD9C1E8952750396AD193649ECA67CB0AE37721DE07E52F02C8605929D10FC6004C720356FA9C5C9FFA1CB35166D93AD9AE073CAA5657904C9FAAD9493A9FF77330A5DCB2FE7B2A9FB4FCB426985C4C283C5E5854620F23C2D11617",
	// 		"00A4678EA437044B4F8BE20EEC2F108234087B76D7D3F3A8FDB9F75158EA39CABCF7E1BFA20AA188117838C4E3D1AA978315EEB7D4D07BCEB8BAADEA596CBA9E7410CD3D194E8B6AE7982D5C9B2D7101E1191E1B489C1A418361D9F0C17F3A38",
	// 		"000BE5A7D14480EFD9CF038C51FE798BF752B0FD5EC00AEE3678D923837D753DA90ADD45154A7BAA28BC237496A9870DF7239CEF71526C73EA066E9AF34E74EBE03B8398E8F46C0A6DAB7842598EF0DCD5EF4BC019F8D20639838F1A8ED65012",
	// 		"0036AB139835FC9D146E519F99A11AAF557E0DF8A971DCB5379A66466C066DB4783716D9C17657F00AD4E68DA91FCBD254F21092F35A4ECE4016B29AE863C19008DE8105503B0BF57AD89B4D707EC51E0EB28D8C49B0AEC5D11C6E7D45713592",
	// 		"00E98C7B72100B181876B1D159643543C803E7B84C8F15BF7DE53062DCF4B922BF7F464BB2A1F6B5891F01143E3930DC8561AEF6D09EDA8DDB03673603C20F5E3FAB9F2E7911792FF789ED9CFEF6BFC45A7AC86B7E8010E7626ED6D2D642AED8",
	// 		"00013D59C89577690D7083AC6BE49C052DDE8CEC2C7A776076644227AC84485372138E60ECB2835069C8059F633C0BA9967F76397861D3A3227D78EF65104E675FC65EC268DF1681F3D0601EB2EB689B644BCAAF5912E1936B70CC017E9F25ED",
	// 	),
	// )
	// if err != nil {
	// 	t.Fatal(err)
	// }

	G1, G2 := g1.One(), g2.One()

	bw6.AddPair(G1, G2)

	actual := bw6.fp6.one()
	bw6.millerLoop(actual)
	if !expected.equal(actual) {
		t.Fatalf("expected != actual")
	}
}
func TestPairingFinalExp(t *testing.T) {
	bw6 := NewEngine()
	g1, g2 := bw6.G1, bw6.G2
	// e(a*G1, b*G2) = e(G1, G2)^c
	expected := new(fe6)
	tmp := new(fe)
	tmp, _ = fromString("00BE64FE0B5406B66F0A022E3580AC6D06CE4120E47DE81EFF70E9C0CF73CF2E4931D5DDA2805079C6383C7D696D6D8B3952B8F1E9EC995B7D6147FC1EE97641CCC27644CD905282B0A87F554A61F4457D29FD1163DD39E019E89F7A1B09D2AB")
	expected[0][0].set(tmp)
	tmp, _ = fromString("0071056F4F5862343DE7A18942A1F6A4B24257BC60D3820437E33C374943153FFC29B1A9812B6F27A704826CEF5E9D8D6412A33443C490464CDB57319EB2A393592EB80140D9F39C68E20EF8138B1375A4EAEE503B1077E0BB7612FF3BA19FDF")
	expected[0][1].set(tmp)
	tmp, _ = fromString("00A7DEBDD2D2B712D9EF7DBD6B8840B9B6ECCE5DE3C631FB14676C849F1D839BCE995A18AF25E0E154C0F5854D81B6ABDC7C18E5E4380111EE5F95A51B1CC08084C6BADF64B80431029911ED13D165A92C5D3A60C4B6B701DB09E5AD8713CF33")
	expected[0][2].set(tmp)
	tmp, _ = fromString("009C94404EFD09DAC985C2B82ECDD08374E82B3BEFD767C997520C277EB6C56DCEDB059CB831C2B393374D95A0438D84FDE4259309349EF86BCD55EC422F3618BB539378E407B89779D3DBF7B7E412C5EFE04E220C10A2790F1A58263F699689")
	expected[1][0].set(tmp)
	tmp, _ = fromString("009263C36801AC6C5626F28D85867F80AB684E3A3E5C5E0E8BA32B876728E8ADCFCB556BA7A2D661F849E985D4FB15909D29C80BE33760C872D6EE16117AE127F39DB7BEC503E3028ABBEC925AE37E40F988E2CCA427AF28C365B51479E83C90")
	expected[1][1].set(tmp)
	tmp, _ = fromString("00477D78FFA08531DF538752849578C78F2C66458DB8A28C27CE7802B03456880B844A03A571DB64B8988BA8B50D6597D561EE93A71D771A529CC56AFDA5C0CDB3C756CD5279D53C3F08E2550F98EE122936E8B6597F9F81E839D01F39CAA971")
	expected[1][2].set(tmp)

	G1, G2 := g1.One(), g2.One()

	actual := bw6.AddPair(G1, G2).Result()

	if !expected.equal(actual) {
		t.Fatalf("expected != actual")
	}
}

func TestPairingNonDegeneracy(t *testing.T) {
	bw6 := NewEngine()
	G1, G2 := bw6.G1, bw6.G2
	g1Zero, g2Zero, g1One, g2One := G1.Zero(), G2.Zero(), G1.One(), G2.One()
	GT := bw6.GT()
	// e(g1^a, g2^b) != 1
	bw6.Reset()
	{
		bw6.AddPair(g1One, g2One)
		e := bw6.Result()
		if e.IsOne() {
			t.Fatal("pairing result is not expected to be one")
		}
		if !GT.IsValid(e) {
			t.Fatal("pairing result is not valid")
		}
	}
	// e(g1^a, 0) == 1
	bw6.Reset()
	{
		bw6.AddPair(g1One, g2Zero)
		e := bw6.Result()
		if !e.IsOne() {
			t.Fatal("pairing result is expected to be one")
		}
	}
	// e(0, g2^b) == 1
	bw6.Reset()
	{
		bw6.AddPair(g1Zero, g2One)
		e := bw6.Result()
		if !e.IsOne() {
			t.Fatal("pairing result is expected to be one")
		}
	}
	//
	bw6.Reset()
	{
		bw6.AddPair(g1Zero, g2One)
		bw6.AddPair(g1One, g2Zero)
		bw6.AddPair(g1Zero, g2Zero)
		e := bw6.Result()
		if !e.IsOne() {
			t.Fatal("pairing result is expected to be one")
		}
	}
}

func TestPairingBilinearity(t *testing.T) {
	bw6 := NewEngine()
	g1, g2 := bw6.G1, bw6.G2
	gt := bw6.GT()
	// e(a*G1, b*G2) = e(G1, G2)^c
	{
		a, b := big.NewInt(17), big.NewInt(117)
		c := new(big.Int).Mul(a, b)
		G1, G2 := g1.One(), g2.One()
		e0 := bw6.AddPair(G1, G2).Result()
		P1, P2 := g1.New(), g2.New()
		g1.MulScalar(P1, G1, a)
		g2.MulScalar(P2, G2, b)
		e1 := bw6.AddPair(P1, P2).Result()
		gt.Exp(e0, e0, c)
		if !e0.Equal(e1) {
			t.Fatal("bad pairing, 1")
		}
	}
	bw6.Reset() // TODO
	// e(a * G1, b * G2) = e((a * b) * G1, G2)
	{
		// scalars
		a, b := big.NewInt(17), big.NewInt(117)
		c := new(big.Int).Mul(a, b)
		// LHS
		G1, G2 := g1.One(), g2.One()
		g1.MulScalar(G1, G1, c)
		bw6.AddPair(G1, G2)
		// RHS
		P1, P2 := g1.One(), g2.One()
		g1.MulScalar(P1, P1, a)
		g2.MulScalar(P2, P2, b)
		bw6.AddPairInv(P1, P2)
		// should be one
		if !bw6.Check() {
			t.Fatal("bad pairing, 2")
		}
	}
}

func TestPairingMulti(t *testing.T) {
	// e(G1, G2) ^ t == e(a01 * G1, a02 * G2) * e(a11 * G1, a12 * G2) * ... * e(an1 * G1, an2 * G2)
	// where t = sum(ai1 * ai2)
	bw6 := NewEngine()
	g1, g2 := bw6.G1, bw6.G2
	numOfPair := 100
	targetExp := new(big.Int)
	// RHS
	for i := 0; i < numOfPair; i++ {
		// (ai1 * G1, ai2 * G2)
		a1, a2 := randScalar(q), randScalar(q)
		P1, P2 := g1.One(), g2.One()
		g1.MulScalar(P1, P1, a1)
		g2.MulScalar(P2, P2, a2)
		bw6.AddPair(P1, P2)
		// accumulate targetExp
		// t += (ai1 * ai2)
		a1.Mul(a1, a2)
		targetExp.Add(targetExp, a1)
	}
	// LHS
	// e(t * G1, G2)
	T1, T2 := g1.One(), g2.One()
	g1.MulScalar(T1, T1, targetExp)
	bw6.AddPairInv(T1, T2)
	if !bw6.Check() {
		t.Fatal("fail multi pairing")
	}
}

func TestPairingEmpty(t *testing.T) {
	bw6 := NewEngine()
	if !bw6.Check() {
		t.Fatal("empty check should be accepted")
	}
	if !bw6.Result().IsOne() {
		t.Fatal("empty pairing result should be one")
	}
}

func TestNaf(t *testing.T) {
	k, ok := new(big.Int).SetString("23ed1347970dec008a442f991fffffffffffffffffffffff", 16)
	if !ok {
		t.Fatalf("invalid hex")
	}
	expected := []int8{
		-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 0, 0, 1, 0, 0, -1, 0, 1, 0, -1, 0, 0, 0, 0, -1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1,
		0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, -1, 0, 0, 0, 0, -1, 0, 0,
		1, 0, 0, 0, -1, 0, 0, -1, 0, 1, 0, -1, 0, 0, 0, 1, 0, 0, 1, 0, -1, 0, 1, 0, 1, 0, 0, 0, 1,
		0, -1, 0, -1, 0, 0, 0, 0, 0, 1, 0, 0, 1}
	result := computeNaf(k)
	for i := 0; i < len(expected); i++ {
		if result[i] != expected[i] {
			t.Fatal("not match")
		}
	}
}

func BenchmarkPairing(t *testing.B) {
	bw6 := NewEngine()
	g1, g2, gt := bw6.G1, bw6.G2, bw6.GT()
	bw6.AddPair(g1.One(), g2.One())
	e := gt.New()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		e = bw6.calculate()
	}
	_ = e
}
