package block

import (
	"fmt"
	"io"
	"strings"

	cmn "github.com/dfinlab/meter/libs/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type QuorumCert struct {
	QCHeight uint64
	QCRound  uint64
	EpochID  uint64

	VoterBitArrayStr string
	VoterMsgHash     [32]byte // [][32]byte
	VoterAggSig      []byte
	voterBitArray    *cmn.BitArray
	VoterViolation   []*Violation
}

func (qc *QuorumCert) String() string {
	if qc != nil {
		bitArray := strings.ReplaceAll(qc.VoterBitArrayStr, "\"", "")
		return fmt.Sprintf("QC(Height:%v, Round:%v, EpochID:%v, BitArray:%v, AggSig:len(%v))",
			qc.QCHeight, qc.QCRound, qc.EpochID, bitArray, len(qc.VoterAggSig))
	}
	return "QC(nil)"
}

func (qc *QuorumCert) CompactString() string {
	if qc != nil {
		return fmt.Sprintf("QC(Height:%v, Round:%v, EpochID:%v)",
			qc.QCHeight, qc.QCRound, qc.EpochID)
	}
	return "QC(nil)"
}

func (qc *QuorumCert) ToBytes() []byte {
	bytes, _ := rlp.EncodeToBytes(qc)
	return bytes
}

// EncodeRLP implements rlp.Encoder.
func (qc *QuorumCert) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		qc.QCHeight,
		qc.QCRound,
		qc.EpochID,
		qc.VoterMsgHash,
		qc.VoterAggSig,
		qc.VoterBitArrayStr,
	})
}

// DecodeRLP implements rlp.Decoder.
func (qc *QuorumCert) DecodeRLP(s *rlp.Stream) error {
	payload := struct {
		QCHeight         uint64
		QCRound          uint64
		EpochID          uint64
		VoterMsgHash     [32]byte
		VoterAggSig      []byte
		VoterBitArrayStr string
	}{}

	if err := s.Decode(&payload); err != nil {
		return err
	}

	*qc = QuorumCert{
		QCHeight:         payload.QCHeight,
		QCRound:          payload.QCRound,
		EpochID:          payload.EpochID,
		VoterMsgHash:     payload.VoterMsgHash,
		VoterAggSig:      payload.VoterAggSig,
		VoterBitArrayStr: payload.VoterBitArrayStr,
	}
	return nil
}

func (qc *QuorumCert) VoterBitArray() *cmn.BitArray {
	bitArray := &cmn.BitArray{}
	err := bitArray.UnmarshalJSON([]byte(qc.VoterBitArrayStr))
	if err != nil {
		return nil
	}
	return bitArray
}

func (qc *QuorumCert) GetViolation() []*Violation {
	return qc.VoterViolation
}

func GenesisQC() *QuorumCert {
	return &QuorumCert{QCHeight: 0, QCRound: 0, EpochID: 0}
}

//--------------
func QCEncodeBytes(qc *QuorumCert) []byte {
	blockBytes, _ := rlp.EncodeToBytes(qc)
	return blockBytes
}

func QCDecodeFromBytes(bytes []byte) (*QuorumCert, error) {
	qc := QuorumCert{}
	err := rlp.DecodeBytes(bytes, &qc)
	return &qc, err
}
