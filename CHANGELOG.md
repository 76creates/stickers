# Changelog

All notable changes to this project will be documented in this file.

## [v1.3] (draft)
### Fixes
* Fixed cursor not moving to the last visible row when filtering
* Fixed margins and borders not being rendered correctly #8
* Additionally fixed margin and border issues on box
### Features
* Converted default cell inheritance of the row style to function `StylePassing` which can be set on the _Box_ and _Row_, if both box and row have style passing enabled, the row will inherit the box style before it passes style to the cells.

## [v1.2](https://github.com/76creates/stickers/compare/v1.1...v1.2) (2022-02-27)
### Features
* Filterin is now availible for `Table` and `TableSingleType` using new methods:
    * `UnsetFilter` remove filtering
    * `SetFilter` sets the filter on a column by index
    * `GetFilter` gets index of filtered column and the value of the filter
* Added `MustGetCellWithIndex`
* Fixed visible table calculations when filtering
* Added filter info to the status box
* Header rendering of sorting and filtering symbols is improved

## [v1.1](https://github.com/76creates/stickers/compare/v1.0...v1.1) (2022-02-26)
### ⚠ BREAKING CHANGES
* Refactored `Table` to support sorting, some methods have changed most notably revolving around adding rows since now its taking [][]any instead of [][]string, initial `Table` is now closer to `TableSingleType[string]`
* Stickers now uses generics, so go1.18 is mandatory

### Fixes
* Fixed recalculation triggering when *FlexBoxCell or *FlexBoxRow is fetched from the FlexBox
* Small lexical changes and tidying up

### Features
* Sorting is now availible for `Table` and `TableSingleType`
* `Table` has been reformatted and now supports **sorting by type**, when `Table` is initialized each colum type is set to `string`, you can now update that using `SetType` method, types supported are located in interface `Ordered`
* Added `TableSingleType` type which locks row type to `string`, this makes it easier for user when adding rows as there is much less errors that can occurr as when using `Table` where all depends on a type
* Added method `OrderByColumn` which envokes sorting for column `n`, for now you cannot explicitly set sorting direction and its switching between `asc` and `desc` when you sort same column 
* Added method `GetCursorLocation` which returns `x`,`y` of the curent cursor location on the table
* Added error types `TableBadTypeError`, `TableRowLenError`, `TableBadCellTypeError`
* Minor preformace enhancements