package cloudwatchevents_test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	events "github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfcloudwatchevents "github.com/hashicorp/terraform-provider-aws/internal/service/cloudwatchevents"
)

func TestAccCloudWatchEventsPermission_basic(t *testing.T) {
	principal1 := "111111111111"
	principal2 := "*"
	statementID := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudwatch_event_permission.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, events.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckCloudWatchEventPermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPermissionResourceBasicConfig("", statementID),
				ExpectError: regexp.MustCompile(`must be \* or a 12 digit AWS account ID`),
			},
			{
				Config:      testAccCheckPermissionResourceBasicConfig(".", statementID),
				ExpectError: regexp.MustCompile(`must be \* or a 12 digit AWS account ID`),
			},
			{
				Config:      testAccCheckPermissionResourceBasicConfig("12345678901", statementID),
				ExpectError: regexp.MustCompile(`must be \* or a 12 digit AWS account ID`),
			},
			{
				Config:      testAccCheckPermissionResourceBasicConfig("abcdefghijkl", statementID),
				ExpectError: regexp.MustCompile(`must be \* or a 12 digit AWS account ID`),
			},
			{
				Config:      testAccCheckPermissionResourceBasicConfig(principal1, ""),
				ExpectError: regexp.MustCompile(`must be between 1 and 64 characters`),
			},
			{
				Config:      testAccCheckPermissionResourceBasicConfig(principal1, sdkacctest.RandString(65)),
				ExpectError: regexp.MustCompile(`must be between 1 and 64 characters`),
			},
			{
				Config:      testAccCheckPermissionResourceBasicConfig(principal1, " "),
				ExpectError: regexp.MustCompile(`must be one or more alphanumeric, hyphen, or underscore characters`),
			},
			{
				Config: testAccCheckPermissionResourceBasicConfig(principal1, statementID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchEventPermissionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action", "events:PutEvents"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "principal", principal1),
					resource.TestCheckResourceAttr(resourceName, "statement_id", statementID),
					resource.TestCheckResourceAttr(resourceName, "event_bus_name", tfcloudwatchevents.DefaultEventBusName),
				),
			},
			{
				Config: testAccCheckPermissionResourceBasicConfig(principal2, statementID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchEventPermissionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "principal", principal2),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccCheckPermissionResourceDefaultEventBusNameConfig(principal2, statementID),
				PlanOnly: true,
			},
		},
	})
}

func TestAccCloudWatchEventsPermission_eventBusName(t *testing.T) {
	principal1 := "111111111111"
	statementID := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	busName := sdkacctest.RandomWithPrefix("tf-acc-test-bus")

	resourceName := "aws_cloudwatch_event_permission.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, events.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckCloudWatchEventPermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPermissionResourceEventBusNameConfig(principal1, busName, statementID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchEventPermissionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action", "events:PutEvents"),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "principal", principal1),
					resource.TestCheckResourceAttr(resourceName, "statement_id", statementID),
					resource.TestCheckResourceAttr(resourceName, "event_bus_name", busName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudWatchEventsPermission_action(t *testing.T) {
	principal := "111111111111"
	statementID := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudwatch_event_permission.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, events.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckCloudWatchEventPermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckPermissionResourceActionConfig("", principal, statementID),
				ExpectError: regexp.MustCompile(`must be between 1 and 64 characters`),
			},
			{
				Config:      testAccCheckPermissionResourceActionConfig(sdkacctest.RandString(65), principal, statementID),
				ExpectError: regexp.MustCompile(`must be between 1 and 64 characters`),
			},
			{
				Config:      testAccCheckPermissionResourceActionConfig("events:", principal, statementID),
				ExpectError: regexp.MustCompile(`must be: events: followed by one or more alphabetic characters`),
			},
			{
				Config:      testAccCheckPermissionResourceActionConfig("events:1", principal, statementID),
				ExpectError: regexp.MustCompile(`must be: events: followed by one or more alphabetic characters`),
			},
			{
				Config: testAccCheckPermissionResourceActionConfig("events:PutEvents", principal, statementID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchEventPermissionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "action", "events:PutEvents"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudWatchEventsPermission_condition(t *testing.T) {
	statementID := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudwatch_event_permission.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, events.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckCloudWatchEventPermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPermissionResourceConditionOrganizationConfig(statementID, "o-1234567890"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchEventPermissionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.key", "aws:PrincipalOrgID"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.type", "StringEquals"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "o-1234567890"),
				),
			},
			{
				Config: testAccCheckPermissionResourceConditionOrganizationConfig(statementID, "o-0123456789"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchEventPermissionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "condition.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.key", "aws:PrincipalOrgID"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.type", "StringEquals"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "o-0123456789"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudWatchEventsPermission_multiple(t *testing.T) {
	principal1 := "111111111111"
	principal2 := "222222222222"
	statementID1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	statementID2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName1 := "aws_cloudwatch_event_permission.test"
	resourceName2 := "aws_cloudwatch_event_permission.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, events.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckCloudWatchEventPermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPermissionResourceBasicConfig(principal1, statementID1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchEventPermissionExists(resourceName1),
					resource.TestCheckResourceAttr(resourceName1, "principal", principal1),
					resource.TestCheckResourceAttr(resourceName1, "statement_id", statementID1),
				),
			},
			{
				Config: testAccCheckPermissionResourceMultipleConfig(principal1, statementID1, principal2, statementID2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchEventPermissionExists(resourceName1),
					testAccCheckCloudWatchEventPermissionExists(resourceName2),
					resource.TestCheckResourceAttr(resourceName1, "principal", principal1),
					resource.TestCheckResourceAttr(resourceName1, "statement_id", statementID1),
					resource.TestCheckResourceAttr(resourceName2, "principal", principal2),
					resource.TestCheckResourceAttr(resourceName2, "statement_id", statementID2),
				),
			},
		},
	})
}

