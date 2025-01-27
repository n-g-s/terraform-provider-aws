package cloudfront_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudfront"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfcloudfront "github.com/hashicorp/terraform-provider-aws/internal/service/cloudfront"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccAWSCloudFrontResponseHeadersPolicy_CorsConfig(t *testing.T) {
	rName1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudfront_response_headers_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(cloudfront.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, cloudfront.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckCloudFrontResponseHeadersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCloudFrontResponseHeadersPolicyCorsConfigConfig(rName1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudFrontResponseHeadersPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", "test comment"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_credentials", "false"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_headers.0.items.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_config.0.access_control_allow_headers.0.items.*", "X-Header1"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_methods.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_methods.0.items.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_config.0.access_control_allow_methods.0.items.*", "GET"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_config.0.access_control_allow_methods.0.items.*", "POST"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_origins.0.items.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_config.0.access_control_allow_origins.0.items.*", "test1.example.com"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_config.0.access_control_allow_origins.0.items.*", "test2.example.com"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_expose_headers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_max_age_sec", "0"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.origin_override", "true"),
					resource.TestCheckResourceAttr(resourceName, "custom_headers_config.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "etag"),
					resource.TestCheckResourceAttr(resourceName, "name", rName1),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccAWSCloudFrontResponseHeadersPolicyCorsConfigUpdatedConfig(rName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudFrontResponseHeadersPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", "test comment updated"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_credentials", "true"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_headers.0.items.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_config.0.access_control_allow_headers.0.items.*", "X-Header2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_config.0.access_control_allow_headers.0.items.*", "X-Header3"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_methods.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_config.0.access_control_allow_methods.0.items.*", "PUT"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_allow_origins.0.items.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_config.0.access_control_allow_origins.0.items.*", "test1.example.com"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_config.0.access_control_allow_origins.0.items.*", "test2.example.com"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_expose_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_expose_headers.0.items.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_config.0.access_control_expose_headers.0.items.*", "HEAD"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.access_control_max_age_sec", "3600"),
					resource.TestCheckResourceAttr(resourceName, "cors_config.0.origin_override", "false"),
					resource.TestCheckResourceAttr(resourceName, "custom_headers_config.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "etag"),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.#", "0"),
				),
			},
		},
	})
}

func TestAccAWSCloudFrontResponseHeadersPolicy_CustomHeadersConfig(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudfront_response_headers_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(cloudfront.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, cloudfront.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckCloudFrontResponseHeadersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCloudFrontResponseHeadersPolicyCustomHeadersConfigConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudFrontResponseHeadersPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", ""),
					resource.TestCheckResourceAttr(resourceName, "cors_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_headers_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_headers_config.0.items.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "custom_headers_config.0.items.*", map[string]string{
						"header":   "X-Header1",
						"override": "true",
						"value":    "value1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "custom_headers_config.0.items.*", map[string]string{
						"header":   "X-Header2",
						"override": "false",
						"value":    "value2",
					}),
					resource.TestCheckResourceAttrSet(resourceName, "etag"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAWSCloudFrontResponseHeadersPolicy_SecurityHeadersConfig(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudfront_response_headers_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(cloudfront.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, cloudfront.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckCloudFrontResponseHeadersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCloudFrontResponseHeadersPolicySecurityHeadersConfigConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudFrontResponseHeadersPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", ""),
					resource.TestCheckResourceAttr(resourceName, "cors_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_headers_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_headers_config.0.items.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "custom_headers_config.0.items.*", map[string]string{
						"header":   "X-Header1",
						"override": "true",
						"value":    "value1",
					}),
					resource.TestCheckResourceAttrSet(resourceName, "etag"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.content_security_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.content_security_policy.0.content_security_policy", "policy1"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.content_security_policy.0.override", "true"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.content_type_options.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.frame_options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.frame_options.0.frame_option", "DENY"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.frame_options.0.override", "false"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.referrer_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.strict_transport_security.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.strict_transport_security.0.access_control_max_age_sec", "90"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.strict_transport_security.0.include_subdomains", "false"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.strict_transport_security.0.override", "true"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.strict_transport_security.0.preload", "true"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.xss_protection.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccAWSCloudFrontResponseHeadersPolicySecurityHeadersConfigUpdatedConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudFrontResponseHeadersPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "comment", ""),
					resource.TestCheckResourceAttr(resourceName, "cors_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "custom_headers_config.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "etag"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.content_security_policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.content_type_options.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.content_type_options.0.override", "true"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.frame_options.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.referrer_policy.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.referrer_policy.0.override", "false"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.referrer_policy.0.referrer_policy", "origin-when-cross-origin"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.strict_transport_security.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.xss_protection.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.xss_protection.0.mode_block", "false"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.xss_protection.0.override", "true"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.xss_protection.0.protection", "true"),
					resource.TestCheckResourceAttr(resourceName, "security_headers_config.0.xss_protection.0.report_uri", "https://example.com/"),
				),
			},
		},
	})
}

