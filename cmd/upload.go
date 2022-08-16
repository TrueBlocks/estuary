/*
Copyright © 2022 TrueBlocks

The use of this program and source code are governed by the text
you will find in the LICENSE file at the root of this repository.
*/
package cmd

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config/scrape"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/pinning"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long:  "",
	RunE:  UploadFile,
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func UploadFile(cmd *cobra.Command, args []string) error {
	local := pinning.Service{}
	localHash, err := local.PinFile("main.go", true)
	if err != nil {
		return err
	}

	pinataKey, pinataSecret, estuaryKey := scrape.PinningKeys("mainnet")
	pinata := pinning.Service{
		Apikey:     pinataKey,
		Secret:     pinataSecret,
		PinUrl:     "https://api.pinata.cloud/pinning/pinFileToIPFS",
		ResultName: "IpfsHash",
		HeaderFunc: PinataHeaders,
	}
	pinataHash, err := pinata.PinFile("main.go", false)
	if err != nil {
		return err
	}
	fmt.Println("local ==> pinata", localHash, pinataHash, (localHash == pinataHash))

	estuary := pinning.Service{
		Apikey:     estuaryKey,
		PinUrl:     "https://api.estuary.tech/content/add",
		HeaderFunc: EstuaryHeaders,
	}
	estuaryHash, err := estuary.PinFile("main.go", false)
	if err != nil {
		return err
	}
	fmt.Println("local ==> estuary", localHash, estuaryHash, (localHash == estuaryHash))

	return nil
}

func PinataHeaders(s *pinning.Service, contentType string) map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = contentType
	headers["pinata_secret_api_key"] = s.Secret
	headers["pinata_api_key"] = s.Apikey
	return headers
}

func EstuaryHeaders(s *pinning.Service, contentType string) map[string]string {
	headers := make(map[string]string)
	headers["Content-Type"] = contentType
	headers["Authorization"] = "Bearer " + s.Apikey
	return headers
}

/*
curl
 -X POST https://api.estuary.tech/content/add \
 -H "Authorization: Bearer REPLACE_ME_WITH_API_KEY" \
 -H "Content-Type: multipart/form-data" \
 -F "data=@PATH_TO_YOUR_FILE"

class Example extends React.Component {
  upload(e) {
    const formData  = new FormData();
    formData.append("data", e.target.files[0]);

    const xhr = new XMLHttpRequest();
    xhr.upload.onprogress = (event) => {
      this.setState({
        loaded: event.loaded,
        total: event.total
      });
    }

    xhr.open(
      "POST",
      "https://api.estuary.tech/content/add"
    );
    xhr.setRequestHeader(
      "Authorization",
      "Bearer REPLACE_ME_WITH_API_KEY"
    );
    xhr.send(formData);
  }

  render() {
    return (
      <React.Fragment>
        <input type="file" onChange={this.upload.bind(this)} />
        <br />
        <pre>{JSON.stringify(this.state, null, 1)}</pre>
      </React.Fragment>
    );
  }
}
*/
