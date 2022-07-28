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
