package main

import (
	"fmt"
	"log"
	"os"

	"github.com/osrg/gobgp/pkg/packet/bgp"
	mrt "github.com/osrg/gobgp/pkg/packet/mrt"
)

func Parser(filepath string, flagWhois bool) (records TRecords) {
	var err error

	var data []byte
	var dataOffset int

	var loopsLogged map[uint32]struct{}
	var prependingsLogged map[uint32]struct{}

	/*** * * ***/

	records = make(TRecords)
	loopsLogged = make(map[uint32]struct{})
	prependingsLogged = make(map[uint32]struct{})

	/*** * * ***/

	data, err = os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	dataOffset = 0

	/*** * * ***/

	for dataOffset < len(data) {
		var err error

		var hdr *mrt.MRTHeader

		var bodyEnd int

		var msg *mrt.MRTMessage

		var rib *mrt.Rib
		var ribOk bool

		/*** * * ***/

		if dataOffset+mrt.MRT_COMMON_HEADER_LEN > len(data) {
			log.Print("Incomplete header at offset", dataOffset)

			break
		}

		/*** * * ***/

		hdr = &mrt.MRTHeader{}

		/*** * * ***/

		err = hdr.DecodeFromBytes(data[dataOffset : dataOffset+mrt.MRT_COMMON_HEADER_LEN])
		if err != nil {
			log.Print("Error decoding header at offset %d: %v\n", dataOffset, err)

			break
		}

		bodyEnd = dataOffset + mrt.MRT_COMMON_HEADER_LEN + int(hdr.Len)
		if bodyEnd > len(data) {
			log.Print("Incomplete body at offset : ", dataOffset)

			break
		}

		msg, err = mrt.ParseMRTBody(hdr, data[dataOffset+mrt.MRT_COMMON_HEADER_LEN:bodyEnd])
		if err != nil {
			log.Print("Error parsing MRT body at offset : ", dataOffset, err)

			dataOffset = bodyEnd

			continue
		}

		// Check if the message body is a RIB
		rib, ribOk = msg.Body.(*mrt.Rib)
		if !ribOk {
			log.Print("Message body is not a RIB at offset :", dataOffset)

			dataOffset = bodyEnd

			continue
		}

		// Find loops
		for _, entry := range rib.Entries {
			for _, pathAttribute := range entry.PathAttributes {
				if pathAttributeAsPath, ok := pathAttribute.(*bgp.PathAttributeAsPath); ok {
					for _, asPathParamInterface := range pathAttributeAsPath.Value {
						ases := asPathParamInterface.GetAS()

						seen := make(map[uint32]bool)

						for i := range ases {
							currentAS := ases[i]

							// Skip first AS
							if i == 0 {
								seen[currentAS] = true
								continue
							}

							if seen[currentAS] {
								var loop bool = false
								var repeatConsecutive bool = false
								var repeatNonConsecutive bool = false

								/*** * * ***/

								if records[currentAS] == nil {
									records[currentAS] = &TRecord{}
								}

								/*** * * ***/

								// Check
								if i > 2 && ases[i-2] == currentAS && ases[i-1] != currentAS {
									loop = true
								} else if ases[i-1] == currentAS {
									repeatConsecutive = true
								} else {
									repeatNonConsecutive = true
								}

								/*** * * ***/

								if loop {
									records[currentAS].Loops++
								} else if repeatConsecutive {
									records[currentAS].Repeats++
									records[currentAS].ConsecutiveRepeats++
								} else if repeatNonConsecutive {
									records[currentAS].Repeats++
									records[currentAS].NonConsecutiveRepeats++
								}

								if records[currentAS].Loops > _CONST_LOG_AFTER_LOOPS_NUM {
									if _, exists := loopsLogged[currentAS]; !exists {
										log.Print("AS"+fmt.Sprintf("%d", currentAS) + " exceeded " + fmt.Sprintf("%d", _CONST_LOG_AFTER_LOOPS_NUM) + " loops.")
										log.Print("AS"+fmt.Sprintf("%d", currentAS) + " loops example: ", ases)

										loopsLogged[currentAS] = struct{}{}
									}
								}

								if records[currentAS].Repeats > _CONST_LOG_AFTER_REPEATS_NUM {
									if _, exists := prependingsLogged[currentAS]; !exists {
										log.Print("AS"+fmt.Sprintf("%d", currentAS) + " exceeded " + fmt.Sprintf("%d", _CONST_LOG_AFTER_REPEATS_NUM) + " repeats.")
										log.Print("AS"+fmt.Sprintf("%d", currentAS) + " prependings example: ", ases)

										prependingsLogged[currentAS] = struct{}{}
									}
								}

								if flagWhois {
									if records[currentAS].Loops > _CONST_WHOIS_AFTER_LOOPS_NUM {
										if len(records[currentAS].Name) == 0 {
											records[currentAS].Name = getASName(currentAS, true)
										}
									}
								}
							} else {
								seen[currentAS] = true
							}
						}
					}
				}
			}
		}

		/*** * * ***/

		dataOffset = bodyEnd
	}

	/*** * * ***/

	return records
}
