
#ifndef INCLUDE_H
#define INCLUDE_H

#ifdef __cplusplus
extern "C" {
#endif

#include "tree_sitter/parser.h"
#ifdef TREE_SITTER_HIDE_SYMBOLS
#define TS_PUBLIC
#elif defined(_WIN32)
#define TS_PUBLIC __declspec(dllexport)
#else
#define TS_PUBLIC __attribute__((visibility("default")))
#endif
TS_PUBLIC const TSLanguage *tree_sitter_kamailio_cfg(void);

#ifdef __cplusplus
}
#endif

#endif /* INCLUDE_H */
