package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

var theJson = strings.TrimSpace(`
{
  "data": [
    {
      "batches": [
        {
          "header": {
            "signer_public_key": "023756357e7eadf66a8866b42b87aa50c8ba77f35488d620c80237f6dd8a06804c",
            "transaction_ids": [
              "0a2633ec68ed8e8117acf80dd38efe957fe4752e830f51a0b53ab12496b308a11107ecafccf8172dab84e74ff107185f8a7d238341fe05a33a10b45caf91dc94"
            ]
          },
          "header_signature": "c77b688795305b45026d33e72f7cc71c1e431b07ac1fc4341e754f4f2fb54fa86c6aeaf56c81eb872edd02d9d07c77c22640737f4a8f6aea5039bbd43b6d3c2d",
          "trace": false,
          "transactions": [
            {
              "header": {
                "batcher_public_key": "023756357e7eadf66a8866b42b87aa50c8ba77f35488d620c80237f6dd8a06804c",
                "dependencies": [],
                "family_name": "sawtooth_settings",
                "family_version": "1.0",
                "inputs": [
                  "000000a87cb5eafdcca6a8cde0fb0dec1400c5ab274474a6aa82c1c0cbf0fbcaf64c0b",
                  "000000a87cb5eafdcca6a8cde0fb0dec1400c5ab274474a6aa82c12840f169a04216b7",
                  "000000a87cb5eafdcca6a8cde0fb0dec1400c5ab274474a6aa82c1918142591ba4e8a7",
                  "000000a87cb5eafdcca6a8cde0fb0dec1400c5ab274474a6aa82c12840f169a04216b7"
                ],
                "nonce": "",
                "outputs": [
                  "000000a87cb5eafdcca6a8cde0fb0dec1400c5ab274474a6aa82c1c0cbf0fbcaf64c0b",
                  "000000a87cb5eafdcca6a8cde0fb0dec1400c5ab274474a6aa82c12840f169a04216b7"
                ],
                "payload_sha512": "89ad0a4248d8f9d2519c69e6ec1949ad55853043854ae3b94e2e87841456465b3e468f23b5ebc42eb6b8b1d060c4ef0ded6bf28166a9626736dd3a8ede059428",
                "signer_public_key": "023756357e7eadf66a8866b42b87aa50c8ba77f35488d620c80237f6dd8a06804c"
              },
              "header_signature": "0a2633ec68ed8e8117acf80dd38efe957fe4752e830f51a0b53ab12496b308a11107ecafccf8172dab84e74ff107185f8a7d238341fe05a33a10b45caf91dc94",
              "payload": "CAESgAEKJnNhd3Rvb3RoLnNldHRpbmdzLnZvdGUuYXV0aG9yaXplZF9rZXlzEkIwMjM3NTYzNTdlN2VhZGY2NmE4ODY2YjQyYjg3YWE1MGM4YmE3N2YzNTQ4OGQ2MjBjODAyMzdmNmRkOGEwNjgwNGMaEjB4YTUyM2M4OTBmZmViNjk2YQ=="
            }
          ]
        }
      ],
      "header": {
        "batch_ids": [
          "c77b688795305b45026d33e72f7cc71c1e431b07ac1fc4341e754f4f2fb54fa86c6aeaf56c81eb872edd02d9d07c77c22640737f4a8f6aea5039bbd43b6d3c2d"
        ],
        "block_num": "0",
        "consensus": "R2VuZXNpcw==",
        "previous_block_id": "0000000000000000",
        "signer_public_key": "03c374673b2150ba5ddd97cea9d748a20f0418291fd3386c7ffcba6ad5d5189581",
        "state_root_hash": "74d482b59034aaa6c2bb607f7833c5f6591940551e81f30e653946656824ddca"
      },
      "header_signature": "e8593bcb16e8baa170ab738f82fb24e97e6803df4fe311fd048df183a75a5bb4391a7915451610ac822fdfdba584bcff286a29e57a137dff5f8d3a6f49cb8a48"
    },
    {
      "batches": [
        {
          "header": {
            "signer_public_key": "0213e4edc8e06e90f4163b10c73385c39a2ba06b665bd7124532ad27e30cafb6d3",
            "transaction_ids": [
              "bd52421c121048bb3abbb2542f4b1f2f6c83a4b260e71ad83393a76147cda4124b3152ab9f0f6ca4f46a2c59f4a52ab2c01f6cb8aef1acb8c9e1caa443c33ebf"
            ]
          },
          "header_signature": "24d38f2f6f533c669a0952858e12fff62c97ef708ae8f0481189662e8702e43c01f9be27690b6740438b6f991701fbd523d47b9749986f254a3abacdd43b541e",
          "trace": false,
          "transactions": [
            {
              "header": {
                "batcher_public_key": "0213e4edc8e06e90f4163b10c73385c39a2ba06b665bd7124532ad27e30cafb6d3",
                "dependencies": [],
                "family_name": "currency",
                "family_version": "1.0",
                "inputs": [
                  "65ec1c0253f016ae8d0e411d4951928a2e74fc949595962a1c452f4846eda8e38196c3"
                ],
                "nonce": "",
                "outputs": [
                  "65ec1c0253f016ae8d0e411d4951928a2e74fc949595962a1c452f4846eda8e38196c3"
                ],
                "payload_sha512": "f3ca7d1e0cfe79d3a09608e8e37e4052bbee02f08a83fe6048a876ca031edaf1bd8239a6a9a7467efba2d7d1cfe43136d69c61113ba222f1e74649ab181f5cdd",
                "signer_public_key": "0213e4edc8e06e90f4163b10c73385c39a2ba06b665bd7124532ad27e30cafb6d3"
              },
              "header_signature": "bd52421c121048bb3abbb2542f4b1f2f6c83a4b260e71ad83393a76147cda4124b3152ab9f0f6ca4f46a2c59f4a52ab2c01f6cb8aef1acb8c9e1caa443c33ebf",
              "payload": "CAU="
            }
          ]
        }
      ],
      "header": {
        "batch_ids": [
          "24d38f2f6f533c669a0952858e12fff62c97ef708ae8f0481189662e8702e43c01f9be27690b6740438b6f991701fbd523d47b9749986f254a3abacdd43b541e"
        ],
        "block_num": "1",
        "consensus": "RGV2bW9kZWLYMLx/HRbOu/8pkQf1V7oKOrRM9hKwg7y/ucIkqR/y",
        "previous_block_id": "e8593bcb16e8baa170ab738f82fb24e97e6803df4fe311fd048df183a75a5bb4391a7915451610ac822fdfdba584bcff286a29e57a137dff5f8d3a6f49cb8a48",
        "signer_public_key": "03c374673b2150ba5ddd97cea9d748a20f0418291fd3386c7ffcba6ad5d5189581",
        "state_root_hash": "7ee80d28c12aee50142329461e012aca4147da65ec619147310ee6470175cc57"
      },
      "header_signature": "d7338b244b5a12d20c3e710e1098d377be26fff645ef49240658ccc06eaad4025312d111f1cb6fcc713f9ef61ea492577b6ed56438a7cf390675d63f9b3fa16d"
    }
  ],
  "head": "d7338b244b5a12d20c3e710e1098d377be26fff645ef49240658ccc06eaad4025312d111f1cb6fcc713f9ef61ea492577b6ed56438a7cf390675d63f9b3fa16d",
  "link": "http://77.93.70.138:8008/blocks?head=d7338b244b5a12d20c3e710e1098d377be26fff645ef49240658ccc06eaad4025312d111f1cb6fcc713f9ef61ea492577b6ed56438a7cf390675d63f9b3fa16d&start=0x0000000000000000&limit=100&reverse",
  "paging": {
    "limit": null,
    "start": null
  }
}
`)

type sawtoothBlocksResponse struct {
	Head string
	Link string
	Data []*blocksData
}

type blocksData struct {
	Header_signature string
}

func main() {
	response := sawtoothBlocksResponse{}
	err := json.Unmarshal([]byte(theJson), &response)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("No error")
	fmt.Println(response.Data[0].Header_signature)
}
