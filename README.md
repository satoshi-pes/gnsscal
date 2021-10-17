# gnsscal
displays a gnss calendar

gnsscal - displays a GNSS calendar inspired by 'gpscal'

# Overview
The gnsscal displays a calendar similar to 'cal' command except for displaying gnss week and doy.

GNSS data analysis often needs calendar with GNSS weeks, doy that used in the names of related files, pre-processing settings, etc. This command privides simple but useful cal-command-like features with the similar interface.

For default, gnsscal displays only the current month. 
If month or year is given, print the specified month / year. In the case only the year is specified, a gnss calender for one year period is displayed.

The gnsscal command is developed inspired by the 'gpscal' created by Dr. Yuki Hatanaka at Geospatial Information Authority of Japan.

# Usage
    gnsscal [Flags] [[month] year]
    
    Flags:
      -h        help for gnsscal
      -n        turns off highlight of today [default: highlight on]
      -3        three-month layout that displays previous, current and next months
      -satsys   satellite system of GNSS week; 'GPS', 'QZS', 'GAL', 'BDS', or 'GLO' [default: GPS]
    
# Example
In default, gnsscal displays current month in a following layout:

    $ gnsscal
    GPS           October 2021
    Week   Sun Mon Tue Wed Thu Fri Sat
    2177                         1   2
                               274 275
    2178     3   4   5   6   7   8   9
           276 277 278 279 280 281 282
    2179    10  11  12  13  14  15  16
           283 284 285 286 287 288 289
    2180    17  18  19  20  21  22  23
           290 291 292 293 294 295 296
    2181    24  25  26  27  28  29  30
           297 298 299 300 301 302 303
    2182    31
           304

The system time can be changed specifing '-satsys' flag. Other layout, e.g. three-month layout, can be also invoked with '-3' flag: 

    $ gnsscal -satsys "GAL" -3
    GAL          September 2021           GAL           October 2021            GAL          November 2021
    Week   Sun Mon Tue Wed Thu Fri Sat    Week   Sun Mon Tue Wed Thu Fri Sat    Week   Sun Mon Tue Wed Thu Fri Sat
    1149                 1   2   3   4    1153                         1   2    1158         1   2   3   4   5   6
                       244 245 246 247                               274 275               305 306 307 308 309 310
    1150     5   6   7   8   9  10  11    1154     3   4   5   6   7   8   9    1159     7   8   9  10  11  12  13
           248 249 250 251 252 253 254           276 277 278 279 280 281 282           311 312 313 314 315 316 317
    1151    12  13  14  15  16  17  18    1155    10  11  12  13  14  15  16    1160    14  15  16  17  18  19  20
           255 256 257 258 259 260 261           283 284 285 286 287 288 289           318 319 320 321 322 323 324
    1152    19  20  21  22  23  24  25    1156    17  18  19  20  21  22  23    1161    21  22  23  24  25  26  27
           262 263 264 265 266 267 268           290 291 292 293 294 295 296           325 326 327 328 329 330 331
    1153    26  27  28  29  30            1157    24  25  26  27  28  29  30    1162    28  29  30
           269 270 271 272 273                   297 298 299 300 301 302 303           332 333 334
                                          1158    31
                                                 304

# Installing
It's easy to install gnsscal. Simply get and install the package:

    go get -u github.com/satoshi-pes/gnsscal  
    go install github.com/satoshi-pes/gnsscal@latest

# License
gnsscal is released under the MIT license. See [LICENSE.txt](https://github.com/satoshi-pes/gnsscal/blob/master/LICENSE.txt)