package main

#ifdef test1
const posidCount = 2000
#endif

#ifdef ipadic
const posidCount = 2000

const unigramUnknownLeftPosid = 1
const unigramUnknownRightPosid = 1
const unigramUnknownWordCost = 10000

const maxConnCost = 25400
#endif

#ifdef unidic
const posidCount = 6000

const unigramUnknownLeftPosid = 5978
const unigramUnknownRightPosid = 5978
const unigramUnknownWordCost = 10000

const maxConnCost = 25400
#endif

