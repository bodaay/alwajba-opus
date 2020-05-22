package main

import (
	"fmt"
	"io"
	"os"

	"github.com/bodaay/alwajba-opus/opus"
	"github.com/mccoyst/ogg"
	"github.com/youpy/go-wav"
)

// UInt64toByteArray unt64 bit to byte array
func UInt64toByteArray(val uint64) []byte {
	r := make([]byte, 8)
	for i := uint64(0); i < 8; i++ {
		r[i] = byte((val >> (i * 8)) & 0xff)
	}
	return r
}

// ByteArraytoUInt64 Byte array to Uint64
func ByteArraytoUInt64(val []byte) uint64 {
	r := uint64(0)
	for i := uint64(0); i < 8; i++ {
		r |= uint64(val[i]) << (8 * i)
	}
	return r
}

// UInt32toByteArray unt32 bit to byte array
func UInt32toByteArray(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

// ByteArraytoUInt32 Byte array to Uint32
func ByteArraytoUInt32(val []byte) uint32 {
	r := uint32(0)
	for i := uint32(0); i < 4; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}
func convertToOpus(filename string) {
	//test wav file load
	file, _ := os.Open(filename)
	reader := wav.NewReader(file)

	wavFormat, err := reader.Format()
	fmt.Printf("%#v\n", wavFormat)
	fmt.Printf("Audio format: %d\n", wavFormat.AudioFormat)
	fmt.Printf("Bits Per Sample: %d\n", wavFormat.BitsPerSample)
	fmt.Printf("Block Align: %d\n", wavFormat.BlockAlign)
	fmt.Printf("Bytes Rate: %d\n", wavFormat.ByteRate)
	fmt.Printf("Channels: %d\n", wavFormat.NumChannels)
	fmt.Printf("SampleRate: %d\n", wavFormat.SampleRate)

	enc, err := opus.NewEncoder(int(wavFormat.SampleRate), int(wavFormat.NumChannels), opus.AppVoIP)
	enc.SetDTX(false)
	enc.SetBitrateToAuto()
	// enc.SetBitrate(32000)
	// enc.SetBitrate(int(wavFormat.ByteRate))
	br, _ := enc.Bitrate()
	sr, _ := enc.SampleRate()
	fmt.Println(br)
	fmt.Println(sr)
	fmt.Println(wavFormat.SampleRate)
	if err != nil {
		panic(err)
	}
	framsSizeWanted := float32(20)
	batchSize := uint32((framsSizeWanted / 1000) * float32(wavFormat.SampleRate)) // we need 60ms per batch, that would be our frame size
	fmt.Printf("Batch size: %d\n", batchSize)
	frameSizeMs := float32(batchSize*uint32(wavFormat.NumChannels)) / float32(wavFormat.NumChannels) * 1000 / float32(wavFormat.SampleRate)
	fmt.Printf("Frame size ms: %f\n", frameSizeMs)

	type encodedFrame struct {
		size int
		data []byte
	}
	var eframes []encodedFrame
	// var testSamples []wav.Sample
	var allsamples []wav.Sample
	for {

		samples, err := reader.ReadSamples() // don't even think of send required samples, the library doesn't work properly, much better load whole shit into ram and then batch process it
		if err == io.EOF {
			break
		}
		for _, s := range samples {
			allsamples = append(allsamples, s)
		}
	}

	fmt.Printf("Total number of samples: %d\n", len(allsamples))
	batchIndex := 0
	for {

		// // fmt.Println(len(samples))
		last := false
		// fmt.Printf("Processing Batch index: %d\n", batchIndex)
		max := batchIndex + int(batchSize)
		if max > len(allsamples) {
			max = len(allsamples)
			last = true
		}
		batchsubset := allsamples[batchIndex:max]
		pcmraw := make([]int16, batchSize*uint32(wavFormat.NumChannels))
		byteindex := 0
		for _, sample := range batchsubset {
			pcmraw[byteindex] = int16(reader.IntValue(sample, 0)) //left channel or mono
			byteindex++
			if wavFormat.NumChannels > 1 {
				pcmraw[byteindex] = int16(reader.IntValue(sample, 1)) //right channel in stereo
				byteindex++
			}
		}
		data := make([]byte, 10000)
		n, err := enc.Encode(pcmraw, data)
		if err != nil {
			panic(err)
		}
		data = data[:n]
		eframes = append(eframes, encodedFrame{
			size: n,
			data: data,
		})
		batchIndex += int(batchSize)
		if last {
			break
		}
	}

	fmt.Printf("Total Frames: %d\n", len(eframes))

	//
	//
	//
	//
	// Decoding shit
	//
	//
	//
	//

	//test write back awv file
	fo, err := os.Create("output.wav")
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	writer := wav.NewWriter(fo, uint32(len(allsamples)), wavFormat.NumChannels, wavFormat.SampleRate, wavFormat.BitsPerSample)

	dec, err := opus.NewDecoder(int(wavFormat.SampleRate), int(wavFormat.NumChannels))
	if err != nil || dec == nil {
		panic(err)
	}

	var samples []wav.Sample
	// var outbytes []byte
	for _, ef := range eframes {
		pcmraw := make([]int16, 10000)
		samplesDecoded, err := dec.Decode(ef.data, pcmraw)
		if err != nil || dec == nil {
			panic(err)
		}
		pcmindex := 0
		for i := 0; i < samplesDecoded; i++ {
			sample := wav.Sample{}
			sample.Values[0] = int(pcmraw[pcmindex])
			pcmindex = pcmindex + 1
			if wavFormat.NumChannels > 1 {
				sample.Values[1] = int(pcmraw[pcmindex])
				pcmindex = pcmindex + 1
			}
			samples = append(samples, sample)

		}
		// for _, b := range pcmraw[:samplesDecoded] {
		// 	outbytes = append(outbytes, byte(b))
		// }

	}
	// writer.Write(outbytes)
	writer.WriteSamples(samples)

	//Test store opus ogg
	oggFile, err := os.Create("output.ogg")
	if err != nil {
		panic(err)
	}
	defer oggFile.Close()
	oggEnc := ogg.NewEncoder(1, oggFile)

	var opusHeader []byte
	opusHeader = append(opusHeader, []byte("OpusHead")...)
	opusHeader = append(opusHeader, byte(0x01))
	opusHeader = append(opusHeader, byte(wavFormat.NumChannels)) //Channels                                                  //version number, 1
	opusHeader = append(opusHeader, []byte{0x00, 0x00}...)
	opusHeader = append(opusHeader, UInt32toByteArray(wavFormat.SampleRate)...) //channels
	opusHeader = append(opusHeader, []byte{0x00, 0x00}...)                      //2 bytes
	opusHeader = append(opusHeader, byte(0x00))
	//identity commoen
	var opusIdentity []byte
	opusIdentity = append(opusIdentity, []byte("OpusTags")...)
	vendorString := "Khalefa"
	opusIdentity = append(opusIdentity, UInt32toByteArray(uint32(len(vendorString)))...) //vendor string length
	opusIdentity = append(opusIdentity, []byte(vendorString)...)                         //vendor string
	opusIdentity = append(opusIdentity, byte(0x00))                                      //user commit list length, 0, we don't have any comments

	opusHeader = append(opusHeader, opusIdentity...)
	err = oggEnc.EncodeBOS(0, opusHeader)
	if err != nil {
		panic(err)
	}

	startGran := int64(0)
	for _, ef := range eframes {
		err := oggEnc.Encode(startGran, ef.data)
		if err != nil {
			panic(err)
		}
		startGran++
	}

	oggEnc.EncodeEOS()

	fmt.Println("Done opus shit")
}
