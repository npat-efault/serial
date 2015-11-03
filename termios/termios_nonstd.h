/* 
 * termios_nonstd.h
 *
 * Non-standard termios mode-flags and cotrol-character indexes. 
 *
 * Copyright (c) 2015, Nick Patavalis (npat@efault.net).
 * All rights reserved.
 * Use of this source code is governed by a BSD-style license that can
 * be found in the LICENSE file.
 */

/* Non-standard flags, set to 0 if missing */

#ifndef IMAXBEL
#define IMAXBEL 0
#endif
#ifndef IUCLC
#define IUCLC 0
#endif
#ifndef IUTF8
#define IUTF8 0
#endif

#ifndef CRTSCTS
#define CRTSCTS 0
#endif
#ifndef CMSPAR
#define CMSPAR 0
#endif

#ifndef PENDIN
#define PENDIN 0
#endif
#ifndef ECHOCTL
#define ECHOCTL 0
#endif
#ifndef ECHOPRT
#define ECHOPRT 0
#endif
#ifndef ECHOKE
#define ECHOKE 0
#endif
#ifndef FLUSHO
#define FLUSHO 0
#endif
#ifndef EXTPROC
#define EXTPROC 0
#endif

/* Non-standard cc indexes. Set to -1 if missing */

#ifndef VREPRINT
#define VREPRINT -1
#endif
#ifndef VDISCARD
#define VDISCARD -1
#endif
#ifndef VWERASE
#define VWERASE -1
#endif
#ifndef VLNEXT
#define VLNEXT -1
#endif
#ifndef VEOL2
#define VEOL2 -1
#endif