func TestAccCloudWatchEventsPermission_disappears(t *testing.T) {
	resourceName := "aws_cloudwatch_event_permission.test"
	principal := "111111111111"
	statementID := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t) },
		ErrorCheck:   acctest.ErrorCheck(t, events.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckCloudWatchEventPermissionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPermissionResourceBasicConfig(principal, statementID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudWatchEventPermissionExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfcloudwatchevents.ResourcePermission(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckCloudWatchEventPermissionExists(pr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).CloudWatchEventsConn
		rs, ok := s.RootModule().Resources[pr]
		if !ok {
			return fmt.Errorf("Not found: %s", pr)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		eventBusName, statementID, err := tfcloudwatchevents.PermissionParseResourceID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error reading CloudWatch Events permission (%s): %w", pr, err)
		}
		input := &events.DescribeEventBusInput{
			Name: aws.String(eventBusName),
		}
		debo, err := conn.DescribeEventBus(input)
		if err != nil {
			return fmt.Errorf("Reading CloudWatch Events bus policy for '%s' failed: %w", pr, err)
		}

		if debo.Policy == nil {
			return fmt.Errorf("Not found: %s", pr)
		}

		var policyDoc tfcloudwatchevents.PermissionPolicyDoc
		err = json.Unmarshal([]byte(*debo.Policy), &policyDoc)
		if err != nil {
			return fmt.Errorf("Reading CloudWatch Events bus policy for '%s' failed: %w", pr, err)
		}

		_, err = tfcloudwatchevents.FindPermissionPolicyStatementByID(&policyDoc, statementID)
		return err
	}
}

func testAccCheckCloudWatchEventPermissionDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).CloudWatchEventsConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cloudwatch_event_permission" {
			continue
		}

		eventBusName, statementID, err := tfcloudwatchevents.PermissionParseResourceID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error reading CloudWatch Events permission (%s): %w", rs.Primary.ID, err)
		}
		input := &events.DescribeEventBusInput{
			Name: aws.String(eventBusName),
		}
		err = resource.Retry(1*time.Minute, func() *resource.RetryError {
			debo, err := conn.DescribeEventBus(input)
			if tfawserr.ErrMessageContains(err, events.ErrCodeResourceNotFoundException, "") {
				return nil
			}
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if debo.Policy == nil {
				return nil
			}

			var policyDoc tfcloudwatchevents.PermissionPolicyDoc
			err = json.Unmarshal([]byte(*debo.Policy), &policyDoc)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("Reading CloudWatch Events permission '%s' failed: %w", rs.Primary.ID, err))
			}

			_, err = tfcloudwatchevents.FindPermissionPolicyStatementByID(&policyDoc, statementID)
			if err == nil {
				return resource.RetryableError(fmt.Errorf("CloudWatch Events permission exists: %s", rs.Primary.ID))
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func testAccCheckPermissionResourceBasicConfig(principal, statementID string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_event_permission" "test" {
  principal    = "%[1]s"
  statement_id = "%[2]s"
}
`, principal, statementID)
}

func testAccCheckPermissionResourceDefaultEventBusNameConfig(principal, statementID string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_event_permission" "test" {
  principal      = %[1]q
  statement_id   = %[2]q
  event_bus_name = "default"
}
`, principal, statementID)
}

func testAccCheckPermissionResourceEventBusNameConfig(principal, busName, statementID string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_event_permission" "test" {
  principal      = %[1]q
  statement_id   = %[2]q
  event_bus_name = aws_cloudwatch_event_bus.test.name
}

resource "aws_cloudwatch_event_bus" "test" {
  name = %[3]q
}
`, principal, statementID, busName)
}

func testAccCheckPermissionResourceActionConfig(action, principal, statementID string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_event_permission" "test" {
  action       = "%[1]s"
  principal    = "%[2]s"
  statement_id = "%[3]s"
}
`, action, principal, statementID)
}

func testAccCheckPermissionResourceConditionOrganizationConfig(statementID, value string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_event_permission" "test" {
  principal    = "*"
  statement_id = %q

  condition {
    key   = "aws:PrincipalOrgID"
    type  = "StringEquals"
    value = %q
  }
}
`, statementID, value)
}

func testAccCheckPermissionResourceMultipleConfig(principal1, statementID1, principal2, statementID2 string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_event_permission" "test" {
  principal    = "%[1]s"
  statement_id = "%[2]s"
}

resource "aws_cloudwatch_event_permission" "test2" {
  principal    = "%[3]s"
  statement_id = "%[4]s"
}
`, principal1, statementID1, principal2, statementID2)
}
