# Changelog

All notable changes to this project will be documented in this file.
## [v1.4.2](https://github.com/76creates/stickers/compare/v1.4.1...v1.4.2) (2024-11-26)
### Fixes
- Fix typo in flexbox/Cell.SetMinHeight
### Deprecations
- `Cell.SetMinHeigth` is now deprecated in favour of SetMinHeight.

## [v1.4.1](https://github.com/76creates/stickers/compare/v1.4.0...v1.4.1) (2024-10-21)
### Features
- Added `OrderByAsc` and `OrderByDesc` methods to `Table` #12 @drmille2
### Deprecations
- `OrderByColumn` is deprecated by `OrderByAsc` and `OrderByDesc` methods. #12 @drmille2
### Updates
- Updated Go to version `1.23`
### Dependencies
- Updated `github.com/charmbracelet/lipgloss` to `v0.13.0`
- Updated `github.com/charmbracelet/bubbletea` to `v1.1.1`
- Updated `github.com/gocarina/gocsv` to `78e41c74b4b1`

## [v1.4.0](https://github.com/76creates/stickers/compare/v1.3.0...v1.4.0) (2024-09-11)
### ⚠ BREAKING CHANGES
- Moved `flexbox` and `table` into separate packages, `github.com/76creates/stickers/flexbox` and `github.com/76creates/stickers/table` respectively. #10 @jon4hz
### Fixes
- Minor lexical fixes
- Fixed repo tags to match go semver format.
### Dependencies
- Updated `github.com/charmbracelet/lipgloss` to `v0.6.0'
### Features
- Added `SetStylePassing` to _Table_ that will pass down the style all the way, from box to cell. No granularity for now.
- Added `HorizontalFlexBox`. #10 @jon4hz
### Updates
- Refactored `FlexBox.GetRow`, `FlexBox.Row`, `FlexBox.MustGetRow`, `FlexBoxRow.Cell`, `FlexBoxRow.GetCellWithID`, `FlexBoxRow.MustGetCellWithIndex`.<br>They are replaced with `FlexBoxRow.GetCell`, `FlexBoxRow.GetCellCopy`, `FlexBox.GetRow`, `FlexBox.GetRowCopy`,`FlexBox.GetRowCellCopy`.<br>Get* now returns pointer and triggers _recalculation_, while one can use Copy* function to get pointer to copied structs which can be used to lookup values without triggering _recalculation_.
- `AddCells` now take cells as a variadic argument. #10 @jon4hz

## [v1.3.0](https://github.com/76creates/stickers/compare/v1.2.0...v1.3.0) (2022-12-28)
### Fixes
* Fixed cursor not moving to the last visible row when filtering
* Fixed margins and borders not being rendered correctly #8
* Additionally, fixed margin and border issues on box
* Allow style override on _Table_ @joejag
### Features
* Converted default cell inheritance of the row style to function `StylePassing` which can be set on the _Box_ and _Row_, if both box and row have style passing enabled, the row will inherit the box style before it passes style to the cells.

## [v1.2.0](https://github.com/76creates/stickers/compare/v1.1.0...v1.2.0) (2022-02-27)
### Features
* Filtering is now available for `Table` and `TableSingleType` using new methods:
    * `UnsetFilter` remove filtering
    * `SetFilter` sets the filter on a column by index
    * `GetFilter` gets index of filtered column and the value of the filter
* Added `MustGetCellWithIndex`
* Fixed visible table calculations when filtering
* Added filter info to the status box
* Header rendering of sorting and filtering symbols is improved

## [v1.1.0](https://github.com/76creates/stickers/compare/v1.0.0...v1.1.0) (2022-02-26)
### ⚠ BREAKING CHANGES
* Refactored `Table` to support sorting, some methods have changed most notably revolving around adding rows since now its taking [][]any instead of [][]string, initial `Table` is now closer to `TableSingleType[string]`
* Stickers now uses generics, so go1.18 is mandatory

### Fixes
* Fixed recalculation triggering when *FlexBoxCell or *FlexBoxRow is fetched from the FlexBox
* Small lexical changes and tidying up

### Features
* Sorting is now available for `Table` and `TableSingleType`
* `Table` has been reformatted and now supports **sorting by type**, when `Table` is initialized each colum type is set to `string`, you can now update that using `SetType` method, types supported are located in interface `Ordered`
* Added `TableSingleType` type which locks row type to `string`, this makes it easier for user when adding rows as there is much fewer errors that can occur as when using `Table` where all depends on a type
* Added method `OrderByColumn` which invokes sorting for column `n`, for now you cannot explicitly set sorting direction and it's switching between `asc` and `desc` when you sort same column 
* Added method `GetCursorLocation` which returns `x`,`y` of the current cursor location on the table
* Added error types `TableBadTypeError`, `TableRowLenError`, `TableBadCellTypeError`
* Minor performance enhancements
