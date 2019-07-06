# GoHeader

A utility to add source code file header.

## Usages

1. Adds the `.header.json` configuration file under your project root directory

  ```json
   {
    "Dir": ".",
    "Template": ".header.txt",
    "Adapter": [
      ".go",
	  ".ex"
	],
    "Excludes": [
      "./test",
	  "./default_excludes.go",
      "./header.go",
      "./main.go",
      "./1.txt"
    ],
    "UseDefaultExcludes": true,
    "Properties": {
        "Year": "2006-2019",
        "Owner": "xiaobo"
    }
}
  ```
2. Adds the `.header.txt` template

  ```
  Copyright (c) {{.Year}}, {{.Owner}}

  This is free software, licensed under the GNU General Public License v3.
  See /LICENSE for more information.

  ```
3. Executes the `goheader` binary
