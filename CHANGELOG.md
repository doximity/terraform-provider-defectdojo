## 0.0.12

BUGFIX:
 - A product with no tags specified would cause a provider error from terraform.

## 0.0.11

FEATURES:
 - Add the following fields to `defectdojo_product` resource:
   - `business_criticality`
   - `enable_full_risk_acceptance`
   - `enable_skip_risk_acceptance`
   - `external_audience`
   - `internet_accessible`
   - `lifecycle`
   - `origin`
   - `platform`
   - `prod_numeric_grade`
   - `regulation_ids`
   - `revenue`
   - `user_records`

## 0.0.10

FEATURES:
 - Add `jira_product_configuration` resource.

## 0.0.9

BUGFIX:
 - Fix delete-drift detection in `product` and `product_type` resources. If the object was deleted outside terraform we remove it from the state.

## 0.0.8

BUGFIX:
 - Don't continue processing after encountering an error, that cause a panic.

## 0.0.7

Initial public release

## 0.0.4

FEATURES:
 - Add basic support for Product Type resource and data source

## 0.0.3

FEATURES:
 - First working version.
 - Basic support for Product resource and data source.
