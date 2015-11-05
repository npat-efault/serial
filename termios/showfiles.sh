#!/bin/sh
#
# showfiles.sh
#
# Show which files will be used to build the package.
# Usage:
#
#   ./showfiles [tag ...]
#
# Copyright (c) 2015, Nick Patavalis (npat@efault.net).
# All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE.txt file.

go list -tags="$@" -f '{{range .GoFiles}}{{println .}}{{end}}{{range .CgoFiles}}{{println .}}{{end}}{{range .SFiles}}{{println .}}{{end}}' .
