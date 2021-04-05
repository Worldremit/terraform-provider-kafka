package kafka

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testProvider *schema.Provider
var testBootstrapServers []string = bootstrapServersFromEnv()

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	meta := testProvider.Meta()
	if meta == nil {
		t.Fatal("Could not construct client")
	}
	client := meta.(*LazyClient)
	if client == nil {
		t.Fatal("No client")
	}
	if err := client.init(); err != nil {
		t.Fatalf("Bad init %v", err)
	}
}

func overrideProvider() (*schema.Provider, error) {
	log.Println("[INFO] Setting up override for a provider")
	provider := Provider()

	diags := provider.Configure(context.Background(), accTestProviderConfig())
	if diags.HasError() {
		log.Printf("[ERROR] Could not configure provider %v", diags)
		return nil, fmt.Errorf("Could not configure provider")
	}

	testProvider = provider
	return provider, nil
}

func accTestProviderConfig() *terraform.ResourceConfig {
	bootstrapServers := bootstrapServersFromEnv()
	bs := make([]interface{}, len(bootstrapServers))

	for i, s := range bootstrapServers {
		bs[i] = s
	}

	ca, err := ioutil.ReadFile("../secrets/ca.crt")
	if err != nil {
		panic(err)
	}
	cert, err := ioutil.ReadFile("../secrets/client.pem")
	if err != nil {
		panic(err)
	}
	key, err := ioutil.ReadFile("../secrets/client.key")
	if err != nil {
		panic(err)
	}

	raw := map[string]interface{}{
		"bootstrap_servers": bs,
		"ca_cert":           string(ca),
		"client_cert":       string(cert),
		"client_key":        string(key),
	}
	return terraform.NewResourceConfigRaw(raw)
}

func bootstrapServersFromEnv() []string {
	fromEnv := strings.Split(os.Getenv("KAFKA_BOOTSTRAP_SERVER"), ",")
	fromEnv = nonEmptyAndTrimmed(fromEnv)

	if len(fromEnv) == 0 {
		fromEnv = []string{"localhost:9092"}
	}

	bootstrapServers := make([]string, 0)
	for _, bs := range fromEnv {
		if bs != "" {
			bootstrapServers = append(bootstrapServers, bs)
		}
	}

	return bootstrapServers
}
