package main

import (
	"context"
	"fmt"

	"github.com/fehmicorp/pkg/v1/cf/dns"
	"github.com/fehmicorp/pkg/v1/utils/os"
)

var CF_ACCOUNT_ID = "c483f3eb849ad7ed42b9889caf62f127"
var CF_API_TOKEN = "cfat_pCXtIMcHsU8anKlydVCYumQuK18CB0j88TkA4krK0ea6708a"

func main() {
	os.SetEnv("CF_ACCOUNT_ID", CF_ACCOUNT_ID)
	os.SetEnv("CF_API_TOKEN", CF_API_TOKEN)
	ctx := context.Background()
	client := dns.Fetch()
	Names := dns.GetZonesNames(ctx, client)
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("Using Zone Name: %s \n", Names)
	Ids := dns.GetZoneIDs(ctx, client)
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("Using Zone Id: %s \n", Ids)

}

// func GetZonewithRecords() {

// 	fmt.Println("Fetching Zones...")
// 	zones, err := client.GetZones(ctx)
// 	if err != nil {
// 		log.Fatalf("Error getting zones: %v", err)
// 	}
// 	if len(zones) == 0 {
// 		log.Fatal("No active zones found in account.")
// 	}
// 	targetZone := zones[0]
// 	fmt.Printf("Using Zone: %s (ID: %s)\n\n", targetZone.Name, targetZone.ID)

// 	fmt.Println("Fetching DNS records...")
// 	records, err := client.ListDNSRecords(ctx, targetZone.ID)
// 	if err != nil {
// 		log.Fatalf("Failed to fetch DNS records: %v", err)
// 	}

// 	// 5. Display fetched records
// 	fmt.Printf("Found %d record(s):\n", len(records))
// 	fmt.Println("------------------------------------------------------------------")
// 	for _, rec := range records {
// 		proxiedStr := "DNS Only"
// 		if rec.Proxied != nil && *rec.Proxied {
// 			proxiedStr = "Proxied"
// 		}
// 		fmt.Printf("[%s]\t%-30s -> %-25s (%s, TTL: %d)\n", rec.Type, rec.Name, rec.Content, proxiedStr, rec.TTL)
// 	}
// 	fmt.Println("------------------------------------------------------------------")
// }
