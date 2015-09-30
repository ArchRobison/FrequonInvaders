package vanity

import (
	"bufio"
	"crypto/rand"
	"crypto/sha1"
	"io"
	"os"
)

func vanityBoardPath() string {
	return "vanity.dat" // FIXME
}

// The cryptographic signature checking can of course be circumvented
// by inspecting the source and generating a fake signature. But I wanted
// to try out Go's support for crypto stuff, so here it is.

func hash(v []Record, salt []byte) [sha1.Size]byte {
	data := []byte{}
	copy(data, salt)
	for i := range v {
		data = append(data, byte(v[i].Score^uint8(i)))
		data = append(data, []byte(v[i].Name)...)
	}
	return sha1.Sum(data)
}

const saltLen = 5

func WriteToFile(v []Record) {
	var sig [saltLen + sha1.Size]byte
	_, err := rand.Reader.Read(sig[:saltLen])
	if err != nil {
		panic(err)
	}
	h := hash(v, sig[:saltLen])
	copy(sig[saltLen:], h[:])

	path := vanityBoardPath()
	file, err := os.Create(path)
	if err == nil {
		defer file.Close()
		f := bufio.NewWriter(file)
		defer f.Flush()
		// Write signature
		f.Write(sig[:])
		// Write each record
		for i, _ := range v {
			f.WriteByte(v[i].Score)
			f.WriteString(v[i].Name)
			f.WriteByte(0)
		}
	}
}

func ReadFromFile() (v []Record, err error) {
	path := vanityBoardPath()
	file, err := os.Open(path)
	// Read records
	v = []Record{}
	if err == nil {
		defer file.Close()
		f := bufio.NewReader(file)
		var sig [saltLen + sha1.Size]byte
		// Read signature
		var n int
		n, err = f.Read(sig[:])
		if err != nil || n != len(sig[:]) {
			return
		}
		for {
			var score byte
			score, err = f.ReadByte()
			if err == io.EOF {
				err = nil
				return
			}
			if err != nil {
				break
			}
			name, err := f.ReadBytes(0)
			if err != nil {
				break
			}
			v = append(v, Record{Score: uint8(score), Name: string(name[:len(name)-1])})
		}
		v = v[:0]
	}
	return
}