func TestAccAWSCloudFrontResponseHeadersPolicy_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_cloudfront_response_headers_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(cloudfront.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, cloudfront.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckCloudFrontResponseHeadersPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSCloudFrontResponseHeadersPolicyCorsConfigConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudFrontResponseHeadersPolicyExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfcloudfront.ResourceResponseHeadersPolicy(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckCloudFrontResponseHeadersPolicyDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).CloudFrontConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_cloudfront_response_headers_policy" {
			continue
		}

		_, err := tfcloudfront.FindResponseHeadersPolicyByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("CloudFront Response Headers Policy %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckCloudFrontResponseHeadersPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CloudFront Response Headers Policy ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).CloudFrontConn

		_, err := tfcloudfront.FindResponseHeadersPolicyByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccAWSCloudFrontResponseHeadersPolicyCorsConfigConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudfront_response_headers_policy" "test" {
  name    = %[1]q
  comment = "test comment"

  cors_config {
    access_control_allow_credentials = false

    access_control_allow_headers {
      items = ["X-Header1"]
    }

    access_control_allow_methods {
      items = ["GET", "POST"]
    }

    access_control_allow_origins {
      items = ["test1.example.com", "test2.example.com"]
    }

    origin_override = true
  }
}
`, rName)
}

func testAccAWSCloudFrontResponseHeadersPolicyCorsConfigUpdatedConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudfront_response_headers_policy" "test" {
  name    = %[1]q
  comment = "test comment updated"

  cors_config {
    access_control_allow_credentials = true

    access_control_allow_headers {
      items = ["X-Header2", "X-Header3"]
    }

    access_control_allow_methods {
      items = ["PUT"]
    }

    access_control_allow_origins {
      items = ["test1.example.com", "test2.example.com"]
    }

    access_control_expose_headers {
      items = ["HEAD"]
    }

    access_control_max_age_sec = 3600

    origin_override = false
  }
}
`, rName)
}

func testAccAWSCloudFrontResponseHeadersPolicyCustomHeadersConfigConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudfront_response_headers_policy" "test" {
  name = %[1]q

  custom_headers_config {
    items {
      header   = "X-Header2"
      override = false
      value    = "value2"
    }

    items {
      header   = "X-Header1"
      override = true
      value    = "value1"
    }
  }
}
`, rName)
}

func testAccAWSCloudFrontResponseHeadersPolicySecurityHeadersConfigConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudfront_response_headers_policy" "test" {
  name = %[1]q

  custom_headers_config {
    items {
      header   = "X-Header1"
      override = true
      value    = "value1"
    }
  }

  security_headers_config {
    content_security_policy {
      content_security_policy = "policy1"
      override                = true
    }

    frame_options {
      frame_option = "DENY"
      override     = false
    }

    strict_transport_security {
      access_control_max_age_sec = 90
      override                   = true
      preload                    = true
    }
  }
}
`, rName)
}

func testAccAWSCloudFrontResponseHeadersPolicySecurityHeadersConfigUpdatedConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_cloudfront_response_headers_policy" "test" {
  name = %[1]q

  security_headers_config {
    content_type_options {
      override = true
    }

    referrer_policy {
      referrer_policy = "origin-when-cross-origin"
      override        = false
    }

    xss_protection {
      override   = true
      protection = true
      report_uri = "https://example.com/"
    }
  }
}
`, rName)
}
