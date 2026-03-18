#!/bin/sh
# Compiler wrapper for controlling SSE mode in go-imagequant compilation
# Controlled by NUGGAN_USE_SSE environment variable:
#   unset or any other value -> SSE enabled on x86_64, auto-disabled on ARM64 (default)
#   0, false, no, off         -> SSE explicitly disabled
#   1, true, yes, on          -> SSE explicitly enabled (may fail on ARM64)

# Detect target architecture
target_arch=""
for arg in "$@"; do
  case "$arg" in
    *arm64*apple*darwin*)
      target_arch="arm64"
      break
      ;;
    *x86_64*apple*darwin*)
      target_arch="x86_64"
      break
      ;;
  esac
done

# Determine if SSE should be disabled
disable_sse=0
case "${NUGGAN_USE_SSE}" in
  0|false|no|off)
    # Explicitly disable SSE
    disable_sse=1
    ;;
  1|true|yes|on)
    # Explicitly enable SSE (even on ARM64, may fail)
    disable_sse=0
    ;;
  *)
    # Default: auto-disable SSE on ARM64, enable on other architectures
    if [ "$target_arch" = "arm64" ]; then
      disable_sse=1
    else
      disable_sse=0
    fi
    ;;
esac

# Find the real C compiler
if command -v clang >/dev/null 2>&1; then
  real_cc="clang"
elif command -v cc >/dev/null 2>&1; then
  real_cc="cc"
elif [ -x /usr/bin/clang ]; then
  real_cc="/usr/bin/clang"
elif [ -x /usr/bin/cc ]; then
  real_cc="/usr/bin/cc"
else
  real_cc="cc"
fi

# Process arguments if SSE should be disabled
if [ "$disable_sse" = "1" ]; then
  out=""
  for a in "$@"; do
    case "$a" in
      -msse)
        # Skip -msse flag
        ;;
      -DUSE_SSE=1)
        # Replace -DUSE_SSE=1 with -DUSE_SSE=0
        out="$out -DUSE_SSE=0"
        ;;
      *)
        out="$out '$a'"
        ;;
    esac
  done
  eval exec "$real_cc" $out
else
  # SSE enabled: pass through all arguments unchanged
  exec "$real_cc" "$@"
fi


