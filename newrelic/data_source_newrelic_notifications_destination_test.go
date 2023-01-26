//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicNotificationsDestinationDataSource_Basic(t *testing.T) {
	resourceName := "newrelic_notifications_destination.foo"
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf-notifications-test-%s", rand)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicNotificationsDestinationDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicNotificationsDestination("data.newrelic_notifications_destination.foo-source"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
				),
			},
		},
	})
}

func testAccNewRelicNotificationsDestinationDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_notification_destination" "foo" {
	account_id = %[1]d
	name = "%[2]s"
	type = "WEBHOOK"
	active = true

	property {
		key = "url"
		value = "https://webhook.site/"
	}
}

data "newrelic_notification_destination" "foo-source" {
	id = newrelic_notification_destination.foo.id
}
`, accountID, name)
}

func testAccNewRelicNotificationsDestination(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get a notification destination from New Relic")
		}

		return nil
	}
}
